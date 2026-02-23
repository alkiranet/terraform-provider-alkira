package alkira

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
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

func deflateSegmentOptions(c alkira.SegmentNameToZone) []map[string]interface{} {
	var options []map[string]interface{}

	for _, outerZoneToGroups := range c {
		for zone, groups := range outerZoneToGroups.ZonesToGroups {
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
	log.Printf("[DEBUG] Convert TypeList %v", in)

	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty TypeList to convert to StringList")
		return nil
	}

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

// importWithReadValidation wraps a Read function for import operations.
// During import, any diagnostic (warning or error) is treated as a failure
// to ensure imports fail clearly when the resource cannot be retrieved.
// This addresses the issue where Read returns diag.Warning for failed API calls,
// which Terraform treats as non-fatal during import, causing "Import successful!"
// messages even when the import actually failed.
func importWithReadValidation(readFunc schema.ReadContextFunc) schema.StateContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
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
