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
	_ resource.ResourceWithConfigure = &alkiraSegmentResource{}
)

type alkiraSegmentResource struct {
	client  *alkira.AlkiraClient
	segment *alkira.AlkiraAPI[alkira.Segment]
}

func NewalkiraSegment() resource.Resource {
	return &alkiraSegmentResource{}
}

// Configure adds the provider configured client to the resource.
func (r *alkiraSegmentResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*alkira.AlkiraClient)
	r.segment = alkira.NewSegment(r.client)
}

// Metadata returns the resource type name.
func (r *alkiraSegmentResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_segment"
}

// Schema defines the schema for the resource.
func (r *alkiraSegmentResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"state": schema.StringAttribute{
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
			"asn": schema.Int64Attribute{
				Description: "The BGP ASN for the segment. Default value is 65514.",
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Int64{
					Int64DefaultValue(types.Int64Value(65514)),
				},
			},
			"cidrs": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The list of CIDR blocks.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "Segment description.",
				Optional:    true,
			},
			"enable_ipv6_to_ipv4_translation": schema.BoolAttribute{
				Description: "Enable IPv6 to IPv4 translation in the " +
					"segment for internet application traffic. (**BETA**)",
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					BoolDefaultValue(types.BoolValue(false)),
				},
			},
			"enterprise_dns_server_ip": schema.StringAttribute{
				Description: "The IP of the DNS server used within the segment. This DNS server " +
					"may be used by the Alkira CXP to resolve the names of LDAP servers for example " +
					"which are configured on the Remote Access Connector. (**BETA**)",
				Optional: true,
			},
			"reserve_public_ips": schema.BoolAttribute{
				Description: "Default value is `false`. When this is set to " +
					"`true`. Alkira reserves public IPs " +
					"which can be used to create underlay tunnels between an " +
					"external service and Alkira. For example the reserved public IPs " +
					"may be used to create tunnels to the Akamai Prolexic. (**BETA**)",
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					BoolDefaultValue(types.BoolValue(false)),
				},
			},
			"src_ipv4_pool_start_ip": schema.StringAttribute{
				Description: "The start IP address of IPv4 pool.",
				Optional:    true,
			},
			"src_ipv4_pool_end_ip": schema.StringAttribute{
				Description: "The end IP address of IPv4 pool.",
				Optional:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *alkiraSegmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alkira.Segment
	var cidrs []string
	srcIpv4PoolList := []alkira.SegmentSrcIpv4PoolList{}
	list := alkira.SegmentSrcIpv4PoolList{}
	var startIp string
	var endIp string

	req.Config.GetAttribute(ctx, path.Root("src_ipv4_pool_start_ip"), &startIp)
	req.Config.GetAttribute(ctx, path.Root("src_ipv4_pool_end_ip"), &endIp)
	if startIp != "" && endIp != "" {
		list.StartIp = startIp
		list.EndIp = endIp
		srcIpv4PoolList = append(srcIpv4PoolList, list)
	} else {
		srcIpv4PoolList = nil
	}

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("asn"), &plan.Asn)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("cidrs"), &cidrs)...)
	req.Config.GetAttribute(ctx, path.Root("enable_ipv6_to_ipv4_translation"), &plan.EnableIpv6ToIpv4Translation)
	req.Config.GetAttribute(ctx, path.Root("reserve_public_ips"), &plan.ReservePublicIPsForUserAndSiteConnectivity)
	req.Config.GetAttribute(ctx, path.Root("description"), &plan.Description)
	req.Config.GetAttribute(ctx, path.Root("enterprise_dns_server_ip"), &plan.EnterpriseDNSServerIP)

	plan.IpBlocks = alkira.SegmentIpBlocks{
		Values: cidrs,
	}
	plan.SrcIpv4PoolList = srcIpv4PoolList

	log.Printf("[DEBUG] plan: %v", plan)

	if resp.Diagnostics.HasError() {
		return
	}

	result, state, err := r.segment.Create(&plan)
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("asn"), result.Asn)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("enable_ipv6_to_ipv4_translation"), result.EnableIpv6ToIpv4Translation)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("reserve_public_ips"), result.ReservePublicIPsForUserAndSiteConnectivity)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cidrs"), result.IpBlocks.Values)...)
	if result.Description != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), result.Description)...)
	}
	if result.EnterpriseDNSServerIP != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("enterprise_dns_server_ip"), result.EnterpriseDNSServerIP)...)
	}
	if result.SrcIpv4PoolList != nil && len(result.SrcIpv4PoolList) > 0 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("src_ipv4_pool_start_ip"), result.SrcIpv4PoolList[0].StartIp)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("src_ipv4_pool_end_ip"), result.SrcIpv4PoolList[0].EndIp)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraSegmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var id int

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.segment.GetById(strconv.Itoa(id))
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), result.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("asn"), result.Asn)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("enable_ipv6_to_ipv4_translation"), result.EnableIpv6ToIpv4Translation)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("reserve_public_ips"), result.ReservePublicIPsForUserAndSiteConnectivity)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cidrs"), result.IpBlocks.Values)...)
	if result.Description != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), result.Description)...)
	}
	if result.EnterpriseDNSServerIP != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("enterprise_dns_server_ip"), result.EnterpriseDNSServerIP)...)
	}
	if result.SrcIpv4PoolList != nil && len(result.SrcIpv4PoolList) > 0 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("src_ipv4_pool_start_ip"), result.SrcIpv4PoolList[0].StartIp)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("src_ipv4_pool_end_ip"), result.SrcIpv4PoolList[0].EndIp)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraSegmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan alkira.Segment
	var id int
	var cidrs []string
	srcIpv4PoolList := []alkira.SegmentSrcIpv4PoolList{}
	list := alkira.SegmentSrcIpv4PoolList{}
	var startIp string
	var endIp string

	// Grab the ID from the state.
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)

	req.Config.GetAttribute(ctx, path.Root("src_ipv4_pool_start_ip"), &startIp)
	req.Config.GetAttribute(ctx, path.Root("src_ipv4_pool_end_ip"), &endIp)
	if startIp != "" && endIp != "" {
		list.StartIp = startIp
		list.EndIp = endIp
		srcIpv4PoolList = append(srcIpv4PoolList, list)
	} else {
		srcIpv4PoolList = nil
	}

	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, path.Root("asn"), &plan.Asn)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &plan.Name)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("cidrs"), &cidrs)...)
	req.Config.GetAttribute(ctx, path.Root("enable_ipv6_to_ipv4_translation"), &plan.EnableIpv6ToIpv4Translation)
	req.Config.GetAttribute(ctx, path.Root("reserve_public_ips"), &plan.ReservePublicIPsForUserAndSiteConnectivity)
	req.Config.GetAttribute(ctx, path.Root("description"), &plan.Description)
	req.Config.GetAttribute(ctx, path.Root("enterprise_dns_server_ip"), &plan.EnterpriseDNSServerIP)

	plan.IpBlocks = alkira.SegmentIpBlocks{
		Values: cidrs,
	}
	plan.SrcIpv4PoolList = srcIpv4PoolList

	log.Printf("[DEBUG] plan: %v", plan)

	if resp.Diagnostics.HasError() {
		return
	}

	state, err := r.segment.Update(strconv.Itoa(id), &plan)
	if err != nil {
		return
	}

	result, err := r.segment.GetById(strconv.Itoa(id))
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("state"), state)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), result.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("asn"), result.Asn)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("enable_ipv6_to_ipv4_translation"), result.EnableIpv6ToIpv4Translation)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("reserve_public_ips"), result.ReservePublicIPsForUserAndSiteConnectivity)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cidrs"), result.IpBlocks.Values)...)
	if result.Description == "" {
		resp.State.SetAttribute(ctx, path.Root("description"), types.StringNull())
	} else {
		resp.State.SetAttribute(ctx, path.Root("description"), result.Description)
	}
	if result.EnterpriseDNSServerIP != "" {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("enterprise_dns_server_ip"), result.EnterpriseDNSServerIP)...)
	}
	if result.SrcIpv4PoolList != nil && len(result.SrcIpv4PoolList) > 0 {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("src_ipv4_pool_start_ip"), result.SrcIpv4PoolList[0].StartIp)...)
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("src_ipv4_pool_end_ip"), result.SrcIpv4PoolList[0].EndIp)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}
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
			"Error Deleting Segment",
			"Could not delete group, unexpected error: "+err.Error(),
		)
		return
	}
}
