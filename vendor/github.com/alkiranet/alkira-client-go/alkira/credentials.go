// Copyright (C) 2020-2025 Alkira Inc. All Rights Reserved.

package alkira

import (
	"encoding/json"
	"fmt"
)

type CredentialType string

const (
	CredentialTypeAkamaiProlexic           CredentialType = "akamai-prolexic"
	CredentialTypeArubaEdgeConnectInstance CredentialType = "aruba-edge-connector-instances"
	CredentialTypeAwsVpc                   CredentialType = "awsvpc"
	CredentialTypeAzureVnet                CredentialType = "azurevnet"
	CredentialTypeChkpFw                   CredentialType = "chkp-fw"
	CredentialTypeChkpFwInstance           CredentialType = "chkp-fw-instance"
	CredentialTypeChkpFwManagement         CredentialType = "chkp-fw-management-server"
	CredentialTypeCiscoFtdv                CredentialType = "cisco-ftdv-fw"
	CredentialTypeCiscoFtdvInstance        CredentialType = "cisco-ftdv-fw-instance"
	CredentialTypeCiscoSdwan               CredentialType = "ciscosdwan"
	CredentialTypeFortinet                 CredentialType = "ftntfw"
	CredentialTypeFortinetInstance         CredentialType = "ftntfw-instance"
	CredentialTypeFortinetSdwanInstance    CredentialType = "ftnt-sdwan-connector-instance"
	CredentialTypeGcpVpc                   CredentialType = "gcpvpc"
	CredentialTypeInfoblox                 CredentialType = "infoblox"
	CredentialTypeInfobloxGridMaster       CredentialType = "infoblox-grid-master"
	CredentialTypeInfobloxInstance         CredentialType = "infoblox-instance"
	CredentialTypeKeyPair                  CredentialType = "keypair"
	CredentialTypeLdap                     CredentialType = "ldap"
	CredentialTypeOciVcn                   CredentialType = "ocivcn"
	CredentialTypePan                      CredentialType = "pan"
	CredentialTypePanInstance              CredentialType = "paninstance"
	CredentialTypePanMasterKey             CredentialType = "pan-masterkey"
	CredentialTypePanRegistration          CredentialType = "pan-registration"
	CredentialTypeVmwareSdwanInstance      CredentialType = "vmware-sdwan-connector-instance"
	CredentialTypeF5Instance               CredentialType = "f5-lb-instance"
	CredentialTypeF5InstanceRegistration   CredentialType = "f5-lb-registration"
	CredentialTypeUserNamePassword         CredentialType = "username-password"
	CredentialTypeApiKey                   CredentialType = "api-key"
)

type CredentialAkamaiProlexic struct {
	BgpAuthenticationKey string `json:"bgpAuthenticationKey"`
}

type CredentialArubaEdgeConnectInstance struct {
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
	Environment    string `json:"environment,omitempty"`
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

type CredentialCiscoFtdv struct {
	Password string `json:"password"`
	Username string `json:"userName"`
}

type CredentialCiscoFtdvInstance struct {
	AdminPassword      string `json:"adminPassword"`
	FmcRegistrationKey string `json:"fmcRegistrationKey"`
	FtvdNatId          string `json:"ftdvNatId,omitempty"`
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

type CredentialFortinetInstance struct {
	LicenseType string `json:"licenseType"`
	LicenseKey  string `json:"licenseKey"`
}

type CredentialFortinetSdwanInstance struct {
	LicenseType string `json:"licenseType"`
	LicenseKey  string `json:"licenseKey"`
	Password    string `json:"password"`
	Username    string `json:"userName"`
}

type CredentialInfoblox struct {
	SharedSecret string `json:"sharedSecret"`
}

type CredentialInfobloxInstance struct {
	Password string `json:"password"`
}

type CredentialInfobloxGridMaster struct {
	Username string `json:"userName"`
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

type CredentialF5Instance struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type CredentialF5InstanceRegistration struct {
	RegistrationKey string `json:"registrationKey"`
}

type CredentialPanMasterKey struct {
	MasterKey string `json:"masterKey"`
}

type CredentialPanRegistration struct {
	RegistrationPinId    string `json:"registrationPinId"`
	RegistrationPinValue string `json:"registrationPinValue"`
}

type CredentialVmwareSdwanInstance struct {
	ActivationCode string `json:"activationCode"`
}

type Credentials struct {
	Name        string      `json:"name"`
	Credentials interface{} `json:"credentials"`
	Expires     int64       `json:"expires,omitempty"`
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

type CredentialUserNamePassword struct {
	Password string `json:"password"`
	Username string `json:"userName"`
}

type CredentialApiKey struct {
	ApiKey string `json:"apiKey"`
}

// CreateCredential create new credential
func (ac *AlkiraClient) CreateCredential(name string, ctype CredentialType, credential interface{}, expires int64) (string, error) {
	uri := fmt.Sprintf("%s/api/credentials/%s", ac.URI, ctype)

	// This body is not the normal JSON format...
	body, err := json.Marshal(Credentials{
		Name:        name,
		Credentials: credential,
		Expires:     expires,
	})

	if err != nil {
		return "", fmt.Errorf("CreateCredential: failed to marshal: %v", err)
	}

	data, _, err, _ := ac.create(uri, body, false)

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
	_, err, _ := ac.delete(uri, false)
	return err
}

// UpdateCredential update a given credential by its Id
func (ac *AlkiraClient) UpdateCredential(id string, name string, ctype CredentialType, credential interface{}, expires int64) error {
	if ctype == CredentialTypeKeyPair || ctype == CredentialTypeArubaEdgeConnectInstance {
		return fmt.Errorf("UpdateCredential: not supported for the credential type")
	}

	uri := fmt.Sprintf("%s/api/credentials/%s/%s", ac.URI, ctype, id)

	// This body is not the normal JSON format...
	body, err := json.Marshal(Credentials{
		Name:        name,
		Credentials: credential,
		Expires:     expires,
	})

	if err != nil {
		return fmt.Errorf("UpdateCredential: failed to marshal: %v", err)
	}

	_, err, _ = ac.update(uri, body, false)

	return err
}

// GetCredentials get all credentials
func (ac *AlkiraClient) GetCredentials() (string, error) {
	uri := fmt.Sprintf("%s/api/credentials/", ac.URI)

	data, _, err := ac.get(uri)
	return string(data), err
}

// GetCredentialById get one credential by its Id
func (ac *AlkiraClient) GetCredentialById(id string) (CredentialResponseDetail, error) {
	uri := fmt.Sprintf("%s/api/credentials/%s", ac.URI, id)

	var credential CredentialResponseDetail

	data, _, err := ac.get(uri)

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
