data "arvan_security_groups" "firewall_list" {
  region = var.region
}

data "arvan_images" "terraform_image" {
  region     = var.region
  image_type = "distributions"
}

data "arvan_plans" "plan_list" {
  region = var.region
}

data "arvan_networks" "terraform_network" {
  region = var.region
}
