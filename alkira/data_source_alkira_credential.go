package alkira

import (
	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlkiraCredential() *schema.Resource {
	return &schema.Resource{
		Description: "Use this data source to get information on an existing credential.",

		Read: dataSourceAlkiraCredentialRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the credentials.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"credential_id": {
				Description: "The ID of the credentials.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
		},
	}
}

func dataSourceAlkiraCredentialRead(d *schema.ResourceData, m interface{}) error {
	client := meta.(*alkira.AlkiraClient)

	credential, err := client.GetCredentialByName(d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(credential.Id)
	d.Set("credential_id", credential.Id)

	return nil
}
