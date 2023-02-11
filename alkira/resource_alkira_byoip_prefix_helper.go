package alkira

import (
	"context"
	"errors"

	"github.com/alkiranet/alkira-client-go/alkira"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func GenerateByoipRequestCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) (*alkira.Byoip, error) {
	var attributes alkira.ByoipExtraAttributes
	var plan alkira.Byoip
	plan.ExtraAttributes = attributes

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("prefix"), &plan.Prefix)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("cxp"), &plan.Cxp)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("do_not_advertise"), &plan.DoNotAdvertise)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("message"), &plan.ExtraAttributes.Message)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("signature"), &plan.ExtraAttributes.Signature)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("public_key"), &plan.ExtraAttributes.PublicKey)...)

	if resp.Diagnostics.HasError() {
		return nil, errors.New("resp diagnostics has error")
	}

	return &plan, nil
}

func GenerateByoipRequestUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) (*alkira.Byoip, error) {
	var attributes alkira.ByoipExtraAttributes
	plan := new(alkira.Byoip)
	plan.ExtraAttributes = attributes

	var id int
	resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("id"), &id)...)
	// plan.Id = json.Number(strconv.FormatInt(int64(id), 10))
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("prefix"), &plan.Prefix)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("cxp"), &plan.Cxp)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("description"), &plan.Description)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("do_not_advertise"), &plan.DoNotAdvertise)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("message"), &plan.ExtraAttributes.Message)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("signature"), &plan.ExtraAttributes.Signature)...)
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("public_key"), &plan.ExtraAttributes.PublicKey)...)

	if resp.Diagnostics.HasError() {
		return nil, errors.New("resp diagnostics has error")
	}

	return plan, nil
}

func SetByoipStateUpdate(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse, plan *alkira.Byoip) error {
	id, _ := plan.Id.Int64()
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("prefix"), plan.Prefix)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cxp"), plan.Cxp)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), plan.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("do_not_advertise"), plan.DoNotAdvertise)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("message"), plan.ExtraAttributes.Message)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("signature"), plan.ExtraAttributes.Signature)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_key"), plan.ExtraAttributes.PublicKey)...)

	if resp.Diagnostics.HasError() {
		return errors.New("resp diagnostics has error")
	}

	return nil
}

func SetByoipStateCreate(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan *alkira.Byoip) error {
	id, _ := plan.Id.Int64()
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("prefix"), plan.Prefix)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cxp"), plan.Cxp)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), plan.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("do_not_advertise"), plan.DoNotAdvertise)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("message"), plan.ExtraAttributes.Message)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("signature"), plan.ExtraAttributes.Signature)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_key"), plan.ExtraAttributes.PublicKey)...)

	if resp.Diagnostics.HasError() {
		return errors.New("resp diagnostics has error")
	}

	return nil
}

func SetByoipStateRead(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse, plan *alkira.Byoip) error {
	id, _ := plan.Id.Int64()
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), id)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("prefix"), plan.Prefix)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("cxp"), plan.Cxp)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("description"), plan.Description)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("do_not_advertise"), plan.DoNotAdvertise)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("message"), plan.ExtraAttributes.Message)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("signature"), plan.ExtraAttributes.Signature)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("public_key"), plan.ExtraAttributes.PublicKey)...)

	if resp.Diagnostics.HasError() {
		return errors.New("resp diagnostics has error")
	}

	return nil
}
