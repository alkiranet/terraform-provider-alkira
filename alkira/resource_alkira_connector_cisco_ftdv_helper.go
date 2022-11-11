package alkira

import (
	"errors"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
)

func expandCiscoFTDvInstances(name string, in []interface{}, m interface{}) ([]alkira.CiscoFTDvInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] invalid Cisco FTDv instance input")
		return nil, errors.New("Invalid Cisco FTDv instance input")
	}

	var adminPassword string
	var fmcRegistrationKey string
	var ftdvNatId string

	instances := make([]alkira.CiscoFTDvInstance, len(in))

	for i, instance := range in {
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

// func expandCiscoFtdvManagementServer(name string, in *schema.Set, m interface{}) (*alkira.CiscoFTDvManagementServer, error) {
// 	client := m.(*alkira.AlkiraClient)

// 	if in == nil || in.Len() > 1 {
// 		log.Printf("[DEBUG] Invalid Cisco FTDv Management Server input.")
// 		return nil, errors.New("Invalid Cisco FTDv Management Server input.")
// 	}

// 	mg := &alkira.CiscoFTDvManagementServer{}

// 	for _, option := range in.List() {
// 		cfg := option.(map[string]interface{})
// 		if v, ok := cfg[""]
// 	}

// }
