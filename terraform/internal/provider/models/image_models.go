package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type ImageItem struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	DistroName types.String `tfsdk:"distro_name"`
	Disk       types.Int64  `tfsdk:"disk"`
	Ram        types.Int64  `tfsdk:"ram"`
	SSHKey     types.Bool   `tfsdk:"ssh_key"`
	Password   types.Bool   `tfsdk:"password"`
}

type ImageDistroListModel struct {
	ImgType       types.String `tfsdk:"image_type"`
	Region        types.String `tfsdk:"region"`
	Distributions []ImageItem  `tfsdk:"distributions"`
}
