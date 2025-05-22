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

data "arvan_images" "image_list" {
  region     = var.region
  image_type = "distributions" // or one of: arvan, private
}

output "images" {
  value = data.arvan_images.image_list.distributions
}
