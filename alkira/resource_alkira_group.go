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
	_ resource.ResourceWithConfigure = &alkiraGroupResource{}
)

type alkiraGroupResource struct {
	client *alkira.AlkiraClient
	group  *alkira.AlkiraAPI[alkira.Group]
}

func NewalkiraGroupResource() resource.Resource {
	return &alkiraGroupResource{}
}

// Configure adds the provider configured client to the resource.
func (r *alkiraGroupResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*alkira.AlkiraClient)
	r.group = alkira.NewGroup(r.client)
}

// Metadata returns the resource type name.
func (r *alkiraGroupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

// Schema defines the schema for the resource.
func (r *alkiraGroupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"state": schema.StringAttribute{
				Description: "Provisioning state of the billing tag.",
				Computed:    true,
			},
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
func (r *alkiraGroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alkira.Group

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, state, err := r.group.Create(&plan)
	if err != nil {
		return
	}

	id, err := result.Id.Int64()
	if err != nil {
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("state"), state)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), result.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), result.Description)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraGroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id int

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.group.GetById(strconv.Itoa(id))
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), result.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), result.Description)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraGroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alkira.Group
	var id int

	// Grab the ID from the state.
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)

	// Grab the name and description from the plan.
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.group.Update(strconv.Itoa(id), &plan)
	if err != nil {
		return
	}

	result, err := r.group.GetById(strconv.Itoa(id))
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("state"), state)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), result.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), result.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("last_updated"), time.Now().Format(time.RFC3339))...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraGroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id int

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.group.Delete(strconv.Itoa(id))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Billing Tag",
			"Could not delete billing tag, unexpected error: "+err.Error(),
		)
		return
	}
}
