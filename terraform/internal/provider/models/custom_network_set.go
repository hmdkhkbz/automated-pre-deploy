package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type CustomNetworkSetType struct {
	basetypes.SetType
}

func (c CustomNetworkSetType) Type(ctx context.Context) attr.Type {
	return CustomNetworkSetType{
		SetType: basetypes.SetType{
			ElemType: c.SetType.ElementType(),
		},
	}
}

func (c CustomNetworkSetType) Equal(o attr.Type) bool {
	other, ok := o.(CustomNetworkSetType)
	if !ok {
		return false
	}
	return other.ElemType.Equal(c.ElemType)
}

func (c CustomNetworkSetType) String() string {
	return c.SetType.String()
}

func (c CustomNetworkSetType) ValueFromSet(ctx context.Context, set basetypes.SetValue) (basetypes.SetValuable, diag.Diagnostics) {
	x := CustomNetworkSetValue{
		SetValue: set,
	}
	return x, nil
}

func (c CustomNetworkSetType) ValueFromTerraform(ctx context.Context, val tftypes.Value) (attr.Value, error) {
	attrVal, err := c.SetType.ValueFromTerraform(ctx, val)
	if err != nil {
		return nil, err
	}
	setVal, ok := attrVal.(basetypes.SetValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrVal)
	}
	setValuable, d := c.ValueFromSet(ctx, setVal)
	if d.HasError() {
		return nil, fmt.Errorf("unexpected error converting CustomNetworkSetValue to SetValuable: %v", d)
	}
	return setValuable, nil
}

func (c CustomNetworkSetType) ValueType(ctx context.Context) attr.Value {
	return CustomNetworkSetValue{}
}

type CustomNetworkSetValue struct {
	basetypes.SetValue
}

func (v CustomNetworkSetValue) Equal(in attr.Value) bool {
	other, ok := in.(CustomNetworkSetValue)
	if !ok {
		return false
	}
	var otherNets []TFNetworkAttachment
	d := other.ElementsAs(context.Background(), &otherNets, false)
	if d.HasError() {
		return false
	}

	var currentNets []TFNetworkAttachment
	d = v.ElementsAs(context.Background(), &currentNets, false)
	if d.HasError() {
		return false
	}

	if len(currentNets) != len(otherNets) {
		return false
	}

	otherNetSet := make(map[string]bool)
	for _, x := range otherNets {
		otherNetSet[x.NetworkID.ValueString()] = true
	}

	currentNetSet := make(map[string]bool)
	for _, x := range currentNets {
		currentNetSet[x.NetworkID.ValueString()] = true
	}

	for nid, _ := range currentNetSet {
		if !otherNetSet[nid] {
			return false
		}
	}

	for nid, _ := range otherNetSet {
		if !currentNetSet[nid] {
			return false
		}
	}

	return true
}

func (v CustomNetworkSetValue) Type(ctx context.Context) attr.Type {
	return CustomNetworkSetType{
		SetType: basetypes.SetType{
			ElemType: v.ElementType(ctx),
		},
	}
}

func (v CustomNetworkSetValue) SetSemanticEquals(ctx context.Context, newVal basetypes.SetValuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	other, ok := newVal.(CustomNetworkSetValue)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", v)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newVal),
		)
		return false, diags
	}
	eq := v.Equal(other)
	tflog.Warn(ctx, "NETWORK_EQUALITY_CHECK", map[string]interface{}{"EQUAL": eq})
	return eq, diags

}

type CustomNetworkAttachmentType struct {
	basetypes.ObjectType
}

func (n CustomNetworkAttachmentType) Type(ctx context.Context) attr.Type {
	return CustomNetworkAttachmentType{
		ObjectType: basetypes.ObjectType{
			AttrTypes: n.AttrTypes,
		},
	}
}

func (n CustomNetworkAttachmentType) Equal(o attr.Type) bool {
	other, ok := o.(CustomNetworkAttachmentType)
	if !ok {
		return false
	}
	return n.ObjectType.Equal(other)
}

func (n CustomNetworkAttachmentType) String() string {
	return n.ObjectType.String()
}

func (n CustomNetworkAttachmentType) ValueFromObject(ctx context.Context, obj basetypes.ObjectValue) (basetypes.ObjectValuable, diag.Diagnostics) {
	x := CustomNetworkAttachmentValue{
		ObjectValue: obj,
	}
	return x, nil
}

func (n CustomNetworkAttachmentType) ValueFromTerraform(ctx context.Context, tfVal tftypes.Value) (attr.Value, error) {
	attrVal, err := n.ObjectType.ValueFromTerraform(ctx, tfVal)
	if err != nil {
		return nil, err
	}
	objVal, ok := attrVal.(basetypes.ObjectValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type of %T", attrVal)
	}
	ret, d := n.ValueFromObject(ctx, objVal)
	if d.HasError() {
		return nil, fmt.Errorf("unexpected error converting CustomNetworkAttachmentType to ObjectValuable: %v", d)
	}
	return ret, nil
}

func (n CustomNetworkAttachmentType) ValueType(ctx context.Context) attr.Value {
	return CustomNetworkAttachmentValue{}
}

type CustomNetworkAttachmentValue struct {
	basetypes.ObjectValue
	EqualityAttributes []string
}

func (a CustomNetworkAttachmentValue) Equal(in attr.Value) bool {
	other, ok := in.(CustomNetworkAttachmentValue)
	if !ok {
		return false
	}
	otherAttrs := other.Attributes()

	currentAttrs := a.Attributes()

	for _, x := range a.EqualityAttributes {
		v1 := otherAttrs[x].(basetypes.StringValue)
		v2 := currentAttrs[x].(basetypes.StringValue)
		if !v2.Equal(v1) {
			return false
		}
	}
	return true

}

func (a CustomNetworkAttachmentValue) Type(ctx context.Context) attr.Type {
	x, ok := a.ObjectValue.Type(ctx).(basetypes.ObjectType)
	if !ok {
		return nil
	}
	return CustomNetworkAttachmentType{
		ObjectType: x,
	}
}
