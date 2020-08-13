package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

type CredentialAwsVpcKey struct {
	Ec2AccessKey    string `json:"ec2AccessKey"`
	Ec2SecretKey    string `json:"ec2SecretKey"`
	Type            string `json:"type"`
}

type CredentialAwsVpcRole struct {
	Ec2RoleArn      string `json:"ec2RoleArn"`
	Ec2ExternalId   string `json:"ec2ExternalId"`
	Type            string `json:"type"`
}

type CredentialAzureVnet struct {
	ApplicationId  string `json:"applicationId"`
	SecretKey      string `json:"secretKey"`
	SubscriptionId string `json:"subscriptionId"`
	TenantId       string `json:"tenantId"`
}

type CredentialGcpVpc struct {
	AuthProvider        string `json:"auth_provider_x509_cert_url"`
	AuthUri             string `json:"auth_uri"`
	ClientEmail         string `json:"client_email"`
	ClientId            string `json:"client_id"`
	ClientX509CertUrl   string `json:"client_x509_cert_url"`
	PrivateKey          string `json:"private_key"`
	PrivateKeyId        string `json:"private_key_id"`
	ProjectId           string `json:"project_id"`
	TokenUri            string `json:"token_uri"`
	Type                string `json:"type"`
}

type CredentialPan struct {
	LicenseKey          string `json:"licenseKey"`
	Password            string `json:"password"`
	Username            string `json:"userName"`
}

type CredentialPanInstance struct {
	AuthKey             string `json:"authKey"`
	AuthCode            string `json:"authCode"`
	LicenseKey          string `json:"licenseKey"`
	Password            string `json:"password"`
	Username            string `json:"userName"`
}

type Credentials struct {
	Name        string      `json:"name"`
	Credentials interface{} `json:"credentials"`
}

type CredentialResponse struct {
	Id string `json:"id"`
}

// Create new Credential
func (ac *AlkiraClient) CreateCredential(name string, credentialType string, credential interface{}) (string, error) {
	credentialEndpoint := ac.URI + "api/credentials/" + credentialType
	credentialId := ""

	// This body is not the normal JSON format...
	body, err := json.Marshal(Credentials{
		Name: name,
		Credentials: credential,
	})

	request, err := http.NewRequest("POST", credentialEndpoint, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return credentialId, err
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result CredentialResponse

	json.Unmarshal([]byte(data), &result)

	credentialId = result.Id

	if response.StatusCode != 200 {
		return credentialId, errors.New("Failed to save credential")
	}

	return credentialId, nil
}


// Delete Credential
func (ac *AlkiraClient) DeleteCredential(id string, credentialType string) (error) {
	credentialEndpoint := ac.URI + "api/credentials/" + credentialType + "/" + id

	request, err := http.NewRequest("DELETE", credentialEndpoint, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to delete credential %s (%v)", credentialType, response.StatusCode)
	}

	return nil
}
