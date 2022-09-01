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

func expandGlobalProtectSegmentOptions(in *schema.Set, m interface{}) (map[string]*alkira.GlobalProtectSegmentName, error) {
	client := m.(*alkira.AlkiraClient)

	if in == nil || in.Len() == 0 {
		return nil, nil
	}

	sgmtOptions := make(map[string]*alkira.GlobalProtectSegmentName)
	for _, sgmtOption := range in.List() {
		r := &alkira.GlobalProtectSegmentName{}
		segmentCfg := sgmtOption.(map[string]interface{})
		var segmentName string

		if v, ok := segmentCfg["segment_id"].(string); ok {
			segment, err := client.GetSegmentById(v)
			if err != nil {
				return nil, err
			}
			segmentName = segment.Name
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
	client := m.(*alkira.AlkiraClient)

	if in == nil || in.Len() == 0 {
		return nil, nil
	}

	sgmtOptions := make(map[string]*alkira.GlobalProtectSegmentNameInstance)
	for _, sgmtOption := range in.List() {
		r := &alkira.GlobalProtectSegmentNameInstance{}
		segmentCfg := sgmtOption.(map[string]interface{})
		var segmentName string

		if v, ok := segmentCfg["segment_id"].(string); ok {
			segment, err := client.GetSegmentById(v)
			if err != nil {
				return nil, err
			}
			segmentName = segment.Name
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
	client := m.(*alkira.AlkiraClient)

	if in == nil {
		return nil, errors.New("invalid SegmentOptions input")
	}

	zoneMap := make([]panZone, in.Len())

	for i, option := range in.List() {
		r := panZone{}
		cfg := option.(map[string]interface{})
		if v, ok := cfg["segment_id"].(string); ok {
			segment, err := client.GetSegmentById(v)
			if err != nil {
				return nil, err
			}
			r.Segment = segment.Name
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
	if in == nil || len(in) == 0 {
		return nil, errors.New("Invalid PAN instance input")
	}

	instances := make([]alkira.ServicePanInstance, len(in))
	for i, instance := range in {
		r := alkira.ServicePanInstance{}
		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["id"].(int); ok {
			r.Id = v
		}
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			r.CredentialId = v
		}
		if v, ok := instanceCfg["global_protect_segment_options"].(*schema.Set); ok {
			options, err := expandGlobalProtectSegmentOptionsInstance(v, m)
			if err != nil {
				return nil, err
			}

			r.GlobalProtectSegmentOptions = options
		}
		instances[i] = r
	}

	return instances, nil
}

// updateOrCreatePanInstanceCreds
// updates and creates pan instance creds, checks if exists first then saves and returns credential ids created in an array of strings.
func updateOrCreatePanInstanceCreds(client *alkira.AlkiraClient, in []interface{}, allCreds []alkira.CredentialResponseDetail) ([]string, error) {
	if in == nil || len(in) == 0 {
		return nil, errors.New("Invalid PAN instance input")
	}
	credItems := make([]string, len(in))
	var err error
	for i, instance := range in {
		r := alkira.CredentialPanInstance{}

		instanceCfg := instance.(map[string]interface{})

		if v, ok := instanceCfg["auth_code"].(string); ok {
			r.AuthCode = v
		}
		if v, ok := instanceCfg["auth_key"].(string); ok {
			r.AuthKey = v
		}

		panCredInName := instanceCfg["name"].(string) + randomNameSuffix()
		if v, ok := instanceCfg["credential_id"].(string); 0 != len(v) && ok {
			var found bool
			for _, g := range allCreds {
				if g.Id == v {
					found = true
					err := client.UpdateCredential(v, panCredInName, alkira.CredentialTypePanInstance, r, 0)
					if err != nil {
						log.Printf("[ERROR] failed to update Pan instance credential, %v", err)
						return nil, err
					}
					credItems[i] = v
				}
			}
			if !found {
				credItems[i], err = client.CreateCredential(panCredInName, alkira.CredentialTypePanInstance, r, 0)
				if err != nil {
					log.Printf("[ERROR] failed to create Pan Instance Credentials, %v", err)
					return nil, err
				}
			}

		} else {
			credId, err := client.CreateCredential(panCredInName, alkira.CredentialTypePanInstance, r, 0)
			if err != nil {
				log.Printf("[ERROR] failed to create Pan Instance Credentials, %v", err)
				return nil, err
			}
			credItems[i] = credId
		}
	}
	return credItems, nil
}

// updateOrCreatePanCred updates or creates pan credential. Checks if exists first.
func updateOrCreatePanCred(client *alkira.AlkiraClient, d *schema.ResourceData, allCreds []alkira.CredentialResponseDetail) (string, error) {

	panCredentialId := d.Get("credential_id").(string)
	panCredName := d.Get("name").(string) + randomNameSuffix()
	panCredential := alkira.CredentialPan{
		Username: d.Get("username").(string),
		Password: d.Get("password").(string),
	}
	var err error
	if 0 != len(panCredentialId) {

		for _, g := range allCreds {
			if g.Id == panCredentialId {
				err = client.UpdateCredential(panCredentialId, panCredName, alkira.CredentialTypePan, panCredential, 0)
				if err != nil {
					log.Printf("[ERROR] failed to update Pan credential, %v", err)
					return panCredentialId, err
				}
				return panCredentialId, nil
			}
		}
	}

	panCredentialId, err = client.CreateCredential(panCredName, alkira.CredentialTypePan, panCredential, 0)

	if err != nil {
		log.Printf("[ERROR] failed to create PAN credentials, %v", err)
		return panCredentialId, err
	}

	return panCredentialId, nil
}
