# Multi-VPC Serverless Kafka

An MSK Serverless cluster reachable from TWO VPCs at once: AWS provisions client-facing network interfaces in every declared placement, so producers in the application VPC and consumers in the analytics VPC each connect privately — no VPC peering, no PrivateLink endpoints, no transit routing.

## What this preset gives you

- One shared Kafka cluster serving workloads in separate VPCs, each through its own subnets and its own security group.
- Per-VPC network policy: each placement's security group owns the port-9098 ingress rule for that VPC's clients, so access is controlled where each client population lives.
- The same fully managed, IAM-authenticated Kafka as the basic preset — capacity scaling, storage, and versions are AWS's problem.

## Before you deploy

- Replace the example IDs: each placement's subnets must belong to the SAME VPC (use two different AZs per placement for production), and each security group must live in that placement's VPC.
- Placement is create-time: adding or removing a VPC later replaces the cluster, so settle the VPC set first.

## Remix ideas

- Reference Planton-managed `AwsSubnet` / `AwsSecurityGroup` nodes with `valueFrom` so both network layers compose into the resource graph.
- Grant each VPC's workloads `kafka-cluster:*` IAM permissions scoped to the exported `cluster_arn` — IAM is the only client-authentication scheme, so access control is pure IAM policy.
