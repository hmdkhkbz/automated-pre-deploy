package utl

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TerraformListStringToSlice(tfList []types.String) []string {
	var ret []string
	for _, x := range tfList {
		ret = append(ret, x.ValueString())
	}
	return ret
}

func DataSourceConfigureError(req *datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	resp.Diagnostics.AddError(
		"Unexpected data source provider data",
		fmt.Sprintf("Expected *api.Clinet, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	)
}

func ResourceConfigureError(req *resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	resp.Diagnostics.AddError(
		"Unexpected data source provider data",
		fmt.Sprintf("Expected *api.Clinet, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	)
}

func GetListDiffs(l1, l2 []types.String) []string {

	m := make(map[string]bool)
	for _, s := range l1 {
		m[s.ValueString()] = true
	}
	var diff []string
	for _, s := range l2 {
		if !m[s.ValueString()] {
			diff = append(diff, s.ValueString())
		}
	}
	return diff
}

func ListToSet(l []types.String) map[string]bool {
	m := make(map[string]bool)
	for _, s := range l {
		m[s.ValueString()] = true
	}
	return m
}

func ListGoStringToSet(l []string) map[string]bool {
	m := make(map[string]bool)
	for _, s := range l {
		m[s] = true
	}
	return m
}

func SetToList(s map[string]bool) []types.String {
	var ret []types.String
	for k, _ := range s {
		ret = append(ret, types.StringValue(k))
	}
	return ret
}

type Operation struct {
	ID string
	Do bool
}

func GetWhatToDo(plan map[string]bool, state map[string]bool) []Operation {
	var ret []Operation
	for k, _ := range plan {
		if !state[k] {
			ret = append(ret, Operation{k, true})
		}
	}

	for k, _ := range state {
		if !plan[k] {
			ret = append(ret, Operation{k, false})
		}
	}
	return ret
}

func AssignStringIfChanged(old *types.String, new string) {
	nv := types.StringValue(new)
	if old.IsNull() && new == "" {
		return
	}
	if !old.Equal(nv) {
		*old = nv
	}
}
