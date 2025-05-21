package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPolicyRouting() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Routing Policy.\n\n" +
			"Configure a routing policy between the Alkira " +
			"CSX and a selected scope with custom rules",
		CreateContext: resourcePolicyRouting,
		ReadContext:   resourcePolicyRoutingRead,
		UpdateContext: resourcePolicyRoutingUpdate,
		DeleteContext: resourcePolicyRoutingDelete,
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
				Description: "The name of the routing policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the routing policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provision state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enabled": {
				Description: "Whether the routing policy is enabled. " +
					"By default, it is set to `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"direction": {
				Description: "The direction of the route, `INBOUND` " +
					"or `OUTBOUND`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"INBOUND", "OUTBOUND"}, false),
			},
			"segment_id": {
				Description: "IDs of segments that will define " +
					"the policy scope.",
				Type:     schema.TypeString,
				Required: true,
			},
			"included_group_ids": {
				Description: "Defines the scope for the policy. Connector associated " +
					"with group IDs metioned here is where this policy would be applied. " +
					"Group IDs that associated with branch/on-premise connectors can be " +
					"used here. These group should not contain any cloud connector.",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Required: true,
			},
			"excluded_group_ids": {
				Description: "Excludes given associated connector from `included_groups`. " +
					"Implicit group ID of a branch/on-premise connector for which a user " +
					"defined group is used in `included_groups` can be used here.",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			"advertise_internet_exit": {
				Description: "Advertise Alkira’s Internet Connector to selected " +
					"scope. This only applies to `OUTBOUND` policy. Default " +
					"value is `true`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"advertise_on_prem_routes": {
				Description: "Advertise routes from other on premise connectors to " +
					"selected scope. Default value is `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"advertise_custom_routes_prefix_id": {
				Description: "Prefix list ID to send aggregates out towards " +
					"on-prem connectors.",
				Type:     schema.TypeInt,
				Optional: true,
			},
			"enable_as_override": {
				Description: "Whether enable AS-override on associated " +
					"connectors. Default value is `true`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the rule.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"sequence_no": {
							Description: "System assigned number for each " +
								"rule starting with `1000`. It defines the " +
								"order of the rules.",
							Type:     schema.TypeInt,
							Computed: true,
						},
						"action": {
							Description: "Action to be set on matched " +
								"routes. Value could be `ALLOW`, " +
								"`DENY` and `ALLOW_W_SET`.",
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ALLOW", "DENY", "ALLOW_W_SET"}, false),
						},
						"match_all": {
							Description: "This acts as match all if enabled" +
								"and should be used as exlusive match option.",
							Type:     schema.TypeBool,
							Optional: true,
						},
						"match_as_path_list_ids": {
							Description: "IDs of a AS Path Lists.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"match_community_list_ids": {
							Description: "IDs of Community Lists.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"match_extended_community_list_ids": {
							Description: "IDs of Extended Community Lists.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"match_prefix_list_ids": {
							Description: "IDs of Prefix Lists.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"match_group_ids": {
							Description: "IDs of groups.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"match_segment_resource_ids": {
							Description: "IDs of segment resources.",
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Optional:    true,
						},
						"match_cxps": {
							Description: "List of CXPs.",
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
						},
						"set_as_path_prepend": {
							Description: "Allows to prepend one or more AS " +
								"numbers to the current AS PATH. Each AS number " +
								"can be a value from 0 through 65535. " +
								"Example - 100 100 100.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"set_community": {
							Description: "Allows to add one or more community " +
								"attributes to the existing communities on the " +
								"route. Community attribute is specified in this " +
								"format: `as-number:community-value`. as-number " +
								"and community-value can be a value from `0` through " +
								"`65535`. Example: `65512:20 65512:21`.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"set_as_path_replace_with_segment_asn": {
							Description: "ASNs that will be replaced with the local segment ASN. " +
								"Accepts a comma-separated string of ASNs or 'ALL'. Can be null. " +
								"This option can be applied only to USERS_AND_SITES connectors.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"set_extended_community": {
							Description: "Allows to add one or more extended " +
								"community attributes to the existing extended " +
								"communities on the route. Extended community " +
								"attribute is specified in this format: " +
								"`type:administrator:assigned-number`. Currently " +
								"only type origin(soo) is supported.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"routes_distribution_type": {
							Description: "Redistribute routes that match with " +
								"this rule match codition to. The value could be " +
								"`ALL`, `LOCAL_ONLY` and `RESTRICTED_CXPS`.",
							Type:     schema.TypeString,
							Optional: true,
						},
						"routes_distribution_restricted_cxps": {
							Description: "List of cxps to which routes " +
								"distribution is restricted. Should be used " +
								"only with distributionType `RESTRICTED_CXPS`.",
							Type:     schema.TypeList,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
						},
						"routes_distribution_as_secondary": {
							Description: "This allows to redistribute routes with " +
								"lower preference to the restrictedCxps. Hence, " +
								"this option can only be used with `RESTRICTED_CXPS` " +
								"distribution_type. Also only 1 CXP is allowed in " +
								"restricted_cxps, when this is set to `true`.",
							Type:     schema.TypeBool,
							Optional: true,
						},
					},
				},
			}, // rule
		},
	}
}

func resourcePolicyRouting(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewRoutePolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyRoutingRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send create request
	response, provState, err, provErr := api.Create(request)

	if err != nil {
		return diag.FromErr(err)
	}

	// Set provision state
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

	d.SetId(string(response.Id))
	return resourcePolicyRoutingRead(ctx, d, m)
}

func resourcePolicyRoutingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewRoutePolicy(m.(*alkira.AlkiraClient))

	policy, provState, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	if policy.AdvertiseInternetExit != nil {
		d.Set("advertise_internet_exit", *policy.AdvertiseInternetExit)
	}

	if policy.EnableASOverride != nil {
		d.Set("enable_as_override", *policy.EnableASOverride)
	}

	d.Set("advertise_custom_routes_prefix_id", policy.AdvertiseCustomRoutesPrefixId)
	d.Set("advertise_on_prem_routes", policy.AdvertiseOnPremRoutes)
	d.Set("description", policy.Description)
	d.Set("direction", policy.Direction)
	d.Set("enabled", policy.Enabled)
	d.Set("excluded_group_ids", policy.ExcludedGroups)
	d.Set("included_group_ids", policy.IncludedGroups)
	d.Set("name", policy.Name)

	//
	// Set segment
	//
	segmentId, err := getSegmentIdByName(policy.Segment, m)

	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("segment_id", segmentId)

	//
	// Set Rule
	//
	err = setPolicyRoutingRules(policy.Rules, d)

	if err != nil {
		return diag.FromErr(err)
	}

	//
	// Set provision state
	//
	if client.Provision == true && provState != "" {
		d.Set("provision_state", provState)
	}

	return nil
}

func resourcePolicyRoutingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewRoutePolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyRoutingRequest(d, m)

	if err != nil {
		return diag.FromErr(err)
	}

	// Send update request
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

	return resourcePolicyRoutingRead(ctx, d, m)
}

func resourcePolicyRoutingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewRoutePolicy(m.(*alkira.AlkiraClient))

	provState, err, provErr := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	if client.Provision == true && provState != "SUCCESS" {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "PROVISION (DELETE) FAILED",
			Detail:   fmt.Sprintf("%s", provErr),
		}}
	}

	d.SetId("")
	return nil
}

func generatePolicyRoutingRequest(d *schema.ResourceData, m interface{}) (*alkira.RoutePolicy, error) {

	//
	// Segment
	//
	segmentName, err := getSegmentNameById(d.Get("segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	//
	// Rule
	//
	rules, err := expandPolicyRoutingRule(d.Get("rule").([]interface{}))

	if err != nil {
		return nil, err
	}

	// This field could be only used when policy is "OUTBOUND".
	advertiseInternetExit := new(bool)

	if d.Get("direction").(string) == "INBOUND" {
		advertiseInternetExit = nil
	} else {
		*advertiseInternetExit = d.Get("advertise_internet_exit").(bool)
	}

	enableASOverride := new(bool)

	if d.Get("direction").(string) == "INBOUND" {
		enableASOverride = nil
	} else {
		*enableASOverride = d.Get("enable_as_override").(bool)
	}

	// Assemble request
	policy := &alkira.RoutePolicy{
		Name:                          d.Get("name").(string),
		Description:                   d.Get("description").(string),
		Direction:                     d.Get("direction").(string),
		Enabled:                       d.Get("enabled").(bool),
		Segment:                       segmentName,
		IncludedGroups:                convertTypeListToIntList(d.Get("included_group_ids").([]interface{})),
		ExcludedGroups:                convertTypeListToIntList(d.Get("excluded_group_ids").([]interface{})),
		AdvertiseInternetExit:         advertiseInternetExit,
		AdvertiseOnPremRoutes:         d.Get("advertise_on_prem_routes").(bool),
		EnableASOverride:              enableASOverride,
		AdvertiseCustomRoutesPrefixId: d.Get("advertise_custom_routes_prefix_id").(int),
		Rules:                         rules,
	}

	return policy, nil
}
