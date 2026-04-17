package alkira

import (
	"errors"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getCheckpointWriteOnlyValue reads a WriteOnly field from raw config.
// WriteOnly fields are never stored in state, so d.GetOk() would always
// return the zero value — only GetRawConfigAt can read the actual value.
func getCheckpointWriteOnlyValue(d *schema.ResourceData, field string) string {
	attrPath := cty.Path{cty.GetAttrStep{Name: field}}
	val, diags := d.GetRawConfigAt(attrPath)

	if !diags.HasError() && !val.IsNull() && val.IsKnown() && val.Type() == cty.String {
		if strVal := val.AsString(); strVal != "" {
			return strVal
		}
	}

	return ""
}

// getNestedWriteOnlyValue reads a WriteOnly field from a nested TypeList block
// via GetRawConfigAt with an indexed cty path.
func getNestedWriteOnlyValue(d *schema.ResourceData, block string, index int, field string) string {
	attrPath := cty.Path{
		cty.GetAttrStep{Name: block},
		cty.IndexStep{Key: cty.NumberIntVal(int64(index))},
		cty.GetAttrStep{Name: field},
	}
	val, diags := d.GetRawConfigAt(attrPath)

	if !diags.HasError() && !val.IsNull() && val.IsKnown() && val.Type() == cty.String {
		if strVal := val.AsString(); strVal != "" {
			return strVal
		}
	}

	return ""
}

// createCheckpointCredential create checkpoint service credential
func createCheckpointCredential(d *schema.ResourceData, c *alkira.AlkiraClient) (string, error) {
	log.Printf("[INFO] Creating Checkpoint service credential")

	credentialName := d.Get("name").(string) + "-" + randomNameSuffix()
	credential := alkira.CredentialCheckPointFwService{AdminPassword: getCheckpointWriteOnlyValue(d, "password")}

	return c.CreateCredential(credentialName, alkira.CredentialTypeChkpFw, credential, 0)
}

// updateCheckpointCredential always updates the checkpoint service
// credential because password is WriteOnly and HasChanges cannot
// detect changes (state is always null for WriteOnly fields).
func updateCheckpointCredential(d *schema.ResourceData, c *alkira.AlkiraClient) error {
	log.Printf("[INFO] Updating Checkpoint service credential")

	credentialId := d.Get("credential_id").(string)
	if credentialId == "" {
		return errors.New("credential_id is empty when updating Checkpoint credential")
	}

	cred, err := c.GetCredentialById(credentialId)
	if err != nil {
		return fmt.Errorf("failed to get checkpoint credential name for ID %s: %w", credentialId, err)
	}

	credential := alkira.CredentialCheckPointFwService{
		AdminPassword: getCheckpointWriteOnlyValue(d, "password"),
	}

	return c.UpdateCredential(credentialId, cred.Name, alkira.CredentialTypeChkpFw, credential, 0)
}

// updateCheckpointManagementServerCredential always updates the
// management server credential because password is WriteOnly.
func updateCheckpointManagementServerCredential(d *schema.ResourceData, c *alkira.AlkiraClient) error {
	log.Printf("[INFO] Updating Checkpoint management server credential")

	// Read credential_id from management_server TypeSet in state
	mgSet := d.Get("management_server").(*schema.Set)
	if mgSet == nil || mgSet.Len() == 0 {
		return nil
	}

	// management_server is validated to have exactly 1 element
	mgList := mgSet.List()
	cfg := mgList[0].(map[string]interface{})

	configMode, _ := cfg["configuration_mode"].(string)
	if configMode != "AUTOMATED" {
		log.Printf("[INFO] Management server configuration_mode is %q, skipping credential update", configMode)
		return nil
	}

	credentialId, _ := cfg["credential_id"].(string)
	password, _ := cfg["password"].(string)

	if credentialId == "" {
		// No credential exists yet — create one (handles mode transition)
		log.Printf("[INFO] management_server credential_id is empty, creating credential")
		name := d.Get("name").(string)
		manServerCredName := name + "-" + randomNameSuffix()
		credential := &alkira.CredentialCheckPointFwManagementServer{Password: password}
		_, err := c.CreateCredential(manServerCredName, alkira.CredentialTypeChkpFwManagement, credential, 0)
		return err
	}

	cred, err := c.GetCredentialById(credentialId)
	if err != nil {
		return fmt.Errorf("failed to get management server credential name for ID %s: %w", credentialId, err)
	}

	credential := &alkira.CredentialCheckPointFwManagementServer{Password: password}

	return c.UpdateCredential(credentialId, cred.Name, alkira.CredentialTypeChkpFwManagement, credential, 0)
}

func expandCheckpointManagementServer(name string, in *schema.Set, m interface{}) (*alkira.CheckpointManagementServer, error) {

	client := m.(*alkira.AlkiraClient)

	if in == nil || in.Len() > 1 {
		log.Printf("[DEBUG] Invalid Checkpoint Firewall Management Server input")
		return nil, errors.New("ERROR: Invalid checkpoint firewall management server input")
	}

	if in.Len() < 1 {
		return nil, nil
	}

	mg := &alkira.CheckpointManagementServer{}
	var manServerPass string

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})
		if v, ok := cfg["configuration_mode"].(string); ok {
			mg.ConfigurationMode = v
		}
		if v, ok := cfg["password"].(string); ok {
			manServerPass = v
		}
		if v, ok := cfg["credential_id"].(string); ok {
			if v == "" && mg.ConfigurationMode == "AUTOMATED" {
				manServerCredName := name + "-" + randomNameSuffix()
				c := &alkira.CredentialCheckPointFwManagementServer{Password: manServerPass}
				credentialId, err := client.CreateCredential(manServerCredName, alkira.CredentialTypeChkpFwManagement, c, 0)
				if err != nil {
					return nil, err
				}
				mg.CredentialId = credentialId
			}

			if v != "" {
				mg.CredentialId = v
			}
		}
		if v, ok := cfg["domain"].(string); ok {
			mg.Domain = v
		}
		if v, ok := cfg["global_cidr_list_id"].(int); ok {
			mg.GlobalCidrListId = v
		}
		if v, ok := cfg["ips"].([]interface{}); ok {
			mg.Ips = convertTypeListToStringList(v)
		}
		if v, ok := cfg["reachability"].(string); ok {
			mg.Reachability = v
		}
		if v, ok := cfg["segment_id"].(string); ok {
			if v != "" {
				segment, err := getSegmentNameById(v, m)

				if err != nil {
					return nil, err
				}

				mg.Segment = segment
			}
		}
		if v, ok := cfg["type"].(string); ok {
			mg.Type = v
		}
		if v, ok := cfg["username"].(string); ok {
			mg.UserName = v
		}
	}
	return mg, nil
}

func expandCheckpointInstances(d *schema.ResourceData, in []interface{}, m interface{}) ([]alkira.CheckpointInstance, error) {

	if in == nil || len(in) == 0 {
		return nil, errors.New("ERROR: Invalid checkpoint firewall instance input")
	}

	client := m.(*alkira.AlkiraClient)

	instances := make([]alkira.CheckpointInstance, len(in))
	for i, instance := range in {
		r := alkira.CheckpointInstance{}
		instanceCfg := instance.(map[string]interface{})

		// Read sic_key from raw config (WriteOnly field not available via d.Get)
		sicKey := getNestedWriteOnlyValue(d, "instance", i, "sic_key")

		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			if v == "" {
				credentialName := r.Name + "-" + randomNameSuffix()
				c := &alkira.CredentialCheckPointFwServiceInstance{SicKey: sicKey}

				log.Printf("[INFO] Creating Credential CheckpointInstance.")
				credentialId, err := client.CreateCredential(
					credentialName,
					alkira.CredentialTypeChkpFwInstance,
					c,
					0)

				if err != nil {
					return nil, err
				}

				r.CredentialId = credentialId
			} else {
				r.CredentialId = v
			}
		}
		if v, ok := instanceCfg["enable_traffic"].(bool); ok {
			r.TrafficEnabled = v
		}
		instances[i] = r
	}

	return instances, nil
}

// Checkpoint expects segment_options to not be empty.
// If segment_options is not defined in the TF file, this function adds the default expected data.
// If segment_options is included, populates it normally.
func expandCheckpointSegmentOptions(segmentName string, in *schema.Set, m interface{}) (alkira.SegmentNameToZone, error) {

	if in == nil || in.Len() == 0 {

		segmentOptions := make(alkira.SegmentNameToZone)
		zonestoGroups := make(alkira.ZoneToGroups)

		z := alkira.OuterZoneToGroups{}
		j := []string{}

		zonestoGroups["DEFAULT"] = j
		z.ZonesToGroups = zonestoGroups

		segmentOptions[segmentName] = z

		return segmentOptions, nil
	}

	return expandSegmentOptions(in, m)

}

func deflateCheckpointManagementServer(mg alkira.CheckpointManagementServer, m interface{}) []map[string]interface{} {
	result := make(map[string]interface{})
	result["configuration_mode"] = mg.ConfigurationMode
	result["credential_id"] = mg.CredentialId
	result["domain"] = mg.Domain
	result["global_cidr_list_id"] = mg.GlobalCidrListId
	result["ips"] = convertStringArrToInterfaceArr(mg.Ips)
	result["reachability"] = mg.Reachability
	result["type"] = mg.Type
	result["username"] = mg.UserName

	// Convert segment name to segment ID for import support
	if mg.Segment != "" && m != nil {
		segmentId, err := getSegmentIdByName(mg.Segment, m)
		if err == nil {
			result["segment_id"] = segmentId
		}
	}

	return []map[string]interface{}{result}
}

func setCheckpointInstances(d *schema.ResourceData, c []alkira.CheckpointInstance) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, value := range d.Get("instance").([]interface{}) {
		cfg := value.(map[string]interface{})

		for _, ins := range c {
			if cfg["id"].(int) == ins.Id || cfg["name"].(string) == ins.Name {
				instance := map[string]interface{}{
					"credential_id":  ins.CredentialId,
					"name":           ins.Name,
					"id":             ins.Id,
					"sic_key":        cfg["sic_key"].(string),
					"enable_traffic": ins.TrafficEnabled,
				}
				instances = append(instances, instance)
				break
			}
		}
	}

	for _, instance := range c {
		new := true

		// Check if the instance already exists in the Terraform config
		for _, ins := range d.Get("instance").([]interface{}) {
			cfg := ins.(map[string]interface{})

			if cfg["id"].(int) == instance.Id || cfg["name"].(string) == instance.Name {
				new = false
				break
			}
		}

		// If the instance is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			instance := map[string]interface{}{
				"credential_id": instance.CredentialId,
				"name":          instance.Name,
				"id":            instance.Id,
			}

			instances = append(instances, instance)
		}
	}

	return instances
}

// generateCheckpointRequest
func generateCheckpointRequest(d *schema.ResourceData, m interface{}) (*alkira.ServiceCheckpoint, error) {

	// Management Server block
	managementServer, err := expandCheckpointManagementServer(d.Get("name").(string), d.Get("management_server").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	//
	// Instances block
	//
	instances, err := expandCheckpointInstances(d, d.Get("instance").([]interface{}), m)

	if err != nil {
		return nil, err
	}

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	//
	// Segment Options
	//
	segmentOptions, err := expandCheckpointSegmentOptions(segmentName, d.Get("segment_options").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	// Assemble request
	return &alkira.ServiceCheckpoint{
		AutoScale:        d.Get("auto_scale").(string),
		BillingTags:      convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		CredentialId:     d.Get("credential_id").(string),
		Cxp:              d.Get("cxp").(string),
		Description:      d.Get("description").(string),
		Instances:        instances,
		LicenseType:      d.Get("license_type").(string),
		ManagementServer: managementServer,
		MinInstanceCount: d.Get("min_instance_count").(int),
		MaxInstanceCount: d.Get("max_instance_count").(int),
		Name:             d.Get("name").(string),
		PdpIps:           convertTypeListToStringList(d.Get("pdp_ips").([]interface{})),
		Segments:         []string{segmentName},
		SegmentOptions:   segmentOptions,
		Size:             d.Get("size").(string),
		TunnelProtocol:   d.Get("tunnel_protocol").(string),
		Version:          d.Get("version").(string),
	}, nil
}
