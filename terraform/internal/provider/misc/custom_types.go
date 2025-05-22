package misc

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"

	"terraform-provider-hashicups-pf/internal/utl"
)

type ListStringType struct {
	basetypes.ListType
}

func (l ListStringType) Type(ctx context.Context) attr.Type {
	// CustomStringType defined in the schema type section
	return ListStringType{
		ListType: basetypes.ListType{
			ElemType: basetypes.StringType{},
		},
	}
}

func (l ListStringType) Equal(o attr.Type) bool {
	other, ok := o.(ListStringType)
	if !ok {
		return false
	}
	_, ok = other.ElemType.(basetypes.StringType)
	if !ok {
		return false
	}
	_, ok = l.ElemType.(basetypes.StringType)
	if !ok {
		return false
	}
	return true
}

func (l ListStringType) String() string {
	return l.ListType.String()
}

func (l ListStringType) ValueFromList(ctx context.Context, list basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	x := ListStringValue{
		ListValue: list,
	}
	return x, nil
}

func (l ListStringType) ValueFromTerraform(ctx context.Context, val tftypes.Value) (attr.Value, error) {
	attrVal, err := l.ListType.ValueFromTerraform(ctx, val)
	if err != nil {
		tflog.Info(ctx, "ERRRRRRRRRRRRRRRRRRRRRRRRRRRRRRR")
		return nil, err
	}
	listVal, ok := attrVal.(basetypes.ListValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrVal)
	}

	listValuable, diags := l.ValueFromList(ctx, listVal)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ListStringValue to ListValuable: %v", diags)
	}
	return listValuable, nil

}
func (l ListStringType) ValueType(ctx context.Context) attr.Value {

	return ListStringValue{}
}

type ListStringValue struct {
	basetypes.ListValue
}

func (l ListStringValue) Equal(in attr.Value) bool {
	other, ok := in.(ListStringValue)
	if !ok {
		return false
	}

	var ll []types.String
	diags := l.ListValue.ElementsAs(context.Background(), &ll, true)
	if diags.HasError() {
		return false
	}
	var ol []types.String
	diags = other.ListValue.ElementsAs(context.Background(), &ol, true)
	if diags.HasError() {
		return false
	}
	if len(ll) != len(ol) {
		return false
	}
	lm := utl.ListToSet(ll)
	om := utl.ListToSet(ol)

	for k, _ := range lm {
		if !om[k] {
			return false
		}
	}

	return true

}

func (l ListStringValue) Type(ctx context.Context) attr.Type {
	lt := basetypes.ListType{
		ElemType: basetypes.StringType{},
	}
	return ListStringType{
		ListType: lt,
	}
}

func (l ListStringValue) ListSemanticEquals(ctx context.Context, newValuable basetypes.ListValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	newVal, ok := newValuable.(ListStringValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", l)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newVal),
		)
		return false, diags
	}

	if l.IsNull() && !newVal.IsNull() && len(newVal.Elements()) == 0 {
		return true, diags
	}
	return l.Equal(newVal), diags
}

func (l ListStringValue) ToStringSlice(ctx context.Context) ([]string, diag.Diagnostics) {
	var ret []string
	d := l.ListValue.ElementsAs(ctx, &ret, true)
	return ret, d
}

func (l ListStringValue) ToTFStringSlice(ctx context.Context) ([]types.String, diag.Diagnostics) {
	var ret []types.String
	d := l.ListValue.ElementsAs(ctx, &ret, true)
	return ret, d
}

func NewListStringValue(elements []string) ListStringValue {
	var in []attr.Value
	for _, x := range elements {
		in = append(in, types.StringValue(x))
	}

	listVal := types.ListValueMust(types.StringType, in)
	return ListStringValue{
		ListValue: listVal,
	}
}

type CustomStringType struct {
	basetypes.StringType
}

func (c CustomStringType) Equal(o attr.Type) bool {
	other, ok := o.(CustomStringType)
	if !ok {
		return false
	}
	return c.StringType.Equal(other.StringType)
}

func (c CustomStringType) String() string {
	return c.StringType.String()
}

func (c CustomStringType) ValueFromString(ctx context.Context, str basetypes.StringValue) (basetypes.StringValuable, diag.Diagnostics) {
	ret := CustomStringValue{
		StringValue: str,
	}
	return ret, nil
}

func (c CustomStringType) ValueFromTerraform(ctx context.Context, val tftypes.Value) (attr.Value, error) {
	tfVal, err := c.StringType.ValueFromTerraform(ctx, val)
	if err != nil {
		return nil, err
	}
	strVal, ok := tfVal.(basetypes.StringValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", tfVal)
	}

	cStingVal, diags := c.ValueFromString(ctx, strVal)
	if diags.HasError() {
		return nil, fmt.Errorf("unexpected error converting ListStringValue to ListValuable: %v", diags)
	}
	return cStingVal, nil
}

type CustomStringValue struct {
	basetypes.StringValue
}

func (c CustomStringValue) Equal(o attr.Value) bool {
	other, ok := o.(CustomStringValue)
	if !ok {

		return false
	}
	return c.StringValue.Equal(other.StringValue)
}

func (c CustomStringValue) Type(ctx context.Context) attr.Type {
	return CustomStringType{}
}

func (c *CustomStringValue) Replace(newVal string) {
	if c.IsNull() && newVal == "" {
		return
	}
	if c.ValueString() != newVal {
		c.StringValue = types.StringValue(newVal)
	}
}

func (c CustomStringValue) StringSemanticEquals(ctx context.Context, newVal basetypes.StringValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	nv, ok := newVal.(CustomStringValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", c)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", nv),
		)
		return false, diags
	}
	tflog.Info(ctx, "StringSemanticEquals", map[string]interface{}{
		"OLD": nv.ValueString(),
		"NEW": c.ValueString(),
	})

	if c.StringValue.IsNull() && !nv.StringValue.IsNull() && nv.ValueString() == "" {
		tflog.Info(ctx, "StringSemanticEquals")
		return true, diags
	}
	if nv.StringValue.IsNull() && !c.StringValue.IsNull() && c.ValueString() == "" {
		tflog.Info(ctx, "StringSemanticEquals")
		return true, diags
	}
	return c.Equal(nv.StringValue), diags
}

func CustomStringValueFromString(ctx context.Context, str string) CustomStringValue {
	s := types.StringValue(str)
	strV, _ := s.ToStringValue(ctx)
	return CustomStringValue{
		StringValue: strV,
	}
}

type CustomObjectListType struct {
	basetypes.ListType
	KeyAttrName string
	Object      basetypes.ObjectType
}

func (c CustomObjectListType) Equal(o attr.Type) bool {
	other, ok := o.(CustomObjectListType)
	if !ok {
		return false
	}
	if other.KeyAttrName != c.KeyAttrName {
		return false
	}
	keyAttrs := strings.Split(c.KeyAttrName, ",")
	for _, k := range keyAttrs {
		ckAttrs, ok := c.Object.AttrTypes[k]
		if !ok {
			return false
		}
		_, ok = ckAttrs.(basetypes.StringType)
		if !ok {
			return false
		}

		okAttr, ok := other.Object.AttrTypes[k]
		if !ok {
			return false
		}
		_, ok = okAttr.(basetypes.StringType)
		if !ok {
			return false
		}
	}

	return c.Object.Equal(c.ListType.ElemType) && other.Object.Equal(other.ListType.ElemType) && c.Object.Equal(other.Object)
}

func (c CustomObjectListType) String() string {
	return c.ListType.String()
}

func (c CustomObjectListType) ValueFromList(ctx context.Context, list basetypes.ListValue) (basetypes.ListValuable, diag.Diagnostics) {
	return CustomObjectListValue{
		ListValue:   list,
		Object:      c.Object,
		KeyAttrName: c.KeyAttrName,
	}, nil

}

func (c CustomObjectListType) ValueFromTerraform(ctx context.Context, val tftypes.Value) (attr.Value, error) {
	nv, err := c.ListType.ValueFromTerraform(ctx, val)
	if err != nil {
		return nil, err
	}

	lisValue, ok := nv.(basetypes.ListValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", lisValue)
	}
	ret, d := c.ValueFromList(ctx, lisValue)
	if d.HasError() {
		return nil, fmt.Errorf("unexpected error %v", d)
	}
	return ret, nil
}

type CustomObjectListValue struct {
	basetypes.ListValue
	KeyAttrName string
	Object      basetypes.ObjectType
}

func (c CustomObjectListValue) Type(ctx context.Context) attr.Type {
	return CustomObjectListType{
		ListType: basetypes.ListType{
			ElemType: c.Object,
		},
		KeyAttrName: c.KeyAttrName,
		Object:      c.Object,
	}
}

func (c CustomObjectListValue) Equal(o attr.Value) bool {
	tflog.Warn(context.Background(), "EQ_OBJECT_LIST")
	other, ok := o.(CustomObjectListValue)
	if !ok {
		return false
	}
	if other.KeyAttrName != c.KeyAttrName {
		return false
	}

	var cl []basetypes.ObjectValue
	for _, x := range c.Elements() {
		o, ok := x.(basetypes.ObjectValue)
		if !ok {
			return false
		}
		if !o.Type(context.Background()).Equal(c.Object) {
			return false
		}
		cl = append(cl, o)
	}

	var ol []basetypes.ObjectValue
	for _, x := range c.Elements() {
		o, ok := x.(basetypes.ObjectValue)
		if !ok {
			return false
		}
		if !o.Type(context.Background()).Equal(c.Object) {
			return false
		}
		ol = append(ol, o)
	}

	if len(ol) != len(cl) {
		tflog.Warn(context.Background(), "NOT_EQ_OBJECT_LIST")
		return false
	}

	var cKeys []types.String
	var oKeys []types.String
	keys := strings.Split(c.KeyAttrName, ",")
	for _, x := range cl {
		attrs := x.Attributes()
		objectKValue := ""
		for _, k := range keys {
			kVal, ok := attrs[k]
			if !ok {
				return false
			}
			kString, ok := kVal.(basetypes.StringValue)
			if !ok {
				return false
			}
			objectKValue += kString.ValueString()

		}
		cKeys = append(cKeys, types.StringValue(objectKValue))

	}

	for _, x := range ol {
		attrs := x.Attributes()
		objectKValue := ""
		for _, k := range keys {
			kVal, ok := attrs[k]
			if !ok {
				return false
			}
			kString, ok := kVal.(basetypes.StringValue)
			if !ok {
				return false
			}
			objectKValue += kString.ValueString()
		}

		oKeys = append(oKeys, types.StringValue(objectKValue))
	}

	oSet := utl.ListToSet(oKeys)
	tflog.Info(context.Background(), "LIST_OBJECT_SEMANTIC_EQUALITY", map[string]interface{}{
		"CURRENT_KEYS": cKeys,
		"OTHER_SET":    oSet,
	})
	for _, x := range cKeys {
		if !oSet[x.ValueString()] {
			tflog.Warn(context.Background(), "LIST_OBJECT_SEMANTIC_EQUALITY", map[string]interface{}{
				"DIFF":         x.ValueString(),
				"CURRENT_KEYS": cKeys,
				"OTHER_SET":    oSet,
			})
			return false
		}
	}
	return true
}

func (c CustomObjectListValue) ListSemanticEquals(ctx context.Context, newValuable basetypes.ListValuable) (bool, diag.Diagnostics) {

	var diags diag.Diagnostics

	other, ok := newValuable.(CustomObjectListValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", c)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", other),
		)
		return false, diags
	}
	eq := c.Equal(other)
	tflog.Warn(ctx, "CUSTOM_OBJECT_LIST_SEMANTIC_EQUALITY", map[string]interface{}{"EQ_CHECK_RESULT": eq})
	return eq, diags
}
