// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Default provision timeout is 240m
const defaultProvTimeout time.Duration = 240 * time.Minute

// Default validation timeout is 10m
const defaultValTimeout time.Duration = 10 * time.Minute

// Default Retry
const defaultRetryInterval time.Duration = 5 * time.Second
const defaultRetryTimeout time.Duration = 10 * time.Second

type AlkiraClient struct {
	Client                *retryablehttp.Client
	URI                   string
	Username              string
	Password              string
	Secret                string
	Authorization         string
	Provision             bool
	Validate              bool
	TenantNetworkId       string
	SerializationEnabled  bool
	serializationTimeout  time.Duration
	apiMutex              sync.Mutex
}

type Session struct {
	jar map[string][]*http.Cookie
}

func (s *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	s.jar[u.Host] = cookies
}

func (s *Session) Cookies(u *url.URL) []*http.Cookie {
	return s.jar[u.Host]
}

// NewAlkiraClient creates a new API client
func NewAlkiraClient(hostname string, username string, password string, secret string, provision bool, validate bool, serializationEnabled bool, serializationTimeout int, auth string) (*AlkiraClient, error) {

	// Construct the portal URI
	url := "https://" + hostname

	logf("DEBUG", "ALKIRA-PROVISION: %v", provision)

	if auth == "header" {
		logf("DEBUG", "ALKIRA-AUTH-METHOD: %v", auth)
		return NewAlkiraClientWithAuthHeader(url, username, password, secret, provision, validate, serializationEnabled, serializationTimeout)
	}

	return NewAlkiraClientInternal(url, username, password, secret, provision, validate, serializationEnabled, serializationTimeout)
}

// NewAlkiraClientWithAuthHeader creates a new internal Alkira client with authentication in header
func NewAlkiraClientWithAuthHeader(url string, username string, password string, secret string, provision bool, validate bool, serializationEnabled bool, serializationTimeout int) (*AlkiraClient, error) {

	// Firstly, construct the portal API based URI
	apiUrl := url + "/api"

	// Parse serialization configuration
	// Use parameters if provided, otherwise fall back to environment variables
	enableSerialization := serializationEnabled
	timeoutSeconds := serializationTimeout

	if !serializationEnabled && os.Getenv("ALKIRA_API_SERIALIZATION_ENABLED") == "true" {
		enableSerialization = true
	}

	if serializationTimeout == 0 {
		if envTimeout := os.Getenv("ALKIRA_API_SERIALIZATION_TIMEOUT"); envTimeout != "" {
			if parsed, err := strconv.Atoi(envTimeout); err == nil && parsed > 0 {
				timeoutSeconds = parsed
			} else {
				timeoutSeconds = 120 // default
			}
		} else {
			timeoutSeconds = 120 // default
		}
	}

	// Generate Authorization header string
	auth := ""

	if len(secret) > 0 {
		authStr := secret
		auth = "api-key " + base64.StdEncoding.EncodeToString([]byte(authStr))
	} else {
		authStr := username + ":" + password
		auth = "basic " + base64.StdEncoding.EncodeToString([]byte(authStr))
	}

	if len(auth) == 0 {
		return nil, fmt.Errorf("invalid credentials to authenticate")
	}

	// Create retry-able HTTP client
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{},
	}

	// Config retry client
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = tr
	retryClient.RetryMax = 5
	retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		// Respect server-specified Retry-After header
		if resp != nil && resp.StatusCode == 429 {
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					logf("DEBUG", "Retry-After: %v\n", seconds)
					return time.Duration(seconds) * time.Second
				}
			}
		}
		return retryablehttp.LinearJitterBackoff(min, max, attemptNum, resp)
	}
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		shouldRetry, e := retryablehttp.DefaultRetryPolicy(ctx, resp, err)

		// Always retry on 429 and 500 status codes
		if resp != nil {
			switch resp.StatusCode {
			case 429, 500:
				return true, fmt.Errorf("retryable status code: %d", resp.StatusCode)
			}
		}
		return shouldRetry, e
	}

	// Get the tenant network ID
	var result []TenantNetworkId
	tenantNetworkUrl := apiUrl + "/tenantnetworksummaries"

	tenantNetworkRequest, _ := retryablehttp.NewRequest("GET", tenantNetworkUrl, nil)
	tenantNetworkRequest.Header.Set("Content-Type", "application/json")
	tenantNetworkRequest.Header.Set("Authorization", auth)
	tenantNetworkResponse, err := retryClient.Do(tenantNetworkRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to make tenant network request, %v", err)
	}

	defer tenantNetworkResponse.Body.Close()

	data, err := io.ReadAll(tenantNetworkResponse.Body)
	if err != nil {
		logf("ERROR", "NewAlkiraClientWithAuthHeader: failed to read tenant network response body: %v", err)
		return nil, fmt.Errorf("NewAlkiraClientWithAuthHeader: failed to read tenant network response body: %v", err)
	}
	logf("TRACE", "Tenant Network Summary: %s\n", string(data))

	if tenantNetworkResponse.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get tenant network (%d)", tenantNetworkResponse.StatusCode)
	}

	json.Unmarshal([]byte(data), &result)

	tenantNetworkId := 0

	if len(result) > 0 {
		tenantNetworkId = result[0].Id
	} else {
		return nil, fmt.Errorf("failed to get tenant network ID")
	}

	// Construct our client with all information
	client := &AlkiraClient{
		Client:               retryClient,
		URI:                  apiUrl,
		Username:             username,
		Password:             password,
		Secret:               secret,
		Authorization:        auth,
		Provision:            provision,
		Validate:             validate,
		TenantNetworkId:      strconv.Itoa(tenantNetworkId),
		SerializationEnabled: enableSerialization,
		serializationTimeout: time.Duration(timeoutSeconds) * time.Second,
	}

	logf("DEBUG", "ALKIRA-API-SERIALIZATION-ENABLED: %v", client.SerializationEnabled)
	logf("DEBUG", "ALKIRA-API-SERIALIZATION-TIMEOUT: %v", client.serializationTimeout)

	return client, nil
}

// NewAlkiraClientInternal creates a new internal Alkira client
func NewAlkiraClientInternal(url string, username string, password string, secret string, provision bool, validate bool, serializationEnabled bool, serializationTimeout int) (*AlkiraClient, error) {

	// Construct the portal URI based on the given endpoint
	apiUrl := url + "/api"

	// Parse serialization configuration
	// Use parameters if provided, otherwise fall back to environment variables
	enableSerialization := serializationEnabled
	timeoutSeconds := serializationTimeout

	if !serializationEnabled && os.Getenv("ALKIRA_API_SERIALIZATION_ENABLED") == "true" {
		enableSerialization = true
	}

	if serializationTimeout == 0 {
		if envTimeout := os.Getenv("ALKIRA_API_SERIALIZATION_TIMEOUT"); envTimeout != "" {
			if parsed, err := strconv.Atoi(envTimeout); err == nil && parsed > 0 {
				timeoutSeconds = parsed
			} else {
				timeoutSeconds = 120 // default
			}
		} else {
			timeoutSeconds = 120 // default
		}
	}

	loginRequestBody, err := json.Marshal(map[string]string{
		"userName": username,
		"password": password,
		"secret":   secret,
	})

	// Create retry-able HTTP client
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{},
	}

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	retryClient.HTTPClient.Transport = tr
	retryClient.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		// Respect server-specified Retry-After header
		if resp != nil && resp.StatusCode == 429 {
			if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
				if seconds, err := strconv.Atoi(retryAfter); err == nil {
					logf("DEBUG", "Retry-After: %v\n", seconds)
					return time.Duration(seconds) * time.Second
				}
			}
		} else {
			return time.Duration(defaultRetryTimeout)
		}
		return retryablehttp.LinearJitterBackoff(min, max, attemptNum, resp)
	}

	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		shouldRetry, e := retryablehttp.DefaultRetryPolicy(ctx, resp, err)

		// Retry on 409 and 429 status codes
		if resp != nil {
			switch resp.StatusCode {
			case 409, 429:
				return true, fmt.Errorf("retryable status code: %d", resp.StatusCode)
			}
		}
		return shouldRetry, e
	}

	// Login to the portal
	jar := &Session{}
	jar.jar = make(map[string][]*http.Cookie)
	retryClient.HTTPClient.Jar = jar

	// User login
	loginUrl := fmt.Sprintf("%s/user/login", apiUrl)

	request, requestErr := retryablehttp.NewRequest("POST", loginUrl, bytes.NewBuffer(loginRequestBody))

	if requestErr != nil {
		return nil, fmt.Errorf("failed to create login request, %v", requestErr)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := retryClient.Do(request)

	if err != nil {
		return nil, fmt.Errorf("failed to make login request, %v", err)
	}

	defer response.Body.Close()

	userAuthData, err := io.ReadAll(response.Body)
	if err != nil {
		logf("ERROR", "NewAlkiraClientInternal: failed to read user auth response body: %v", err)
		return nil, fmt.Errorf("NewAlkiraClientInternal: failed to read user auth response body: %v", err)
	}

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("failed to login to portal (%d)", response.StatusCode)
	}

	// Obtain the session
	sessionUrl := apiUrl + "/sessions"

	sessionRequest, _ := retryablehttp.NewRequest("POST", sessionUrl, bytes.NewBuffer(userAuthData))
	sessionRequest.Header.Set("Content-Type", "application/json")
	sessionResponse, err := retryClient.Do(sessionRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to make session request, %v", err)
	}

	defer sessionResponse.Body.Close()

	sessionData, err := io.ReadAll(sessionResponse.Body)
	if err != nil {
		logf("ERROR", "NewAlkiraClientInternal: failed to read session response body: %v", err)
		return nil, fmt.Errorf("NewAlkiraClientInternal: failed to read session response body: %v", err)
	}
	logf("TRACE", "session data: %s\n", string(sessionData))

	if sessionResponse.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get session (%d)", sessionResponse.StatusCode)
	}

	// Get the tenant network ID
	var result []TenantNetworkId
	tenantNetworkUrl := apiUrl + "/tenantnetworksummaries"

	tenantNetworkRequest, _ := retryablehttp.NewRequest("GET", tenantNetworkUrl, nil)
	tenantNetworkRequest.Header.Set("Content-Type", "application/json")
	tenantNetworkResponse, err := retryClient.Do(tenantNetworkRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to make tenant network request, %v", err)
	}

	defer tenantNetworkResponse.Body.Close()

	data, err := io.ReadAll(tenantNetworkResponse.Body)
	if err != nil {
		logf("ERROR", "NewAlkiraClientInternal: failed to read tenant network response body: %v", err)
		return nil, fmt.Errorf("NewAlkiraClientInternal: failed to read tenant network response body: %v", err)
	}
	logf("TRACE", "Tenant Network Summary: %s\n", string(data))

	if tenantNetworkResponse.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get tenant network (%d)", tenantNetworkResponse.StatusCode)
	}

	err = json.Unmarshal([]byte(data), &result)
	if err != nil {
		logf("ERROR", "NewAlkiraClientInternal: failed to unmarshal tenant network data: %v", err)
		return nil, fmt.Errorf("NewAlkiraClientInternal: failed to unmarshal tenant network data: %v", err)
	}

	tenantNetworkId := 0

	if len(result) > 0 {
		tenantNetworkId = result[0].Id
	} else {
		return nil, fmt.Errorf("failed to get tenant network ID")
	}

	// Construct our client with all information
	client := &AlkiraClient{
		URI:                  apiUrl,
		Username:             username,
		Password:             password,
		TenantNetworkId:      strconv.Itoa(tenantNetworkId),
		Client:               retryClient,
		Provision:            provision,
		Validate:             validate,
		SerializationEnabled: enableSerialization,
		serializationTimeout: time.Duration(timeoutSeconds) * time.Second,
	}

	logf("DEBUG", "ALKIRA-API-SERIALIZATION-ENABLED: %v", client.SerializationEnabled)
	logf("DEBUG", "ALKIRA-API-SERIALIZATION-TIMEOUT: %v", client.serializationTimeout)

	return client, nil
}

// get retrieve a resource by sending a GET request
func (ac *AlkiraClient) get(uri string) ([]byte, string, error) {
	logf("DEBUG", "client-get URI: %s\n", uri)

	requestId := "client-" + uuid.New().String()
	request, _ := retryablehttp.NewRequest("GET", uri, nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", ac.Authorization)
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, "", fmt.Errorf("client-get(%s) failed to send request, %v", requestId, err)
	}

	defer response.Body.Close()
	logf("DEBUG", "client-get(%s): received response with status: %d", requestId, response.StatusCode)
	logf("DEBUG", "client-get(%s): response headers: %v", requestId, response.Header)
	data, err := io.ReadAll(response.Body)
	if err != nil {
		logf("ERROR", "client-get(%s): failed to read response body: %v", requestId, err)
		return nil, "", fmt.Errorf("client-get(%s) failed to read response body: %v", requestId, err)
	}
	logf("DEBUG", "client-get(%s) %d RSP: %s", requestId, response.StatusCode, string(data))
	logf("DEBUG", "client-get(%s): response body length: %d", requestId, len(data))

	if response.StatusCode != 200 {
		return nil, "", fmt.Errorf("client-get(%s): %d %s", requestId, response.StatusCode, string(data))
	}

	//
	// If provision is enabled, try to grab the provisioning status
	//
	provisionState := response.Header.Get("x-provision-request-state")

	return data, provisionState, nil
}

// get retrieve a resource by sending a GET request
func (ac *AlkiraClient) getByName(uri string) ([]byte, string, error) {
	logf("DEBUG", "client-get URI: %s\n", uri)

	requestId := "client-" + uuid.New().String()
	request, _ := retryablehttp.NewRequest("GET", uri, nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", ac.Authorization)
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, "", fmt.Errorf("client-get(%s): failed to send request, %v", requestId, err)
	}

	defer response.Body.Close()
	logf("DEBUG", "client-get(%s): received response with status: %d", requestId, response.StatusCode)
	logf("DEBUG", "client-get(%s): response headers: %v", requestId, response.Header)
	data, err := io.ReadAll(response.Body)
	if err != nil {
		logf("ERROR", "client-get(%s): failed to read response body: %v", requestId, err)
		return nil, "", fmt.Errorf("client-get(%s): failed to read response body: %v", requestId, err)
	}
	logf("DEBUG", "client-get(%s) %d RSP: %s", requestId, response.StatusCode, string(data))
	logf("DEBUG", "client-get(%s): response body length: %d", requestId, len(data))

	if response.StatusCode != 200 {
		return nil, "", fmt.Errorf("client-get-by-name(%s): %d %s", requestId, response.StatusCode, string(data))
	}

	//
	// If provision is enabled, try to grab the provisioning status
	//
	// Also, this header is only available when `?name=` query
	// parameter is used.
	//
	provisionState := response.Header.Get("x-provision-request-state")

	return data, provisionState, nil
}

// executeWithMutex executes a function while holding the API mutex if serialization is enabled
// Returns an error if the mutex cannot be acquired within the configured timeout
func (ac *AlkiraClient) executeWithMutex(fn func() error) error {
	// If serialization is disabled, execute immediately
	if !ac.SerializationEnabled {
		return fn()
	}

	// Channel to signal mutex acquisition
	mutexAcquired := make(chan struct{})

	// Try to acquire the mutex in a goroutine
	go func() {
		ac.apiMutex.Lock()
		close(mutexAcquired)
	}()

	// Wait for either mutex acquisition or timeout
	select {
	case <-mutexAcquired:
		// Mutex acquired successfully
		defer ac.apiMutex.Unlock()
		logf("DEBUG", "API mutex acquired, executing request")
		return fn()
	case <-time.After(ac.serializationTimeout):
		// Timeout occurred
		return fmt.Errorf("failed to acquire API mutex within timeout (%v)", ac.serializationTimeout)
	}
}

// create send a POST request to create resource
func (ac *AlkiraClient) create(uri string, body []byte, provision bool) ([]byte, string, error, error, error) {
	logf("DEBUG", "client-create REQ: %s", string(body))

	//
	// There are two knobs here to support turning provision on/off
	// globally through ENV var and to support APIs that doesn't need
	// to provision.
	//
	parsedURL, urlErr := url.Parse(uri)
    if urlErr != nil {
        return nil, "", fmt.Errorf("client-create: failed to parse URI: %w", urlErr), nil, nil
    }

    query := parsedURL.Query()

    if ac.Provision && provision {
        logf("DEBUG", "client-create: enable provision")
       	query.Set("provision", "true")
    }
    if ac.Validate {
       	logf("DEBUG", "client-create: enable async validation")
       	query.Set("async", "true")
    }

    parsedURL.RawQuery = query.Encode()
    uri = parsedURL.String()

	requestId := "client-" + uuid.New().String()
	request, _ := retryablehttp.NewRequest("POST", uri, bytes.NewBuffer(body))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", ac.Authorization)
	request.Header.Set("x-ak-request-id", requestId)

	// Execute the HTTP request with serialization if enabled
	var response *http.Response
	var err error
	mutexErr := ac.executeWithMutex(func() error {
		response, err = ac.Client.Do(request)
		return err
	})

	if mutexErr != nil {
		return nil, "", fmt.Errorf("client-create(%s): %v", requestId, mutexErr), nil, nil
	}

	if err != nil {
		return nil, "", fmt.Errorf("client-create(%s): failed to send request, %v", requestId, err), nil, nil
	}

	defer response.Body.Close()
	logf("DEBUG", "client-create(%s): received response with status: %d", requestId, response.StatusCode)
	logf("DEBUG", "client-create(%s): response headers: %v", requestId, response.Header)
	data, err := io.ReadAll(response.Body)
	if err != nil {
		logf("ERROR", "client-create(%s): failed to read response body: %v", requestId, err)
		return nil, "", fmt.Errorf("client-create(%s): failed to read response body: %v", requestId, err), nil, nil
	}

	logf("DEBUG", "client-create(%s) %d RSP: %s", requestId, response.StatusCode, string(data))
	logf("DEBUG", "client-create(%s): response body length: %d", requestId, len(data))

	if response.StatusCode != 201 && response.StatusCode != 200 && response.StatusCode != 202 {
		return nil, "", fmt.Errorf("client-create(%s): %d %s.", requestId, response.StatusCode, string(data)), nil, nil
	}

	// Handle validation if enabled and response is 202
	if ac.Validate && response.StatusCode == 202 {
		// Handle validation
		err := ac.handleValidation(response)
		if err != nil {
			return data, "", nil, err, nil
		}

		return data, "", nil, nil, nil
	}

	//
	// If provision is enabled, wait for provision to finish and
	// return the provision state
	//
	if ac.Provision == true && provision == true {
		provisionRequestId := response.Header.Get("x-provision-request-id")

		if provisionRequestId == "" {
			return data, "FAILED", nil, nil, fmt.Errorf("client-create(%s): failed to get provision request ID", requestId)
		}

		err := wait.Poll(10*time.Second, defaultProvTimeout, func() (bool, error) {
			request, err := ac.GetTenantNetworkProvisionRequest(provisionRequestId)

			if err != nil {
				return false, err
			}

			if request.State == "SUCCESS" {
				return true, nil
			} else if request.State == "FAILED" || request.State == "PARTIAL_SUCCESS" {
				return false, fmt.Errorf("client-create(%s): provision request %s failed", requestId, provisionRequestId)
			}

			logf("DEBUG", "client-create(%s): waiting for provision request %s to finish. (state: %s)", requestId, provisionRequestId, request.State)
			return false, nil
		})

		if err != nil {
			if err == wait.ErrWaitTimeout {
				return data, "FAILED", nil, nil, fmt.Errorf("client-create(%s): provision request %s timed out", requestId, provisionRequestId)
			}

			return data, "FAILED", nil, nil, err
		}

		return data, "SUCCESS", nil, nil, nil
	}

	return data, "", nil, nil, nil
}

// delete send a DELETE request to delete a resource
func (ac *AlkiraClient) delete(uri string, provision bool) (string, error, error, error) {
	logf("DEBUG", "client-delete: URI %s\n", uri)

	//
	// There are two knobs here to support turning provision on/off
	// globally through ENV var and to support APIs that doesn't need
	// to provision.
	//
	parsedURL, urlErr := url.Parse(uri)
    if urlErr != nil {
        return "", fmt.Errorf("client-delete: failed to parse URI: %w", urlErr), nil, nil
    }

    query := parsedURL.Query()

    if ac.Provision && provision {
        logf("DEBUG", "client-delete: enable provision")
       	query.Set("provision", "true")
    }
    if ac.Validate {
       	logf("DEBUG", "client-delete: enable async validation")
       	query.Set("async", "true")
    }

    parsedURL.RawQuery = query.Encode()
    uri = parsedURL.String()

	requestId := "client-" + uuid.New().String()
	request, _ := retryablehttp.NewRequest("DELETE", uri, nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", ac.Authorization)
	request.Header.Set("x-ak-request-id", requestId)

	// Execute the HTTP request with serialization if enabled
	var response *http.Response
	var err error
	mutexErr := ac.executeWithMutex(func() error {
		response, err = ac.Client.Do(request)
		return err
	})

	if mutexErr != nil {
		return "", fmt.Errorf("client-delete(%s): %v", requestId, mutexErr), nil, nil
	}

	if err != nil {
		return "", fmt.Errorf("client-delete(%s): failed to send request, %v", requestId, err), nil, nil
	}

	defer response.Body.Close()
	logf("DEBUG", "client-delete(%s): received response with status: %d", requestId, response.StatusCode)
	logf("DEBUG", "client-delete(%s): response headers: %v", requestId, response.Header)
	data, err := io.ReadAll(response.Body)
	if err != nil {
		logf("ERROR", "client-delete(%s): failed to read response body: %v", requestId, err)
		return "", fmt.Errorf("client-delete(%s): failed to read response body: %v", requestId, err), nil, nil
	}

	logf("DEBUG", "client-delete(%s): %d RSP: %s\n", requestId, response.StatusCode, string(data))
	logf("DEBUG", "client-delete(%s): response body length: %d", requestId, len(data))

	// Handle validation if enabled and response is 202
	if ac.Validate && response.StatusCode == 202 {
		// Handle validation
		err := ac.handleValidation(response)
		if err != nil {
			return "", nil, err, nil
		}

		return "", nil, nil, nil
	}

	if response.StatusCode < 200 || response.StatusCode > 299 {
		if response.StatusCode == 404 {
			logf("INFO", "client-delete(%s): %d resource was already deleted.\n", requestId, response.StatusCode)
			return "", nil, nil, nil
		}

		return "", fmt.Errorf("client-delete(%s): %d %s", requestId, response.StatusCode, string(data)), nil, nil
	}

	// If provision is enabled, wait for provision to finish and
	// return the proper provision state
	if ac.Provision == true && provision == true {
		provisionRequestId := response.Header.Get("x-provision-request-id")

		if provisionRequestId == "" {
			return "FAILED", nil, nil, fmt.Errorf("client-delete(%s): failed to get provision request ID", requestId)
		}

		err := wait.Poll(10*time.Second, defaultProvTimeout, func() (bool, error) {
			request, err := ac.GetTenantNetworkProvisionRequest(provisionRequestId)

			if err != nil {
				return false, err
			}

			if request.State == "SUCCESS" {
				return true, nil
			} else if request.State == "FAILED" {
				return false, fmt.Errorf("client-delete(%s): provision request %s failed", requestId, provisionRequestId)
			}

			logf("DEBUG", "client-delete(%s): waiting for provision request %s to finish. (state: %s)", requestId, provisionRequestId, request.State)
			return false, nil
		})

		if err != nil {
			if err == wait.ErrWaitTimeout {
				return "FAILED", nil, nil, fmt.Errorf("client-delete(%s): provision request %s timed out", requestId, provisionRequestId)
			}

			return "FAILED", nil, nil, err
		}

		return "SUCCESS", nil, nil, nil
	}

	return "", nil, nil, nil
}

// update send a PUT request to update a resource
func (ac *AlkiraClient) update(uri string, body []byte, provision bool) (string, error, error, error) {
	logf("DEBUG", "client-update: REQ: %s\n", string(body))

	//
	// There are two knobs here to support turning provision on/off
	// globally through the flag and to support APIs that doesn't need
	// to provision.
	//
	parsedURL, urlErr := url.Parse(uri)
    if urlErr != nil {
        return "", fmt.Errorf("client-update: failed to parse URI: %w", urlErr), nil, nil
    }

    query := parsedURL.Query()

    if ac.Provision && provision {
        logf("DEBUG", "client-update: enable provision")
       	query.Set("provision", "true")
    }
    if ac.Validate {
       	logf("DEBUG", "client-update: enable async validation")
       	query.Set("async", "true")
    }

    parsedURL.RawQuery = query.Encode()
    uri = parsedURL.String()

	requestId := "client-" + uuid.New().String()
	request, _ := retryablehttp.NewRequest("PUT", uri, bytes.NewBuffer(body))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", ac.Authorization)
	request.Header.Set("x-ak-request-id", requestId)

	// Execute the HTTP request with serialization if enabled
	var response *http.Response
	var err error
	mutexErr := ac.executeWithMutex(func() error {
		response, err = ac.Client.Do(request)
		return err
	})

	if mutexErr != nil {
		return "", fmt.Errorf("client-update(%s): %v", requestId, mutexErr), nil, nil
	}

	if err != nil {
		return "", fmt.Errorf("client-update(%s): failed to send request, %v", requestId, err), nil, nil
	}

	defer response.Body.Close()
	logf("DEBUG", "client-update(%s): received response with status: %d", requestId, response.StatusCode)
	logf("DEBUG", "client-update(%s): response headers: %v", requestId, response.Header)
	data, err := io.ReadAll(response.Body)
	if err != nil {
		logf("ERROR", "client-update(%s): failed to read response body: %v", requestId, err)
		return "", fmt.Errorf("client-update(%s): failed to read response body: %v", requestId, err), nil, nil
	}

	logf("DEBUG", "client-update(%s): %d RSP: %s\n", requestId, response.StatusCode, string(data))
	logf("DEBUG", "client-update(%s): response body length: %d", requestId, len(data))

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return "", fmt.Errorf("client-update(%s): %d %s", requestId, response.StatusCode, string(data)), nil, nil
	}

	// Handle validation if enabled and response is 202
	if ac.Validate && response.StatusCode == 202 {
		// Handle validation
		err := ac.handleValidation(response)
		if err != nil {
			return "", nil, err, nil
		}

		return "", nil, nil, nil
	}

	//
	// If provision is enabled, wait for provision to finish and return the proper state
	//
	if ac.Provision == true && provision == true {
		provisionRequestId := response.Header.Get("x-provision-request-id")

		if provisionRequestId == "" {
			return "FAILED", nil, nil, fmt.Errorf("client-update(%s): failed to get provision request ID", requestId)
		}

		err := wait.Poll(10*time.Second, defaultProvTimeout, func() (bool, error) {
			request, err := ac.GetTenantNetworkProvisionRequest(provisionRequestId)

			if err != nil {
				return false, err
			}

			if request.State == "SUCCESS" {
				return true, nil
			} else if request.State == "FAILED" {
				return false, fmt.Errorf("client-update(%s): provision request %s failed", requestId, provisionRequestId)
			}

			logf("DEBUG", "client-update(%s): waiting for provision request %s to finish. (state: %s)", requestId, provisionRequestId, request.State)
			return false, nil
		})

		if err != nil {
			if err == wait.ErrWaitTimeout {
				return "FAILED", nil, nil, fmt.Errorf("client-update(%s): provision request %s timed out", requestId, provisionRequestId)
			}

			return "FAILED", nil, nil, err
		}

		return "SUCCESS", nil, nil, nil
	}

	return "", nil, nil, nil
}
