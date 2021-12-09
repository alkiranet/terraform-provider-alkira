package alkira

import (
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlkiraCloudVisorAccount() *schema.Resource {
	return &schema.Resource{
		Description: "Manage CloudVisor Accounts",
		Create:      resourceCloudVisorAccount,
		Read:        resourceCloudVisorAccountRead,
		Update:      resourceCloudVisorAccountUpdate,
		Delete:      resourceCloudVisorAccountDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the account.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"credential_id": {
				Description: "Credential Id to be used for the account.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cloud_provider": {
				Description:  "Cloud provider of the account, currently, `AWS` and `AZURE` are supported.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"AWS", "AZURE"}, false),
			},
			"auto_sync": {
				Description:  "The interval at which the account should be auto synced. The value could be `NONE`, `DAILY`, `WEEKLY` and `MONTHLY`.",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"NONE", "DAILY", "WEEKLY", "MONTHLY"}, false),
			},
			"native_id": {
				Description: "The native cloud provider account Id.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceCloudVisorAccount(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	request := generateCloudVisorAccountRequest(d, m)

	log.Printf("[INFO] Creating cloudvisor-account")
	id, err := client.CreateCloudProviderAccount(request)

	if err != nil {
		return err
	}

	d.SetId(id)

	return resourceCloudVisorAccountRead(d, m)
}

func resourceCloudVisorAccountRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)

	account, err := client.GetCloudProviderAccountById(d.Id())

	if err != nil {
		return err
	}

	d.Set("name", account.Name)
	d.Set("credential_id", account.CredentialId)
	d.Set("cloud_provider", account.CloudProvider)
	d.Set("auto_sync", account.AutoSync)
	d.Set("native_id", account.NativeId)

	return nil
}

func resourceCloudVisorAccountUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*alkira.AlkiraClient)
	request := generateCloudVisorAccountRequest(d, m)

	log.Printf("[INFO] Updating cloudvisor-account (%s)", d.Id())
	err := client.UpdateCloudProviderAccount(d.Id(), request)

	if err != nil {
		return err
	}

	return nil
}

func resourceCloudVisorAccountDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	log.Printf("[INFO] Deleting CloudVisorAccount (%s)", d.Id())
	err := client.DeleteCloudProviderAccount(d.Id())

	if err != nil {
		return err
	}

	log.Printf("[INFO] Deleted CloudVisorAccount (%s)", d.Id())
	d.SetId("")
	return nil
}

func generateCloudVisorAccountRequest(d *schema.ResourceData, m interface{}) alkira.CloudProviderAccount {
	return alkira.CloudProviderAccount{
		Name:          d.Get("name").(string),
		CredentialId:  d.Get("credential_id").(string),
		CloudProvider: d.Get("cloud_provider").(string),
		AutoSync:      d.Get("auto_sync").(string),
		NativeId:      d.Get("native_id").(string),
	}
}