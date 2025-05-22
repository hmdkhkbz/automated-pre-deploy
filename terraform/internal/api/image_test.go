package api

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func getRequester() *Requester {
	c := http.Client{
		Timeout: 1 * time.Minute,
	}
	r := &Requester{
		client: &c,
		apiKey: "apikey 82e8ceae-b333-5fc9-8de8-53961f71fa65",
	}
	return r
}

func TestGetImageList(t *testing.T) {
	r := getRequester()
	imgC := NewImageClient(r)
	ret, err := imgC.ListImages(context.Background(), "ir-thr-ba1", "distributions")
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range ret.Data {
		t.Log(x.Name)
	}
}

func TestGetInstances(t *testing.T) {
	r := getRequester()
	instC := NewInstanceClient(r)
	ret, err := instC.ListInstances(context.Background(), "ir-thr-fr1")
	if err != nil {
		t.Fatal(err)
	}
	for _, x := range ret {
		t.Log(x)
	}
}
