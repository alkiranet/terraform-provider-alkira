// Copyright (C) 2020-2021 Alkira Inc. All Rights Reserved.

package alkira

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

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
func NewAlkiraClient(url string, username string, password string) (*AlkiraClient, error) {

	// Construct the complete URI based on the given endpoint
	apiUrl := "https://" + url + "/api"

	loginRequestBody, err := json.Marshal(map[string]string{
		"userName": username,
		"password": password,
	})

	// Login to the portal
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// Using a client to set a timeout. This is alkira service. It
	// should not take that long
	var httpClient = &http.Client{
		Timeout:   time.Second * 30,
		Transport: tr,
	}

	jar := &Session{}
	jar.jar = make(map[string][]*http.Cookie)
	httpClient.Jar = jar

	loginUrl := fmt.Sprintf("%s/user/login", apiUrl)

	request, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(loginRequestBody))
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

	sessionRequest, err := http.NewRequest("POST", sessionUrl, bytes.NewBuffer(userAuthData))
	sessionRequest.Header.Set("Content-Type", "application/json")
	sessionResponse, err := httpClient.Do(sessionRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to make session request, %v", err)
	}

	defer sessionResponse.Body.Close()

	sessionData, _ := ioutil.ReadAll(sessionResponse.Body)

	if sessionResponse.StatusCode != 200 {
		log.Println(string(sessionData))
		return nil, fmt.Errorf("failed to get session (%d)", sessionResponse.StatusCode)
	}

	// Get the tenant network ID
	var result []TenantNetworkId
	tenantNetworkUrl := apiUrl + "/tenantnetworks"

	tenantNetworkRequest, err := http.NewRequest("GET", tenantNetworkUrl, nil)
	tenantNetworkRequest.Header.Set("Content-Type", "application/json")
	tenantNetworkResponse, err := httpClient.Do(tenantNetworkRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to make tenant network request, %v", err)
	}

	defer tenantNetworkResponse.Body.Close()

	data, _ := ioutil.ReadAll(tenantNetworkResponse.Body)

	if tenantNetworkResponse.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get session (%d)", sessionResponse.StatusCode)
	}

	json.Unmarshal([]byte(data), &result)

	tenantNetworkId := 0

	if len(result) > 0 {
		tenantNetworkId = result[0].Id
	} else {
		return nil, fmt.Errorf("failed to get Tenant Network Id")
	}

	// Construct our client with all information
	client := &AlkiraClient{URI: apiUrl, Username: username, Password: password, TenantNetworkId: strconv.Itoa(tenantNetworkId), Client: httpClient}

	return client, nil
}

// get retrieve a resource by sending a GET request
func (ac *AlkiraClient) get(uri string) ([]byte, error) {
	request, err := http.NewRequest("GET", uri, nil)

	if err != nil {
		return nil, fmt.Errorf("request(GET) failed: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("request(GET) failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return nil, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return data, nil
}

// create send a POST request to create resource
func (ac *AlkiraClient) create(uri string, body []byte) ([]byte, error) {
	logf("DEBUG", "request(POST): %s\n", string(body))

	request, err := http.NewRequest("POST", uri, bytes.NewBuffer(body))

	if err != nil {
		return nil, fmt.Errorf("request(POST) failed: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := ac.Client.Do(request)

	if err != nil {
		return nil, fmt.Errorf("request(POST) failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 201 {
		return nil, fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return data, nil
}

// delete send a DELETE request to delete a resource
func (ac *AlkiraClient) delete(uri string) error {
	logf("DEBUG", "request(DELETE) uri: %s\n", uri)

	request, err := http.NewRequest("DELETE", uri, nil)

	if err != nil {
		return fmt.Errorf("request(DELETE) failed: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("request(DELETE) failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
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

	request, err := http.NewRequest("PUT", uri, bytes.NewBuffer(body))

	if err != nil {
		return fmt.Errorf("request(PUT) failed: %v", err)
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return fmt.Errorf("request(PUT): failed, %v", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	if response.StatusCode != 200 {
		return fmt.Errorf("(%d) %s", response.StatusCode, string(data))
	}

	return nil
}
