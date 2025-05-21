# Abrak Infrastructure on Arvan IaaS with Terraform

This repository contains the Terraform configuration for deploying and managing Abraks on Arvan IaaS.

## Overview

This configuration automates the creation and management of the following resources on Arvan IaaS:

* **Abraks** Virtual Machines (Instances)
* **Private Networks** Virtual Networks and Subnets
* **IaaS Firewalls** Security Groups (Firewall Rules)
* **Disks of Abrak** Volumes (Storage)

## Prerequisites

Before using this configuration, ensure you have the following:

* **Terraform CLI** installed on your local machine (version >= X.Y.Z). You can find installation instructions [here](https://www.terraform.io/downloads.html).
* An **Arvan IaaS account** with the necessary permissions to create and manage resources.
* Your **Arvan API credentials** configured for Terraform. This typically involves setting environment variables or using a provider configuration file. Refer to the [Arvan Terraform Provider documentation](<Arvan Provider Documentation Link - Replace with actual link>) for details.
* pre-installed **Arvan Terraform** source codes: https://git.arvancloud.ir/arvancloud/iaas/terraform-provider
## Getting Started

1.  **Clone Arvancloud IaaS terraform:**
    before run tf files, you should have installed Arvan IaaS Terraform: 
     https://git.arvancloud.ir/arvancloud/iaas/terraform-provider
git clone https://git.arvancloud.ir/arvancloud/iaas/terraform-provider.git



2.  **Initialize Terraform:**
    ```bash
    terraform init
    terraform apply -auto-approve
    ```
    This command downloads the necessary Arvan provider plugins.

3.  **Review and modify the `variables.tf` file:**
    Provide the required values for the variables defined in `variables.tf`. This file typically includes:
    * Arvan region
    * Instance sizes
    * Network CIDR blocks
    * Security group rules


4.  **Plan the infrastructure:**
    ```bash
    terraform plan 
    ```
    This command shows you the changes that Terraform will apply to your Arvan IaaS environment without actually making them. Review the output carefully to ensure it aligns with your expectations.

5.  **Apply the configuration:**
    ```bash
    terraform apply -auto-approve
    ```
    This command creates or modifies the resources as defined in your configuration. You will be prompted to confirm the actions before they are executed.

## Configuration Files

* `main.tf`: Contains the main resource definitions for your AbrAk infrastructure.
* `variables.tf`: Defines the input variables used in the configuration.
* `outputs.tf`: Defines the output values that will be displayed after the deployment.
* **(Add any other relevant configuration files or directories)**

## Usage

After successfully applying the configuration, your AbrAk infrastructure will be running on Arvan IaaS. You can access the deployed resources based on their configurations (e.g., SSH to instances using their public IPs, access web applications via load balancer IPs/DNS names).

To make changes to your infrastructure, modify the configuration files and run `terraform plan` and `terraform apply` again.

## Destroying the Infrastructure

To remove all the resources created by this configuration, run the following command:

```bash
terraform destroy
