package alkira

import (
	"context"
	"fmt"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraCredentialPrismaSDWAN() *schema.Resource {
	return &schema.Resource{
		Description:   "Provides Prisma SD-WAN connector instance credential resource.",
		CreateContext: resourceCredentialPrismaSDWANCreate,
		ReadContext:   resourceCredentialPrismaSDWANRead,
		UpdateContext: resourceCredentialPrismaSDWANUpdate,
		DeleteContext: resourceCredentialPrismaSDWANDelete,
		Importer: &schema.ResourceImporter{
			StateContext: importWithReadValidation(resourceCredentialPrismaSDWANRead),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the credential.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ion_token": {
				Description: "The ION token for Prisma SD-WAN device registration.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"ion_secret": {
				Description: "The ION secret for Prisma SD-WAN device registration.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

func resourceCredentialPrismaSDWANCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialPrismaSDWANInstance{
		IonToken:  d.Get("ion_token").(string),
		IonSecret: d.Get("ion_secret").(string),
	}

	id, err := client.CreateCredential(d.Get("name").(string), alkira.CredentialTypePrismaSDWANInstance, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(id)
	return resourceCredentialPrismaSDWANRead(ctx, d, meta)
}

func resourceCredentialPrismaSDWANRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func resourceCredentialPrismaSDWANUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	c := alkira.CredentialPrismaSDWANInstance{
		IonToken:  d.Get("ion_token").(string),
		IonSecret: d.Get("ion_secret").(string),
	}

	err := client.UpdateCredential(d.Id(), d.Get("name").(string), alkira.CredentialTypePrismaSDWANInstance, c, 0)

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceCredentialPrismaSDWANRead(ctx, d, meta)
}

func resourceCredentialPrismaSDWANDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*alkira.AlkiraClient)

	err := client.DeleteCredential(d.Id(), alkira.CredentialTypePrismaSDWANInstance)

	if err != nil {
		name, _ := d.GetOk("name")
		if nameStr, ok := name.(string); ok && nameStr != "" {
			return diag.FromErr(fmt.Errorf("%w alkira_credential_prisma_sdwan (name=%q id=%s)", err, nameStr, d.Id()))
		}
		return diag.FromErr(fmt.Errorf("%w alkira_credential_prisma_sdwan (id=%s)", err, d.Id()))
	}

	d.SetId("")
	return nil
}
