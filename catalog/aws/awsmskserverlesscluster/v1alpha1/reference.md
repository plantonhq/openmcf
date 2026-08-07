# AwsMskServerlessCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsMskServerlessClusterSpec defines the desired state of an Amazon MSK
Serverless cluster. Unlike a provisioned MSK cluster, capacity is fully
managed by AWS: there are no brokers, instance types, storage volumes, or
Kafka version to declare -- AWS scales compute and storage automatically and
bills per throughput and storage consumed. The whole declaration is WHERE
the cluster lives (subnets + security groups) and the rest is AWS's.

Two properties shape this spec:

  1. The resource is effectively IMMUTABLE. Every field below is create-time
     (ForceNew) in the AWS provider -- only tags can change in place.
     Changing anything else replaces the cluster.
  2. SASL/IAM is the ONLY client authentication scheme serverless MSK
     supports, and it is mandatory. Clients authenticate with AWS IAM
     credentials on port 9098; there is no SCRAM, mTLS, or unauthenticated
     access. Both IaC modules enable it unconditionally, so it is not a
     field here -- there is nothing to choose.

Network ingress is composed, never embedded: the cluster attaches the
referenced security_group_ids directly, and the ingress rule that opens the
SASL/IAM port (9098) lives on those first-class AwsSecurityGroup nodes.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsMskServerlessCluster
metadata:
  name: test-msk-serverless
spec:
  region: us-west-2
  subnetIds:
    - value: "<private-subnet-id-az1>"
    - value: "<private-subnet-id-az2>"
  securityGroupIds:
    - value: "<broker-security-group-id>"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.subnetIds` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.subnetIds

`[]string | valueFrom` · required

subnet_ids are the VPC subnets where the cluster places its network
interfaces. Provide subnets in at least two Availability Zones for
production use; clients connect through these ENIs.
ForceNew: changing subnets forces cluster replacement.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.securityGroupIds

`[]string | valueFrom`

security_group_ids are the security groups ATTACHED to the cluster
network interfaces -- they define what can reach the brokers. The ingress
rule for the SASL/IAM listener port (9098) belongs on these referenced
AwsSecurityGroup nodes. Maximum 5. If omitted, AWS attaches the VPC's
default security group.
ForceNew: adding or removing entries after creation forces replacement.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: {"repeated":{"maxItems":"5"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsMskServerlessCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_arn` | `string` | cluster_arn is the Amazon Resource Name of the MSK Serverless cluster -- also the resource identifier AWS uses for it. Referenced in IAM policies (kafka-cluster:* actions) and Lambda event source mappings. |
| `status.outputs.cluster_name` | `string` | cluster_name is the human-readable name of the cluster. |
| `status.outputs.cluster_uuid` | `string` | cluster_uuid is the unique identifier extracted from the cluster ARN. |
| `status.outputs.bootstrap_brokers_sasl_iam` | `string` | bootstrap_brokers_sasl_iam is the comma-separated SASL/IAM broker endpoint list (port 9098) -- the only connection string serverless MSK exposes. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.subnetIds` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |

## See Also

- [Overview](../README.md)
