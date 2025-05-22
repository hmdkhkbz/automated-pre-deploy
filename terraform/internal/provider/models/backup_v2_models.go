package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-hashicups-pf/internal/api"
	"time"
)

var (
	BackupV2DetailItemType = map[string]attr.Type{
		"provisioned_size": types.Int64Type,
		"used_size":        types.Float64Type,
		"created_at":       types.StringType,
		"backup_id":        types.StringType,
		"status":           types.StringType,
		"slot_name":        types.StringType,
		"fail_reason":      types.StringType,
	}

	BackupV2ItemType = map[string]attr.Type{
		"backup_name":   types.StringType,
		"instance_id":   types.StringType,
		"instance_name": types.StringType,
		"status":        types.StringType,
		"quota":         types.Int64Type,
		"occupancy":     types.Int64Type,
		"next_backup":   types.StringType,
		"labels": types.ListType{
			ElemType: types.StringType,
		},
		"details": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: BackupV2DetailItemType,
			},
		},
	}
)

type TFBackupV2 struct {
	Region  types.String `tfsdk:"region"`
	Backups types.List   `tfsdk:"backups"`
}

func (b *TFBackupV2) SetBackups(items []BackupV2Item) diag.Diagnostics {
	var vals []attr.Value
	for _, x := range items {
		vals = append(vals, types.ObjectValueMust(BackupV2ItemType, map[string]attr.Value{
			"backup_name":   x.Name,
			"instance_id":   x.InstanceID,
			"instance_name": x.InstanceName,
			"status":        x.Status,
			"quota":         x.Quota,
			"occupancy":     x.Occupancy,
			"next_backup":   x.NextBackup,
			"labels":        x.Labels,
			"details":       x.Details,
		}))
	}
	l, d := types.ListValue(types.ObjectType{AttrTypes: BackupV2ItemType}, vals)
	if d.HasError() {
		return d
	}
	b.Backups = l
	return d
}

type BackupV2Item struct {
	Name         types.String `tfsdk:"backup_name"`
	InstanceID   types.String `tfsdk:"instance_id"`
	InstanceName types.String `tfsdk:"instance_name"`
	Status       types.String `tfsdk:"status"`
	Quota        types.Int64  `tfsdk:"quota"`
	Occupancy    types.Int64  `tfsdk:"occupancy"`
	NextBackup   types.String `tfsdk:"next_backup"`
	Labels       types.List   `tfsdk:"labels"`
	Details      types.List   `tfsdk:"details"`
}

func (bi *BackupV2Item) SetLabels(ctx context.Context, labels []string) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, types.StringType, labels)
	if d.HasError() {
		return d
	}
	bi.Labels = l
	return d
}

func (bi *BackupV2Item) SetDetails(ctx context.Context, details []api.BackupDetailsData) diag.Diagnostics {
	var tfDetails []attr.Value
	for _, x := range details {

		cAt := time.Unix(x.CreateAt/1000, 0).Format(time.RFC3339)
		tfDetails = append(tfDetails, types.ObjectValueMust(BackupV2DetailItemType, map[string]attr.Value{
			"provisioned_size": types.Int64Value(int64(x.ProvisionedSize)),
			"used_size":        types.Float64Value(x.UsedSize),
			"created_at":       types.StringValue(cAt),
			"backup_id":        types.StringValue(x.BackupID),
			"status":           types.StringValue(x.Status),
			"slot_name":        types.StringValue(x.SlotName),
			"fail_reason":      types.StringValue(x.FailReason),
		}))
	}
	l, d := types.ListValue(types.ObjectType{AttrTypes: BackupV2DetailItemType}, tfDetails)
	if d.HasError() {
		return d
	}
	bi.Details = l
	return d
}

type BackupV2DetailItem struct {
	ProvisionedSize types.Int64   `tfsdk:"provisioned_size"`
	UsedSize        types.Float64 `tfsdk:"used_size"`
	CreatedAt       types.String  `tfsdk:"created_at"`
	BackupID        types.String  `tfsdk:"backup_id"`
	Status          types.String  `tfsdk:"status"`
	SlotName        types.String  `tfsdk:"slot_name"`
	FailReason      types.String  `tfsdk:"fail_reason"`
}
