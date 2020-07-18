package alkira

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/alkiranet/terraform-provider-alkira/alkira/internal"
)

func resourceAlkiraSegment() *schema.Resource {
	return &schema.Resource{
		Create: resourceSegment,
		Read:   resourceSegmentRead,
		Update: resourceSegmentUpdate,
		Delete: resourceSegmentDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the segment",
			},
			"asn": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ASN of the segement",

			},
			"cidr": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "The CIDR Block of the segment",
			},
			"segment_id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the segment",
			},
		},
	}
}

func resourceSegment(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*internal.AlkiraClient)
	name   := d.Get("name").(string)

	log.Printf("[INFO] Segment Creating")
	id, statusCode := client.CreateSegment(name, d.Get("asn").(string), d.Get("cidr").(string))
	log.Printf("[INFO] Segment ID: %d", id)

	if statusCode != 200 {
		fmt.Printf("ERROR: failed to create segment")
	}

	d.SetId(strconv.Itoa(id))
	d.Set("segment_id", strconv.Itoa(id))
	return resourceSegmentRead(d, meta)
}

func resourceSegmentRead(d *schema.ResourceData, meta interface{}) error {
        return nil
}

func resourceSegmentUpdate(d *schema.ResourceData, meta interface{}) error {
        return resourceSegmentRead(d, meta)
}

func resourceSegmentDelete(d *schema.ResourceData, meta interface{}) error {
	client    := meta.(*internal.AlkiraClient)
	segmentId := d.Id()

	log.Printf("[INFO] Deleting Segment %s", segmentId)
	statusCode := client.DeleteSegment(segmentId)

	if statusCode != 202 {
	 	return fmt.Errorf("failed to delete segment %s", segmentId)
	}

	return nil
}
