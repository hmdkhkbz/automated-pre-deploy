terraform {
  required_providers {
    arvan = {
      source = "terraform.arvancloud.ir/arvancloud/iaas"
    }
  }
}

provider "arvan" {
  api_key = "apikey"
}

variable "region" {
  type        = string
  description = "The chosen region for resources"
  default     = "ir-thr-ba1"
}

resource "arvan_server_snapshot" "terraform_server_snapshot" {
  region      = var.region
  description = "Terraform-created server snapshot"
  name        = "tf_server_snapshot"
  server_id   = "id of an existing instance"
}
