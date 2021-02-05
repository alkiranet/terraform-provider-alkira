package alkira

import (
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraSegment() *schema.Resource {
	return &schema.Resource{
		Description: "This resource manages segments.\n\n" +
			"A Segment is a section of a network isolated from one another" +
			"to make it possible to more effectively control who has access" +
			"to what. Segmentation also allows for segregation of resources" +
			"between segments for security and isolation purposes.",

		Create: resourceSegment,
		Read:   resourceSegmentRead,
		Update: resourceSegmentUpdate,
		Delete: resourceSegmentDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the segment.",
				Type:     schema.TypeString,
				Required: true,
			},
			"asn": {
				Description: "The BGP ASN for the segment.",
				Type:     schema.TypeString,
				Required: true,
			},
			"cidr": {
				Description: "The CIDR block.",
				Type:     schema.TypeString,
				Required: true,
			},
			"segment_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceSegment(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	name := d.Get("name").(string)

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
	client := meta.(*alkira.AlkiraClient)
	segmentId := d.Get("segment_id").(int)

	log.Printf("[INFO] Updateing Segment %d", segmentId)
	err := client.UpdateSegment(segmentId, d.Get("name").(string), d.Get("asn").(string), d.Get("cidr").(string))

	if err != nil {
		return err
	}

	return nil
}

func resourceSegmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)
	segmentId := d.Get("segment_id").(int)

	log.Printf("[INFO] Deleting Segment %d", segmentId)
	err := client.DeleteSegment(segmentId)

	if err != nil {
		return err
	}

	return nil
}
