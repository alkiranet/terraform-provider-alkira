// Copyright (C) 2022 Alkira Inc. All Rights Reserved.
package alkira

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
)

func createCheckpointCredential(name, password string, client *alkira.AlkiraClient) (cid string, err error) {
	c := &alkira.CredentialCheckPointFwService{AdminPassword: password}

	log.Printf("[INFO] Creating Credential (Checkpoint)")
	return client.CreateCredential(name, alkira.CredentialTypeChkpFw, c)
}

func createCheckpointCredentialInstances(sicKeys []string, client *alkira.AlkiraClient) error {
	var err error

	for i, v := range sicKeys {
		cInstance := &alkira.CredentialCheckPointFwServiceInstance{SicKey: v}
		log.Printf("[INFO] Creating Credential (CheckpointInstance)")

		name := "checkpoint-instance-credential-" + strconv.Itoa(i)
		_, credentialErr := client.CreateCredential(name, alkira.CredentialTypeChkpFwInstance, cInstance)

		if credentialErr != nil {
			err = fmt.Errorf("%w:", err)
		}
	}

	return err
}

func createCheckpointCredentialManagementServer(name, password string, client *alkira.AlkiraClient) error {
	c := &alkira.CredentialCheckPointFwManagementServer{Password: password}

	log.Printf("[INFO] Creating Credential (Checkpoint Management Server)")
	_, err := client.CreateCredential(name, alkira.CredentialTypeChkpFwManagement, c)

	return err
}

func updateCheckpointCredential(cid, name, password string, client *alkira.AlkiraClient) error {
	c := &alkira.CredentialCheckPointFwService{AdminPassword: password}

	log.Printf("[INFO] Updating Credential (Checkpoint)")

	return client.UpdateCredential(cid, name, alkira.CredentialTypeChkpFw, c)
}

func updateCheckpointCredentialManagementServerByName(name, password string, client *alkira.AlkiraClient) error {
	respDetail, err := client.GetCredentialByName(name)
	if err != nil {
		return err
	}

	c := &alkira.CredentialCheckPointFwManagementServer{Password: password}
	log.Printf("[INFO] Updating Credential (Checkpoint Management Server)")

	return client.UpdateCredential(respDetail.Id, name, alkira.CredentialTypeChkpFwManagement, c)
}

func deleteCheckpointCredential(cid string, client *alkira.AlkiraClient) error {

	log.Printf("[INFO] Deleting Credential (Checkpoint %s)\n", cid)
	err := client.DeleteCredential(cid, alkira.CredentialTypeChkpFw)

	if err != nil {
		log.Printf("[INFO] Credential (Checkpoint %s) was already deleted\n", cid)
	}

	return nil
}

func deleteCheckpointCredentialInstances(client *alkira.AlkiraClient) error {
	var err error
	var credentials []alkira.CredentialResponseDetail

	js, err := client.GetCredentials()
	if err != nil {
		return err
	}

	err = json.Unmarshal([]byte(js), &credentials)
	if err != nil {
		log.Printf("[INFO] Failed Unmarshalling Credential (CheckpointInstance)")
		return err
	}

	for _, v := range credentials {
		if v.Type == "CHKPFW_INSTANCE" {
			delErr := client.DeleteCredential(v.Id, alkira.CredentialTypeChkpFwInstance)
			if delErr != nil {
				err = fmt.Errorf("%w: ", err)
			}
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
