package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraSegment() *schema.Resource {
	return &schema.Resource{
		Description: "Provides segment resource.",
		Create:      resourceSegment,
		Read:        resourceSegmentRead,
		Update:      resourceSegmentUpdate,
		Delete:      resourceSegmentDelete,
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
			"description": {
				Description: "The description of the segment.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"provision_state": {
				Description: "The provisioning state of the resource.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"enable_ipv6_to_ipv4_translation": {
				Description: "Enable IPv6 to IPv4 translation in the " +
					"segment for internet application traffic. (**BETA**)",
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enterprise_dns_server_ip": {
				Description: "The IP of the DNS server used within the segment. This DNS server " +
					"may be used by the Alkira CXP to resolve the names of LDAP servers for example " +
					"which are configured on the Remote Access Connector. (**BETA**)",
				Type:     schema.TypeString,
				Optional: true,
			},
			"reserve_public_ips": {
				Description: "Default value is `false`. When this is set to " +
					"`true`. Alkira reserves public IPs " +
					"which can be used to create underlay tunnels between an " +
					"external service and Alkira. For example the reserved public IPs " +
					"may be used to create tunnels to the Akamai Prolexic. (**BETA**)",
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
	api := alkira.NewSegment(client)

	// Construct request
	segment, err := generateSegmentRequest(d)

	if err != nil {
		return err
	}

	// Send create request
	response, provisionState, err := api.Create(segment)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	d.SetId(string(response.Id))
	return resourceSegmentRead(d, m)
}

func resourceSegmentRead(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegment(client)

	// Get the resource
	segment, err := api.GetById(d.Id())

	if err != nil {
		return err
	}

	d.Set("asn", segment.Asn)
	d.Set("description", segment.Description)
	d.Set("enable_ipv6_to_ipv4_translation", segment.EnableIpv6ToIpv4Translation)
	d.Set("enterprise_dns_server_ip", segment.EnterpriseDNSServerIP)
	d.Set("name", segment.Name)
	d.Set("reserve_public_ips", segment.ReservePublicIPsForUserAndSiteConnectivity)

	if segment.SrcIpv4PoolList != nil && len(segment.SrcIpv4PoolList) > 0 {
		d.Set("src_ipv4_pool_start_ip", segment.SrcIpv4PoolList[0].StartIp)
		d.Set("src_ipv4_pool_end_ip", segment.SrcIpv4PoolList[0].EndIp)
	}

	setCidrsSegmentRead(d, segment)

	// Set provision state
	_, provisionState, err := api.GetByName(d.Get("name").(string))

	if client.Provision == true && provisionState != "" {
		d.Set("provision_state", provisionState)
	}

	return nil
}

func resourceSegmentUpdate(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegment(client)

	// Construct request
	segment, err := generateSegmentRequest(d)

	if err != nil {
		return err
	}

	// Send update request
	provisionState, err := api.Update(d.Id(), segment)

	if err != nil {
		return err
	}

	// Set provision state
	if client.Provision == true {
		d.Set("provision_state", provisionState)
	}

	return resourceSegmentRead(d, m)
}

func resourceSegmentDelete(d *schema.ResourceData, m interface{}) error {

	client := m.(*alkira.AlkiraClient)
	api := alkira.NewSegment(client)

	provisionState, err := api.Delete(d.Id())

	if err != nil {
		return err
	}

	if client.Provision == true && provisionState != "SUCCESS" {
		return fmt.Errorf("failed to delete segment %s, provision failed", d.Id())
	}

	d.SetId("")
	return nil
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
		Description:                 d.Get("description").(string),
		EnableIpv6ToIpv4Translation: d.Get("enable_ipv6_to_ipv4_translation").(bool),
		EnterpriseDNSServerIP:       d.Get("enterprise_dns_server_ip").(string),
		Name:                        d.Get("name").(string),
		ReservePublicIPsForUserAndSiteConnectivity: d.Get("reserve_public_ips").(bool),
		IpBlocks: alkira.SegmentIpBlocks{
			Values: cidrs,
		},
		SrcIpv4PoolList: srcIpv4PoolList,
	}

	return seg, nil
}

func setCidrsSegmentRead(d *schema.ResourceData, segment *alkira.Segment) {
	if segment.IpBlock == "" || stringInSlice(segment.IpBlock, segment.IpBlocks.Values) {
		d.Set("cidrs", segment.IpBlocks.Values)
	} else {
		// segment.IpBlocks.Values could be empty doesn't matter we just get an empty slice to append to
		c := segment.IpBlocks.Values[0:]
		c = append(c, segment.IpBlock)
		d.Set("cidrs", c)
	}
}
