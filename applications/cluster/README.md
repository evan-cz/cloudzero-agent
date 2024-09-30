# AWS EKS Training Cluster

This provisions a multi-node AWS EKS cluster foundation ideal for deploying, chart training.

## Prerequisites
1. [Cloudzero profile tool](https://github.com/Cloudzero/dev-tools?tab=readme-ov-file#useful-functions-and-aliases)
2. [Install Pulumi](https://www.pulumi.com/docs/get-started/install/)
3. [Configure Pulumi for Python](https://www.pulumi.com/docs/intro/languages/python/)
4. [Configure AWS Credentials](https://www.pulumi.com/docs/intro/cloud-providers/aws/setup/)
5. [`brew install aws-iam-authenticator`](https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html)
6. [Install `kubectl`](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
7. [Install helm](https://helm.sh/docs/intro/install/)

## 1. Creating an AWS EKS Cluster

After cloning this repo, run these commands from the working directory:

1. Make sure you are authenticated for the Cloudzero AWS `research` Account:
    ```sh
    use_profile cz-research.Engineering
    ```

3. Setup the virtual python environment
    ```sh
    make init
    ```

4. Create a new stack, which is an isolated deployment target for this example:

    ```sh
    pulumi stack init py-eks-workshop
    ```

5. Setup your configuration values.

    <details collapse=true>
    <summary>Configure your pulumi deployment</summary>


    | Value Name                   | Description                                          | Example                   |
    | ---------------------------- | ---------------------------------------------------- | ------------------------- |
    | `aws:region`                 | The AWS region to deploy into                        | us-east-2                 |
    | `cz:owner`                   | The owner email address                              | joe.barnett@cloudzero.com |
    | `cz:team`                    | The team name                                        | cirrus                    |
    | `cz:namespace`               | The environment namespace                            | jb                        |
    | `cz:purpose`                 | The purpose for the cluster                          | iac-workshop              |
    | `machine:type`               | The type of EC2 instance to use                      | t3.small                  |
    | `machine:min`                | The minimum number of EC2 instances in the node group| 1                         |
    | `machine:desired`            | The desired number of EC2 instances in the node group| 2                         |
    | `machine:max`                | The maximum number of EC2 instances in the node group| 2                         |

    You can use the `pulumi config set` command to set the values, for example:

    ```sh
    pulumi config set aws:region us-east-2
    pulumi config set cz:owner joe.barnett@cloudzero.com
    pulumi config set cz:team cirrus
    pulumi config set cz:namespace jb
    pulumi config set cz:purpose training
    pulumi config set machine:type t3.small
    pulumi config set machine:min 1
    pulumi config set machine:desired 1
    pulumi config set machine:max 1
    pulumi config set machine:max_pods 110
    ```

    </details>

6. Execute the Pulumi program to create our EKS Cluster:

	```sh
	pulumi up
	```
    > Note use the `--yes` flag if you do not want the prompt confirmation shown below


    <details collapse=true>
    <summary>Example Output</summary>

    **example output**
    ```sh
    Previewing update (py-eks-workshop)

    View in Browser (Ctrl+O): https://app.pulumi.com/josephbarnett/cloudzero-eks/py-eks-workshop/previews/f47825cf-0ca9-4673-8b33-dc1770dd88d9

    Downloading plugin: 275.36 MiB / 275.36 MiB [======================] 100.00% 47s

    [resource plugin aws-6.46.0] installing
        Type                             Name                                     Plan       
    +   pulumi:pulumi:Stack              cloudzero-eks-py-eks-workshop            create     
    +   ├─ aws:iam:Role                  aws-jb-cirrus-iac-workshop-eks-role      create     
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-rpa-0         create     
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-ngpa-0        create     
    +   ├─ aws:iam:Role                  aws-jb-cirrus-iac-workshop-nodegrp-role  create     
    +   ├─ aws:ec2:SecurityGroup         aws-jb-cirrus-training-cluster-sg        create     
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-rpa-1         create     
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-ngpa-1        create     
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-ngpa-2        create     
    +   ├─ aws:eks:Cluster               aws-jb-cirrus-training-cluster           create     
    +   └─ aws:eks:NodeGroup             aws-jb-cirrus-iac-workshop-ng            create     

    Outputs:
        clusterName: "aws-jb-cirrus-training-cluster"

    Resources:
        + 11 to create

    Do you want to perform this update? yes
    Updating (py-eks-workshop)

    View in Browser (Ctrl+O): https://app.pulumi.com/josephbarnett/cloudzero-eks/py-eks-workshop/updates/1

        Type                             Name                                     Status              
    +   pulumi:pulumi:Stack              cloudzero-eks-py-eks-workshop            created (660s)      
    +   ├─ aws:iam:Role                  aws-jb-cirrus-iac-workshop-eks-role      created (1s)        
    +   ├─ aws:iam:Role                  aws-jb-cirrus-iac-workshop-nodegrp-role  created (1s)        
    +   ├─ aws:ec2:SecurityGroup         aws-jb-cirrus-training-cluster-sg    created (4s)        
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-rpa-1         created (0.87s)     
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-rpa-0         created (1s)        
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-ngpa-0        created (1s)        
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-ngpa-1        created (1s)        
    +   ├─ aws:iam:RolePolicyAttachment  aws-jb-cirrus-iac-workshop-ngpa-2        created (1s)        
    +   ├─ aws:eks:Cluster               aws-jb-cirrus-training-cluster       created (511s)      
    +   └─ aws:eks:NodeGroup             aws-jb-cirrus-iac-workshop-ng            created (139s)      

    Outputs:
        clusterName: "aws-jb-cirrus-training-cluster"

    Resources:
        + 11 created

    Duration: 11m3s
    ```

    > Use the `up` arrow to select `yes` then hit the enter button.

    </details>

7. After 10-15 minutes, your cluster will be ready, then we can configure `kubectl` using the following AWS command:

    <details collapse=true>
    <summary>Update k8s kubectl configuration using AWS CLI</summary>

    ```sh
    aws eks update-kubeconfig --region us-east-2 --name aws-jb-cirrus-training-cluster
    ```
    > Make sure to use the generated cluster name matches the output of the pulumi deployment

    **verify the context is set for kubectl**
    ```sh
    kubectl config get-contexts
    ```

    **example output**
    ```sh
    CURRENT   NAME                                                                            CLUSTER                                                                         AUTHINFO                                                                        NAMESPACE
    *         arn:aws:eks:us-east-2:975482786146:cluster/aws-jb-cirrus-training-cluster   arn:aws:eks:us-east-2:975482786146:cluster/aws-jb-cirrus-training-cluster   arn:aws:eks:us-east-2:975482786146:cluster/aws-jb-cirrus-training-cluster
    ```

    **Test kubectl works as expected**
    ```sh
    kubectl get nodes
    ```

    </details>

    <details collapse=true>
    <summary>Get local k8s kubectl configuration from pulumi deployment</summary>

    ```sh
    pulumi stack output kubeConfig --show-secrets > kubeconfig.yaml
    ```

    Then to use it with kubectl, do the following:
    ```sh
    KUBECONFIG=./kubeconfig.yaml kubectl get nodes
    ```
    </details>


# Resource cleanup

To save money, always remember to cleanup your resources

1. Cleanup your kubectl context configuration:

    ```sh
    kubectl config delete-context arn:aws:eks:us-east-2:975482786146:cluster/aws-jb-cirrus-py-iac-workshop-cluster
    ```

2. Delete the AWS Resources

    ```sh
    pulumi destroy --yes
    pulumi stack rm --yes
    ```
