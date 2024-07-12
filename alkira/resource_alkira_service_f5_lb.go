package alkira

import (
	"context"
	"fmt"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraF5LoadBalancer() *schema.Resource {
	return &schema.Resource{
		Description:   "",
		CreateContext: resourceF5LoadBalancerCreate,
		ReadContext:   resourceF5LoadBalancerRead,
		UpdateContext: resourceF5LoadBalancerUpdate,
		DeleteContext: resourceF5LoadBalancerDelete,
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
				Description: "Name of the service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "Description of the service.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"cxp": {
				Description: "The CXP where the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "The size of the service, one of" +
					" `SMALL`, `MEDIUM`, `LARGE`" +
					" `2LARGE`, `5LARGE`.",
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"SMALL", "MEDIUM", "LARGE", "2LARGE", "5LARGE"}, false),
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"billing_tag_ids": {
				Description: "IDs of billing tags to associate with " +
					"the service.",
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt}},
			"elb_cidrs": {
				Description: "",
				Type:        schema.TypeSet,
				Required:    true,
			},
			"big_ip_allow_list": {
				Description: "",
				Type:        schema.TypeSet,
				Required:    true,
			},
			"instances": {
				Description: "An array containing the properties for each F5 load" +
					" balancer instance that needs to be deployed.",
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{Schema: map[string]*schema.Schema{
					"name": {
						Description: "Name of the F5 load balancer instance.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"registration_credential_id": {
						Description: "The ID of the F5 load balancer registration",
						Type:        schema.TypeString,
						Required:    true,
					},
					"credential_id": {
						Description: "The ID of the F5 load balancer credential.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"license_type": {
						Description: "The type of license used for the F5 load balancer instance.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"version": {
						Description: "The version of the F5 load balancer instance.",
						Type:        schema.TypeString,
						Required:    true,
					},
					"deployment_option": {
						Description: "The deployment option used for the F5 load balancer instance.",
						Type:        schema.TypeString,
						Optional:    true,
					},
					"deployment_type": {
						Description: "The deployment type used for the F5 load balancer instance.",
						Type:        schema.TypeString,
						Required:    true,
					},
				},
				},
			},
		},
	}
}

func resourceF5LoadBalancerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceF5Lb(client)

	request, err := generateRequestF5Lb(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(response.Id))

	if client.Provision {
		d.Set("provision_state", provState)
		if provState == "FAILED" {
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "PROVISION (CREATE) FAILED",
					Detail:   fmt.Sprintf("%s", provErr),
				},
			}
		}
	}
	return resourceF5LoadBalancerRead(ctx, d, m)
}
func resourceF5LoadBalancerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceF5Lb(m.(*alkira.AlkiraClient))

	lb, provState, err := api.GetById(d.Id())
	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", lb.Name)
	d.Set("description", lb.Description)
	d.Set("cxp", lb.Cxp)
	d.Set("size", lb.Size)
	d.Set("elb_cidrs", lb.ElbCidrs)
	d.Set("big_ip_allow_list", lb.BigIpAllowList)
	d.Set("instances", lb.Instances)
	d.Set("billing_tag_ids", lb.BillingTags)

	// Set segments
	segments := make([]int, len(lb.Segments))

	for _, seg := range lb.Segments {
		seg, err := getSegmentIdByName(seg, m)

		if err != nil {
			return diag.FromErr(err)
		}
		segId, _ := strconv.Atoi(seg)
		segments = append(segments, segId)
	}
	d.Set("segment_ids", segments)

	// Set provision state
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil

}
func resourceF5LoadBalancerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceF5Lb(m.(*alkira.AlkiraClient))

	request, err := generateRequestF5Lb(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	provState, err, provErr := api.Update(d.Id(), request)

	if client.Provision == true {
		d.Set("provision_state", provState)
		if provState == "FAILED" {
			return diag.Diagnostics{
				{
					Severity: diag.Warning,
					Summary:  "PROVISION (UPDATE) FAILED",
					Detail:   fmt.Sprintf("%s", provErr),
				},
			}
		}
	}
	return nil
}
func resourceF5LoadBalancerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*alkira.AlkiraClient)
	api := alkira.NewServiceF5Lb(m.(*alkira.AlkiraClient))

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
