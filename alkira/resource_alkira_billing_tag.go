package alkira

import (
	"context"
	"strconv"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraBillingTagResource{}
)

type alkiraBillingTagResource struct {
	client *alkira.AlkiraClient
}

func NewalkiraBillingTagResource() resource.Resource {
	return &alkiraBillingTagResource{}
}

// Configure adds the provider configured client to the resource.
func (r *alkiraBillingTagResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*alkira.AlkiraClient)
}

// Metadata returns the resource type name.
func (r *alkiraBillingTagResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_tag"
}

// Schema defines the schema for the resource.
func (r *alkiraBillingTagResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID billing tag.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the billing tag.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Billing tag description.",
				Optional:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *alkiraBillingTagResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alkira.BillingTag

	// resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := r.client.CreateBillingTag(plan.Name, plan.Description)
	if err != nil {
		return
	}

	plan.Id, _ = strconv.Atoi(id)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.Id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), plan.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), plan.Description)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraBillingTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan alkira.BillingTag

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &plan.Id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan, err := r.client.GetBillingTagById(strconv.Itoa(plan.Id))
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), plan.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), plan.Description)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraBillingTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alkira.BillingTag

	// Grab the ID from the state.
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &plan.Id)...)

	// Grab the name and description from the plan.
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.UpdateBillingTag(strconv.Itoa(plan.Id), plan.Name, plan.Description)
	if err != nil {
		return
	}

	plan, err = r.client.GetBillingTagById(strconv.Itoa(plan.Id))
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), plan.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), plan.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("last_updated"), time.Now().Format(time.RFC3339))...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraBillingTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var plan alkira.BillingTag

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &plan.Id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteBillingTag(strconv.Itoa(plan.Id))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Billing Tag",
			"Could not delete billing tag, unexpected error: "+err.Error(),
		)
		return
	}
}
