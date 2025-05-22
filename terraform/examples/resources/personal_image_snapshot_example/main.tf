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

resource "arvan_snapshot_image" "terraform_personal_image" {
  region      = var.region
  description = "Terraform-created personal image"
  name        = "tf_personal_image"
  snapshot_id = "id of an exsiting snapshot"
}
