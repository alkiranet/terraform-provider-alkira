package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandConnectorRemoteAccessLdapSettings
func expandConnectorRemoteAccessLdapSettings(in *schema.Set, m interface{}) *alkira.ConnectorRemoteAccessLdapSettings {

	if in == nil || in.Len() == 0 {
		return nil
	}

	var ldapSettings alkira.ConnectorRemoteAccessLdapSettings

	for _, input := range in.List() {
		setting := input.(map[string]interface{})

		if v, ok := setting["bind_user_domain"].(string); ok {
			ldapSettings.BindUserDomain = v
		}
		if v, ok := setting["destination_address"].(string); ok {
			ldapSettings.DestinationAddress = v
		}
		if v, ok := setting["ldap_type"].(string); ok {
			ldapSettings.LdapType = v
		}
		if v, ok := setting["management_segment_id"].(int); ok {
			ldapSettings.ManagementSegmentId = v
		}
		if v, ok := setting["search_scope_domain"].(string); ok {
			ldapSettings.SearchScopeDomain = v
		}
	}

	return &ldapSettings
}

// expandConnectorRemoteAccessAuthorization
func expandConnectorRemoteAccessAuhtorization(in *schema.Set, cxp string, m interface{}) ([]alkira.ConnectorRemoteAccessSegmentOptions, error) {

	if in == nil || in.Len() == 0 {
		return nil, fmt.Errorf("Invalid connector-remote-acceauthorization")
	}

	var segmentOptions []alkira.ConnectorRemoteAccessSegmentOptions

	for _, input := range in.List() {
		auth := input.(map[string]interface{})

		//
		// segOption contains mapping
		// mapping contains cxpMapping
		//
		segOption := alkira.ConnectorRemoteAccessSegmentOptions{}
		mapping := alkira.ConnectorRemoteAccessUserGroupMappings{}
		cxpMapping := alkira.ConnectorRemoteAccessCxpToSubnetMapping{}

		if v, ok := auth["user_group_name"].(string); ok {
			mapping.Name = v
		}
		if v, ok := auth["segment_id"].(int); ok {
			segOption.SegmentId = v
		}
		if v, ok := auth["split_tunneling"].(bool); ok {
			mapping.SplitTunneling = v
		}
		if v, ok := auth["prefix_list_id"].(int); ok {
			mapping.PrefixListId = v
		}
		if v, ok := auth["subnet"].(string); ok {
			cxpMapping.Cxp = cxp
			cxpMapping.Subnets = []string{v}

			mapping.CxpToSubnetsMapping = []alkira.ConnectorRemoteAccessCxpToSubnetMapping{cxpMapping}
		}
		if v, ok := auth["billing_tag_id"].(int); ok {
			mapping.BillingTag = v
		}

		segOption.UserGroupMappings = []alkira.ConnectorRemoteAccessUserGroupMappings{mapping}
		segmentOptions = append(segmentOptions, segOption)
	}

	return segmentOptions, nil
}

// convertSegmentIdSetToStringList
func convertSegmentIdSetToStringList(in *schema.Set, m interface{}) []string {

	if in == nil || in.Len() == 0 {
		return nil
	}

	strList := make([]string, in.Len())

	for i, value := range in.List() {
		if value != nil {
			segmentName, err := getSegmentNameById(value.(string), m)

			if err != nil {
				return nil
			}

			strList[i] = segmentName
		} else {
			strList[i] = ""
		}
	}

	return strList
}

// generateConnectorAwsVpcRequest generate request for connector_aws_vpc
func generateConnectorRemoteAccessRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorRemoteAccessTemplate, error) {

	// Segment
	segmentNames := convertSegmentIdSetToStringList(d.Get("segment_ids").(*schema.Set), m)

	// Process Auth Options
	var authOptions alkira.ConnectorRemoteAccessAuthOptions

	if d.Get("authentication_mode").(string) != "LDAP" {
		authOptions = alkira.ConnectorRemoteAccessAuthOptions{
			SupportedModes: []string{d.Get("authentication_mode").(string)},
			LdapSettings:   nil,
		}
	} else {
		ldapSettings := expandConnectorRemoteAccessLdapSettings(d.Get("ldap_settings").(*schema.Set), m)

		authOptions = alkira.ConnectorRemoteAccessAuthOptions{
			SupportedModes: []string{d.Get("authentication_mode").(string)},
			LdapSettings:   ldapSettings,
		}
	}

	// Process authorization blocks
	segOptions, err := expandConnectorRemoteAccessAuhtorization(d.Get("authorization").(*schema.Set), d.Get("cxp").(string), m)

	if err != nil {
		return nil, err
	}

	// Assembly the request
	request := &alkira.ConnectorRemoteAccessTemplate{
		AdvancedOptions: alkira.ConnectorRemoteAccessAdvancedOptions{
			EnableDynamicRegionMapping: d.Get("enable_dynamic_region_mapping").(bool),
			MaxActiveUsersThreshold:    d.Get("concurrent_sessions_alert_threshold").(int),
			NameServer:                 d.Get("name_server").(string),
		},
		Arguments: []alkira.ConnectorRemoteAccessArguments{alkira.ConnectorRemoteAccessArguments{
			BillingTags: convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
			Cxp:         d.Get("cxp").(string),
			Size:        d.Get("size").(string),
		}},
		AuthenticationOptions: authOptions,
		Name:                  d.Get("name").(string),
		Segments:              segmentNames,
		SegmentOptions:        segOptions,
		BannerText:            d.Get("banner_text").(string),
	}

	return request, nil
}

// setAuthorization
func setAuthorization(d *schema.ResourceData, segmentOptions []alkira.ConnectorRemoteAccessSegmentOptions) {

	var authorizations []map[string]interface{}

	for _, option := range segmentOptions {

		if len(option.UserGroupMappings) != 1 ||
			len(option.UserGroupMappings[0].CxpToSubnetsMapping) != 1 {
			log.Printf("[ERROR] Invalid SegmentOptions in connector-remote-access")
			continue
		}

		auth := map[string]interface{}{
			"segment_id":      option.SegmentId,
			"user_group_name": option.UserGroupMappings[0].Name,
			"split_tunneling": option.UserGroupMappings[0].SplitTunneling,
			"prefix_list_id":  option.UserGroupMappings[0].PrefixListId,
			"billing_tag_id":  option.UserGroupMappings[0].BillingTag,
			"subnet":          option.UserGroupMappings[0].CxpToSubnetsMapping[0].Subnets,
		}

		authorizations = append(authorizations, auth)
	}

	d.Set("authorization", authorizations)
}

// setConnectorRemoteAccess
func setConnectorRemoteAccess(connector *alkira.ConnectorRemoteAccessTemplate, d *schema.ResourceData, m interface{}) error {

	d.Set("authentication_mode", connector.AuthenticationOptions.SupportedModes)

	// Set ldap_settings block
	if connector.AuthenticationOptions.LdapSettings != nil {
		var settings []map[string]interface{}
		setting := map[string]interface{}{
			"bind_user_domain":      connector.AuthenticationOptions.LdapSettings.BindUserDomain,
			"ldap_type":             connector.AuthenticationOptions.LdapSettings.LdapType,
			"destination_address":   connector.AuthenticationOptions.LdapSettings.DestinationAddress,
			"management_segment_id": connector.AuthenticationOptions.LdapSettings.ManagementSegmentId,
			"search_scope_domain":   connector.AuthenticationOptions.LdapSettings.SearchScopeDomain,
		}

		settings = append(settings, setting)
		d.Set("ldap_settings", settings)
	}

	d.Set("enable_dynamic_region_mapping", connector.AdvancedOptions.EnableDynamicRegionMapping)
	d.Set("name_server", connector.AdvancedOptions.NameServer)
	d.Set("concurrent_sessions_alert_threshold", connector.AdvancedOptions.MaxActiveUsersThreshold)
	d.Set("cxp", connector.Arguments[0].Cxp)
	d.Set("billing_tag_ids", connector.Arguments[0].BillingTags)
	d.Set("size", connector.Arguments[0].Size)
	d.Set("name", connector.Name)
	d.Set("banner_text", connector.BannerText)

	// Set segment_ids
	var segmentIds []string

	for _, segment := range connector.Segments {
		segmentId, err := getSegmentIdByName(segment, m)

		if err != nil {
			return err
		}

		segmentIds = append(segmentIds, segmentId)
	}

	d.Set("segment_ids", segmentIds)

	// Set authorization block
	setAuthorization(d, connector.SegmentOptions)

	return nil
}
