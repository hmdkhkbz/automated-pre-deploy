// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// testAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance testing. The factory function will be invoked for every Terraform
// CLI command executed to create a provider server to which the CLI can
// reattach.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"scaffolding": providerserver.NewProtocol6WithError(New("test")()),
}

var testAccProtoV6ArvanProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"arvan": providerserver.NewProtocol6WithError(NewArvanProvider("test")()),
}

func testAccPreCheck(t *testing.T) {
	if k := os.Getenv("TF_VAR_API_KEY"); k == "" {
		t.Fatal("TF_VAR_API_KEY environment variable must be set")
	}
}
