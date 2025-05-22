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

variable "chosen_distro_name" {
  type        = string
  description = " The chosen distro name for image"
  default     = "ubuntu"
}

variable "chosen_name" {
  type        = string
  description = "The chosen release for image"
  default     = "22.04"
}

variable "chosen_plan_id" {
  type        = string
  description = "The chosen ID of plan"
  default     = "g2-4-2-0"
}

variable "chosen_snapshot_id" {
  type        = string
  description = "The chosen ID of snapshot"
  default     = ""
}

data "arvan_security_groups" "default_security_groups" {
  region = var.region
}

data "arvan_images" "terraform_image" {
  region     = var.region
  image_type = "distributions" // or one of: arvan, private
}

data "arvan_plans" "plan_list" {
  region = var.region
}

locals {
  chosen_image = try([for image in data.arvan_images.terraform_image.distributions : image
  if image.distro_name == var.chosen_distro_name && image.name == var.chosen_name][0], null)
  selected_plan = try([for plan in data.arvan_plans.plan_list.plans : plan if plan.id == var.chosen_plan_id][0], null)
}

resource "arvan_volume" "terraform_volume" {
  region      = var.region
  description = "Terraform-created volume"
  name        = "tf_volume"
  size        = 9
}

data "arvan_networks" "terraform_network" {
  region = var.region
}

resource "arvan_network" "terraform_private_network" {
  region      = var.region
  description = "Terraform-created private network"
  name        = "tf_private_network"
  dhcp_range = {
    start = "10.255.255.19"
    end   = "10.255.255.150"
  }
  dns_servers    = ["8.8.8.8", "1.1.1.1"]
  enable_dhcp    = true
  enable_gateway = true
  cidr           = "10.255.255.0/24"
  gateway_ip     = "10.255.255.1"
}

resource "arvan_abrak" "built_by_terraform" {
  depends_on = [arvan_volume.terraform_volume, arvan_network.terraform_private_network]
  timeouts {
    create = "1h30m"
    update = "2h"
    delete = "20m"
    read   = "10m"
  }
  region       = var.region
  name         = "built_from_snapshot_by_terraform_${count.index + 1}"
  count        = 1
  image_id     = local.chosen_image.id
  flavor_id    = local.selected_plan.id
  disk_size    = 25
  snapshot_id  = var.chosen_snapshot_id
  enable_ipv4  = true // optional, default: true
  enable_ipv6  = true
  networks = [
    {
      network_id = arvan_network.terraform_private_network.network_id
    }
  ]
  security_groups = [data.arvan_security_groups.default_security_groups.groups[0].id]
  volumes         = [arvan_volume.terraform_volume.id]
}
