package alkira

import (
	"errors"
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandCiscoFTDvInstances(in []interface{}, m interface{}) ([]alkira.CiscoFTDvInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] invalid Cisco FTDv instance input")
		return nil, errors.New("ERROR: Invalid Cisco FTDv instance input")
	}

	var adminPassword string
	var fmcRegistrationKey string
	var ftdvNatId string

	instances := make([]alkira.CiscoFTDvInstance, 0, len(in))

	for _, instance := range in {

		r := alkira.CiscoFTDvInstance{}
		instanceCfg := instance.(map[string]interface{})

		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["hostname"].(string); ok {
			r.Hostname = v
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
				credentialName := r.Hostname + "-" + randomNameSuffix()

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
		if v, ok := instanceCfg["enable_traffic"].(bool); ok {
			r.TrafficEnabled = v
		}

		instances = append(instances, r)
	}

	return instances, nil
}

func expandCiscoFtdvManagementServer(in *schema.Set, m interface{}) (string, []string, alkira.CiscoFTDvManagementServer, error) {
	client := m.(*alkira.AlkiraClient)

	var credentialId string
	var ipAllowList = []string{}
	var managementServer = alkira.CiscoFTDvManagementServer{}

	if in == nil || in.Len() != 1 {
		log.Printf("[DEBUG] Invalid Cisco FTDv Management Server input.")
		return credentialId, ipAllowList, managementServer, errors.New("ERROR: Invalid Cisco FTDv Management Server input")
	}

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})

		var username string
		var password string

		if v, ok := cfg["username"].(string); ok {
			username = v
		}
		if v, ok := cfg["password"].(string); ok {
			password = v
		}
		if v, ok := cfg["credential_id"].(string); ok {
			if v == "" {
				credentialName := "cisco-fdtv-" + randomNameSuffix()
				c := alkira.CredentialCiscoFtdv{Username: username, Password: password}

				log.Printf("[INFO] Creating Cisco FTDv Firewall Service Credentials")
				cId, err := client.CreateCredential(credentialName, alkira.CredentialTypeCiscoFtdv, c, 0)

				if err != nil {
					return credentialId, ipAllowList, managementServer, err
				}

				credentialId = cId
			} else {
				credentialId = v
			}
		}
		if v, ok := cfg["server_ip"].(string); ok {
			managementServer.IPAddress = v
		}
		if v, ok := cfg["ip_allow_list"].([]interface{}); ok {
			ipAllowList = convertTypeListToStringList(v)
		}
		if v, ok := cfg["segment_id"].(string); ok {
			segmentName, _ := getSegmentNameById(v, m)
			managementServer.Segment = segmentName
		}
	}

	return credentialId, ipAllowList, managementServer, nil
}

func expandCiscoFtdvSegmentOptions(in *schema.Set, m interface{}) (alkira.SegmentNameToZone, error) {

	segmentOptions := make(alkira.SegmentNameToZone, in.Len())

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})

		segmentName, err := getSegmentNameById(cfg["segment_id"].(string), m)

		if err != nil {
			return nil, err
		}

		groupList := convertTypeListToStringList(cfg["groups"].([]interface{}))

		if groupList == nil {
			groupList = []string{}
		}

		if val, ok := segmentOptions[segmentName]; ok {
			val.ZonesToGroups[cfg["zone_name"].(string)] = groupList

		} else {
			zonestoGroups := make(alkira.ZoneToGroups)
			zonestoGroups[cfg["zone_name"].(string)] = groupList

			segmentId, _ := strconv.Atoi(cfg["segment_id"].(string))

			outerZoneToGroups := alkira.OuterZoneToGroups{
				SegmentId:     segmentId,
				ZonesToGroups: zonestoGroups,
			}

			segmentOptions[segmentName] = outerZoneToGroups
		}

	}

	return segmentOptions, nil
}

func deflateCiscoFTDvManagementServer(service *alkira.ServiceCiscoFTDv, m interface{}) []map[string]interface{} {

	result := make(map[string]interface{})

	result["credential_id"] = service.CredentialId
	result["ip_allow_list"] = service.IpAllowList
	result["server_ip"] = service.ManagementServer.IPAddress

	// Convert segment name to segment ID for import support
	if service.ManagementServer.Segment != "" && m != nil {
		segmentId, err := getSegmentIdByName(service.ManagementServer.Segment, m)
		if err == nil {
			result["segment_id"] = segmentId
		}
	}

	return []map[string]interface{}{result}
}

// setCiscoFTDvInstances
func setCiscoFTDvInstances(d *schema.ResourceData, c []alkira.CiscoFTDvInstance) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, value := range d.Get("instance").([]interface{}) {
		cfg := value.(map[string]interface{})

		for _, ins := range c {
			if cfg["id"].(int) == ins.Id || cfg["hostname"].(string) == ins.Hostname {
				instance := map[string]interface{}{
					"admin_password":       cfg["admin_password"].(string),
					"credential_id":        ins.CredentialId,
					"fmc_registration_key": cfg["fmc_registration_key"].(string),
					"ftdv_nat_id":          cfg["ftdv_nat_id"].(string),
					"hostname":             ins.Hostname,
					"id":                   ins.Id,
					"license_type":         ins.LicenseType,
					"version":              ins.Version,
					"enable_traffic":       ins.TrafficEnabled,
				}
				instances = append(instances, instance)
				break
			}
		}
	}

	for _, instance := range c {
		new := true

		// Check if the instance already exists in the Terraform config
		for _, ins := range d.Get("instance").([]interface{}) {
			cfg := ins.(map[string]interface{})

			if cfg["id"].(int) == instance.Id || cfg["hostname"].(string) == instance.Hostname {
				new = false
				break
			}
		}

		// If the instance is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			instance := map[string]interface{}{
				"credential_id":  instance.CredentialId,
				"hostname":       instance.Hostname,
				"id":             instance.Id,
				"license_type":   instance.LicenseType,
				"version":        instance.Version,
				"enable_traffic": instance.TrafficEnabled,
			}

			instances = append(instances, instance)
			break
		}
	}

	return instances
}
