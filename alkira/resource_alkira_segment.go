package alkira

import (
	"log"

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
				Type:        schema.TypeString,
				Required:    true,
			},
			"asn": {
				Description: "The BGP ASN for the segment.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cidr": {
				Description: "The CIDR block.",
				Type:        schema.TypeString,
				Required:    true,
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

	d.SetId(id)

	return resourceSegmentRead(d, meta)
}

func resourceSegmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	segment, err := client.GetSegmentById(d.Id())

	if err != nil {
		return err
	}

	d.Set("asn", segment.Asn)
	d.Set("cidr", segment.IpBlock)
	d.Set("name", segment.Name)

	return nil
}

func resourceSegmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Updateing Segment %s", d.Id())
	err := client.UpdateSegment(d.Id(), d.Get("name").(string), d.Get("asn").(string), d.Get("cidr").(string))

	return err
}

func resourceSegmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Segment %s", d.Id())
	err := client.DeleteSegment(d.Id())

	return err
}
