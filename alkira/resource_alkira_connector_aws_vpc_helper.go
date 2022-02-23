package alkira

import (
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandAwsVpcRouteTables expand AWS-VPC route tables
func expandAwsVpcRouteTables(in *schema.Set) []alkira.RouteTables {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] Empty VPC route table input")
		return []alkira.RouteTables{}
	}

	tables := make([]alkira.RouteTables, in.Len())
	for i, table := range in.List() {
		r := alkira.RouteTables{}
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

// generateUserInputPrefixes generate UserInputPrefixes used in AWS-VPC connector
func generateUserInputPrefixes(cidr []interface{}, subnets *schema.Set) ([]alkira.InputPrefixes, error) {

	if len(cidr) == 0 && subnets == nil {
		return nil, fmt.Errorf("ERROR: either \"vpc_subnet\" or \"vpc_cidr\" must be specified.")
	}

	// Processing "vpc_cidr"
	if len(cidr) > 0 {
		log.Printf("[DEBUG] Processing vpc_cidr %v", cidr)
		cidrList := make([]alkira.InputPrefixes, len(cidr))

		for i, value := range cidr {
			cidrList[i].Value = value.(string)
			cidrList[i].Type = "CIDR"
		}

		return cidrList, nil
	}

	// Processing VPC subnets
	log.Printf("[DEBUG] Processing vpc_subnet")
	if subnets == nil || subnets.Len() == 0 {
		log.Printf("[DEBUG] Empty vpc_subnet")
		return nil, fmt.Errorf("ERROR: Invalid vpc_subnet.")
	}

	prefixes := make([]alkira.InputPrefixes, subnets.Len())
	for i, subnet := range subnets.List() {
		r := alkira.InputPrefixes{}
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
