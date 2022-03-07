// Copyright (C) 2020-2022 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type CredentialType string

const (
	CredentialTypeAkamaiProlexic   CredentialType = "akamai-prolexic"
	CredentialTypeArubaEdgeConnect                = "aruba-edge-connector-instances"
	CredentialTypeAwsVpc                          = "awsvpc"
	CredentialTypeAzureVnet                       = "azurevnet"
	CredentialTypeChkpFw                          = "chkp-fw"
	CredentialTypeChkpFwInstance                  = "chkp-fw-instance"
	CredentialTypeChkpFwManagement                = "chkp-fw-management-server"
	CredentialTypeCiscoSdwan                      = "ciscosdwan"
	CredentialTypeGcpVpc                          = "gcpvpc"
	CredentialTypeKeyPair                         = "keypair"
	CredentialTypeLdap                            = "ldap"
	CredentialTypeOciVcn                          = "ocivcn"
	CredentialTypePan                             = "pan"
	CredentialTypePanInstance                     = "paninstance"
	CredentialTypeFortinet                        = "ftntfw"
)

type CredentialAkamaiProlexic struct {
	BgpAuthenticationKey string `json:"bgpAuthenticationKey"`
}

type CredentialArubaEdgeConnect struct {
	AccountKey string `json:"accountKey"`
}

type CredentialAwsVpcKey struct {
	Ec2AccessKey string `json:"ec2AccessKey"`
	Ec2SecretKey string `json:"ec2SecretKey"`
	Type         string `json:"type"`
}

type CredentialAwsVpcRole struct {
	Ec2RoleArn    string `json:"ec2RoleArn"`
	Ec2ExternalId string `json:"ec2ExternalId"`
	Type          string `json:"type"`
}

type CredentialAzureVnet struct {
	ApplicationId  string `json:"applicationId"`
	SecretKey      string `json:"secretKey"`
	SubscriptionId string `json:"subscriptionId"`
	TenantId       string `json:"tenantId"`
}

type CredentialCheckPointFwService struct {
	AdminPassword string `json:"adminPassword"`
}

type CredentialCheckPointFwServiceInstance struct {
	SicKey string `json:"sicKey"`
}

type CredentialCheckPointFwManagementServer struct {
	Password string `json:"password"`
}

type CredentialCiscoSdwan struct {
	Password string `json:"password"`
	Username string `json:"userName"`
}

type CredentialGcpVpc struct {
	AuthProvider      string `json:"auth_provider_x509_cert_url"`
	AuthUri           string `json:"auth_uri"`
	ClientEmail       string `json:"client_email"`
	ClientId          string `json:"client_id"`
	ClientX509CertUrl string `json:"client_x509_cert_url"`
	PrivateKey        string `json:"private_key"`
	PrivateKeyId      string `json:"private_key_id"`
	ProjectId         string `json:"project_id"`
	TokenUri          string `json:"token_uri"`
	Type              string `json:"type"`
}

type CredentialFortinet struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type CredentialKeyPair struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Type       string `json:"type"`
}

type CredentialLdap struct {
	BindPassword   string `json:"bindPassword"`
	TlsCertificate string `json:"tlsCertificate"`
}

type CredentialOciVcn struct {
	UserId      string `json:"userId"`
	FingerPrint string `json:"fingerPrint"`
	Key         string `json:"key"`
	TenantId    string `json:"tenantId"`
}

type CredentialPan struct {
	LicenseKey string `json:"licenseKey"`
	Password   string `json:"password"`
	Username   string `json:"userName"`
}

type CredentialPanInstance struct {
	AuthKey    string `json:"authKey"`
	AuthCode   string `json:"authCode"`
	LicenseKey string `json:"licenseKey"`
	Password   string `json:"password"`
	Username   string `json:"userName"`
}

type Credentials struct {
	Name        string      `json:"name"`
	Credentials interface{} `json:"credentials"`
}

type CredentialResponse struct {
	Id string `json:"id"`
}

type CredentialResponseDetail struct {
	Id      string `json:"credentialId"`
	Type    string `json:"credentialType"`
	Name    string `json:"name"`
	SubType string `json:"subType"`
}

// CreateCredential create new credential
func (ac *AlkiraClient) CreateCredential(name string, ctype CredentialType, credential interface{}) (string, error) {
	uri := fmt.Sprintf("%s/api/credentials/%s", ac.URI, ctype)

	// This body is not the normal JSON format...
	body, err := json.Marshal(Credentials{
		Name:        name,
		Credentials: credential,
	})

	if err != nil {
		return "", fmt.Errorf("CreateCredential: failed to marshal: %v", err)
	}

	data, err := ac.create(uri, body)

	if err != nil {
		return "", err
	}

	var result CredentialResponse
	json.Unmarshal([]byte(data), &result)

	return result.Id, nil
}

// DeleteCredential delete credential by its Id
func (ac *AlkiraClient) DeleteCredential(id string, ctype CredentialType) error {
	uri := fmt.Sprintf("%s/api/credentials/%s/%s", ac.URI, ctype, id)
	return ac.delete(uri)
}

// UpdateCredential update a given credential by its Id
func (ac *AlkiraClient) UpdateCredential(id string, name string, ctype CredentialType, credential interface{}) error {
	if ctype == CredentialTypeKeyPair || ctype == CredentialTypeArubaEdgeConnect {
		return fmt.Errorf("UpdateCredential: not supported for the credential type")
	}

	uri := fmt.Sprintf("%s/api/credentials/%s/%s", ac.URI, ctype, id)

	// This body is not the normal JSON format...
	body, err := json.Marshal(Credentials{
		Name:        name,
		Credentials: credential,
	})

	if err != nil {
		return fmt.Errorf("UpdateCredential: failed to marshal: %v", err)
	}

	return ac.update(uri, body)
}

// GetCredentials get all credentials
func (ac *AlkiraClient) GetCredentials() (string, error) {
	uri := fmt.Sprintf("%s/api/credentials/", ac.URI)

	data, err := ac.get(uri)
	return string(data), err
}

// GetCredentialById get one credential by its Id
func (ac *AlkiraClient) GetCredentialById(id string) (CredentialResponseDetail, error) {
	uri := fmt.Sprintf("%s/api/credentials/%s", ac.URI, id)

	var credential CredentialResponseDetail

	data, err := ac.get(uri)

	if err != nil {
		return credential, err
	}

	err = json.Unmarshal([]byte(data), &credential)

	if err != nil {
		return credential, fmt.Errorf("GetCredentialById: failed to unmarshal: %v", err)
	}

	return credential, nil
}

// GetCredentialByName get the credential by its name
func (ac *AlkiraClient) GetCredentialByName(name string) (CredentialResponseDetail, error) {
	var credential CredentialResponseDetail

	if len(name) == 0 {
		return credential, fmt.Errorf("Invalid credential name input")
	}

	credentials, err := ac.GetCredentials()

	if err != nil {
		return credential, err
	}

	var result []CredentialResponseDetail
	json.Unmarshal([]byte(credentials), &result)

	for _, g := range result {
		if g.Name == name {
			return g, nil
		}
	}

	return credential, fmt.Errorf("Failed to find the credential by %s", name)
}
