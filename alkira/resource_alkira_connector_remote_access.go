package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorRemoteAccess() *schema.Resource {
	return &schema.Resource{
		Description:   "Provide Connector Remote Access resource.",
		CreateContext: resourceConnectorRemoteAccess,
		ReadContext:   resourceConnectorRemoteAccessRead,
		UpdateContext: resourceConnectorRemoteAccessUpdate,
		DeleteContext: resourceConnectorRemoteAccessDelete,
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
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the connector.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"billing_tag_ids": {
				Description: "Billing tags to be associated with " +
					"the resource. (see resource `alkira_billing_tag`).",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the connector should be " +
					"provisioned.",
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Description: "The size of the connector, one of `SMALL`, " +
					"`MEDIUM`, `LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM", "LARGE"}, false),
			},
			"authentication_mode": {
				Description: "Authentication mode, the value could be " +
					"`LOCAL`, `LDAP` and `SAML`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_ids": {
				Description: "Segments that are associated with the " +
					"connector.",
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"ldap_settings": {
				Description: "LDAP Settings when `authentication_mode` " +
					"is `LDAP`.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_user_domain": {
							Description: "The domain.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"ldap_type": {
							Description: "The LDAP type.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"destination_address": {
							Description: "Destination dddress.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"management_segment_id": {
							Description: "The management segment.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"search_scope_domain": {
							Description: "Base DN to query and validate " +
								"remote users that will connect to the " +
								"connector.",
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"enable_dynamic_region_mapping": {
				Description: "Enable dynamic region mapping. Default value " +
					"is `true`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"name_server": {
				Description: "Name server.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"concurrent_sessions_alert_threshold": {
				Description: "The threshhold for concurrent sessions alert.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     80,
			},
			"authorization": {
				Description: "Map Segments of the selected CXP regions to one " +
					"or more User Groups and client subnets.",
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "The segment (`alkira_segment`) to be mapped.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"user_group_name": {
							Description: "User group (`alkira_group_user`) name.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"split_tunneling": {
							Description: "Enable split tunneling to send " +
								"traffic destined to only IP addresses in " +
								"the prefix list over the VPN tunnel. Default " +
								"is `false`",
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
						"prefix_list_id": {
							Description: "The ID of the prefix list (" +
								"`alkira_policy_prefix_list`).",
							Type:     schema.TypeInt,
							Optional: true,
						},
						"billing_tag_id": {
							Description: "Billing tag (`alkira_billing_tag`).",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"subnet": {
							Description: "The client subnet.",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
			},
			"banner_text": {
				Description: "The user provided connectors banner text.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceConnectorRemoteAccess(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorRemoteAccessTemplate(m.(*alkira.AlkiraClient))

	request, err := generateConnectorRemoteAccessRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// CREATE
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	// Set the provision state
	if client.Provision == true {
		d.Set("provision_state", provState)

		if provErr != nil {
			return diag.Diagnostics{{
				Severity: diag.Warning,
				Summary:  "PROVISION (CREATE) FAILED",
				Detail:   fmt.Sprintf("%s", provErr),
			}}
		}
	}

	return resourceConnectorRemoteAccessRead(ctx, d, m)
}

func resourceConnectorRemoteAccessRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorRemoteAccessTemplate(m.(*alkira.AlkiraClient))

	// Get
	connector, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	err = setConnectorRemoteAccess(connector, d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourceConnectorRemoteAccessUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	// INIT
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorRemoteAccessTemplate(m.(*alkira.AlkiraClient))

	request, err := generateConnectorRemoteAccessRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// UPDATE
	provState, err, provErr := api.Update(d.Id(), request)

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

	return nil
}

func resourceConnectorRemoteAccessDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewConnectorRemoteAccessTemplate(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	return nil
}
