package alkira

import (
	"encoding/json"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func generateHTTPProbeRequest(d *schema.ResourceData) (*alkira.Probe, error) {
	probe := &alkira.Probe{
		Name:             d.Get("name").(string),
		Type:             "HTTP",
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

	// HTTP Parameters
	httpParams := &alkira.HttpProbe{
		URI: d.Get("uri").(string),
	}

	if validators, ok := d.GetOk("validators"); ok {
		httpParams.Validators = expandValidatorsHTTP(validators.([]interface{}))
	}

	paramsJson, _ := json.Marshal(httpParams)
	probe.Parameters = paramsJson

	return probe, nil
}

func setHTTPProbeState(probe *alkira.Probe, d *schema.ResourceData) error {
	d.Set("name", probe.Name)
	d.Set("enabled", probe.Enabled)
	d.Set("failure_threshold", probe.FailureThreshold)
	d.Set("success_threshold", probe.SuccessThreshold)
	d.Set("period_seconds", probe.PeriodSeconds)
	d.Set("timeout_seconds", probe.TimeoutSeconds)

	// Network Entity
	d.Set("network_entity_id", probe.NetworkEntity.ID)

	// HTTP Parameters
	var params alkira.HttpProbe
	if err := json.Unmarshal(probe.Parameters, &params); err != nil {
		return err
	}

	d.Set("uri", params.URI)
	// d.Set("headers", params.Headers)
	d.Set("validators", flattenValidatorsHTTP(params.Validators))

	return nil
}

func expandValidatorsHTTP(input []interface{}) []alkira.ProbeValidator {
	validators := make([]alkira.ProbeValidator, 0)

	for _, v := range input {
		val := v.(map[string]interface{})
		validator := alkira.ProbeValidator{
			Type: val["type"].(string),
		}

		var params interface{}
		switch validator.Type {
		case "HTTP_STATUS_CODE":
			params = alkira.ProbeStatusCodeValidator{
				StatusCode: val["status_code"].(string),
			}
		case "HTTP_RESPONSE_BODY":
			params = alkira.ProbeResponseBodyValidator{
				Regex: val["regex"].(string),
			}
		}

		paramsJson, _ := json.Marshal(params)
		validator.Parameters = paramsJson
		validators = append(validators, validator)
	}

	return validators
}

func flattenValidatorsHTTP(validators []alkira.ProbeValidator) []interface{} {
	result := make([]interface{}, 0)

	for _, v := range validators {
		val := map[string]interface{}{
			"type": v.Type,
		}

		switch v.Type {
		case "HTTP_STATUS_CODE":
			var sc alkira.ProbeStatusCodeValidator
			json.Unmarshal(v.Parameters, &sc)
			val["status_code"] = sc.StatusCode
		case "HTTP_RESPONSE_BODY":
			var rb alkira.ProbeResponseBodyValidator
			json.Unmarshal(v.Parameters, &rb)
			val["regex"] = rb.Regex
		}

		result = append(result, val)
	}

	return result
}

func convertMapToStringMapHTTP(input map[string]interface{}) map[string]any {
	result := make(map[string]any)
	for k, v := range input {
		result[k] = v.(string)
	}
	return result
}
