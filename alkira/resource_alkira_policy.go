package alkira

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
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
			"policy_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rule_list_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"to_groups": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourcePolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	segmentIds := convertTypeListToStringList(d.Get("segment_ids").([]interface{}))
	fromGroups := convertTypeListToStringList(d.Get("from_groups").([]interface{}))
	toGroups   := convertTypeListToStringList(d.Get("to_groups").([]interface{}))

	policy := &alkira.PolicyRequest{
		Description: d.Get("description").(string),
		Enabled:     d.Get("enabled").(string),
		FromGroups:  fromGroups,
		Name:        d.Get("name").(string),
		RuleListId:  d.Get("rule_list_id").(string),
		SegmentIds:  segmentIds,
		ToGroups:    toGroups,
	}

	log.Printf("[INFO] Policy Creating")
	id, err := client.CreatePolicy(policy)
	log.Printf("[INFO] Policy ID: %d", id)

	if id == 0 || err != nil {
		log.Printf("[ERROR] Failed to create policy")
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("policy_id", id)

	return resourcePolicyRead(d, meta)
}

func resourcePolicyRead(d *schema.ResourceData, meta interface{}) error {
        return nil
}

func resourcePolicyUpdate(d *schema.ResourceData, meta interface{}) error {
        return resourcePolicyRead(d, meta)
}

func resourcePolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client    := meta.(*alkira.AlkiraClient)
	policyId  := d.Get("policy_id").(int)

	log.Printf("[INFO] Deleting Policy %s", policyId)
	err := client.DeletePolicy(policyId)

	if err != nil {
	 	return err
	}

	return nil
}
