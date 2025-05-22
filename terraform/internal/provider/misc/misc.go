package misc

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/utl"
)

func RemoveResourceIfNotFound(ctx context.Context, resp *resource.ReadResponse, err error) bool {
	if respErr, ok := err.(*api.ResponseError); ok {
		if respErr.Code == 404 {
			resp.State.RemoveResource(ctx)
			return true
		}
	}
	return false
}

type Configurable interface {
	SetAPIClient(client *api.Client)
}

func ConfigureResource(ctx context.Context, req *resource.ConfigureRequest, resp *resource.ConfigureResponse, res Configurable) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		utl.ResourceConfigureError(req, resp)
		return
	}
	res.SetAPIClient(client)
}

func ConfigureDatasource(ctx context.Context, req *datasource.ConfigureRequest, resp *datasource.ConfigureResponse, ds Configurable) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		utl.DataSourceConfigureError(req, resp)
		return
	}
	ds.SetAPIClient(client)
}
