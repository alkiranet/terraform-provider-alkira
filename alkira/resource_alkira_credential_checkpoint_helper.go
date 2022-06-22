// Copyright (C) 2022 Alkira Inc. All Rights Reserved.
package alkira

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/alkiranet/alkira-client-go/alkira"
)

func createCheckpointCredential(name, password string, client *alkira.AlkiraClient) (cid string, err error) {
	c := &alkira.CredentialCheckPointFwService{AdminPassword: password}

	log.Printf("[INFO] Creating Credential (Checkpoint)")
	return client.CreateCredential(name, alkira.CredentialTypeChkpFw, c, 0)
}

func createCheckpointCredentialInstances(sicKeys []string, client *alkira.AlkiraClient) error {
	var err error

	for i, v := range sicKeys {
		cInstance := &alkira.CredentialCheckPointFwServiceInstance{SicKey: v}
		log.Printf("[INFO] Creating Credential (CheckpointInstance)")

		name := "checkpoint-instance-credential-" + strconv.Itoa(i)
		_, credentialErr := client.CreateCredential(name, alkira.CredentialTypeChkpFwInstance, cInstance, 0)

		if credentialErr != nil {
			err = fmt.Errorf("%w:", err)
		}
	}

	return err
}

func createCheckpointCredentialManagementServer(name, password string, client *alkira.AlkiraClient) error {
	c := &alkira.CredentialCheckPointFwManagementServer{Password: password}

	log.Printf("[INFO] Creating Credential (Checkpoint Management Server)")
	_, err := client.CreateCredential(name, alkira.CredentialTypeChkpFwManagement, c, 0)

	return err
}

func updateCheckpointCredential(cid, name, password string, client *alkira.AlkiraClient) error {
	c := &alkira.CredentialCheckPointFwService{AdminPassword: password}

	log.Printf("[INFO] Updating Credential (Checkpoint)")

	return client.UpdateCredential(cid, name, alkira.CredentialTypeChkpFw, c, 0)
}

func updateCheckpointCredentialManagementServerByName(name, password string, client *alkira.AlkiraClient) error {
	respDetail, err := client.GetCredentialByName(name)
	if err != nil {
		return err
	}

	c := &alkira.CredentialCheckPointFwManagementServer{Password: password}
	log.Printf("[INFO] Updating Credential (Checkpoint Management Server)")

	return client.UpdateCredential(respDetail.Id, name, alkira.CredentialTypeChkpFwManagement, c, 0)
}

func deleteCheckpointCredential(cid string, client *alkira.AlkiraClient) error {

	log.Printf("[INFO] Deleting Credential (Checkpoint %s)", cid)
	err := client.DeleteCredential(cid, alkira.CredentialTypeChkpFw)

	if err != nil {
		log.Printf("[INFO] Credential (Checkpoint %s) was already deleted", cid)
	}

	return nil
}

func deleteCheckpointCredentialInstances(client *alkira.AlkiraClient) error {
	var err error

	credentials, err := getAllCheckpointCredentialInstances(client)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleting All Credential (CheckpointInstances)")
	for _, v := range credentials {
		delErr := client.DeleteCredential(v.Id, alkira.CredentialTypeChkpFwInstance)
		if delErr != nil {
			err = fmt.Errorf("%w: ", err)
		}
	}

	return err
}

func deleteCheckpointCredentialManagementServerByName(name string, client *alkira.AlkiraClient) error {
	managementServerCredential, err := client.GetCredentialByName(name)
	if err != nil {
		return err
	}

	return client.DeleteCredential(managementServerCredential.Id, alkira.CredentialTypeChkpFwManagement)
}

func getAllCheckpointCredentials(client *alkira.AlkiraClient) ([]alkira.CredentialResponseDetail, error) {
	credentials, err := getAllCredentialsAsCredentialResponseDetails(client)
	if err != nil {
		return nil, err
	}

	var checkpointCredentials []alkira.CredentialResponseDetail
	for _, v := range credentials {
		if strings.Contains(v.Type, "CHKPFW") {
			checkpointCredentials = append(checkpointCredentials, v)
		}
	}

	return credentials, nil
}

func getAllCheckpointCredentialInstances(client *alkira.AlkiraClient) ([]alkira.CredentialResponseDetail, error) {
	details, err := getAllCheckpointCredentials(client)
	if err != nil {
		return nil, err
	}

	return parseAllCheckpointCredentialInstances(details), nil
}

func parseAllCheckpointCredentialInstances(credentials []alkira.CredentialResponseDetail) []alkira.CredentialResponseDetail {
	var checkpointCredentials []alkira.CredentialResponseDetail

	for _, v := range credentials {
		if v.Type == "CHKPFW_INSTANCE" {
			checkpointCredentials = append(checkpointCredentials, v)
		}
	}

	return checkpointCredentials
}

//if strings.Contains(v.Type, "CHKPFW") == "CHKPFW_INSTANCE" {
func parseCheckpointCredentialManagementServer(credentials []alkira.CredentialResponseDetail) *alkira.CredentialResponseDetail {
	for _, v := range credentials {
		if v.Type == "CHKPFW_MANAGEMENT_SERVER" {
			return &v
		}
	}

	return nil
}

func fromCheckpointCredentialRespDetailsToCheckpointInstance(credentials []alkira.CredentialResponseDetail) []alkira.CheckpointInstance {
	var checkpointInstances []alkira.CheckpointInstance
	for _, v := range credentials {
		checkpointInstances = append(checkpointInstances, alkira.CheckpointInstance{
			CredentialId: v.Id,
			Name:         v.Name,
		})
	}

	return checkpointInstances
}
