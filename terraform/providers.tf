terraform {
  required_providers {
    arvan = {
      source = "terraform.arvancloud.ir/arvancloud/iaas"
    }
    template = {
      source = "registry.terraform.io/hashicorp/template"
      version = "2.2.0"
    }
    local = {
      source = "registry.terraform.io/hashicorp/local"
      version = "2.5.2"
    }
  }
}

provider "arvan" {
  api_key = "apikey xxxx"
}
