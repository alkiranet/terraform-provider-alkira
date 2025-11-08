package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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
