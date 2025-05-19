locals {
  chosen_network = try([for network in tolist(data.arvan_networks.terraform_network.networks) : network][0], null)
  
  chosen_image = try([for image in data.arvan_images.terraform_image.distributions : image
    if image.distro_name == var.chosen_distro_name && image.name == var.chosen_name][0], null)

  controller_plan = try([for plan in data.arvan_plans.plan_list.plans : plan
    if plan.id == var.controller_plan_id][0], null)

  compute_plan = try([for plan in data.arvan_plans.plan_list.plans : plan
    if plan.id == var.compute_plan_id][0], null)

  network_plan = try([for plan in data.arvan_plans.plan_list.plans : plan
    if plan.id == var.network_plan_id][0], null)

  chosen_security_group = try([for sg in data.arvan_security_groups.firewall_list.groups : sg
    if sg.name == "default"][0], null)

}
