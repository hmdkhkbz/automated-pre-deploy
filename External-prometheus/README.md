# OpenStack Monitoring with Prometheus, Grafana, and Alertmanager

This repository provides a comprehensive solution for monitoring an OpenStack environment using a powerful stack of open-source tools: Prometheus for metrics collection, Grafana for visualization, and Alertmanager for robust alerting. It leverages the `openstack-exporter` to gather OpenStack service metrics and `node_exporter` for host-level metrics.


## Project Overview

Monitoring an OpenStack cloud is crucial for maintaining performance, identifying issues, and ensuring service availability. This project automates the deployment and configuration of a monitoring stack designed specifically for OpenStack, providing deep insights into both the OpenStack services themselves and the underlying infrastructure nodes.

## Prerequisites

Before you begin, ensure you have the following installed on your monitoring server:

* **Docker:** [Installation Guide](https://docs.docker.com/get-docker/)
* **Docker Compose:** [Installation Guide](https://docs.docker.com/compose/install/)
* ** OpenStack nodes:** For deploying and configuring `node_exporter` on openstack nodes with enables this parameter in kolla-ansible globals.yml.

## Architecture

The monitoring solution consists of the following components:

* **Prometheus Server:** Pulls metrics from exporters.
* **Alertmanager:** Processes alerts from Prometheus and sends notifications.
* **Grafana:** Queries Prometheus for metrics and renders dashboards.
* **OpenStack Exporter:** Runs on a machine with OpenStack client access, exposes OpenStack service metrics.
* **Node Exporter:** Runs on each OpenStack compute, controller, and storage node, exposes host-level metrics.
