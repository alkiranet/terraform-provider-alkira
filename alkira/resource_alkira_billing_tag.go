package alkira

import (
	"context"
	"log"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraBillingTagResource{}
)

type alkiraBillingTagResource struct {
	client *alkira.AlkiraClient
}

type alkiraBillingTagResourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	LastUpdated types.String `tfsdk:"last_updated"`
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
			"id": schema.StringAttribute{
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
	var plan alkiraBillingTagResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)

	id, err := r.client.CreateBillingTag(plan.Name.ValueString(), plan.Description.ValueString())
	if err != nil {
		return
	}

	// Set state
	plan.Id = types.StringValue(id)
	diags := resp.State.Set(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraBillingTagResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state alkiraBillingTagResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	billingTag, err := r.client.GetBillingTagById(state.Id.String())
	if err != nil {
		return
	}

	// Set state
	state.Name = types.StringValue(billingTag.Name)
	state.Description = types.StringValue(billingTag.Description)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraBillingTagResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alkiraBillingTagResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// log.Println("[ERROR] plan: %s %s %S", plan.Id.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())

	err := r.client.UpdateBillingTag(plan.Id.ValueString(), plan.Name.ValueString(), plan.Description.ValueString())
	log.Printf("ABABA")
	if err != nil {
		return
	}

	billingTag, err := r.client.GetBillingTagById(plan.Id.ValueString())
	if err != nil {
		return
	}

	plan.Name = types.StringValue(billingTag.Name)
	plan.Description = types.StringValue(billingTag.Description)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraBillingTagResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state alkiraBillingTagResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteBillingTag(state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Billing Tag",
			"Could not delete billing tag, unexpected error: "+err.Error(),
		)
		return
	}
}
