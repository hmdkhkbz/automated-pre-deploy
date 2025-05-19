resource "arvan_network" "mgmt_net" {
  region      = var.region
  name        = "mgmt"
  description = "Management private network"
  cidr        = "192.168.88.0/24"
  enable_dhcp = true
  enable_gateway = false

  dhcp_range = {
    start = "192.168.88.10"
    end   = "192.168.88.200"
  }
  dns_servers = ["8.8.8.8", "1.1.1.1"]
}
resource "arvan_network" "internal_net" {
  region      = var.region
  name        = "internal"
  description = "Internal private network"
  cidr        = "172.16.18.0/24"
  enable_dhcp = true
  enable_gateway = false

  dhcp_range = {
    start = "172.16.18.10"
    end   = "172.16.18.200"
  }
  dns_servers = ["8.8.8.8", "1.1.1.1"]
}
resource "arvan_network" "tenant_net" {
  region      = var.region
  name        = "tenant"
  description = "Tenant private network"
  cidr        = "10.0.0.0/24"
  enable_dhcp = true
  enable_gateway = false

  dhcp_range = {
    start = "10.0.0.10"
    end   = "10.0.0.200"
  }
  dns_servers = ["8.8.8.8", "1.1.1.1"]
}
resource "arvan_network" "external_provider_net" {
  region      = var.region
  name        = "external-provider"
  description = "External provider private network"
  cidr        = "192.168.99.0/24" 
  enable_dhcp = true
  enable_gateway = false

  dhcp_range = {
    start = "192.168.99.10"
    end   = "192.168.99.200"
  }
  dns_servers = ["8.8.8.8", "1.1.1.1"]
}
resource "arvan_abrak" "controllers" {
  timeouts {
    create = "30m"
    update = "20m"
    delete = "20m"
    read   = "10m"
  }
  region     = var.region
  name       = "controller0${count.index + 1}"
  ssh_key_name = "ary"
  image_id   = local.chosen_image.id
  flavor_id  = local.controller_plan.id
  disk_size  = 60
  networks = [
    {
      network_id = local.chosen_network.network_id
    },
    {
      network_id = arvan_network.mgmt_net.network_id
    },
    {
      network_id = arvan_network.internal_net.network_id
    }
  ]
  security_groups = [local.chosen_security_group.id]
  count           = var.num_controller_instances 
}
resource "arvan_abrak" "computes" {
  timeouts {
    create = "30m"
    update = "20m"
    delete = "20m"
    read   = "10m"
  }
  region     = var.region
  name       = "compute0${count.index + 1}"
  ssh_key_name = "ary"
  image_id   = local.chosen_image.id
  flavor_id  = local.compute_plan.id
  disk_size  = 80
  networks = [
    {
      network_id = local.chosen_network.network_id
    },
    {
      network_id = arvan_network.mgmt_net.network_id
    },
    {
      network_id = arvan_network.internal_net.network_id
    },
    {
      network_id = arvan_network.tenant_net.network_id
    }
  ]
  security_groups = [local.chosen_security_group.id]
  count           = var.num_compute_instances
}
resource "arvan_abrak" "networks" {
  timeouts {
    create = "30m"
    update = "20m"
    delete = "20m"
    read   = "10m"
  }
  region     = var.region
  name       = "network0${count.index + 1}"
  ssh_key_name = "ary"
  image_id   = local.chosen_image.id
  flavor_id  = local.network_plan.id
  disk_size  = 40
  networks = [
    {
      network_id = local.chosen_network.network_id
    },
    {
      network_id = arvan_network.mgmt_net.network_id
    },
    {
      network_id = arvan_network.internal_net.network_id
    },
    {
      network_id = arvan_network.tenant_net.network_id
    },
    {
      network_id = arvan_network.external_provider_net.network_id
    }
  ]
  security_groups = [local.chosen_security_group.id]
  count           = var.num_network_instances
}
locals {
  controllers_info = [
    for c in arvan_abrak.controllers : {
      name         = c.name
      access_ip    = [for net in c.networks : net.ip if net.network_id == arvan_network.internal_net.network_id][0]
      ansible_host = [for net in c.networks : net.ip if net.network_id == arvan_network.mgmt_net.network_id][0]
    }
  ]

  computes_info = [
    for c in arvan_abrak.computes : {
      name         = c.name
      access_ip    = [for net in c.networks : net.ip if net.network_id == arvan_network.internal_net.network_id][0]
      ansible_host = [for net in c.networks : net.ip if net.network_id == arvan_network.mgmt_net.network_id][0]
    }
  ]

  networks_info = [
    for c in arvan_abrak.networks : {
      name         = c.name
      access_ip    = [for net in c.networks : net.ip if net.network_id == arvan_network.internal_net.network_id][0]
      ansible_host = [for net in c.networks : net.ip if net.network_id == arvan_network.mgmt_net.network_id][0]
    }
  ]
}
locals {
  ansible_inventory_rendered = templatefile("${path.module}/templates/inventory.tpl", {
    controllers = local.controllers_info
    computes    = local.computes_info
    networks    = local.networks_info
  })
}
resource "local_file" "ansible_inventory_yaml" {
  filename = "${path.module}/host.yaml"
  content  = local.ansible_inventory_rendered
}
