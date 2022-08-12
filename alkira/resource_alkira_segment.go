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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"cidrs": {
				Description: "The list of CIDR blocks.",
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"enable_ipv6_to_ipv4_translation": {
				Description: "Enable IPv6 to IPv4 translation in the " +
					"segment for internet application traffic.",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
			"src_ipv4_pool_start_ip": {
				Description: "The start IP address of IPv4 pool.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"src_ipv4_pool_end_ip": {
				Description: "The end IP address of IPv4 pool.",
				Type:        schema.TypeString,
				Optional:    true,
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
	d.Set("enable_ipv6_to_ipv4_translation", segment.EnableIpv6ToIpv4Translation)
	d.Set("name", segment.Name)
	d.Set("reserve_public_ips", segment.ReservePublicIPsForUserAndSiteConnectivity)

	if segment.SrcIpv4PoolList != nil && len(segment.SrcIpv4PoolList) > 0 {
		d.Set("src_ipv4_pool_start_ip", segment.SrcIpv4PoolList[0].StartIp)
		d.Set("src_ipv4_pool_end_ip", segment.SrcIpv4PoolList[0].EndIp)
	}

	setCidrsSegmentRead(d, segment)

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
	cidrs := convertTypeListToStringList(d.Get("cidrs").([]interface{}))

	// Special handle for pool list, otherwise, request will simply fail
	srcIpv4PoolList := []alkira.SegmentSrcIpv4PoolList{}
	list := alkira.SegmentSrcIpv4PoolList{}

	if d.Get("src_ipv4_pool_start_ip") != "" && d.Get("src_ipv4_pool_end_ip") != "" {
		list.StartIp = d.Get("src_ipv4_pool_start_ip").(string)
		list.EndIp = d.Get("src_ipv4_pool_end_ip").(string)

		srcIpv4PoolList = []alkira.SegmentSrcIpv4PoolList{list}
	} else {
		srcIpv4PoolList = nil
	}

	seg := &alkira.Segment{
		Asn:                         d.Get("asn").(int),
		EnableIpv6ToIpv4Translation: d.Get("enable_ipv6_to_ipv4_translation").(bool),
		Name:                        d.Get("name").(string),
		ReservePublicIPsForUserAndSiteConnectivity: d.Get("reserve_public_ips").(bool),
		IpBlocks: alkira.SegmentIpBlocks{
			Values: cidrs,
		},
		SrcIpv4PoolList: srcIpv4PoolList,
	}

	return seg, nil
}

func setCidrsSegmentRead(d *schema.ResourceData, segment alkira.Segment) {
	if segment.IpBlock == "" || stringInSlice(segment.IpBlock, segment.IpBlocks.Values) {
		d.Set("cidrs", segment.IpBlocks.Values)
	} else {
		// segment.IpBlocks.Values could be empty doesn't matter we just get an empty slice to append to
		c := segment.IpBlocks.Values[0:]
		c = append(c, segment.IpBlock)
		d.Set("cidrs", c)
	}
}
