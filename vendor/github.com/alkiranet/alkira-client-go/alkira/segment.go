package alkira

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Segment struct {
	Id              int         `json:"id"`
	Name            string      `json:"name"`
}


// Get all segments from the given tenant network
func (ac *AlkiraClient) GetSegments() ([]byte, int) {
	segmentEndpoint := ac.URI + "tenantnetworks/" + ac.TenantNetworkId + "/segments"

	request, err := http.NewRequest("GET", segmentEndpoint, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		log.Printf("Error : %s", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	//log.Println(response.StatusCode)
	//log.Println(string(data))

	return data, response.StatusCode
}

// Create a new Segment
func (ac *AlkiraClient) CreateSegment(name string, asn string, ipBlock string) (int, int) {
	var result Segment

	segmentEndpoint := ac.URI + "tenantnetworks/" + ac.TenantNetworkId + "/segments"

	body, err := json.Marshal(map[string]string{
		"name":    name,
		"asn":     asn,
		"ipBlock": ipBlock,
	})

	request, err := http.NewRequest("POST", segmentEndpoint, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		log.Printf("Error : %s", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	json.Unmarshal([]byte(data), &result)

	return result.Id, response.StatusCode
}

// Delete a segment
func (ac *AlkiraClient) DeleteSegment(segmentId string) (int) {
	segmentEndpoint := ac.URI + "tenantnetworks/" + ac.TenantNetworkId + "/segments/" + segmentId
	log.Println(segmentEndpoint)

	request, err := http.NewRequest("DELETE", segmentEndpoint, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		log.Printf("Error : %s", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	log.Println(response.StatusCode)
	log.Println(string(data))

	return response.StatusCode
}
