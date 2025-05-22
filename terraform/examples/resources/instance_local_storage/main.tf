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
  default     = "ls2-12-4-128" // before choose flavor, check which region has this plan
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

resource "arvan_abrak" "built_by_terraform" {
  depends_on = [arvan_volume.terraform_volume]
  timeouts {
    create = "1h30m"
    update = "2h"
    delete = "20m"
    read   = "10m"
  }
  region       = var.region
  name         = "built_by_terraform_${count.index + 1}"
  # ssh_key_name = "your-sshkey-name"
  count        = 1
  image_id     = local.chosen_image.id
  flavor_id    = local.selected_plan.id
  disk_size    = 25
  enable_ipv4  = true // optional, default: true
  enable_ipv6  = true
  security_groups = [data.arvan_security_groups.default_security_groups.groups[0].id]
  volumes         = [arvan_volume.terraform_volume.id]
}
