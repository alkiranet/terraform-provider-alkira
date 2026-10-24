package alkira

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// expandMerakiSDWANSegmentOptions expand Meraki SD-WAN target segment
func expandMerakiSDWANSegmentOptions(in []interface{}) []alkira.ConnectorMerakiSDWANSegmentOption {

	if len(in) == 0 {
		log.Printf("[DEBUG] Empty target_segment")
		return []alkira.ConnectorMerakiSDWANSegmentOption{}
	}

	options := make([]alkira.ConnectorMerakiSDWANSegmentOption, len(in))

	for i, option := range in {
		r := alkira.ConnectorMerakiSDWANSegmentOption{}
		t := option.(map[string]interface{})

		if v, ok := t["advertise_on_prem_routes"].(bool); ok {
			r.AdvertiseOnPremRoutes = v
		}
		if v, ok := t["advertise_default_route"].(bool); ok {
			r.DisableInternetExit = !v
		}
		if v, ok := t["meraki_bgp_asn"].(int); ok {
			r.MerakiBgpAsn = int64(v)
		}
		if v, ok := t["segment_id"].(int); ok {
			r.SegmentId = v
		}

		options[i] = r
	}

	return options
}

// merakiSDWANInstanceHash identify an instance by its host_name, which is
// unique in the connector and can't change after provisioning. auth_token
// is in the hash too, because a set shows no change for a field that is
// not in the hash.
func merakiSDWANInstanceHash(v interface{}) int {
	m := v.(map[string]interface{})
	hostName, _ := m["host_name"].(string)
	authToken, _ := m["auth_token"].(string)
	return schema.HashString(hostName + "\x00" + authToken)
}

// expandMerakiSDWANInstances expand Meraki SD-WAN connector instances. An
// instance is matched to the prior state by its host_name. A new instance
// gets a credential created from its auth_token. An existing instance
// keeps its ID and credential; when its auth_token differs from the prior
// state, the credential is updated in place, keeping its ID and name.
// After an import the prior token is empty, so the configured token is
// written once. A host_name used twice fails before any credential call.
// A token can be used once, so a token shared by two instances, or a new
// instance with the token of an existing one, also fails before it.
//
// It returns the IDs of the credentials it created and the instances whose
// token it replaced, also when it fails part way. If it fails, the
// credentials created here are deleted.
func expandMerakiSDWANInstances(ac *alkira.AlkiraClient, in []interface{}, old []interface{}) ([]alkira.ConnectorMerakiSDWANInstance, []string, []string, error) {

	if len(in) == 0 {
		log.Printf("[DEBUG] Empty instances")
		return []alkira.ConnectorMerakiSDWANInstance{}, nil, nil, nil
	}

	// Instances from the prior state, keyed by host_name, and their tokens
	oldInstances := make(map[string]map[string]interface{})
	oldTokens := make(map[string]string)

	for _, inst := range old {
		o, ok := inst.(map[string]interface{})
		if !ok {
			continue
		}

		hostName, _ := o["host_name"].(string)
		oldInstances[hostName] = o

		if token, _ := o["auth_token"].(string); token != "" {
			oldTokens[token] = hostName
		}
	}

	// Check the host names and tokens before any credential call, so a
	// rejected apply leaves the credentials unchanged
	hostNames := make(map[string]bool)
	tokens := make(map[string]string)

	for _, inst := range in {
		t := inst.(map[string]interface{})
		hostName, _ := t["host_name"].(string)
		token, _ := t["auth_token"].(string)

		if hostNames[hostName] {
			return nil, nil, nil, fmt.Errorf("instances: host_name %q is used by more than one instance. "+
				"Each instance needs a unique host_name", hostName)
		}
		hostNames[hostName] = true

		if other, found := tokens[token]; found {
			return nil, nil, nil, fmt.Errorf("instances: %q and %q have the same auth_token. "+
				"A token can be used once: generate a token for each instance", other, hostName)
		}
		tokens[token] = hostName

		if _, found := oldInstances[hostName]; found {
			continue
		}

		if other, found := oldTokens[token]; found {
			return nil, nil, nil, fmt.Errorf("instances: the new instance %q has the auth_token of instance %q. "+
				"A token can be used once: generate a new token for %q. "+
				"A new host_name adds a new vMX", hostName, other, hostName)
		}
	}

	instances := make([]alkira.ConnectorMerakiSDWANInstance, len(in))
	var created, replaced []string

	for i, inst := range in {
		r := alkira.ConnectorMerakiSDWANInstance{}
		t := inst.(map[string]interface{})

		var authToken string

		if v, ok := t["host_name"].(string); ok {
			r.HostName = v
		}
		if v, ok := t["auth_token"].(string); ok {
			authToken = v
		}

		o, found := oldInstances[r.HostName]

		if !found {
			log.Printf("[DEBUG] Creating the credential of Meraki SD-WAN instance %s", r.HostName)
			credentialId, err := ac.CreateCredential(
				r.HostName+randomNameSuffix(),
				alkira.CredentialTypeMerakiSDWANInstance,
				alkira.CredentialMerakiSDWANInstance{AuthToken: authToken},
				0,
			)

			if err != nil {
				deleteMerakiSDWANCredentials(ac, created)
				return nil, nil, replaced, fmt.Errorf("failed to create the credential of instance %s: %w", r.HostName, err)
			}

			created = append(created, credentialId)
			r.CredentialId = credentialId
		} else {
			r.CredentialId, _ = o["credential_id"].(string)

			if v, ok := o["id"].(int); ok && v != 0 {
				r.Id = json.Number(fmt.Sprintf("%d", v))
			}

			if o["auth_token"] != authToken {
				log.Printf("[DEBUG] Replacing the auth_token of Meraki SD-WAN instance %s (credential %s)", r.HostName, r.CredentialId)

				if err := updateMerakiSDWANInstanceCredential(ac, r.CredentialId, authToken); err != nil {
					deleteMerakiSDWANCredentials(ac, created)
					return nil, nil, replaced, fmt.Errorf("failed to replace the auth_token of instance %s: %w", r.HostName, err)
				}

				replaced = append(replaced, fmt.Sprintf("%s (credential %s)", r.HostName, r.CredentialId))
			}
		}

		instances[i] = r
	}

	return instances, created, replaced, nil
}

// updateMerakiSDWANInstanceCredential update the token of an existing
// instance credential, keeping its name
func updateMerakiSDWANInstanceCredential(ac *alkira.AlkiraClient, credentialId string, authToken string) error {

	existing, err := ac.GetCredentialById(credentialId)

	if err != nil {
		return fmt.Errorf("failed to get Meraki SD-WAN instance credential %s: %w", credentialId, err)
	}

	err = ac.UpdateCredential(
		credentialId,
		existing.Name,
		alkira.CredentialTypeMerakiSDWANInstance,
		alkira.CredentialMerakiSDWANInstance{AuthToken: authToken},
		0,
	)

	if err != nil {
		return fmt.Errorf("failed to update Meraki SD-WAN instance credential %s: %w", credentialId, err)
	}

	return nil
}

// merakiSDWANReplacedTokensError say which tokens were already replaced
// when the update fails after replacing them. The credentials hold the new
// tokens, while the state keeps the old ones until an update succeeds.
func merakiSDWANReplacedTokensError(replaced []string, err error) error {
	if len(replaced) == 0 {
		return err
	}

	return fmt.Errorf("the auth_token was replaced for %s, but the connector update failed; apply again to finish: %w",
		strings.Join(replaced, ", "), err)
}

// deleteMerakiSDWANCredentials delete the credentials created in this
// apply, when the apply fails before a connector references them
func deleteMerakiSDWANCredentials(ac *alkira.AlkiraClient, ids []string) {
	for _, id := range ids {
		if err := ac.DeleteCredential(id, alkira.CredentialTypeMerakiSDWANInstance); err != nil {
			log.Printf("[WARN] failed to delete Meraki SD-WAN instance credential %s: %v", id, err)
		}
	}
}

// setMerakiSDWANInstances set the instances from the API response.
// auth_token is not returned by the API, so it is kept from the state of
// the instance with the same host_name. An instance not in the state gets
// an empty token, as after an import.
func setMerakiSDWANInstances(d *schema.ResourceData, connector *alkira.ConnectorMerakiSDWAN) {
	tokens := make(map[string]interface{})

	for _, inst := range d.Get("instances").(*schema.Set).List() {
		i := inst.(map[string]interface{})
		tokens[i["host_name"].(string)] = i["auth_token"]
	}

	var instances []map[string]interface{}

	for _, info := range connector.Instances {
		authToken, found := tokens[info.HostName]
		if !found {
			authToken = ""
		}

		id, _ := info.Id.Int64()
		instances = append(instances, map[string]interface{}{
			"auth_token":    authToken,
			"az":            info.Az,
			"credential_id": info.CredentialId,
			"host_name":     info.HostName,
			"id":            int(id),
			"name":          info.Name,
			"public_ip":     info.PublicIp,
			"wan_ip":        info.WanIp,
		})
	}

	d.Set("instances", instances)
}

// validateMerakiSDWANBgpAsn check the ASN is in 1-4294967295. The range
// is checked as int64 because int is 32-bit on some provider builds.
func validateMerakiSDWANBgpAsn(i interface{}, k string) ([]string, []error) {
	v, ok := i.(int)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %s to be int", k)}
	}

	if int64(v) < 1 || int64(v) > 4294967295 {
		return nil, []error{fmt.Errorf("expected %s to be in the range (1 - 4294967295), got %d", k, v)}
	}

	return nil, nil
}
