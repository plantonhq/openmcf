# Composed Serverless Kafka (Full References)

An MSK Serverless cluster where every cross-resource input resolves from Planton-managed resources: subnets from `AwsSubnet` nodes and the attached security group from an `AwsSecurityGroup` node. The cluster joins the resource graph — redeploying the network layer flows new IDs into the cluster automatically.

## What this preset gives you

- The same fully managed, IAM-authenticated Kafka as the basic preset, with zero hardcoded infrastructure IDs.
- Network policy as a first-class node: the `kafka-broker-sg` security group owns the port-9098 ingress rule (typically allowing the application tier's group), so it can be shared, audited, and evolved independently of the cluster.

## Before you deploy

- Ensure the referenced `AwsSubnet` and `AwsSecurityGroup` resources exist in the same environment (the cluster declares them as prerequisites).
- Placement is create-time: changing subnets or security groups later replaces the cluster, so settle network placement first.

## Remix ideas

- Point the security group's ingress at your service's security group for tier-to-tier access without CIDR management.
- Add an `AwsIamPolicy` granting the producing/consuming workloads `kafka-cluster:*` actions on the exported `cluster_arn`.
