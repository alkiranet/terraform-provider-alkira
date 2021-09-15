package alkira

import (
	"encoding/json"
	"fmt"
)

type InventoryCXP struct {
	Id                string             `json:"id"`
	Name              string             `json:"name"`
	Provider          string             `json:"provider"`
	ProviderRegion    string             `json:"providerRegion"`
	State             string             `json:"state"`
	AvailabilityZones map[string]string  `json:"availabilityZones"`
	Geolocation       map[string]float64 `json:"geolocation"`
}

// Get all cxps from.
func (ac *AlkiraClient) GetCXPs() ([]InventoryCXP, error) {
	uri := fmt.Sprintf("%s/inventory/cxps", ac.URI)
	cxps := []InventoryCXP{}

	data, err := ac.get(uri)
	if err != nil {
		return cxps, fmt.Errorf("GetCXPs: failed to get cxps: %v", err)
	}

	err = json.Unmarshal([]byte(data), &cxps)
	if err != nil {
		return cxps, fmt.Errorf("GetCXPs: failed to unmarshal: %v", err)
	}

	return cxps, err
}

// GetCXPById get the cxp by id.
func (ac *AlkiraClient) GetCXPById(cxpId string) (InventoryCXP, error) {
	uri := fmt.Sprintf("%s/inventory/cxps/%s", ac.URI, cxpId)

	var cxp InventoryCXP
	data, err := ac.get(uri)

	if err != nil {
		return cxp, err
	}

	err = json.Unmarshal([]byte(data), &cxp)
	if err != nil {
		return cxp, fmt.Errorf("GetCXPById: failed to unmarshal: %v", err)
	}

	return cxp, nil
}

// GetCXPByName get the cxp by its name.
func (ac *AlkiraClient) GetCXPByName(name string) (InventoryCXP, error) {
	var cxps []InventoryCXP
	if len(name) == 0 {
		return cxps[0], fmt.Errorf("invalid cxp name input")
	}

	uri := fmt.Sprintf("%s/inventory/cxps?name=%s", ac.URI, name)
	data, err := ac.get(uri)

	if err != nil {
		return cxps[0], err
	}

	err = json.Unmarshal([]byte(data), &cxps)

	if err != nil {
		return cxps[0], fmt.Errorf("GetCXPByName: failed to unmarshal: %v", err)
	}

	if len(cxps) > 1 {
		return cxps[0], fmt.Errorf("multiple cxps found with same name %s: %v", name, cxps)
	}

	return cxps[0], nil
}

// CreateCXP create a new CXP.
func (ac *AlkiraClient) CreateCXP(cxpRequest *InventoryCXP) (string, error) {
	uri := fmt.Sprintf("%s/inventory/cxps", ac.URI)

	body, err := json.Marshal(cxpRequest)
	if err != nil {
		return "", fmt.Errorf("CreateCXP: failed to marshal payload: %v", err)
	}

	data, err := ac.create(uri, body)
	if err != nil {
		return "", err
	}

	var result InventoryCXP
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		return "", fmt.Errorf("CreateCXP: failed to unmarshal: %v", err)
	}

	return result.Id, nil
}

// DeleteCXP delete a cxp by id.
func (ac *AlkiraClient) DeleteCXP(cxpId string) error {
	uri := fmt.Sprintf("%s/inventory/cxps/%s", ac.URI, cxpId)
	return ac.delete(uri)
}

// UpdateCXP update cxp by id.
// currently only state can be updated
func (ac *AlkiraClient) UpdateCXP(cxpId string, state string) error {
	uri := fmt.Sprintf("%s/inventory/cxps/%s", ac.URI, cxpId)
	body, err := json.Marshal(map[string]string{
		"state": state,
	})

	if err != nil {
		return fmt.Errorf("UpdateCxp: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}
