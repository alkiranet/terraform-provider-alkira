package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraSegment() *schema.Resource {
	return &schema.Resource{
		Description: "The segment data source allows a segment to be retrieved by its name.",
		Read:        dataSourceAlkiraSegmentRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the segment.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_id": {
				Description: "The ID of the segment.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func dataSourceAlkiraSegmentRead(d *schema.ResourceData, m interface{}) error {
	api := alkira.NewSegment(m.(*alkira.AlkiraClient))

	segment, _, err := api.GetByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(string(segment.Id))
	d.Set("segment_id", segment.Id)

	return nil
}
