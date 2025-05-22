package api

import (
	"context"
	"testing"
)

func TestCreateVolume(t *testing.T) {
	c := VolumeClient{r: getRequester()}
	resp, err := c.CreateVolume(context.Background(), "ir-thr-fr1", &ServerVolume{
		Name:        "test_api_sdk",
		Description: "asdasd",
		Size:        9,
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func TestGetAttachments(t *testing.T) {
	c := NewClient("apikey d48e00ad-75a5-5ed7-b988-e0415693f823")

	netsa, err := c.Subnet.GetAllNetworks(context.Background(), "ir-thr-fr1")
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range netsa {

		if v.ID == "f8c08f20-d156-4bcd-bab2-bedb13146db6" {
			if len(v.Subnets) > 0 {
				for _, x := range v.Subnets[0].Servers {
					for _, y := range x.IPs {
						//t.Log(v.Name, y.IP, y.PortID, y.SubnetID, y.SubnetID, v.Subnets[0].SubnetID)
						if y.SubnetID == v.Subnets[0].ID && x.ID == "175f312b-06e5-41b7-93bc-5bc4cc6d44e1" {
							t.Log(x.Name, x.ID, y.IP, y.PortID, y.SubnetID, y.SubnetID, v.Subnets[0].ID)
						}

					}
				}
			}
		}

	}
}
