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

data "arvan_server_groups" "server_group_list" {
  region = var.region
}

output "server_group_list" {
  value = data.arvan_server_groups.server_group_list.server_groups
}
