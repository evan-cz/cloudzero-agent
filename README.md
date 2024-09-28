# CloudZero Team Cirrus Development, Stress, and Load Testing

This guide provides instructions for setting up a Kubernetes cluster for development and testing of Team Cirrus services.

By following the steps outlined below, you will be able to deploy necessary monitoring tools, ingest data, and verify the performance of your Kubernetes cluster.

## Step-by-Step

### 1. **Kubernetes Cluster**

Follow the instructions in the [cluster README](applications/cluster/README.md) to set up a Kubernetes cluster suitable for load testing. The value of using the deployment is you can control **_Pod density per node_**, which allows one to overload a node, triggering resource issues.

### 2. Deploy a Kubernetes Dashboard

Follow the steps outlined in the [Kubernetes Dashboard setup guide](documents/k8s-dashboard.md) to deploy a Kubernetes Dashboard in your cluster.

### 3. Setup your Monitors

Follow the steps outlined in the [Installing Monitors](documents/installing-monitors.md) to deploy the CloudZero agent and related monitors.

### 4. Deploy Some stress load to your Cluster

Follow the instructions in the [stress tester README](applications/stress/README.md) to build the stress tester.

Next, refer to the [stress testing load guide](documents/stress-testing-load.md) for detailed instructions on deploying the stress load to your cluster.


### 5. Forcing ingest of metrics data to allow quering Prometheus Staging Tables

Ingest data immediately to allow for verification queries:

```sh
$ python -m scripts.cirrusops org ingest-status -o 02fa7d30-c3de-4e0a-8f1e-2de120e7fd23

$ date -u +"%Y-%m-%dT%H:%M:%SZ"
2024-09-26T21:35:56Z

$ python -m scripts.cirrusops org reingest -o 02fa7d30-c3de-4e0a-8f1e-2de120e7fd23 -m hour --start "2024-09-26 18" --end  "2024-09-26 22"
```

### 6. Quering the data

Inspect at data using the following "queries" - you may need to edit the dates, clusters and org ID info

* [Which metrics are coming from your new cluster](./validation/snowflake_queries/01-external/metric_names.sql)
* [Which labels are coming from your new cluster](./validation/snowflake_queries/01-external/labels.sql)
* [Are there gaps in CPU records which should be every 60 seconds?](./validation/snowflake_queries/02-message-frequency/cpu-usage-minutes-example.sql)
* [Are there gaps in Memory records which should be every 60 seconds?](./validation/snowflake_queries/02-message-frequency/memory-usage-minutes-example.sql)
* [Are there gaps in POD records which should be every 60 seconds?](./validation/snowflake_queries/02-message-frequency/pod-usage-minutes-example.sql)
* [Are there gaps in Node records which should be every 60 seconds?](./validation/snowflake_queries/02-message-frequency/node-usage-minutes-example.sql)


---

# Don't Forget to Clean-up!

Don't forget to bring your cluster down!
