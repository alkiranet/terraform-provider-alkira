package alkira

import (
	"encoding/json"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// TCP Probe Helper Functions
func generateTCPProbeRequest(d *schema.ResourceData) (*alkira.Probe, error) {
	probe := &alkira.Probe{
		Name:             d.Get("name").(string),
		Type:             "TCP",
		Enabled:          d.Get("enabled").(bool),
		FailureThreshold: d.Get("failure_threshold").(int),
		SuccessThreshold: d.Get("success_threshold").(int),
		PeriodSeconds:    d.Get("period_seconds").(int),
		TimeoutSeconds:   d.Get("timeout_seconds").(int),
	}

	// Network Entity
	probe.NetworkEntity = alkira.ProbeNetworkEntity{
		Type: "INTERNET_APPLICATION", // Hardcode type
		ID:   d.Get("network_entity_id").(string),
	}

	// TCP Parameters
	tcpParams := &alkira.TcpProbe{
		Port: d.Get("port").(int),
	}
	paramsJson, _ := json.Marshal(tcpParams)
	probe.Parameters = paramsJson

	return probe, nil
}

func setTCPProbeState(probe *alkira.Probe, d *schema.ResourceData) error {
	d.Set("name", probe.Name)
	d.Set("enabled", probe.Enabled)
	d.Set("failure_threshold", probe.FailureThreshold)
	d.Set("success_threshold", probe.SuccessThreshold)
	d.Set("period_seconds", probe.PeriodSeconds)
	d.Set("timeout_seconds", probe.TimeoutSeconds)
	d.Set("network_entity_id", &probe.NetworkEntity.ID)

	// TCP Parameters
	var params alkira.TcpProbe
	if err := json.Unmarshal(probe.Parameters, &params); err != nil {
		return err
	}

	d.Set("port", params.Port)

	return nil
}
