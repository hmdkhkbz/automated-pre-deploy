package ds

import (
	"context"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/models"
	"terraform-provider-hashicups-pf/internal/utl"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type InstanceDatasource struct {
	client *api.Client
}

func (i *InstanceDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_abraks"
}

func (i *InstanceDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		utl.DataSourceConfigureError(&req, resp)
		return
	}
	i.client = client
}

func (i *InstanceDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"instances": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"flavor": schema.ObjectAttribute{
							Computed: true,
							AttributeTypes: map[string]attr.Type{
								"id":   types.StringType,
								"name": types.StringType,
							},
						},
						"status": schema.StringAttribute{
							Computed: true,
						},
						"image": schema.ObjectAttribute{
							Computed: true,
							AttributeTypes: map[string]attr.Type{
								"id":         types.StringType,
								"name":       types.StringType,
								"os":         types.StringType,
								"os_version": types.StringType,
								"metadata": types.MapType{
									ElemType: types.StringType,
								},
							},
						},
						"created": schema.StringAttribute{
							Computed: true,
						},
						"password": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"task_state": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"key_name": schema.StringAttribute{
							Computed: true,
							Optional: true,
						},
						"security_groups": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
									"description": schema.StringAttribute{
										Computed: true,
									},
									"default": schema.BoolAttribute{
										Computed: true,
									},
									"ip_addresses": schema.ListAttribute{
										Computed:    true,
										ElementType: types.StringType,
									},
									"rules": schema.ListNestedAttribute{
										Computed: true,
										NestedObject: schema.NestedAttributeObject{
											Attributes: map[string]schema.Attribute{
												"id": schema.StringAttribute{
													Computed: true,
												},
												"description": schema.StringAttribute{
													Computed: true,
												},
												"direction": schema.StringAttribute{
													Computed: true,
												},
												"ether_type": schema.StringAttribute{
													Computed: true,
												},
												"group_id": schema.StringAttribute{
													Computed: true,
												},
												"ip": schema.StringAttribute{
													Computed: true,
													Optional: true,
												},
												"port_start": schema.Int64Attribute{
													Computed: true,
												},
												"port_end": schema.Int64Attribute{
													Computed: true,
												},
												"protocol": schema.StringAttribute{
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"addresses": schema.MapAttribute{
							Computed: true,
							ElementType: types.ListType{
								ElemType: types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"mac":       types.StringType,
										"version":   types.StringType,
										"address":   types.StringType,
										"type":      types.StringType,
										"is_public": types.BoolType,
									},
								},
							},
						},
						"tags": schema.ListNestedAttribute{
							Computed: true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										Computed: true,
									},
									"name": schema.StringAttribute{
										Computed: true,
									},
								},
							},
						},
						"ha_enabled": schema.BoolAttribute{
							Computed: true,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func (i *InstanceDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var tfData models.TFInstanceDatasourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &tfData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := i.client.Instance.ListInstances(ctx, tfData.Region.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching abraks", err.Error())
		return
	}

	for _, s := range apiResp {
		tfInst := models.TFInstanceDetails{
			ID:        types.StringValue(s.ID),
			Name:      types.StringValue(s.Name),
			Status:    types.StringValue(s.Status),
			Created:   types.StringValue(s.Created),
			Password:  types.StringValue(s.Password),
			KeyName:   types.StringValue(s.KeyName),
			HAEnabled: types.BoolValue(s.HAEnabled),
		}

		if s.TaskState != nil {
			tfInst.TaskState = types.StringValue(*s.TaskState)
		}

		if s.Flavor != nil {
			tfInst.Flavor = models.TFServerFlavor{
				ID:   types.StringValue(s.Flavor.ID),
				Name: types.StringValue(s.Flavor.Name),
			}
		}
		if s.Image != nil {
			tfInst.Image = models.TFServerImage{
				ID:        types.StringValue(s.Image.ID),
				Name:      types.StringValue(s.Name),
				OS:        types.StringValue(s.Image.OS),
				OSVersion: types.StringValue(s.Image.OSVersion),
			}
			if s.Image.MetaData != nil {
				met := make(map[string]attr.Value)
				for k, v := range s.Image.MetaData {
					met[k] = types.StringValue(v)
				}
				tfMeta, diags := types.MapValue(types.StringType, met)
				resp.Diagnostics.Append(diags...)
				if resp.Diagnostics.HasError() {
					return
				}
				tfInst.Image.Metadata = tfMeta
			}
		}

		for _, sg := range s.SecurityGroups {
			tfSg := models.TFSecurityGroup{
				ID:          types.StringValue(sg.ID),
				Name:        types.StringValue(sg.Name),
				Description: types.StringValue(sg.Description),
				Default:     types.BoolValue(sg.Default),
				ReadOnly:    types.BoolValue(sg.ReadOnly),
			}
			for _, ip := range sg.IPAddresses {
				tfSg.IPAddresses = append(tfSg.IPAddresses, types.StringValue(ip))
			}
			for _, ru := range sg.Rules {
				if ru != nil {
					tfSg.Rules = append(tfSg.Rules, models.TFSecGroupRule{
						ID:          types.StringValue(ru.ID),
						Description: types.StringValue(ru.Description),
						Direction:   types.StringValue(ru.Direction),
						EtherType:   types.StringValue(ru.EtherType),
						GroupID:     types.StringValue(ru.GroupID),
						IP:          types.StringValue(ru.IP),
						PortStart:   types.Int64Value(int64(ru.PortStart)),
						PortEnd:     types.Int64Value(int64(ru.PortEnd)),
						Protocol:    types.StringValue(ru.Protocol),
					})
				}
			}
			tfInst.SecurityGroups = append(tfInst.SecurityGroups, tfSg)
		}

		tfInst.Addresses = make(map[string][]models.TFServerAddress)
		if s.Addresses != nil {
			for k, v := range s.Addresses {
				for _, addr := range v {
					tfInst.Addresses[k] = append(tfInst.Addresses[k], models.TFServerAddress{
						MAC:      types.StringValue(addr.MAC),
						Version:  types.StringValue(addr.Version),
						Addr:     types.StringValue(addr.Addr),
						Type:     types.StringValue(addr.Type),
						IsPublic: types.BoolValue(addr.IsPublic),
					})
				}

			}
		}

		for _, t := range s.Tags {
			tfInst.Tags = append(tfInst.Tags, models.TFTag{
				ID:   types.StringValue(t.ID),
				Name: types.StringValue(t.Name),
			})
		}

		tfData.Instances = append(tfData.Instances, tfInst)

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &tfData)...)

}

func NewInstanceDatasource() datasource.DataSource {
	return &InstanceDatasource{}
}
