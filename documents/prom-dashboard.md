# Accessing the Prometheus Dashboard Using Port Forwarding

To access the Prometheus dashboard, follow these steps:

1. **Identify the Agent Pod Name**:

    First, you need to identify the name of the agent pod in the `monitors` namespace. You can list all pods in the namespace using the following command:

    ```sh
    kubectl get pods -n monitors
    ```

2. **Run the Port Forwarding Command**:

    Once you have identified the agent pod name, run the following command to set up port forwarding. Replace `cloudzero-agent-server-5dcdb884f-82fsg` with the actual pod name you identified:

    ```sh
    kubectl port-forward -n monitors pod/cloudzero-agent-server-5dcdb884f-82fsg 9090:9090
    ```

3. **Open the Browser**:

    Open your browser and navigate to [http://localhost:9090](http://localhost:9090) to access the Prometheus dashboard.

