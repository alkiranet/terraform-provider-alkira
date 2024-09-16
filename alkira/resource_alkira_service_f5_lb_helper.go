package alkira

import (
	"errors"
	"log"

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

		if name, ok := tfInstance["name"]; ok {

			instanceStruct.Name = name.(string)
		}
		if licenseType, ok := tfInstance["license_type"]; ok {
			instanceStruct.LicenseType = licenseType.(string)
			// if the license type is BRING_YOUR_OWN, we need a registration credential.
			if licenseType == "BRING_YOUR_OWN" {
				if regCredId, ok := tfInstance["registration_credential_id"].(string); ok {
					if regCredId == "" {
						credentialName := instanceStruct.Name + "registration" + randomNameSuffix()
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
		if deploymentOption, ok := tfInstance["deployment_option"]; ok {
			instanceDeployment.Option = deploymentOption.(string)
		}
		instanceStruct.Deployment = instanceDeployment

		if credId, ok := tfInstance["credential_id"].(string); ok {

			if credId == "" {
				credentialName := instanceStruct.Name + randomNameSuffix()
				credentialF5Instance := alkira.CredentialF5Instance{
					// UserName: tfInstance["f5_username"].(string),
					// hardcode the username to admin for now.
					UserName: "admin",
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

		segmentOptions[segmentName] = subOption
	}

	return segmentOptions, nil
}

func deflateF5SegmentOptions(in alkira.F5SegmentOption, m interface{}) ([]map[string]interface{}, error) {
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

		segmentOptions = append(segmentOptions, option)
	}

	return segmentOptions, nil
}

func setF5Instances(d *schema.ResourceData, c []alkira.F5Instance) []map[string]interface{} {
	var instances []map[string]interface{}

	for _, instance := range d.Get("instances").([]interface{}) {
		tfInstance := instance.(map[string]interface{})

		for _, apiInstance := range c {
			if tfInstance["id"].(int) == apiInstance.Id || tfInstance["name"].(string) == apiInstance.Name {
				instanceStruct := map[string]interface{}{
					"name":                       apiInstance.Name,
					"id":                         apiInstance.Id,
					"credential_id":              apiInstance.CredentialId,
					"registration_credential_id": apiInstance.RegistrationCredentialId,
					"license_type":               apiInstance.LicenseType,
					"version":                    apiInstance.Version,
					"hostname_fqdn":              apiInstance.HostNameFqdn,
					"deployment_option":          apiInstance.Deployment.Option,
					"deployment_type":            apiInstance.Deployment.Type,
					"f5_username":                "admin",
					"f5_password":                tfInstance["f5_password"].(string),
				}
				instances = append(instances, instanceStruct)
				break
			}
		}
	}

	for _, apiInstance := range c {
		new := true

		for _, instance := range d.Get("instances").([]interface{}) {
			configInstance := instance.(map[string]interface{})

			if configInstance["id"].(int) == apiInstance.Id || configInstance["name"].(string) == apiInstance.Name {
				new = false
				break
			}
		}
		if new {
			instanceStruct := map[string]interface{}{
				"credential_id":              apiInstance.CredentialId,
				"registration_credential_id": apiInstance.RegistrationCredentialId,
				"license_type":               apiInstance.LicenseType,
				"version":                    apiInstance.Version,
				"hostname_fqdn":              apiInstance.HostNameFqdn,
				"deployment_option":          apiInstance.Deployment.Option,
				"deployment_type":            apiInstance.Deployment.Type,
			}
			instances = append(instances, instanceStruct)
			break
		}

	}
	return instances
}

// generateRequestF5Lb generates the request payload for creating an F5 Load Balancer service.
func generateRequestF5Lb(d *schema.ResourceData, m interface{}) (*alkira.ServiceF5Lb, error) {

	billingTagIds := convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set))

	instances, err := expandF5Instances(
		d.Get("instances").([]interface{}), m)
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
