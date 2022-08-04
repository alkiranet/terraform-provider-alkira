package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPolicyRoute() *schema.Resource {
	return &schema.Resource{
		Description: "Manage Route Policy.",
		Create:      resourcePolicyRoute,
		Read:        resourcePolicyRouteRead,
		Update:      resourcePolicyRouteUpdate,
		Delete:      resourcePolicyRouteDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the route policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the route policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Is the route policy enabled.",
				Type:        schema.TypeBool,
				Optional:    true,
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
				Description: "Defines the scope for the policy. Connector associated" +
					"with group IDs metioned here is where this policy would be applied." +
					"Group IDs that associated with branch/on-premise connectors can be" +
					"used here. These group should not contain any cloud connector.",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Required: true,
			},
			"excluded_group_ids": {
				Description: "Excludes given associated connector from `included_groups`." +
					"Implicit group ID of a branch/on-premise connector for which a user" +
					"defined group is used in `included_groups` can be used here.",
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
			"advertise_internet_exit": {
				Description: "Advertise Alkira’s Internet Connector to selected scope. Default value is `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"advertise_on_prem_routes": {
				Description: "Advertise routes from other on premise connectors to selected scope. Default value is `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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
							Description: "Action to be set on matched routes. Value could be `ALLOW`, `DENY` and `ALLOW_W_SET`.",
							Type:        schema.TypeString,
							Required:    true,
							ValidateFunc: validation.StringInSlice([]string{"ALLOW", "DENY", "ALLOW_W_SET"}, false),
						},
						"match": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"all": {
										Description: "This acts as match all if enabled" +
											"and should be used as exlusive match option.",
										Type:     schema.TypeBool,
										Required: true,
									},
									"as_path_list_ids": {
										Description: "IDs of a AS Path Lists.",
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Optional:    true,
									},
									"community_list_ids": {
										Description: "IDs of Community Lists.",
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Optional:    true,
									},
									"extended_community_list_ids": {
										Description: "IDs of Extended Community Lists.",
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Optional:    true,
									},
									"prefix_list_ids": {
										Description: "IDs of Prefix Lists.",
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Optional:    true,
									},
									"group_ids": {
										Description: "IDs of groups.",
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeInt},
										Optional:    true,
									},
									"cxps": {
										Description: "List of CXPs.",
										Type:        schema.TypeList,
										Elem:        &schema.Schema{Type: schema.TypeString},
										Optional:    true,
									},
								},
							},
						}, // match
						"set": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"as_path_prepend": {
										Description: "Allows to prepend one or more AS " +
											"numbers to the current AS PATH. Each AS number " +
											"can be a value from 0 through 65535. " +
											"Example - 100 100 100.",
										Type:     schema.TypeString,
										Required: true,
									},
									"community": {
										Description: "Allows to add one or more community " +
											"attributes to the existing communities on the " +
											"route. Community attribute is specified in this " +
											"format: `as-number:community-value`. as-number " +
											"and community-value can be a value from `0` through " +
											"`65535`. Example: `65512:20 65512:21`.",
										Type:     schema.TypeString,
										Required: true,
									},
									"extended_community": {
										Description: "Allows to add one or more extended " +
											"community attributes to the existing extended " +
											"communities on the route. Extended community " +
											"attribute is specified in this format: " +
											"`type:administrator:assigned-number`. Currently " +
											"only type origin(soo) is supported.",
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						}, // set
						"inter_cxp_routes_redistribution": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"distribution_type": {
										Description: "Redistribute routes that match with " +
											"this rule match codition to. The value could be " +
											"`ALL`, `LOCAL_ONLY` and `RESTRICTED_CXPS`.",
										Type:     schema.TypeString,
										Required: true,
									},
									"restricted_cxps": {
										Description: "List of cxps to which routes" +
											"distribution is restricted. Should be used" +
											"only with distributionType `RESTRICTED_CXPS`.",
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
										Optional: true,
									},
									"redistribute_as_secondary": {
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
						}, // inter_cxp_routes_redistribution
					},
				},
			}, // rule
		},
	}
}

func resourcePolicyRoute(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyRouteRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy")
		return err
	}

	id, err := client.CreateRoutePolicy(request)

	if err != nil {
		log.Printf("[ERROR] Failed to create route policy")
		return err
	}

	d.SetId(id)
	return resourcePolicyRouteRead(d, m)
}

func resourcePolicyRouteRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	policy, err := client.GetRoutePolicy(d.Id())

	if err != nil {
		log.Printf("[ERROR] Failed to read policy %s", d.Id())
		return err
	}

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("included_group_ids", policy.IncludedGroups)
	d.Set("excluded_group_ids", policy.ExcludedGroups)

	segment, err := client.GetSegmentByName(policy.Segment)

	if err != nil {
		return err
	}
	d.Set("segment_id", segment.Id)

	return nil
}

func resourcePolicyRouteUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyRouteRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy")
		return err
	}

	err = client.UpdateRoutePolicy(d.Id(), request)

	if err != nil {
		log.Printf("[ERROR] Failed to update route policy")
		return err
	}

	return resourcePolicyRouteRead(d, m)
}

func resourcePolicyRouteDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	return client.DeleteRoutePolicy(d.Id())
}

func generatePolicyRouteRequest(d *schema.ResourceData, m interface{}) (*alkira.RoutePolicy, error) {

	client := m.(*alkira.AlkiraClient)

	inGroups := convertTypeListToIntList(d.Get("included_group_ids").([]interface{}))
	exGroups := convertTypeListToIntList(d.Get("excluded_group_ids").([]interface{}))

	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	policy := &alkira.RoutePolicy{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Segment:        segment.Name,
		IncludedGroups: inGroups,
		ExcludedGroups: exGroups,
	}

	return policy, nil
}
