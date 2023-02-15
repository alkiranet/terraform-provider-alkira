package alkira

import (
	"context"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.ResourceWithConfigure = &alkiraCredentialKeyPairResource{}
)

type alkiraCredentialKeyPairResource struct {
	client *alkira.AlkiraClient
}

func NewalkiraCredentialKeyPair() resource.Resource {
	return &alkiraCredentialKeyPairResource{}
}

// Configure adds the provider configured client to the resource.
func (r *alkiraCredentialKeyPairResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*alkira.AlkiraClient)
}

// Metadata returns the resource type name.
func (r *alkiraCredentialKeyPairResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_credential_ssh_key_pair"
}

// Schema defines the schema for the resource.
func (r *alkiraCredentialKeyPairResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"state": schema.StringAttribute{
				Description: "Provisioning state of the ssh credential keypair.",
				Computed:    true,
			},
			"credential_id": schema.StringAttribute{
				Description: "The ID ssh credential keypair.",
				Computed:    true,
			},
			"name": schema.StringAttribute{
				Description: "The name of the ssh credential keypair.",
				Required:    true,
			},
			"public_key": schema.StringAttribute{
				Description: "The public key.",
				Optional:    true,
			},
			"last_updated": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *alkiraCredentialKeyPairResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan alkira.CredentialKeyPair
	var name string

	plan.Type = "IMPORTED"
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("name"), &name)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("public_key"), &plan.PublicKey)...)
	if resp.Diagnostics.HasError() {
		return
	}

	credentialId, err := r.client.CreateCredential(name, alkira.CredentialTypeKeyPair, plan, 0)
	if err != nil {
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), name)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_key"), plan.PublicKey)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("credential_id"), credentialId)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *alkiraCredentialKeyPairResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *alkiraCredentialKeyPairResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *alkiraCredentialKeyPairResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var credentialId string

	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("credential_id"), &credentialId)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteCredential(credentialId, alkira.CredentialTypeKeyPair)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Credential Key Pair",
			"Could not delete credential key pair, unexpected error: "+err.Error(),
		)
		return
	}
}
