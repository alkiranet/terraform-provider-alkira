package alkira

import (
	"errors"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// UNUSED: Commented out to suppress linter warnings
// type panZone struct {
// 	Segment string
// 	Zone    string
// 	Groups  interface{}
// }

// Helper functions for PAN credentials
func createPanCredential(d *schema.ResourceData, c *alkira.AlkiraClient) (string, error) {
	log.Printf("[INFO] Creating PAN Credential")

	credentialName := d.Get("name").(string) + randomNameSuffix()
	credential := alkira.CredentialPan{
		Username:   d.Get("pan_username").(string),
		Password:   d.Get("pan_password").(string),
		LicenseKey: d.Get("pan_license_key").(string),
	}
	d.Set("pan_credential_name", credentialName)

	return c.CreateCredential(credentialName, alkira.CredentialTypePan, credential, 0)
}

func updatePanCredential(d *schema.ResourceData, c *alkira.AlkiraClient) error {
	log.Printf("[INFO] Updating PAN Credential")

	if d.HasChanges("pan_username", "pan_password", "pan_license_key") {
		log.Printf("[INFO] PAN credential has changed")

		if d.Get("pan_credential_id") == nil {
			return errors.New("pan_credential_id is empty when updating PAN credential")
		} else {
			if d.Get("pan_credential_name") == nil || d.Get("pan_credential_name").(string) == "" {
				return errors.New("pan_credential_name is empty when updating PAN credential")
			}

			credentialId := d.Get("pan_credential_id").(string)
			credentialName := d.Get("pan_credential_name").(string)
			credential := alkira.CredentialPan{
				Username:   d.Get("pan_username").(string),
				Password:   d.Get("pan_password").(string),
				LicenseKey: d.Get("pan_license_key").(string),
			}
			return c.UpdateCredential(credentialId, credentialName, alkira.CredentialTypePan, credential, 0)
		}
	}

	return nil
}

// UNUSED: Commented out to suppress linter warnings
// func deletePanCredential(id string, c *alkira.AlkiraClient) error {
// 	log.Printf("[INFO] Deleting PAN Credential")
// 	return c.DeleteCredential(id, alkira.CredentialTypePan)
// }

// Helper functions for PAN Registration Credentials
func createPanRegistrationCredential(d *schema.ResourceData, c *alkira.AlkiraClient) (string, error) {
	var credentialExpiry int64
	var expiryValue string

	if v, ok := d.GetOk("registration_pin_expiry"); ok {
		expiryValue = v.(string)
		log.Printf("[INFO] Creating PAN Registration Credential with expiry %v", expiryValue)
		var err error
		credentialExpiry, err = convertInputTimeToEpoch(expiryValue)
		if err != nil {
			log.Printf("[ERROR] failed to parse 'registration_pin_expiry', %v", err)
			return "", err
		}
	} else {
		log.Printf("[INFO] Creating PAN Registration Credential with no expiry")
		credentialExpiry = 0
	}

	credentialName := d.Get("name").(string) + randomNameSuffix()
	credential := alkira.CredentialPanRegistration{
		RegistrationPinId:    d.Get("registration_pin_id").(string),
		RegistrationPinValue: d.Get("registration_pin_value").(string),
	}

	return c.CreateCredential(credentialName, alkira.CredentialTypePanRegistration, credential, credentialExpiry)
}

// UNUSED: Commented out to suppress linter warnings
// func deletePanRegistrationCredential(id string, c *alkira.AlkiraClient) error {
// 	log.Printf("[INFO] Deleting PAN Registration Credential")
// 	return c.DeleteCredential(id, alkira.CredentialTypePanRegistration)
// }

// Helper function for PAN Master Key Credential
func createPanMasterKeyCredential(d *schema.ResourceData, c *alkira.AlkiraClient) (string, error) {
	log.Printf("[INFO] Creating PAN Master Key Credential %v", d.Get("master_key_expiry").(string))

	if !d.Get("master_key_enabled").(bool) {
		log.Printf("[INFO] PAN master key is not enabled, skip creating credential")
		return "", nil
	}

	credentialName := d.Get("name").(string) + randomNameSuffix()
	credential := alkira.CredentialPanMasterKey{
		MasterKey: d.Get("master_key").(string),
	}

	credentialExpiry, err := convertInputTimeToEpoch(d.Get("master_key_expiry").(string))

	if err != nil {
		log.Printf("[ERROR] failed to parse 'master_key_expiry', %v", err)
		return "", err
	}

	if credentialExpiry == 0 {
		log.Printf("[ERROR] argument 'master_key_expiry' is required when master key was enabled.")
		return "", err
	}

	return c.CreateCredential(credentialName, alkira.CredentialTypePanMasterKey, credential, credentialExpiry)
}

// UNUSED: Commented out to suppress linter warnings
// func deletePanMasterKeyCredential(id string, c *alkira.AlkiraClient) error {
// 	log.Printf("[INFO] Deleting PAN Master Key Credential")
// 	return c.DeleteCredential(id, alkira.CredentialTypePanMasterKey)
// }

// Create all credentials of PAN service
//
// - PAN Credential
// - PAN Registration Credential
// - PAN Master Key Credential
func createCredentials(d *schema.ResourceData, c *alkira.AlkiraClient) error {

	// Create PAN credentail
	panCredentialId, err := createPanCredential(d, c)
	if err != nil {
		return err
	}

	d.Set("pan_credential_id", panCredentialId)

	// Create PAN Registration Credential
	panRegistrationCredentialId, err := createPanRegistrationCredential(d, c)
	if err != nil {
		return err
	}
	d.Set("pan_registration_credential_id", panRegistrationCredentialId)

	// Create PAN Master Key Credential
	panMasterKeyCredentialId, err := createPanMasterKeyCredential(d, c)
	if err != nil {
		return err
	}
	d.Set("pan_master_key_credential_id", panMasterKeyCredentialId)

	return nil
}

func updateCredentials(d *schema.ResourceData, c *alkira.AlkiraClient) error {

	// Update PAN credentail
	err := updatePanCredential(d, c)
	if err != nil {
		return err
	}

	return nil
}

// Global Protect Segment Options
func expandGlobalProtectSegmentOptions(in *schema.Set, m interface{}) (map[string]*alkira.GlobalProtectSegmentName, error) {

	if in == nil || in.Len() == 0 {
		return nil, nil
	}

	sgmtOptions := make(map[string]*alkira.GlobalProtectSegmentName)
	for _, sgmtOption := range in.List() {
		r := &alkira.GlobalProtectSegmentName{}
		segmentCfg := sgmtOption.(map[string]interface{})
		var segmentName string

		if v, ok := segmentCfg["segment_id"].(string); ok {
			segmentName, _ = getSegmentNameById(v, m)
		}
		if v, ok := segmentCfg["remote_user_zone_name"].(string); ok {
			r.RemoteUserZoneName = v
		}
		if v, ok := segmentCfg["portal_fqdn_prefix"].(string); ok {
			r.PortalFqdnPrefix = v
		}
		if v, ok := segmentCfg["service_group_name"].(string); ok {
			r.ServiceGroupName = v
		}

		sgmtOptions[segmentName] = r
	}

	return sgmtOptions, nil
}

func expandGlobalProtectSegmentOptionsInstance(in *schema.Set, m interface{}) (map[string]*alkira.GlobalProtectSegmentNameInstance, error) {

	if in == nil || in.Len() == 0 {
		return nil, nil
	}

	sgmtOptions := make(map[string]*alkira.GlobalProtectSegmentNameInstance)
	for _, sgmtOption := range in.List() {
		r := &alkira.GlobalProtectSegmentNameInstance{}
		segmentCfg := sgmtOption.(map[string]interface{})
		var segmentName string

		if v, ok := segmentCfg["segment_id"].(string); ok {
			segmentName, _ = getSegmentNameById(v, m)
		}
		if v, ok := segmentCfg["portal_enabled"].(bool); ok {
			r.PortalEnabled = v
		}
		if v, ok := segmentCfg["gateway_enabled"].(bool); ok {
			r.GatewayEnabled = v
		}
		if v, ok := segmentCfg["prefix_list_id"].(int); ok {
			r.PrefixListId = v
		}

		sgmtOptions[segmentName] = r
	}

	return sgmtOptions, nil
}

func flattenGlobalProtectSegmentOptions(in map[string]*alkira.GlobalProtectSegmentName, m interface{}) []map[string]interface{} {

	if len(in) == 0 {
		return nil
	}

	var options []map[string]interface{}

	for segmentName, option := range in {
		segmentId, err := getSegmentIdByName(segmentName, m)
		if err != nil {
			log.Printf("[WARNING] Failed to get segment ID for segment name %s: %v", segmentName, err)
			continue
		}

		opt := map[string]interface{}{
			"segment_id":            segmentId,
			"remote_user_zone_name": option.RemoteUserZoneName,
			"portal_fqdn_prefix":    option.PortalFqdnPrefix,
			"service_group_name":    option.ServiceGroupName,
		}
		options = append(options, opt)
	}

	return options
}

func flattenGlobalProtectSegmentOptionsInstance(in map[string]*alkira.GlobalProtectSegmentNameInstance, m interface{}) []map[string]interface{} {

	if len(in) == 0 {
		return nil
	}

	var options []map[string]interface{}

	for segmentName, option := range in {
		segmentId, err := getSegmentIdByName(segmentName, m)
		if err != nil {
			log.Printf("[WARNING] Failed to get segment ID for segment name %s: %v", segmentName, err)
			continue
		}

		opt := map[string]interface{}{
			"segment_id":      segmentId,
			"portal_enabled":  option.PortalEnabled,
			"gateway_enabled": option.GatewayEnabled,
			"prefix_list_id":  option.PrefixListId,
		}
		options = append(options, opt)
	}

	return options
}

// UNUSED: Commented out to suppress linter warnings
// func expandPanSegmentOptions(in *schema.Set, m interface{}) (map[string]interface{}, error) {
//
// 	if in == nil {
// 		return nil, errors.New("invalid SegmentOptions input")
// 	}
//
// 	zoneMap := make([]panZone, in.Len())
//
// 	for i, option := range in.List() {
// 		r := panZone{}
// 		cfg := option.(map[string]interface{})
// 		if v, ok := cfg["segment_id"].(string); ok {
// 			segmentName, err := getSegmentNameById(v, m)
//
// 			if err != nil {
// 				return nil, err
// 			}
// 			r.Segment = segmentName
// 		}
// 		if v, ok := cfg["zone_name"].(string); ok {
// 			r.Zone = v
// 		}
//
// 		r.Groups = cfg["groups"]
//
// 		zoneMap[i] = r
// 	}
//
// 	segmentOptions := make(map[string]interface{})
//
// 	for _, x := range zoneMap {
// 		zone := make(map[string]interface{})
// 		zone[x.Zone] = x.Groups
//
// 		for _, y := range zoneMap {
// 			if x.Segment == y.Segment {
// 				zone[y.Zone] = y.Groups
// 			}
// 		}
//
// 		zonesToGroups := make(map[string]interface{})
// 		zonesToGroups["zonesToGroups"] = zone
//
// 		segmentOptions[x.Segment] = zonesToGroups
// 	}
//
// 	return segmentOptions, nil
// }

// expand "instance" block from config to generate request payload
func expandPanInstances(in []interface{}, m interface{}) ([]alkira.ServicePanInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if len(in) == 0 {
		return nil, errors.New("invalid PAN instance input")
	}

	instances := make([]alkira.ServicePanInstance, len(in))
	for i, instance := range in {
		r := alkira.ServicePanInstance{}
		instanceCfg := instance.(map[string]interface{})

		var authCode string
		var authKey string
		var authExpiry int64

		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := instanceCfg["auth_code"].(string); ok {
			authCode = v
		}
		if v, ok := instanceCfg["auth_key"].(string); ok {
			authKey = v
		}
		if v, ok := instanceCfg["auth_expiry"].(string); ok {
			if v != "" {
				var err error
				authExpiry, err = convertInputTimeToEpoch(v)

				if err != nil {
					log.Printf("[ERROR] failed to parse 'auth_expiry', %v", err)
					return nil, err
				}
			}
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			if v == "" {
				credentialName := r.Name + randomNameSuffix()
				credentialPanInstance := alkira.CredentialPanInstance{
					AuthCode: authCode,
					AuthKey:  authKey,
				}

				log.Printf("[INFO] Creating PAN Instance Credential %s", credentialName)
				credentialId, err := client.CreateCredential(
					credentialName,
					alkira.CredentialTypePanInstance,
					credentialPanInstance,
					authExpiry,
				)

				if err != nil {
					return nil, err
				}

				r.CredentialId = credentialId
			} else {
				r.CredentialId = v
			}
		}
		if v, ok := instanceCfg["global_protect_segment_options"].(*schema.Set); ok {
			options, err := expandGlobalProtectSegmentOptionsInstance(v, m)
			if err != nil {
				return nil, err
			}

			r.GlobalProtectSegmentOptions = options
		}
		if v, ok := instanceCfg["enable_traffic"].(bool); ok {
			r.TrafficEnabled = v
		}
		instances[i] = r
	}

	return instances, nil
}

// generate request payload
func generateServicePanRequest(d *schema.ResourceData, m interface{}) (*alkira.ServicePan, error) {

	panoramaDeviceGroup := d.Get("panorama_device_group").(string)
	panoramaIpAddresses := convertTypeListToStringList(d.Get("panorama_ip_addresses").([]interface{}))
	panoramaTemplate := d.Get("panorama_template").(string)

	//
	// Construct instances
	//
	instances, err := expandPanInstances(d.Get("instance").([]interface{}), m)

	if err != nil {
		return nil, err
	}

	//
	// Construct segment options
	//
	segmentOptions, err := expandSegmentOptions(d.Get("segment_options").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	//
	// Construct global protect
	//
	globalProtectSegmentOptions, err := expandGlobalProtectSegmentOptions(d.Get("global_protect_segment_options").(*schema.Set), m)

	if err != nil {
		return nil, err
	}

	service := &alkira.ServicePan{
		BillingTagIds:               convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set)),
		Bundle:                      d.Get("bundle").(string),
		CXP:                         d.Get("cxp").(string),
		CredentialId:                d.Get("pan_credential_id").(string),
		GlobalProtectEnabled:        d.Get("global_protect_enabled").(bool),
		GlobalProtectSegmentOptions: globalProtectSegmentOptions,
		Instances:                   instances,
		LicenseType:                 d.Get("license_type").(string),
		SubLicenseType:              d.Get("license_sub_type").(string),
		MasterKeyCredentialId:       d.Get("pan_master_key_credential_id").(string),
		MasterKeyEnabled:            d.Get("master_key_enabled").(bool),
		MaxInstanceCount:            d.Get("max_instance_count").(int),
		MinInstanceCount:            d.Get("min_instance_count").(int),
		ManagementSegmentId:         d.Get("management_segment_id").(int),
		Name:                        d.Get("name").(string),
		PanoramaEnabled:             d.Get("panorama_enabled").(bool),
		PanoramaDeviceGroup:         &panoramaDeviceGroup,
		PanoramaIpAddresses:         panoramaIpAddresses,
		PanoramaTemplate:            &panoramaTemplate,
		RegistrationCredentialId:    d.Get("pan_registration_credential_id").(string),
		SegmentOptions:              segmentOptions,
		SegmentIds:                  convertTypeSetToIntList(d.Get("segment_ids").(*schema.Set)),
		TunnelProtocol:              d.Get("tunnel_protocol").(string),
		Size:                        d.Get("size").(string),
		Type:                        d.Get("type").(string),
		Version:                     d.Get("version").(string),
		Description:                 d.Get("description").(string),
	}

	return service, nil
}

// Set "instance" blocks from API response
func setPanInstances(d *schema.ResourceData, c []alkira.ServicePanInstance, m interface{}) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, ins := range c {

		// locate the hidden credential from existing Terraform state
		var authKey string
		var authCode string
		var authExpiry string

		for _, value := range d.Get("instance").([]interface{}) {
			cfg := value.(map[string]interface{})
			if cfg["id"].(int) == ins.Id || cfg["name"].(string) == ins.Name {
				authKey = cfg["auth_key"].(string)
				authCode = cfg["auth_code"].(string)
				authExpiry = cfg["auth_expiry"].(string)
			}
		}

		instance := map[string]interface{}{
			"name":                           ins.Name,
			"id":                             ins.Id,
			"auth_key":                       authKey,
			"auth_code":                      authCode,
			"auth_expiry":                    authExpiry,
			"credential_id":                  ins.CredentialId,
			"enable_traffic":                 ins.TrafficEnabled,
			"global_protect_segment_options": flattenGlobalProtectSegmentOptionsInstance(ins.GlobalProtectSegmentOptions, m),
		}
		instances = append(instances, instance)
	}

	return instances
}
