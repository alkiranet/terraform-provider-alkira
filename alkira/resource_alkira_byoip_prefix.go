package alkira

import (
	"context"
	"encoding/json"
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
			"id": schema.NumberAttribute{
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
	plan, err := GenerateByoipRequestCreate(ctx, req, resp)
	if err != nil {
		return
	}

	id, err := r.client.CreateByoip(plan)
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(id), &plan.Id)
	if err != nil {
		return
	}

	err = SetByoipStateCreate(ctx, req, resp, plan)
	if err != nil {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraByoipPrefixResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id int

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan, err := r.client.GetByoipById(strconv.Itoa(id))
	if err != nil {
		return
	}

	err = SetByoipStateRead(ctx, req, resp, plan)
	if err != nil {
		return
	}
}

// Byoip Prefix does not support update
func (r *alkiraByoipPrefixResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraByoipPrefixResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id int

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteByoip(strconv.Itoa(id))
	if err != nil {
		return
	}
}
