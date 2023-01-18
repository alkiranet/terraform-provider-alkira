package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraPolicyNat() *schema.Resource {
	return &schema.Resource{
		Description: "Manage NAT policy.",
		Create:      resourcePolicyNat,
		Read:        resourcePolicyNatRead,
		Update:      resourcePolicyNatUpdate,
		Delete:      resourcePolicyNatDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"type": {
				Description:  "The type of NAT policy, currently only `INTRA_SEGMENT`is supported.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"INTRA_SEGMENT"}, false),
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
			"nat_rule_ids": {
				Description: "The list of NAT rules to be applied by the policy.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"category": {
				Description:  "The category of NAT policy, options are `DEFAULT` or `INTERNET_CONNECTOR`.",
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"DEFAULT", "INTERNET_CONNECTOR"}, false),
			},
		},
	}
}

func resourcePolicyNat(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewNatPolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyNatRequest(d, m)

	if err != nil {
		return err
	}

	// Send request
	resource, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	d.SetId(string(resource.Id))
	d.Set("provision_state", provisionState)

	return resourcePolicyNatRead(d, m)
}

func resourcePolicyNatRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewNatPolicy(m.(*alkira.AlkiraClient))

	policy, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("type", policy.Type)
	d.Set("included_group_ids", policy.IncludedGroups)
	d.Set("excluded_group_ids", policy.ExcludedGroups)
	d.Set("nat_rule_ids", policy.NatRuleIds)

	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, _, err := segmentApi.GetByName(policy.Segment)

	if err != nil {
		return err
	}

	d.Set("segment_id", segment.Id)

	return nil
}

func resourcePolicyNatUpdate(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewNatPolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyNatRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	d.Set("provision_state", provisionState)

	return resourcePolicyNatRead(d, m)
}

func resourcePolicyNatDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewNatPolicy(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if provisionState != "SUCCESS" {
	}

	d.SetId("")
	return nil
}

func generatePolicyNatRequest(d *schema.ResourceData, m interface{}) (*alkira.NatPolicy, error) {

	inGroups := convertTypeListToIntList(d.Get("included_group_ids").([]interface{}))
	exGroups := convertTypeListToIntList(d.Get("excluded_group_ids").([]interface{}))
	natRules := convertTypeListToIntList(d.Get("nat_rule_ids").([]interface{}))

	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, err := segmentApi.GetById(strconv.Itoa(d.Get("segment_id").(int)))

	if err != nil {
		log.Printf("[ERROR] failed to get segment by Id: %d", d.Get("segment_id"))
		return nil, err
	}

	policy := &alkira.NatPolicy{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Type:           d.Get("type").(string),
		Segment:        segment.Name,
		IncludedGroups: inGroups,
		ExcludedGroups: exGroups,
		NatRuleIds:     natRules,
		Category:       d.Get("category").(string),
	}

	return policy, nil
}
