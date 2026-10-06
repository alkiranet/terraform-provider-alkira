package alkira

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandSegmentOptions(in *schema.Set, m interface{}) (alkira.SegmentNameToZone, error) {
	// As segment options are optional we don't care if none are
	// provided
	if in == nil || in.Len() == 0 {
		return nil, nil
	}

	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segmentOptions := make(alkira.SegmentNameToZone)

	for _, options := range in.List() {
		optionsCfg := options.(map[string]interface{})
		zonesToGroups := make(alkira.ZoneToGroups)
		z := alkira.OuterZoneToGroups{}

		var zoneName *string
		var segment *alkira.Segment
		var groups []string

		if v, ok := optionsCfg["zone_name"].(string); ok {
			zoneName = &v
		}

		if v, ok := optionsCfg["segment_id"].(string); ok {
			if err := validateSegmentId(v); err != nil {
				return nil, err
			}

			seg, _, err := segmentApi.GetById(v)

			if err != nil {
				return nil, err
			}
			segment = seg
		}

		if v, ok := optionsCfg["groups"].([]interface{}); ok && len(v) > 0 {
			groups = convertTypeListToStringList(v)
		} else {
			groups = []string{}
		}

		if zoneName == nil || segment == nil {
			return nil, errors.New("segment_option zone_name and segment_id cannot be nil")
		}

		if v, ok := segmentOptions[segment.Name]; ok {
			v.ZonesToGroups[*zoneName] = groups
		} else {
			zonesToGroups[*zoneName] = groups
			z.ZonesToGroups = zonesToGroups

			segId, _ := strconv.Atoi(string(segment.Id))
			z.SegmentId = segId

			segmentOptions[segment.Name] = z
		}
	}

	return segmentOptions, nil
}

// alkiraManagementZoneName is the backend-injected management zone
// that should be filtered from Terraform state to prevent perpetual
// drift. The backend auto-creates this zone for all firewall services
// but users never configure it in HCL.
const alkiraManagementZoneName = "ALKIRA_MGMT_ZONE"

func deflateSegmentOptions(c alkira.SegmentNameToZone) []map[string]interface{} {
	var options []map[string]interface{}

	for _, outerZoneToGroups := range c {
		for zone, groups := range outerZoneToGroups.ZonesToGroups {
			if zone == alkiraManagementZoneName {
				log.Printf("[DEBUG] Filtering out backend-injected %s from segment_options", alkiraManagementZoneName)
				continue
			}
			i := map[string]interface{}{
				"segment_id": strconv.Itoa(outerZoneToGroups.SegmentId),
				"zone_name":  zone,
				"groups":     groups,
			}
			options = append(options, i)
		}
	}

	return options
}

// convertTypeListToIntList convert a TypeList into a list of int
func convertTypeListToIntList(in []interface{}) []int {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty TypeList to convert to IntList")
		return nil
	}

	intList := make([]int, len(in))

	for i, value := range in {
		intList[i] = value.(int)
	}

	return intList
}

// convertTypeListToStringList convert a TypeList into a list of string
func convertTypeListToStringList(in []interface{}) []string {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty TypeList to convert to StringList")
		return nil
	}

	// Log only the entry count, never the values. This is a generic converter
	// with no way to know whether its input is secret, and IPSec pre-shared keys
	// are among the values routed through it.
	log.Printf("[DEBUG] Convert TypeList with %d entries to StringList", len(in))

	strList := make([]string, len(in))

	for i, value := range in {
		if value != nil {
			strList[i] = value.(string)
		} else {
			strList[i] = ""
		}
	}

	return strList
}

// convertTypeSetToIntList convert a TypeSet into a list of int
func convertTypeSetToIntList(in *schema.Set) []int {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] empty TypeSet to convert to IntList")
		return nil
	}

	intList := make([]int, in.Len())

	for i, value := range in.List() {
		intList[i] = value.(int)
	}

	return intList
}

// convertTypeSetToStringList convert a TypeSet into a list of string
func convertTypeSetToStringList(in *schema.Set) []string {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] empty TypeSet to convert to StringList")
		return nil
	}

	intList := make([]string, in.Len())

	for i, value := range in.List() {
		intList[i] = value.(string)
	}

	return intList
}

func convertStringArrToInterfaceArr(sArr []string) []interface{} {
	iArr := make([]interface{}, len(sArr))
	for i, v := range sArr {
		iArr[i] = v
	}

	return iArr
}

func getAllCredentialsAsCredentialResponseDetails(client *alkira.AlkiraClient) ([]alkira.CredentialResponseDetail, error) {
	credentials, err := client.GetCredentials()
	if err != nil {
		log.Printf("[INFO] Failed getting Credential list")
		return nil, err
	}

	var result []alkira.CredentialResponseDetail
	err = json.Unmarshal([]byte(credentials), &result)
	if err != nil {
		log.Printf("[INFO] Failed Unmarshalling Credentials")
		return nil, err
	}

	return result, nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// For infoblox if there is a failed POST for infoblox the backend does not clean up the
// credentials that were created in preparation for creating the infoblox service. This means
// if you make the same attempt to create an infoblox there will likely already be a credential name
// that exists. This throws an error. To avoid that this function will be used to add a random suffix
// of a-zA-z to the end of the credential name. That way each time an attempt and subsequent failure
// occurs when creating the infoblox there will be no clash with existing credentials. This is only
// necessary because the infoblox credentials are not exposed in the UI. Otherwise the user could
// manage the credentials themselves.
func randomNameSuffix() string {
	possibleChars := []rune("abcdefghijklmnopqrstuvwxyzABXDEFGHIJKLMNOPQRSTUVWXYZ")

	min := 0
	max := len(possibleChars)
	var sb strings.Builder
	var lengthNewStr = 20

	for i := 0; i < lengthNewStr; i++ {
		j := rand.Intn(max-min) + min
		sb.WriteRune(possibleChars[j])
	}

	return sb.String()
}

func convertInputTimeToEpoch(t string) (int64, error) {
	layout := "2006-01-02"
	timeInput, err := time.Parse(layout, t)

	if err != nil {
		log.Printf("[ERROR] Failed parse the time input.")
		return 0, err
	}

	return timeInput.Unix(), nil
}

// importIDPattern allow-lists the characters accepted in a `terraform import`
// id. Every alkira_* resource id is portal-issued and is always a short
// numeric or alphanumeric token (see the resource docs' import examples) --
// this pattern accepts all such ids. It rejects ids containing path
// metacharacters ("/", ".."), query-string delimiters ("?", "&"), or other
// characters that would otherwise be interpolated unescaped into the request
// URI the vendored alkira-client-go SDK builds
// (e.g. `fmt.Sprintf("%s/%s", a.Uri, id)`), which has no escaping of its
// own.
//
// Scope: this check runs once, at `terraform import` time, via
// importWithReadValidation below, and also on the ordinary plan/apply path
// via validateReferenceId (below) for cross-resource id references such as
// `segment_id`. It does not run again for an already-imported resource's
// own id on subsequent Read/Update/Delete -- those read the id straight
// back out of terraform.tfstate via d.Id() with no revalidation. An id
// written directly into state (hand-edited, restored from a tampered
// backend, or produced by a state-manipulation tool) reaches the same
// unescaped SDK call unguarded. Closing that gap means validating at the
// point of use (Read/Update/Delete) rather than only at the point of entry
// (import); that is a larger, cross-cutting change spanning every resource
// file and is tracked as a follow-up, not done in this change.
//
// The upper bound is a sanity limit, not a contract: it is not derived from
// any API guarantee about maximum id length, so it is set well above every
// id shape observed today (numeric ids, short alphanumeric tokens, UUIDs)
// rather than tight to it, to avoid rejecting a legitimate id if the portal
// ever issues a longer one.
//
// For the six alkira_credential_* resources (aws_vpc, azure_vnet, gcp_vpc,
// oci_vcn, prisma_sdwan, ssh_key_pair), Read is a no-op that always returns
// nil, so this check is the only validation an import id gets for those
// resources -- there is no second line of defense from a 404 on a bad id
// the way there is for every other resource's Read/GetById call.
var importIDPattern = regexp.MustCompile(`^[A-Za-z0-9_-]{1,256}$`)

// validateReferenceId checks a cross-resource id reference (e.g. a
// segment_id given in HCL config for another resource) against the same
// allow-list as importIDPattern, before it is used to look up that
// resource server-side.
//
// This guards the ordinary plan/apply path, not just import: a field such
// as `segment_id` is an ordinary schema.TypeString with no `Validate*`,
// and its value is passed straight to a GetById-style lookup whose URI is
// built with unescaped string interpolation (see importIDPattern above).
// Without this check, a crafted value (e.g. containing "../" and a "#"
// fragment to neutralize an appended query string) reaches that lookup on
// a plain `terraform apply` -- no import required. Every resource or
// helper that takes a portal-issued id from config and passes it to such
// a lookup should call this first.
//
// Named distinctly from the test-only validateResourceId (test_utils.go),
// which checks a *returned* resource id is purely numeric -- a different,
// stricter check for a different purpose (asserting on API responses in
// tests, not validating attacker-reachable config input).
func validateReferenceId(id string) error {
	if !importIDPattern.MatchString(id) {
		return fmt.Errorf("invalid resource id %q; expected an alphanumeric id (letters, digits, '-', '_'), 1-256 characters", id)
	}
	return nil
}

// importWithReadValidation wraps a Read function for import operations.
// During import, any diagnostic (warning or error) is treated as a failure
// to ensure imports fail clearly when the resource cannot be retrieved.
// This addresses the issue where Read returns diag.Warning for failed API calls,
// which Terraform treats as non-fatal during import, causing "Import successful!"
// messages even when the import actually failed.
func importWithReadValidation(readFunc schema.ReadContextFunc) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		// Reject ids that don't look like a portal-issued resource id before
		// they ever reach the Read function (and, downstream, the
		// unescaped URI interpolation in the vendored SDK).
		if !importIDPattern.MatchString(d.Id()) {
			return nil, fmt.Errorf("import failed: invalid resource id %q; expected an alphanumeric id (letters, digits, '-', '_'), 1-256 characters", d.Id())
		}

		// Call the Read function to populate state
		diags := readFunc(ctx, d, m)

		// During import, treat any diagnostic (warning or error) as failure
		if diags.HasError() || len(diags) > 0 {
			var msgs []string
			for _, diagnostic := range diags {
				msgs = append(msgs, diagnostic.Summary+": "+diagnostic.Detail)
			}
			return nil, errors.New("import failed: " + strings.Join(msgs, "; "))
		}

		return []*schema.ResourceData{d}, nil
	}
}

// toInt converts a value to int, handling both int and string representations
// that may appear in raw state maps.
func toInt(v interface{}) int {
	switch val := v.(type) {
	case int:
		return val
	case float64:
		return int(val)
	case string:
		if i, err := strconv.Atoi(strings.TrimSpace(val)); err == nil {
			return i
		}
	}
	return 0
}

// typeSetHash returns a schema.SchemaSetFunc that computes a hash
// for TypeSet elements using a key extractor function. The key extractor
// builds a string from the element's fields, which is then hashed to
// produce the set key. This allows elements to be matched by content
// rather than position, preventing spurious updates when elements are
// added, removed, or reordered.
//
// IMPORTANT: TypeSet compares elements solely by hash. The key extractor
// MUST include ALL fields of the block. If any field is omitted, changes
// to that field will be invisible to Terraform — plan will show
// "No changes" even when the value has changed. The corresponding Read
// helper that populates the set from API data must also always set every
// field (including empty/zero values) so that hashes match the config.
//
// For single-field blocks: return just that field (e.g., m["hostname"])
// For multi-field blocks: combine all fields with fmt.Sprintf
func typeSetHash(keyExtractor func(map[string]interface{}) string) schema.SchemaSetFunc {
	return func(v interface{}) int {
		var buf bytes.Buffer
		m := v.(map[string]interface{})
		fmt.Fprintf(&buf, "%s-", keyExtractor(m))
		return schema.HashString(buf.String())
	}
}

// warnOnFailedStateUpdate wraps a resource's UpdateContextFunc and emits a
// non-fatal warning when an update carrying config changes is applied
// against a resource in FAILED provision state. In that case, the backend
// skips the config update and re-provisions the previously saved config
// (the retry-by-reapply mechanism), so the requested changes are not saved.
func warnOnFailedStateUpdate(update schema.UpdateContextFunc) schema.UpdateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
		client := m.(*alkira.AlkiraClient)

		// Compute the condition up front for clarity. HasChangesExcept
		// filters out the retry-only diff forced by CustomizeDiff.
		oldState, _ := d.GetChange("provision_state")

		warn := client.Provision &&
			oldState == "FAILED" &&
			d.HasChangesExcept("provision_state")

		diags := update(ctx, d, m)

		// Suppress the warning if the update itself failed - the
		// skip message would be misleading in that case.
		if warn && !diags.HasError() {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "CONFIGURATION CHANGES SKIPPED",
				Detail: "Resource was in FAILED provision state; the backend " +
					"re-provisions the previously saved configuration and skips " +
					"configuration changes until the resource recovers. See the " +
					"PROVISIONING section in the provider documentation.",
			})
		}

		return diags
	}
}

// validateIPv4CidrOrIP validates that a string is a valid IPv4 address
// or IPv4 CIDR. IPv6 is rejected.
func validateIPv4CidrOrIP(i interface{}, k string) ([]string, []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	// A colon means an IPv6 literal.
	if strings.Contains(v, ":") {
		return nil, []error{fmt.Errorf("%q must be an IPv4 address or IPv4 CIDR, got %q (IPv6 is not supported)", k, v)}
	}

	if strings.Contains(v, "/") {
		ip, ipnet, err := net.ParseCIDR(v)
		if err != nil || ip.To4() == nil {
			return nil, []error{fmt.Errorf("%q must be a valid IPv4 CIDR, got %q", k, v)}
		}
		// Reject host bits rather than silently masking them.
		if !ip.Equal(ipnet.IP) {
			return nil, []error{fmt.Errorf("%q must be a CIDR network address with no host bits set, got %q (did you mean %q?)", k, v, ipnet.String())}
		}
		return nil, nil
	}

	if ip := net.ParseIP(v); ip == nil || ip.To4() == nil {
		return nil, []error{fmt.Errorf("%q must be a valid IPv4 address or IPv4 CIDR, got %q", k, v)}
	}

	return nil, nil
}
