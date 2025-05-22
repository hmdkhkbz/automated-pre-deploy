package misc

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/praserx/ipconv"
	"net"
	//"terraform-provider-hashicups-pf/internal/provider/models"
)

type RequiredWhenAttributeHasBoolValue[T comparable] struct {
	expected T
	attr     path.Path
}

func (r *RequiredWhenAttributeHasBoolValue[T]) Description(ctx context.Context) string {
	return fmt.Sprintf("%s must have value %v for this attribute to be set", r.attr.String(), r.expected)
}

func (r *RequiredWhenAttributeHasBoolValue[T]) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("%s must have value %v for this attribute to be set", r.attr.String(), r.expected)
}

func (r *RequiredWhenAttributeHasBoolValue[T]) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var actual T
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, r.attr, &actual)...)
	if resp.Diagnostics.HasError() {
		return
	}
	if actual != r.expected && !req.ConfigValue.IsNull() {
		resp.Diagnostics.AddAttributeError(req.Path, "validation error", fmt.Sprintf(
			"%s must be set only when %s has %v value", req.Path.String(), r.attr.String(), r.expected))
		return
	}

	if actual == r.expected && req.ConfigValue.IsNull() {
		resp.Diagnostics.AddAttributeError(req.Path, "validation error", fmt.Sprintf(
			"%s must be set when %s has %v value", req.Path.String(), r.attr.String(), r.expected))
		return
	}
}

func NewRequiredWhenAttributeHasBoolValue[T comparable](expectedValue T, attributePath path.Path) validator.String {
	return &RequiredWhenAttributeHasBoolValue[T]{
		expected: expectedValue,
		attr:     attributePath,
	}
}

type RequiredToBeEqualToAtLeastOneOfAList[T any] struct {
	listPath path.Path
	selector func([]T) map[string]bool
}

func NewRequiredToBeEqualToAtLeastOneOfAList[T any](lp path.Path, selector func([]T) map[string]bool) validator.String {
	return &RequiredToBeEqualToAtLeastOneOfAList[T]{
		listPath: lp,
		selector: selector,
	}
}

func (r *RequiredToBeEqualToAtLeastOneOfAList[T]) Description(ctx context.Context) string {
	return "checks if the value provided is in a set or list"
}

func (r *RequiredToBeEqualToAtLeastOneOfAList[T]) MarkdownDescription(ctx context.Context) string {
	return "checks if the value provided is in a set or list"
}

func (r *RequiredToBeEqualToAtLeastOneOfAList[T]) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var val []T

	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, r.listPath, &val)...)
	if resp.Diagnostics.HasError() {
		return
	}
	list := r.selector(val)
	cv := req.ConfigValue.ValueString()
	if !list[cv] {
		resp.Diagnostics.AddAttributeError(req.Path, "invalid value", fmt.Sprintf("%s must be one of %v", cv, list))
		return
	}

}

type IPAddressList struct {
}

func (i *IPAddressList) Description(ctx context.Context) string {
	return "validates ipv4 address"
}

func (i *IPAddressList) MarkdownDescription(ctx context.Context) string {
	return "validates ipv4 address"
}

func (i *IPAddressList) ValidateList(ctx context.Context, req validator.ListRequest, resp *validator.ListResponse) {
	var ips []string
	resp.Diagnostics.Append(req.ConfigValue.ElementsAs(ctx, &ips, true)...)
	if resp.Diagnostics.HasError() {
		return
	}
	for _, addr := range ips {
		ip := net.ParseIP(addr)
		if ip == nil {
			resp.Diagnostics.AddAttributeError(req.Path, "invalid attribute value", fmt.Sprintf("%s is not a valid ip address", addr))
			return
		}
	}
}

type IPRangeValidator struct {
	startPath path.Path
	endPath   path.Path
	cidrPath  path.Path
}

func (r *IPRangeValidator) Description(ctx context.Context) string {
	return "validates ip range"
}

func (r *IPRangeValidator) MarkdownDescription(ctx context.Context) string {
	return "validates ip range"
}

func (r *IPRangeValidator) ValidateObject(ctx context.Context, req validator.ObjectRequest, resp *validator.ObjectResponse) {
	var cidr string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, r.cidrPath, &cidr)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var start string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, r.startPath, &start)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "DHCP_START", map[string]interface{}{"START_IP": start})

	var end string
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, r.endPath, &end)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "DHCP_END", map[string]interface{}{"END_IP": start})

	startIP := net.ParseIP(start)
	if startIP == nil {
		resp.Diagnostics.AddAttributeError(r.startPath, "invalid attribute value", fmt.Sprintf("%s is not a valid ip address", startIP))
	}
	endIP := net.ParseIP(end)
	if endIP == nil {
		resp.Diagnostics.AddAttributeError(r.endPath, "invalid attribute value", fmt.Sprintf("%s is not a valid ip address", endIP))
	}

	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		resp.Diagnostics.AddAttributeError(r.cidrPath, "invalid attribute value", err.Error())
	}
	if resp.Diagnostics.HasError() {
		return
	}

	if !ipNet.Contains(startIP) {
		resp.Diagnostics.AddAttributeError(r.startPath, "invalid attribute value", fmt.Sprintf("%s is not in range %s", start, cidr))
	}
	if !ipNet.Contains(endIP) {
		resp.Diagnostics.AddAttributeError(r.startPath, "invalid attribute value", fmt.Sprintf("%s is not in range %s", end, cidr))
	}
	if resp.Diagnostics.HasError() {
		return
	}

	sInt, err := ipconv.IPv4ToInt(startIP)
	if err != nil {
		resp.Diagnostics.AddAttributeError(r.startPath, "unexpected error", err.Error())
		return
	}
	eInt, err := ipconv.IPv4ToInt(endIP)

	if err != nil {
		resp.Diagnostics.AddAttributeError(r.endPath, "unexpected error", err.Error())
		return
	}

	if sInt >= eInt {
		tflog.Info(ctx, "IP_VALIDATION", map[string]interface{}{"START_END": []uint32{binary.BigEndian.Uint32(startIP), binary.BigEndian.Uint32(endIP)}})
		resp.Diagnostics.AddAttributeError(r.startPath, "invalid attribute value", "start ip address must be less than end ip address")
	}

}

func NewIPRangeValidator(cidr, start, end path.Path) validator.Object {
	return &IPRangeValidator{
		startPath: start,
		endPath:   end,
		cidrPath:  cidr,
	}
}
