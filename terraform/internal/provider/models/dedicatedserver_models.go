package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	DedicatedServerItemType = map[string]attr.Type{
		"id":           types.StringType,
		"name":         types.StringType,
		"type_id":      types.StringType,
		"sockets":      types.Int64Type,
		"vcpus":        types.Int64Type,
		"vcpus_used":   types.Int64Type,
		"memory":       types.Int64Type,
		"memory_used":  types.Int64Type,
		"disk":         types.Int64Type,
		"disk_used":    types.Int64Type,
		"instances":    types.Int64Type,
		"status":       types.StringType,
		"cluster_name": types.StringType,
		"created_at":   types.StringType,
		"labels": types.ListType{
			ElemType: types.StringType,
		},
	}
)

type TFDedicatedServer struct {
	Region           types.String `tfsdk:"region"`
	DedicatedServers types.List   `tfsdk:"dedicated_servers"`
}

func (t *TFDedicatedServer) SetDedicatedServers(items []DedicatedServerItem) diag.Diagnostics {
	var vals []attr.Value
	for _, x := range items {
		vals = append(vals, types.ObjectValueMust(DedicatedServerItemType, map[string]attr.Value{
			"id":           x.ID,
			"name":         x.Name,
			"type_id":      x.TypeID,
			"sockets":      x.Sockets,
			"vcpus":        x.VCPUs,
			"vcpus_used":   x.VCPUsUsed,
			"memory":       x.Memory,
			"memory_used":  x.MemoryUsed,
			"disk":         x.Disk,
			"disk_used":    x.DiskUsed,
			"instances":    x.Instances,
			"status":       x.Status,
			"cluster_name": x.ClusterName,
			"created_at":   x.CreatedAt,
			"labels":       x.Labels,
		}))
	}
	l, d := types.ListValue(types.ObjectType{AttrTypes: DedicatedServerItemType}, vals)
	if d.HasError() {
		return d
	}
	t.DedicatedServers = l
	return d
}

type DedicatedServerItem struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	TypeID      types.String `tfsdk:"type_id"`
	Sockets     types.Int64  `tfsdk:"sockets"`
	VCPUs       types.Int64  `tfsdk:"vcpus"`
	VCPUsUsed   types.Int64  `tfsdk:"vcpus_used"`
	Memory      types.Int64  `tfsdk:"memory"`
	MemoryUsed  types.Int64  `tfsdk:"memory_used"`
	Disk        types.Int64  `tfsdk:"disk"`
	DiskUsed    types.Int64  `tfsdk:"disk_used"`
	Instances   types.Int64  `tfsdk:"instances"`
	Status      types.String `tfsdk:"status"`
	ClusterName types.String `tfsdk:"cluster_name"`
	CreatedAt   types.String `tfsdk:"created_at"`
	Labels      types.List   `tfsdk:"labels"`
}

func (i *DedicatedServerItem) SetLabels(ctx context.Context, labels []string) diag.Diagnostics {
	l, d := types.ListValueFrom(ctx, types.StringType, labels)
	if d.HasError() {
		return d
	}
	i.Labels = l
	return d
}