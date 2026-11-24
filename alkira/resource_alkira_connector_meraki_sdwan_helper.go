package alkira

import (
	"encoding/json"
	"fmt"
	"log"
	"slices"
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

// suppressMerakiSDWANImportedToken hide the auth_token diff of an imported
// instance. The API doesn't return the token, so the state token is empty
// after an import while the instance has a credential. The diff is shown
// when another host_name takes the place of the instance in the list: that
// is a new instance and its token must be sent.
func suppressMerakiSDWANImportedToken(k, old, new string, d *schema.ResourceData) bool {
	prefix := k[:strings.LastIndex(k, ".")+1]
	oldHostName, newHostName := d.GetChange(prefix + "host_name")
	if oldHostName != newHostName {
		return false
	}

	return suppressCredentialSecretDiff("credential_id", true)(k, old, new, d)
}

// expandMerakiSDWANInstances expand Meraki SD-WAN connector instances. An
// instance is matched to the prior state by its host_name and keeps its ID.
// The instance credential can't be updated, so a new instance or a changed
// auth_token gets a new credential. An unchanged auth_token keeps the
// credential. An empty prior token (after an import, or a credential
// changed outside the state) keeps the current credential, and the
// host_name is returned in adopted. A host_name
// used twice fails before any credential call. A token can be used once,
// so a token shared by two instances, or an instance with the prior token
// of another instance, also fails before it.
//
// It returns the IDs of the credentials it created. If it fails, the
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

		if other, found := tokens[token]; found && token != "" {
			return nil, nil, nil, fmt.Errorf("instances: %q and %q have the same auth_token. "+
				"A token can be used once: generate a token for each instance", other, hostName)
		}
		tokens[token] = hostName

		if other, found := oldTokens[token]; found && other != hostName {
			return nil, nil, nil, fmt.Errorf("instances: %q has the auth_token that instance %q had. "+
				"A token can be used once: generate a new token for %q. "+
				"A new host_name adds a new vMX", hostName, other, hostName)
		}
	}

	instances := make([]alkira.ConnectorMerakiSDWANInstance, len(in))
	var created, adopted []string

	for i, inst := range in {
		r := alkira.ConnectorMerakiSDWANInstance{}
		t := inst.(map[string]interface{})

		r.HostName, _ = t["host_name"].(string)
		authToken, _ := t["auth_token"].(string)

		if o, found := oldInstances[r.HostName]; found {
			if v, ok := o["id"].(int); ok && v != 0 {
				r.Id = json.Number(fmt.Sprintf("%d", v))
			}

			oldToken, _ := o["auth_token"].(string)
			oldCredentialId, _ := o["credential_id"].(string)

			switch {
			case oldCredentialId == "":
			case oldToken == authToken:
				r.CredentialId = oldCredentialId
			case oldToken == "":
				r.CredentialId = oldCredentialId
				adopted = append(adopted, r.HostName)
			}
		}

		if r.CredentialId == "" {
			log.Printf("[DEBUG] Creating a credential for Meraki SD-WAN instance %s", r.HostName)
			credentialId, err := ac.CreateCredential(
				r.HostName+randomNameSuffix(),
				alkira.CredentialTypeMerakiSDWANInstance,
				alkira.CredentialMerakiSDWANInstance{AuthToken: authToken},
				0,
			)

			if err != nil {
				deleteMerakiSDWANCredentials(ac, created)
				return nil, nil, nil, fmt.Errorf("failed to create the credential of instance %s: %w", r.HostName, err)
			}

			created = append(created, credentialId)
			r.CredentialId = credentialId
		}

		instances[i] = r
	}

	return instances, created, adopted, nil
}

// merakiSDWANSavedCredentials return the credential ID of each saved
// instance, keyed by host_name
func merakiSDWANSavedCredentials(api *alkira.AlkiraAPI[alkira.ConnectorMerakiSDWAN], id string) (map[string]string, error) {
	connector, _, err := api.GetById(id)
	if err != nil {
		return nil, err
	}

	saved := make(map[string]string, len(connector.Instances))
	for _, inst := range connector.Instances {
		saved[inst.HostName] = inst.CredentialId
	}

	return saved, nil
}

// merakiSDWANSkippedInstances return the host names whose sent credential
// is not the saved one. With provisioning on, a connector in FAILED state
// is not updated: the backend provisions the saved configuration instead.
func merakiSDWANSkippedInstances(request *alkira.ConnectorMerakiSDWAN, saved map[string]string) []string {
	var skipped []string
	for _, inst := range request.Instances {
		if saved[inst.HostName] != inst.CredentialId {
			skipped = append(skipped, inst.HostName)
		}
	}
	return skipped
}

// deleteUnusedMerakiSDWANCredentials delete the created credentials that no
// saved instance uses
func deleteUnusedMerakiSDWANCredentials(ac *alkira.AlkiraClient, created []string, saved map[string]string) {
	used := make(map[string]bool, len(saved))
	for _, credentialId := range saved {
		used[credentialId] = true
	}

	var unused []string
	for _, id := range created {
		if !used[id] {
			unused = append(unused, id)
		}
	}

	deleteMerakiSDWANCredentials(ac, unused)
}

// deleteMerakiSDWANCredentials delete the given credentials created in this
// apply. A failed delete is only logged.
func deleteMerakiSDWANCredentials(ac *alkira.AlkiraClient, ids []string) {
	for _, id := range ids {
		if err := ac.DeleteCredential(id, alkira.CredentialTypeMerakiSDWANInstance); err != nil {
			log.Printf("[WARN] failed to delete Meraki SD-WAN instance credential %s: %v", id, err)
		}
	}
}

// setMerakiSDWANInstances set the instances from the API response, in the
// order of the state, then the instances the state doesn't have.
// auth_token is not returned by the API, so it is kept from the state of
// the instance with the same host_name. An instance not in the state gets
// an empty token, as after an import. So does an instance whose credential
// changed outside the state (an apply that didn't save its state, or the
// UI): the state token may not be the deployed one, so the next apply
// keeps the credential and doesn't send a spent token again.
func setMerakiSDWANInstances(d *schema.ResourceData, connector *alkira.ConnectorMerakiSDWAN) {
	tokens := make(map[string]interface{})
	credentialIds := make(map[string]string)

	for _, inst := range d.Get("instances").([]interface{}) {
		i := inst.(map[string]interface{})
		hostName := i["host_name"].(string)
		tokens[hostName] = i["auth_token"]
		credentialIds[hostName], _ = i["credential_id"].(string)
	}

	byHostName := make(map[string]alkira.ConnectorMerakiSDWANInstance)
	for _, info := range connector.Instances {
		byHostName[info.HostName] = info
	}

	ordered := make([]alkira.ConnectorMerakiSDWANInstance, 0, len(connector.Instances))
	for _, inst := range d.Get("instances").([]interface{}) {
		hostName, _ := inst.(map[string]interface{})["host_name"].(string)
		if info, found := byHostName[hostName]; found {
			ordered = append(ordered, info)
			delete(byHostName, hostName)
		}
	}
	for _, info := range connector.Instances {
		if _, found := byHostName[info.HostName]; found {
			ordered = append(ordered, info)
		}
	}

	var instances []map[string]interface{}

	for _, info := range ordered {
		authToken, found := tokens[info.HostName]
		if !found {
			authToken = ""
		}
		if id := credentialIds[info.HostName]; id != "" && id != info.CredentialId {
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

// syncMerakiSDWANSentInstances set the computed fields of each instance,
// matched by host_name: credential_id and id as sent in the update, and
// name, az and wan_ip from the prior state. The planned value of a computed
// field stays at its list index, so an instance that moved or got a new
// token has another instance's or its old values. Read would take a wrong
// credential_id as a credential changed outside Terraform and clear the
// token, and a path that saves the state without a read would keep the
// wrong values. A new instance has no id, name, az or wan_ip until a read.
func syncMerakiSDWANSentInstances(d *schema.ResourceData, request *alkira.ConnectorMerakiSDWAN, old []interface{}) {
	oldInstances := make(map[string]map[string]interface{}, len(old))
	for _, inst := range old {
		if o, ok := inst.(map[string]interface{}); ok {
			hostName, _ := o["host_name"].(string)
			oldInstances[hostName] = o
		}
	}

	sent := make(map[string]alkira.ConnectorMerakiSDWANInstance, len(request.Instances))
	for _, inst := range request.Instances {
		sent[inst.HostName] = inst
	}

	instances := d.Get("instances").([]interface{})
	for _, inst := range instances {
		i := inst.(map[string]interface{})
		hostName, _ := i["host_name"].(string)
		s, found := sent[hostName]
		if !found {
			continue
		}

		o := oldInstances[hostName]
		id, _ := s.Id.Int64()
		name, _ := o["name"].(string)
		az, _ := o["az"].(int)
		wanIp, _ := o["wan_ip"].(string)

		i["credential_id"] = s.CredentialId
		i["id"] = int(id)
		i["name"] = name
		i["az"] = az
		i["wan_ip"] = wanIp
	}

	d.Set("instances", instances)
}

// recordMerakiSDWANConfigTokens record the configured token of each instance
// whose planned token is empty, as after an import, and return their host
// names. The token was not sent: the instance kept its credential. Without
// it in the state, a later token change would not be planned.
func recordMerakiSDWANConfigTokens(d *schema.ResourceData) []string {
	tokens := make(map[string]string)
	for _, inst := range merakiSDWANConfigTokens(d) {
		tokens[inst.hostName] = inst.token
	}

	var recorded []string
	instances := d.Get("instances").([]interface{})
	for _, inst := range instances {
		i := inst.(map[string]interface{})
		hostName, _ := i["host_name"].(string)
		if token, _ := i["auth_token"].(string); token == "" && tokens[hostName] != "" {
			i["auth_token"] = tokens[hostName]
			recorded = append(recorded, hostName)
		}
	}

	if len(recorded) > 0 {
		d.Set("instances", instances)
	}

	return recorded
}

// appendMerakiSDWANHostNames append the host names that names doesn't have
func appendMerakiSDWANHostNames(names []string, more []string) []string {
	for _, name := range more {
		if !slices.Contains(names, name) {
			names = append(names, name)
		}
	}
	return names
}

type merakiSDWANConfigToken struct {
	hostName string
	token    string
}

// merakiSDWANConfigTokens is the host_name and auth_token of each instance
// in the raw config, in config order. Unknown and null values are skipped.
func merakiSDWANConfigTokens(d *schema.ResourceData) []merakiSDWANConfigToken {
	config := d.GetRawConfig()
	if config.IsNull() || !config.IsKnown() {
		return nil
	}

	configured := config.GetAttr("instances")
	if configured.IsNull() || !configured.IsKnown() {
		return nil
	}

	var result []merakiSDWANConfigToken
	for it := configured.ElementIterator(); it.Next(); {
		_, inst := it.Element()
		hostName, token := inst.GetAttr("host_name"), inst.GetAttr("auth_token")
		if hostName.IsKnown() && !hostName.IsNull() && token.IsKnown() && !token.IsNull() {
			result = append(result, merakiSDWANConfigToken{hostName.AsString(), token.AsString()})
		}
	}

	return result
}

// checkMerakiSDWANConfigTokens check that no two instances have the same
// configured token. The planned token of an imported instance is empty, so
// expandMerakiSDWANInstances can't see a token it shares.
func checkMerakiSDWANConfigTokens(d *schema.ResourceData) error {
	seen := make(map[string]string)
	for _, inst := range merakiSDWANConfigTokens(d) {
		if other, found := seen[inst.token]; found && inst.token != "" {
			return fmt.Errorf("instances: %q and %q have the same auth_token. "+
				"A token can be used once: generate a token for each instance", other, inst.hostName)
		}
		seen[inst.token] = inst.hostName
	}

	return nil
}
