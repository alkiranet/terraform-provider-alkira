package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/stretchr/testify/assert"
)

func TestInternetApplicationsFieldNameMatch(t *testing.T) {
	t.Run("schema uses target field name", func(t *testing.T) {
		resourceSchema := resourceAlkiraInternetApplication().Schema

		// Verify "target" field exists in schema
		targetField, exists := resourceSchema["target"]
		assert.True(t, exists, "Schema must have 'target' field")
		assert.NotNil(t, targetField, "target field must not be nil")

		// Verify "targets" (plural - the bug) does NOT exist in schema
		_, wrongFieldExists := resourceSchema["targets"]
		assert.False(t, wrongFieldExists, "Schema should NOT have 'targets' field (bug was using plural)")
	})
}

// AK-71658: Read did not persist inbound_connector_type, so `terraform import`
// left it null in state and the next plan showed a spurious
// `+ inbound_connector_type = "DEFAULT"`.
func TestInternetApplicationSetFields(t *testing.T) {
	t.Run("inbound_connector_type is written to state", func(t *testing.T) {
		d := resourceAlkiraInternetApplication().TestResourceData()

		setInternetApplicationFields(d, &alkira.InternetApplication{
			InboundConnectorType: "AKAMAI_PROLEXIC",
		})

		assert.Equal(t, "AKAMAI_PROLEXIC", d.Get("inbound_connector_type"),
			"Read must persist inbound_connector_type, otherwise import leaves it null")
	})

	t.Run("api response fields are written to state", func(t *testing.T) {
		d := resourceAlkiraInternetApplication().TestResourceData()

		app := &alkira.InternetApplication{
			BiDirectionalAvailabilityZone: "AZ1",
			ByoipId:                       42,
			ConnectorId:                   3490,
			ConnectorType:                 "AWS_VPC",
			Description:                   "test ifa",
			FqdnPrefix:                    "ifaapp",
			IlbCredentialId:               "cred-1",
			InboundConnectorType:          "DEFAULT",
			InternetProtocol:              "IPV4",
			Name:                          "ifa-app",
			PublicIps:                     []string{"1.2.3.4"},
			Size:                          "SMALL",
		}

		setInternetApplicationFields(d, app)

		assert.Equal(t, app.BiDirectionalAvailabilityZone, d.Get("bi_directional_az"))
		assert.Equal(t, app.ByoipId, d.Get("byoip_id"))
		assert.Equal(t, app.ConnectorId, d.Get("connector_id"))
		assert.Equal(t, app.ConnectorType, d.Get("connector_type"))
		assert.Equal(t, app.Description, d.Get("description"))
		assert.Equal(t, app.FqdnPrefix, d.Get("fqdn_prefix"))
		assert.Equal(t, app.IlbCredentialId, d.Get("ilb_credential_id"))
		assert.Equal(t, app.InboundConnectorType, d.Get("inbound_connector_type"))
		assert.Equal(t, app.InternetProtocol, d.Get("internet_protocol"))
		assert.Equal(t, app.Name, d.Get("name"))
		assert.Equal(t, app.Size, d.Get("size"))
	})
}
