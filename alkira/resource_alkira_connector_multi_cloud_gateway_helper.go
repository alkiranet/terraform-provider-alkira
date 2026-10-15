package alkira

import (
	"sort"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandConnectorMultiCloudGatewaySegmentOptions(in []interface{}, m interface{}) ([]alkira.ConnectorMultiCloudGatewaySegmentOption, error) {
	options := make([]alkira.ConnectorMultiCloudGatewaySegmentOption, len(in))
	for i, v := range in {
		cfg := v.(map[string]interface{})
		segmentName, err := getSegmentNameById(cfg["segment_id"].(string), m)
		if err != nil {
			return nil, err
		}
		options[i] = alkira.ConnectorMultiCloudGatewaySegmentOption{
			SegmentName:           segmentName,
			AdvertiseOnPremRoutes: cfg["advertise_on_prem_routes"].(bool),
			DisableInternetExit:   cfg["disable_internet_exit"].(bool),
		}
	}
	return options, nil
}

func flattenConnectorMultiCloudGatewaySegmentOptions(options []alkira.ConnectorMultiCloudGatewaySegmentOption, m interface{}) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, len(options))
	for i, o := range options {
		segmentId, err := getSegmentIdByName(o.SegmentName, m)
		if err != nil {
			return nil, err
		}
		result[i] = map[string]interface{}{
			"segment_id":               segmentId,
			"advertise_on_prem_routes": o.AdvertiseOnPremRoutes,
			"disable_internet_exit":    o.DisableInternetExit,
		}
	}
	return result, nil
}

// instanceIdsByName maps each stored instance's name to its id, so an update can match a
// config block to its stored instance by name instead of by list position.
func instanceIdsByName(instances []interface{}) map[string]int {
	ids := make(map[string]int, len(instances))
	for _, v := range instances {
		cfg := v.(map[string]interface{})
		if id, ok := cfg["id"].(int); ok && id != 0 {
			ids[cfg["name"].(string)] = id
		}
	}
	return ids
}

// expandConnectorMultiCloudGatewayInstances builds the request instances, carrying over the
// stored id of any instance whose name is unchanged; a renamed or new instance gets no id, so
// the backend treats it as new and removes the old one.
func expandConnectorMultiCloudGatewayInstances(in []interface{}, ids map[string]int) []alkira.ConnectorMultiCloudGatewayInstance {
	instances := make([]alkira.ConnectorMultiCloudGatewayInstance, len(in))
	for i, v := range in {
		cfg := v.(map[string]interface{})
		name := cfg["name"].(string)
		instances[i] = alkira.ConnectorMultiCloudGatewayInstance{Name: name, Id: ids[name]}
	}
	return instances
}

func flattenConnectorMultiCloudGatewayInstances(instances []alkira.ConnectorMultiCloudGatewayInstance) []map[string]interface{} {
	result := make([]map[string]interface{}, len(instances))
	for i, in := range instances {
		result[i] = map[string]interface{}{
			"id":                  in.Id,
			"name":                in.Name,
			"bandwidth_in_mbps":   in.BandwidthInMbps,
			"on_ramp":             in.OnRamp,
			"lumen_connection_id": in.LumenConnectionId,
		}
	}
	return result
}

// orderInstancesByPriorState reorders flattened instances to match the position each one has
// in priorOrder (matched by name), appending any instance not found there at the end in API
// order, so an API re-ordering of unchanged instances doesn't show as a permanent diff.
func orderInstancesByPriorState(flattened []map[string]interface{}, priorOrder []interface{}) []map[string]interface{} {
	position := make(map[string]int, len(priorOrder))
	for i, v := range priorOrder {
		cfg := v.(map[string]interface{})
		position[cfg["name"].(string)] = i
	}

	result := make([]map[string]interface{}, len(flattened))
	copy(result, flattened)

	sort.SliceStable(result, func(i, j int) bool {
		pi, iOk := position[result[i]["name"].(string)]
		pj, jOk := position[result[j]["name"].(string)]
		if iOk && jOk {
			return pi < pj
		}
		// Unknown instances keep their relative (API) order, after any known ones.
		return iOk && !jOk
	})

	return result
}

func generateConnectorMultiCloudGatewayRequest(d *schema.ResourceData, m interface{}) (*alkira.ConnectorMultiCloudGateway, error) {
	segmentOptions, err := expandConnectorMultiCloudGatewaySegmentOptions(d.Get("segment_options").([]interface{}), m)
	if err != nil {
		return nil, err
	}
	old, new := d.GetChange("instance")
	ids := instanceIdsByName(old.([]interface{}))

	return &alkira.ConnectorMultiCloudGateway{
		Name:                d.Get("name").(string),
		Description:         d.Get("description").(string),
		Cxp:                 d.Get("cxp").(string),
		Group:               d.Get("group").(string),
		Size:                d.Get("size").(string),
		Enabled:             d.Get("enabled").(bool),
		BillingTags:         convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		MultiCloudGatewayId: d.Get("multi_cloud_gateway_id").(int),
		LoopbackPrefix:      d.Get("loopback_prefix").(string),
		SegmentOptions:      segmentOptions,
		Instances:           expandConnectorMultiCloudGatewayInstances(new.([]interface{}), ids),
	}, nil
}

func setConnectorMultiCloudGateway(d *schema.ResourceData, connector *alkira.ConnectorMultiCloudGateway, m interface{}) diag.Diagnostics {
	segmentOptions, err := flattenConnectorMultiCloudGatewaySegmentOptions(connector.SegmentOptions, m)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("name", connector.Name)
	d.Set("description", connector.Description)
	d.Set("cxp", connector.Cxp)
	d.Set("group", connector.Group)
	d.Set("size", connector.Size)
	d.Set("enabled", connector.Enabled)
	d.Set("billing_tag_ids", connector.BillingTags)
	d.Set("multi_cloud_gateway_id", connector.MultiCloudGatewayId)
	d.Set("loopback_prefix", connector.LoopbackPrefix)
	d.Set("implicit_group_id", connector.ImplicitGroupId)
	if err := d.Set("segment_options", segmentOptions); err != nil {
		return diag.FromErr(err)
	}
	ordered := orderInstancesByPriorState(flattenConnectorMultiCloudGatewayInstances(connector.Instances), d.Get("instance").([]interface{}))
	if err := d.Set("instance", ordered); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
