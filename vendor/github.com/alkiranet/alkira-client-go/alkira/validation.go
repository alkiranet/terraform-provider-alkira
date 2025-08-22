// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/wait"
)

// ValidationState represents the state of a validation operation
type ValidationState struct {
	State        string       `json:"state"`
	ErrorDetails ErrorDetails `json:"errorDetails,omitempty"`
}

type ErrorDetails struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// GetValidationState retrieves the validation state by its ID
func (ac *AlkiraClient) GetValidationState(id string) (*ValidationState, error) {
	uri := fmt.Sprintf("%s/tenantnetworks/%s/validation-states/%s", ac.URI, ac.TenantNetworkId, id)
	logf("DEBUG", "GetValidationState: requesting URI: %s", uri)
	data, _, err := ac.get(uri)
	if err != nil {
		return nil, err
	}

	// Log detailed information about the response data
	logf("DEBUG", "GetValidationState: raw response data length: %d", len(data))
	logf("DEBUG", "GetValidationState: raw response data as string: %s", string(data))
	logf("DEBUG", "GetValidationState: raw response data as bytes: %v", data)

	var result ValidationState
	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		// Log additional debugging information
		logf("ERROR", "GetValidationState: JSON unmarshal failed with data length: %d", len(data))
		logf("ERROR", "GetValidationState: JSON unmarshal failed with data as string: %s", string(data))
		logf("ERROR", "GetValidationState: JSON unmarshal error: %v", err)
		return nil, fmt.Errorf("GetValidationState: failed to unmarshal: %v", err)
	}

	logf("DEBUG", "GetValidationState: successfully unmarshaled validation state: %+v", result)
	return &result, nil
}

// handleValidation processes validation for a response
func (ac *AlkiraClient) handleValidation(response *http.Response) error {
	logf("DEBUG", "handleValidation: called with response status: %d", response.StatusCode)

	if !ac.Validate {
		logf("DEBUG", "handleValidation: validation is disabled")
		return nil
	}

	validationStateID := response.Header.Get("x-ak-validation-state-id")
	logf("DEBUG", "handleValidation: validation state ID from header: %s", validationStateID)

	if validationStateID != "" {
		logf("DEBUG", "handleValidation: starting validation polling for ID: %s", validationStateID)
		err := wait.Poll(10*time.Second, defaultValTimeout, func() (bool, error) {
			logf("DEBUG", "handleValidation: polling validation state for ID: %s", validationStateID)
			validationState, err := ac.GetValidationState(validationStateID)
			if err != nil {
				logf("ERROR", "handleValidation: failed to get validation state for ID %s: %v", validationStateID, err)
				return false, err
			}

			logf("DEBUG", "handleValidation: received validation state: %+v", validationState)

			switch validationState.State {
			case "SUCCESS":
				logf("DEBUG", "handleValidation: validation succeeded for ID: %s", validationStateID)
				return true, nil
			case "FAILED":
				logf("DEBUG", "handleValidation: validation failed for ID: %s", validationStateID)
				if len(validationState.ErrorDetails.Message) > 0 && validationState.ErrorDetails.Message != "" {
					return false, fmt.Errorf("validation failed: %s", validationState.ErrorDetails.Message)
				} else if validationState.ErrorDetails.Code != "" {
					return false, fmt.Errorf("validation failed: %s", validationState.ErrorDetails.Code)
				}
				return false, fmt.Errorf("validation failed with state: %s", validationState.State)
			}

			logf("DEBUG", "handleValidation: waiting for validation %s to finish. (state: %s)", validationStateID, validationState.State)
			return false, nil
		})

		if err != nil {
			if err == wait.ErrWaitTimeout {
				logf("ERROR", "handleValidation: validation %s timed out", validationStateID)
				return fmt.Errorf("validation %s timed out", validationStateID)
			}
			logf("ERROR", "handleValidation: validation %s failed with error: %v", validationStateID, err)
			return err
		}
		logf("DEBUG", "handleValidation: completed validation for ID: %s", validationStateID)
	} else {
		logf("DEBUG", "handleValidation: no validation state ID found in response headers")
	}

	return nil
}
