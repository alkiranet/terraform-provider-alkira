package alkira

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// segmentIdPattern matches a segment ID. Leading zeros are excluded: the
// backend accepts GET /segments/0690 and answers with id 690, which Read then
// writes back to state against a config that still says 0690.
var segmentIdPattern = regexp.MustCompile(`^[1-9][0-9]*$`)

const segmentIdValidationMessage = "must be a segment ID rather than a segment name, for example alkira_segment.example.id"

// validateSegmentId rejects a value that is not a segment ID before it reaches
// GET /segments/<value>. The backend answers a segment name there with a 500
// that the client retries five times with escalating backoff, so the request
// costs minutes before surfacing a type-conversion error from Java.
//
// validateReferenceId runs first because it carries the URI-safety guarantee
// its own comment describes. The ParseInt check bounds the magnitude: the
// backend parses the path segment as a Java Long, and a value that overflows
// it draws the same 500 and the same retries as a name.
func validateSegmentId(id string) error {
	if err := validateReferenceId(id); err != nil {
		return err
	}

	if !segmentIdPattern.MatchString(id) {
		return fmt.Errorf("invalid segment id %q; expected a segment's numeric id rather than its name, for example alkira_segment.example.id", id)
	}

	if _, err := strconv.ParseInt(id, 10, 64); err != nil {
		return fmt.Errorf("invalid segment id %q; value is too large to be a segment id", id)
	}

	return nil
}

// getSegmentNameById get a segment name by its ID
func getSegmentNameById(id string, m interface{}) (string, error) {

	if err := validateSegmentId(id); err != nil {
		return "", err
	}

	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, _, err := segmentApi.GetById(id)

	if err != nil {
		return "", err
	}

	return segment.Name, err
}

// getSegmentIdByName get a segment ID by its name
func getSegmentIdByName(name string, m interface{}) (string, error) {

	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, _, err := segmentApi.GetByName(name)

	if err != nil {
		return "", err
	}

	return string(segment.Id), err
}

// convertSegmentIdsToSegmentNames
func convertSegmentIdsToSegmentNames(in *schema.Set, m interface{}) ([]string, error) {

	if in == nil || in.Len() == 0 {
		log.Printf("[DEBUG] empty SegmentIds to convert to SegmentNames")
		return nil, nil
	}

	segmentNames := make([]string, in.Len())

	for i, id := range in.List() {
		segmentName, err := getSegmentNameById(id.(string), m)

		if err != nil {
			log.Printf("[DEBUG] failed to get segment name by ID %s.", id)
			return nil, err
		}

		segmentNames[i] = segmentName
	}

	return segmentNames, nil
}

// convertSegmentNamesToSegmentIds
func convertSegmentNamesToSegmentIds(names []string, m interface{}) ([]string, error) {
	api := alkira.NewSegment(m.(*alkira.AlkiraClient))

	var segmentIds []string
	for _, name := range names {
		seg, _, err := api.GetByName(name)
		if err != nil {
			log.Printf("[DEBUG] failed to get segment. %s does not exist: ", name)
			return nil, err
		}

		segmentIds = append(segmentIds, string(seg.Id))
	}

	return segmentIds, nil
}
