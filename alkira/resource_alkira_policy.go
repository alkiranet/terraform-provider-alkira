package alkira

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraPolicy{}
)

type alkiraPolicyModel struct {
	ProvisionState types.String `tfsdk:"provision_state"`
	Id             types.Int64  `tfsdk:"id"`
	Description    types.String `tfsdk:"description"`
	Enabled        types.Bool   `tfsdk:"enabled"`
	FromGroups     types.List   `tfsdk:"from_groups"`
	Name           types.String `tfsdk:"name"`
	RuleListId     types.Int64  `tfsdk:"rule_list_id"`
	SegmentIds     types.List   `tfsdk:"segment_ids"`
	ToGroups       types.List   `tfsdk:"to_groups"`
	LastUpdated    types.String `tfsdk:"last_updated"`
}

type alkiraPolicy struct {
	client *alkira.AlkiraClient
	policy *alkira.AlkiraAPI[alkira.TrafficPolicy]
}

func NewalkiraPolicy() resource.Resource {
	return &alkiraPolicy{}
}

// Configure adds the provider configured client to the resource.
func (r *alkiraPolicy) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*alkira.AlkiraClient)
	r.policy = alkira.NewTrafficPolicy(r.client)
}

// Metadata returns the resource type name.
func (r *alkiraPolicy) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy"
}

// Schema defines the schema for the resource.
func (r *alkiraPolicy) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"provision_state": schema.StringAttribute{
				Description: "Provisioning state of the Segment.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "The ID Segment.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the Segment.",
				Required:    true,
			},
			"enabled": schema.BoolAttribute{
				Description: "Whether the policy is enabled.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Segment description.",
				Optional:    true,
			},
			"from_groups": schema.ListAttribute{
				Description: "IDs of the groups that will define source in the policy scope.",
				Required:    true,
				ElementType: types.Int64Type,
			},
			"rule_list_id": schema.Int64Attribute{
				Description: "The `rulelist` that will be used by the policy.",
				Required:    true,
			},
			"segment_ids": schema.ListAttribute{
				Description: "IDs of the segments that will define the policy scope.",
				Required:    true,
				ElementType: types.Int64Type,
			},
			"to_groups": schema.ListAttribute{
				Description: "IDs of groups that will define destination in the policy scope.",
				Required:    true,
				ElementType: types.Int64Type,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *alkiraPolicy) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alkiraPolicyModel
	var segIds []int
	var fromGroups []int
	var toGroups []int

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	resp.Diagnostics.Append(plan.SegmentIds.ElementsAs(ctx, &segIds, true)...)
	resp.Diagnostics.Append(plan.FromGroups.ElementsAs(ctx, &fromGroups, true)...)
	resp.Diagnostics.Append(plan.ToGroups.ElementsAs(ctx, &toGroups, true)...)

	log.Printf("[DEBUG] Creating Policy: %v", plan)

	if resp.Diagnostics.HasError() {
		return
	}

	policy := &alkira.TrafficPolicy{
		Name:        plan.Name.ValueString(),
		Enabled:     plan.Enabled.ValueBool(),
		Description: plan.Description.ValueString(),
		FromGroups:  fromGroups,
		RuleListId:  int(plan.RuleListId.ValueInt64()),
		SegmentIds:  segIds,
		ToGroups:    toGroups,
	}

	result, state, err := r.policy.Create(policy)
	if err != nil {
		return
	}

	id, err := result.Id.Int64()
	if err != nil {
		return
	}

	plan.Id = types.Int64Value(id)
	plan.ProvisionState = types.StringValue(state)
	plan.Description = types.StringValue(result.Description)
	plan.Enabled = types.BoolValue(result.Enabled)
	plan.RuleListId = types.Int64Value(int64(result.RuleListId))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("from_groups"), result.FromGroups)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("segment_ids"), result.SegmentIds)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("to_groups"), result.ToGroups)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraPolicy) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan alkiraPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.policy.GetById(strconv.Itoa(int(plan.Id.ValueInt64())))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Policy Resource",
			"Could not read policy resource, unexpected error: "+err.Error(),
		)
		return
	}

	id, err := result.Id.Int64()
	if err != nil {
		return
	}

	plan.Id = types.Int64Value(id)
	plan.Description = types.StringValue(result.Description)
	plan.Enabled = types.BoolValue(result.Enabled)
	plan.RuleListId = types.Int64Value(int64(result.RuleListId))
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("from_groups"), result.FromGroups)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("segment_ids"), result.SegmentIds)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("to_groups"), result.ToGroups)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraPolicy) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// state, err := r.policy.Update(strconv.Itoa(id), &plan)
	// if err != nil {
	// 	return
	// }

	// result, err := r.policy.GetById(strconv.Itoa(id))
	// if err != nil {
	// 	return
	// }

	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraPolicy) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id int

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.policy.Delete(strconv.Itoa(id))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Segment",
			"Could not delete group, unexpected error: "+err.Error(),
		)
		return
	}
}
