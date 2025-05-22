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

resource "arvan_volume" "terraform_volume" {
  region      = var.region
  description = "Terraform-created volume"
  name        = "tf_volume"
  size        = 9
}
