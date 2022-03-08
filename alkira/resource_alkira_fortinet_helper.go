package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandFortinetInstances(in *schema.Set) []alkira.FortinetInstance {
	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] invalid Fortinet instance input")
		return nil
	}

	instances := make([]alkira.FortinetInstance, in.Len())
	for i, instance := range in.List() {
		r := alkira.FortinetInstance{}
		instanceCfg := instance.(map[string]interface{})
		if v, ok := instanceCfg["name"].(string); ok {
			r.Name = v
			r.HostName = v
		}
		if v, ok := instanceCfg["serial_number"].(string); ok {
			r.SerialNumber = v
		}
		if v, ok := instanceCfg["credential_id"].(string); ok {
			r.CredentialId = v
		}
		instances[i] = r
	}

	return instances
}
