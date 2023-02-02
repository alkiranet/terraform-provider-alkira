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
	_ datasource.DataSource              = &alkiraBillingTagDataSource{}
	_ datasource.DataSourceWithConfigure = &alkiraBillingTagDataSource{}
)

type alkiraBillingTagDataSource struct {
	client *alkira.AlkiraClient
}

func NewAlkiraBillingTagDataSource() datasource.DataSource {
	return &alkiraBillingTagDataSource{}
}

// Configure adds the provider configured client to the data source.
func (d *alkiraBillingTagDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*alkira.AlkiraClient)
}

func (d *alkiraBillingTagDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_billing_tag"
}

func (d *alkiraBillingTagDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
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

func (d *alkiraBillingTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state alkira.BillingTag

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &state.Name)...)
	state, err := d.client.GetBillingTagByName(string(state.Name))
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), state.Id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), state.Name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), state.Description)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// var state alkiraBillingTagDataSourceModel

	// resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &state.Name)...)

	// billingTag, err := d.client.GetBillingTagByName(state.Name.ValueString())

	// if err != nil {
	// 	return
	// }

	// // Set state
	// state.Id = types.Int64Value(int64(billingTag.Id))
	// diags := resp.State.Set(ctx, &state)

	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}
