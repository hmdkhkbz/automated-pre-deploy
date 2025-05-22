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

resource "arvan_security_group" "terraform_security_group" {
  region      = var.region
  description = "Terraform-created security group"
  name        = "tf_security_group"
  rules = [
    {
      direction = "ingress"
      port_from = "12000"
      port_to   = "15000"
      protocol  = "tcp"
    },
    {
      direction = "egress"
      port_from = "1000"
      port_to   = "2000"
      protocol  = "udp"
      ip        = "192.168.0.240/32"
    }
  ]
}
