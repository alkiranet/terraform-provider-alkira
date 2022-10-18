package alkira

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandFortinetInstances(licenseType string, in []interface{}, m interface{}) ([]alkira.FortinetInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] invalid Fortinet instance input")
		return nil, errors.New("Invalid Fortinet instance input")
	}

	var licenseKeyPath string
	var licenseKeyLiteral string
	//var err error

	instances := make([]alkira.FortinetInstance, len(in))
	for i, instance := range in {
		r := alkira.FortinetInstance{}
		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
			r.HostName = v
		}
		if v, ok := instanceCfg["serial_number"].(string); ok {
			r.SerialNumber = v
		}
		if v, ok := instanceCfg["license_key_file_path"].(string); ok {
			licenseKeyPath = v
		}
		if v, ok := instanceCfg["license_key"].(string); ok {
			licenseKeyLiteral = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			if v == "" {

				lk, err := extractLicenseKey(licenseKeyLiteral, licenseKeyPath)
				if err != nil {
					return nil, err
				}
				c := alkira.CredentialFortinetInstance{
					LicenseKey:  lk,
					LicenseType: licenseType,
				}

				//r.SerialNumber, err = extractLicenseKey(
				//	r.SerialNumber,
				//	licenseKeyPath,
				//)
				//if err != nil {
				//	return nil, err
				//}
				//c := alkira.CredentialFortinetInstance{
				//	LicenseKey:  r.SerialNumber,
				//	LicenseType: licenseType,
				//}

				credentialName := r.Name + randomNameSuffix()

				log.Printf("[INFO] Creating Fortinet Instance Credential %s", credentialName)

				credentialId, err := client.CreateCredential(
					credentialName,
					alkira.CredentialTypeFortinetInstance,
					c,
					0,
				)
				if err != nil {
					return nil, err
				}

				r.CredentialId = credentialId
			}

			if v != "" {
				r.CredentialId = v
			}
		}
		instances[i] = r
	}

	return instances, nil
}

func expandFortinetZone(in *schema.Set) map[string][]string {
	zonesToGroups := make(map[string][]string)

	for _, zone := range in.List() {
		zoneCfg := zone.(map[string]interface{})
		var name *string
		var groups []string

		if v, ok := zoneCfg["name"].(string); ok {
			name = &v
		}

		if v, ok := zoneCfg["groups"].([]interface{}); ok {
			groups = convertTypeListToStringList(v)
		}

		zonesToGroups[*name] = groups
	}

	return zonesToGroups
}

// extractLicenseKey takes two string values. The order of the string parameters matters. After
// validation, if both fields have are not empty strings extractLicenseKey will default to using
// licenseKey as the return value. Otherwise extractLicenseKey will read from the licenseKeyPath
// and return the output as a string.
func extractLicenseKey(licenseKey string, licenseKeyPath string) (string, error) {
	// if both params are empty
	if licenseKey == "" && licenseKeyPath == "" {
		return "", errors.New("either license_key or license_key_file_path must be populated")
	}

	if licenseKey != "" {
		return licenseKey, nil
	}

	fmt.Println("HARPO: ", licenseKeyPath)
	if _, err := os.Stat(licenseKeyPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("HARPOHARPO2")
		return "", fmt.Errorf("file not found at %s: %w", licenseKeyPath, err)
	}

	b, err := os.ReadFile(licenseKeyPath)
	if err != nil {
		fmt.Println("HARPOHARPO1")
		return "", err
	}

	return string(b), nil
}
