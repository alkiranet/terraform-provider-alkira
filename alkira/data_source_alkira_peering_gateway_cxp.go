package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraPeeringGatewayCxp() *schema.Resource {
	return &schema.Resource{
		Description: "This data source allows to retrieve an existing " +
			"Cxp Peering Gateway by its name.",

		Read: dataSourceAlkiraPeeringGatewayCxpRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the resource.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cxp": {
				Description: "The CXP to which the Gateway is attached.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cloud_provider": {
				Description: "The cloud provider where this resource is created.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cloud_region": {
				Description: "The region of the specified cloud provider.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"segment_id": {
				Description: "The ID of the segment that is associated with the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: "The state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"metadata": {
				Description: "Metadata information available once the Peering Gateway is in ACTIVE state.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ilb_ip_address": {
							Description: "Internal Load Balancer IP Address.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"guest_email": {
							Description: "Guest email address.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"vnet_resource_id": {
							Description: "Azure VNET Resource ID associated with the ATH.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"vnet_name": {
							Description: "Azure VNET name associated with the ATH.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"subscription_id": {
							Description: "Azure Subscription ID.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"resource_group": {
							Description: "Azure Resource Group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlkiraPeeringGatewayCxpRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewPeeringGatewayCxp(m.(*alkira.AlkiraClient))

	resource, _, err := api.GetByName(d.Get("name").(string))
	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))

	segmentId, err := getSegmentIdByName(resource.Segment, m)
	if err != nil {
		return err
	}

	d.Set("description", resource.Description)
	d.Set("cxp", resource.Cxp)
	d.Set("cloud_provider", resource.CloudProvider)
	d.Set("cloud_region", resource.CloudRegion)
	d.Set("segment_id", segmentId)
	d.Set("state", resource.State)

	// Set metadata if available (only populated when state is ACTIVE)
	if resource.Metadata != nil {
		metadata := []interface{}{
			map[string]interface{}{
				"ilb_ip_address":   resource.Metadata.IlbIpAddress,
				"guest_email":      resource.Metadata.GuestEmail,
				"vnet_resource_id": resource.Metadata.VnetResourceId,
				"vnet_name":        resource.Metadata.VnetName,
				"subscription_id":  resource.Metadata.SubscriptionId,
				"resource_group":   resource.Metadata.ResourceGroup,
			},
		}
		d.Set("metadata", metadata)
	}

	return nil
}
