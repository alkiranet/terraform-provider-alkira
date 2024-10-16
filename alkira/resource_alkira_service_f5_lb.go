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
		Description:   "F5 Load Balancer Service.",
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
				Description: "CXP on which the service should be provisioned.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"size": {
				Description: "Size of the service, one of" +
					" `SMALL`, `MEDIUM`, `LARGE`" +
					" `2LARGE`, `5LARGE`.",
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_ids": {
				Description: "IDs of segments associated with the service.",
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"billing_tag_ids": {
				Description: "IDs of billing tags to associate with" +
					" the service.",
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			"global_cidr_list_id": {
				Description: "ID of global CIDR list from which subnets" +
					" will be allocated for the external network interfaces of" +
					" instances. These interfaces host the public IP" +
					" addresses needed for virtual IPs.",
				Type:     schema.TypeInt,
				Required: true,
			},
			"prefix_list_id": {
				Description: "ID of prefix list to use for IP allowlist",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"segment_options": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "The segment options as used by your F5 Load Balancer.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Description: "ID of the segment.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"nat_pool_prefix_length": {
							Description: "Prefix length of subnets for the segment.",
							Type:        schema.TypeInt,
							Optional:    true,
						},
						"elb_nic_count": {
							Description: "Number of nics to allocate for the segment.",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
			},
			"service_group_name": {
				Description: "Name of the service group to be associated with the service.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"instances": {
				Description: "An array containing the properties for each F5 load" +
					" balancer instance.",
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "Name of the F5 load balancer instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"id": {
							Description: "ID of the F5 load balancer instance.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"license_type": {
							Description: "The type of license used for the F5 load balancer instance." +
								" Can be one of `BRING_YOUR_OWN` or `PAY_AS_YOU_GO`",
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice(
								[]string{"BRING_YOUR_OWN", "PAY_AS_YOU_GO"},
								false),
							Required: true,
						},
						"registration_credential_id": {
							Description: "ID of the F5 load balancer registration credential." +
								" If the `registration_credential_id` is not passed, `f5_registration_key`" +
								" is required to create new credentials." +
								" Only required if `license_type` is `BRING_YOUR_OWN`.",
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"f5_registration_key": {
							Description: "Registration key for the F5 load balancer." +
								" Only required if `license_type` is `BRING_YOUR_OWN`." +
								" This can also be set by `ALKIRA_F5_REGISTRATION_KEY`" +
								" environment variable.",
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							DefaultFunc: envDefaultFunc("ALKIRA_F5_REGISTRATION_KEY"),
						},
						"f5_username": {
							Description: "Username for the F5 load balancer." +
								" Username is `admin` and cannot be changed.",
							Type:     schema.TypeString,
							Computed: true,
							// read only for now, so the environment variable is not set
							// Optional:    true,
							// Sensitive:   true,
							// DefaultFunc: envDefaultFunc("ALKIRA_F5_USERNAME"),
						},
						"f5_password": {
							Description: "Password for the F5 load balancer." +
								" This can also be set by `ALKIRA_F5_PASSWORD`" +
								" environment variable.",
							Type:        schema.TypeString,
							Optional:    true,
							Sensitive:   true,
							DefaultFunc: envDefaultFunc("ALKIRA_F5_PASSWORD"),
						},
						"credential_id": {
							Description: "ID of the F5 load balancer credential." +
								" If the `credential_id` is not passed," +
								" `f5_username` and `f5_password` is required" +
								" to create new credentials.",
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"version": {
							Description: "The version of the F5 load balancer instance.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"deployment_option": {
							Description: "The deployment option of the F5 LB instance," +
								" can be one of `ONE_BOOT_LOCATION` or `TWO_BOOT_LOCATION`." +
								" Only required when license_type is `BRING_YOUR_OWN`.",
							ValidateFunc: validation.StringInSlice(
								[]string{"ONE_BOOT_LOCATION", "TWO_BOOT_LOCATION"},
								false),
							Type:     schema.TypeString,
							Optional: true,
						},
						"deployment_type": {
							Description: "The deployment type used for the F5 load balancer instance." +
								" Can be one of `GOOD`, `BETTER`, `BEST`, `LTM_DNS` or `ALL`," +
								" deployment types `GOOD`, `BETTER` and `BEST`" +
								" are only applicable to license_type `PAY_AS_YOU_GO`" +
								" `LTM_DNS` and `ALL` are only applicable to license_type `BRING_YOUR_OWN`",
							Type: schema.TypeString,
							ValidateFunc: validation.StringInSlice(
								[]string{"GOOD", "BETTER", "BEST", "LTM_DNS", "ALL"},
								false),
							Required: true,
						},
						"hostname_fqdn": {
							Description: "The FQDN defined in route 53.",
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
	d.Set("instances", lb.Instances)
	d.Set("billing_tag_ids", lb.BillingTags)
	d.Set("global_cidr_list_id", lb.GlobalCidrListId)
	d.Set("prefix_list_id", lb.PrefixListId)

	segmentOptions, err := deflateF5SegmentOptions(lb.SegmentOptions, m)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("segment_options", segmentOptions)

	instance := setF5Instances(d, lb.Instances)
	d.Set("instances", instance)

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