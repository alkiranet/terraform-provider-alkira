package alkira

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-cty/cty"
	ctymsgpack "github.com/hashicorp/go-cty/cty/msgpack"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	azureVnetTypeName = "alkira_connector_azure_vnet"
	cpgAttr           = "peering_gateway_cxp_id"
)

// azureVnetProvider drives the resource through the same gRPC calls Terraform
// core makes, so plan results include DiffSuppressFunc and Computed handling.
type azureVnetProvider struct {
	res    *schema.Resource
	server tfprotov5.ProviderServer
	ty     cty.Type
	blocks map[string]bool
}

func newAzureVnetProvider(client *alkira.AlkiraClient) *azureVnetProvider {
	res := resourceAlkiraConnectorAzureVnet()
	p := &schema.Provider{ResourcesMap: map[string]*schema.Resource{azureVnetTypeName: res}}
	p.SetMeta(client)

	coreSchema := res.CoreConfigSchema()
	blocks := map[string]bool{}
	for name := range coreSchema.BlockTypes {
		blocks[name] = true
	}

	return &azureVnetProvider{
		res:    res,
		server: schema.NewGRPCProviderServer(p),
		ty:     coreSchema.ImpliedType(),
		blocks: blocks,
	}
}

// config returns a minimal configuration. A null cpg means the attribute is
// omitted.
func (p *azureVnetProvider) config(cpg cty.Value) cty.Value {
	m := map[string]cty.Value{}
	for name, attrType := range p.ty.AttributeTypes() {
		switch {
		case p.blocks[name] && attrType.IsListType():
			m[name] = cty.ListValEmpty(attrType.ElementType())
		case p.blocks[name]:
			m[name] = cty.SetValEmpty(attrType.ElementType())
		default:
			m[name] = cty.NullVal(attrType)
		}
	}
	m["name"] = cty.StringVal("vnet")
	m["azure_vnet_id"] = cty.StringVal("vnet-id")
	m["credential_id"] = cty.StringVal("cred")
	m["cxp"] = cty.StringVal("USCENTRAL-AZURE-3")
	m["segment_id"] = cty.StringVal("5")
	m["size"] = cty.StringVal("SMALL")
	m["connection_mode"] = cty.StringVal("VNET_PEERING")
	m[cpgAttr] = cpg
	return cty.ObjectVal(m)
}

func (p *azureVnetProvider) state(cpg cty.Value) cty.Value {
	m := p.config(cpg).AsValueMap()
	m["id"] = cty.StringVal("1")
	return cty.ObjectVal(m)
}

// proposedNew mirrors core: a Computed attribute that config leaves null takes
// the prior value.
func (p *azureVnetProvider) proposedNew(config, prior cty.Value) cty.Value {
	m := config.AsValueMap()
	for name, s := range p.res.Schema {
		if s.Computed && m[name].IsNull() && !prior.IsNull() {
			m[name] = prior.GetAttr(name)
		}
	}
	m["id"] = cty.NullVal(cty.String)
	if !prior.IsNull() {
		m["id"] = prior.GetAttr("id")
	}
	return cty.ObjectVal(m)
}

func (p *azureVnetProvider) encode(t *testing.T, v cty.Value) *tfprotov5.DynamicValue {
	b, err := ctymsgpack.Marshal(v, p.ty)
	require.NoError(t, err)
	return &tfprotov5.DynamicValue{MsgPack: b}
}

func (p *azureVnetProvider) decode(t *testing.T, v *tfprotov5.DynamicValue) cty.Value {
	out, err := ctymsgpack.Unmarshal(v.MsgPack, p.ty)
	require.NoError(t, err)
	return out
}

// planCpg returns the planned peering_gateway_cxp_id.
func (p *azureVnetProvider) planCpg(t *testing.T, prior, config cty.Value) string {
	resp, err := p.server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
		TypeName:         azureVnetTypeName,
		PriorState:       p.encode(t, prior),
		ProposedNewState: p.encode(t, p.proposedNew(config, prior)),
		Config:           p.encode(t, config),
	})
	require.NoError(t, err)
	require.Empty(t, resp.Diagnostics)
	return cpgString(p.decode(t, resp.PlannedState).GetAttr(cpgAttr))
}

func cpgString(v cty.Value) string {
	switch {
	case !v.IsKnown():
		return "unknown"
	case v.IsNull():
		return "null"
	default:
		return v.AsBigFloat().String()
	}
}

// azureVnetConnectorHandler serves the connector GET and the segment lookup
// Read makes. A nil cpg leaves cxpPeeringGatewayId out of the response.
func azureVnetConnectorHandler(cpg *int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasSuffix(r.URL.Path, "/segments"):
			w.Write([]byte(`[{"id":5,"name":"seg"}]`))
		case strings.Contains(r.URL.Path, "/segments/"):
			w.Write([]byte(`{"id":5,"name":"seg"}`))
		default:
			connector := map[string]interface{}{
				"id":             1,
				"name":           "vnet",
				"cxp":            "USCENTRAL-AZURE-3",
				"credentialId":   "cred",
				"vnetId":         "vnet-id",
				"size":           "SMALL",
				"segments":       []string{"seg"},
				"connectionMode": "VNET_PEERING",
				"enabled":        true,
				"vnetRouting": map[string]interface{}{
					"importFromCXPOptions": map[string]interface{}{"routeImportMode": "ADV_ONLY"},
				},
			}
			if cpg != nil {
				connector["cxpPeeringGatewayId"] = *cpg
			}
			body, _ := json.Marshal(connector)
			w.Write(body)
		}
	}
}

// An explicit peering_gateway_cxp_id = 0 (as written by terraform import
// -generate-config-out) plans the same as an omitted value.
func TestPeeringGatewayCxpIdZeroPlansAsUnset(t *testing.T) {
	p := newAzureVnetProvider(&alkira.AlkiraClient{})
	zero := cty.NumberIntVal(0)
	omitted := cty.NullVal(cty.Number)
	prior := p.state(cty.NumberIntVal(163))

	t.Run("zero keeps the assigned gateway", func(t *testing.T) {
		assert.Equal(t, "163", p.planCpg(t, prior, p.config(zero)))
	})

	t.Run("omitted keeps the assigned gateway", func(t *testing.T) {
		assert.Equal(t, "163", p.planCpg(t, prior, p.config(omitted)))
	})

	t.Run("explicit id still plans a change", func(t *testing.T) {
		assert.Equal(t, "200", p.planCpg(t, prior, p.config(cty.NumberIntVal(200))))
	})

	t.Run("zero on create is known after apply", func(t *testing.T) {
		assert.Equal(t, "unknown", p.planCpg(t, cty.NullVal(p.ty), p.config(zero)))
	})
}

// Read leaves peering_gateway_cxp_id null when the API has no gateway, so an
// imported connector does not generate peering_gateway_cxp_id = 0.
func TestPeeringGatewayCxpIdImportWithoutGatewayIsNull(t *testing.T) {
	assigned := 163
	cases := []struct {
		name string
		cpg  *int
		want string
	}{
		{"no gateway", nil, "null"},
		{"assigned gateway", &assigned, "163"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := newAzureVnetProvider(createMockAlkiraClient(t, azureVnetConnectorHandler(c.cpg)))

			imported, err := p.server.ImportResourceState(context.Background(), &tfprotov5.ImportResourceStateRequest{
				TypeName: azureVnetTypeName,
				ID:       "1",
			})
			require.NoError(t, err)
			require.Empty(t, imported.Diagnostics)
			require.Len(t, imported.ImportedResources, 1)

			// Terraform refreshes an imported resource before writing its state.
			refreshed, err := p.server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
				TypeName:     azureVnetTypeName,
				CurrentState: imported.ImportedResources[0].State,
			})
			require.NoError(t, err)
			require.Empty(t, refreshed.Diagnostics)

			assert.Equal(t, c.want, cpgString(p.decode(t, refreshed.NewState).GetAttr(cpgAttr)))
		})
	}
}

// apply plans config against prior, applies the plan and returns the new state.
func (p *azureVnetProvider) apply(t *testing.T, prior, config cty.Value) cty.Value {
	planned, err := p.server.PlanResourceChange(context.Background(), &tfprotov5.PlanResourceChangeRequest{
		TypeName:         azureVnetTypeName,
		PriorState:       p.encode(t, prior),
		ProposedNewState: p.encode(t, p.proposedNew(config, prior)),
		Config:           p.encode(t, config),
	})
	require.NoError(t, err)
	require.Empty(t, planned.Diagnostics)

	applied, err := p.server.ApplyResourceChange(context.Background(), &tfprotov5.ApplyResourceChangeRequest{
		TypeName:       azureVnetTypeName,
		PriorState:     p.encode(t, prior),
		PlannedState:   planned.PlannedState,
		Config:         p.encode(t, config),
		PlannedPrivate: planned.PlannedPrivate,
	})
	require.NoError(t, err)
	require.Empty(t, applied.Diagnostics)
	return p.decode(t, applied.NewState)
}

// requestRecorder answers PUT and POST with a connector id and keeps the last
// request body. Other requests go to next.
type requestRecorder struct {
	body map[string]interface{}
}

func (rec *requestRecorder) handler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			next(w, r)
			return
		}
		rec.body = map[string]interface{}{}
		json.NewDecoder(r.Body).Decode(&rec.body)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"id":1}`))
	}
}

// An update with peering_gateway_cxp_id omitted or 0 sends the assigned
// gateway, and an explicit id sends that id.
func TestPeeringGatewayCxpIdUpdateSendsGateway(t *testing.T) {
	cases := []struct {
		name string
		cpg  cty.Value
		want int
	}{
		{"zero", cty.NumberIntVal(0), 163},
		{"omitted", cty.NullVal(cty.Number), 163},
		{"explicit id", cty.NumberIntVal(200), 200},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assigned := 163
			rec := &requestRecorder{}
			p := newAzureVnetProvider(createMockAlkiraClient(t, rec.handler(azureVnetConnectorHandler(&assigned))))

			config := p.config(c.cpg).AsValueMap()
			config["description"] = cty.StringVal("updated")
			p.apply(t, p.state(cty.NumberIntVal(163)), cty.ObjectVal(config))

			require.NotNil(t, rec.body)
			assert.EqualValues(t, c.want, rec.body["cxpPeeringGatewayId"])
		})
	}
}

// A create with peering_gateway_cxp_id omitted or 0 leaves cxpPeeringGatewayId
// out of the request, and the gateway the API assigns lands in state.
func TestPeeringGatewayCxpIdCreateOmitsUnsetGateway(t *testing.T) {
	cases := []struct {
		name string
		cpg  cty.Value
	}{
		{"omitted", cty.NullVal(cty.Number)},
		{"zero", cty.NumberIntVal(0)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assigned := 163
			rec := &requestRecorder{}
			p := newAzureVnetProvider(createMockAlkiraClient(t, rec.handler(azureVnetConnectorHandler(&assigned))))

			state := p.apply(t, cty.NullVal(p.ty), p.config(c.cpg))

			require.NotNil(t, rec.body)
			assert.NotContains(t, rec.body, "cxpPeeringGatewayId")
			assert.Equal(t, "163", cpgString(state.GetAttr(cpgAttr)))
		})
	}
}

// Planned peering_gateway_cxp_id for each prior state and config. A prior of 0
// is state written while an explicit 0 still planned as 0.
func TestPeeringGatewayCxpIdPlanMatrix(t *testing.T) {
	p := newAzureVnetProvider(&alkira.AlkiraClient{})
	noPrior := cty.NullVal(p.ty)
	omitted := cty.NullVal(cty.Number)
	id := func(v int64) cty.Value { return cty.NumberIntVal(v) }

	cases := []struct {
		name   string
		prior  cty.Value
		config cty.Value
		want   string
	}{
		{"create, omitted", noPrior, omitted, "unknown"},
		{"create, 0", noPrior, id(0), "unknown"},
		{"create, 200", noPrior, id(200), "200"},
		{"prior null, omitted", p.state(omitted), omitted, "null"},
		{"prior null, 0", p.state(omitted), id(0), "null"},
		{"prior null, 200", p.state(omitted), id(200), "200"},
		{"prior 0, omitted", p.state(id(0)), omitted, "0"},
		{"prior 0, 0", p.state(id(0)), id(0), "0"},
		{"prior 0, 200", p.state(id(0)), id(200), "200"},
		{"prior 163, omitted", p.state(id(163)), omitted, "163"},
		{"prior 163, 0", p.state(id(163)), id(0), "163"},
		{"prior 163, 163", p.state(id(163)), id(163), "163"},
		{"prior 163, 200", p.state(id(163)), id(200), "200"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, p.planCpg(t, c.prior, p.config(c.config)))
		})
	}
}

// Refresh keeps peering_gateway_cxp_id in state when the API stops returning
// cxpPeeringGatewayId, and picks up a gateway the API assigns after state was
// written without one.
func TestPeeringGatewayCxpIdRefresh(t *testing.T) {
	assigned := 163
	cases := []struct {
		name  string
		prior cty.Value
		cpg   *int
		want  string
	}{
		{"api omits gateway", cty.NumberIntVal(163), nil, "163"},
		{"api assigns gateway to null state", cty.NullVal(cty.Number), &assigned, "163"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := newAzureVnetProvider(createMockAlkiraClient(t, azureVnetConnectorHandler(c.cpg)))

			refreshed, err := p.server.ReadResource(context.Background(), &tfprotov5.ReadResourceRequest{
				TypeName:     azureVnetTypeName,
				CurrentState: p.encode(t, p.state(c.prior)),
			})
			require.NoError(t, err)
			require.Empty(t, refreshed.Diagnostics)

			assert.Equal(t, c.want, cpgString(p.decode(t, refreshed.NewState).GetAttr(cpgAttr)))
		})
	}
}
