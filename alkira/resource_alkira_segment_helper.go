package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getSegmentNamebyId get a segment name by its ID
func getSegmentNameById(id string, m interface{}) (string, error) {

	segmentApi := alkira.NewSegment(m.(*alkira.AlkiraClient))
	segment, err := segmentApi.GetById(id)

	if err != nil {
		return "", err
	}

	return segment.Name, err
}

// getSegmentIdbyName get a segment ID by its name
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
