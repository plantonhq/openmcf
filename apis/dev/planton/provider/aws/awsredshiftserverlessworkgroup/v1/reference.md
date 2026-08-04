# AwsRedshiftServerlessWorkgroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1`

AwsRedshiftServerlessWorkgroupSpec defines an Amazon Redshift
Serverless workgroup -- the COMPUTE plane of the serverless warehouse:
Redshift Processing Unit (RPU) capacity, VPC placement, network
reachability, and query-level configuration. A workgroup computes; the
data it serves lives on the AwsRedshiftServerlessNamespace it attaches
to by name.

Many workgroups can serve one namespace -- e.g. a capped dev workgroup
and an autoscaling production workgroup over the same data -- and each
is created and destroyed without touching the namespace. That split is
AWS's own resource model, mirrored here as two composable nodes.
Billing follows the compute: RPU-hours accrue only while queries
execute, so an idle workgroup costs nothing.

The workgroup name comes from metadata.name (create-time immutable in
AWS). Subnets and security groups compose by reference -- warehouse
ingress rules belong on the referenced AwsSecurityGroup nodes, never
inside this workgroup.

## Example

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsRedshiftServerlessWorkgroup
metadata:
  name: awsredshiftserverlessworkgroup-demo
spec:
  region: us-west-2
  namespaceName:
    value: awsredshiftserverlessnamespace-demo
  baseCapacity: 8
  subnetIds:
    - value: subnet-0a1b2c3d4e5f60001
    - value: subnet-0a1b2c3d4e5f60002
    - value: subnet-0a1b2c3d4e5f60003
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.namespaceName` | `string \| valueFrom` | yes |  | AwsRedshiftServerlessNamespace (`status.outputs.namespace_name`) |
| `spec.baseCapacity` | `int32` |  |  |  |
| `spec.maxCapacity` | `int32` |  |  |  |
| `spec.pricePerformanceTarget` | `AwsRedshiftServerlessWorkgroupPricePerformanceTarget` |  |  |  |
| `spec.pricePerformanceTarget.enabled` | `bool` |  |  |  |
| `spec.pricePerformanceTarget.level` | `int32` |  |  |  |
| `spec.subnetIds` | `[]string \| valueFrom` |  |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.securityGroupIds` | `[]string \| valueFrom` |  |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.enhancedVpcRouting` | `bool` |  |  |  |
| `spec.publiclyAccessible` | `bool` |  |  |  |
| `spec.port` | `int32` |  |  |  |
| `spec.configParameters` | `[]AwsRedshiftServerlessWorkgroupConfigParameter` |  |  |  |
| `spec.configParameters[].name` | `string` | yes |  |  |
| `spec.configParameters[].value` | `string` | yes |  |  |
| `spec.trackName` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region the workgroup is created in. Must match the region
of the namespace and of the subnets and security groups it
references.
Example: "us-west-2", "eu-west-1".

- rule: {"string":{"minLen":"1"}}

### spec.namespaceName

`string | valueFrom` · required

The namespace this workgroup serves -- the data plane behind this
compute. Create-time only: moving a workgroup to another namespace
replaces it. Reference an AwsRedshiftServerlessNamespace
namespace_name output or pass a literal namespace name.

- references: AwsRedshiftServerlessNamespace (`status.outputs.namespace_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsRedshiftServerlessNamespace, name: <that resource's name>, fieldPath: status.outputs.namespace_name}} -- a bare string does not parse

### spec.baseCapacity

`int32`

The baseline compute the workgroup starts every query with, in
Redshift Processing Units. 0 keeps the AWS default (128 RPU).
AWS currently accepts 4-1024 RPU (values above 512 in units of 8;
small values in the 4/8-RPU steps AWS publishes) -- the exact
increments have changed over time, so they are validated by AWS at
deploy rather than frozen here. Higher base = faster first query,
higher floor cost per second of execution. Mutually exclusive with
an enabled price_performance_target, where AWS picks the baseline.

### spec.maxCapacity

`int32`

A hard ceiling on the compute the workgroup may scale to, in RPUs.
0 leaves scaling uncapped (the AWS default) -- the workgroup grows
to whatever the query mix demands. Set it to bound worst-case
spend; must be at least base_capacity when both are set.

### spec.pricePerformanceTarget

`AwsRedshiftServerlessWorkgroupPricePerformanceTarget`

Let AWS choose capacity against a price-performance dial instead of
a fixed base: level 1 leans cheapest, 100 leans fastest, 50 is
balanced. When enabled, base_capacity must stay unset (AWS owns the
baseline); max_capacity still caps spend.

### spec.pricePerformanceTarget.enabled

`bool`

Turn price-performance targeting on. While enabled, AWS owns the
capacity baseline and base_capacity must stay unset.

### spec.pricePerformanceTarget.level

`int32`

Where on the dial AWS should aim: 1 (cheapest), 25, 50 (balanced),
75, or 100 (fastest). 0 keeps the AWS default (50).

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"in":[1,25,50,75,100]}}

### spec.subnetIds

`[]string | valueFrom`

Subnets the workgroup places its compute (and managed VPC endpoint)
in. AWS requires at least THREE subnets spanning three distinct
availability zones; the free-IP requirement per subnet scales with
base capacity. Empty lets AWS use the account's default VPC (only
meaningful in accounts that still have one). Reference AwsSubnet
subnet_id outputs or pass literal subnet IDs.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.securityGroupIds

`[]string | valueFrom`

Security groups attached to the workgroup's endpoint. Empty uses
the VPC's default security group (the AWS default). Reference
AwsSecurityGroup security_group_id outputs or pass literal SG IDs
-- warehouse ingress rules (e.g. port 5439 from your BI tooling)
belong on the referenced AwsSecurityGroup node, never inside this
workgroup.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.enhancedVpcRouting

`bool`

Force all COPY and UNLOAD traffic between the workgroup and data
repositories (S3, DynamoDB, ...) through the VPC instead of the
public internet -- enabling VPC flow logs, endpoints, and other
network controls to see and govern warehouse data movement.

### spec.publiclyAccessible

`bool`

Give the workgroup's endpoint a public IP so it can be reached from
outside the VPC. Off by default -- serverless warehouses almost
always stay private behind VPC routing (and Query Editor / private
BI reach them fine).

### spec.port

`int32`

The port the workgroup accepts connections on. 0 keeps the AWS
default (5439). Redshift Serverless only accepts ports within
5431-5455 or 8191-8215.

### spec.configParameters

`[]AwsRedshiftServerlessWorkgroupConfigParameter`

Query-level configuration parameters (e.g. "require_ssl",
"max_query_execution_time", "search_path"), applied directly to the
workgroup -- serverless has no parameter groups, so there is
nothing to fold or reference. The name list mirrors what the
Redshift Serverless API accepts.

### spec.configParameters[].name

`string` · required

The parameter name. The Redshift Serverless API accepts exactly
this set; workload-limit parameters (max_query_*, max_scan_*,
max_*_row_count) implement query monitoring rules.

- rule: {"required":true,"string":{"in":["auto_mv","datestyle","enable_case_sensitive_identifier","enable_user_activity_logging","query_group","search_path","max_query_cpu_time","max_query_blocks_read","max_scan_row_count","max_query_execution_time","max_query_queue_time","max_query_cpu_usage_percent","max_query_temp_blocks_to_disk","max_join_row_count","max_nested_loop_join_row_count","require_ssl","use_fips_ssl"]}}

### spec.configParameters[].value

`string` · required

The parameter value. Required.

- rule: {"required":true}

### spec.trackName

`string`

The release track the workgroup follows: "current" (the AWS
default -- the latest certified release) or "trailing" (one
certified release behind); AWS also accepts named preview tracks.
Empty keeps the AWS default.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"maxLen":"256","pattern":"^[a-zA-Z0-9_]+$"}}

## Validation Rules

- `base_capacity_xor_price_performance`: base_capacity cannot be set when price_performance_target is enabled -- either fix the baseline yourself or let AWS pick it against the price-performance dial
- `base_capacity_positive`: base_capacity cannot be negative -- 0 keeps the AWS default (128 RPU)
- `max_capacity_positive`: max_capacity cannot be negative -- 0 leaves scaling uncapped (the AWS default)
- `max_capacity_covers_base`: max_capacity must be at least base_capacity -- a ceiling below the baseline leaves the workgroup nothing to run on
- `subnets_span_three_azs`: provide at least three subnet_ids (in three distinct availability zones) -- Redshift Serverless refuses a workgroup with fewer; leave the list empty only to use the account's default VPC
- `port_in_serverless_ranges`: port must be within 5431-5455 or 8191-8215 -- the only ranges Redshift Serverless accepts; 0 keeps the AWS default (5439)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsRedshiftServerlessWorkgroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.workgroupName` | `string` | The workgroup name -- the handle the Redshift Serverless APIs, the credentials API (GetCredentials), and custom domain associations address the workgroup by. |
| `status.outputs.workgroupId` | `string` | The unique identifier AWS assigns to the workgroup. |
| `status.outputs.arn` | `string` | The Amazon Resource Name of the workgroup, for IAM policies and usage limits. |
| `status.outputs.endpointAddress` | `string` | The DNS hostname SQL clients connect to. |
| `status.outputs.port` | `int32` | The port the workgroup accepts connections on. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespaceName` | AwsRedshiftServerlessNamespace | `status.outputs.namespace_name` |
| `spec.subnetIds` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
