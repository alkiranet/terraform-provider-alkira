package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourcePolicy,
		Read:   resourcePolicyRead,
		Update: resourcePolicyUpdate,
		Delete: resourcePolicyDelete,

		Schema: map[string]*schema.Schema{
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"from_groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_list_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"segment_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Required: true,
			},
			"to_groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Required: true,
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

	log.Printf("[INFO] Creating Policy")
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

	log.Printf("[INFO] Updating Policy")
	err = client.UpdatePolicy(d.Id(), request)

	if err != nil {
		log.Printf("[ERROR] Failed to update policy")
		return err
	}

	return resourcePolicyRead(d, m)
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Policy %s", d.Id())
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
