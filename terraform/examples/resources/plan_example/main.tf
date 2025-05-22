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

data "arvan_plans" "plan_list" {
  region = var.region
}

output "plans" {
  value = data.arvan_plans.plan_list.plans
}
