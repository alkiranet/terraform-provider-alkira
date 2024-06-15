package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
)

// getDirectInterConnectorGroupNameByID gets the direct interconnector group name by its ID.

func getDirectInterConnectorGroupNameByID(id string, m interface{}) (string, error) {
	directInterConnectorGroupApi := alkira.NewInterConnectorCommunicationGroup(m.(*alkira.AlkiraClient))

	directInterConnectorGroup, _, err := directInterConnectorGroupApi.GetById(id)

	if err != nil {
		return "", err
	}

	return directInterConnectorGroup.Name, nil
}

// getDirectInterConnectorGroupIDByName gets a direct interconnector group ID by its name
func getDirectInterConnectorGroupIDByName(name string, m interface{}) (string, error) {
	directInterConnectorGroupApi := alkira.NewInterConnectorCommunicationGroup(m.(*alkira.AlkiraClient))

	directInterConnectorGroup, _, err := directInterConnectorGroupApi.GetByName(name)

	if err != nil {
		return "", err
	}

	return string(directInterConnectorGroup.Id), nil

}
