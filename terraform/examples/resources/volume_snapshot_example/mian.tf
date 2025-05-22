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

resource "arvan_volume_snapshot" "terraform_volume_snapshot" {
  region      = var.region
  description = "Terraform-created volume snapshot"
  name        = "tf_volume_snapshot"
  volume_id   = "id of an existing volume"
}
