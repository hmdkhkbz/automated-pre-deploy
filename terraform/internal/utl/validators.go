package utl

import (
	"context"
	"net"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func CIDRValidator() cidrValidator {
	return cidrValidator{}
}

type cidrValidator struct {}

func (v cidrValidator) ValidateString(c context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	_, _, err := net.ParseCIDR(req.ConfigValue.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Value must be in CIDR notation format", err.Error())
	}
}

func (v cidrValidator) Description(c context.Context) string {
	return ""
}

func (v cidrValidator) MarkdownDescription(context.Context) string {
	return ""
}



func PortValidator() portValidator {
	return portValidator{}
}

type portValidator struct {}

func (v portValidator) ValidateString(c context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	p, err := strconv.Atoi(req.ConfigValue.ValueString())
	if err != nil || p < 0 || p > 65535 {
		resp.Diagnostics.AddError("Value must be a number between 0 and 65535", "Port values are required to be in 0-65535 range")
	}
}

func (v portValidator) Description(c context.Context) string {
	return ""
}

func (v portValidator) MarkdownDescription(context.Context) string {
	return ""
}