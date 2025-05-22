package rs

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"terraform-provider-hashicups-pf/internal/api"
	"terraform-provider-hashicups-pf/internal/provider/misc"
	"terraform-provider-hashicups-pf/internal/provider/models"
	"terraform-provider-hashicups-pf/internal/utl"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type InstanceResource struct {
	client *api.Client
}

func (i *InstanceResource) SetAPIClient(c *api.Client) {
	i.client = c
}

func (i *InstanceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_abrak"
}

func (i *InstanceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	misc.ConfigureResource(ctx, &req, resp, i)
}

func (i *InstanceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	resp.Schema = schema.Schema{
		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Create: true,
				Read:   true,
				Update: true,
				Delete: true,
			}),
		},
		Attributes: map[string]schema.Attribute{
			"region": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"task_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"image_id": schema.StringAttribute{
				Required: true,
			},
			"enable_ipv4": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(true),
			},
			"enable_ipv6": schema.BoolAttribute{
				Optional: true,
			},
			"networks": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					PlanModifiers: []planmodifier.Object{
						misc.RemovePlanModifier{},
					},
					Attributes: map[string]schema.Attribute{
						"subnet_id": schema.StringAttribute{
							Computed: true,
							Optional: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"ip": schema.StringAttribute{
							Optional: true,
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"network_id": schema.StringAttribute{
							Required: true,
						},
						"port_id": schema.StringAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"is_public": schema.BoolAttribute{
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"port_security_enabled": schema.BoolAttribute{
							Optional: true,
							Default:  booldefault.StaticBool(true),
							Computed: true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"flavor_id": schema.StringAttribute{
				Required: true,
			},
			"security_groups": schema.SetAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"disk_size": schema.Int64Attribute{
				Required:   true,
				Validators: []validator.Int64{int64validator.AtLeast(25)},
			},
			"init_script": schema.StringAttribute{
				Optional: true,
			},
			"volumes": schema.SetAttribute{
				Optional: true,
				//Computed:    true,
				ElementType: types.StringType,
			},
			"ssh_key_name": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"floating_ip": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"floating_ip_id": schema.StringAttribute{
						Required: true,
					},
					"network_id": schema.StringAttribute{
						Required: true,
						Validators: []validator.String{
							misc.NewRequiredToBeEqualToAtLeastOneOfAList(path.Root("networks"), func(networks []models.TFNetworkAttachment) map[string]bool {
								m := make(map[string]bool)
								for _, x := range networks {
									m[x.NetworkID.ValueString()] = true
								}
								return m
							}),
						},
					},
				},
			},
			"revert_to": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"snapshot_id": schema.StringAttribute{
						Required: true,
					},
				},
			},
			"server_group_id": schema.StringAttribute{
				Optional: true,
			},
			"dedicated_server_id": schema.StringAttribute{
				Optional: true,
			},
			"cluster_id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"snapshot_id": schema.StringAttribute{
				Optional: true,
			},

		},
	}
}

func (i *InstanceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.TFInstanceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.Snapshot.IsNull() {
		resp.Diagnostics.AddError("invalid operation", "revert_to attribute can only be set after instance creation")
		return
	}

	apiCreateReq := api.InstanceCreateRequest{
		Name:              data.Name.ValueString(),
		Count:             1,
		ImageID:           data.ImageID.ValueString(),
		FlavorID:          data.FlavorID.ValueString(),
		SSHKey:            !data.SSHKeyName.IsNull(),
		DiskSize:          int(data.DiskSize.ValueInt64()),
		OSVolumeID:        "",
		InitScript:        data.InitScript.ValueString(),
		DedicatedServerID: data.DedicatedServerID.ValueString(),
		EnableIPv4:        data.EnableIPv4.ValueBool(),
	}

	if !data.ServerGroupID.IsNull() {
		apiCreateReq.ServerGroupID = data.ServerGroupID.ValueString()
	}

	if !data.DedicatedServerID.IsNull() {
		apiCreateReq.DedicatedServerID = data.DedicatedServerID.ValueString()
	}

	if !data.SnapshotID.IsNull() {
		apiCreateReq.SnapshotID = data.SnapshotID.ValueString()
	}

	if !data.EnableIPv6.IsNull() {
		apiCreateReq.EnableIPv6 = data.EnableIPv6.ValueBool()
	}

	tfNets, d := data.GetNetworkAttachments(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	for _, n := range tfNets {
		apiCreateReq.NetworkIDs = append(apiCreateReq.NetworkIDs, n.NetworkID.ValueString())
	}

	if apiCreateReq.SSHKey {
		apiCreateReq.KeyName = data.SSHKeyName.ValueString()
	}

	sgs, d := data.GetSecurityGroups(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, sg := range sgs {
		apiCreateReq.SecurityGroups = append(apiCreateReq.SecurityGroups, api.SecGroupName{
			Name: sg,
		})
	}

	createTimeout, diags := data.Timeouts.Create(ctx, 20*time.Minute)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()
	apiResp, err := i.client.Instance.CreateInstanceAsync(ctx, data.Region.ValueString(), &apiCreateReq)
	if err != nil {
		resp.Diagnostics.AddError("error creating instance", err.Error())
		return
	}
	data.ID = types.StringValue(apiResp.Data.ID)
	data.TaskID = types.StringValue(apiResp.Data.TaskID)
	data.Password = types.StringValue(apiResp.Data.Password)
	data.Status = types.StringValue(apiResp.Data.Status)

	err = i.client.WaitForCondition(ctx, createTimeout, func() (bool, error) {
		var detail *api.ServerDetail
		var err error

		if data.TaskID.ValueString() != "" {
			detail, err = i.client.Instance.InquiryInstance(ctx, data.Region.ValueString(), data.TaskID.ValueString())
		} else {
			detail, err = i.client.Instance.GetInstance(ctx, data.Region.ValueString(), data.ID.ValueString())
		}


		if err != nil {
			return false, err
		}
		tflog.Info(ctx, "STATUS", map[string]interface{}{"STATUS": detail.Status})
		if detail.Status == "ACTIVE" {
			data.Status = types.StringValue(detail.Status)
			return true, nil
		}

		if detail.Status == "ERROR" {
			data.Status = types.StringValue(detail.Status)
			return false, errors.New("instance status transitioned into invalid state ERROR")
		}

		data.ID = types.StringValue(detail.ID)
		data.ClusterID = types.StringValue(detail.ClusterID)
		return false, nil
	})

	if err != nil {
		tflog.Info(ctx, "STATUS", map[string]interface{}{"err": err})
		resp.Diagnostics.AddError("instance power on error", err.Error())
		return
	}

	var networkIds []string
	for idx := range tfNets {
		networkIds = append(networkIds, tfNets[idx].NetworkID.ValueString())
	}

	err = i.client.WaitForCondition(ctx, createTimeout, func() (bool, error) {
		attachments, err := i.client.FillNetworkData(ctx, networkIds, data.Region.ValueString(), data.ID.ValueString())
		if err != nil {
			return false, err
		}

		tflog.Info(ctx, "ATTACHMENT_COUNT", map[string]interface{}{"COUNT": len(attachments)})
		var conditions []bool
		for _, v := range attachments {
			conditions = append(conditions, v.PortID != "" && v.SubnetID != "")
		}

		for _, x := range conditions {
			if !x {
				return false, nil
			}
		}
		if len(attachments) != len(tfNets) {
			return false, nil
		}

		for idx := 0; idx < len(tfNets); idx++ {
			if a, ok := attachments[tfNets[idx].NetworkID.ValueString()]; ok {
				tflog.Info(ctx, "ATTACHMENT", map[string]interface{}{"AT": a})
				tfNets[idx].IP = types.StringValue(a.IP)
				tfNets[idx].PortID = types.StringValue(a.PortID)
				tfNets[idx].SubnetID = types.StringValue(a.SubnetID)
				tfNets[idx].IsPublic = types.BoolValue(a.IsPublic)
				tfNets[idx].PortSecurityEnabled = types.BoolValue(a.PortSecurityEnabled)

			}

		}
		return true, nil
	})
	if err != nil {
		resp.Diagnostics.AddError("waiting for condition failed", err.Error())
		resp.Diagnostics.Append(d...)
		return
	}

	d = data.SetNetworkAttachments(ctx, tfNets)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	volIds, d := data.GetVolumes(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, v := range volIds {
		_, err = i.client.Volume.AttachVolume(ctx, data.Region.ValueString(), &api.VolumeAttachDetach{
			ServerID: data.ID.ValueString(),
			VolumeID: v,
		})
		if err != nil {
			resp.Diagnostics.AddError("error attaching volume", err.Error())
			return
		}

	}

	if !data.FloatingIP.IsNull() {
		att, d := data.GetFloatingIPAttachment(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, n := range tfNets {
			if n.NetworkID.Equal(att.NetworkID) {
				info, err := i.client.FIPClient.GetServerIPInfo(ctx, data.Region.ValueString())
				if err != nil {
					resp.Diagnostics.AddError("error fetching server ip info", err.Error())
					return
				}
				sInfo, ok := info[data.ID.ValueString()]
				if !ok {
					resp.Diagnostics.AddError("unexpected error", "server ip info no found")
					return
				}
				if sInfo.HasPublicIP {
					resp.Diagnostics.AddError("invalid operation", "floating ip can only be attached to servers without public ip")
					return
				}
				err = i.client.FIPClient.AttachFloatingIP(ctx, data.Region.ValueString(), att.FloatingID.ValueString(), &api.AttachReq{
					ServerID: data.ID.ValueString(),
					SubnetID: n.SubnetID.ValueString(),
					PortID:   n.PortID.ValueString(),
				})
				if err != nil {
					resp.Diagnostics.AddError("error attaching floating ip", err.Error())
					return
				}
			}
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (i *InstanceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.TFInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiResp, err := i.client.Instance.GetInstance(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok {
			if respErr.Code == 404 {
				resp.State.RemoveResource(ctx)
				return
			}
		}
		resp.Diagnostics.AddError("error fetching instance", err.Error())
		return
	}
	data.Status = types.StringValue(apiResp.Status)
	data.Name = types.StringValue(apiResp.Name)
	data.ClusterID = types.StringValue(apiResp.ClusterID)

	if apiResp.DedicatedServerID != "" {
		data.DedicatedServerID = types.StringValue(apiResp.DedicatedServerID)
	}

	if apiResp.Image != nil {
		data.ImageID = types.StringValue(apiResp.Image.ID)
	}
	if apiResp.Flavor != nil {
		data.FlavorID = types.StringValue(apiResp.Flavor.ID)
	}
	
	if len(apiResp.SecurityGroups) > 0 {
		var sgIds = make(map[string]bool)
		for _, sg := range apiResp.SecurityGroups {
			sgIds[sg.ID] = true
		}
		resp.Diagnostics.Append(data.SetSecurityGroups(ctx, sgIds)...)
		if resp.Diagnostics.HasError() {
			return
		}

	}

	attachments, err := i.client.GetNetworkAttachments(ctx, apiResp, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error getting network attachments", err.Error())
		return
	}

	var tfNets []models.TFNetworkAttachment
	for _, a := range attachments {
		tfX := models.TFNetworkAttachment{
			IP:                  types.StringValue(a.IP),
			SubnetID:            types.StringValue(a.SubnetID),
			NetworkID:           types.StringValue(a.NetworkID),
			PortID:              types.StringValue(a.PortID),
			IsPublic:            types.BoolValue(a.IsPublic),
			PortSecurityEnabled: types.BoolValue(a.PortSecurityEnabled),
		}
		/*if !a.PortSecurityEnabled {
			tfX.PortSecurityEnabled = types.BoolNull()
		}*/

		tfNets = append(tfNets, tfX)
	}

	oldAttachments, d := data.GetNetworkAttachments(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	var toAdd []models.TFNetworkAttachment
	var existing []models.TFNetworkAttachment

	for _, x := range tfNets {
		var e bool
		for _, y := range oldAttachments {
			if x.NetworkID.Equal(y.NetworkID) {
				existing = append(existing, y)
				e = true
				break
			}

		}
		if !e {
			toAdd = append(toAdd, x)
		}

	}
	existing = append(existing, toAdd...)

	d = data.SetNetworkAttachments(ctx, oldAttachments)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	serverVolumes, err := i.client.Volume.GetServerVolumes(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error fetching server volumes", err.Error())
		return
	}

	d = data.SetVolumes(ctx, serverVolumes)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	fipInfo, err := i.client.GetServerFloatingIPInfo(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			resp.Diagnostics.Append(data.SetFloatingIPAttachment(ctx, nil)...)
			if resp.Diagnostics.HasError() {
				return
			}
			resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
			return
		}
		resp.Diagnostics.AddError("error fetching floating ip", err.Error())
		return
	}
	d = data.SetFloatingIPAttachment(ctx, &models.TFFloatingIPAttachment{
		FloatingID: types.StringValue(fipInfo.ID),
		NetworkID:  types.StringValue(fipInfo.PrivateNetworkID),
	})
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (i *InstanceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.TFInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	err := i.client.Instance.DeleteInstance(ctx, data.Region.ValueString(), data.ID.ValueString())
	if err != nil {
		if respErr, ok := err.(*api.ResponseError); ok && respErr.Code == 404 {
			return
		}
		resp.Diagnostics.AddError("error deleting instance", err.Error())
		return
	}

	deleteTimeout, d := data.Timeouts.Delete(ctx, time.Minute*2)
	resp.Diagnostics.Append(d...)

	err = i.client.WaitForCondition(ctx, deleteTimeout, func() (bool, error) {
		_, err := i.client.Instance.GetInstance(ctx, data.Region.ValueString(), data.ID.ValueString())
		if err != nil && err.(*api.ResponseError).Code == 404 {
			return true, nil
		}
		
		if err != nil {
			return false, err
		}

		return false, nil
	})
	if err != nil {
		resp.Diagnostics.AddError("error deleting instance", err.Error())
	}
}

func (i *InstanceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "UPDATING")
	var planData models.TFInstanceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var stateData models.TFInstanceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "PLAN_NETWORK", map[string]interface{}{"NETWORK": planData.Networks})

	planData.Status = stateData.Status

	i.handleSimpleChanges(ctx, &stateData, &planData, resp)
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.Append(resp.State.Set(ctx, &stateData)...)
		return
	}

	i.handleRename(ctx, &stateData, &planData, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	i.handleSecurityGroups(ctx, &stateData, &planData, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	i.handleFlavorResize(ctx, &stateData, &planData, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	i.handleVolumeAttachments(ctx, &stateData, &planData, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	i.handleNetworkAttachments(ctx, &stateData, &planData, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	i.handleRootVolumeResize(ctx, &stateData, &planData, resp)

	if resp.Diagnostics.HasError() {
		return
	}

	i.handleFloatingIP(ctx, &stateData, &planData, resp)

	if resp.Diagnostics.HasError() {
		return

	}

	i.handleRevertToSnapshot(ctx, &stateData, &planData, resp)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planData)...)

}

func (i *InstanceResource) handleRevertToSnapshot(ctx context.Context, stateData, planData *models.TFInstanceResourceModel, resp *resource.UpdateResponse) {
	if !stateData.Snapshot.Equal(planData.Snapshot) {
		snp, d := planData.GetSnapshot(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		if snp != nil {
			err := i.client.SnapshotClient.Revert(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), snp.SnapshotID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("error reverting instance to snapshot", err.Error())
				return
			}

		}
	}
}

func (i *InstanceResource) handleSimpleChanges(ctx context.Context, stateData, planData *models.TFInstanceResourceModel, resp *resource.UpdateResponse) {
	const unsupportedOperation = "unsupported operation"

	if !planData.Status.Equal(stateData.Status) {
		resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("status"), planData.Status.ValueString())...)
	}

	if !planData.ImageID.Equal(stateData.ImageID) {
		resp.Diagnostics.AddError(unsupportedOperation, "image id can not be changed")
	}

	if strings.HasPrefix(stateData.FlavorID.ValueString(), "ls") && !planData.FlavorID.Equal(stateData.FlavorID) {
		resp.Diagnostics.AddError(unsupportedOperation, "instances with local storage can not be resized")
	}

	if !planData.InitScript.Equal(stateData.InitScript) {
		resp.Diagnostics.AddError(unsupportedOperation, "init script can only be set at creation time")
	}

	if !planData.SSHKeyName.Equal(stateData.SSHKeyName) {
		resp.Diagnostics.AddError(unsupportedOperation, "ssh key name can only be set at creation time")
	}

	if !planData.Region.Equal(stateData.Region) {
		resp.Diagnostics.AddError(unsupportedOperation, "region can only be set at creation time")
	}

	if !planData.ServerGroupID.Equal(stateData.ServerGroupID) {
		resp.Diagnostics.AddError(unsupportedOperation, "server group id can only be set at creation time")
	}

	if !planData.DedicatedServerID.IsNull() && !planData.DedicatedServerID.Equal(stateData.DedicatedServerID) {
		resp.Diagnostics.AddError(unsupportedOperation, "dedicated server id can only be set at creation time")
	}

	if !planData.EnableIPv4.Equal(stateData.EnableIPv4) {
		resp.Diagnostics.AddError(unsupportedOperation, "setting public IPs can only happen at creation time")
	}

	if !planData.EnableIPv6.IsNull() && !planData.EnableIPv6.Equal(stateData.EnableIPv6) {
		resp.Diagnostics.AddError(unsupportedOperation, "setting public IPs can only happen at creation time")
	}

	return
}

func (i *InstanceResource) handleRename(ctx context.Context, stateData, planData *models.TFInstanceResourceModel, resp *resource.UpdateResponse) {
	if !stateData.Name.Equal(planData.Name) {
		err := i.client.Instance.RenameInstance(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), planData.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error renaming instance", err.Error())
			return
		}
	}
}

func (i *InstanceResource) handleFlavorResize(ctx context.Context, stateData, planData *models.TFInstanceResourceModel, resp *resource.UpdateResponse) {
	if planData.FlavorID.Equal(stateData.FlavorID) {
		return
	}

	resp.Diagnostics.AddWarning("instance power off", "during resize operation your instance powers off")
	err := i.client.Instance.PowerOffInstance(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error powering off instance", err.Error())
		return
	}

	updateTimeout, d := planData.Timeouts.Update(ctx, 5*time.Minute)
	resp.Diagnostics.Append(d...)

	err = i.client.WaitForCondition(ctx, updateTimeout, func() (bool, error) {
		det, err := i.client.Instance.GetInstance(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
		if err != nil {
			return false, err
		}
		if det.Status == "SHUTOFF" {
			planData.Status = types.StringValue(det.Status)
			return true, nil
		}
		return false, nil
	})
	if err != nil {
		resp.Diagnostics.AddError("error powering off instance", err.Error())
		return
	}

	err = i.client.Instance.ResizeInstance(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), planData.FlavorID.String())
	if err != nil {
		resp.Diagnostics.AddError("error resizing instance", err.Error())
		return
	}

	err = i.client.Instance.PowerOnInstance(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("error powering on instance", err.Error())
		return
	}

	err = i.client.WaitForCondition(ctx, updateTimeout, func() (bool, error) {
		det, err := i.client.Instance.GetInstance(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
		if err != nil {
			return false, err
		}
		if det.Status == "ACTIVE" {
			planData.Status = types.StringValue(det.Status)
			return true, nil
		}

		if det.Status == "ERROR" {
			planData.Status = types.StringValue(det.Status)
			return false, errors.New("instance state transitioned to ERROR")
		}

		return false, nil
	})
	if err != nil {
		resp.Diagnostics.AddError("instance power on failed", err.Error())
		return
	}

}

func (i *InstanceResource) handleVolumeAttachments(ctx context.Context, stateData, planData *models.TFInstanceResourceModel, resp *resource.UpdateResponse) {
	planVols, d := planData.GetVolumes(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateVols, d := stateData.GetVolumes(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !planData.Volumes.Equal(stateData.Volumes) {

		stateVolSet := utl.ListGoStringToSet(stateVols)
		tflog.Info(ctx, "VOLUME_STATE", map[string]interface{}{"VOLUMES": stateVolSet})
		planVolSet := utl.ListGoStringToSet(planVols)

		tflog.Info(ctx, "VOLUME_PLAN", map[string]interface{}{"VOLUMES": planVolSet})

		wtd := utl.GetWhatToDo(planVolSet, stateVolSet)

		for _, w := range wtd {
			if w.Do {
				_, err := i.client.Volume.AttachVolume(ctx, planData.Region.ValueString(), &api.VolumeAttachDetach{
					ServerID: planData.ID.ValueString(),
					VolumeID: w.ID,
				})
				if err != nil {
					resp.Diagnostics.AddError("error attaching volume", err.Error())
					return
				}
				continue
			}
			err := i.client.Volume.DetachVolume(ctx, planData.Region.ValueString(), &api.VolumeAttachDetach{
				ServerID: planData.ID.ValueString(),
				VolumeID: w.ID,
			})
			if err != nil {
				resp.Diagnostics.AddError("error detaching volume", err.Error())
				return
			}

		}
	}
}

func (i *InstanceResource) handleNetworkAttachments(ctx context.Context, stateData, planData *models.TFInstanceResourceModel, resp *resource.UpdateResponse) {

	tfPlanNets, d := planData.GetNetworkAttachments(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	tfStateNets, d := stateData.GetNetworkAttachments(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, stateNet := range tfStateNets {
		if attachment, _ := planData.GetNetworkAttachment(ctx, stateNet.NetworkID.ValueString()); attachment == nil {
			tflog.Warn(ctx, "detaching network non-existent in plan", map[string]interface{}{"NETWORK_ID": stateNet.NetworkID.ValueString()})
			err := i.client.Subnet.DetachServerFromNetwork(ctx, stateData.Region.ValueString(), stateNet.PortID.ValueString(), stateData.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("error detaching server from network", err.Error())
				return
			}
		}
	}

	newNetStates := make([]models.TFNetworkAttachment, 0)
	for _, planNet := range tfPlanNets {
		if ok, _ := stateData.HasEqualNetworkAttachment(ctx, planNet); ok {
			tflog.Warn(ctx, fmt.Sprintf("network %s needs no changes", planNet.NetworkID.ValueString()))
			newNetStates = append(newNetStates, planNet)
			continue
		}

		if currentAttachment, _ := stateData.GetNetworkAttachment(ctx, planNet.NetworkID.ValueString()); currentAttachment != nil {
			tflog.Warn(ctx, fmt.Sprintf("network %s exists but needs changes, detaching", planNet.NetworkID.ValueString()))

			tflog.Warn(ctx, "NETWORK_DETACH", map[string]interface{}{"NETWORK_ID": planNet.NetworkID.ValueString()})
			err := i.client.Subnet.DetachServerFromNetwork(ctx, stateData.Region.ValueString(), currentAttachment.PortID.ValueString(), stateData.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("error detaching server from network", err.Error())
				return
			}
		}

		tflog.Warn(ctx, "NETWORK_ATTACH", map[string]interface{}{"NETWORK_ID": planNet.NetworkID.ValueString()})
		s, err := i.client.Subnet.GetNetworkSubnet(ctx, stateData.Region.ValueString(), planNet.NetworkID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error fetching subnet", err.Error())
			return
		}

		req := &api.AttachServerToNetworkRequest{
			ServerID:           stateData.ID.ValueString(),
			SubnetID:           s.ID,
			IP:                 planNet.IP.ValueString(),
			EnablePortSecurity: planNet.PortSecurityEnabled.ValueBool(),
		}

		apiResp, err := i.client.Subnet.AttachServerToNetwork(ctx, stateData.Region.ValueString(), planNet.NetworkID.ValueString(), req)
		if err != nil {
			resp.Diagnostics.AddError("error attaching network", err.Error())
			return
		}

		newNetStates = append(newNetStates, models.TFNetworkAttachment{
			PortID:              types.StringValue(apiResp.ID),
			IP:                  types.StringValue(apiResp.IPAddress),
			SubnetID:            types.StringValue(apiResp.SubnetID),
			NetworkID:           types.StringValue(apiResp.NetworkID),
			PortSecurityEnabled: planNet.PortSecurityEnabled,
			IsPublic:            types.BoolValue(false),
		})
	}

	// keep the networks null value instead of empty list
	if len(newNetStates) != 0 {
		resp.Diagnostics.Append(planData.SetNetworkAttachments(ctx, newNetStates)...)
	}
}

func (i *InstanceResource) handleRootVolumeResize(ctx context.Context, stateData, planData *models.TFInstanceResourceModel, resp *resource.UpdateResponse) {
	if !planData.DiskSize.Equal(stateData.DiskSize) {
		if strings.HasPrefix(planData.FlavorID.ValueString(), "ls") {
			resp.Diagnostics.AddError("invalid operation", "local storage instance root volume can not be resized")
			return
		}
		if planData.DiskSize.ValueInt64() < stateData.DiskSize.ValueInt64() {
			resp.Diagnostics.AddError("invalid operation", "new root volume size must be greater than old root volume disk size")
			return
		}

		if stateData.Status.ValueString() == "ACTIVE" {
			resp.Diagnostics.AddWarning("instance power off", "during root volume resize operation your instance powers off")
			err := i.client.Instance.PowerOffInstance(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError("error powering off instance", err.Error())
				return
			}
		}

		updateTimeout, d := planData.Timeouts.Update(ctx, 5*time.Minute)
		resp.Diagnostics.Append(d...)

		err := i.client.WaitForCondition(ctx, updateTimeout, func() (bool, error) {
			det, err := i.client.Instance.GetInstance(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
			if err != nil {
				return false, err
			}
			if det.Status == "SHUTOFF" {
				planData.Status = types.StringValue(det.Status)
				return true, nil
			}

			return false, nil
		})
		if err != nil {
			resp.Diagnostics.AddError("error powering off instance", err.Error())
			return
		}

		err = i.client.Instance.ResizeRootVolume(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), planData.DiskSize.ValueInt64())
		if err != nil {
			resp.Diagnostics.AddError("error resizing root volume", err.Error())
			return
		}

		err = i.client.WaitForCondition(ctx, updateTimeout, func() (bool, error) {
			det, err := i.client.Instance.GetInstance(ctx, stateData.Region.ValueString(), stateData.ID.ValueString())
			if err != nil {
				return false, err
			}

			if det.Status == "ACTIVE" {
				planData.Status = types.StringValue(det.Status)
				return true, nil
			}

			return false, nil
		})
		if err != nil {
			resp.Diagnostics.AddError("error", "instance power on took too long")
			return
		}

	}

}

func (i *InstanceResource) handleSecurityGroups(ctx context.Context, stateData, planData *models.TFInstanceResourceModel, resp *resource.UpdateResponse) {
	stateSGIDs, d := stateData.GetSecurityGroups(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	planSGIDs, d := planData.GetSecurityGroups(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !planData.SecurityGroups.Equal(stateData.SecurityGroups) {
		stateSgSet := utl.ListGoStringToSet(stateSGIDs)
		planSgSet := utl.ListGoStringToSet(planSGIDs)

		whtd := utl.GetWhatToDo(planSgSet, stateSgSet)
		for _, w := range whtd {
			if w.Do {
				err := i.client.Firewall.AddServerToGroup(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), w.ID)
				if err != nil {
					resp.Diagnostics.AddError("error adding server to security group", err.Error())
					return
				}
				continue
			}
			err := i.client.Firewall.RemoveServerFromGroup(ctx, stateData.Region.ValueString(), stateData.ID.ValueString(), w.ID)
			if err != nil {
				resp.Diagnostics.AddError("error removing server from security group", err.Error())
				return
			}
		}
	}
}

func (i *InstanceResource) handleFloatingIP(ctx context.Context, stateData, planData *models.TFInstanceResourceModel, resp *resource.UpdateResponse) {
	if planData.FloatingIP.Equal(stateData.FloatingIP) {
		return
	}

	tfStateNets, d := stateData.GetNetworkAttachments(ctx)
	resp.Diagnostics.Append(d...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !stateData.FloatingIP.IsNull() {
		var portID string
		stateFip, d := stateData.GetFloatingIPAttachment(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, n := range tfStateNets {
			if n.NetworkID.Equal(stateFip.NetworkID) {
				portID = n.PortID.ValueString()
			}
			break
		}

		err := i.client.FIPClient.DetachFloatingIP(ctx, stateData.Region.ValueString(), portID)
		if err != nil {
			resp.Diagnostics.AddError("error detaching floating ip", err.Error())
			return
		}
	}

	if !planData.FloatingIP.IsNull() {
		tfPlanNets, d := planData.GetNetworkAttachments(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		var planPortID string
		var planSubnetID string
		planFip, d := planData.GetFloatingIPAttachment(ctx)
		resp.Diagnostics.Append(d...)
		if resp.Diagnostics.HasError() {
			return
		}
		for _, n := range tfPlanNets {
			if n.NetworkID.Equal(planFip.NetworkID) {
				planPortID = n.PortID.ValueString()
				planSubnetID = n.SubnetID.ValueString()
				break
			}
		}

		ips, err := i.client.FIPClient.GetServerIPInfo(ctx, stateData.Region.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("error fetching server ip info", err.Error())
			return
		}

		if info, ok := ips[stateData.ID.ValueString()]; ok {
			if info.HasPublicIP {
				resp.Diagnostics.AddError("unsupported operation", "floating ip can only be attached to instances without public ip address")
				return
			}

		}

		err = i.client.FIPClient.AttachFloatingIP(ctx, stateData.Region.ValueString(), planFip.FloatingID.ValueString(), &api.AttachReq{
			ServerID: stateData.ID.ValueString(),
			SubnetID: planSubnetID,
			PortID:   planPortID,
		})
		if err != nil {
			resp.Diagnostics.AddError("error attaching floating ip to instance", err.Error())
			return
		}
	}

}

func NewInstanceResource() resource.Resource {
	return &InstanceResource{}
}
