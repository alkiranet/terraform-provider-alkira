package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicy() *schema.Resource {
	return &schema.Resource{
		Description: "Manage policy.",
		Create:      resourcePolicy,
		Read:        resourcePolicyRead,
		Update:      resourcePolicyUpdate,
		Delete:      resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": {
				Description: "The description of the policy.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"enabled": {
				Description: "Whether the policy is enabled.",
				Type:        schema.TypeBool,
				Required:    true,
			},
			"from_groups": {
				Description: "IDs of groups that will define source in the policy scope",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"name": {
				Description: "The name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"rule_list_id": {
				Description: "The `rulelist` that will be used by the policy.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"segment_ids": {
				Description: "IDs of segments that will define the policy scope.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
			"to_groups": {
				Description: "IDs of groups that will define destination in the policy scope.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Required:    true,
			},
		},
	}
}

func resourcePolicy(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy")
		return err
	}

	id, err := client.CreatePolicy(request)

	if err != nil {
		log.Printf("[ERROR] Failed to create policy")
		return err
	}

	d.SetId(id)
	return resourcePolicyRead(d, m)
}

func resourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	policy, err := client.GetPolicy(d.Id())

	if err != nil {
		log.Printf("[ERROR] Failed to read policy %s", d.Id())
		return err
	}

	d.Set("description", policy.Description)
	d.Set("enabled", policy.Enabled)
	d.Set("name", policy.Name)
	d.Set("rule_list_id", policy.RuleListId)
	d.Set("segment_ids", policy.SegmentIds)
	d.Set("from_groups", policy.FromGroups)
	d.Set("to_groups", policy.ToGroups)

	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	request, err := generatePolicyRequest(d, m)

	if err != nil {
		log.Printf("[ERROR] Failed to generate policy")
		return err
	}

	err = client.UpdatePolicy(d.Id(), request)

	if err != nil {
		log.Printf("[ERROR] Failed to update policy")
		return err
	}

	return resourcePolicyRead(d, m)
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	return client.DeletePolicy(d.Id())
}

func generatePolicyRequest(d *schema.ResourceData, m interface{}) (*alkira.Policy, error) {

	segmentIds := convertTypeListToIntList(d.Get("segment_ids").([]interface{}))
	fromGroups := convertTypeListToIntList(d.Get("from_groups").([]interface{}))
	toGroups := convertTypeListToIntList(d.Get("to_groups").([]interface{}))

	policy := &alkira.Policy{
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(bool),
		FromGroups:  fromGroups,
		Name:        d.Get("name").(string),
		RuleListId:  d.Get("rule_list_id").(int),
		SegmentIds:  segmentIds,
		ToGroups:    toGroups,
	}

	return policy, nil
}
