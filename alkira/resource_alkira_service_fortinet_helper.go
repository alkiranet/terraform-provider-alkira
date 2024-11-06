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

// expandFortinetInstances expand instance blocks to construct the request payload
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
				credentialId, err := createFortinetInstanceCredential(client, r.Name, licenseType, licenseKeyLiteral, licenseKeyPath)
				if err != nil {
					return nil, err
				}

				r.CredentialId = credentialId
			} else {
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

// createFortinetCredential
func createFortinetCredential(d *schema.ResourceData, c *alkira.AlkiraClient) (string, error) {

	log.Printf("[INFO] Creating Fortinet Credential")

	fortinetCredName := d.Get("name").(string) + "_" + randomNameSuffix()
	fortinetCred := alkira.CredentialPan{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}

	return c.CreateCredential(fortinetCredName, alkira.CredentialTypeFortinet, fortinetCred, 0)
}

// updateFortinetCredential
func updateFortinetCredential(d *schema.ResourceData, c *alkira.AlkiraClient) error {
	if d.HasChanges("username", "password") {
		log.Printf("[INFO] Fortinet credential has changed")
		id, err := createFortinetCredential(d, c)
		if err != nil {
			return err
		}
		d.Set("credential_id", id)
	}
	return nil
}

// deleteFortinetCredential
func deleteFortinetCredential(id string, c *alkira.AlkiraClient) error {

	log.Printf("[INFO] Deleting Fortinet Credential")
	return c.DeleteCredential(id, alkira.CredentialTypeFortinet)
}

// createFortinetInstanceCredential
func createFortinetInstanceCredential(c *alkira.AlkiraClient, name string, licenseType string, licenseKey string, licenseKeyPath string) (string, error) {

	log.Printf("[INFO] Creating Fortinet Instance Credential")

	// When license type is "PAY_AS_YOU_GO", license key is optional
	if licenseKey == "" && licenseKeyPath == "" {
		if licenseType == "PAY_AS_YOU_GO" {
			return "", nil
		}

		return "", errors.New("either 'license_key' or 'license_key_file_path' must be provided")
	}

	// If the license_key is provided directly in the config, use it,
	// otherwise, try to read it from the given license key file
	if licenseKey == "" {
		if _, err := os.Stat(licenseKeyPath); errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("file not found at %s: %w", licenseKeyPath, err)
		}

		b, err := os.ReadFile(licenseKeyPath)
		if err != nil {
			return "", err
		}
		licenseKey = string(b)

		if licenseKey == "" {
			return "", errors.New("'license_key' of 'service_fortinet_instance' is invalid")
		}
	}

	id, err := c.CreateCredential(
		name+randomNameSuffix(),
		alkira.CredentialTypeFortinetInstance,
		alkira.CredentialFortinetInstance{
			LicenseKey:  licenseKey,
			LicenseType: licenseType,
		},
		0)

	if err != nil {
		return "", err
	}

	return id, nil
}

func updateFortinetInstanceCredential(d *schema.ResourceData, c *alkira.AlkiraClient) error {
	if d.HasChanges("username", "password") {
		log.Printf("[INFO] FOrtinet credential has changed")
		id, err := createFortinetCredential(d, c)
		if err != nil {
			return err
		}
		d.Set("credential_id", id)
	}
	return nil
}

func generateFortinetRequest(d *schema.ResourceData, m interface{}) (*alkira.ServiceFortinet, error) {

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
		Scheme:           d.Get("license_scheme").(string),
		ManagementServer: managementServer,
		MaxInstanceCount: d.Get("max_instance_count").(int),
		MinInstanceCount: d.Get("min_instance_count").(int),
		Name:             d.Get("name").(string),
		Segments:         segmentNames,
		SegmentOptions:   segmentOptions,
		Size:             d.Get("size").(string),
		TunnelProtocol:   d.Get("tunnel_protocol").(string),
		Version:          d.Get("version").(string),
		Description:      d.Get("description").(string),
	}

	return service, nil
}
