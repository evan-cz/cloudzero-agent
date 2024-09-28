# Guide to Editing Configuration Maps within Kubernetes Deployments

## Introduction
ConfigMaps in Kubernetes are used to store configuration data in key-value pairs. This guide will walk you through the steps to edit ConfigMaps within your Kubernetes deployments.

## Prerequisites
- Kubernetes cluster access
- `kubectl` command-line tool installed and configured

## Steps to Edit ConfigMaps

### 1. List Existing ConfigMaps
To list all ConfigMaps in a specific namespace, use:
```sh
kubectl get configmaps -n <namespace>
```

### 2. View a ConfigMap
To view the details of a specific ConfigMap, use:
```sh
kubectl describe configmap <configmap-name> -n <namespace>
```

### 3. Edit a ConfigMap
To edit a ConfigMap, use:
```sh
kubectl edit configmap -n <namespace> <configmap-name>
```
This command opens the ConfigMap in your default text editor. Make the necessary changes and save the file.

### 4. Apply Changes
After editing, Kubernetes will automatically apply the changes. However, you may need to restart the pods that use the ConfigMap to pick up the new configuration:
```sh
kubectl rollout restart deployment -n <namespace> <deployment-name> 
```

> NOTE: for the cloudzero prometheus agent - this is not necessary, the `reloader` will activate the change.

---
# Local Editting

### 5. Save and Edit Locally

To save a ConfigMap locally and edit it, use:

```sh
kubectl get configmap -n <namespace> <configmap-name> -o yaml > configmap.yaml
```

This command saves the ConfigMap as a YAML file named `configmap.yaml`.

Open the `configmap.yaml` file in your preferred text editor and make the necessary changes.

### 6. Update the ConfigMap

After editing the local file, update the ConfigMap in the cluster with:

```sh
kubectl apply -f configmap.yaml -n <namespace>
```

This command applies the changes from the local file to the ConfigMap in the specified namespace.