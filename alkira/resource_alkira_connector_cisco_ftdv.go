package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraConnectorCiscoFTDv() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Cisco FTDv Connector. (**BETA**)",

		Create: resourceConnectorCiscoFTDvCreate,
		Read:   resourceConnectorCiscoFTDvRead,
		Update: resourceConnectorCiscoFTDvUpdate,
		Delete: resourceConnectorCiscoFTDvDelete,

		Schema: map[string]*schema.Schema{
			// "id": {
			// 	Description: "The ID of the Cisco FTDv Firewall.",
			// 	Type:        schema.TypeInt,
			// 	Computed:    true,
			// },
			"name": {
				Description: "The name of the connector.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"credential_id": {
				Description: "ID of Checkpoint Firewall credential.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"auto_scale": {
				Description: "Indicate if `auto_scale` should be enabled for your Cisco FTDv connector." +
					" `ON` and `OFF` are accepted values. `OFF` is the default if " +
					"field is omitted.",
				Type:         schema.TypeString,
				Default:      "OFF",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ON", "OFF"}, false),
			},
			"size": {
				Description:  "The size of the connector, one of `SMALL`, `MEDIUM`, `LARGE`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"SMALL", `MEDIUM`, `LARGE`}, false),
			},
			"tunnel_protocol": {
				Description:  "The tunnel protocol. One of `IPSEC`. Default is `IPSEC`",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "IPSEC",
				ValidateFunc: validation.StringInSlice([]string{"IPSEC"}, false),
			},
			"cxp": {
				Description: "The CXP where the connector should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"global_cidr_list_id": {
				Description: "The ID of the global cidr list to be associated with " +
					"the management server.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_instance_count": {
				Description: "The maximum number of Cisco FTDv instances that should be deployed.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"min_instance_count": {
				Description: "The minimum number of Cisco FTDv instances that should be deployed.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"ip_allow_list": {
				Description: "",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"billing_tag_ids": {
				Description: "IDs of Billing Tags.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"username": {
				Description: "",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "",
				Type:        schema.TypeString,
				Required:    true,
			},
			"management_server": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"fmc_ip": {
							Description: "",
							Type:        schema.TypeString,
							Required:    true,
						},
						"segment_name": {
							Description: "",
							Type:        schema.TypeString,
							Required:    true,
						},
						"segment_id": {
							Description: "The ID of the segment to be used to access the management server.",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
			},
			"instances": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The ID of the Cisco FTDv instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"credential_id": {
							Description: "",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"internal_name": {
							Description: "",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"state": {
							Description: "",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"hostname": {
							Description: "",
							Type:        schema.TypeString,
							Optional:    true,
						},
						"version": {
							Description: "",
							Type:        schema.TypeString,
							Required:    true,
						},
						"license_type": {
							Description:  ", either `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"}, false),
						},
						"admin_password": {
							Description: "",
							Type:        schema.TypeString,
							Required:    true,
						},
						"fmc_registration_key": {
							Description: "",
							Type:        schema.TypeString,
							Required:    true,
						},
						"ftdv_nat_id": {
							Description: "",
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
			"segment_options": {
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "The ID of the segment.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"zone_name": {
							Description: "The name of the associated zone.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"groups": {
							Description: "The list of Groups associated with the zone.",
							Type:        schema.TypeList,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

// resourceConnectorCiscoFTDvCreate create a Cisco FTDv connector
func resourceConnectorCiscoFTDvCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorCiscoFTDvRequest(d, m)

	if err != nil {
		return err
	}

	id, err := client.CreateConnectorCiscoFTDv(connector)

	if err != nil {
		return err
	}

	d.SetId(id)
	return resourceConnectorCiscoFTDvRead(d, m)
}

// resourceConnectorCiscoFTDvRead get and save a Cisco FTDv connectors
func resourceConnectorCiscoFTDvRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	connector, err := client.GetConnectorCiscoFTDv(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", connector.Name)
	d.Set("auto_scale", connector.AutoScale)
	d.Set("size", connector.Size)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("cxp", connector.Cxp)
	d.Set("tunnel_protocol", connector.TunnelProtocol)
	d.Set("max_instance_count", connector.MaxInstanceCount)
	d.Set("min_instance_count", connector.MinInstanceCount)
	d.Set("ip_allow_list", connector.IpAllowList)
	d.Set("global_cidr_list_id", connector.GlobalCidrListId)
	d.Set("credential_id", connector.CredentialId)
	d.Set("management_server", connector.ManagementServer)

	var instances []map[string]interface{}
	for _, instance := range connector.Instances {
		i := map[string]interface{}{
			"id":            instance.Id,
			"credential_id": instance.CredentialId,
			"internal_name": instance.InternalName,
			"state":         instance.State,
			"hostname":      instance.Hostname,
			"version":       instance.Version,
			"license_type":  instance.LicenseType,
		}
		instances = append(instances, i)
	}
	d.Set("instances", instances)

	return nil
}

// resourceConnectorCiscoFTDvUpdate update a Cisco FTDv connector
func resourceConnectorCiscoFTDvUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	connector, err := generateConnectorCiscoFTDvRequest(d, m)

	if err != nil {
		return fmt.Errorf("UpdateConnectorCiscoFTDv: failed to marshal: %v", err)
	}

	err = client.UpdateConnectorCiscoFTDv(d.Id(), connector)

	if err != nil {
		return err
	}

	return resourceConnectorCiscoFTDvRead(d, m)
}

func resourceConnectorCiscoFTDvDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	err := client.DeleteConnectorCiscoFTDv((d.Id()))

	return err
}

// generateConnectorCiscoFTDvRequest generate a request for Azure ExpressRoute connector
func generateConnectorCiscoFTDvRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorCiscoFTDv, error) {
	client := m.(*alkira.AlkiraClient)

	ciscoFTDvCredId := d.Get("credential_id").(string)
	if 0 == len(ciscoFTDvCredId) {
		log.Printf("[INFO] Creating Cisco FTDv Firewall Service Credentials")
		ciscoFTDvName := d.Get("name").(string) + "-" + randomNameSuffix()
		c := alkira.CredentialCiscoFtdv{Username: d.Get("username").(string), Password: d.Get("password").(string)}
		credentialId, err := client.CreateCredential(ciscoFTDvName, alkira.CredentialTypeCiscoFtdv, c, 0)
		if err != nil {
			return nil, err
		}
		d.Set("credential_id", credentialId)

	}

	managementServer, err := expandCiscoFtdvManagementServer(d.Get("management_server").(*schema.Set))
	if err != nil {
		return nil, err
	}

	ids := convertTypeListToStringList(d.Get("segment_ids").([]interface{}))
	segment_names, err := convertSegmentIdsToSegmentNames(ids, m)
	if err != nil {
		return nil, err
	}

	// segmentOptions, err := expandCiscoFtdvSegmentOptions(d.Get("segment_options").(*schema.Set), m)
	// if err != nil {
	// 	return nil, err
	// }

	billingTags := convertTypeListToIntList(d.Get("billing_tag_ids").([]interface{}))

	instances, err := expandCiscoFTDvInstances(d.Get("name").(string), d.Get("instances").(*schema.Set), m)
	if err != nil {
		return nil, err
	}

	request := &alkira.ConnectorCiscoFTDv{
		// Id:               d.Get("id").(int),
		Name:             d.Get("name").(string),
		GlobalCidrListId: d.Get("global_cidr_list_id").(int),
		Size:             d.Get("size").(string),
		CredentialId:     d.Get("credential_id").(string),
		Cxp:              d.Get("cxp").(string),
		ManagementServer: managementServer,
		IpAllowList:      convertTypeListToStringList(d.Get("ip_allow_list").([]interface{})),
		MaxInstanceCount: d.Get("max_instance_count").(int),
		MinInstanceCount: d.Get("min_instance_count").(int),
		Segments:         segment_names,
		// SegmentOptions:   segmentOptions,
		AutoScale:      d.Get("auto_scale").(string),
		TunnelProtocol: d.Get("tunnel_protocol").(string),
		BillingTags:    billingTags,
		Instances:      instances,
	}

	return request, nil
}
