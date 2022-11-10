package alkira

import "github.com/alkiranet/alkira-client-go/alkira"

func expandCiscoFTDvInstances(c []alkira.CheckpointInstance) (alkira.CiscoFTDvInstance, error) {

	var instances []map[string]interface{}

	for _, instance := range c {
		i := map[string]interface{}{
			"name":          instance.Name,
			"credential_id": instance.CredentialId,
		}
		instances = append(instances, i)
	}

	return instances, nil

}
