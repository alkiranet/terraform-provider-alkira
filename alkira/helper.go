package alkira

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
)

type createCredential = func(name string, ctype alkira.CredentialType, credential interface{}) (string, error)

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

type getSegmentById = func(id string) (alkira.Segment, error)
type getSegmentByName = func(name string) (alkira.Segment, error)

func convertSegmentIdsToSegmentNames(getSegById getSegmentById, ids []string) ([]string, error) {
	var segmentNames []string
	for _, id := range ids {
		seg, err := getSegById(id)
		if err != nil {
			log.Printf("[DEBUG] failed to segment. %s does not exist: ", id)
			return nil, err
		}

		segmentNames = append(segmentNames, seg.Name)
	}

	return segmentNames, nil
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
