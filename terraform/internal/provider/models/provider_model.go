package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type ArvanProviderDataModel struct {
	ApiKey types.String `tfsdk:"api_key"`
}
