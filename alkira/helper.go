package alkira

import (
	"encoding/json"
	"log"

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

func convertStringArrToInterfaceArr(sArr []string) []interface{} {
	iArr := make([]interface{}, len(sArr))
	for i, v := range sArr {
		iArr[i] = v
	}

	return iArr
}
