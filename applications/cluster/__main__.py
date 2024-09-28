# Copyright (c) 2024 CloudZero - ALL RIGHTS RESERVED - PROPRIETARY AND CONFIDENTIAL
# Unauthorized copying of this file and/or project, via any medium is strictly prohibited.
# Direct all questions to legal@cloudzero.com
# Adaption for CloudZero use from
# https://github.com/pulumi/pulumi-eks/tree/master/examples/cluster-py
# https://github.com/pulumi/pulumi-eks/blob/master/examples/managed-nodegroups-py
import json
import pulumi
from pulumi_aws import ec2, iam
import pulumi_eks as eks


project_name = pulumi.get_project()


# Standard tags to apply to all resources
required_tags = {
    'cz:feature': 'container-analysis',
    'cz:description': 'Cirrus Kubernetes Development',
    'cz:owner': 'joe.barnett@cloudzero.com',
    'cz:team': 'cirrus@cloudzero.com',
    'cz:repo': 'none',
    'cz:env': 'research',
    'cz:namespace': 'jb',
    'cz:customer-data': 'false',
}


managed_policy_arns = [
    'arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy',
    'arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy',
    'arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly',
]

# Creates a role and attaches the EKS worker node IAM managed policies
def create_role(name: str) -> iam.Role:
    role = iam.Role(name, assume_role_policy=json.dumps({
        'Version': '2012-10-17',
        'Statement': [
            {
                'Sid': 'AllowAssumeRole',
                'Effect': 'Allow',
                'Principal': {
                    'Service': 'ec2.amazonaws.com',
                },
                'Action': 'sts:AssumeRole'
            }
        ]
    }))

    for i, policy in enumerate(managed_policy_arns):
        # Create RolePolicyAttachment without returning it.
        iam.RolePolicyAttachment(f'{name}-policy-{i}',
            policy_arn=policy,
            role=role.id)

    return role

def main():
    # Load configurations
    cz_cfg = pulumi.Config('cz')
    team = cz_cfg.require('team')
    owner = cz_cfg.require('owner')
    namespace = cz_cfg.require('namespace')
    

    machine_cfg = pulumi.Config('machine')
    node_group_instance_type = machine_cfg.require('type')
    node_group_min = machine_cfg.require_int('min')
    node_group_desired = machine_cfg.require_int('desired')
    node_group_max = machine_cfg.require_int('max')
    max_pods_per_node = machine_cfg.require_int('max_pods')

    # Update required tags
    required_tags['cz:owner'] = owner
    required_tags['cz:team'] = team
    required_tags['cz:env'] = namespace

    tags = {k: v for k, v in required_tags.items()}

    # Lookup VPC
    vpc = ec2.get_vpc(default=True)

    # Get subnets
    subnet = ec2.get_subnets(filters=[{'name': 'vpc-id', 'values': [vpc.id]}])
    
    print(f'VPC ID: {vpc.id}')
    print(f'Public Subnet IDs: {subnet.ids}')
    
    # IAM roles for the node groups.
    ng_role = create_role(f'{project_name}-ng-role')

    # Add Name key to the tags dictionary
    cluster_name = f'{project_name}-cluster'
    tags['Name'] = cluster_name
    eks_cluster = eks.Cluster(
        resource_name=cluster_name,
        vpc_id=vpc.id,
        public_subnet_ids=subnet.ids,
        public_access_cidrs=['0.0.0.0/0'],
        instance_roles=[ng_role],
        skip_default_node_group=True,
        # fargate=True, # Enable to also enable fargate
        tags=tags,
    )

    
    eks.NodeGroup(
        f'{project_name}-node-group',
        cluster=eks_cluster,
        instance_type=node_group_instance_type,
        desired_capacity=node_group_desired,
        min_size=node_group_min,
        max_size=node_group_max,
        node_root_volume_type='gp3',
        kubelet_extra_args=f'--max-pods {max_pods_per_node}',
        labels={'ondemand': 'true'},
        instance_profile=iam.InstanceProfile(
            f'{project_name}-ng-inst-profile', role=ng_role.name
        ),
    )

    # TODO: Add support configuration to allow different NodeGroup types
    # eg. Spot instances
    # labels={'preemptible': 'true'}, # spot
    
    pulumi.export('kubeconfig', eks_cluster.kubeconfig)
    pulumi.export('clusterName', cluster_name)

if __name__ == '__main__':
    main()
