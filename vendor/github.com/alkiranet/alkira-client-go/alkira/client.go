// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/go-retryablehttp"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Default provision timeout is 240m
const defaultProvTimeout time.Duration = 240 * time.Minute

// Default Retry
const defaultRetryInterval time.Duration = 5 * time.Second
const defaultRetryTimeout time.Duration = 60 * time.Second

type AlkiraClient struct {
	Client          *retryablehttp.Client
	URI             string
	Username        string
	Password        string
	Secret          string
	Authorization   string
	Provision       bool
	TenantNetworkId string
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
func NewAlkiraClient(hostname string, username string, password string, secret string, provision bool, auth string) (*AlkiraClient, error) {

	// Construct the portal URI
	url := "https://" + hostname

	logf("DEBUG", "ALKIRA-PROVISION: %v", provision)

	if auth == "header" {
		logf("DEBUG", "ALKIRA-AUTH-METHOD: %v", auth)
		return NewAlkiraClientWithAuthHeader(url, username, password, secret, provision)
	}

	return NewAlkiraClientInternal(url, username, password, secret, provision)
}

// NewAlkiraClientWithAuthHeader creates a new internal Alkira client with authentication in header
func NewAlkiraClientWithAuthHeader(url string, username string, password string, secret string, provision bool) (*AlkiraClient, error) {

	// Firstly, construct the portal API based URI
	apiUrl := url + "/api"

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
	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {

		shouldRetry, e := retryablehttp.DefaultRetryPolicy(ctx, resp, err)

		if !shouldRetry {
			if resp.StatusCode == 500 {
				return true, fmt.Errorf("%s", resp.Status)
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

	data, _ := ioutil.ReadAll(tenantNetworkResponse.Body)
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
		Client:          retryClient,
		URI:             apiUrl,
		Username:        username,
		Password:        password,
		Secret:          secret,
		Authorization:   auth,
		Provision:       provision,
		TenantNetworkId: strconv.Itoa(tenantNetworkId),
	}

	return client, nil
}

// NewAlkiraClientInternal creates a new internal Alkira client
func NewAlkiraClientInternal(url string, username string, password string, secret string, provision bool) (*AlkiraClient, error) {

	// Construct the portal URI based on the given endpoint
	apiUrl := url + "/api"

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

	retryClient.CheckRetry = func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		shouldRetry, e := retryablehttp.DefaultRetryPolicy(ctx, resp, err)

		// In addition, retry on 409 as well due to DELETE
		if !shouldRetry && resp != nil {
			if resp.StatusCode == 409 {
				return true, fmt.Errorf("%s", resp.Status)
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

	userAuthData, _ := ioutil.ReadAll(response.Body)

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

	sessionData, _ := ioutil.ReadAll(sessionResponse.Body)
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

	data, _ := ioutil.ReadAll(tenantNetworkResponse.Body)
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
		URI:             apiUrl,
		Username:        username,
		Password:        password,
		TenantNetworkId: strconv.Itoa(tenantNetworkId),
		Client:          retryClient,
		Provision:       provision,
	}

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
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "client-get(%s) %d RSP: %s", requestId, response.StatusCode, string(data))

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
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "client-get(%s) %d RSP: %s", requestId, response.StatusCode, string(data))

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

// create send a POST request to create resource
func (ac *AlkiraClient) create(uri string, body []byte, provision bool) ([]byte, string, error, error) {
	logf("DEBUG", "client-create REQ: %s", string(body))

	//
	// There are two knobs here to support turning provision on/off
	// globally through ENV var and to support APIs that doesn't need
	// to provision.
	//
	if ac.Provision == true && provision == true {
		logf("DEBUG", "client-create: enable provision")
		uri = fmt.Sprintf("%s?provision=true", uri)
	}

	requestId := "client-" + uuid.New().String()
	request, _ := retryablehttp.NewRequest("POST", uri, bytes.NewBuffer(body))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", ac.Authorization)
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, "", fmt.Errorf("client-create(%s): failed to send request, %v", requestId, err), nil
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	logf("DEBUG", "client-create(%s) %d RSP: %s", requestId, response.StatusCode, string(data))

	if response.StatusCode != 201 && response.StatusCode != 200 {
		return nil, "", fmt.Errorf("client-create(%s): %d %s.", requestId, response.StatusCode, string(data)), nil
	}

	//
	// If provision is enabled, wait for provision to finish and
	// return the provision state
	//
	if ac.Provision == true && provision == true {
		provisionRequestId := response.Header.Get("x-provision-request-id")

		if provisionRequestId == "" {
			return data, "FAILED", nil, fmt.Errorf("client-create(%s): failed to get provision request ID", requestId)
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
				return data, "FAILED", nil, fmt.Errorf("client-create(%s): provision request %s timed out", requestId, provisionRequestId)
			}

			return data, "FAILED", nil, err
		}

		return data, "SUCCESS", nil, nil
	}

	return data, "", nil, nil
}

// delete send a DELETE request to delete a resource
func (ac *AlkiraClient) delete(uri string, provision bool) (string, error, error) {
	logf("DEBUG", "client-delete: URI %s\n", uri)

	//
	// There are two knobs here to support turning provision on/off
	// globally through ENV var and to support APIs that doesn't need
	// to provision.
	//
	if ac.Provision == true && provision == true {
		logf("DEBUG", "client-delete: enable provision")
		uri = fmt.Sprintf("%s?provision=true", uri)
	}

	requestId := "client-" + uuid.New().String()
	request, _ := retryablehttp.NewRequest("DELETE", uri, nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", ac.Authorization)
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("client-delete(%s): failed to send request, %v", requestId, err), nil
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	logf("DEBUG", "client-delete(%s): %d RSP: %s\n", requestId, response.StatusCode, string(data))

	if response.StatusCode < 200 || response.StatusCode > 299 {
		if response.StatusCode == 404 {
			logf("INFO", "client-delete(%s): %d resource was already deleted.\n", requestId, response.StatusCode)
			return "", nil, nil
		}

		return "", fmt.Errorf("client-delete(%s): %d %s", requestId, response.StatusCode, string(data)), nil
	}

	// If provision is enabled, wait for provision to finish and
	// return the proper provision state
	if ac.Provision == true && provision == true {
		provisionRequestId := response.Header.Get("x-provision-request-id")

		if provisionRequestId == "" {
			return "FAILED", nil, fmt.Errorf("client-delete(%s): failed to get provision request ID", requestId)
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
				return "FAILED", nil, fmt.Errorf("client-delete(%s): provision request %s timed out", requestId, provisionRequestId)
			}

			return "FAILED", nil, err
		}

		return "SUCCESS", nil, nil
	}

	return "", nil, nil
}

// update send a PUT request to update a resource
func (ac *AlkiraClient) update(uri string, body []byte, provision bool) (string, error, error) {
	logf("DEBUG", "client-update: REQ: %s\n", string(body))

	//
	// There are two knobs here to support turning provision on/off
	// globally through the flag and to support APIs that doesn't need
	// to provision.
	//
	if ac.Provision == true && provision == true {
		logf("DEBUG", "client-update: enable provision")
		uri = fmt.Sprintf("%s?provision=true", uri)
	}

	requestId := "client-" + uuid.New().String()
	request, _ := retryablehttp.NewRequest("PUT", uri, bytes.NewBuffer(body))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", ac.Authorization)
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("client-update(%s): failed to send request, %v", requestId, err), nil
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	logf("DEBUG", "client-update(%s): %d RSP: %v\n", requestId, response.StatusCode, data)

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return "", fmt.Errorf("client-update(%s): %d %s", requestId, response.StatusCode, string(data)), nil
	}

	//
	// If provision is enabled, wait for provision to finish and return the proper state
	//
	if ac.Provision == true && provision == true {
		provisionRequestId := response.Header.Get("x-provision-request-id")

		if provisionRequestId == "" {
			return "FAILED", nil, fmt.Errorf("client-update(%s): failed to get provision request ID", requestId)
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
				return "FAILED", nil, fmt.Errorf("client-update(%s): provision request %s timed out", requestId, provisionRequestId)
			}

			return "FAILED", nil, err
		}

		return "SUCCESS", nil, nil
	}

	return "", nil, nil
}
