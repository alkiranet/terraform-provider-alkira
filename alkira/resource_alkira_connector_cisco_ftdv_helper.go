package alkira

import (
	"errors"
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandCiscoFTDvInstances(name string, in *schema.Set, m interface{}) ([]alkira.CiscoFTDvInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid Cisco FTDv instance input")
		return nil, errors.New("Invalid Cisco FTDv instance input")
	}

	var adminPassword string
	var fmcRegistrationKey string
	var ftdvNatId string

	instances := make([]alkira.CiscoFTDvInstance, in.Len())

	for i, instance := range in.List() {
		r := alkira.CiscoFTDvInstance{}
		instanceCfg := instance.(map[string]interface{})

		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["hostname"].(string); ok {
			r.Hostname = v + "-" + randomNameSuffix()
		} else {
			r.Hostname = name + "-" + randomNameSuffix()
		}
		if v, ok := instanceCfg["admin_password"].(string); ok {
			adminPassword = v
		}
		if v, ok := instanceCfg["fmc_registration_key"].(string); ok {
			fmcRegistrationKey = v
		}
		if v, ok := instanceCfg["ftdv_nat_id"].(string); ok {
			ftdvNatId = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			if v == "" {
				credentialName := name + "-" + randomNameSuffix()
				c := alkira.CredentialCiscoFtdvInstance{
					AdminPassword:      adminPassword,
					FmcRegistrationKey: fmcRegistrationKey,
					FtvdNatId:          ftdvNatId,
				}
				log.Printf("[INFO] Creating Credential Cisco FTDv Instance.")
				credentialId, err := client.CreateCredential(
					credentialName,
					alkira.CredentialTypeCiscoFtdvInstance,
					c,
					0)
				if err != nil {
					return nil, err
				}
				r.CredentialId = credentialId
			} else {
				r.CredentialId = v
			}
		}
		if v, ok := instanceCfg["version"].(string); ok {
			r.Version = v
		}
		if v, ok := instanceCfg["license_type"].(string); ok {
			r.LicenseType = v
		}

		instances[i] = r
	}

	return instances, nil
}

func expandCiscoFtdvManagementServer(in *schema.Set, m interface{}) (alkira.CiscoFTDvManagementServer, error) {
	mg := alkira.CiscoFTDvManagementServer{}
	if in == nil || in.Len() != 1 {
		log.Printf("[DEBUG] Invalid Cisco FTDv Management Server input.")
		return mg, errors.New("Invalid Cisco FTDv Management Server input.")
	}

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})
		if v, ok := cfg["fmc_ip"].(string); ok {
			mg.IPAddress = v
		}
		if v, ok := cfg["segment_id"].(int); ok {
			mg.SegmentId = v
		}
		segName, err := convertSegmentIdToSegmentName(strconv.Itoa(mg.SegmentId), m)
		if err != nil {
			return mg, err
		}
		mg.Segment = segName
	}

	return mg, nil
}

func expandCiscoFtdvSegmentOptions(in *schema.Set, m interface{}) (alkira.SegmentNameToZone, error) {

	segmentOptions := make(alkira.SegmentNameToZone, in.Len())

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})

		zonestoGroups := make(alkira.ZoneToGroups)
		j := []string{}
		zonestoGroups[cfg["zone_name"].(string)] = j

		outerZoneToGroups := alkira.OuterZoneToGroups{
			SegmentId:     cfg["segment_id"].(int),
			ZonesToGroups: zonestoGroups,
		}

		segmentName, err := convertSegmentIdToSegmentName(strconv.Itoa(cfg["segment_id"].(int)), m)
		if err != nil {
			return nil, errors.New("Segment could not be found")
		}

		segmentOptions[segmentName] = outerZoneToGroups

	}

	return segmentOptions, nil
}

func deflateCiscoFTDvManagementServer(mg alkira.CiscoFTDvManagementServer) []map[string]interface{} {
	m := make(map[string]interface{})
	m["fmc_ip"] = mg.IPAddress
	m["segment_name"] = mg.Segment
	m["segment_id"] = mg.SegmentId

	return []map[string]interface{}{m}
}
