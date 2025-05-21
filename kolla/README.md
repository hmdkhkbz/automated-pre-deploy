# Advanced customized configuration for Kolla 

This repository contains custom configurations and playbooks for deploying OpenStack using Kolla Ansible, with a focus on integrating External Ceph RBD, customizing Fluentd for logging, customizing Internal certificate and nova.conf and enhancing OVN (Open Virtual Network) logging details.

## Repository Structure

The core customization files are organized as follows, aligning with your provided folder structure:

.
├── cinder/

├── fluentd/

├── glance/

├── haproxy/

├── nova/

├── openvswitch/

## Features

This repository provides the following key customizations:

* **Ceph Pools:** Configures `cinder-volume`, `cinder-backup`, `nova`, and `glance` to use Ceph for storage, including `ceph.conf` and keyring management.
* **Customized Fluentd:** Defines custom input and output configurations for Fluentd to tailor log aggregation.
* **Detailed OVN Logging:** Enhances logging for OVN Northbound (NB) and Southbound (SB) databases for better debugging and operational insights.

## Prerequisites

Before using these customizations, ensure you have a functional Kolla Ansible setup. You should be familiar with:

* **Kolla Ansible:** Basic understanding of Kolla Ansible deployment process.
* **Ceph:** Knowledge of Ceph concepts and how to manage Ceph clusters.
* **OpenStack:** Familiarity with Cinder, Glance, and Nova services.
* **Fluentd:** Understanding of Fluentd configuration.

## Customization Details and Usage

Here's how each customization is implemented and how you can integrate it into your Kolla Ansible deployment.

### 1. Ceph Configuration for OpenStack Services

This section details how `cinder`, `glance`, and `nova` services are configured to use Ceph.

* **Ceph Configuration File (`ceph.conf`):**
    You will typically place your customized `ceph.conf` in a location that Kolla Ansible can pick up, such as `etc/kolla/config/ceph.conf`. This file will contain the necessary global Ceph settings, including `fsid`, `mon_host`, and client-specific configurations.

* **Ceph Keyrings:**
    Each OpenStack service that interacts with Ceph requires its own keyring with appropriate capabilities. These keyrings should be placed in the correct paths that Kolla Ansible mounts into the containers. For example, you might place them under `etc/kolla/config/ceph/` or directly mount them.

    * **Cinder (cinder-volume, cinder-backup):**
        The `cinder` folder likely contains overrides for Cinder configuration. To enable Ceph for Cinder, you'll modify `cinder.conf` (often via `etc/kolla/config/cinder/cinder.conf` or a `config.json` entry in Kolla).


        Ensure `ceph.client.cinder.keyring` with read/write capabilities on the `volumes` and `images` (if glance is using it) pools.

    * **Glance:**
        The `glance` folder will contain Glance configuration. To enable Ceph for Glance image storage, you'll typically modify `glance-api.conf`.


        Ensure `ceph.client.glance.keyring` with read/write capabilities on the `images` pool.

    * **Nova (nova-compute for ephemeral volumes):**
        The `nova` folder contains Nova configurations. For Nova instances to use Ceph RBD for ephemeral disks, `nova.conf` needs to be configured. The `nova.conf` file at the root of your repository is likely intended for this.


        Ensure `ceph.client.nova.keyring` (or the one associated with the `rbd_user`) with read/write capabilities on the `volumes` (or ephemeral) pool. You'll also need to create a libvirt secret on each compute node that references this keyring.

**Deployment Steps for Ceph:**
1.  Place your `ceph.conf` and keyring files in appropriate paths in your Kolla Ansible configuration directory (e.g., `etc/kolla/config/`).
2.  Modify the respective service configuration files (`cinder.conf`, `glance-api.conf`, `nova.conf`) either directly in this repository (if using custom config overrides) or by configuring Kolla Ansible's `kolla_extra_config` or `*_config_options` variables in `globals.yml` or your inventory.
3.  Ensure the Ceph client packages are installed on your hosts if Kolla Ansible doesn't handle it for you (though it generally manages container dependencies well).
4.  Run `kolla-ansible deploy` or `kolla-ansible reconfigure` after making changes.

### 2. Customized Fluentd Input and Output

The `fluentd` folder is where your custom Fluentd configurations reside. This allows you to define specific log sources (inputs) and destinations (outputs) beyond Kolla's defaults.

* **Custom Fluentd Configuration:**
    You'll typically create files like `etc/kolla/config/fluentd/fluent.conf` or separate configuration snippets. These files will be mounted into the Fluentd containers.

**Deployment Steps for Fluentd:**
1.  Place your customized Fluentd configuration files within the `fluentd/` directory of this repository, mimicking the path Kolla Ansible expects (e.g., `etc/kolla/config/fluentd/fluent.conf`).
2.  Kolla Ansible will mount these files into the Fluentd containers.
3.  Run `kolla-ansible deploy` or `kolla-ansible reconfigure` to apply the changes.

### 3. Detailed OVN SB and NB Logging

The `openvswitch` folder likely contains configurations related to OVN and Open vSwitch. To get more detailed logs for OVN Southbound (SB) and Northbound (NB) databases, you'll need to adjust the logging levels for the `ovn-controller` and `ovn-northd` services.

This is typically done by modifying the `ovn-controller` and `ovn-northd` container environment variables or by passing specific arguments. In Kolla Ansible, you might achieve this through `kolla_extra_config` for OVN services or by adding arguments to `ovn-controller` and `ovn-northd` startup commands.

* **How to Set Logging Levels:**
    You'll usually modify the `ovn-controller` and `ovn-northd` services' startup parameters to include verbose logging flags. 

**Deployment Steps for OVN Logging:**
1.  Modify your Kolla Ansible `globals.yml` or inventory file to include the `openvswitch_ovn_northd_extra_args` and `openvswitch_ovn_controller_extra_args` variables with the desired logging levels (e.g., `--verbose=dbg`).
2.  Run `kolla-ansible reconfigure` or `kolla-ansible deploy` to apply the changes to the OVN containers.
3.  Verify the logs in the `/var/log/openvswitch/` directory within the respective containers.

