package alkira

import (
	"context"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraByoipPrefixResource{}
)

type alkiraByoipPrefixResource struct {
	client *alkira.AlkiraClient
}

func NewalkiraByoipPrefixResource() resource.Resource {
	return &alkiraByoipPrefixResource{}
}

// Configure adds the provider configured client to the resource.
func (r *alkiraByoipPrefixResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*alkira.AlkiraClient)
}

// Metadata returns the resource type name.
func (r *alkiraByoipPrefixResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_byoip_prefix"
}

// Schema defines the schema for the resource.
func (r *alkiraByoipPrefixResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Description: "The ID of Byoip Prefix.",
				Computed:    true,
			},
			"prefix": schema.StringAttribute{
				Description: "Prefix for BYOIP.",
				Required:    true,
			},
			"cxp": schema.StringAttribute{
				Description: "CXP region.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Description for the list.",
				Optional:    true,
			},
			"do_not_advertise": schema.BoolAttribute{
				Description: "Do not advertise.",
				Optional:    true,
				// set default value to false
			},
			"message": schema.StringAttribute{
				Description: "Message from AWS BYOIP.",
				Required:    true,
			},
			"signature": schema.StringAttribute{
				Description: "Signature from AWS BYOIP.",
				Required:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "Public key from AWS BYOIP.",
				Required:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *alkiraByoipPrefixResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	plan, err := GenerateByoipRequest(ctx, req, resp)
	if err != nil {
		return
	}

	id, err := r.client.CreateByoip(plan)
	if err != nil {
		return
	}

	plan.Id, _ = strconv.Atoi(id)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), plan.Id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("prefix"), plan.Prefix)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cxp"), plan.Cxp)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), plan.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("do_not_advertise"), plan.DoNotAdvertise)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("message"), plan.ExtraAttributes.Message)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("signature"), plan.ExtraAttributes.Signature)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_key"), plan.ExtraAttributes.PublicKey)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraByoipPrefixResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraByoipPrefixResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraByoipPrefixResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}
