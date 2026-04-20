package alkira

import (
	"testing"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func remoteAccessTestSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"authentication_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cxp": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeString,
				Required: true,
			},
			"billing_tag_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"enable_dynamic_region_mapping": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"name_server": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"fallback_to_tcp": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"concurrent_sessions_alert_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  80,
			},
			"banner_text": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"segment_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ldap_settings": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"bind_user_domain": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ldap_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"destination_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"management_segment_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"search_scope_domain": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"authorization": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"segment_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"user_group_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"split_tunneling": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"prefix_list_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"billing_tag_id": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"subnet": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

// TestSetConnectorRemoteAccess_SingleMode verifies that a single-element
// SupportedModes slice is correctly unwrapped to a TypeString attribute.
func TestSetConnectorRemoteAccess_SingleMode(t *testing.T) {
	d := remoteAccessTestSchema().TestResourceData()

	connector := &alkira.ConnectorRemoteAccessTemplate{
		Name: "test-connector",
		AuthenticationOptions: alkira.ConnectorRemoteAccessAuthOptions{
			SupportedModes: []string{"LOCAL"},
		},
		AdvancedOptions: alkira.ConnectorRemoteAccessAdvancedOptions{
			EnableDynamicRegionMapping: true,
			MaxActiveUsersThreshold:    80,
			FallbackToTcp:              false,
		},
		Arguments: []alkira.ConnectorRemoteAccessArguments{
			{Cxp: "US-WEST", Size: "SMALL", BillingTags: []int{}},
		},
		Segments:       []string{},
		SegmentOptions: []alkira.ConnectorRemoteAccessSegmentOptions{},
	}

	err := setConnectorRemoteAccess(connector, d, nil)
	assert.NoError(t, err)
	assert.Equal(t, "LOCAL", d.Get("authentication_mode").(string))
}

// TestSetConnectorRemoteAccess_MultipleModes verifies that when the API
// returns multiple modes, the first one is used and no error occurs.
func TestSetConnectorRemoteAccess_MultipleModes(t *testing.T) {
	d := remoteAccessTestSchema().TestResourceData()

	connector := &alkira.ConnectorRemoteAccessTemplate{
		Name: "test-connector-multi",
		AuthenticationOptions: alkira.ConnectorRemoteAccessAuthOptions{
			SupportedModes: []string{"LOCAL", "LDAP"},
		},
		AdvancedOptions: alkira.ConnectorRemoteAccessAdvancedOptions{
			EnableDynamicRegionMapping: true,
			MaxActiveUsersThreshold:    80,
			FallbackToTcp:              false,
		},
		Arguments: []alkira.ConnectorRemoteAccessArguments{
			{Cxp: "US-WEST", Size: "SMALL", BillingTags: []int{}},
		},
		Segments:       []string{},
		SegmentOptions: []alkira.ConnectorRemoteAccessSegmentOptions{},
	}

	err := setConnectorRemoteAccess(connector, d, nil)
	assert.NoError(t, err)
	assert.Equal(t, "LOCAL", d.Get("authentication_mode").(string))
}

// TestSetConnectorRemoteAccess_EmptyModes verifies that when the API
// returns an empty SupportedModes, authentication_mode is left unset.
func TestSetConnectorRemoteAccess_EmptyModes(t *testing.T) {
	d := remoteAccessTestSchema().TestResourceData()

	connector := &alkira.ConnectorRemoteAccessTemplate{
		Name: "test-connector-empty",
		AuthenticationOptions: alkira.ConnectorRemoteAccessAuthOptions{
			SupportedModes: []string{},
		},
		AdvancedOptions: alkira.ConnectorRemoteAccessAdvancedOptions{
			EnableDynamicRegionMapping: true,
			MaxActiveUsersThreshold:    80,
			FallbackToTcp:              false,
		},
		Arguments: []alkira.ConnectorRemoteAccessArguments{
			{Cxp: "US-WEST", Size: "SMALL", BillingTags: []int{}},
		},
		Segments:       []string{},
		SegmentOptions: []alkira.ConnectorRemoteAccessSegmentOptions{},
	}

	err := setConnectorRemoteAccess(connector, d, nil)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Get("authentication_mode").(string))
}

// TestSetConnectorRemoteAccess_EmptyArguments verifies that an empty
// Arguments slice returns an error instead of panicking.
func TestSetConnectorRemoteAccess_EmptyArguments(t *testing.T) {
	d := remoteAccessTestSchema().TestResourceData()

	connector := &alkira.ConnectorRemoteAccessTemplate{
		Name: "test-connector-no-args",
		AuthenticationOptions: alkira.ConnectorRemoteAccessAuthOptions{
			SupportedModes: []string{"LOCAL"},
		},
		AdvancedOptions: alkira.ConnectorRemoteAccessAdvancedOptions{
			EnableDynamicRegionMapping: true,
			MaxActiveUsersThreshold:    80,
			FallbackToTcp:              false,
		},
		Arguments:      []alkira.ConnectorRemoteAccessArguments{},
		Segments:       []string{},
		SegmentOptions: []alkira.ConnectorRemoteAccessSegmentOptions{},
	}

	err := setConnectorRemoteAccess(connector, d, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty arguments")
}

// TestSetConnectorRemoteAccess_FieldMapping verifies that all fields
// from the API response are correctly mapped to resource data.
func TestSetConnectorRemoteAccess_FieldMapping(t *testing.T) {
	d := remoteAccessTestSchema().TestResourceData()

	connector := &alkira.ConnectorRemoteAccessTemplate{
		Name:       "my-remote-access",
		BannerText: "Welcome",
		AuthenticationOptions: alkira.ConnectorRemoteAccessAuthOptions{
			SupportedModes: []string{"SAML"},
		},
		AdvancedOptions: alkira.ConnectorRemoteAccessAdvancedOptions{
			EnableDynamicRegionMapping: false,
			MaxActiveUsersThreshold:    50,
			NameServer:                 "8.8.8.8",
			FallbackToTcp:              true,
		},
		Arguments: []alkira.ConnectorRemoteAccessArguments{
			{Cxp: "US-EAST", Size: "MEDIUM", BillingTags: []int{1, 2}},
		},
		Segments:       []string{},
		SegmentOptions: []alkira.ConnectorRemoteAccessSegmentOptions{},
	}

	err := setConnectorRemoteAccess(connector, d, nil)
	assert.NoError(t, err)

	assert.Equal(t, "SAML", d.Get("authentication_mode").(string))
	assert.Equal(t, "my-remote-access", d.Get("name").(string))
	assert.Equal(t, "Welcome", d.Get("banner_text").(string))
	assert.Equal(t, "US-EAST", d.Get("cxp").(string))
	assert.Equal(t, "MEDIUM", d.Get("size").(string))
	assert.Equal(t, false, d.Get("enable_dynamic_region_mapping").(bool))
	assert.Equal(t, 50, d.Get("concurrent_sessions_alert_threshold").(int))
	assert.Equal(t, "8.8.8.8", d.Get("name_server").(string))
	assert.Equal(t, true, d.Get("fallback_to_tcp").(bool))
}

// TestSetAuthorization_SingleAuth verifies that a single authorization
// block is correctly set in resource data.
func TestSetAuthorization_SingleAuth(t *testing.T) {
	d := remoteAccessTestSchema().TestResourceData()

	segOptions := []alkira.ConnectorRemoteAccessSegmentOptions{
		{
			SegmentId: 867,
			UserGroupMappings: []alkira.ConnectorRemoteAccessUserGroupMappings{
				{
					Name:           "test-group",
					SplitTunneling: false,
					PrefixListId:   0,
					BillingTag:     0,
					CxpToSubnetsMapping: []alkira.ConnectorRemoteAccessCxpToSubnetMapping{
						{Cxp: "US-WEST", Subnets: []string{"10.99.0.0/24"}},
					},
				},
			},
		},
	}

	setAuthorization(d, segOptions)

	auths := d.Get("authorization").(*schema.Set).List()
	assert.Len(t, auths, 1)

	auth := auths[0].(map[string]interface{})
	assert.Equal(t, 867, auth["segment_id"])
	assert.Equal(t, "test-group", auth["user_group_name"])
	assert.Equal(t, false, auth["split_tunneling"])
	assert.Equal(t, "10.99.0.0/24", auth["subnet"])
}

// TestSetAuthorization_MultipleAuths verifies that multiple authorization
// blocks are correctly set in resource data.
func TestSetAuthorization_MultipleAuths(t *testing.T) {
	d := remoteAccessTestSchema().TestResourceData()

	segOptions := []alkira.ConnectorRemoteAccessSegmentOptions{
		{
			SegmentId: 100,
			UserGroupMappings: []alkira.ConnectorRemoteAccessUserGroupMappings{
				{
					Name:           "group-a",
					SplitTunneling: true,
					CxpToSubnetsMapping: []alkira.ConnectorRemoteAccessCxpToSubnetMapping{
						{Cxp: "US-WEST", Subnets: []string{"10.1.0.0/24"}},
					},
				},
			},
		},
		{
			SegmentId: 200,
			UserGroupMappings: []alkira.ConnectorRemoteAccessUserGroupMappings{
				{
					Name:           "group-b",
					SplitTunneling: false,
					CxpToSubnetsMapping: []alkira.ConnectorRemoteAccessCxpToSubnetMapping{
						{Cxp: "US-EAST", Subnets: []string{"10.2.0.0/24"}},
					},
				},
			},
		},
	}

	setAuthorization(d, segOptions)

	auths := d.Get("authorization").(*schema.Set).List()
	assert.Len(t, auths, 2)

	groupNames := make([]string, 0, 2)
	for _, a := range auths {
		auth := a.(map[string]interface{})
		groupNames = append(groupNames, auth["user_group_name"].(string))
	}
	assert.ElementsMatch(t, []string{"group-a", "group-b"}, groupNames)
}

// TestSetAuthorization_InvalidMapping verifies that an authorization
// with invalid mapping structure is skipped without panic.
func TestSetAuthorization_InvalidMapping(t *testing.T) {
	d := remoteAccessTestSchema().TestResourceData()

	segOptions := []alkira.ConnectorRemoteAccessSegmentOptions{
		{
			SegmentId:         100,
			UserGroupMappings: []alkira.ConnectorRemoteAccessUserGroupMappings{},
		},
	}

	setAuthorization(d, segOptions)

	auths := d.Get("authorization").(*schema.Set).List()
	assert.Len(t, auths, 0)
}
