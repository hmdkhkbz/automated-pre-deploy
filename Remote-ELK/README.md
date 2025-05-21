# ELK Stack with Fluentd (Kolla-Ansible on OpenStack Integration)

---

This repository provides a Docker Compose setup for an ELK (Elasticsearch, Logstash, Kibana) stack, specifically designed to integrate with Fluentd instances deployed on OpenStack nodes via Kolla-Ansible. This configuration is ideal for centralizing logs from your OpenStack environment, making them searchable and visualizable within Kibana.


## Overview

This setup leverages Docker and Docker Compose to orchestrate the core ELK services. The key distinction here is that the Fluentd component is *not* part of this Docker Compose file for deployment. Instead, it assumes Fluentd instances will be deployed **on your OpenStack compute/controller nodes using Kolla-Ansible**, collecting logs locally and forwarding them to the Logstash service defined in this setup.

The services managed by this Docker Compose are:

* **Logstash:** Receives logs from the Fluentd instances deployed on your OpenStack nodes, processes them, and then sends the refined logs to Elasticsearch.
* **Elasticsearch:** A distributed search and analytics engine that stores the logs. It's the core of the indexing and search capabilities.
* **Kibana:** A data visualization dashboard for Elasticsearch. It allows you to explore, analyze, and visualize your OpenStack log data in various charts and graphs.

Fluentd, as deployed by Kolla-Ansible, acts as the lightweight and efficient data collector on each OpenStack node, forwarding logs to this centralized Logstash instance.

---

## Components

The `docker-compose.yml` file defines the following services for the centralized ELK stack:

* `logstash`: The log processor, receiving from external Fluentd instances.
* `elasticsearch`: The data store for logs.
* `kibana`: The visualization layer for logs.

The Fluentd configuration provided here (managed by kolla-ansible) would typically deploy with Kolla-Ansible on your OpenStack nodes.


