package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraConnectorIpsecTunnelProfile() *schema.Resource {
	return &schema.Resource{
		Description:   "Manages IPSec Tunnel Profile.",
		CreateContext: resourceConnectorIpsecTunnelProfile,
		ReadContext:   resourceConnectorIpsecTunnelProfileRead,
		UpdateContext: resourceConnectorIpsecTunnelProfileUpdate,
		DeleteContext: resourceConnectorIpsecTunnelProfileDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the profile.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the tunnel profile.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ipsec_encryption_algorithm": {
				Description: "ESP encryption algorithm of the IPSec Tunnel. " +
					"The value could be: `AES256CBC`, `AES192CBC`, `AES128CBC`, " +
					"`3DESCBC`, `AES256GCM16`, `AES192GCM16`, `AES128GCM16`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"ipsec_integrity_algorithm": {
				Description: "ESP integrity algorithm of the IPSec tunnel. " +
					"The value could be: `SHA1`, `SHA256`, `SHA384`, `SHA512` " +
					"and `MD5`.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"ipsec_dh_group": {
				Description: "ESP DH group number of the IPSec tunnel. The value " +
					"could be: `MODP1024`, `MODP2048`, `MODP3072`, `MODP4096`, " +
					"`MODP6144`, `MODP8192`, `ECP256`, `ECP384`, `ECP521`, " +
					"`CURVE25519`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"ike_encryption_algorithm": {
				Description: "ESP encryption algorithm used with IKE tunnel. " +
					"The value could be: `AES256CBC`, `AES192CBC`, `AES128CBC`, " +
					"`3DESCBC`, `AES256GCM16`, `AES192GCM16`, `AES128GCM16`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"ike_integrity_algorithm": {
				Description: "ESP integrity algorithm used with IKE tunnel. " +
					"The value could be: `SHA1`, `SHA256`, `SHA384`, `SHA512` " +
					"and `MD5`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"ike_dh_group": {
				Description: "ESP DH group number of the IKE tunnel. The value " +
					"could be: `MODP1024`, `MODP2048`, `MODP3072`, `MODP4096`, " +
					"`MODP6144`, `MODP8192`, `ECP256`, `ECP384`, `ECP521`, " +
					"`CURVE25519`.",
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceConnectorIpsecTunnelProfile(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorIPSecTunnelProfile(client)

	// Construct request
	req, err := generateConnectorIpsecTunnelProfileRequest(d)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, provErr := api.Create(req)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provState == "FAILED" {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorIpsecTunnelProfileRead(ctx, d, m)
}

func resourceConnectorIpsecTunnelProfileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorIPSecTunnelProfile(client)

	// Get the resource
	profile, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", profile.Name)
	d.Set("description", profile.Description)
	d.Set("ipsec_encryption_algorithm", profile.IpSecConfiguration.EncryptionAlgorithm)
	d.Set("ipsec_integrity_algorithm", profile.IpSecConfiguration.IntegrityAlgorithm)
	d.Set("ipsec_dh_group", profile.IpSecConfiguration.DhGroup)
	d.Set("ike_encryption_algorithm", profile.IkeConfiguration.EncryptionAlgorithm)
	d.Set("ike_integrity_algorithm", profile.IkeConfiguration.IntegrityAlgorithm)
	d.Set("ike_dh_group", profile.IkeConfiguration.DhGroup)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorIpsecTunnelProfileUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorIPSecTunnelProfile(client)

	// Construct request
	profile, err := generateConnectorIpsecTunnelProfileRequest(d)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
	provState, err, provErr := api.Update(d.Id(), profile)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (UPDATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorIpsecTunnelProfileRead(ctx, d, m)
}

func resourceConnectorIpsecTunnelProfileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorIPSecTunnelProfile(client)

	// Delete
	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	// Check provisions state
	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}

func generateConnectorIpsecTunnelProfileRequest(d *schema.ResourceData) (*alkira.ConnectorIPSecTunnelProfile, error) {

	profile := &alkira.ConnectorIPSecTunnelProfile{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		IpSecConfiguration: alkira.ConnectorIPSecTunnelProfileIpSecConfiguration{
			EncryptionAlgorithm: d.Get("ipsec_encryption_algorithm").(string),
			IntegrityAlgorithm:  d.Get("ipsec_integrity_algorithm").(string),
			DhGroup:             d.Get("ipsec_dh_group").(string),
		},
		IkeConfiguration: alkira.ConnectorIPSecTunnelProfileIkeConfiguration{
			EncryptionAlgorithm: d.Get("ike_encryption_algorithm").(string),
			IntegrityAlgorithm:  d.Get("ike_integrity_algorithm").(string),
			DhGroup:             d.Get("ike_dh_group").(string),
		},
	}

	return profile, nil
}
