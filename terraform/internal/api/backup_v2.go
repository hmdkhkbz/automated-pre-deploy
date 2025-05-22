package api

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type S3Data struct {
	Progress   int    `json:"progress"`
	BackupID   string `json:"backup_id"`
	Bucket     string `json:"bucket"`
	Region     string `json:"region"`
	FailReason string `json:"fail_reason,omitempty"`
}

type ListBackup struct {
	Meta Meta             `json:"meta"`
	Data []BackupListData `json:"data"`
}

type BackupListData struct {
	Occupancy    int      `json:"occupancy"`
	Quota        int      `json:"quota"`
	NextBackup   string   `json:"next_backup"`
	BackupName   string   `json:"backup_name"`
	InstanceID   string   `json:"instance_id"`
	InstanceName string   `json:"instance_name"`
	Status       string   `json:"status"`
	Labels       []string `json:"labels,omitempty"`
	S3           *S3Data  `json:"s3,omitempty"`
}

type BackupDetails struct {
	Data []BackupDetailsData `json:"data"`
}

type BackupDetailsData struct {
	ProvisionedSize int     `json:"provisioned_size"`
	UsedSize        float64 `json:"used_size"`
	CreateAt        int64   `json:"created_at"`
	BackupID        string  `json:"backup_id"`
	Status          string  `json:"status"`
	SlotName        string  `json:"slot_name"`
	FailReason      string  `json:"fail_reason,omitempty"`
}

type Meta struct {
	Total int `json:"total"`
}

type ListVolumeSnapshots struct {
	Data []ListVolumeSnapshotsData `json:"data"`
}

type ListVolumeSnapshotsData struct {
	VolumeID               string `json:"volume_id"`
	VolumeName             string `json:"volume_name"`
	SnapshotCount          int    `json:"snapshots_count"`
	Status                 string `json:"status"`
	Progress               int    `json:"progress"`
	InProgressSnapshotID   string `json:"in_progress_snapshot_id"`
	InProgressSnapshotName string `json:"in_progress_snapshot_name"`
}

type SnapshotDetailsList struct {
	ID        string                `json:"id"`
	Snapshots []SnapshotDetailsData `json:"snapshots"`
}

type SnapshotDetailsData struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Size         int64    `json:"size"`
	CreatedAt    int64    `json:"created_at"`
	Status       string   `json:"status"`
	Progress     int      `json:"progress"`
	CurrentState bool     `json:"current_state"`
	Labels       []string `json:"labels"`
}

func (s *SnapshotDetailsData) GetFormattedTime() string {
	return time.Unix(s.CreatedAt/1000, 0).Format(time.RFC3339)
}

type EditSnapshotName struct {
	Name string `json:"name"`
}

type EditSnapshotNameResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message,omitempty"`
	Items   [][]string `json:"errors,omitempty"`
}

type EditSnapshotLabels struct {
	Labels []string `json:"labels"`
}

type EditSnapshotLabelsResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

type DeleteSnapshot struct {
	SnapshotIDs []string `json:"snapshot_ids"`
}

type DeleteSnapshotResponse struct {
	Code    int        `json:"code"`
	Message string     `json:"message,omitempty"`
	Items   [][]string `json:"errors,omitempty"`
}

type CreateVolumeSnapshot struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	VolumeID    string `json:"volume_id"`
}

type CreateVolumeSnapshotResponse struct {
	VolumeID   string     `json:"volume_id"`
	SnapshotID string     `json:"snapshot_id"`
	Message    string     `json:"message,omitempty"`
	Items      [][]string `json:"errors,omitempty"`
}

type CreateVolumeFromSnapshot struct {
	Name string `json:"name"`
}

type CreateVolumeFromSnapshotResponse struct {
	VolumeID   string `json:"id"`
	VolumeName string `json:"name"`
	VolumeSize int    `json:"size"`
	Code       int    `json:"code"`
	Message    string `json:"message,omitempty"`
}

type ListInstanceSnapshots struct {
	Data []ListInstanceSnapshotsData `json:"data"`
}

type ListInstanceSnapshotsData struct {
	InstanceID             string `json:"instance_id"`
	InstanceName           string `json:"instance_name"`
	SnapshotCount          int    `json:"snapshots_count"`
	Status                 string `json:"status"`
	Progress               int    `json:"progress"`
	InProgressSnapshotID   string `json:"in_progress_snapshot_id"`
	InProgressSnapshotName string `json:"in_progress_snapshot_name"`
}

type BackupV2Client struct {
	r *Requester
}

func NewBackupV2Client(r *Requester) *BackupV2Client {
	return &BackupV2Client{
		r: r,
	}
}

func (b *BackupV2Client) ListBackups(ctx context.Context, region string) (*ListBackup, error) {
	uri := fmt.Sprintf("%s/backup/%s/list", bpV2, region)
	data, err := b.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret ListBackup
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) BackupDetails(ctx context.Context, region, instanceID string) (*BackupDetails, error) {
	uri := fmt.Sprintf("%s/backup/%s/details/%s", bpV2, region, instanceID)
	data, err := b.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret BackupDetails
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) ListVolumeSnapshots(ctx context.Context, region string) (*ListVolumeSnapshots, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/volume/list", bpV2, region)
	data, err := b.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret ListVolumeSnapshots
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) VolumeSnapshotDetails(ctx context.Context, region, volumeID string) (*SnapshotDetailsList, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/volume/%s/details", bpV2, region, volumeID)
	data, err := b.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret SnapshotDetailsList
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) EditSnapshotName(ctx context.Context, region, snapshotID string, req *EditSnapshotName) (*EditSnapshotNameResponse, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/%s/name", bpV2, region, snapshotID)
	data, err := b.r.DoRequest(ctx, "PUT", uri, req)
	if err != nil {
		return nil, err
	}
	var ret EditSnapshotNameResponse
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) EditSnapshotLabels(ctx context.Context, region, snapshotID string, req *EditSnapshotLabels) (*EditSnapshotLabelsResponse, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/%s/labels", bpV2, region, snapshotID)
	data, err := b.r.DoRequest(ctx, "PUT", uri, req)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	var ret EditSnapshotLabelsResponse
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) DeleteSnapshot(ctx context.Context, region string, req *DeleteSnapshot) (*DeleteSnapshotResponse, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/delete", bpV2, region)
	data, err := b.r.DoRequest(ctx, "POST", uri, req)
	if err != nil {
		return nil, err
	}
	var ret DeleteSnapshotResponse
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) CreateVolumeSnapshot(ctx context.Context, region string, req *CreateVolumeSnapshot) (*CreateVolumeSnapshotResponse, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/volume/create", bpV2, region)
	data, err := b.r.DoRequest(ctx, "POST", uri, req)
	if err != nil {
		return nil, err
	}
	var ret CreateVolumeSnapshotResponse
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) CreateVolumeFromSnapshot(ctx context.Context, region, snapshotID string, req *CreateVolumeFromSnapshot) (*CreateVolumeFromSnapshotResponse, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/%s/create-volume", bpV2, region, snapshotID)
	data, err := b.r.DoRequest(ctx, "POST", uri, req)
	if err != nil {
		return nil, err
	}
	var ret CreateVolumeFromSnapshotResponse
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) GetVolumeSnapshot(ctx context.Context, region, volumeID, snapshotID string) (*SnapshotDetailsData, error) {
	all, err := b.VolumeSnapshotDetails(ctx, region, volumeID)
	if err != nil {
		return nil, err
	}
	for _, x := range all.Snapshots {
		if x.ID == snapshotID {
			return &x, nil
		}
	}
	return nil, &ResponseError{
		Code:    404,
		Message: "snapshot not found",
	}
}

func (b *BackupV2Client) ListInstanceSnapshots(ctx context.Context, region string) (*ListInstanceSnapshots, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/instance/list", bpV2, region)
	data, err := b.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret ListInstanceSnapshots
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

func (b *BackupV2Client) InstanceSnapshotDetails(ctx context.Context, region, instanceID string) (*SnapshotDetailsList, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/instance/%s/details", bpV2, region, instanceID)
	data, err := b.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret SnapshotDetailsList
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

type DeleteInstanceSnapshotsRequest struct {
	InstanceIDs []string `json:"instance_ids"`
}

type DeleteInstanceSnapshotsResponse struct {
	Message string     `json:"message,omitempty"`
	Items   [][]string `json:"errors,omitempty"`
}

func (b *BackupV2Client) DeleteInstanceSnapshot(ctx context.Context, region string, req *DeleteInstanceSnapshotsRequest) (*DeleteInstanceSnapshotsResponse, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/instance/delete", bpV2, region)
	data, err := b.r.DoRequest(ctx, "POST", uri, req)
	if err != nil {
		return nil, err
	}
	var ret DeleteInstanceSnapshotsResponse
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

type CreateInstanceSnapshotRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InstanceID  string `json:"instance_id"`
}

type CreateInstanceSnapshotResponse struct {
	InstanceID string     `json:"instance_id"`
	SnapshotID string     `json:"snapshot_id"`
	Message    string     `json:"message,omitempty"`
	Items      [][]string `json:"errors,omitempty"`
}

func (b *BackupV2Client) CreateInstanceSnapshot(ctx context.Context, region string, req *CreateInstanceSnapshotRequest) (*CreateInstanceSnapshotResponse, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/instance/create", bpV2, region)
	data, err := b.r.DoRequest(ctx, "POST", uri, req)
	if err != nil {
		return nil, err
	}
	var ret CreateInstanceSnapshotResponse
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}

type SnapshotDetailsResponse struct {
	Data SnapshotDetailsData `json:"data"`
}

func (b *BackupV2Client) GetSnapshotDetails(ctx context.Context, region, snapshotID string) (*SnapshotDetailsResponse, error) {
	uri := fmt.Sprintf("%s/snapshot/%s/%s", bpV2, region, snapshotID)
	data, err := b.r.DoRequest(ctx, "GET", uri, nil)
	if err != nil {
		return nil, err
	}
	var ret SnapshotDetailsResponse
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return &ret, nil
}
