package alkira

import (
	"context"
	"log"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraByoipPrefixResource{}
)

// PlanModifyBool is a plan modifier for types.Bool attributes.
func BoolDefaultValue(v types.Bool) planmodifier.Bool {
	return &boolDefaultValuePlanModifier{v}
}

type boolDefaultValuePlanModifier struct {
	DefaultValue types.Bool
}

var _ planmodifier.Bool = (*boolDefaultValuePlanModifier)(nil)

func (apm *boolDefaultValuePlanModifier) Description(ctx context.Context) string {
	/* ... */
	return ""
}

func (apm *boolDefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	/* ... */
	return ""
}

func (apm *boolDefaultValuePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, res *planmodifier.BoolResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		log.Printf("[ERROR] AAAAA")
		return
	}
	log.Printf("[ERROR] BBBBB")
	res.PlanValue = apm.DefaultValue
}

type alkiraByoipPrefixResource struct {
	client      *alkira.AlkiraClient
	byoipPrefix *alkira.AlkiraAPI[alkira.Byoip]
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
	r.byoipPrefix = alkira.NewByoip(r.client)
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
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Bool{
					BoolDefaultValue(types.BoolValue(false)),
				},
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

	if resp.Diagnostics.HasError() {
		return
	}

	result, _, err := r.byoipPrefix.Create(plan)
	if err != nil {
		return
	}

	err = SetByoipStateCreate(ctx, req, resp, result)
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

	result, err := r.byoipPrefix.GetById(strconv.Itoa(id))
	if err != nil {
		return
	}

	err = SetByoipStateRead(ctx, req, resp, result)
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

	_, err := r.byoipPrefix.Delete(strconv.Itoa(id))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Byoip Prefix",
			"Could not delete Byoip Prefix, unexpected error: "+err.Error(),
		)
		return
	}
}
