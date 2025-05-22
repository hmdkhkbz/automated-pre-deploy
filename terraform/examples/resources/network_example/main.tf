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

data "arvan_networks" "network_list" {
  region = var.region
  filters {                                  //optional
    public     = true                        //optional
    name       = "network name"              //optional
    network_id = "id of an existing network" //optional
  }
}

output "network_list" {
  value = data.arvan_networks.network_list.networks
}
