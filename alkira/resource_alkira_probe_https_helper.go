package alkira

import (
	"encoding/json"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// HTTPS Probe Helper Functions
func generateHTTPSProbeRequest(d *schema.ResourceData) (*alkira.Probe, error) {
	probe := &alkira.Probe{
		Name:             d.Get("name").(string),
		Type:             "HTTPS",
		Enabled:          d.Get("enabled").(bool),
		FailureThreshold: d.Get("failure_threshold").(int),
		SuccessThreshold: d.Get("success_threshold").(int),
		PeriodSeconds:    d.Get("period_seconds").(int),
		TimeoutSeconds:   d.Get("timeout_seconds").(int),
	}

	// Network Entity
	networkEntity := expandNetworkEntityHTTPS(d.Get("network_entity").([]interface{}))
	probe.NetworkEntity = *networkEntity

	// HTTPS Parameters
	httpsParams := &alkira.HttpsProbe{
		URI:                   d.Get("uri").(string),
		ServerName:            d.Get("server_name").(string),
		DisableCertValidation: d.Get("disable_cert_validation").(bool),
		CaCertificate:         d.Get("ca_certificate").(string),
	}

	// if headers, ok := d.GetOk("headers"); ok {
	// 	httpsParams.Headers = convertMapToStringMapHTTPS(headers.(map[string]interface{}))
	// }

	if validators, ok := d.GetOk("validators"); ok {
		httpsParams.Validators = expandValidatorsHTTPS(validators.([]interface{}))
	}

	paramsJson, _ := json.Marshal(httpsParams)
	probe.Parameters = paramsJson

	return probe, nil
}

func setHTTPSProbeState(probe *alkira.Probe, d *schema.ResourceData) error {
	d.Set("name", probe.Name)
	d.Set("enabled", probe.Enabled)
	d.Set("failure_threshold", probe.FailureThreshold)
	d.Set("success_threshold", probe.SuccessThreshold)
	d.Set("period_seconds", probe.PeriodSeconds)
	d.Set("timeout_seconds", probe.TimeoutSeconds)

	// Network Entity
	d.Set("network_entity", flattenNetworkEntityHTTPS(&probe.NetworkEntity))

	// HTTPS Parameters
	var params alkira.HttpsProbe
	if err := json.Unmarshal(probe.Parameters, &params); err != nil {
		return err
	}

	d.Set("uri", params.URI)
	d.Set("server_name", params.ServerName)
	d.Set("disable_cert_validation", params.DisableCertValidation)
	d.Set("ca_certificate", params.CaCertificate)
	// d.Set("headers", params.Headers)
	d.Set("validators", flattenValidatorsHTTPS(params.Validators))

	return nil
}

func expandNetworkEntityHTTPS(input []interface{}) *alkira.ProbeNetworkEntity {
	if len(input) == 0 {
		return nil
	}

	entity := input[0].(map[string]interface{})
	return &alkira.ProbeNetworkEntity{
		Type: entity["type"].(string),
		ID:   entity["id"].(string),
	}
}

func flattenNetworkEntityHTTPS(entity *alkira.ProbeNetworkEntity) []interface{} {
	return []interface{}{map[string]interface{}{
		"type": entity.Type,
		"id":   entity.ID,
	}}
}

func expandValidatorsHTTPS(input []interface{}) []alkira.ProbeValidator {
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

func flattenValidatorsHTTPS(validators []alkira.ProbeValidator) []interface{} {
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

func convertMapToStringMapHTTPS(input map[string]interface{}) map[string]any {
	result := make(map[string]any)
	for k, v := range input {
		result[k] = v.(string)
	}
	return result
}
