package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraPolicyPrefixList() *schema.Resource {
	return &schema.Resource{
		Description: "Manage policy prefix list.",
		Create:      resourcePolicyPrefixList,
		Read:        resourcePolicyPrefixListRead,
		Update:      resourcePolicyPrefixListUpdate,
		Delete:      resourcePolicyPrefixListDelete,
		CustomizeDiff: func(ctx context.Context, d *schema.ResourceDiff, m interface{}) error {
			client := m.(*alkira.AlkiraClient)

			old, _ := d.GetChange("provision_state")

			if client.Provision == true && old == "FAILED" {
				d.SetNew("provision_state", "SUCCESS")
			}

			return nil
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the prefix list.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the prefix list.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"prefixes": {
				Description: "A list of prefixes.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"prefix_range": {
				Description: "A valid prefix range that could be used to " +
					"define a prefix of type `ROUTE`.",
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"prefix": {
							Description: "A valid CIDR as prefix in " +
								"`x.x.x.x/m` format.",
							Type:     schema.TypeString,
							Required: true,
						},
						"le": {
							Description: "Integer less than `32` but " +
								"greater than mask `m` in prefix",
							Type:     schema.TypeInt,
							Optional: true,
						},
						"ge": {
							Description: "Integer less than `32` but " +
								"greater than mask `m` in prefix and less than `le`.",
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourcePolicyPrefixList(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyPrefixListRequest(d, m)

	if err != nil {
		return err
	}

	// Send request
	response, provisionState, err := api.Create(request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	d.SetId(string(response.Id))
	return resourcePolicyPrefixListRead(d, m)
}

func resourcePolicyPrefixListRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	list, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", list.Name)
	d.Set("description", list.Description)
	d.Set("prefixes", list.Prefixes)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourcePolicyPrefixListUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	// Construct request
	request, err := generatePolicyPrefixListRequest(d, m)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), request)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourcePolicyPrefixListRead(d, m)
}

func resourcePolicyPrefixListDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewPolicyPrefixList(m.(*alkira.AlkiraClient))

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete policy_prefix_list %s, provision failed", d.Id())
	}

	d.SetId("")
	return nil
}

func generatePolicyPrefixListRequest(d *schema.ResourceData, m interface{}) (*alkira.PolicyPrefixList, error) {

	prefixRanges, err := expandPrefixListPrefixRanges(d.Get("prefix_range").([]interface{}))

	if err != nil {
		return nil, err
	}

	list := &alkira.PolicyPrefixList{
		Description:  d.Get("description").(string),
		Name:         d.Get("name").(string),
		Prefixes:     convertTypeListToStringList(d.Get("prefixes").([]interface{})),
		PrefixRanges: prefixRanges,
	}

	return list, nil
}
