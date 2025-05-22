package misc

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

type RemovePlanModifier struct {
}

func (r RemovePlanModifier) Description(context.Context) string {
	return ""
}

func (r RemovePlanModifier) MarkdownDescription(context.Context) string {
	return ""
}

func (r RemovePlanModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	/*tflog.Warn(ctx, "PLAN_MODIFIER", map[string]interface{}{
		"PATH": req.Path.String(),
	})
	attrs := req.StateValue.Attributes()
	tflog.Warn(ctx, "STATE_VALUE", map[string]interface{}{
		"ATTRIBUTES": attrs,
	})*/
	if req.StateValue.IsNull() {
		return
	}

	// Do nothing if there is a known planned value.
	if !req.PlanValue.IsUnknown() {
		return
	}

	// Do nothing if there is an unknown configuration value, otherwise interpolation gets messed up.
	if req.ConfigValue.IsUnknown() {
		return
	}

	resp.PlanValue = req.StateValue
}
