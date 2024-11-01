package alkira

import (
	"errors"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type panZone struct {
	Segment string
	Zone    string
	Groups  interface{}
}

// Helper functions for credentials. There are 3 types of credentials used in PAN:
//
// - PAN credential
// - PAN Registration
// - PAN Master Key
func createPanCredential(d *schema.ResourceData, client *alkira.AlkiraClient) (string, error) {
	log.Printf("[INFO] Creating PAN Credential")

	panCredName := d.Get("name").(string) + randomNameSuffix()
	panCredential := alkira.CredentialPan{
		Username:   d.Get("pan_username").(string),
		Password:   d.Get("pan_password").(string),
		LicenseKey: d.Get("pan_license_key").(string),
	}

	return client.CreateCredential(
		panCredName,
		alkira.CredentialTypePan,
		panCredential,
		0,
	)

}

func deletePanCredential(credentialId string, client *alkira.AlkiraClient) error {
	log.Printf("[INFO] Deleting PAN Credential")

	return client.DeleteCredential(credentialId, alkira.CredentialTypePan)
}

func createPanRegistrationCredential(d *schema.ResourceData, client *alkira.AlkiraClient) (string, error) {
	log.Printf("[INFO] Creating PAN Registration Credential")

	credentialName := d.Get("name").(string) + randomNameSuffix()
	credential := alkira.CredentialPanRegistration{
		RegistrationPinId:    d.Get("registration_pin_id").(string),
		RegistrationPinValue: d.Get("registration_pin_value").(string),
	}
	credentialExpiry, err := convertInputTimeToEpoch(d.Get("registration_pin_expiry").(string))

	if err != nil {
		log.Printf("[ERROR] failed to parse 'registration_pin_exiry', %v", err)
		return "", err
	}

	return client.CreateCredential(credentialName, alkira.CredentialTypePanRegistration, credential, credentialExpiry)
}

func deletePanRegistrationCredential(credentialId string, client *alkira.AlkiraClient) error {
	log.Printf("[INFO] Deleting PAN Registration Credential")

	return client.DeleteCredential(credentialId, alkira.CredentialTypePanRegistration)
}

func createPanMasterKeyCredential(d *schema.ResourceData, client *alkira.AlkiraClient) (string, error) {
	log.Printf("[INFO] Creating PAN Master Key Credential")

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

	return client.CreateCredential(credentialName, alkira.CredentialTypePanMasterKey, credential, credentialExpiry)
}

func deletePanMasterKeyCredential(credentialId string, client *alkira.AlkiraClient) error {
	log.Printf("[INFO] Deleting PAN Master Key Credential")

	return client.DeleteCredential(credentialId, alkira.CredentialTypePanMasterKey)
}

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

func expandPanSegmentOptions(in *schema.Set, m interface{}) (map[string]interface{}, error) {

	if in == nil {
		return nil, errors.New("invalid SegmentOptions input")
	}

	zoneMap := make([]panZone, in.Len())

	for i, option := range in.List() {
		r := panZone{}
		cfg := option.(map[string]interface{})
		if v, ok := cfg["segment_id"].(string); ok {
			segmentName, err := getSegmentNameById(v, m)

			if err != nil {
				return nil, err
			}
			r.Segment = segmentName
		}
		if v, ok := cfg["zone_name"].(string); ok {
			r.Zone = v
		}

		r.Groups = cfg["groups"]

		zoneMap[i] = r
	}

	segmentOptions := make(map[string]interface{})

	for _, x := range zoneMap {
		zone := make(map[string]interface{})
		zone[x.Zone] = x.Groups

		for _, y := range zoneMap {
			if x.Segment == y.Segment {
				zone[y.Zone] = y.Groups
			}
		}

		zonesToGroups := make(map[string]interface{})
		zonesToGroups["zonesToGroups"] = zone

		segmentOptions[x.Segment] = zonesToGroups
	}

	return segmentOptions, nil
}

func expandPanInstances(in []interface{}, m interface{}) ([]alkira.ServicePanInstance, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		return nil, errors.New("Invalid PAN instance input")
	}

	instances := make([]alkira.ServicePanInstance, len(in))
	for i, instance := range in {
		r := alkira.ServicePanInstance{}
		instanceCfg := instance.(map[string]interface{})

		var AuthCode string
		var AuthKey string

		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := instanceCfg["auth_code"].(string); ok {
			AuthCode = v
		}
		if v, ok := instanceCfg["auth_key"].(string); ok {
			AuthKey = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			if v == "" {
				credentialName := r.Name + randomNameSuffix()
				credentialPanInstance := alkira.CredentialPanInstance{
					AuthCode: AuthCode,
					AuthKey:  AuthKey,
				}

				log.Printf("[INFO] Creating PAN Instance Credential %s", credentialName)
				credentialId, err := client.CreateCredential(
					credentialName,
					alkira.CredentialTypePanInstance,
					credentialPanInstance,
					0,
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

func setPanInstances(d *schema.ResourceData, c []alkira.ServicePanInstance) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, value := range d.Get("instance").([]interface{}) {
		cfg := value.(map[string]interface{})

		for _, ins := range c {
			if cfg["id"].(int) == ins.Id || cfg["name"].(string) == ins.Name {
				instance := map[string]interface{}{
					"name":           ins.Name,
					"id":             ins.Id,
					"credential_id":  ins.CredentialId,
					"auth_key":       cfg["auth_key"].(string),
					"auth_code":      cfg["auth_code"].(string),
					"enable_traffic": ins.TrafficEnabled,
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

			if cfg["id"].(int) == instance.Id || cfg["name"].(string) == instance.Name {
				new = false
				break
			}
		}

		// If the instance is new, add it to the tail of the list,
		// this will generate a diff
		if new {
			instance := map[string]interface{}{
				"credential_id":  instance.CredentialId,
				"name":           instance.Name,
				"id":             instance.Id,
				"enable_traffic": instance.TrafficEnabled,
			}

			instances = append(instances, instance)
			break
		}
	}

	return instances
}
