package alkira

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// setInstance set instance block values
func setInstance(d *schema.ResourceData, service *alkira.ServiceFortinet) {
	var instances []map[string]interface{}

	//
	// Go through all instance blocks from the config firstly to find a
	// match, instance's ID should be uniquely identifying an instance
	// block.
	//
	// On the first read call at the end of the create call, Terraform
	// didn't track any instance IDs yet.
	//
	for _, instance := range d.Get("instances").([]interface{}) {
		instanceConfig := instance.(map[string]interface{})

		for _, info := range service.Instances {

			if instanceConfig["id"].(int) == info.Id ||
				instanceConfig["name"].(string) == info.Name {

				log.Printf("[DEBUG] Found instance [%v|%v]", info.Id, info.Name)
				instance := map[string]interface{}{
					"id":                    info.Id,
					"credential_id":         info.CredentialId,
					"name":                  info.Name,
					"license_key_file_path": instanceConfig["license_key_file_path"].(string),
					"license_key":           instanceConfig["license_key"].(string),
					"serial_number":         info.SerialNumber,
				}

				instances = append(instances, instance)
				break
			}
		}
	}

	//
	// Go through all instances from the API response one more
	// time to find any instance that has not been tracked from Terraform
	// config.
	//
	for _, info := range service.Instances {
		new := true

		// Check if the instance already exists in the Terraform config
		for _, instance := range d.Get("instances").([]interface{}) {
			instanceConfig := instance.(map[string]interface{})

			if instanceConfig["id"].(int) == info.Id ||
				instanceConfig["name"].(string) == info.Name {
				new = false
				break
			}
		}

		// If the instance is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			i := map[string]interface{}{
				"id":            info.Id,
				"credential_id": info.CredentialId,
				"name":          info.Name,
				"serial_number": info.SerialNumber,
			}
			instances = append(instances, i)
			break
		}
	}

	d.Set("instances", instances)
}

func expandFortinetInstances(licenseType string, in []interface{}, m interface{}) ([]alkira.FortinetInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] invalid Fortinet instance input")
		return nil, errors.New("Invalid Fortinet instance input")
	}

	var licenseKeyPath string
	var licenseKeyLiteral string

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

				lk, err := extractLicenseKey(licenseType, licenseKeyLiteral, licenseKeyPath)
				if err != nil {
					return nil, err
				}
				c := alkira.CredentialFortinetInstance{
					LicenseKey:  lk,
					LicenseType: licenseType,
				}

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

// extractLicenseKey takes two string values. The order of the string
// parameters matters. After validation, if both fields have are not
// empty strings extractLicenseKey will default to using licenseKey as
// the return value. Otherwise extractLicenseKey will read from the
// licenseKeyPath and return the output as a string
func extractLicenseKey(licenseType string, licenseKey string, licenseKeyPath string) (string, error) {
	// if both params are empty
	if licenseKey == "" && licenseKeyPath == "" {

		// license key is optional for PAY_AS_YOU_GO
		if licenseType == "PAY_AS_YOU_GO" {
			return "", nil
		}

		return "", errors.New("either 'license_key' or 'icense_key_file_path' must be populated")
	}

	if licenseKey != "" {
		return licenseKey, nil
	}

	if _, err := os.Stat(licenseKeyPath); errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("file not found at %s: %w", licenseKeyPath, err)
	}

	b, err := os.ReadFile(licenseKeyPath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func generateFortinetRequest(d *schema.ResourceData, m interface{}) (*alkira.ServiceFortinet, error) {

	client := m.(*alkira.AlkiraClient)
	fortinetCredId := d.Get("credential_id").(string)

	if 0 == len(fortinetCredId) {
		log.Printf("[INFO] Creating Fortinet FW Credential")

		fortinetCredName := d.Get("name").(string) + randomNameSuffix()
		fortinetCred := alkira.CredentialPan{
			Username: d.Get("username").(string),
			Password: d.Get("password").(string),
		}

		credentialId, err := client.CreateCredential(
			fortinetCredName,
			alkira.CredentialTypeFortinet,
			fortinetCred,
			0,
		)
		if err != nil {
			return nil, err
		}
		d.Set("credential_id", credentialId)
	}

	billingTagIds := convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set))

	mgmtSegName, err := getSegmentNameById(d.Get("management_server_segment_id").(string), m)

	if err != nil {
		return nil, err
	}

	managementServer := &alkira.FortinetManagmentServer{
		IpAddress: d.Get("management_server_ip").(string),
		Segment:   mgmtSegName,
	}

	instances, err := expandFortinetInstances(
		d.Get("license_type").(string),
		d.Get("instances").([]interface{}),
		m,
	)
	if err != nil {
		return nil, err
	}

	// Convert segment IDs to segment names
	segmentNames, err := convertSegmentIdsToSegmentNames(d.Get("segment_ids").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	segmentOptions, err := expandSegmentOptions(d.Get("segment_options").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	service := &alkira.ServiceFortinet{
		AutoScale:        d.Get("auto_scale").(string),
		BillingTags:      billingTagIds,
		CredentialId:     d.Get("credential_id").(string),
		Cxp:              d.Get("cxp").(string),
		Instances:        instances,
		LicenseType:      d.Get("license_type").(string),
		ManagementServer: managementServer,
		MaxInstanceCount: d.Get("max_instance_count").(int),
		MinInstanceCount: d.Get("min_instance_count").(int),
		Name:             d.Get("name").(string),
		Segments:         segmentNames,
		SegmentOptions:   segmentOptions,
		Size:             d.Get("size").(string),
		TunnelProtocol:   d.Get("tunnel_protocol").(string),
		Version:          d.Get("version").(string),
	}

	return service, nil
}
