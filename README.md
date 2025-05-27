# automated-pre-deploy HA Openstack
automated-pre-deploy Openstack environment (including: 3x Controller, 2x Compute and 2x Network nodes)

![alt text](https://thamed.s3.ir-tbz-sh1.arvanstorage.ir/hope.png)


# OpenStack Multi-Node HA Deployment: Network Design

This repository outlines the network design for a robust, highly available (HA) multi-node OpenStack deployment. The core principle driving this design is **network segregation**, ensuring optimal performance, security, and scalability for various OpenStack traffic types.

## 🌐 Network Types

A well-designed OpenStack deployment leverages multiple distinct networks, each serving a specific purpose. This section details the key network types utilized in our setup.

---

### 1. Management Network

* **Purpose:** This network is dedicated to internal OpenStack service communication. It carries API requests, RPC (Remote Procedure Call) messages between OpenStack services (e.g., Nova, Neutron, Cinder, Glance, Keystone), and synchronization traffic for databases and message queues.
* **Traffic Examples:**
    * `nova-api` communicating with `nova-scheduler`.
    * `neutron-server` interacting with `neutron-agents`.
    * Database replication traffic (e.g., Galera cluster).
    * RabbitMQ/ActiveMQ inter-node communication.
    * SSH access for administrative tasks to compute, controller, and storage nodes.
* **Considerations:**
    * **Criticality:** Extremely vital for OpenStack's operation.
    * **Security:** Should be isolated and highly secured, typically not exposed to tenant networks or the external world.
    * **IP Addressing:** Static IP addresses are commonly used for all OpenStack service endpoints on this network.

---

### 2. Tenant Network

* **Purpose:** This network carries the virtual machine (VM) traffic within a tenant's private network segments. Instances launched by users communicate with each other over this network.
* **Traffic Examples:**
    * VM-to-VM communication within the same tenant's private network.
    * Communication between VMs and router interfaces within a tenant's virtual router.
* **Technologies:**
    * **VLANs (Virtual Local Area Networks):** Used for layer-2 segmentation within a physical network. Each tenant network can be assigned a unique VLAN ID.
    * **VXLAN (Virtual Extensible LAN) / GRE (Generic Routing Encapsulation):** Overlay network technologies that encapsulate Layer 2 Ethernet frames in UDP (VXLAN) or IP (GRE) packets, allowing for larger scale multi-tenancy across Layer 3 networks.
* **Considerations:**
    * **Isolation:** Strict isolation between different tenants' networks is paramount.
    * **Scalability:** Must support a large number of tenant networks and instances.
    * **Bandwidth:** Can experience high traffic volumes, especially with data-intensive applications.

---

### 3. Internal Network (API / Data Plane)

* **Purpose:** This network is often a high-bandwidth, low-latency network primarily used for the data plane of services that require significant throughput, such as storage communication (e.g., Ceph) and sometimes the Nova instance migration traffic. In some designs, the OpenStack API endpoints for services like Glance and Cinder might also be exposed on a dedicated "internal" network for other OpenStack services to consume, keeping them off the general management network for improved performance or isolation.
* **Traffic Examples:**
    * **Storage Traffic:** Ceph replication, OSD (Object Storage Daemon) communication, iSCSI traffic between compute nodes and storage arrays.
    * **Live Migration:** Data transfer during Nova live migrations between compute nodes.
    * (Optional) Internal API endpoints for services like Glance, Cinder, Nova, etc., consumed by other OpenStack services.
* **Considerations:**
    * **Performance:** Requires high bandwidth and low latency for optimal storage and migration performance.
    * **Reliability:** Crucial for data integrity and consistent service operation.
    * **Security:** Typically isolated from external access, but might be accessible to specific OpenStack components.

---

### 4. External Network (Public / Provider Network)

* **Purpose:** This network provides connectivity between OpenStack instances and the outside world (e.g., the internet, corporate intranet). It's where floating IPs are allocated and associated with instances, allowing external access to them. It also provides the gateway for tenant networks to reach external destinations.

## External Ceph RBD (using cephadm)
- 3x mons
- 3x osds

## Ceph Networks
- Public Network: for connection with Storage clients (openstack)
- Cluster Network: Dedicated for replication of Internal osd workloads.

## Ceph RBD Pools
- volumes: for cinder-volume
- images: for glance
- backups: for cinder-backup
- vms: for nova

## ELK Logging Stack

collects all fluentd output from openstack nodes

## External Prom/Grafana Monitoring Stack & alerting

scrapes all host-level and service-level metrics from all openstack nodes.

