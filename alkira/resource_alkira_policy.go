package alkira

import (
	"fmt"
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
			"description": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The description of the policy",

			},
			"enabled": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ASN of the segement",

			},
			"from_groups": &schema.Schema{
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the policy",
			},
			"rule_list_id": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_ids": &schema.Schema{
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
			},
			"to_groups": &schema.Schema{
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required:    true,
			},
		},
	}
}

func resourcePolicy(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	segmentIds := expandStringFromList(d.Get("segment_ids").([]interface{}))
	fromGroups := expandStringFromList(d.Get("from_groups").([]interface{}))
	toGroups   := expandStringFromList(d.Get("to_groups").([]interface{}))

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
	}

	d.SetId(strconv.Itoa(id))
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
	policyId  := d.Id()

	log.Printf("[INFO] Deleting Policy %s", policyId)
	err := client.DeletePolicy(policyId)

	if err != nil {
	 	return fmt.Errorf("failed to delete policy %s", policyId)
	}

	return nil
}
