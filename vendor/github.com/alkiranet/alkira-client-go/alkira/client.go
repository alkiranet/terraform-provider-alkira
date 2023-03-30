// Copyright (C) 2020-2023 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/wait"
)

// Default client timeout is 60s and provision timeout is 240m
const defaultClientTimeout time.Duration = 60 * time.Second
const defaultProvTimeout time.Duration = 240 * time.Minute

type AlkiraClient struct {
	Client          *http.Client
	Password        string
	Provision       bool
	TenantNetworkId string
	URI             string
	Username        string
}

type Session struct {
	jar map[string][]*http.Cookie
}

func (s *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
	logf("TRACE", "SetCookies URL : %s\n", u.String())
	logf("TRACE", "SetCookies: %s\n", cookies)
	s.jar[u.Host] = cookies
}

func (s *Session) Cookies(u *url.URL) []*http.Cookie {
	logf("TRACE", "Cookie URL is : %s\n", u.String())
	logf("TRACE", "Cookie being returned is : %s\n", s.jar[u.Host])
	return s.jar[u.Host]
}

// NewAlkiraClient creates a new API client
func NewAlkiraClient(hostname string, username string, password string, provision bool) (*AlkiraClient, error) {

	// Construct the portal URI
	url := "https://" + hostname

	// Set the client timeout
	clientTimeout := defaultClientTimeout

	if t := os.Getenv("ALKIRA_CLIENT_TIMEOUT"); t != "" {
		var err error
		clientTimeout, err = time.ParseDuration(t)

		if err != nil {
			return nil, fmt.Errorf("failed to parse ENV variable ALKIRA_CLIENT_TIMEOUT, %v", err)
		}
	}

	logf("DEBUG", "PROVISION: %v", provision)
	return NewAlkiraClientInternal(url, username, password, clientTimeout, provision)
}

// NewAlkiraClientInternal creates a new internal Alkira client
func NewAlkiraClientInternal(url string, username string, password string, timeout time.Duration, provision bool) (*AlkiraClient, error) {

	// Construct the portal URI based on the given endpoint
	apiUrl := url + "/api"

	loginRequestBody, err := json.Marshal(map[string]string{
		"userName": username,
		"password": password,
	})

	// Login to the portal
	tr := &http.Transport{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	var httpClient = &http.Client{
		Timeout:   timeout,
		Transport: tr,
	}

	jar := &Session{}
	jar.jar = make(map[string][]*http.Cookie)
	httpClient.Jar = jar

	// User login
	loginUrl := fmt.Sprintf("%s/user/login", apiUrl)

	request, requestErr := http.NewRequest("POST", loginUrl, bytes.NewBuffer(loginRequestBody))

	if requestErr != nil {
		return nil, fmt.Errorf("failed to create login request, %v", requestErr)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := httpClient.Do(request)

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

	sessionRequest, _ := http.NewRequest("POST", sessionUrl, bytes.NewBuffer(userAuthData))
	sessionRequest.Header.Set("Content-Type", "application/json")
	sessionResponse, err := httpClient.Do(sessionRequest)

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
	tenantNetworkUrl := apiUrl + "/tenantnetworks"

	tenantNetworkRequest, _ := http.NewRequest("GET", tenantNetworkUrl, nil)
	tenantNetworkRequest.Header.Set("Content-Type", "application/json")
	tenantNetworkResponse, err := httpClient.Do(tenantNetworkRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to make tenant network request, %v", err)
	}

	defer tenantNetworkResponse.Body.Close()

	data, _ := ioutil.ReadAll(tenantNetworkResponse.Body)
	logf("TRACE", "tenant network: %s\n", string(data))

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
		Client:          httpClient,
		Provision:       provision,
	}

	return client, nil
}

// get retrieve a resource by sending a GET request
func (ac *AlkiraClient) get(uri string) ([]byte, string, error) {
	logf("DEBUG", "client-get URI: %s\n", uri)

	requestId := "client-" + uuid.New().String()
	request, _ := http.NewRequest("GET", uri, nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, "", fmt.Errorf("client-get %s failed, %v", requestId, err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "client-get(%s) RSP: %s\n", requestId, string(data))

	if response.StatusCode != 200 {
		return nil, "", fmt.Errorf("%s(%d): %s", requestId, response.StatusCode, string(data))
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
	request, _ := http.NewRequest("GET", uri, nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, "", fmt.Errorf("client-get(%s): failed, %v", requestId, err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "client-get(%s) RSP: %s\n", requestId, string(data))

	if response.StatusCode != 200 {
		return nil, "", fmt.Errorf("%s(%d): %s", requestId, response.StatusCode, string(data))
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
	logf("DEBUG", "client-create REQ: %s\n", string(body))

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
	request, _ := http.NewRequest("POST", uri, bytes.NewBuffer(body))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, "", fmt.Errorf("client-create(%s): failed to send request, %v", requestId, err), nil
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	logf("DEBUG", "client-create(%s) RSP: %s\n", requestId, string(data))

	if response.StatusCode != 201 && response.StatusCode != 200 {
		return nil, "", fmt.Errorf("%s(%d) %s", requestId, response.StatusCode, string(data)), nil
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

		if err == wait.ErrWaitTimeout {
			return data, "FAILED", nil, fmt.Errorf("client-create(%s): provision request %s timed out", requestId, provisionRequestId)
		}

		if err != nil {
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
	request, _ := http.NewRequest("DELETE", uri, nil)

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("client-delete(%s): failed, %v", requestId, err), nil
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	logf("DEBUG", "client-delete(%s): RSP: %s\n", requestId, string(data))

	if response.StatusCode != 200 && response.StatusCode != 202 {
		if response.StatusCode == 404 {
			logf("INFO", "client-delete(%s): resource was already deleted.\n", requestId)
			return "", nil, nil
		}

		// Retry several more times and see if the delete goes through
		err := wait.Poll(2*time.Second, 10*time.Second, func() (bool, error) {
			response, err = ac.Client.Do(request)
			if err != nil {
				return false, err
			}

			data, _ = ioutil.ReadAll(response.Body)

			if response.StatusCode == 200 || response.StatusCode == 202 || response.StatusCode == 404 {
				return true, nil
			}

			logf("INFO", "client-delete(%s): retry delete (%d).", requestId, response.StatusCode)
			return false, nil
		})

		if err == wait.ErrWaitTimeout {
			return "", fmt.Errorf("client-delete(%s): retry timeout, %s", requestId, string(data)), nil
		}

		if err != nil {
			return "", fmt.Errorf("%s(%d): %s", requestId, response.StatusCode, string(data)), nil
		}
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

		if err == wait.ErrWaitTimeout {
			return "FAILED", nil, fmt.Errorf("client-delete(%s): provision request %s timed out", requestId, provisionRequestId)
		}

		if err != nil {
			return "FAILED", nil, err
		}

		return "SUCCESS", nil, nil
	}

	return "", nil, nil
}

// update send a PUT request to update a resource
func (ac *AlkiraClient) update(uri string, body []byte, provision bool) (string, error, error) {
	logf("DEBUG", "client-update: REQUEST: %s\n", string(body))

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
	request, _ := http.NewRequest("PUT", uri, bytes.NewBuffer(body))

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("x-ak-request-id", requestId)

	response, err := ac.Client.Do(request)

	if err != nil {
		return "", fmt.Errorf("client-update(%s): failed, %v", requestId, err), nil
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	logf("DEBUG", "client-update(%s): RSP: %s\n", requestId, string(data))

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return "", fmt.Errorf("%s(%d) %s", requestId, response.StatusCode, string(data)), nil
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

		if err == wait.ErrWaitTimeout {
			return "FAILED", nil, fmt.Errorf("client-update(%s): provision request %s timed out", requestId, provisionRequestId)
		}

		if err != nil {
			return "FAILED", nil, err
		}

		return "SUCCESS", nil, nil
	}

	return "", nil, nil
}
