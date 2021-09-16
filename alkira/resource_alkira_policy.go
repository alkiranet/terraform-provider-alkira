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
				Type:     schema.TypeString,
				Required: true,
			},
			"from_groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"rule_list_id": {
				Type:     schema.TypeString,
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

func resourcePolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	segmentIds := convertTypeListToIntList(d.Get("segment_ids").([]interface{}))
	fromGroups := convertTypeListToIntList(d.Get("from_groups").([]interface{}))
	toGroups := convertTypeListToIntList(d.Get("to_groups").([]interface{}))

	policy := &alkira.Policy{
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(string),
		FromGroups:  fromGroups,
		Name:        d.Get("name").(string),
		RuleListId:  d.Get("rule_list_id").(int),
		SegmentIds:  segmentIds,
		ToGroups:    toGroups,
	}

	log.Printf("[INFO] Policy Creating")
	id, err := client.CreatePolicy(policy)

	if err != nil {
		log.Printf("[ERROR] Failed to create policy")
		return err
	}

	d.SetId(id)
	return resourcePolicyRead(d, meta)
}

func resourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourcePolicyRead(d, meta)
}

func resourcePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Policy %s", d.Id())
	return client.DeletePolicy(d.Id())
}
