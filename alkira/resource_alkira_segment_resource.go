package alkira

import (
	"context"
	"strconv"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraSegmentResource{}
)

type alkiraSegmentResourceModel struct {
	ProvisionState  types.String `tfsdk:"provision_state"`
	Id              types.Int64  `tfsdk:"id"`
	Name            types.String `tfsdk:"name"`
	SegmentId       types.String `tfsdk:"segment_id"`
	ImplicitGroupId types.Int64  `tfsdk:"implicit_group_id"`
	LastUpdated     types.String `tfsdk:"last_updated"`
	GroupPrefixes   types.Set    `tfsdk:"group_prefix"`
}

type segmentResourceGroupPrefixModel struct {
	GroupId      types.Int64 `tfsdk:"group_id"`
	PrefixListId types.Int64 `tfsdk:"prefix_list_id"`
}

type alkiraSegmentResource struct {
	client  *alkira.AlkiraClient
	segment *alkira.AlkiraAPI[alkira.SegmentResource]
}

func NewalkiraSegmentResource() resource.Resource {
	return &alkiraSegmentResource{}
}

// Configure adds the provider configured client to the resource.
func (r *alkiraSegmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*alkira.AlkiraClient)
	r.segment = alkira.NewSegmentResource(r.client)
}

// Metadata returns the resource type name.
func (r *alkiraSegmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment_resource"
}

// Schema defines the schema for the resource.
func (r *alkiraSegmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"provision_state": schema.StringAttribute{
				Description: "Provisioning state of the group tag.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "The ID group.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the group tag.",
				Required:    true,
			},
			"segment_id": schema.StringAttribute{
				Description: "The ID of the segment.",
				Optional:    true,
			},
			"implicit_group_id": schema.Int64Attribute{
				Description: "The ID of automatically created implicit group.",
				Computed:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
		Blocks: map[string]schema.Block{
			"group_prefix": schema.SetNestedBlock{
				Validators: []validator.Set{
					setvalidator.IsRequired(),
				},
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group_id": schema.Int64Attribute{
							Description: "The connector group ID associated with segment resource.",
							Optional:    true,
						},
						"prefix_list_id": schema.Int64Attribute{
							Description: "The prefix list ID associated with segment resource.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *alkiraSegmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alkiraSegmentResourceModel
	var prefixes []segmentResourceGroupPrefixModel
	var groupPrefixes []alkira.SegmentResourceGroupPrefix

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	m := alkira.NewSegment(r.client)
	segment, err := m.GetById(plan.SegmentId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Segment Resource",
			"Could not get segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(plan.GroupPrefixes.ElementsAs(ctx, &prefixes, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, prefix := range prefixes {
		groupPrefixes = append(groupPrefixes, alkira.SegmentResourceGroupPrefix{
			GroupId:      int(prefix.GroupId.ValueInt64()),
			PrefixListId: int(prefix.PrefixListId.ValueInt64()),
		})
	}

	resource := alkira.SegmentResource{
		Name:          plan.Name.ValueString(),
		Segment:       segment.Name,
		GroupPrefixes: groupPrefixes,
	}

	result, provState, err := r.segment.Create(&resource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Segment Resource",
			"Could not create segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	id, _ := result.Id.Int64()
	plan.Id = types.Int64Value(id)
	plan.ProvisionState = types.StringValue(provState)
	plan.SegmentId = types.StringValue(segment.Id.String())
	plan.ImplicitGroupId = types.Int64Value(int64(result.GroupId))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraSegmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan alkiraSegmentResourceModel
	var prefixes []segmentResourceGroupPrefixModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	result, err := r.segment.GetById(strconv.Itoa(int(plan.Id.ValueInt64())))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Segment Resource",
			"Could not read segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	id, _ := result.Id.Int64()
	plan.Id = types.Int64Value(id)
	plan.Name = types.StringValue(result.Name)
	plan.ImplicitGroupId = types.Int64Value(int64(result.GroupId))

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("segment_id"), &plan.SegmentId)...)
	m := alkira.NewSegment(r.client)
	segment, err := m.GetById(plan.SegmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Segment Resource",
			"Could not read segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	plan.SegmentId = types.StringValue(segment.Id.String())

	for _, prefix := range result.GroupPrefixes {
		prefixes = append(prefixes, segmentResourceGroupPrefixModel{
			GroupId:      types.Int64Value(int64(prefix.GroupId)),
			PrefixListId: types.Int64Value(int64(prefix.PrefixListId)),
		})
	}

	resp.Diagnostics.Append(req.State.Set(ctx, plan)...)
	resp.Diagnostics.Append(req.State.SetAttribute(ctx, path.Root("group_prefix"), prefixes)...)

	if resp.Diagnostics.HasError() {
		return
	}

}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraSegmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alkiraSegmentResourceModel
	var prefixes []segmentResourceGroupPrefixModel
	var groupPrefixes []alkira.SegmentResourceGroupPrefix

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &plan.Id)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("provision_state"), &plan.ProvisionState)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("implicit_group_id"), &plan.ImplicitGroupId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	m := alkira.NewSegment(r.client)
	segment, err := m.GetById(plan.SegmentId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Segment Resource",
			"Could not get segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(plan.GroupPrefixes.ElementsAs(ctx, &prefixes, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, prefix := range prefixes {
		groupPrefixes = append(groupPrefixes, alkira.SegmentResourceGroupPrefix{
			GroupId:      int(prefix.GroupId.ValueInt64()),
			PrefixListId: int(prefix.PrefixListId.ValueInt64()),
		})
	}

	resource := alkira.SegmentResource{
		Name:          plan.Name.ValueString(),
		Segment:       segment.Name,
		GroupPrefixes: groupPrefixes,
	}

	_, err = r.segment.Update(plan.Id.String(), &resource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Segment Resource",
			"Could not update segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_prefix"), prefixes)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("last_updated"), time.Now().Format(time.RFC3339))...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraSegmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id int

	req.State.GetAttribute(ctx, path.Root("id"), &id)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.segment.Delete(strconv.Itoa(id))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Segment Resource",
			"Could not delete segment resource, unexpected error: "+err.Error(),
		)
		return
	}
}
