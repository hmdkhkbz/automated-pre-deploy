# ELK Stack with Fluentd (Kolla-Ansible on OpenStack Integration)

---

This repository provides a Docker Compose setup for an ELK (Elasticsearch, Logstash, Kibana) stack, specifically designed to integrate with Fluentd instances deployed on OpenStack nodes via Kolla-Ansible. This configuration is ideal for centralizing logs from your OpenStack environment, making them searchable and visualizable within Kibana.

## Table of Contents

* [Overview](#overview)
* [Components](#components)
* [Prerequisites](#prerequisites)
---

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

The Fluentd configuration provided here (`fluentd/fluent.conf`) is a *template* for what you would typically deploy with Kolla-Ansible on your OpenStack nodes.

---

## Prerequisites

Before you begin, ensure you have the following installed:

* **Docker:** [Install Docker](https://docs.docker.com/get-docker/)
* **Docker Compose:** [Install Docker Compose](https://docs.docker.com/compose/install/) (usually comes with Docker Desktop)
* **Familiarity with Kolla-Ansible:** Understanding how Kolla-Ansible deploys and manages OpenStack services, including custom configurations for log aggregation, is crucial for integrating Fluentd.
* **OpenStack Nodes with Kolla-Ansible:** Your OpenStack environment should be deployed or ready for deployment using Kolla-Ansible.

---

## Consider Kolla Ansible version compatibility.



What version of Kolla Ansible are you running? This directly impacts the container image versions it deploys, including Fluentd and its plugins.
Once you align the Fluentd client version with your Elasticsearch/OpenSearch server version, your Fluentd container should start and remain running.
