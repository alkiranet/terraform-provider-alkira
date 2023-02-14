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
			"provision_state": {
				Description: "The provision state of the policy.",
				Type:        schema.TypeString,
				Computed:    true,
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
	api := alkira.NewTrafficPolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyRequest(d, m)

	if err != nil {
		return err
	}

	// Send request
	resource, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	d.Set("provision_state", provisionState)
	d.SetId(string(resource.Id))

	return resourcePolicyRead(d, m)
}

func resourcePolicyRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewTrafficPolicy(m.(*alkira.AlkiraClient))

	policy, err := api.GetById(d.Id())

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
	api := alkira.NewTrafficPolicy(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	d.Set("provision_state", provisionState)
	return resourcePolicyRead(d, m)
}

func resourcePolicyDelete(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewTrafficPolicy(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if provisionState != "SUCCESS" {
	}

	d.SetId("")
	return nil
}

func generatePolicyRequest(d *schema.ResourceData, m interface{}) (*alkira.TrafficPolicy, error) {

	segmentIds := convertTypeListToIntList(d.Get("segment_ids").([]interface{}))
	fromGroups := convertTypeListToIntList(d.Get("from_groups").([]interface{}))
	toGroups := convertTypeListToIntList(d.Get("to_groups").([]interface{}))

	policy := &alkira.TrafficPolicy{
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
