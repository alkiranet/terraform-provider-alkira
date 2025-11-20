package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraServiceBluecat() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get an existing Bluecat service.",
		ReadContext: dataSourceServiceBluecatRead,
		Schema: map[string]*schema.Schema{
			"bddsAnycast": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Defines the BDDS AnyCast policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ips": {
							Description: "The IPs to be used when AnyCast is enabled.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"backup_cxps": {
							Description: "The backup CXPs to be used when the current Bluecat service is not available.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"edgeAnycast": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Defines the Edge AnyCast policy.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ips": {
							Description: "The IPs to be used when AnyCast is enabled.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"backup_cxps": {
							Description: "The backup CXPs to be used when the current Bluecat service is not available.",
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"billing_tag_ids": {
				Description: "IDs of billing tags associated with the service.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"cxp": {
				Description: "The CXP where the service is provisioned.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"description": {
				Description: "The description of the Bluecat service.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"global_cidr_list_id": {
				Description: "The ID of the global cidr list associated with the Bluecat service.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"id": {
				Description: "ID of the Bluecat service.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"instance": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The properties pertaining to each individual instance of the Bluecat service.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the Bluecat instance.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"id": {
							Description: "The ID of the Bluecat instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"bddsOptions": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "BDDS options for the instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"clientId": {
										Description: "The license clientId of the Bluecat BDDS instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"activationKey": {
										Description: "The license activationKey of the Bluecat BDDS instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"hostname": {
										Description: "The host name of the instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"model": {
										Description: "The model of the Bluecat BDDS instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"version": {
										Description: "The version of the Bluecat BDDS instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"edgeOptions": {
							Type:        schema.TypeSet,
							Computed:    true,
							Description: "Edge options for the instance.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"configData": {
										Description: "The config data of the Bluecat Edge instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"hostname": {
										Description: "The host name of the Edge instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
									"version": {
										Description: "The version of the Bluecat Edge instance.",
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"type": {
							Description: "The type of the Bluecat instance (BDDS or EDGE).",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"license_type": {
				Description: "Bluecat license type.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Name of the Bluecat service.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"service_group_name": {
				Description: "The name of the service group associated with the service.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"service_group_id": {
				Description: "The ID of the service group associated with the service.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"service_group_implicit_group_id": {
				Description: "The ID of the implicit group associated with the service.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func dataSourceServiceBluecatRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewServiceBluecat(m.(*alkira.AlkiraClient))

	serviceBluecat, provState, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(serviceBluecat.Id))
	setAllBluecatResourceFields(d, serviceBluecat)

	// Set provision state
	d.Set("provision_state", provState)

	return nil
}
