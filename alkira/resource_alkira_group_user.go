package alkira

import (
	"context"
	"fmt"
	"log"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlkiraGroupUser() *schema.Resource {
	return &schema.Resource{
		Description:   "Manage user groups\n\n",
		CreateContext: resourceGroupUser,
		ReadContext:   resourceGroupUserRead,
		UpdateContext: resourceGroupUserUpdate,
		DeleteContext: resourceGroupUserDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Description: "The name of the user group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the user group.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}

func resourceGroupUser(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewUserGroup(m.(*alkira.AlkiraClient))

	group := &alkira.UserGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	resource, _, err, _ := api.Create(group)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(string(resource.Id))

	return resourceGroupUserRead(ctx, d, m)
}

func resourceGroupUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewUserGroup(m.(*alkira.AlkiraClient))

	group, _, err := api.GetById(d.Id())

	if err != nil {
		return diag.Diagnostics{{
			Severity: diag.Warning,
			Summary:  "FAILED TO GET RESOURCE",
			Detail:   fmt.Sprintf("%s", err),
		}}
	}

	d.Set("name", group.Name)
	d.Set("description", group.Description)

	return nil
}

func resourceGroupUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewUserGroup(m.(*alkira.AlkiraClient))

	group := &alkira.UserGroup{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	log.Printf("[INFO] Updating User Group (%s)", d.Id())
	_, err, _ := api.Update(d.Id(), group)

	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGroupUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	api := alkira.NewUserGroup(m.(*alkira.AlkiraClient))

	_, err, _ := api.Delete(d.Id())

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
