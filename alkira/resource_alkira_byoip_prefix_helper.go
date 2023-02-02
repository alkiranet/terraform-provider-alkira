package alkira

import (
	"context"
	"errors"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func GenerateByoipRequest(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) (*alkira.Byoip, error) {
	var attributes alkira.ByoipExtraAttributes
	plan := new(alkira.Byoip)
	plan.ExtraAttributes = attributes

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("prefix"), &plan.Prefix)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("cxp"), &plan.Cxp)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("do_not_advertise"), &plan.DoNotAdvertise)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("message"), &plan.ExtraAttributes.Message)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("signature"), &plan.ExtraAttributes.Signature)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("public_key"), &plan.ExtraAttributes.PublicKey)...)

	if resp.Diagnostics.HasError() {
		return nil, errors.New("resp diagnostics has error.")
	}

	return plan, nil
}

func SetByoipState(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
}
