package alkira

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/require"
)

func TestGenerateCheckpointRequestService(t *testing.T) {
	seedInstance := alkira.CheckpointInstance{
		CredentialId: "testInstanceCredentialId",
		Name:         "testInstanceName",
	}
	expectedInstances := makeNumCheckpointInstances(3, seedInstance)

	expectedManagementServer := &alkira.CheckpointManagementServer{
		ConfigurationMode: "MANUAL",
		CredentialId:      "testCredentialIdMg",
		Domain:            "domain",
		GlobalCidrListId:  0,
		Ips:               []string{"1", "2", "3"},
		Reachability:      "reachability",
		Segment:           "segmentName",
		SegmentId:         0,
		Type:              "type",
		UserName:          "userName",
	}

	expectedZoneName := "zoneName"
	expectedGroups := []string{"1", "2", "3", "4"}
	expectedSegment := &alkira.Segment{
		Id:      0,
		Name:    "0",
		IpBlock: "10.255.254.0/24",
		Asn:     65001,
	}

	zonesToGroups := make(alkira.CheckpointZoneToGroups)
	zonesToGroups[expectedZoneName] = expectedGroups

	z := alkira.OuterZoneToGroups{
		SegmentId:     expectedSegment.Id,
		ZonesToGroups: zonesToGroups,
	}

	expectedSegmentOptions := make(map[string]alkira.OuterZoneToGroups)
	expectedSegmentOptions[expectedSegment.Name] = z

	expectedCheckpoint := &alkira.Checkpoint{
		AutoScale:        "ON",
		BillingTags:      []int{1, 2, 3},
		CredentialId:     "testCredentialId",
		Cxp:              "US-WEST",
		Description:      "testDescription",
		Instances:        expectedInstances,
		LicenseType:      "BRING_YOUR_OWN",
		ManagementServer: expectedManagementServer,
		MinInstanceCount: 1,
		MaxInstanceCount: 1,
		Name:             "name",
		PdpIps:           []string{"10.10.10.10"},
		Segments:         []string{"segmentName0", "segmentName1"},
		SegmentOptions:   expectedSegmentOptions,
		Size:             "size",
		TunnelProtocol:   "tunnel",
		Version:          "version",
	}

	ac := serveCheckpoint(t, expectedCheckpoint)

	r := resourceAlkiraCheckpoint()
	d := r.TestResourceData()

	s := newSetFromCheckpointResource(convertCheckpointInstanceToArrayInterface(expectedInstances))
	d.Set("instances", s)

	err := resourceCheckpointRead(d, ac)
	require.Nil(t, err)

	actualInstances := expandCheckpointInstances(s)
	require.ElementsMatch(t, actualInstances, expectedInstances)

	actualSegmentOptions, err := expandCheckpointSegmentOptions(d.Get("segment_options").(*schema.Set), getCheckpointSegmentInTest)
	require.Nil(t, err)
	require.Contains(t, actualSegmentOptions, "segmentName")
}

//
// TEST HELPER
//
func serveCheckpoint(t *testing.T, c *alkira.Checkpoint) *alkira.AlkiraClient {
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			json.NewEncoder(w).Encode(c)
			w.Header().Set("Content-Type", "application/json")
		},
	))
	t.Cleanup(server.Close)

	return &alkira.AlkiraClient{
		URI:             server.URL,
		TenantNetworkId: "0",
		Client:          &http.Client{Timeout: time.Duration(1) * time.Second},
	}
}
