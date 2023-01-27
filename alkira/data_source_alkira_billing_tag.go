package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &alkiraBillingTagDataSource{}
	_ datasource.DataSourceWithConfigure = &alkiraBillingTagDataSource{}
)

type alkiraBillingTagDataSource struct {
	client *alkira.AlkiraClient
}

type alkiraBillingTagDataSourceModel struct {
	Name types.String `tfsdk:"name"`
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
			"name": schema.StringAttribute{
				Description: "The name of the billing tag.",
				Required:    true,
			},
		},
	}
}

func (d *alkiraBillingTagDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state alkiraBillingTagDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	billingTag, err := d.client.GetBillingTagByName(state.Name.String())

	if err != nil {
		return
	}

	// Set state
	diags := resp.State.Set(ctx, &billingTag)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
