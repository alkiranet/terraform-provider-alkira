package alkira

import (
	"context"
	"strconv"
	"time"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraPolicyRule{}
)

type alkiraPolicyRuleModel struct {
	ProvisionState        types.String `tfsdk:"provision_state"`
	Id                    types.Int64  `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Description           types.String `tfsdk:"description"`
	ApplicationList       types.List   `tfsdk:"application_ids"`
	ApplicationFamilyList types.List   `tfsdk:"application_family_ids"`
	Dscp                  types.String `tfsdk:"dscp"`
	DstIp                 types.String `tfsdk:"dst_ip"`
	SrcIp                 types.String `tfsdk:"src_ip"`
	SrcPortList           types.List   `tfsdk:"src_ports"`
	DstPortList           types.List   `tfsdk:"dst_ports"`
	SrcPrefixListId       types.Int64  `tfsdk:"src_prefix_list_id"`
	DstPrefixListId       types.Int64  `tfsdk:"dst_prefix_list_id"`
	InternetApplicationId types.Int64  `tfsdk:"internet_application_id"`
	Protocol              types.String `tfsdk:"protocol"`
	Action                types.String `tfsdk:"rule_action"`
	ServiceTypeList       types.List   `tfsdk:"rule_action_service_types"`
	ServiceList           types.List   `tfsdk:"rule_action_service_ids"`
	LastUpdated           types.String `tfsdk:"last_updated"`
}

type alkiraPolicyRule struct {
	client *alkira.AlkiraClient
	policy *alkira.AlkiraAPI[alkira.TrafficPolicyRule]
}

func NewalkiraPolicyRule() resource.Resource {
	return &alkiraPolicyRule{}
}

// Configure adds the provider configured client to the resource.
func (r *alkiraPolicyRule) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*alkira.AlkiraClient)
	r.policy = alkira.NewTrafficPolicyRule(r.client)
}

// Metadata returns the resource type name.
func (r *alkiraPolicyRule) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_policy_rule"
}

// Schema defines the schema for the resource.
func (r *alkiraPolicyRule) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"provision_state": schema.StringAttribute{
				Description: "Provisioning state of the policy rule.",
				Computed:    true,
			},
			"id": schema.Int64Attribute{
				Description: "The ID Segment.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the policy rule.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Segment description.",
				Optional:    true,
			},
			"application_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"application_family_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"dscp": schema.StringAttribute{
				Description: "The dscp value can be any or between 0 and 63 inclusive.",
				Required:    true,
			},
			"dst_ip": schema.StringAttribute{
				Description: "A single destination IP as the match condition of the rule.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Root("dst_prefix_list_id").Expression()),
				},
			},
			"src_ip": schema.StringAttribute{
				Description: "A single source IP as the match condition of the rule.",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.Root("src_prefix_list_id").Expression()),
				},
			},
			"src_ports": schema.ListAttribute{
				Description: "Source ports that can take values: any or 1 to 65535.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"dst_ports": schema.ListAttribute{
				Description: "Destination ports that can take values: any or 1 to 65535.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"src_prefix_list_id": schema.Int64Attribute{
				Description: "The ID of the source prefix list.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.Root("src_ip").Expression()),
				},
			},
			"dst_prefix_list_id": schema.Int64Attribute{
				Description: "The ID of the destination prefix list associated with the rule.",
				Optional:    true,
				Validators: []validator.Int64{
					int64validator.ConflictsWith(path.Root("dst_ip").Expression()),
				},
			},
			"internet_application_id": schema.Int64Attribute{
				Description: "The ID of the internet_application associated with the " +
					"rule. When an internet applciation is selected, destination IP " +
					"and port will be the private IP and port of the application.",
				Optional: true,
			},
			"protocol": schema.StringAttribute{
				Description: "The following protocols are supported. icmp, tcp, udp or any.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.OneOf("icmp", "tcp", "udp", "any"),
				},
			},
			"rule_action": schema.StringAttribute{
				Description: "The action that is applied on matched traffic, " +
					"either ALLOW or DROP. The default value is ALLOW.",
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					StringDefaultValue(types.StringValue("ALLOW")),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("ALLOW", "DROP"),
				},
			},
			"rule_action_service_types": schema.ListAttribute{
				Description: "Based on the service type, traffic is routed to service " +
					"of the given type. For service chaining, both PAN and ZIA service " +
					"types can be selected but must follow order.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"rule_action_service_ids": schema.ListAttribute{
				Description: "Based on the service IDs, traffic is routed to the " +
					"specified services. For service chaining, both service_pan " +
					"and service_zscaler's IDs can be added here, but ID of " +
					"service_pan must be by followed by ID of service_zscaler.",
				ElementType: types.Int64Type,
				Optional:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *alkiraPolicyRule) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alkiraPolicyRuleModel
	var srcPortList []string
	var dstPortList []string
	var appIds []int
	var appFamList []int
	var serviceTypeList []string
	var serviceList []int

	resp.Diagnostics.Append(req.Config.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	resp.Diagnostics.Append(plan.SrcPortList.ElementsAs(ctx, &srcPortList, true)...)
	resp.Diagnostics.Append(plan.DstPortList.ElementsAs(ctx, &dstPortList, true)...)
	resp.Diagnostics.Append(plan.ApplicationList.ElementsAs(ctx, &appIds, true)...)
	resp.Diagnostics.Append(plan.ApplicationFamilyList.ElementsAs(ctx, &appFamList, true)...)
	resp.Diagnostics.Append(plan.ServiceTypeList.ElementsAs(ctx, &serviceTypeList, true)...)
	resp.Diagnostics.Append(plan.ServiceList.ElementsAs(ctx, &serviceList, true)...)

	if resp.Diagnostics.HasError() {
		return
	}

	policy := &alkira.TrafficPolicyRule{
		Name:        plan.Name.ValueString(),
		Description: plan.Description.ValueString(),
		MatchCondition: alkira.PolicyRuleMatchCondition{
			SrcIp:                 plan.SrcIp.ValueString(),
			DstIp:                 plan.DstIp.ValueString(),
			Dscp:                  plan.Dscp.ValueString(),
			Protocol:              plan.Protocol.ValueString(),
			SrcPortList:           srcPortList,
			DstPortList:           dstPortList,
			SrcPrefixListId:       int(plan.SrcPrefixListId.ValueInt64()),
			DstPrefixListId:       int(plan.DstPrefixListId.ValueInt64()),
			InternetApplicationId: int(plan.InternetApplicationId.ValueInt64()),
			ApplicationList:       appIds,
			ApplicationFamilyList: appFamList,
		},
		RuleAction: alkira.PolicyRuleAction{
			Action:          plan.Action.ValueString(),
			ServiceTypeList: serviceTypeList,
			ServiceList:     serviceList,
		},
	}

	result, provState, err := r.policy.Create(policy)
	if err != nil {
		return
	}

	id, err := result.Id.Int64()
	if err != nil {
		return
	}

	plan.Id = types.Int64Value(id)
	plan.ProvisionState = types.StringValue(provState)
	plan.Description = types.StringValue(result.Description)
	plan.SrcIp = types.StringValue(result.MatchCondition.SrcIp)
	plan.DstIp = types.StringValue(result.MatchCondition.DstIp)
	plan.Dscp = types.StringValue(result.MatchCondition.Dscp)
	plan.Protocol = types.StringValue(result.MatchCondition.Protocol)
	if !plan.SrcPrefixListId.IsNull() {
		plan.SrcPrefixListId = types.Int64Value(int64(result.MatchCondition.SrcPrefixListId))
	}
	if result.MatchCondition.DstPrefixListId == 0 {
		plan.DstPrefixListId = types.Int64Null()
	} else {
		plan.DstPrefixListId = types.Int64Value(int64(result.MatchCondition.DstPrefixListId))
	}
	if result.MatchCondition.InternetApplicationId == 0 {
		plan.InternetApplicationId = types.Int64Null()
	} else {
		plan.InternetApplicationId = types.Int64Value(int64(result.MatchCondition.InternetApplicationId))
	}
	plan.Action = types.StringValue(result.RuleAction.Action)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.State.Set(ctx, plan)
	resp.State.SetAttribute(ctx, path.Root("src_ports"), result.MatchCondition.SrcPortList)
	resp.State.SetAttribute(ctx, path.Root("dst_ports"), result.MatchCondition.DstPortList)
	if !plan.ApplicationFamilyList.IsNull() {
		resp.State.SetAttribute(ctx, path.Root("application_family_ids"), result.MatchCondition.ApplicationFamilyList)
	}
	if !plan.ApplicationList.IsNull() {
		resp.State.SetAttribute(ctx, path.Root("application_ids"), result.MatchCondition.ApplicationList)
	}
	resp.State.SetAttribute(ctx, path.Root("rule_action_service_types"), result.RuleAction.ServiceTypeList)
	resp.State.SetAttribute(ctx, path.Root("rule_action_service_ids"), result.RuleAction.ServiceList)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraPolicyRule) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var plan alkiraPolicyRuleModel

	resp.Diagnostics.Append(req.State.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.policy.GetById(strconv.Itoa(int(plan.Id.ValueInt64())))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Policy Rule Resource",
			"Could not read policy rule resource, unexpected error: "+err.Error(),
		)
		return
	}

	id, err := result.Id.Int64()
	if err != nil {
		return
	}

	// plan.Id = types.Int64Value(id)
	// plan.Description = types.StringValue(result.Description)
	// plan.SrcIp = types.StringValue(result.MatchCondition.SrcIp)
	// plan.DstIp = types.StringValue(result.MatchCondition.DstIp)
	// plan.Dscp = types.StringValue(result.MatchCondition.Dscp)
	// plan.Protocol = types.StringValue(result.MatchCondition.Protocol)
	// plan.SrcPrefixListId = types.Int64Value(int64(result.MatchCondition.SrcPrefixListId))
	// plan.DstPrefixListId = types.Int64Value(int64(result.MatchCondition.DstPrefixListId))
	// plan.InternetApplicationId = types.Int64Value(int64(result.MatchCondition.InternetApplicationId))
	// plan.Action = types.StringValue(result.RuleAction.Action)

	// resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("src_ports"), result.MatchCondition.SrcPortList)...)
	// resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("dst_ports"), result.MatchCondition.DstPortList)...)
	// resp.Diagnostics.Append(
	// 	resp.State.SetAttribute(ctx, path.Root("application_ids"), result.MatchCondition.ApplicationList)...)
	// resp.Diagnostics.Append(
	// 	resp.State.SetAttribute(ctx, path.Root("application_family_ids"), result.MatchCondition.ApplicationFamilyList)...)
	// resp.Diagnostics.Append(
	// 	resp.State.SetAttribute(ctx, path.Root("rule_action_service_types"), result.RuleAction.ServiceTypeList)...)
	// resp.Diagnostics.Append(
	// 	resp.State.SetAttribute(ctx, path.Root("rule_action_service_ids"), result.RuleAction.ServiceList)...)

	plan.Id = types.Int64Value(id)
	plan.Description = types.StringValue(result.Description)
	plan.SrcIp = types.StringValue(result.MatchCondition.SrcIp)
	plan.DstIp = types.StringValue(result.MatchCondition.DstIp)
	plan.Dscp = types.StringValue(result.MatchCondition.Dscp)
	plan.Protocol = types.StringValue(result.MatchCondition.Protocol)
	if !plan.SrcPrefixListId.IsNull() {
		plan.SrcPrefixListId = types.Int64Value(int64(result.MatchCondition.SrcPrefixListId))
	}
	if result.MatchCondition.DstPrefixListId == 0 {
		plan.DstPrefixListId = types.Int64Null()
	} else {
		plan.DstPrefixListId = types.Int64Value(int64(result.MatchCondition.DstPrefixListId))
	}
	if result.MatchCondition.InternetApplicationId == 0 {
		plan.InternetApplicationId = types.Int64Null()
	} else {
		plan.InternetApplicationId = types.Int64Value(int64(result.MatchCondition.InternetApplicationId))
	}
	plan.Action = types.StringValue(result.RuleAction.Action)
	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	resp.State.Set(ctx, plan)
	resp.State.SetAttribute(ctx, path.Root("src_ports"), result.MatchCondition.SrcPortList)
	resp.State.SetAttribute(ctx, path.Root("dst_ports"), result.MatchCondition.DstPortList)
	if !plan.ApplicationFamilyList.IsNull() {
		resp.State.SetAttribute(ctx, path.Root("application_family_ids"), result.MatchCondition.ApplicationFamilyList)
	}
	if !plan.ApplicationList.IsNull() {
		resp.State.SetAttribute(ctx, path.Root("application_ids"), result.MatchCondition.ApplicationList)
	}
	resp.State.SetAttribute(ctx, path.Root("rule_action_service_types"), result.RuleAction.ServiceTypeList)
	resp.State.SetAttribute(ctx, path.Root("rule_action_service_ids"), result.RuleAction.ServiceList)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraPolicyRule) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
func (r *alkiraPolicyRule) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
