package alkira

import (
	"regexp"
	"sort"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// secretAttrPattern matches attribute names that look like they hold secret
// material. It is intentionally broad and word-boundary aware (via `_` or
// start/end of string as the separator, since schema keys are snake_case) so
// it catches new secret-shaped attributes as they are added, rather than
// pinning today's known list.
var secretAttrPattern = regexp.MustCompile(
	`(?i)(^|_)(password|secret|private_key|preshared|shared_secret|auth_key|license_key|api_key|token|passphrase)($|_)`,
)

// allowedNonSensitive is the reviewed exception list: attributes that match
// secretAttrPattern by name but are not secret material, so they are
// deliberately left unmarked. Every entry must carry a one-line justification.
var allowedNonSensitive = map[string]string{
	// A filesystem path to a license key file, not the key material itself.
	"license_key_file_path": "path to a license key file, not key material",
	// A GCP OAuth2 token endpoint URL (e.g. https://oauth2.googleapis.com/token),
	// not a token value.
	"token_uri": "GCP OAuth2 endpoint URL, not a token value",
	// A username, matched only because it shares the "pan_" prefix convention
	// with pan_password/pan_license_key; contains no secret material.
	"pan_username": "PAN Panorama username, not secret material",
	// A GCP service-account key identifier (fingerprint used for key
	// rotation/reference). It does not grant access on its own — the actual
	// secret is the sibling "private_key" attribute, which is Sensitive.
	"private_key_id": "GCP service-account key identifier, not key material (see PR body)",
}

// walkSchema recurses into a resource schema, including nested TypeList/TypeSet
// blocks whose Elem is a *schema.Resource, and reports every attribute whose
// name matches secretAttrPattern. path identifies the attribute's location for
// failure messages, e.g. "alkira_connector_ipsec/endpoint/bgp_auth_key".
func walkSchema(resourceName, path string, s map[string]*schema.Schema, out *[]string, sensitiveOK *[]bool) {
	for key, attr := range s {
		attrPath := path + "/" + key
		if secretAttrPattern.MatchString(key) {
			*out = append(*out, resourceName+attrPath)
			*sensitiveOK = append(*sensitiveOK, attr.Sensitive)
		}
		if attr.Elem != nil {
			if r, ok := attr.Elem.(*schema.Resource); ok {
				walkSchema(resourceName, attrPath, r.Schema, out, sensitiveOK)
			}
		}
	}
}

// TestSecretResourceAttributesAreSensitive walks every resource in
// Provider().ResourcesMap (including nested TypeList/TypeSet sub-resources)
// and fails if any attribute whose name looks secret-shaped lacks
// Sensitive: true, unless it is named in allowedNonSensitive. This is
// deny-by-default: it does not enumerate the attributes it expects to find,
// so it also catches secret attributes added after this test was written.
func TestSecretResourceAttributesAreSensitive(t *testing.T) {
	provider := Provider()

	var flagged []string
	var sensitiveOK []bool

	resourceNames := make([]string, 0, len(provider.ResourcesMap))
	for name := range provider.ResourcesMap {
		resourceNames = append(resourceNames, name)
	}
	sort.Strings(resourceNames)

	for _, name := range resourceNames {
		res := provider.ResourcesMap[name]
		if res == nil {
			continue
		}
		walkSchema(name, "", res.Schema, &flagged, &sensitiveOK)
	}

	var failures []string
	for i, attrPath := range flagged {
		// attrPath looks like "<resource>/<nested>.../<attr>"; the leaf
		// segment after the last "/" is the actual schema key.
		leaf := attrPath
		if idx := lastSlash(attrPath); idx >= 0 {
			leaf = attrPath[idx+1:]
		}
		if _, ok := allowedNonSensitive[leaf]; ok {
			continue
		}
		if !sensitiveOK[i] {
			failures = append(failures, attrPath)
		}
	}

	if len(failures) > 0 {
		sort.Strings(failures)
		t.Fatalf("the following attributes look secret-shaped (matched against %q) "+
			"but are not marked Sensitive, so their values render in cleartext in "+
			"terraform plan/apply output and CI job logs. Either add Sensitive: true, "+
			"or add the attribute name to allowedNonSensitive with a one-line "+
			"justification if it is genuinely not secret material:\n  %s",
			secretAttrPattern.String(), joinLines(failures))
	}
}

func lastSlash(s string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == '/' {
			return i
		}
	}
	return -1
}

func joinLines(lines []string) string {
	out := ""
	for i, l := range lines {
		if i > 0 {
			out += "\n  "
		}
		out += l
	}
	return out
}
