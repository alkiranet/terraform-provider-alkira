package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &alkiraGroupDataSource{}
	_ datasource.DataSourceWithConfigure = &alkiraGroupDataSource{}
)

type alkiraGroupDataSource struct {
	client *alkira.AlkiraClient
	group  *alkira.AlkiraAPI[alkira.Group]
}

func NewAlkiraGroupDataSource() datasource.DataSource {
	return &alkiraGroupDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *alkiraGroupDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*alkira.AlkiraClient)
	d.group = alkira.NewGroup(d.client)
}

func (d *alkiraGroupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *alkiraGroupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"state": schema.StringAttribute{
				Description: "Provisioning state of the billing tag.",
				Computed:    true,
			},
			"id": schema.NumberAttribute{
				Description: "The ID billing tag.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the billing tag.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the billing tag.",
				Computed:    true,
			},
		},
	}
}

func (d *alkiraGroupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var name string

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, state, err := d.group.GetByName(name)
	if err != nil {
		return
	}

	id, _ := result.Id.Int64()
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("state"), state)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), result.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), result.Description)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
