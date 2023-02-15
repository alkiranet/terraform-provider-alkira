package alkira

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Int64DefaultValue(v types.Int64) planmodifier.Int64 {
	return &int64DefaultValuePlanModifier{v}
}

type int64DefaultValuePlanModifier struct {
	DefaultValue types.Int64
}

var _ planmodifier.Int64 = (*int64DefaultValuePlanModifier)(nil)

func (apm *int64DefaultValuePlanModifier) Description(ctx context.Context) string {
	return ""
}

func (apm *int64DefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (apm *int64DefaultValuePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, res *planmodifier.Int64Response) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}

	res.PlanValue = apm.DefaultValue
}

//------------------------------------------------------------------------

// PlanModifyBool is a plan modifier for types.Bool attributes.
func BoolDefaultValue(v types.Bool) planmodifier.Bool {
	return &boolDefaultValuePlanModifier{v}
}

type boolDefaultValuePlanModifier struct {
	DefaultValue types.Bool
}

var _ planmodifier.Bool = (*boolDefaultValuePlanModifier)(nil)

func (apm *boolDefaultValuePlanModifier) Description(ctx context.Context) string {
	return ""
}

func (apm *boolDefaultValuePlanModifier) MarkdownDescription(ctx context.Context) string {
	return ""
}

func (apm *boolDefaultValuePlanModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, res *planmodifier.BoolResponse) {
	// If the attribute configuration is not null, we are done here
	if !req.ConfigValue.IsNull() {
		return
	}

	// If the attribute plan is "known" and "not null", then a previous plan modifier in the sequence
	// has already been applied, and we don't want to interfere.
	if !req.PlanValue.IsUnknown() && !req.PlanValue.IsNull() {
		return
	}
	res.PlanValue = apm.DefaultValue
}
