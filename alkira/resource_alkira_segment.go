package alkira

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceAlkiraSegment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSegment,
		Read:   resourceSegmentRead,
		Update: resourceSegmentUpdate,
		Delete: resourceSegmentDelete,

		Schema: map[string]*schema.Schema{
			"address": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSegment(d *schema.ResourceData, m interface{}) error {
        return resourceSegmentRead(d, m)
}

func resourceSegmentRead(d *schema.ResourceData, m interface{}) error {
        return nil
}

func resourceSegmentUpdate(d *schema.ResourceData, m interface{}) error {
        return resourceSegmentRead(d, m)
}

func resourceSegmentDelete(d *schema.ResourceData, m interface{}) error {
        return nil
}
