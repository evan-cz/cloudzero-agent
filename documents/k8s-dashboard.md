# Kubernetes Dashboard Deployment Guide

## Prerequisites

- [A running Kubernetes cluster](../applications/cluster/README.md)
- `kubectl` command-line tool configured to communicate with your cluster

## Step-by-Step Guide

### 1. Deploy the Kubernetes Dashboard

Add the Kubernetes Dashboard Helm repository and install the dashboard:

```sh
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/
helm upgrade --install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard --create-namespace --namespace kubernetes-dashboard
```

### 2. Create a Service Account

Apply the service account configuration:

```sh
kubectl apply -f deployments/kubernetes-dashboard/user.yaml
```

### 3. Obtain the Bearer Token

Generate a token for the admin user:

```sh
kubectl -n kubernetes-dashboard create token admin-user
```

### 4. Access the Dashboard

Forward the port to the dashboard container:

```sh
kubectl -n kubernetes-dashboard port-forward svc/kubernetes-dashboard-kong-proxy 8443:443
```

Open your browser and navigate to:

```
https://localhost:8443
```
