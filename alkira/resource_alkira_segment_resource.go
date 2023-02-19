package alkira

import (
	"context"
	"strconv"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraSegmentResource{}
)

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
			"segment_id": schema.Int64Attribute{
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

// Create creates the resource and sets the initial Terraform state.
func (r *alkiraSegmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alkira.SegmentResource
	var segId int

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("segment_id"), &segId)...)
	m := alkira.NewSegment(r.client)
	segment, err := m.GetById(strconv.Itoa(segId))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Segment Resource",
			"Could not get segment resource, unexpected error: "+err.Error(),
		)
		return
	}
	plan.Segment = segment.Name

	prefixPathsExpr := path.MatchRoot("group_prefix").AtAnySetValue()
	prefixPaths, _ := req.Config.PathMatches(ctx, prefixPathsExpr)

	if len(prefixPaths) == 0 {
		resp.Diagnostics.AddError(
			"invalid input for segment resource group prefix.",
			"Could not create segment resource, unexpected error: ",
		)
		return
	}

	prefixes := make([]alkira.SegmentResourceGroupPrefix, len(prefixPaths))
	for i, path := range prefixPaths {
		req.Config.GetAttribute(ctx, path.AtName("group_id"), prefixes[i].GroupId)
		req.Config.GetAttribute(ctx, path.AtName("prefix_list_id"), prefixes[i].PrefixListId)
	}

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	plan.GroupPrefixes = prefixes

	result, prov_state, err := r.segment.Create(&plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Segment Resource",
			"Could not create segment resource, unexpected error: "+err.Error(),
		)
		return
	}

	resp.State.SetAttribute(ctx, path.Root("provision_state"), prov_state)
	resp.State.SetAttribute(ctx, path.Root("id"), result.Id)
	resp.State.SetAttribute(ctx, path.Root("name"), result.Name)
	resp.State.SetAttribute(ctx, path.Root("segment_id"), result.Segment)
	for i, path := range prefixPaths {
		resp.State.SetAttribute(ctx, path.AtName("group_id"), result.GroupPrefixes[i].GroupId)
		resp.State.SetAttribute(ctx, path.AtName("prefix_list_id"), result.GroupPrefixes[i].PrefixListId)
	}

	if resp.Diagnostics.HasError() {
		return
	}

}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraSegmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraSegmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraSegmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var id int

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
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
