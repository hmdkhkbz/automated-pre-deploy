# Advanced customized configuration for Kolla 

This repository contains custom configurations and playbooks for deploying OpenStack using Kolla Ansible, with a focus on integrating Ceph storage, customizing Fluentd for logging, and enhancing OVN (Open Virtual Network) logging details.

## Repository Structure

The core customization files are organized as follows, aligning with your provided folder structure:

.
├── cinder/
├── fluentd/
├── glance/
├── haproxy/
├── nova/
├── openvswitch/
└── nova.conf

## Features

This repository provides the following key customizations:

* **Ceph Integration:** Configures `cinder-volume`, `cinder-backup`, `nova`, and `glance` to use Ceph for storage, including `ceph.conf` and keyring management.
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

    ```ini
    # Example content for etc/kolla/config/ceph.conf
    [global]
    fsid = <your_ceph_cluster_fsid>
    mon_host = <ceph_monitor_ip_1>,<ceph_monitor_ip_2> # ...
    auth_cluster_required = cephx
    auth_service_required = cephx
    auth_client_required = cephx
    # ... other global settings

    [client]
    rbd_cache = true
    rbd_cache_max_mb = 256
    rbd_cache_writethrough_until_flush = true
    keyring = /etc/ceph/ceph.client.nova.keyring # Example, refer to service specific keyrings
    # ... other client settings
    ```

* **Ceph Keyrings:**
    Each OpenStack service that interacts with Ceph requires its own keyring with appropriate capabilities. These keyrings should be placed in the correct paths that Kolla Ansible mounts into the containers. For example, you might place them under `etc/kolla/config/ceph/` or directly mount them.

    * **Cinder (cinder-volume, cinder-backup):**
        The `cinder` folder likely contains overrides for Cinder configuration. To enable Ceph for Cinder, you'll modify `cinder.conf` (often via `etc/kolla/config/cinder/cinder.conf` or a `config.json` entry in Kolla).

        ```ini
        # Example for cinder.conf to use Ceph RBD driver
        [DEFAULT]
        # ...
        enabled_backends = rbd

        [rbd]
        volume_backend_name = rbd
        volume_driver = cinder.volume.drivers.rbd.RBDDriver
        rbd_cluster_name = ceph
        rbd_pool = volumes
        rbd_ceph_conf = /etc/ceph/ceph.conf
        rbd_secret_uuid = <your_libvirt_secret_uuid_for_cinder> # If using Cinder with Nova for ephemeral volumes
        rbd_user = cinder
        rbd_keyring_path = /etc/ceph/ceph.client.cinder.keyring
        ```
        Ensure `ceph.client.cinder.keyring` with read/write capabilities on the `volumes` and `images` (if glance is using it) pools.

    * **Glance:**
        The `glance` folder will contain Glance configuration. To enable Ceph for Glance image storage, you'll typically modify `glance-api.conf`.

        ```ini
        # Example for glance-api.conf to use Ceph RBD backend
        [glance_store]
        stores = file,http,rbd
        default_store = rbd
        rbd_store_pool = images
        rbd_store_user = glance
        rbd_store_ceph_conf = /etc/ceph/ceph.conf
        rbd_store_chunk_size = 8
        rbd_store_keyring_path = /etc/ceph/ceph.client.glance.keyring
        ```
        Ensure `ceph.client.glance.keyring` with read/write capabilities on the `images` pool.

    * **Nova (nova-compute for ephemeral volumes):**
        The `nova` folder contains Nova configurations. For Nova instances to use Ceph RBD for ephemeral disks, `nova.conf` needs to be configured. The `nova.conf` file at the root of your repository is likely intended for this.

        ```ini
        # Example for nova.conf to use Ceph RBD for ephemeral storage
        [libvirt]
        images_type = rbd
        images_rbd_pool = volumes # Or a dedicated ephemeral pool
        images_rbd_ceph_conf = /etc/ceph/ceph.conf
        images_rbd_user = nova
        images_rbd_secret_uuid = <your_libvirt_secret_uuid_for_nova>
        ```
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

    ```ini
    # Example content for fluentd/fluent.conf
    # This example adds a new input for a custom application log
    <source>
      @type tail
      path /var/log/my_app/app.log
      pos_file /var/log/td-agent/my_app.log.pos
      tag my_app.log
      <parse>
        @type none # Or use json, regexp, etc.
      </parse>
    </source>

    # This example adds a new output to an external Splunk instance
    <match my_app.log>
      @type splunkhec
      hec_host <splunk_hec_host>
      hec_port 8088
      hec_token <splunk_hec_token>
      # ... other Splunk HEC options
    </match>

    # You can also customize existing outputs, e.g., sending OpenStack logs to a different target
    <match kolla.**>
      @type elasticsearch
      host <your_elasticsearch_host>
      port 9200
      logstash_format true
      # ... other Elasticsearch options
    </match>
    ```

**Deployment Steps for Fluentd:**
1.  Place your customized Fluentd configuration files within the `fluentd/` directory of this repository, mimicking the path Kolla Ansible expects (e.g., `etc/kolla/config/fluentd/fluent.conf`).
2.  Kolla Ansible will mount these files into the Fluentd containers.
3.  Run `kolla-ansible deploy` or `kolla-ansible reconfigure` to apply the changes.

### 3. Detailed OVN SB and NB Logging

The `openvswitch` folder likely contains configurations related to OVN and Open vSwitch. To get more detailed logs for OVN Southbound (SB) and Northbound (NB) databases, you'll need to adjust the logging levels for the `ovn-controller` and `ovn-northd` services.

This is typically done by modifying the `ovn-controller` and `ovn-northd` container environment variables or by passing specific arguments. In Kolla Ansible, you might achieve this through `kolla_extra_config` for OVN services or by adding arguments to `ovn-controller` and `ovn-northd` startup commands.

* **How to Set Logging Levels:**
    You'll usually modify the `ovn-controller` and `ovn-northd` services' startup parameters to include verbose logging flags. For example:

    ```bash
    # For ovn-northd
    --log-file=/var/log/openvswitch/ovn-northd.log --verbose=dbg

    # For ovn-controller
    --log-file=/var/log/openvswitch/ovn-controller.log --verbose=dbg
    ```

    In Kolla Ansible, you can add these to the `openvswitch_ovn_northd_extra_args` and `openvswitch_ovn_controller_extra_args` variables in your `globals.yml` or inventory:

    ```yaml
    # Example for globals.yml or inventory
    openvswitch_ovn_northd_extra_args: "--verbose=dbg"
    openvswitch_ovn_controller_extra_args: "--verbose=dbg"
    ```

    Alternatively, you might need to create a custom entrypoint script or a custom `kolla-ansible` override if direct argument injection isn't sufficient for very specific logging configurations.

**Deployment Steps for OVN Logging:**
1.  Modify your Kolla Ansible `globals.yml` or inventory file to include the `openvswitch_ovn_northd_extra_args` and `openvswitch_ovn_controller_extra_args` variables with the desired logging levels (e.g., `--verbose=dbg`).
2.  Run `kolla-ansible reconfigure` or `kolla-ansible deploy` to apply the changes to the OVN containers.
3.  Verify the logs in the `/var/log/openvswitch/` directory within the respective containers.

