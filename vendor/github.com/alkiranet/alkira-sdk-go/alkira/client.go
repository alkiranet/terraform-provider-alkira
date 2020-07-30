package alkira

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type AlkiraClient struct {
	Username         string
	Password         string
	URI              string
	TenantNetworkId  string
	Client           *http.Client
}

type Session struct {
    jar map[string] []*http.Cookie
}

func (s *Session) SetCookies(u *url.URL, cookies []*http.Cookie) {
    //log.Printf("The URL is : %s\n", u.String())
    //log.Printf("The cookie being set is : %s\n", cookies)
    s.jar [u.Host] = cookies
}

func (s *Session) Cookies(u *url.URL) []*http.Cookie {
    //log.Printf("The URL is : %s\n", u.String())
    //log.Printf("Cookie being returned is : %s\n", s.jar[u.Host])
    return s.jar[u.Host]
}


// New API client creates a new API client
func NewAlkiraClient(url string, username string, password string) (*AlkiraClient, error) {

	// Construct the complete URI based on the given endpoint
	apiUrl := "https://" + url + "/api/"

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
		Timeout:   time.Second * 4,
		Transport: tr,
	}

	jar := &Session{}
    jar.jar = make(map[string] []*http.Cookie)
    httpClient.Jar = jar

	loginUrl := apiUrl + "user/login"

	request, err := http.NewRequest("POST", loginUrl, bytes.NewBuffer(loginRequestBody))
	request.Header.Set("Content-Type", "application/json")
	response, err := httpClient.Do(request)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer response.Body.Close()

	userAuthData, _ := ioutil.ReadAll(response.Body)

	sessionUrl := apiUrl + "sessions"

	sessionRequest, err := http.NewRequest("POST", sessionUrl, bytes.NewBuffer(userAuthData))
	sessionRequest.Header.Set("Content-Type", "application/json")
	sessionResponse, err := httpClient.Do(sessionRequest)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer sessionResponse.Body.Close()

	// Get the tenant network ID
	var result []TenantNetworkId
	tenantNetworkUrl := apiUrl + "tenantnetworks"

	tenantNetworkRequest, err := http.NewRequest("GET", tenantNetworkUrl, nil)
	tenantNetworkRequest.Header.Set("Content-Type", "application/json")
	tenantNetworkResponse, err := httpClient.Do(tenantNetworkRequest)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer tenantNetworkResponse.Body.Close()

	data, _ := ioutil.ReadAll(tenantNetworkResponse.Body)

	json.Unmarshal([]byte(data), &result)

	tenantNetworkId := 0

	if len(result) > 0 {
		tenantNetworkId = result[0].Id
	}

	// Construct our client with all information
	client := &AlkiraClient{URI: apiUrl, Username: username, Password: password, TenantNetworkId: strconv.Itoa(tenantNetworkId), Client: httpClient}

	return client, nil
}

