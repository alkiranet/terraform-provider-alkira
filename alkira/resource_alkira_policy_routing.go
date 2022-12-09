package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPolicyRouting() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Routing Policy.\n\n" +
			"Configure a routing policy between the Alkira " +
			"CSX and a selected scope with custom rules",
		Create: resourcePolicyRouting,
		Read:   resourcePolicyRoutingRead,
		Update: resourcePolicyRoutingUpdate,
		Delete: resourcePolicyRoutingDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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
			"enabled": {
				Description: "Whether the routing policy is enabled. By default, it is set to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"direction": {
				Description:  "The direction of the route, `INBOUND` or `OUTBOUND`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"INBOUND", "OUTBOUND"}, false),
			},
			"segment_id": {
				Description: "IDs of segments that will define the policy scope.",
				Type:        schema.TypeInt,
				Required:    true,
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
				Description: "Advertise Alkira’s Internet Connector to selected scope. Default " +
					"value is `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"advertise_on_prem_routes": {
				Description: "Advertise routes from other on premise connectors to selected scope. " +
					"Default value is `false`.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"advertise_custom_routes_prefix_id": {
				Description: "Prefix list ID to send aggregates out towards on-prem connectors.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"rule": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Description: "The name of the rule.",
							Type:        schema.TypeString,
							Required:    true,
						},
						"action": {
							Description:  "Action to be set on matched routes. Value could be `ALLOW`, `DENY` and `ALLOW_W_SET`.",
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY", "ALLOW_W_SET"}, false),
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

func resourcePolicyRouting(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyRoutingRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate routing policy request.")
		return err
	}

	id, err := client.CreateRoutePolicy(request)

	if err != nil {
		log.Printf("[ERROR] Failed to create routing policy.")
		return err
	}

	d.SetId(id)
	return resourcePolicyRoutingRead(d, m)
}

func resourcePolicyRoutingRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	policy, err := client.GetRoutePolicy(d.Id())

	if err != nil {
		log.Printf("[ERROR] Failed to read routing policy %s.", d.Id())
		return err
	}

	d.Set("advertise_custom_routes_prefix_id", policy.AdvertiseCustomRoutesPrefixId)
	d.Set("advertise_internet_exit", policy.AdvertiseInternetExit)
	d.Set("advertise_on_prem_routes", policy.AdvertiseOnPremRoutes)
	d.Set("description", policy.Description)
	d.Set("direction", policy.Direction)
	d.Set("enabled", policy.Enabled)
	d.Set("excluded_group_ids", policy.ExcludedGroups)
	d.Set("included_group_ids", policy.IncludedGroups)
	d.Set("name", policy.Name)

	segment, err := client.GetSegmentByName(policy.Segment)

	if err != nil {
		return err
	}
	d.Set("segment_id", segment.Id)

	var rules []map[string]interface{}

	for _, rule := range policy.Rules {
		i := map[string]interface{}{
			"name":      rule.Name,
			"action":    rule.Action,
			"match_all": rule.Match.All,
		}

		if rule.Match.AsPathListIds != nil {
			i["match_as_path_list_ids"] = rule.Match.AsPathListIds
		}
		if rule.Match.CommunityListIds != nil {
			i["match_community_list_ids"] = rule.Match.CommunityListIds
		}
		if rule.Match.ExtendedCommunityListIds != nil {
			i["match_extended_community_list_ids"] = rule.Match.ExtendedCommunityListIds
		}
		if rule.Match.PrefixListIds != nil {
			i["match_prefix_list_ids"] = rule.Match.PrefixListIds
		}
		if rule.Match.ConnectorGroupIds != nil {
			i["match_group_ids"] = rule.Match.ConnectorGroupIds
		}
		if rule.Match.Cxps != nil {
			i["cxps"] = rule.Match.Cxps
		}

		if rule.Set != nil {
			i["set_as_path_prepend"] = rule.Set.AsPathPrepend
			i["set_community"] = rule.Set.Community
			i["set_extended_community"] = rule.Set.ExtendedCommunity
		}

		if rule.InterCxpRoutesRedistribution != nil {
			i["routes_distribution_restricted_cxps"] = rule.InterCxpRoutesRedistribution.RestrictedCxps
			i["routes_distribution_type"] = rule.InterCxpRoutesRedistribution.DistributionType
			i["routes_distribution_as_secondary"] = rule.InterCxpRoutesRedistribution.RedistributeAsSecondary
		}

		rules = append(rules, i)
	}

	d.Set("rule", rules)

	return nil
}

func resourcePolicyRoutingUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyRoutingRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate routing policy request.")
		return err
	}

	err = client.UpdateRoutePolicy(d.Id(), request)

	if err != nil {
		log.Printf("[ERROR] Failed to update routing policy.")
		return err
	}

	return resourcePolicyRoutingRead(d, m)
}

func resourcePolicyRoutingDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	return client.DeleteRoutePolicy(d.Id())
}

func generatePolicyRoutingRequest(d *schema.ResourceData, m interface{}) (*alkira.RoutePolicy, error) {

	client := m.(*alkira.AlkiraClient)

	inGroups := convertTypeListToIntList(d.Get("included_group_ids").([]interface{}))
	exGroups := convertTypeListToIntList(d.Get("excluded_group_ids").([]interface{}))

	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by ID: %d", d.Get("segment_id"))
		return nil, err
	}

	rules, err := expandPolicyRoutingRule(d.Get("rule").(*schema.Set))

	if err != nil {
		log.Printf("[ERROR] failed to expand routing policy rules.")
		return nil, err
	}

	policy := &alkira.RoutePolicy{
		Name:                          d.Get("name").(string),
		Description:                   d.Get("description").(string),
		Direction:                     d.Get("direction").(string),
		Enabled:                       d.Get("enabled").(bool),
		Segment:                       segment.Name,
		IncludedGroups:                inGroups,
		ExcludedGroups:                exGroups,
		AdvertiseInternetExit:         d.Get("advertise_internet_exit").(bool),
		AdvertiseOnPremRoutes:         d.Get("advertise_on_prem_routes").(bool),
		AdvertiseCustomRoutesPrefixId: d.Get("advertise_custom_routes_prefix_id").(int),
		Rules:                         rules,
	}

	return policy, nil
}
