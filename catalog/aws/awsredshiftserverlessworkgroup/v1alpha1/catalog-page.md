# AWS Redshift Serverless Workgroup

Deploys an Amazon Redshift Serverless workgroup -- the compute plane of the serverless warehouse: RPU capacity (fixed baseline or an AWS-managed price-performance target), VPC placement across three availability zones, network reachability, and query-level configuration. The data it serves lives on the `AwsRedshiftServerlessNamespace` it attaches to by reference; subnets and security groups also compose by reference. RPU-hours accrue only while queries execute, so an idle workgroup costs nothing.

## What Gets Created

When you deploy an AwsRedshiftServerlessWorkgroup resource, Planton provisions:

- **Redshift Serverless Workgroup** — a `redshiftserverless.Workgroup` with the specified capacity contract, VPC attachment, reachability, config parameters, and release track, attached to the referenced namespace

## Prerequisites

- **AWS credentials** configured via environment variables or Planton provider config
- **A namespace** -- an `AwsRedshiftServerlessNamespace` (or an existing namespace name) in the same region
- **Three subnets in three distinct Availability Zones** -- the Redshift Serverless minimum; omit subnets only in accounts that still have a default VPC
- **Security group IDs** if the endpoint should not use the VPC's default security group -- ingress rules (e.g. port 5439 from BI tooling) belong on the referenced `AwsSecurityGroup` nodes

## Quick Start

Create a file `redshift-serverless-workgroup.yaml`:

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsRedshiftServerlessWorkgroup
metadata:
  name: my-analytics-compute
spec:
  region: us-west-2
  namespaceName:
    valueFrom:
      kind: AwsRedshiftServerlessNamespace
      name: my-analytics-data
      fieldPath: status.outputs.namespace_name
  baseCapacity: 8
  subnetIds:
    - value: "<private-subnet-id-az1>"
    - value: "<private-subnet-id-az2>"
    - value: "<private-subnet-id-az3>"
```

Deploy:

```shell
planton apply -f redshift-serverless-workgroup.yaml
```

This creates a private 8-RPU-baseline workgroup over the referenced namespace's data. Point SQL clients at the `endpoint_address` output on port 5439.

## Configuration Reference

### Required Fields

| Field | Type | Description |
| --- | --- | --- |
| `region` | string | AWS region; must match the namespace, subnets, and security groups |
| `namespaceName` | ref → AwsRedshiftServerlessNamespace | The data plane this compute serves (create-time only) |

### Capacity

| Field | Type | Description |
| --- | --- | --- |
| `baseCapacity` | int | The RPU floor each query starts from; 0 keeps the AWS default (128). Mutually exclusive with an enabled price-performance target |
| `maxCapacity` | int | Hard RPU ceiling; 0 leaves scaling uncapped |
| `pricePerformanceTarget` | object | `enabled` + `level` (1 cheapest … 100 fastest; 0 keeps 50): AWS owns the baseline |

### Networking

| Field | Type | Description |
| --- | --- | --- |
| `subnetIds` | list of refs → AwsSubnet | At least three subnets in three distinct AZs; empty uses the account's default VPC |
| `securityGroupIds` | list of refs → AwsSecurityGroup | Empty uses the VPC's default security group |
| `publiclyAccessible` | bool | Give the endpoint a public IP (off by default) |
| `enhancedVpcRouting` | bool | Force COPY/UNLOAD data movement through the VPC |
| `port` | int | Only 5431-5455 or 8191-8215; 0 keeps 5439 |

### Configuration

| Field | Type | Description |
| --- | --- | --- |
| `configParameters` | list | Engine parameters (`require_ssl`, `search_path`, `max_query_execution_time`, ...) applied directly to the workgroup |
| `trackName` | string | Release track: `current` (default), `trailing`, or a named track |

## Stack Outputs

| Output | Description |
| --- | --- |
| `workgroup_name` | The handle the serverless APIs and credentials API address |
| `workgroup_id` | The unique identifier AWS assigns |
| `arn` | The workgroup ARN, for IAM policies and usage limits |
| `endpoint_address` | The DNS hostname SQL clients connect to |
| `port` | The port the workgroup accepts connections on |
