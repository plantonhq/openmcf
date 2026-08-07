---
title: "Redshift Serverless Workgroup"
description: "Redshift Serverless Workgroup deployment documentation"
icon: "package"
order: 100
componentName: "awsredshiftserverlessworkgroup"
---

# AWS Redshift Serverless Workgroup

Deploys an Amazon Redshift Serverless workgroup — the compute plane of the serverless warehouse: Redshift Processing Unit (RPU) capacity, VPC placement, network reachability, and query-level configuration. A workgroup computes; the data it serves lives on the [AwsRedshiftServerlessNamespace](/cloud-catalog/aws-redshift-serverless-namespace) it attaches to by name. Billing follows the compute — RPU-hours accrue only while queries execute, so an idle workgroup costs nothing. Many workgroups can serve one namespace (a capped dev endpoint and an autoscaling production endpoint over the same data), and each is created and destroyed without touching what is stored. The workgroup integrates with Planton's Provider Connections for AWS credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Redshift Serverless Workgroup** -- the compute plane whose name is the resource name (create-time immutable); SQL clients connect to its endpoint
- **Capacity Posture** -- a fixed RPU baseline you choose, or the price-performance dial where AWS owns the baseline; plus the optional hard spend ceiling
- **VPC Placement** -- compute and a managed VPC endpoint in the subnets you reference (three AZs minimum), guarded by the security groups you attach
- **Reachability Posture** -- enhanced VPC routing for governed COPY/UNLOAD data movement, private-by-default endpoint exposure, and the connection port
- **Query Configuration** -- session parameters and query monitoring guardrails applied directly to the workgroup (serverless has no parameter groups), plus the release track
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **The namespace first** -- deploy the [AwsRedshiftServerlessNamespace](/cloud-catalog/aws-redshift-serverless-namespace) this workgroup serves; the workgroup references its `namespace_name` output.

### AWS Account

- **Three subnets, three AZs** -- Redshift Serverless refuses a workgroup with fewer than three subnets spanning three distinct Availability Zones (leave the list empty only to use the account's default VPC). Each subnet needs free IPs in proportion to base capacity.
- **Ingress on the security groups** -- warehouse ingress rules (e.g. port 5439 from BI tooling) belong on the referenced [AwsSecurityGroup](/cloud-catalog/aws-security-group) nodes, never inside the workgroup.

## Deploy

### Console

Open the deployment store, find **AWS Redshift Serverless Workgroup**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Capped Dev** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRedshiftServerlessWorkgroup
metadata:
  name: analytics-dev
  org: acme-corp
  env: dev
spec:
  region: us-west-2
  namespaceName:
    valueFrom:
      kind: AwsRedshiftServerlessNamespace
      name: analytics-data
      fieldPath: status.outputs.namespace_name
  baseCapacity: 8
  maxCapacity: 32
  subnetIds:
    - value: subnet-0a1b2c3d4e5f60001
    - value: subnet-0a1b2c3d4e5f60002
    - value: subnet-0a1b2c3d4e5f60003
```

```shell
planton apply -f redshift-serverless-workgroup.yaml
```

This creates a cost-bounded dev workgroup over the referenced namespace. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, the namespace deploys first, then workgroups reference it — and applications consume the workgroup's endpoint:

```yaml
# This workgroup's connect string, from another resource's env wiring:
valueFrom:
  kind: AwsRedshiftServerlessWorkgroup
  name: analytics-prod
  fieldPath: status.outputs.endpoint_address
```

## Key Configuration

These are the most important decisions when configuring a serverless workgroup. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**The capacity model** -- mutually exclusive by design: fix the RPU baseline yourself (`baseCapacity`, 0 keeps the AWS default 128), or enable `pricePerformanceTarget` and AWS picks and adjusts the baseline against a cost/speed dial (1 cheapest, 50 balanced, 100 fastest). The console's mode selector makes the exclusivity structural. `maxCapacity` applies on BOTH models — the worst-case-spend guardrail production workgroups should always set.

**The data/compute split** -- destroying or recreating this workgroup never touches the namespace's data. Run several workgroups over one namespace when different consumers need different capacity or network postures; the attachment itself is create-time immutable.

**Network posture** -- private by default (Query Editor and in-VPC BI need no public IP). `enhancedVpcRouting` forces COPY/UNLOAD data movement through the VPC where flow logs and endpoints govern it — the usual data-governance ask. The port accepts only 5431-5455 and 8191-8215 (default 5439).

**Query guardrails** -- parameters apply directly to the workgroup from the exact list the API accepts: `require_ssl` and `max_query_execution_time` are the production classics; the `max_*` family implements query monitoring rules that cancel runaway work. `enable_user_activity_logging` pairs with the namespace's `useractivitylog` export — the workgroup produces the trail, the namespace delivers it.

## Outputs and Dependencies

### What This Component Consumes

| Field | Referenced Kind | Purpose |
|-------|-----------------|---------|
| `namespaceName` | [AwsRedshiftServerlessNamespace](/cloud-catalog/aws-redshift-serverless-namespace) | The data plane this compute serves (required) |
| `subnetIds[]` | [AwsSubnet](/cloud-catalog/aws-subnet) | Compute and endpoint placement (three AZs minimum) |
| `securityGroupIds[]` | [AwsSecurityGroup](/cloud-catalog/aws-security-group) | Who can reach the endpoint |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `endpoint_address` | The DNS hostname SQL clients connect to | Application/BI connection strings |
| `port` | The port the workgroup accepts connections on | Paired with the endpoint address |
| `workgroup_name` | The workgroup name (the resource name) | GetCredentials, custom domain associations |
| `workgroup_id` | The unique identifier AWS assigns | Account-level automation and audits |
| `arn` | Amazon Resource Name of the workgroup | IAM policies and usage limits |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Capped dev workgroup** -- the smallest practical baseline (8 RPU) with a hard 32-RPU ceiling: cheap, bounded, and safe to leave running (idle costs nothing). Start from the **Capped Dev** preset.

**Price-performance production workgroup** -- AWS owns the baseline at the balanced level, a 512-RPU cap bounds spend, enhanced VPC routing governs data movement, TLS is required, and a four-hour query limit guards runaway work. Start from the **Price-Performance Production** preset.

## Works With

- [**AWS Redshift Serverless Namespace**](/cloud-catalog/aws-redshift-serverless-namespace) -- the data plane this workgroup computes for (references `namespace_name`)
- [**AWS Subnet**](/cloud-catalog/aws-subnet) -- placement for the compute and its managed VPC endpoint
- [**AWS Security Group**](/cloud-catalog/aws-security-group) -- carries the warehouse's ingress rules
- [**AWS Redshift Cluster**](/cloud-catalog/aws-redshift-cluster) -- the provisioned alternative when steady, predictable load makes reserved capacity cheaper
