package alkira

import (
	"context"
	"os"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &alkiraProvider{}
)

// alkiraProviderModel maps provider schema data to a Go type.
type alkiraProviderModel struct {
	Portal    types.String `tfsdk:"portal"`
	Username  types.String `tfsdk:"username"`
	Password  types.String `tfsdk:"password"`
	Provision types.Bool   `tfsdk:"provision"`
}

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &alkiraProvider{}
}

// alkiraProvider is the provider implementation.
type alkiraProvider struct{}

// Metadata returns the provider type name.
func (p *alkiraProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "alkira"
}

// Schema defines the provider-level schema for configuration data.
func (p *alkiraProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"portal": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"provision": schema.BoolAttribute{
				Required: true,
			},
		},
	}
}

// Configure prepares a Alkira API client for data sources and resources.
func (p *alkiraProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config alkiraProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Portal.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Alkira Portal",
			"The provider cannot create the Alkira API client as there is an unknown configuration value for the Alkira API portal. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ALKIRA_PORTAL environment variable.",
		)
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Unknown Alkira Username",
			"The provider cannot create the Alkira API client as there is an unknown configuration value for the Alkira API username. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ALKIRA_USERNAME environment variable.",
		)
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Unknown Alkira Password",
			"The provider cannot create the Alkira API client as there is an unknown configuration value for the Alkira API password. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the ALKIRA_PASSWORD environment variable.",
		)
	}

	if config.Provision.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("provision"),
			"Unknown Alkira Provision",
			"The provider cannot create the Alkira API client as there is an unknown configuration value for the Alkira API provision.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	portal := os.Getenv("ALKIRA_PORTAL")
	username := os.Getenv("ALKIRA_USERNAME")
	password := os.Getenv("ALKIRA_PASSWORD")

	if !config.Portal.IsNull() {
		portal = config.Portal.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	provision := config.Provision.ValueBool()

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if portal == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Alkira API Host",
			"The provider cannot create the Alkira API client as there is a missing or empty value for the Alkira API host. "+
				"Set the host value in the configuration or use the ALKIRA_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"),
			"Missing Alkira API Username",
			"The provider cannot create the Alkira API client as there is a missing or empty value for the Alkira API username. "+
				"Set the username value in the configuration or use the ALKIRA_USERNAME environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"),
			"Missing Alkira API Password",
			"The provider cannot create the Alkira API client as there is a missing or empty value for the Alkira API password. "+
				"Set the password value in the configuration or use the ALKIRA_PASSWORD environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Create a new Alkira client using the configuration values
	client, err := alkira.NewAlkiraClient(portal, username, password, provision)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Create Alkira API Client",
			"An unexpected error occurred when creating the Alkira API client. "+
				"If the error is not clear, please contact the provider developers.\n\n"+
				"Alkira Client Error: "+err.Error(),
		)
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

// DataSources defines the data sources implemented in the provider.
func (p *alkiraProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewAlkiraBillingTagDataSource,
		NewAlkiraGroupDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *alkiraProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewalkiraBillingTag,
		NewalkiraByoipPrefix,
		NewalkiraGroup,
		NewalkiraCredentialKeyPair,
		NewalkiraSegment,
		NewalkiraSegmentResource,
	}
}
