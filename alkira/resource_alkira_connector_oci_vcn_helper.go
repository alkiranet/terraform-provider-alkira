package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setConnectorOciVcnRouting sets vcn_cidr, vcn_subnet, and vcn_route_table in state
// by parsing the VcnRouting interface{} returned from the API.
func setConnectorOciVcnRouting(d *schema.ResourceData, vcnRouting interface{}) {
	if vcnRouting == nil {
		return
	}

	routing, ok := vcnRouting.(map[string]interface{})
	if !ok {
		return
	}

	// Parse exportToCXPOptions → vcn_cidr or vcn_subnet
	if exportOpts, ok := routing["exportToCXPOptions"].(map[string]interface{}); ok {
		if prefixes, ok := exportOpts["userInputPrefixes"].([]interface{}); ok {
			var cidrList []string
			var subnetList []map[string]interface{}

			for _, p := range prefixes {
				prefix, ok := p.(map[string]interface{})
				if !ok {
					continue
				}
				prefixType, _ := prefix["type"].(string)
				value, _ := prefix["value"].(string)
				id, _ := prefix["id"].(string)

				if prefixType == "CIDR" {
					cidrList = append(cidrList, value)
				} else if prefixType == "SUBNET" {
					subnetList = append(subnetList, map[string]interface{}{
						"id":   id,
						"cidr": value,
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
	}

	// Parse importFromCXPOptions → vcn_route_table
	if importOpts, ok := routing["importFromCXPOptions"].(map[string]interface{}); ok {
		if routeTables, ok := importOpts["routeTables"].([]interface{}); ok {
			var tables []map[string]interface{}

			for _, rt := range routeTables {
				table, ok := rt.(map[string]interface{})
				if !ok {
					continue
				}
				id, _ := table["id"].(string)
				mode, _ := table["routeImportMode"].(string)

				var prefixListIds []int
				if pids, ok := table["prefixListIds"].([]interface{}); ok {
					for _, pid := range pids {
						if v, ok := pid.(float64); ok {
							prefixListIds = append(prefixListIds, int(v))
						}
					}
				}

				tables = append(tables, map[string]interface{}{
					"id":              id,
					"options":         mode,
					"prefix_list_ids": prefixListIds,
				})
			}

			if len(tables) > 0 {
				d.Set("vcn_route_table", tables)
			}
		}
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
