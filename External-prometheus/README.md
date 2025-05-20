# OpenStack Monitoring with Prometheus, Grafana, and Alertmanager

This repository provides a comprehensive solution for monitoring an OpenStack environment using a powerful stack of open-source tools: Prometheus for metrics collection, Grafana for visualization, and Alertmanager for robust alerting. It leverages the `openstack-exporter` to gather OpenStack service metrics and `node_exporter` for host-level metrics.

## Table of Contents

* [Project Overview](#project-overview)
* [Features](#features)
* [Prerequisites](#prerequisites)
* [Architecture](#architecture)
## Project Overview

Monitoring an OpenStack cloud is crucial for maintaining performance, identifying issues, and ensuring service availability. This project automates the deployment and configuration of a monitoring stack designed specifically for OpenStack, providing deep insights into both the OpenStack services themselves and the underlying infrastructure nodes.

## Features

* **Comprehensive Metric Collection:** Gathers metrics from OpenStack services (e.g., Nova, Neutron, Cinder, Keystone) via `openstack-exporter` and host-level metrics (CPU, memory, disk, network) via `node_exporter`.
* **Powerful Data Storage & Querying:** Utilizes Prometheus for efficient time-series data storage and a flexible query language (PromQL).
* **Rich Visualization:** Leverages Grafana to create interactive and insightful dashboards for visualizing your OpenStack infrastructure's health and performance.
* **Advanced Alerting:** Configures Alertmanager to send notifications based on predefined alert rules, integrating with various notification channels (e.g., email, Slack, PagerDuty).
* **Easy Deployment:** Provides `docker-compose` files for quick and consistent deployment of the entire monitoring stack.

## Prerequisites

Before you begin, ensure you have the following installed on your monitoring server:

* **Docker:** [Installation Guide](https://docs.docker.com/get-docker/)
* **Docker Compose:** [Installation Guide](https://docs.docker.com/compose/install/)
* **SSH access to OpenStack nodes:** For deploying and configuring `node_exporter`.

## Architecture

The monitoring solution consists of the following components:

* **Prometheus Server:** Pulls metrics from exporters.
* **Alertmanager:** Processes alerts from Prometheus and sends notifications.
* **Grafana:** Queries Prometheus for metrics and renders dashboards.
* **OpenStack Exporter:** Runs on a machine with OpenStack client access, exposes OpenStack service metrics.
* **Node Exporter:** Runs on each OpenStack compute, controller, and storage node, exposes host-level metrics.
