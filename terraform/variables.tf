variable "region" {
  type        = string
  description = "The chosen region for resources"
  default     = "ir-thr-ba1"
}

variable "chosen_distro_name" {
  type        = string
  description = "The chosen distro name for image"
  default     = "ubuntu"
}

variable "chosen_name" {
  type        = string
  description = "The chosen release for image"
  default     = "22.04"
}

variable "controller_plan_id" {
  type        = string
  description = "The chosen ID for controller plan"
  default     = "g2-16-8-0"
}

variable "compute_plan_id" {
  type        = string
  description = "The chosen ID for compute node plan"
  default     = "g2-16-8-0"
}

variable "network_plan_id" {
  type        = string
  description = "The chosen ID for network node plan"
  default     = "g2-8-4-0"
}

variable "num_controller_instances" {
  type        = number
  description = "Number of controller node instances to create"
  default     = 3
}

variable "num_compute_instances" {
  type        = number
  description = "Number of compute node instances to create"
  default     = 2
}

variable "num_network_instances" {
  type        = number
  description = "Number of network node instances to create"
  default     = 2
}


variable "api_key" {
  type        = string
  description = "API key for accessing the provider"
  sensitive   = true
  default     = "apikey xxxx"
}
