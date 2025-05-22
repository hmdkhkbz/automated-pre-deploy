package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	networkType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"subnet_id":             types.StringType,
			"ip":                    types.StringType,
			"network_id":            types.StringType,
			"port_id":               types.StringType,
			"is_public":             types.BoolType,
			"port_security_enabled": types.BoolType,
		},
	}

	floatingIPType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"floating_ip_id": types.StringType,
			"network_id":     types.StringType,
		},
	}

	revertToSnapshotObjType = basetypes.ObjectType{
		AttrTypes: map[string]attr.Type{
			"revert_to": types.StringType,
		},
	}
)

type TFServerFlavor struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type TFServerImage struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	OS        types.String `tfsdk:"os"`
	OSVersion types.String `tfsdk:"os_version"`
	Metadata  types.Map    `tfsdk:"metadata"`
}

type TFServerAddress struct {
	MAC      types.String `tfsdk:"mac"`
	Version  types.String `tfsdk:"version"`
	Addr     types.String `tfsdk:"address"`
	Type     types.String `tfsdk:"type"`
	IsPublic types.Bool   `tfsdk:"is_public"`
}

type TFTag struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type TFServerRevertTo struct {
	SnapshotID types.String `tfsdk:"snapshot_id"`
}

type TFInstanceDetails struct {
	ID             types.String                 `tfsdk:"id"`
	Name           types.String                 `tfsdk:"name"`
	Flavor         TFServerFlavor               `tfsdk:"flavor"`
	Status         types.String                 `tfsdk:"status"`
	Image          TFServerImage                `tfsdk:"image"`
	Created        types.String                 `tfsdk:"created"`
	Password       types.String                 `tfsdk:"password"`
	TaskState      types.String                 `tfsdk:"task_state"`
	KeyName        types.String                 `tfsdk:"key_name"`
	SecurityGroups []TFSecurityGroup            `tfsdk:"security_groups"`
	Addresses      map[string][]TFServerAddress `tfsdk:"addresses"`
	Tags           []TFTag                      `tfsdk:"tags"`
	HAEnabled      types.Bool                   `tfsdk:"ha_enabled"`
}

type TFInstanceDatasourceModel struct {
	Region    types.String        `tfsdk:"region"`
	Instances []TFInstanceDetails `tfsdk:"instances"`
}

type TFServerVolume struct {
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Size        types.Int64  `tfsdk:"size"`
}

type TFInstanceResourceModel struct {
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	Region            types.String   `tfsdk:"region"`
	ID                types.String   `tfsdk:"id"`
	TaskID            types.String   `tfsdk:"task_id"`
	Name              types.String   `tfsdk:"name"`
	ImageID           types.String   `tfsdk:"image_id"`
	Networks          types.List     `tfsdk:"networks"`
	FlavorID          types.String   `tfsdk:"flavor_id"`
	SecurityGroups    types.Set      `tfsdk:"security_groups"`
	DiskSize          types.Int64    `tfsdk:"disk_size"`
	InitScript        types.String   `tfsdk:"init_script"`
	Volumes           types.Set      `tfsdk:"volumes"`
	SSHKeyName        types.String   `tfsdk:"ssh_key_name"`
	Password          types.String   `tfsdk:"password"`
	Status            types.String   `tfsdk:"status"`
	FloatingIP        types.Object   `tfsdk:"floating_ip"`
	Snapshot          types.Object   `tfsdk:"revert_to"`
	ServerGroupID     types.String   `tfsdk:"server_group_id"`
	DedicatedServerID types.String   `tfsdk:"dedicated_server_id"`
	ClusterID         types.String   `tfsdk:"cluster_id"`
	SnapshotID        types.String   `tfsdk:"snapshot_id"`
	EnableIPv4        types.Bool     `tfsdk:"enable_ipv4"`
	EnableIPv6        types.Bool     `tfsdk:"enable_ipv6"`
}

func (i *TFInstanceResourceModel) GetSnapshot(ctx context.Context) (*TFServerRevertTo, diag.Diagnostics) {
	var ret TFServerRevertTo
	var d diag.Diagnostics
	if i.Snapshot.IsNull() || i.Snapshot.IsUnknown() {
		return nil, d
	}
	d = i.Snapshot.As(ctx, &ret, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	return &ret, d
}

func (i *TFInstanceResourceModel) GetFloatingIPAttachment(ctx context.Context) (*TFFloatingIPAttachment, diag.Diagnostics) {
	var ret TFFloatingIPAttachment
	var d diag.Diagnostics
	if i.FloatingIP.IsNull() || i.FloatingIP.IsUnknown() {
		return nil, d
	}
	d = i.FloatingIP.As(ctx, &ret, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    true,
		UnhandledUnknownAsEmpty: true,
	})
	if d.HasError() {
		return nil, d
	}
	return &ret, d
}

func (i *TFInstanceResourceModel) SetFloatingIPAttachment(ctx context.Context, fAt *TFFloatingIPAttachment) diag.Diagnostics {
	var d diag.Diagnostics
	if fAt == nil {
		i.FloatingIP = types.ObjectNull(floatingIPType.AttrTypes)
		return d
	}
	obj, d := types.ObjectValue(floatingIPType.AttrTypes, map[string]attr.Value{
		"floating_ip_id": fAt.FloatingID,
		"network_id":     fAt.NetworkID,
	})
	if d.HasError() {
		return d
	}
	i.FloatingIP = obj
	return d
}

func (i *TFInstanceResourceModel) GetNetworkAttachments(ctx context.Context) ([]TFNetworkAttachment, diag.Diagnostics) {
	var ret []TFNetworkAttachment
	d := i.Networks.ElementsAs(ctx, &ret, false)
	return ret, d
}

func (i *TFInstanceResourceModel) SetNetworkAttachments(ctx context.Context, attachments []TFNetworkAttachment) diag.Diagnostics {
	lv, diags := types.ListValueFrom(ctx, networkType, attachments)
	if diags.HasError() {
		return diags
	}
	i.Networks = lv
	return nil
}

func (i *TFInstanceResourceModel) SetNetworkAttachmentsIfNotEqual(ctx context.Context, attachments []TFNetworkAttachment) diag.Diagnostics {
	lv, diags := types.ListValueFrom(ctx, networkType, attachments)
	if diags.HasError() {
		return diags
	}

	if !i.Networks.Equal(lv) {
		i.Networks = lv
	}
	return nil
}

func (i *TFInstanceResourceModel) SetVolumes(ctx context.Context, volumes []string) diag.Diagnostics {
	sv, d := types.SetValueFrom(ctx, types.StringType, volumes)
	if d.HasError() {
		return d
	}
	i.Volumes = sv
	return nil
}

func (i *TFInstanceResourceModel) GetVolumes(ctx context.Context) ([]string, diag.Diagnostics) {
	var ret []string
	d := i.Volumes.ElementsAs(ctx, &ret, true)
	if d.HasError() {
		return nil, d
	}
	return ret, nil
}

func (i *TFInstanceResourceModel) SetSecurityGroups(ctx context.Context, sgs map[string]bool) diag.Diagnostics {
	tflog.Info(ctx, "SECGS", map[string]interface{}{"IDS": sgs})
	var sgSlc []string
	for k, _ := range sgs {
		sgSlc = append(sgSlc, k)
	}
	sv, d := types.SetValueFrom(ctx, types.StringType, sgSlc)
	if d.HasError() {
		return d
	}
	i.SecurityGroups = sv
	return nil
}

func (i *TFInstanceResourceModel) GetSecurityGroups(ctx context.Context) ([]string, diag.Diagnostics) {
	var ret []string
	d := i.SecurityGroups.ElementsAs(ctx, &ret, true)
	if d.HasError() {
		return nil, d
	}
	return ret, nil
}

// HasNetworkAttachment checks wether the instance has a network attachment completely equal to the given network attachment
func (i *TFInstanceResourceModel) HasEqualNetworkAttachment(ctx context.Context, target TFNetworkAttachment) (bool, diag.Diagnostics) {
	attachments, d := i.GetNetworkAttachments(ctx)
	for _, a := range attachments {
		if a.Equals(target) {
			return true, d
		}
	}

	return false, d
}

// IsAttachedToNetwork returns the network attachment for the instance with
// a specific network, nil if not exists
func (i *TFInstanceResourceModel) GetNetworkAttachment(ctx context.Context, networkID string) (*TFNetworkAttachment, diag.Diagnostics) {
	attachments, d := i.GetNetworkAttachments(ctx)
	for _, a := range attachments {
		if a.NetworkID.ValueString() == networkID {
			tflog.Info(ctx, "NETOWKRKRKR", map[string]interface{}{"networkID": a.NetworkID.ValueString()})
			return &a, d
		}
	}

	return nil, d
}

type TFNetworkAttachment struct {
	IP                  types.String `tfsdk:"ip"`
	SubnetID            types.String `tfsdk:"subnet_id"`
	NetworkID           types.String `tfsdk:"network_id"`
	PortID              types.String `tfsdk:"port_id"`
	IsPublic            types.Bool   `tfsdk:"is_public"`
	PortSecurityEnabled types.Bool   `tfsdk:"port_security_enabled"`
}

func (a *TFNetworkAttachment) Equals(t TFNetworkAttachment) bool {
	return *a == t
}

type TFFloatingIPAttachment struct {
	FloatingID types.String `tfsdk:"floating_ip_id"`
	NetworkID  types.String `tfsdk:"network_id"`
}
