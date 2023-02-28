package alkira

import (
	"context"
	"strconv"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraSegmentResourceShare{}
)

type alkiraSegmentResourceShareModel struct {
	ProvisionState      types.String `tfsdk:"provision_state"`
	Id                  types.Int64  `tfsdk:"id"`
	Name                types.String `tfsdk:"name"`
	ServiceList         types.List   `tfsdk:"service_ids"`
	DesignatedSegmentId types.String `tfsdk:"designated_segment_id"`
	EndAResources       types.List   `tfsdk:"end_a_segment_resource_ids"`
	EndBResources       types.List   `tfsdk:"end_b_segment_resource_ids"`
	EndARouteLimit      types.Int64  `tfsdk:"end_a_route_limit"`
	EndBRouteLimit      types.Int64  `tfsdk:"end_b_route_limit"`
	Direction           types.String `tfsdk:"traffic_direction"`
	LastUpdated         types.String `tfsdk:"last_updated"`
}

type alkiraSegmentResourceShare struct {
	client  *alkira.AlkiraClient
	segment *alkira.AlkiraAPI[alkira.SegmentResourceShare]
}

func NewalkiraSegmentResourceShare() resource.Resource {
	return &alkiraSegmentResourceShare{}
}

// Configure adds the provider configured client to the resource.
func (r *alkiraSegmentResourceShare) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	_ *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*alkira.AlkiraClient)
	r.segment = alkira.NewSegmentResourceShare(r.client)
}

// Metadata returns the resource type name.
func (r *alkiraSegmentResourceShare) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_segment_resource_share"
}

// func (r *alkiraSegmentResourceShare) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
// 	var data alkiraSegmentResourceShare

// 	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

// 	if resp.Diagnostics.HasError() {
// 		return
// 	}

// 	// If attribute_one is not configured, return without warning.
// 	if data..IsNull() || data.AttributeOne.IsUnknown() {
// 		return
// 	}

// 	// If attribute_two is not null, return without warning.
// 	if !data.AttributeTwo.IsNull() {
// 		return
// 	}

// 	resp.Diagnostics.AddAttributeWarning(
// 		path.Root("attribute_two"),
// 		"Missing Attribute Configuration",
// 		"Expected attribute_two to be configured with attribute_one. "+
// 			"The resource may return unexpected results.",
// 	)
// }

// Schema defines the schema for the resource.
func (r *alkiraSegmentResourceShare) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"service_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "The list of the service IDs.",
				Required:    true,
			},
			"designated_segment_id": schema.StringAttribute{
				Description: "The designmated segment ID.",
				Required:    true,
			},

			"end_a_segment_resource_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "The End-A segment resource IDs. All " +
					"segment resources must be on the same segment.",
				Required: true,
			},
			"end_b_segment_resource_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "The End-B segment resource IDs. All " +
					"segment resources must be on the same segment.",
				Required: true,
			},
			"end_a_route_limit": schema.Int64Attribute{
				Description: "The End-A route limit. The default " +
					"value is 100.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(100)),
				},
			},
			"end_b_route_limit": schema.Int64Attribute{
				Description: "The End-B route limit. The default " +
					"value is 100.",
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(100)),
				},
			},
			"traffic_direction": schema.StringAttribute{
				Description: "Specify the direction in which traffic " +
					"is orignated at both Resource End-A and Resource " +
					"End-B. The default value is BIDIRECTIONAL.",
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					StringDefaultValue(types.StringValue("BIDIRECTIONAL")),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("BIDIRECTIONAL", "UNIDIRECTIONAL"),
				},
			},

			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *alkiraSegmentResourceShare) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alkiraSegmentResourceShareModel
	var serviceList []int
	var endAResources []int
	var endBResources []int

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(plan.ServiceList.ElementsAs(ctx, &serviceList, true)...)
	resp.Diagnostics.Append(plan.EndAResources.ElementsAs(ctx, &endAResources, true)...)
	resp.Diagnostics.Append(plan.EndBResources.ElementsAs(ctx, &endBResources, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	m := alkira.NewSegment(r.client)
	segment, err := m.GetById(plan.DesignatedSegmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Segment Resource",
			"Could not get segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	resource := &alkira.SegmentResourceShare{
		Name:              plan.Name.ValueString(),
		ServiceList:       serviceList,
		DesignatedSegment: segment.Name,
		EndAResources:     endAResources,
		EndBResources:     endBResources,
		EndARouteLimit:    int(plan.EndARouteLimit.ValueInt64()),
		EndBRouteLimit:    int(plan.EndBRouteLimit.ValueInt64()),
		Direction:         plan.Direction.ValueString(),
	}

	segmentResourceShare, provState, err := r.segment.Create(resource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Segment Resource",
			"Could not create segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	id, err := segmentResourceShare.Id.Int64()
	if err != nil {
		return
	}

	plan.ProvisionState = types.StringValue(provState)
	plan.Id = types.Int64Value(id)
	plan.Name = types.StringValue(segmentResourceShare.Name)
	plan.EndARouteLimit = types.Int64Value(int64(segmentResourceShare.EndARouteLimit))
	plan.EndBRouteLimit = types.Int64Value(int64(segmentResourceShare.EndBRouteLimit))
	plan.Direction = types.StringValue(segmentResourceShare.Direction)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.State.Set(ctx, &plan)
	resp.State.SetAttribute(ctx, path.Root("service_ids"), segmentResourceShare.ServiceList)
	resp.State.SetAttribute(ctx, path.Root("end_a_segment_resource_ids"), segmentResourceShare.EndAResources)
	resp.State.SetAttribute(ctx, path.Root("end_b_segment_resource_ids"), segmentResourceShare.EndBResources)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraSegmentResourceShare) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan alkiraSegmentResourceShareModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.segment.GetById(strconv.Itoa(int(plan.Id.ValueInt64())))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Segment Resource",
			"Could not read segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	// m := alkira.NewSegment(r.client)
	// segment, provState, err := m.GetByName(plan.DesignatedSegmentId.ValueString())
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error Reading Segment Resource",
	// 		"Could not read segment resource, unexpected error: "+err.Error(),
	// 	)
	// 	return
	// }

	id, err := result.Id.Int64()
	if err != nil {
		return
	}

	// if r.client.Provision && provState != "" {
	// 	plan.ProvisionState = types.StringValue(provState)
	// }

	plan.Id = types.Int64Value(id)
	plan.Name = types.StringValue(result.Name)
	// plan.DesignatedSegmentId = types.StringValue(segment.Name)
	plan.EndARouteLimit = types.Int64Value(int64(result.EndARouteLimit))
	plan.EndBRouteLimit = types.Int64Value(int64(result.EndBRouteLimit))
	plan.Direction = types.StringValue(result.Direction)

	resp.Diagnostics.Append(req.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.SetAttribute(ctx, path.Root("service_ids"), result.ServiceList)...)
	resp.Diagnostics.Append(req.State.SetAttribute(ctx, path.Root("end_a_segment_resource_ids"), result.EndAResources)...)
	resp.Diagnostics.Append(req.State.SetAttribute(ctx, path.Root("end_b_segment_resource_ids"), result.EndBResources)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraSegmentResourceShare) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alkiraSegmentResourceShareModel
	// var id int64
	var serviceList []int
	var endAResources []int
	var endBResources []int

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &plan.Id)...)
	resp.Diagnostics.Append(plan.ServiceList.ElementsAs(ctx, &serviceList, true)...)
	resp.Diagnostics.Append(plan.EndAResources.ElementsAs(ctx, &endAResources, true)...)
	resp.Diagnostics.Append(plan.EndBResources.ElementsAs(ctx, &endBResources, true)...)
	if resp.Diagnostics.HasError() {
		return
	}

	m := alkira.NewSegment(r.client)
	segment, err := m.GetById(plan.DesignatedSegmentId.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Segment Resource",
			"Could not read segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	resource := &alkira.SegmentResourceShare{
		Name:              plan.Name.ValueString(),
		ServiceList:       serviceList,
		DesignatedSegment: segment.Name,
		EndAResources:     endAResources,
		EndBResources:     endBResources,
		EndARouteLimit:    int(plan.EndARouteLimit.ValueInt64()),
		EndBRouteLimit:    int(plan.EndBRouteLimit.ValueInt64()),
		Direction:         plan.Direction.ValueString(),
	}

	provState, err := r.segment.Update(plan.Id.String(), resource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Segment Resource",
			"Could not update segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))
	plan.ProvisionState = types.StringValue(provState)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("service_ids"), resource.ServiceList)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("end_a_segment_resource_ids"), resource.EndAResources)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("end_b_segment_resource_ids"), resource.EndBResources)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraSegmentResourceShare) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id int

	req.State.GetAttribute(ctx, path.Root("id"), &id)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.segment.Delete(strconv.Itoa(id))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Segment Resource share",
			"Could not delete segment resource share, unexpected error: "+err.Error(),
		)
		return
	}
}
