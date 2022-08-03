// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

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

	"k8s.io/apimachinery/pkg/util/wait"
)

// Default client timeout is 60s
const defaultClientTimeout time.Duration = 60 * time.Second

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
	logf("DEBUG", "SetCookies URL : %s\n", u.String())
	logf("DEBUG", "SetCookies: %s\n", cookies)
	s.jar[u.Host] = cookies
}

func (s *Session) Cookies(u *url.URL) []*http.Cookie {
	logf("DEBUG", "Cookie URL is : %s\n", u.String())
	logf("DEBUG", "Cookie being returned is : %s\n", s.jar[u.Host])
	return s.jar[u.Host]
}

// NewAlkiraClient creates a new API client
func NewAlkiraClient(hostname string, username string, password string) (*AlkiraClient, error) {

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

	// With provision
	provision := false

	if os.Getenv("ALKIRA_PROVISION") == "TRUE" {
		provision = true
	}

	logf("DEBUG", "PROVISION: %v", provision)
	return NewAlkiraClientInternal(url, username, password, clientTimeout, provision)
}

// NewAlkiraClientInternal creates a new client
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
	logf("DEBUG", "session data: %s\n", string(sessionData))

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
	logf("DEBUG", "tenant network: %s\n", string(data))

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
	client := &AlkiraClient{URI: apiUrl, Username: username, Password: password, TenantNetworkId: strconv.Itoa(tenantNetworkId), Client: httpClient, Provision: provision}

	return client, nil
}

// get retrieve a resource by sending a GET request
func (ac *AlkiraClient) get(uri string) ([]byte, error) {
	logf("DEBUG", "request(GET) URI: %s\n", uri)

	request, _ := http.NewRequest("GET", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("request(GET) failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "request(GET) RSP: %s\n", string(data))

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return data, nil
}

// create send a POST request to create resource
func (ac *AlkiraClient) create(uri string, body []byte, provision bool) ([]byte, error) {
	logf("DEBUG", "client-create: REQUEST: %s\n", string(body))

	//
	// There are two knobs here to support turning provision on/off
	// globally through ENV var and to support APIs that doesn't need
	// to provision.
	//
	if ac.Provision == true && provision == true {
		logf("DEBUG", "client-create: enable provision")
		uri = fmt.Sprintf("%s?provision=true", uri)
	}

	request, _ := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("client-create: failed to send request, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "client-create: RESPONSE: %s\n", string(data))

	if response.StatusCode != 201 && response.StatusCode != 200 {
		return nil, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	// If provision is enabled, wait for provision to finish
	if ac.Provision == true && provision == true {
		provisionRequestId := response.Header.Get("x-provision-request-id")

		if provisionRequestId == "" {
			return nil, fmt.Errorf("client-create: failed to get provision request ID")
		}

		err := wait.Poll(10*time.Second, 120*time.Minute, func() (bool, error) {
			request, err := ac.GetTenantNetworkProvisionRequest(provisionRequestId)

			if err != nil {
				return false, err
			}

			if request.State == "SUCCESS" || request.State == "PARTIAL_SUCCESS" {
				return true, nil
			}
			if request.State == "FAILED" {
				return false, fmt.Errorf("client-create: provision request %s failed", provisionRequestId)
			}

			logf("DEBUG", "client-create: waiting for provision to finish.")
			return false, nil
		})

		if err == wait.ErrWaitTimeout {
			return nil, fmt.Errorf("client-create: timed out waiting for provision to complete")
		}

		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

// delete send a DELETE request to delete a resource
func (ac *AlkiraClient) delete(uri string, provision bool) error {
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

	request, _ := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("client-delete: failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "client-delete: RESPONSE: %s\n", string(data))

	if response.StatusCode != 200 && response.StatusCode != 202 {
		if response.StatusCode == 404 {
			logf("INFO", "client-delete: resource was already deleted.\n")
			return nil
		}
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	// If provision is enabled, wait for provision to finish
	if ac.Provision == true && provision == true {
		provisionRequestId := response.Header.Get("x-provision-request-id")

		if provisionRequestId == "" {
			return fmt.Errorf("client-delete: failed to get provision request ID")
		}

		err := wait.Poll(10*time.Second, 120*time.Minute, func() (bool, error) {
			request, err := ac.GetTenantNetworkProvisionRequest(provisionRequestId)

			if err != nil {
				return false, err
			}

			if request.State == "SUCCESS" || request.State == "PARTIAL_SUCCESS" {
				return true, nil
			}
			if request.State == "FAILED" {
				return false, fmt.Errorf("client-delete: provision request %s failed", provisionRequestId)
			}

			logf("DEBUG", "client-delete: waiting for provision to finish.")
			return false, nil
		})

		if err == wait.ErrWaitTimeout {
			return fmt.Errorf("client-delete: timed out waiting for provision to complete")
		}

		return err
	}

	return nil
}

// update send a PUT request to update a resource
func (ac *AlkiraClient) update(uri string, body []byte, provision bool) error {
	logf("DEBUG", "client-update: REQUEST: %s\n", string(body))

	//
	// There are two knobs here to support turning provision on/off
	// globally through ENV var and to support APIs that doesn't need
	// to provision.
	//
	if ac.Provision == true && provision == true {
		logf("DEBUG", "client-update: enable provision")
		uri = fmt.Sprintf("%s?provision=true", uri)
	}

	request, _ := http.NewRequest("PUT", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("client-update: failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "client-update: RESPONSE: %s\n", string(data))

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	// If provision is enabled, wait for provision to finish
	if ac.Provision == true && provision == true {
		provisionRequestId := response.Header.Get("x-provision-request-id")

		if provisionRequestId == "" {
			return fmt.Errorf("client-update: failed to get provision request ID")
		}

		err := wait.Poll(10*time.Second, 120*time.Minute, func() (bool, error) {
			request, err := ac.GetTenantNetworkProvisionRequest(provisionRequestId)

			if err != nil {
				return false, err
			}

			// The provision states could be misleading in certain
			// cases. For "PARTIAL_SUCCESS", provisioning of some
			// resources actaully failed. For "FAILED" state, some
			// resources may get provisioned successfully.
			if request.State == "SUCCESS" || request.State == "PARTIAL_SUCCESS" {
				return true, nil
			}
			if request.State == "FAILED" {
				return false, fmt.Errorf("client-update: provision request %s failed", provisionRequestId)
			}

			logf("DEBUG", "client-update: waiting for provision to finish.")
			return false, nil
		})

		if err == wait.ErrWaitTimeout {
			return fmt.Errorf("client-update: timed out waiting for provision to complete")
		}

		return err
	}

	return nil
}
