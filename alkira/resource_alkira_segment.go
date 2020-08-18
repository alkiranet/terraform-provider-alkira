package alkira

import (
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/alkira-client-go/alkira"
)

func resourceAlkiraSegment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSegment,
		Read:   resourceSegmentRead,
		Update: resourceSegmentUpdate,
		Delete: resourceSegmentDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
			},
			"asn": {
				Type:        schema.TypeString,
				Required:    true,
			},
			"cidr": {
				Type:        schema.TypeString,
				Required:    true,
			},
			"segment_id": {
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceSegment(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	name   := d.Get("name").(string)

	log.Printf("[INFO] Segment Creating")
	id, err := client.CreateSegment(name, d.Get("asn").(string), d.Get("cidr").(string))
	log.Printf("[INFO] Segment ID: %d", id)

	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))
	d.Set("segment_id", id)

	return resourceSegmentRead(d, meta)
}

func resourceSegmentRead(d *schema.ResourceData, meta interface{}) error {
        return nil
}

func resourceSegmentUpdate(d *schema.ResourceData, meta interface{}) error {
        return resourceSegmentRead(d, meta)
}

func resourceSegmentDelete(d *schema.ResourceData, meta interface{}) error {
	client    := meta.(*alkira.AlkiraClient)
	segmentId := d.Get("segment_id").(int)

	log.Printf("[INFO] Deleting Segment %d", segmentId)
	err := client.DeleteSegment(segmentId)

	if err != nil {
	 	return err
	}

	return nil
}
