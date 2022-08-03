package alkira

import (
	"encoding/json"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
)

func getInternetApplicationGroup(client *alkira.AlkiraClient) int {
	groups, err := client.GetConnectorGroups()

	if err != nil {
		log.Printf("[ERROR] failed to get groups")
		return 0
	}

	var result []alkira.ConnectorGroup
	json.Unmarshal([]byte(groups), &result)

	for _, group := range result {
		if group.Name == "ALK-INB-INT-GROUP" {
			return group.Id
		}
	}

	return 0
}

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

func convertTypeListToStringList(in []interface{}) []string {
	if in == nil || len(in) == 0 {
		log.Printf("[DEBUG] empty TypeList to convert to StringList")
		return nil
	}

	strList := make([]string, len(in))

	for i, value := range in {
		strList[i] = value.(string)
	}

	return strList
}

func convertSegmentIdsToSegmentNames(ids []string, m interface{}) ([]string, error) {
	client := m.(*alkira.AlkiraClient)

	var segmentNames []string
	for _, id := range ids {
		seg, err := client.GetSegmentById(id)
		if err != nil {
			log.Printf("[DEBUG] failed to get segment. %s does not exist: ", id)
			return nil, err
		}

		segmentNames = append(segmentNames, seg.Name)
	}

	return segmentNames, nil
}

func convertSegmentNamesToSegmentIds(names []string, m interface{}) ([]string, error) {
	client := m.(*alkira.AlkiraClient)

	var segmentIds []string
	for _, name := range names {
		seg, err := client.GetSegmentByName(name)
		if err != nil {
			log.Printf("[DEBUG] failed to get segment. %s does not exist: ", name)
			return nil, err
		}

		segmentIds = append(segmentIds, strconv.Itoa(seg.Id))
	}

	return segmentIds, nil
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

//For infoblox if there is a failed POST for infoblox the backend does not clean up the
//credentials that were created in preparation for creating the infoblox service. This means
//if you make the same attempt to create an infoblox there will likely already be a credential name
//that exists. This throws an error. To avoid that this function will be used to add a random suffix
//of a-zA-z to the end of the credential name. That way each time an attempt and subsequent failure
//occurs when creating the infoblox there will be no clash with existing credentials. This is only
//neccesary because the infoblox credentials are not exposed in the UI. Otherwise the user could
//manage the credentials themselves.
func randomNameSuffix() string {
	possibleChars := []rune("abcdefghijklmnopqrstuvwxyzABXDEFGHIJKLMNOPQRSTUVWXYZ")

	rand.Seed(time.Now().UnixNano())
	min := 0
	max := len(possibleChars)
	var sb strings.Builder
	var lengthNewStr int = 20

	for i := 0; i < lengthNewStr; i++ {
		j := rand.Intn(max-min) + min
		s := string(possibleChars[j])
		sb.WriteString(s)
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
