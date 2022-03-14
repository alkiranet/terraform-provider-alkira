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
				Description: "The BGP ASN for the segment. Default value is `65514`.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     "65514",
			},
			"cidr": {
				Description: "The CIDR block.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"reserve_public_ips": {
				Description: "Default value is `false`. When this is set to " +
					"`true`. Alkira reserves public IPs " +
					"which can be used to create underlay tunnels between an " +
					"external service and Alkira. For example the reserved public IPs " +
					"may be used to create tunnels to the Akamai Prolexic.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceSegment(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	segment, err := generateSegmentRequest(d)

	if err != nil {
		return err
	}

	id, err := client.CreateSegment(segment)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceSegmentRead(d, m)
}

func resourceSegmentRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	segment, err := client.GetSegmentById(d.Id())

	if err != nil {
		return err
	}

	d.Set("asn", segment.Asn)
	d.Set("cidr", segment.IpBlock)
	d.Set("name", segment.Name)
	d.Set("reserve_public_ips", segment.ReservePublicIPsForUserAndSiteConnectivity)

	return nil
}

func resourceSegmentUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	segment, err := generateSegmentRequest(d)

	if err != nil {
		return err
	}

	log.Printf("[INFO] Updateing Segment %s", d.Id())
	err = client.UpdateSegment(d.Id(), segment)

	return err
}

func resourceSegmentDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting Segment %s", d.Id())
	err := client.DeleteSegment(d.Id())

	return err
}

func generateSegmentRequest(d *schema.ResourceData) (*alkira.Segment, error) {
	seg := &alkira.Segment{
		Asn:     d.Get("asn").(int),
		IpBlock: d.Get("cidr").(string),
		Name:    d.Get("name").(string),
		ReservePublicIPsForUserAndSiteConnectivity: d.Get("reserve_public_ips").(bool),
	}

	return seg, nil
}
