package api

import (
	"context"
	"encoding/json"
	"fmt"
)

type SSHKey struct {
	Name      string `json:"name"`
	PublicKey string `json:"public_key"`
}

type SSHKeyClient struct {
	r *Requester
}

func NewSSHKeyClient(r *Requester) *SSHKeyClient {
	return &SSHKeyClient{
		r: r,
	}
}

func (s *SSHKeyClient) GetSSHKeys(ctx context.Context, region string) ([]*SSHKey, error) {
	type response struct {
		Data []*SSHKey `json:"data"`
	}
	url := fmt.Sprintf("%s/%s/ssh", basePathV2, region)
	data, err := s.r.DoRequest(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	var ret response
	err = json.Unmarshal(data, &ret)
	if err != nil {
		return nil, err
	}
	return ret.Data, nil

}
