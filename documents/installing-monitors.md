# CloudZero Agent and Service Discovery
This tutorial guides you through the process of installing the CloudZero agent, _when an existing `kube-state-metrics` are deployed in a different `kubernetes namespace_`.

One valuable aspect of this guide highlights prometheus's ability for service discovery.

An additional twist of this guide is using value.yaml files to override default values in a Helm chart deployment. This skill is valuable when thinking about automation across multiple environments.

## Prerequisites

1. [`kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
2. [Install `Helm`](https://helm.sh/docs/intro/install/)
3. [A Kubernetes Cluster](../../aws/README.md)

## Step-by-Step Guide

### Deployment Steps

1. **Create a different namespace for `kube-state-metrics`, and `prometheus-node-exporter`:**

    ```sh
    kubectl apply -f deployments/monitors/namespace.yaml
    ```

2. **Install `kube-state-metrics`**:

    First, add the required repository if not already added:

    ```sh
    helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
    helm repo update
    ```

    Now, install the chart:

    ```sh
    helm upgrade --install kube-state-metrics prometheus-community/kube-state-metrics -n monitors -f deployments/monitors/kube-state-metrics-values.yaml
    ```

3. Edit the deployment for your [cloudzero-agent deployment variables](../../deployments/monitors/cloudzero-agent-values.yaml) in `deployments/monitors/cloudzero-agent-values.yaml`:

    ```yaml
    cloudAccountId: |-
      975482786146
    clusterName: cloudzero-eks-cluster-eksCluster-b4e6994
    region: us-east-2
    existingSecretName: cz-api-token

    kube-state-metrics:
    enabled: false
    prometheus-node-exporter:
    enabled: false

    server:
    args:
    # ENABLE UI FOR DEPLOYMENT - TO ALLOW VERIFYING SCRAPE CONFIG CHANGES
    - --config.file=/etc/config/prometheus/configmaps/prometheus.yml
    - --web.enable-lifecycle
    - --web.console.libraries=/etc/prometheus/console_libraries
    - --web.console.templates=/etc

    validator:
    image:
        tag: 0.6.0
    ```

4. **Create a Secret for the CloudZero API Key**:

   Ensure your [CloudZero API token](https://app.cloudzero.com/organization/api-keys) is ready. This token allows the CloudZero Agent to securely communicate with your CloudZero account. If you do not already have an API token, you may need to generate one from the [CloudZero API Key page](https://app.cloudzero.com/organization/api-keys).

   Export your API token as an environment variable:

    ```sh
    export CZ_API_TOKEN='your-api-token'
    ```

   Replace `your-api-token` with the actual API token provided by CloudZero.

   Now, create the Kubernetes secret in the designated namespace:

    ```sh
    kubectl -n monitors create secret generic cz-api-token --from-literal=value=${CZ_API_TOKEN}
    ```

    **Details**:
    - `cz-api-token`: This is the name of the Kubernetes secret where your API token will be stored.
    - `--from-literal=value=${CZ_API_TOKEN}`: This command sets the secret's `value` field to the value of your API token.

    > **Note**: This step assumes you have the `kubectl` command line tool installed and configured to communicate with your Kubernetes cluster.

   Verify the secret's presence:

    ```sh
    kubectl -n monitors get secrets
    ```

5. **Install the CloudZero Agent**:

    Add the CloudZero repository if you haven't already:

    ```sh
    pushd charts/charts/cloudzero-agent
    ```

    Check out your development `branch` which you'd like to work with:

    ```sh
    git checkout my-branch
    ```

    Invoke helm build dependencies to allow the local installation:

    ```sh
    helm dependency build
    ```

    Install the CloudZero Agent using the values file:

    ```sh
    helm upgrade --install cloudzero-agent . --namespace monitors -f ../../../deployments/monitors/cloudzero-agent-values.yaml
    ```

### Validation

1. **Verify all components are running**:

    ```sh
    kubectl -n monitors get pods
    ```

    You should see pods for `kube-state-metrics`, `node-exporter`, and `cloudzero-agent` running.

2. **Check connectivity and data flow**:

    Log into [CloudZero](https://app.cloudzero.com/organization/k8s-integration) to confirm that the cluster is correctly reporting data in the Kubernetes Integrations Page.
