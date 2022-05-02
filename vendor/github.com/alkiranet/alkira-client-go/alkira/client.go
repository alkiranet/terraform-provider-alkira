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
)

// Default client timeout is 60s
const defaultClientTimeout time.Duration = 60 * time.Second

type AlkiraClient struct {
	Username        string
	Password        string
	URI             string
	TenantNetworkId string
	Client          *http.Client
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

	return NewAlkiraClientInternal(url, username, password, clientTimeout)
}

// NewAlkiraClientInternal creates a new client
func NewAlkiraClientInternal(url string, username string, password string, timeout time.Duration) (*AlkiraClient, error) {

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
	client := &AlkiraClient{URI: apiUrl, Username: username, Password: password, TenantNetworkId: strconv.Itoa(tenantNetworkId), Client: httpClient}

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
func (ac *AlkiraClient) create(uri string, body []byte) ([]byte, error) {
	logf("DEBUG", "request(POST): %s\n", string(body))

	request, _ := http.NewRequest("POST", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("request(POST) failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "request(POST) RSP: %s\n", string(data))

	if response.StatusCode != 201 && response.StatusCode != 200 {
		return nil, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return data, nil
}

// delete send a DELETE request to delete a resource
func (ac *AlkiraClient) delete(uri string) error {
	logf("DEBUG", "request(DEL) uri: %s\n", uri)

	request, _ := http.NewRequest("DELETE", uri, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("request(DEL) failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "request(DEL) RSP: %s\n", string(data))

	if response.StatusCode != 200 && response.StatusCode != 202 {
		if response.StatusCode == 404 {
			logf("INFO", "resource was already deleted.\n")
			return nil
		}
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}

// update send a PUT request to update a resource
func (ac *AlkiraClient) update(uri string, body []byte) error {
	logf("DEBUG", "request(PUT): %s\n", string(body))

	request, _ := http.NewRequest("PUT", uri, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("request(PUT): failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)
	logf("DEBUG", "request(PUT) RSP: %s\n", string(data))

	if response.StatusCode != 200 && response.StatusCode != 202 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
