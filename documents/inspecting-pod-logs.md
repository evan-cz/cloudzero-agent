
# Inspecting Multi-Container Pod Logs

When working with Kubernetes, you may need to inspect logs from different containers within the same pod. Below are the steps to follow using `kubectl` commands.

## Example Commands

### Getting a List of Containers in the Pod

To get a list of containers within a specific pod, use the following command:

```sh
kubectl -n monitors get pod cloudzero-agent-server-7684899c67-6nvb6 -o jsonpath='{.spec.containers[*].name}'
```

In this example:
- `-n monitors` specifies the namespace.
- `cloudzero-agent-server-7684899c67-6nvb6` is the pod name.
- `-o jsonpath='{.spec.containers[*].name}'` extracts the container names.

### Getting Pod Information

To get detailed information about a specific pod, use the following command:

```sh
kubectl -n monitors get pod cloudzero-agent-server-7684899c67-6nvb6 -o yaml
```

In this example:
- `-n monitors` specifies the namespace.
- `cloudzero-agent-server-7684899c67-6nvb6` is the pod name.
- `-o yaml` outputs the pod information in YAML format.

### Describing the Pod

To describe a specific pod and get detailed status and event information, use the following command:

```sh
kubectl -n monitors describe pod cloudzero-agent-server-7684899c67-6nvb6
```

In this example:
- `-n monitors` specifies the namespace.
- `cloudzero-agent-server-7684899c67-6nvb6` is the pod name.

### Inspecting Logs from a Specific Container

To inspect logs from a specific container within a pod, use the following command:

```sh
kubectl -n monitors logs -f -c env-validator cloudzero-agent-server-7684899c67-6nvb6
```

In this example:
- `-n monitors` specifies the namespace.
- `-f` enables log streaming.
- `-c env-validator` specifies the container name.
- `cloudzero-agent-server-7684899c67-6nvb6` is the pod name.

### Inspecting Logs from Another Container

If you need to inspect logs from another container within the same pod, use the following command:

```sh
kubectl -n monitors logs -c cloudzero-agent-server -f cloudzero-agent-server-7684899c67-6nvb6
```

In this example:
- `-n monitors` specifies the namespace.
- `-c cloudzero-agent-server` specifies the container name.
- `-f` enables log streaming.
- `cloudzero-agent-server-7684899c67-6nvb6` is the pod name.

By using these commands, you can effectively monitor and troubleshoot your multi-container pods in Kubernetes.