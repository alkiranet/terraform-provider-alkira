package alkira

import (
	"errors"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandInstanceDeployment converts the input data to an InstanceDeployment struct.
func expandInstanceDeployment(in []interface{}) (*alkira.InstanceDeployment, error) {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] Invalid Instance Deployment")
		return nil, errors.New("Invalid Instance Deployment.")
	}

	deployment := &alkira.InstanceDeployment{}
	for _, v := range in {
		deploymentConfig := v.(map[string]interface{})
		if option, ok := deploymentConfig["deployment_option"]; ok {
			deployment.Option = option.(string)
		}
		if dtype, ok := deploymentConfig["deployment_type"]; ok {
			deployment.Type = dtype.(string)
		}
	}

	return deployment, nil
}

// expandF5Instances converts the input data to a slice of F5Instances structs.
func expandF5Instances(in []interface{}) ([]alkira.F5Instances, error) {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] Invalid F5 Load Balancer instance input")
		return nil, errors.New("Invalid F5 Load Balancer instance input")
	}

	instances := make([]alkira.F5Instances, len(in))

	for i, instance := range in {
		instanceConfig := instance.(map[string]interface{})
		f5lb := alkira.F5Instances{}

		if name, ok := instanceConfig["name"]; ok {
			f5lb.Name = name.(string)
		}
		if licenseType, ok := instanceConfig["license_type"]; ok {
			f5lb.LicenseType = licenseType.(string)
		}
		if version, ok := instanceConfig["version"]; ok {
			f5lb.Version = version.(string)
		}
		if regCredId, ok := instanceConfig["registration_credential_id"]; ok {
			f5lb.RegistrationCredentialId = regCredId.(string)
		}
		if credId, ok := instanceConfig["credential_id"]; ok {
			f5lb.CredentialId = credId.(string)
		}
		if deployment, ok := instanceConfig["deployment"]; ok {
			deploymentStruct, err := expandInstanceDeployment(deployment.([]interface{}))
			if err != nil {
				return nil, err
			}
			f5lb.Deployment = deploymentStruct
		}

		instances[i] = f5lb
	}

	return instances, nil
}

// generateRequestF5Lb generates the request payload for creating an F5 Load Balancer service.
func generateRequestF5Lb(d *schema.ResourceData, m interface{}) (*alkira.ServiceF5Lb, error) {

	billingTagIds := convertTypeSetToIntList(d.Get("billing_tag_ids").(*schema.Set))

	instances, err := expandF5Instances(
		d.Get("instances").([]interface{}),
	)
	if err != nil {
		return nil, err
	}

	// Convert segment IDs to segment names
	segmentNames, err := convertSegmentIdsToSegmentNames(d.Get("segment_ids").(*schema.Set), m)
	if err != nil {
		return nil, err
	}

	service := &alkira.ServiceF5Lb{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Cxp:            d.Get("cxp").(string),
		Size:           d.Get("size").(string),
		Segments:       segmentNames,
		BillingTags:    billingTagIds,
		ElbCidrs:       convertTypeSetToStringList(d.Get("elb_cidrs").(*schema.Set)),
		BigIpAllowList: convertTypeSetToStringList(d.Get("big_ip_allow_list").(*schema.Set)),
		Instances:      instances,
	}

	return service, nil
}
