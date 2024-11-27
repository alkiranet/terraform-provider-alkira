package alkira

import (
	"errors"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getAwsDirectConnectSegmentOptions set "segment_options" block from API response
func getAwsDirectConnectSegmentOptions(instance alkira.ConnectorAwsDirectConnectInstance, m interface{}) ([]map[string]interface{}, error) {

	segmentOptions := instance.SegmentOptions

	if segmentOptions == nil {
		return nil, errors.New("invalid \"segment_options{}\"")
	}

	var segmentOptionBlocks []map[string]interface{}

	for _, option := range segmentOptions {

		segmentId, err := getSegmentIdByName(option.SegmentName, m)
		if err != nil {
			return nil, err
		}

		segmentOption := map[string]interface{}{
			"segment_id":                            segmentId,
			"on_prem_segment_asn":                   option.CustomerAsn,
			"customer_loopback_ip":                  option.CustomerLoopbackIp,
			"alkira_loopback_ip1":                   option.AlkLoopbackIp1,
			"alkira_loopback_ip2":                   option.AlkLoopbackIp2,
			"loopback_subnet":                       option.LoopbackSubnet,
			"advertise_on_prem_routes":              option.AdvertiseOnPremRoutes,
			"advertise_default_routes":              !option.DisableInternetExit,
			"number_of_customer_loopback_ips":       option.NumOfCustomerLoopbackIps,
			"tunnel_count_per_customer_loopback_ip": option.TunnelCountPerCustomerLoopbackIp,
		}
		segmentOptionBlocks = append(segmentOptionBlocks, segmentOption)
	}

	return segmentOptionBlocks, nil
}

func setAwsDirectConnectInstance(d *schema.ResourceData, m interface{}, connector *alkira.ConnectorAwsDirectConnect) error {
	var instances []map[string]interface{}

	for _, ins := range connector.Instances {

		// Firstly, get segment_options of the instance
		segmentOptions, err := getAwsDirectConnectSegmentOptions(ins, m)

		if err != nil {
			return err
		}

		instance := map[string]interface{}{
			"name":                ins.Name,
			"id":                  ins.Id,
			"connection_id":       ins.ConnectionId,
			"dx_asn":              ins.DcGatewayAsn,
			"dx_gateway_ip":       ins.AwsUnderlayIp,
			"on_prem_asn":         ins.UnderlayAsn,
			"on_prem_gateway_ip":  ins.OnPremUnderlayIp,
			"underlay_prefix":     ins.UnderlayPrefix,
			"bgp_auth_key":        ins.BgpAuthKey,
			"bgp_auth_key_alkira": ins.BgpAuthKeyAlkira,
			"vlan_id":             ins.Vlan,
			"aws_region":          ins.CustomerRegion,
			"credential_id":       ins.CredentialId,
			"gateway_mac_address": ins.GatewayMacAddress,
			"segment_options":     segmentOptions,
		}
		instances = append(instances, instance)
	}

	d.Set("instances", instances)
	return nil
}

// expandAwsDirectConnectSegmentOptions expand "segment_options" block
// in "instance" block to generate request payload.
func expandAwsDirectConnectSegmentOptions(in *schema.Set, m interface{}) ([]alkira.ConnectorAwsDirectConnectSegmentOption, error) {

	if in == nil || in.Len() == 0 {
		return nil, errors.New("[ERROR] invalid connector_aws_directconnect segment_options.")
	}

	segmentOptions := make([]alkira.ConnectorAwsDirectConnectSegmentOption, in.Len())

	for i, block := range in.List() {
		cfg := block.(map[string]interface{})
		option := alkira.ConnectorAwsDirectConnectSegmentOption{}

		if v, ok := cfg["segment_id"].(string); ok {
			segmentName, err := getSegmentNameById(v, m)
			if err != nil {
				return nil, err
			}
			option.SegmentName = segmentName
		}
		if v, ok := cfg["on_prem_segment_asn"].(int); ok {
			option.CustomerAsn = v
		}
		if v, ok := cfg["customer_loopback_ip"].(string); ok {
			option.CustomerLoopbackIp = v
		}
		if v, ok := cfg["alkira_loopback_ip1"].(string); ok {
			option.AlkLoopbackIp1 = v
		}
		if v, ok := cfg["alkira_loopback_ip2"].(string); ok {
			option.AlkLoopbackIp2 = v
		}
		if v, ok := cfg["loopback_subnet"].(string); ok {
			option.LoopbackSubnet = v
		}
		if v, ok := cfg["advertise_on_prem_routes"].(bool); ok {
			option.AdvertiseOnPremRoutes = v
		}
		if v, ok := cfg["advertise_default_routes"].(bool); ok {
			option.DisableInternetExit = !v
		}
		if v, ok := cfg["number_of_customer_loopback_ips"].(int); ok {
			option.NumOfCustomerLoopbackIps = v
		}
		if v, ok := cfg["tunnel_count_per_customer_loopback_ip"].(int); ok {
			option.TunnelCountPerCustomerLoopbackIp = v
		}

		segmentOptions[i] = option
	}

	return segmentOptions, nil
}

// expandAwsDirectConnectInstances expand instance block to generate request payload
func expandAwsDirectConnectInstances(in []interface{}, m interface{}) ([]alkira.ConnectorAwsDirectConnectInstance, error) {

	if in == nil || len(in) == 0 {
		return nil, errors.New("[ERROR] Invalid AWS DX instance input")
	}

	instances := make([]alkira.ConnectorAwsDirectConnectInstance, len(in))

	for i, instance := range in {

		cfg := instance.(map[string]interface{})
		ins := alkira.ConnectorAwsDirectConnectInstance{}

		if v, ok := cfg["name"].(string); ok {
			ins.Name = v
		}
		if v, ok := cfg["id"].(int); ok {
			ins.Id = v
		}
		if v, ok := cfg["connection_id"].(string); ok {
			ins.ConnectionId = v
		}
		if v, ok := cfg["dx_asn"].(int); ok {
			ins.DcGatewayAsn = v
		}
		if v, ok := cfg["dx_gateway_ip"].(string); ok {
			ins.AwsUnderlayIp = v
		}
		if v, ok := cfg["on_prem_asn"].(int); ok {
			ins.UnderlayAsn = v
		}
		if v, ok := cfg["on_prem_gateway_ip"].(string); ok {
			ins.OnPremUnderlayIp = v
		}
		if v, ok := cfg["underlay_prefix"].(string); ok {
			ins.UnderlayPrefix = v
		}
		if v, ok := cfg["bgp_auth_key"].(string); ok {
			ins.BgpAuthKey = v
		}
		if v, ok := cfg["bgp_auth_key_alkira"].(string); ok {
			ins.BgpAuthKeyAlkira = v
		}
		if v, ok := cfg["vlan_id"].(int); ok {
			ins.Vlan = v
		}
		if v, ok := cfg["aws_region"].(string); ok {
			ins.CustomerRegion = v
		}
		if v, ok := cfg["credential_id"].(string); ok {
			ins.CredentialId = v
		}
		if v, ok := cfg["gateway_mac_address"].(string); ok {
			ins.GatewayMacAddress = v
		}
		if v, ok := cfg["segment_options"].(*schema.Set); ok {
			segmentOptions, err := expandAwsDirectConnectSegmentOptions(v, m)

			if err != nil {
				return nil, err
			}
			ins.SegmentOptions = segmentOptions
		}

		instances[i] = ins
	}
	return instances, nil
}

func generateAwsDirectConnectRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorAwsDirectConnect, error) {

	// Expand instances
	instances, err := expandAwsDirectConnectInstances(d.Get("instance").([]interface{}), m)

	if err != nil {
		return nil, err
	}

	// Assemble request
	connector := &alkira.ConnectorAwsDirectConnect{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Cxp:            d.Get("cxp").(string),
		Enabled:        d.Get("enabled").(bool),
		Group:          d.Get("group").(string),
		TunnelProtocol: d.Get("tunnel_protocol").(string),
		BillingTags:    convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Size:           d.Get("size").(string),
		Instances:      instances,
	}

	return connector, nil
}
