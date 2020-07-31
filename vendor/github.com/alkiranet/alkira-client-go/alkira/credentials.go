package alkira

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type CredentialAwsVpc struct {
	Ec2AccessKey    string `json:"ec2AccessKey"`
	Ec2SecretKey    string `json:"ec2SecretKey"`
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

type Credentials struct {
	Name        string      `json:"name"`
	Credentials interface{} `json:"credentials"`
}

type CredentialResponse struct {
	Id string `json:"id"`
}

// Create new Credential for AWS-VPC Connector
func (ac *AlkiraClient) CreateCredentialAwsVpc(name string, accessKey string, secretKey string, authType string) (string, error) {
	credentialEndpoint := ac.URI + "api/credentials/awsvpc"

	// This body is not the normal JSON format...
	body, err := json.Marshal(Credentials{
		Name: name,
		Credentials: CredentialAwsVpc{
			Ec2AccessKey: accessKey,
			Ec2SecretKey: secretKey,
			Type: authType,
		},
	})

	request, err := http.NewRequest("POST", credentialEndpoint, bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		log.Printf("Error : %s", err)
	}

	defer response.Body.Close()
	data, _ := ioutil.ReadAll(response.Body)

	var result CredentialResponse

	json.Unmarshal([]byte(data), &result)

	if response.StatusCode != 200 {
		return result.Id, errors.New("Failed to save credential")
	}

	return result.Id, nil
}

// Delete Credential for AWS-VPC Connector
func (ac *AlkiraClient) DeleteCredentialAwsVpc(id string) (error) {
	credentialEndpoint := ac.URI + "api/credentials/awsvpc/" + id

	request, err := http.NewRequest("DELETE", credentialEndpoint, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to delete credential-aws-vpc (%v)", response.StatusCode)
	}

	return nil
}


// Create new Credential for AZURE-VNET Connector
func (ac *AlkiraClient) CreateCredentialAzureVnet(name string, applicationId string, secretKey string, subscriptionId string, tenantId string) (string, error) {
	credentialEndpoint := ac.URI + "api/credentials/azurevnet"
	credentialId := ""

	// This body is not the normal JSON format...
	body, err := json.Marshal(Credentials{
		Name: name,
		Credentials: CredentialAzureVnet{
			ApplicationId:  applicationId,
			SecretKey:      secretKey,
			SubscriptionId: subscriptionId,
			TenantId:       tenantId,
		},
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


// Delete Credential for AZURE-VNET Connector
func (ac *AlkiraClient) DeleteCredentialAzureVnet(id string) (error) {
	credentialEndpoint := ac.URI + "api/credentials/azurevnet/" + id

	request, err := http.NewRequest("DELETE", credentialEndpoint, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to delete credential-azure-vnet (%v)", response.StatusCode)
	}

	return nil
}

// Create new Credential for GCP-VPC Connector
func (ac *AlkiraClient) CreateCredentialGcpVpc(name string, credential *CredentialGcpVpc) (string, error) {
	credentialEndpoint := ac.URI + "api/credentials/gcpvpc"
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

// Delete Credential for AZURE-VNET Connector
func (ac *AlkiraClient) DeleteCredentialGcpVpc(id string) (error) {
	credentialEndpoint := ac.URI + "api/credentials/gcpvpc/" + id

	request, err := http.NewRequest("DELETE", credentialEndpoint, nil)
	request.Header.Set("Content-Type", "application/json")
	response, err := ac.Client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Failed to delete credential-gcp-vpc (%v)", response.StatusCode)
	}

	return nil
}
