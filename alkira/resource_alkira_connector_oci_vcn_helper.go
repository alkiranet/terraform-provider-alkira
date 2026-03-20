package alkira

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setConnectorOciVcnRouting sets vcn_cidr, vcn_subnet, and vcn_route_table in state
// by unmarshaling the VcnRouting interface{} returned from the API into typed structs.
func setConnectorOciVcnRouting(d *schema.ResourceData, vcnRouting interface{}) {
	if vcnRouting == nil {
		return
	}

	routingJSON, err := json.Marshal(vcnRouting)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal VcnRouting: %v", err)
		return
	}

	var routing alkira.ConnectorOciVcnRouting
	if err := json.Unmarshal(routingJSON, &routing); err != nil {
		log.Printf("[ERROR] Failed to unmarshal VcnRouting: %v", err)
		return
	}

	setOciVcnExportPrefixes(routing.Export, d)
	setOciVcnImportRouteTables(routing.Import, d)
}

// setOciVcnExportPrefixes sets vcn_cidr or vcn_subnet from export options.
func setOciVcnExportPrefixes(exportOptions interface{}, d *schema.ResourceData) {
	if exportOptions == nil {
		return
	}

	exportJSON, err := json.Marshal(exportOptions)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal ExportOptions: %v", err)
		return
	}

	var export alkira.ConnectorOciVcnExportOptions
	if err := json.Unmarshal(exportJSON, &export); err != nil {
		log.Printf("[ERROR] Failed to unmarshal ExportOptions: %v", err)
		return
	}

	var cidrList []string
	var subnetList []map[string]interface{}

	for _, prefix := range export.Prefixes {
		switch prefix.Type {
		case "CIDR":
			cidrList = append(cidrList, prefix.Value)
		case "SUBNET":
			subnetList = append(subnetList, map[string]interface{}{
				"id":   prefix.Id,
				"cidr": prefix.Value,
			})
		}
	}

	if len(cidrList) > 0 {
		d.Set("vcn_cidr", cidrList)
	}
	if len(subnetList) > 0 {
		d.Set("vcn_subnet", subnetList)
	}
}

// setOciVcnImportRouteTables sets vcn_route_table from import options.
func setOciVcnImportRouteTables(importOptions interface{}, d *schema.ResourceData) {
	if importOptions == nil {
		return
	}

	importJSON, err := json.Marshal(importOptions)
	if err != nil {
		log.Printf("[ERROR] Failed to marshal ImportOptions: %v", err)
		return
	}

	var importOpts alkira.ConnectorOciVcnImportOptions
	if err := json.Unmarshal(importJSON, &importOpts); err != nil {
		log.Printf("[ERROR] Failed to unmarshal ImportOptions: %v", err)
		return
	}

	var tables []map[string]interface{}
	for _, rt := range importOpts.RouteTables {
		tables = append(tables, map[string]interface{}{
			"id":              rt.Id,
			"options":         rt.Mode,
			"prefix_list_ids": rt.PrefixListIds,
		})
	}

	if len(tables) > 0 {
		d.Set("vcn_route_table", tables)
	}
}

// expandConnectorOciVcnRouteTables expand OCI-VCN route tables
func expandConnectorOciVcnRouteTables(in *schema.Set) []alkira.ConnectorOciVcnRouteTables {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty VCN route table input")
		return []alkira.ConnectorOciVcnRouteTables{}
	}

	tables := make([]alkira.ConnectorOciVcnRouteTables, in.Len())
	for i, table := range in.List() {
		r := alkira.ConnectorOciVcnRouteTables{}
		t := table.(map[string]interface{})
		if v, ok := t["id"].(string); ok {
			r.Id = v
		}
		if v, ok := t["options"].(string); ok {
			r.Mode = v
		}

		r.PrefixListIds = convertTypeListToIntList(t["prefix_list_ids"].([]interface{}))
		tables[i] = r
	}

	return tables
}

// generateConnectorOciVcnUserInputPrefixes generate UserInputPrefixes used in connector-oci-vcn
func generateConnectorOciVcnUserInputPrefixes(cidr []interface{}, subnets *schema.Set) ([]alkira.ConnectorOciVcnInputPrefixes, error) {

	if len(cidr) == 0 && subnets == nil {
		return nil, fmt.Errorf("ERROR: either `vcn_subnet` or `vcn_cidr` must be specified")
	}

	// Processing "vcn_cidr"
	if len(cidr) > 0 {
		log.Printf("[DEBUG] Processing vcn_cidr %v", cidr)
		cidrList := make([]alkira.ConnectorOciVcnInputPrefixes, len(cidr))

		for i, value := range cidr {
			cidrList[i].Value = value.(string)
			cidrList[i].Type = "CIDR"
		}

		return cidrList, nil
	}

	// Processing VCN subnets
	log.Printf("[DEBUG] Processing vcn_subnet")
	if subnets == nil || subnets.Len() == 0 {
		log.Printf("[DEBUG] Empty vcn_subnet")
		return nil, fmt.Errorf("ERROR: Invalid vcn_subnet")
	}

	prefixes := make([]alkira.ConnectorOciVcnInputPrefixes, subnets.Len())
	for i, subnet := range subnets.List() {
		r := alkira.ConnectorOciVcnInputPrefixes{}
		t := subnet.(map[string]interface{})
		if v, ok := t["id"].(string); ok {
			r.Id = v
		}
		if v, ok := t["cidr"].(string); ok {
			r.Value = v
		}

		r.Type = "SUBNET"
		prefixes[i] = r
	}

	return prefixes, nil
}
