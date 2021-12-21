package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyNat() *schema.Resource {
	return &schema.Resource{
		Description: "Manage NAT policy.",
		Create:      resourcePolicyNat,
		Read:        resourcePolicyNatRead,
		Update:      resourcePolicyNatUpdate,
		Delete:      resourcePolicyNatDelete,

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
				Description: "The type of NAT policy, currently only `INTRA_SEGMENT`is supported.",
				Type:        schema.TypeBool,
				Required:    true,
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
				Required: true,
			},
			"nat_rule_ids": {
				Description: "The list of NAT rules to be applied by the policy.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
		},
	}
}

func resourcePolicyNat(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyNatRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy")
		return err
	}

	id, err := client.CreateNatPolicy(request)

	if err != nil {
		log.Printf("[ERROR] Failed to create policy")
		return err
	}

	d.SetId(id)
	return resourcePolicyNatRead(d, m)
}

func resourcePolicyNatRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	policy, err := client.GetNatPolicy(d.Id())

	if err != nil {
		log.Printf("[ERROR] Failed to read policy %s", d.Id())
		return err
	}

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("type", policy.Type)
	d.Set("included_group_ids", policy.IncludedGroups)
	d.Set("excluded_group_ids", policy.ExcludedGroups)
	d.Set("nat_rule_ids", policy.NatRuleIds)

	segment, err := client.GetSegmentByName(policy.Segment)

	if err != nil {
		return err
	}
	d.Set("segment_id", segment.Id)

	return nil
}

func resourcePolicyNatUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyNatRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy")
		return err
	}

	err = client.UpdateNatPolicy(d.Id(), request)

	if err != nil {
		log.Printf("[ERROR] Failed to update policy")
		return err
	}

	return resourcePolicyNatRead(d, m)
}

func resourcePolicyNatDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	return client.DeleteNatPolicy(d.Id())
}

func generatePolicyNatRequest(d *schema.ResourceData, m interface{}) (*alkira.NatPolicy, error) {

	client := m.(*alkira.AlkiraClient)

	inGroups := convertTypeListToIntList(d.Get("include_group_ids").([]interface{}))
	exGroups := convertTypeListToIntList(d.Get("exclude_group_ids").([]interface{}))
	natRules := convertTypeListToIntList(d.Get("nat_rule_ids").([]interface{}))

	segment, err := client.GetSegmentById(strconv.Itoa(d.Get("segment_id").(int)))

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
	}

	return policy, nil
}
