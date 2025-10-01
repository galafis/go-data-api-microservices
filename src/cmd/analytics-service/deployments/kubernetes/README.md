# Kubernetes Deployment

Kubernetes manifests for the Analytics Service.

## Overview

Contains Kubernetes YAML files for deploying and managing the Analytics Service in a Kubernetes cluster.

## Contents

- `deployment.yaml` - Service deployment configuration
- `service.yaml` - Service exposure and networking
- `configmap.yaml` - Configuration management
- `secret.yaml` - Sensitive configuration data
- `hpa.yaml` - Horizontal Pod Autoscaler
- `ingress.yaml` - External access configuration

## Usage

Apply manifests using `kubectl apply -f` for production deployments.
