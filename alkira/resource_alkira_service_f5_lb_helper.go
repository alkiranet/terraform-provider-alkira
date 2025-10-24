package alkira

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandF5Instances converts the input data to a slice of F5Instances structs.
func expandF5Instances(in []interface{}, m interface{}) ([]alkira.F5Instance, error) {

	client := m.(*alkira.AlkiraClient)

	if in == nil || len(in) == 0 {
		log.Printf("[ERROR] Invalid F5 Load Balancer instance input.")
		return nil, errors.New("Invalid F5 Load Balancer instance input.")
	}

	instances := make([]alkira.F5Instance, len(in))
	instanceDeployment := alkira.F5InstanceDeployment{}

	for i, instance := range in {
		tfInstance := instance.(map[string]interface{})
		instanceStruct := alkira.F5Instance{}
		if instanceId, ok := tfInstance["id"]; ok {

			instanceStruct.Id = instanceId.(int)
		}
		if name, ok := tfInstance["name"]; ok {

			instanceStruct.Name = name.(string)
		}
		if licenseType, ok := tfInstance["license_type"]; ok {
			instanceStruct.LicenseType = licenseType.(string)
			// if the license type is BRING_YOUR_OWN, we need a registration credential.
			if licenseType == "BRING_YOUR_OWN" {
				if rawRegCredId, ok := tfInstance["registration_credential_id"]; ok {
					regCredId := rawRegCredId.(string)
					if regCredId == "" {
						credentialName := instanceStruct.Name + "registration" + randomNameSuffix()
						if len(instanceStruct.Name) > 255 {
							credentialName = instanceStruct.Name[0:225] + "registration" + randomNameSuffix()
						}
						credentialF5Registration := alkira.CredentialF5InstanceRegistration{
							RegistrationKey: tfInstance["f5_registration_key"].(string),
						}

						log.Printf("[INFO] Creating F5 Load Balancer Instance Registration Credential %s", credentialName)
						credentialId, err := client.CreateCredential(
							credentialName,
							alkira.CredentialTypeF5InstanceRegistration,
							credentialF5Registration,
							0,
						)
						if err != nil {
							return nil, err
						}
						instanceStruct.RegistrationCredentialId = credentialId
					} else {
						instanceStruct.RegistrationCredentialId = regCredId
					}
				}
			}
		}
		if version, ok := tfInstance["version"]; ok {
			instanceStruct.Version = version.(string)
		}
		if fqdn, ok := tfInstance["hostname_fqdn"]; ok {
			instanceStruct.HostNameFqdn = fqdn.(string)
		}

		if deploymentType, ok := tfInstance["deployment_type"]; ok {
			instanceDeployment.Type = deploymentType.(string)
		}
		instanceStruct.Deployment = instanceDeployment

		if rawCredId, ok := tfInstance["credential_id"]; ok {
			credId := rawCredId.(string)

			if credId == "" {
				credentialName := instanceStruct.Name + "credential" + randomNameSuffix()
				if len(instanceStruct.Name) > 255 {
					credentialName = instanceStruct.Name[0:225] + "credential" + randomNameSuffix()
				}
				credentialF5Instance := alkira.CredentialF5Instance{
					UserName: tfInstance["f5_username"].(string),
					Password: tfInstance["f5_password"].(string),
				}

				log.Printf("[INFO] Creating F5 Load Balancer Instance Credential %s", credentialName)
				credentialId, err := client.CreateCredential(
					credentialName,
					alkira.CredentialTypeF5Instance,
					credentialF5Instance,
					0,
				)

				if err != nil {
					return nil, err
				}

				instanceStruct.CredentialId = credentialId
			} else {
				instanceStruct.CredentialId = credId
			}
		}
		if availabilityZone, ok := tfInstance["availability_zone"]; ok && strings.TrimSpace(availabilityZone.(string)) != "" {
			instanceStruct.AvailabilityZone = json.Number(availabilityZone.(string))
		}
		instances[i] = instanceStruct
	}

	return instances, nil
}

func expandF5SegmentOptions(in *schema.Set, m interface{}) (alkira.F5SegmentOption, error) {
	segmentOptions := make(alkira.F5SegmentOption)

	for _, option := range in.List() {
		cfg := option.(map[string]interface{})
		segmentId := cfg["segment_id"].(string)

		segmentName, err := getSegmentNameById(segmentId, m)
		if err != nil {
			return nil, err
		}

		subOption := alkira.F5SegmentSubOption{
			ElbNicCount: cfg["elb_nic_count"].(int),
		}

		if natPoolPrefixLength, ok := cfg["nat_pool_prefix_length"]; ok {
			subOption.NatPoolPrefixLength = natPoolPrefixLength.(int)
		}

		if advertiseToCxpPrefixListId, ok := cfg["elb_bgp_options_advertise_to_cxp_prefix_list_id"]; ok && advertiseToCxpPrefixListId != 0 {
			bgpOptions := &alkira.ElbBgpOptions{}
			bgpOptions.AdvertiseToCXPPrefixListId = advertiseToCxpPrefixListId.(int)
			subOption.ElbBgpOptions = bgpOptions
		}

		segmentOptions[segmentName] = subOption
	}

	return segmentOptions, nil
}

func setF5SegmentOptions(in alkira.F5SegmentOption, m interface{}) ([]map[string]interface{}, error) {
	if in == nil {
		return nil, errors.New("[ERROR] Segment options is nil.")
	}

	var segmentOptions []map[string]interface{}

	for segmentName, subOption := range in {
		segmentId, err := getSegmentIdByName(segmentName, m)
		if err != nil {
			return nil, err
		}

		option := map[string]interface{}{
			"segment_id":    segmentId,
			"elb_nic_count": subOption.ElbNicCount,
		}

		if subOption.NatPoolPrefixLength != 0 {
			option["nat_pool_prefix_length"] = subOption.NatPoolPrefixLength
		}

		if subOption.ElbBgpOptions != nil {
			option["elb_bgp_options_advertise_to_cxp_prefix_list_id"] = subOption.ElbBgpOptions.AdvertiseToCXPPrefixListId
		}

		segmentOptions = append(segmentOptions, option)
	}

	return segmentOptions, nil
}

// Set `instance` block from the API response, except the creds.
func setF5Instances(d *schema.ResourceData, ins []alkira.F5Instance) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, in := range ins {
		// fetch the creds from the terraform state.
		f5Username := ""
		f5Password := ""
		f5RegistrationKey := ""
		for _, value := range d.Get("instance").([]interface{}) {
			cfg := value.(map[string]interface{})
			if cfg["id"].(int) == in.Id || cfg["name"].(string) == in.Name {
				f5Username = cfg["f5_username"].(string)
				f5Password = cfg["f5_password"].(string)
				f5RegistrationKey = cfg["f5_registration_key"].(string)
			}
		}
		instance := map[string]interface{}{
			"name":                       in.Name,
			"id":                         in.Id,
			"license_type":               in.LicenseType,
			"registration_credential_id": in.RegistrationCredentialId,
			"credential_id":              in.CredentialId,
			"version":                    in.Version,
			"deployment_type":            in.Deployment.Type,
			"hostname_fqdn":              in.HostNameFqdn,
			"f5_registration_key":        f5RegistrationKey,
			"f5_username":                f5Username,
			"f5_password":                f5Password,
			"availability_zone":          in.AvailabilityZone,
		}

		instances = append(instances, instance)
	}
	return instances
}

// generateRequestF5Lb generates the request payload for creating an F5 Load Balancer service.
func generateRequestF5Lb(d *schema.ResourceData, m interface{}) (*alkira.ServiceF5Lb, error) {

	billingTagIds := convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set))

	instances, err := expandF5Instances(
		d.Get("instance").([]interface{}), m)
	if err != nil {
		return nil, err
	}

	// Convert segment IDs to segment names
	segmentNames, err := convertSegmentIdsToSegmentNames(d.Get("segment_ids").(*schema.Set), m)
	if err != nil {
		return nil, err
	}

	segmentOptions, err := expandF5SegmentOptions(d.Get("segment_options").(*schema.Set), m)
	if err != nil {
		return nil, err
	}

	service := &alkira.ServiceF5Lb{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		Cxp:              d.Get("cxp").(string),
		Size:             d.Get("size").(string),
		ServiceGroupName: d.Get("service_group_name").(string),
		Segments:         segmentNames,
		BillingTags:      billingTagIds,
		Instances:        instances,
		SegmentOptions:   segmentOptions,
		PrefixListId:     d.Get("prefix_list_id").(int),
		GlobalCidrListId: d.Get("global_cidr_list_id").(int),
	}
	return service, nil
}
