# Deploying Stress Load on your Cluster

## 1. Deploy the Stress Pod

To deploy the stress pod, run the following command:

```sh
kubectl apply -f deployments/stress/namespace.yaml -f deployments/stress/deployment.yaml
```

## 2. Scale the Stress Pod

To scale the stress pod to 100 replicas, use the following command:

```sh
kubectl scale -n stress deployment cpu-stressor-deployment --replicas=100
```
