# AwsMskCluster

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsMskClusterSpec defines the desired state of an Amazon MSK (Managed Streaming for Apache Kafka) cluster.
MSK is a fully managed service that provisions, configures, and maintains Apache Kafka clusters,
handling broker infrastructure, coordination (ZooKeeper or KRaft, depending on the Kafka version),
and storage so teams can focus on producing and consuming streaming data rather than operating
Kafka infrastructure.

Network ingress is composed, never embedded: brokers attach the referenced
security_group_ids directly, and the ingress rules that open the Kafka/ZooKeeper
ports live on those first-class AwsSecurityGroup nodes where they can be shared,
audited, and evolved independently of the cluster.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsMskCluster
metadata:
  name: test-msk-cluster
spec:
  region: us-west-2
  kafkaVersion: "3.6.0"
  numberOfBrokerNodes: 3
  instanceType: kafka.m5.large
  subnetIds:
    - value: subnet-0aaa1111bbb222333
    - value: subnet-0bbb2222ccc333444
    - value: subnet-0ccc3333ddd444555
  securityGroupIds:
    - value: sg-0abc1234def567890
  authentication:
    saslIamEnabled: true
  serverProperties:
    auto.create.topics.enable: "false"
    default.replication.factor: "3"
    min.insync.replicas: "2"
  topics:
    - name: orders.events
      partitionCount: 6
      replicationFactor: 3
      configs:
        retention.ms: "604800000"
        cleanup.policy: delete
        min.insync.replicas: "2"
    - name: orders.events.dlq
      partitionCount: 1
      replicationFactor: 3
  clusterPolicy:
    Version: "2012-10-17"
    Statement:
      - Sid: AllowPartnerVpcConnection
        Effect: Allow
        Principal:
          AWS: arn:aws:iam::444455556666:root
        Action:
          - kafka:CreateVpcConnection
          - kafka:GetBootstrapBrokers
          - kafka:DescribeCluster
          - kafka:DescribeClusterV2
        Resource: "*"
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.kafkaVersion` | `string` | yes |  |  |
| `spec.numberOfBrokerNodes` | `int32` | yes |  |  |
| `spec.instanceType` | `string` | yes |  |  |
| `spec.subnetIds` | `[]string \| valueFrom` | yes |  | AwsSubnet (`status.outputs.subnet_id`) |
| `spec.securityGroupIds` | `[]string \| valueFrom` | yes |  | AwsSecurityGroup (`status.outputs.security_group_id`) |
| `spec.publicAccessType` | `string` |  |  |  |
| `spec.vpcConnectivity` | `AwsMskClusterVpcConnectivity` |  |  |  |
| `spec.vpcConnectivity.saslIamEnabled` | `bool` |  |  |  |
| `spec.vpcConnectivity.saslScramEnabled` | `bool` |  |  |  |
| `spec.vpcConnectivity.tlsEnabled` | `bool` |  |  |  |
| `spec.networkType` | `string` |  |  |  |
| `spec.ebsVolumeSizeGib` | `int32` |  |  |  |
| `spec.provisionedThroughputEnabled` | `bool` |  |  |  |
| `spec.provisionedThroughputMbs` | `int32` |  |  |  |
| `spec.storageMode` | `string` |  |  |  |
| `spec.kmsKeyArn` | `string \| valueFrom` |  |  | AwsKmsKey (`status.outputs.key_arn`) |
| `spec.clientBrokerEncryption` | `string` |  | `TLS` |  |
| `spec.inClusterEncryption` | `bool` |  | `true` |  |
| `spec.authentication` | `AwsMskClusterAuthentication` |  |  |  |
| `spec.authentication.saslIamEnabled` | `bool` |  |  |  |
| `spec.authentication.saslScramEnabled` | `bool` |  |  |  |
| `spec.authentication.tlsEnabled` | `bool` |  |  |  |
| `spec.authentication.tlsCertificateAuthorityArns` | `[]string \| valueFrom` |  |  |  |
| `spec.authentication.unauthenticated` | `bool` |  |  |  |
| `spec.scramSecretArns` | `[]string` |  |  |  |
| `spec.clusterPolicy` | `object` |  |  |  |
| `spec.configurationArn` | `string` |  |  |  |
| `spec.configurationRevision` | `int32` |  |  |  |
| `spec.serverProperties` | `map<string, string>` |  |  |  |
| `spec.logging` | `AwsMskClusterLogging` |  |  |  |
| `spec.logging.cloudwatchLogs` | `AwsMskClusterCloudwatchLogging` |  |  |  |
| `spec.logging.cloudwatchLogs.enabled` | `bool` |  |  |  |
| `spec.logging.cloudwatchLogs.logGroup` | `string \| valueFrom` |  |  | AwsCloudwatchLogGroup (`status.outputs.log_group_name`) |
| `spec.logging.firehose` | `AwsMskClusterFirehoseLogging` |  |  |  |
| `spec.logging.firehose.enabled` | `bool` |  |  |  |
| `spec.logging.firehose.deliveryStream` | `string \| valueFrom` |  |  | AwsKinesisFirehose (`status.outputs.delivery_stream_name`) |
| `spec.logging.s3` | `AwsMskClusterS3Logging` |  |  |  |
| `spec.logging.s3.enabled` | `bool` |  |  |  |
| `spec.logging.s3.bucket` | `string \| valueFrom` |  |  | AwsS3Bucket (`status.outputs.bucket_id`) |
| `spec.logging.s3.prefix` | `string` |  |  |  |
| `spec.enhancedMonitoring` | `string` |  |  |  |
| `spec.jmxExporterEnabled` | `bool` |  |  |  |
| `spec.nodeExporterEnabled` | `bool` |  |  |  |
| `spec.rebalancingStatus` | `string` |  |  |  |
| `spec.topics` | `[]AwsMskClusterTopic` |  |  |  |
| `spec.topics[].name` | `string` | yes |  |  |
| `spec.topics[].partitionCount` | `int32` | yes |  |  |
| `spec.topics[].replicationFactor` | `int32` | yes |  |  |
| `spec.topics[].configs` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.kafkaVersion

`string` · required

kafka_version is the Apache Kafka version for all brokers in the cluster.
Examples: "3.6.0", "3.5.1", "3.4.0", "2.8.1".
Upgrades are applied in place via rolling restart (an UpdateClusterKafkaVersion
operation); a version DOWNGRADE cannot be performed in place and forces cluster
replacement.

- rule: {"required":true}

### spec.numberOfBrokerNodes

`int32` · required

number_of_broker_nodes is the total number of Kafka broker nodes in the cluster.
Must be a multiple of the number of subnets provided in subnet_ids so that
brokers are evenly distributed across Availability Zones. Broker count can be
INCREASED in place (an UpdateBrokerCount operation); AWS does not support
decreasing the broker count.

- rule: {"required":true,"int32":{"gte":1}}

### spec.instanceType

`string` · required

instance_type determines the compute and memory capacity of each broker node.
Standard types: kafka.m5.large, kafka.m5.xlarge, kafka.m5.2xlarge, kafka.m5.4xlarge.
Graviton types: kafka.m7g.large, kafka.m7g.xlarge (better price-performance).
Express brokers: express.m7g.large and up (AWS-managed storage, faster scaling,
intelligent rebalancing -- see rebalancing_status).
Small/dev types: kafka.t3.small (no tiered storage, no public access, no
provisioned throughput).
Updatable in place: changing the instance type is a rolling UpdateBrokerType
operation, not a replacement.

- rule: {"required":true}

### spec.subnetIds

`[]string | valueFrom` · required

subnet_ids are the VPC subnets where broker nodes are placed.
Brokers are distributed round-robin across subnets. The number of broker nodes
must be a multiple of the number of subnets for even AZ distribution.
ForceNew: changing subnets forces cluster replacement.

- references: AwsSubnet (`status.outputs.subnet_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.securityGroupIds

`[]string | valueFrom` · required

security_group_ids are the security groups ATTACHED to the broker network
interfaces -- they define what can reach the brokers. Ingress rules for the
Kafka listener ports (9092 plaintext, 9094 TLS, 9096 SASL/SCRAM, 9098
SASL/IAM) and ZooKeeper (2181-2182) belong on these referenced
AwsSecurityGroup nodes.
ForceNew: adding or removing entries after creation forces cluster replacement.

- references: AwsSecurityGroup (`status.outputs.security_group_id`)
- rule: {"repeated":{"minItems":"1"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.publicAccessType

`string`

public_access_type controls whether the cluster is reachable from the public internet.
"DISABLED" (default): brokers are only reachable within the VPC.
"SERVICE_PROVIDED_EIPS": AWS assigns public IPs to brokers.
AWS only allows turning public access ON for an EXISTING cluster: creation always
starts DISABLED, and the provider applies SERVICE_PROVIDED_EIPS as a follow-up
connectivity update. Public access also requires real client authentication
(SASL/IAM, SASL/SCRAM, or mTLS -- unauthenticated must be off) and TLS-only
client_broker_encryption.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["DISABLED","SERVICE_PROVIDED_EIPS"]}}

### spec.vpcConnectivity

`AwsMskClusterVpcConnectivity`

vpc_connectivity enables multi-VPC private connectivity (AWS PrivateLink) so
clients in OTHER VPCs or accounts connect to the brokers without peering or
public exposure. Enable at least one authentication scheme here; each scheme
requires the same scheme to be enabled in `authentication` (an IAM-only
cluster cannot offer SCRAM over PrivateLink). AWS activates PrivateLink
connectivity as a follow-up update after the cluster is created; the consumer
side (an aws_msk_vpc_connection in the client VPC) is a separate surface.

### spec.vpcConnectivity.saslIamEnabled

`bool`

sasl_iam_enabled offers SASL/IAM authentication to PrivateLink clients.

### spec.vpcConnectivity.saslScramEnabled

`bool`

sasl_scram_enabled offers SASL/SCRAM authentication to PrivateLink clients.

### spec.vpcConnectivity.tlsEnabled

`bool`

tls_enabled offers mutual TLS authentication to PrivateLink clients.

### spec.networkType

`string`

network_type selects the IP addressing of the broker network interfaces.
"IPV4" (default): IPv4-only.
"DUAL": dual-stack IPv4 + IPv6 (requires dual-stack subnets).
One-way: AWS only supports updating IPV4 -> DUAL; going back forces replacement.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["IPV4","DUAL"]}}

### spec.ebsVolumeSizeGib

`int32` · optional (explicit presence)

ebs_volume_size_gib is the size of the EBS volume per broker, in GiB.
Range: 1-16384. If omitted, AWS uses the instance-type-specific default.
Volume size can be increased in place (an UpdateBrokerStorage operation).
Not applicable to express.* instance types (AWS manages Express broker storage).

- rule: {"int32":{"lte":16384,"gte":1}}

### spec.provisionedThroughputEnabled

`bool`

provisioned_throughput_enabled enables provisioned EBS throughput for higher streaming performance.
Only supported on kafka.m5.4xlarge and larger instance types with ebs_volume_size_gib >= 10 GiB.

### spec.provisionedThroughputMbs

`int32`

provisioned_throughput_mbs is the provisioned EBS throughput in MiB/s per broker.
Range: 250-2375. Required when provisioned_throughput_enabled is true.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"lte":2375,"gte":250}}

### spec.storageMode

`string`

storage_mode controls the data storage strategy.
"LOCAL" (default): all data on broker EBS volumes.
"TIERED": hot data on EBS, warm data automatically offloaded to low-cost storage.
Tiered storage requires Kafka 2.8.2.tiered+ and supported instance types
(not kafka.t3.small). One-way: moving TIERED -> LOCAL forces cluster replacement.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["LOCAL","TIERED"]}}

### spec.kmsKeyArn

`string | valueFrom`

kms_key_arn is the KMS key ARN for encrypting data at rest on broker EBS volumes.
If omitted, AWS uses the default aws/msk service key.
ForceNew: changing the KMS key forces cluster replacement.

- references: AwsKmsKey (`status.outputs.key_arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKmsKey, name: <that resource's name>, fieldPath: status.outputs.key_arn}} -- a bare string does not parse

### spec.clientBrokerEncryption

`string` · optional (explicit presence)

client_broker_encryption controls encryption for data in transit between clients and brokers.
"TLS" (default, recommended): all client-broker traffic is TLS-encrypted (port 9094).
"TLS_PLAINTEXT": both TLS (9094) and plaintext (9092) are available.
"PLAINTEXT": all client-broker traffic is unencrypted (port 9092).

- default: `TLS`
- rule: {"string":{"in":["TLS","TLS_PLAINTEXT","PLAINTEXT"]}}

### spec.inClusterEncryption

`bool` · optional (explicit presence)

in_cluster_encryption enables TLS encryption for data in transit between brokers.
Strongly recommended for production. ForceNew: changing this forces cluster replacement.

- default: `true`

### spec.authentication

`AwsMskClusterAuthentication`

authentication configures client authentication methods for the cluster.
Multiple methods can be enabled simultaneously (e.g., SASL/IAM + TLS).
If no authentication is configured, the cluster accepts unauthenticated connections.

- rule: tls_certificate_authority_arns is required when tls_enabled is true

### spec.authentication.saslIamEnabled

`bool`

sasl_iam_enabled enables SASL/IAM authentication.
Recommended for most workloads. Clients authenticate using AWS IAM credentials
(access key, role assumption, or instance profiles). No password management required.
Brokers listen on port 9098 for SASL/IAM connections.

### spec.authentication.saslScramEnabled

`bool`

sasl_scram_enabled enables SASL/SCRAM-SHA-512 authentication.
Clients authenticate with username/password stored in AWS Secrets Manager.
Useful for non-AWS clients that cannot use IAM. Brokers listen on port 9096.
Associate the credential secrets via scram_secret_arns on the spec.

### spec.authentication.tlsEnabled

`bool`

tls_enabled enables mutual TLS (mTLS) authentication.
Clients present X.509 certificates signed by a private Certificate Authority.
Requires tls_certificate_authority_arns to specify trusted CAs.
Brokers listen on port 9094 for TLS connections.

### spec.authentication.tlsCertificateAuthorityArns

`[]string | valueFrom`

tls_certificate_authority_arns are the ARNs of ACM Private Certificate Authority resources
trusted for client certificate validation. Required when tls_enabled is true.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.authentication.unauthenticated

`bool`

unauthenticated allows clients to connect without any authentication.
Not recommended for production. Useful for development or when network-level
security (VPC + security groups) provides sufficient access control.

### spec.scramSecretArns

`[]string`

scram_secret_arns are AWS Secrets Manager secret ARNs associated with the
cluster for SASL/SCRAM username/password authentication. Each secret holds a
{"username": ..., "password": ...} JSON document, MUST be named with the
AmazonMSK_ prefix, and MUST be encrypted with a customer-managed KMS key
(AWS rejects secrets on the default aws/secretsmanager key). Associations are
a cluster-keyed setting (added/removed in place); the secrets themselves are
managed outside this resource.

- rule: {"repeated":{"unique":true,"items":{"string":{"pattern":"^arn:aws[a-zA-Z-]*:secretsmanager:[a-z0-9-]+:\\d{12}:secret:AmazonMSK_.+$"}}}}

### spec.clusterPolicy

`object`

cluster_policy is a resource-based IAM policy attached to the cluster, as a
structured policy document -- the mechanism behind cross-account PrivateLink
access (granting kafka:CreateVpcConnection and the Get*ForVpcConnection
actions to consumer principals). A cluster setting keyed by the cluster ARN,
updated in place. The IaC modules serialize this document to JSON for the
provider.

### spec.configurationArn

`string`

configuration_arn is the ARN of an externally managed MSK Configuration resource.
MSK Configurations hold Apache Kafka server.properties overrides (e.g., replication factor,
min ISR, log retention). Mutually exclusive with server_properties.

### spec.configurationRevision

`int32`

configuration_revision is the revision number of the external MSK Configuration.
Required when configuration_arn is set. Must be >= 1.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","int32":{"gte":1}}

### spec.serverProperties

`map<string, string>`

server_properties defines Apache Kafka server.properties overrides as key-value pairs.
When provided, a module-managed MSK Configuration resource is created and associated
with the cluster (the folded shape; bring an existing one via configuration_arn).
Common properties: auto.create.topics.enable, default.replication.factor,
min.insync.replicas, num.partitions, log.retention.hours, log.retention.bytes.
Mutually exclusive with configuration_arn.

### spec.logging

`AwsMskClusterLogging`

logging configures broker log delivery to one or more destinations.
All three destinations (CloudWatch Logs, Kinesis Data Firehose, S3) can be
enabled simultaneously for different operational workflows.

### spec.logging.cloudwatchLogs

`AwsMskClusterCloudwatchLogging`

cloudwatch_logs configures delivery of broker logs to CloudWatch Logs.

- rule: log_group is required when CloudWatch logging is enabled

### spec.logging.cloudwatchLogs.enabled

`bool`

enabled controls whether broker logs are delivered to CloudWatch Logs.

### spec.logging.cloudwatchLogs.logGroup

`string | valueFrom`

log_group is the CloudWatch Logs group that receives broker logs.
Required when enabled is true.

- references: AwsCloudwatchLogGroup (`status.outputs.log_group_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsCloudwatchLogGroup, name: <that resource's name>, fieldPath: status.outputs.log_group_name}} -- a bare string does not parse

### spec.logging.firehose

`AwsMskClusterFirehoseLogging`

firehose configures delivery of broker logs to a Kinesis Data Firehose stream.

- rule: delivery_stream is required when Firehose logging is enabled

### spec.logging.firehose.enabled

`bool`

enabled controls whether broker logs are delivered to Kinesis Data Firehose.

### spec.logging.firehose.deliveryStream

`string | valueFrom`

delivery_stream is the Firehose delivery stream that receives broker logs.
Required when enabled is true.

- references: AwsKinesisFirehose (`status.outputs.delivery_stream_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsKinesisFirehose, name: <that resource's name>, fieldPath: status.outputs.delivery_stream_name}} -- a bare string does not parse

### spec.logging.s3

`AwsMskClusterS3Logging`

s3 configures delivery of broker logs to an S3 bucket.

- rule: bucket is required when S3 logging is enabled

### spec.logging.s3.enabled

`bool`

enabled controls whether broker logs are delivered to S3.

### spec.logging.s3.bucket

`string | valueFrom`

bucket is the S3 bucket that receives broker logs.
Required when enabled is true.

- references: AwsS3Bucket (`status.outputs.bucket_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AwsS3Bucket, name: <that resource's name>, fieldPath: status.outputs.bucket_id}} -- a bare string does not parse

### spec.logging.s3.prefix

`string`

prefix is an optional S3 key prefix for log objects.
Example: "msk-logs/my-cluster/" produces keys like "msk-logs/my-cluster/broker-1/...".

### spec.enhancedMonitoring

`string`

enhanced_monitoring sets the level of CloudWatch metrics published by the cluster.
"DEFAULT": cluster-level and topic-level metrics.
"PER_BROKER": adds per-broker metrics.
"PER_TOPIC_PER_BROKER": adds per-topic-per-broker metrics.
"PER_TOPIC_PER_PARTITION": most granular, adds per-partition metrics.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["DEFAULT","PER_BROKER","PER_TOPIC_PER_BROKER","PER_TOPIC_PER_PARTITION"]}}

### spec.jmxExporterEnabled

`bool`

jmx_exporter_enabled enables the Prometheus JMX Exporter on all brokers.
When enabled, JMX metrics are available on port 11001 for Prometheus scraping.
Provides detailed JVM and Kafka broker metrics.

### spec.nodeExporterEnabled

`bool`

node_exporter_enabled enables the Prometheus Node Exporter on all brokers.
When enabled, host-level metrics (CPU, memory, disk, network) are available
on port 11002 for Prometheus scraping.

### spec.rebalancingStatus

`string`

rebalancing_status controls intelligent rebalancing on Express-broker clusters
(instance_type express.*): MSK automatically redistributes partitions when
brokers are added or removed. "ACTIVE" (the AWS default for Express clusters)
or "PAUSED". Not applicable to standard kafka.* instance types.

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["ACTIVE","PAUSED"]}}

### spec.topics

`[]AwsMskClusterTopic`

topics declares Apache Kafka topics managed WITH the cluster through the
MSK topic API (CreateTopic/UpdateTopic/DeleteTopic) -- no Kafka client,
bootstrap connectivity, or client authentication setup is needed; the
control plane manages them via IAM alone. Each entry is keyed by its
topic name and follows the cluster's lifecycle.

Declare the topics your applications depend on here (the chart-wiring
shape: a cluster and its contract topics deploy as one unit); leave
application-owned ad-hoc topics to the applications themselves, or set
auto.create.topics.enable in server_properties.

Deleting a declared topic entry deletes the topic and its data. Topic
deletion requires delete.topic.enable=true on the cluster (MSK's default;
only relevant if server_properties overrides it to false, in which case
destroy operations on declared topics hang until the server property is
restored).

- rule: a topic name cannot be '.' or '..'

### spec.topics[].name

`string` · required

name is the Kafka topic name. Letters, digits, '.', '_', and '-' only,
up to 249 characters (Kafka's own naming contract). Names beginning with
"__" are reserved for Kafka-internal topics (e.g. __consumer_offsets) --
avoid declaring them.
ForceNew: a Kafka topic cannot be renamed; changing the name replaces the
topic and its data.

- rule: {"required":true,"string":{"pattern":"^[a-zA-Z0-9._-]{1,249}$"}}

### spec.topics[].partitionCount

`int32` · required

partition_count is the number of partitions for the topic.
Kafka only supports INCREASING the partition count in place; declaring a
lower count than the topic currently has fails the update (decrease is
not possible without replacing the topic).

- rule: {"required":true,"int32":{"gte":1}}

### spec.topics[].replicationFactor

`int32` · required

replication_factor is the number of replicas per partition. Cannot exceed
the cluster's number_of_broker_nodes (CEL-enforced on the spec). Use 3 for
production durability, matching MSK's default.replication.factor guidance.
ForceNew: changing the replication factor replaces the topic.

- rule: {"required":true,"int32":{"gte":1}}

### spec.topics[].configs

`map<string, string>`

configs are topic-level Apache Kafka configuration overrides, e.g.
"retention.ms", "cleanup.policy", "min.insync.replicas",
"retention.bytes". Keys and values follow Kafka's topic-config vocabulary;
values are strings exactly as Kafka accepts them ("-1" for unlimited
retention.bytes). Updated in place via the MSK topic API.

## Validation Rules

- `provisioned_throughput_requires_mbs`: provisioned_throughput_mbs must be set (250-2375) when provisioned_throughput_enabled is true
- `configuration_mutual_exclusion`: configuration_arn and server_properties are mutually exclusive
- `configuration_revision_required`: configuration_revision is required (>= 1) when configuration_arn is set
- `scram_secrets_require_scram_auth`: scram_secret_arns requires authentication.sasl_scram_enabled to be true
- `vpc_connectivity_iam_requires_cluster_iam`: vpc_connectivity.sasl_iam_enabled requires authentication.sasl_iam_enabled on the cluster
- `vpc_connectivity_scram_requires_cluster_scram`: vpc_connectivity.sasl_scram_enabled requires authentication.sasl_scram_enabled on the cluster
- `vpc_connectivity_tls_requires_cluster_tls`: vpc_connectivity.tls_enabled requires authentication.tls_enabled on the cluster
- `public_access_requires_authenticated_tls`: public access requires TLS-only client_broker_encryption and at least one of SASL/IAM, SASL/SCRAM, or mTLS authentication with unauthenticated disabled
- `topic_names_unique`: topic names must be unique within the cluster
- `topic_replication_within_brokers`: a topic's replication_factor cannot exceed number_of_broker_nodes

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsMskCluster, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.cluster_arn` | `string` | cluster_arn is the Amazon Resource Name of the MSK cluster, used in IAM policies, event source mappings, and cross-service references. |
| `status.outputs.cluster_name` | `string` | cluster_name is the human-readable name of the cluster. |
| `status.outputs.cluster_uuid` | `string` | cluster_uuid is the unique identifier extracted from the cluster ARN. |
| `status.outputs.current_version` | `string` | current_version is the cluster version string, required for update operations. This value changes after each successful cluster modification. |
| `status.outputs.bootstrap_brokers` | `string` | bootstrap_brokers is a comma-separated list of broker endpoints for plaintext connections (port 9092). Empty when client_broker_encryption is "TLS" (plaintext disabled). |
| `status.outputs.bootstrap_brokers_tls` | `string` | bootstrap_brokers_tls is a comma-separated list of broker endpoints for TLS connections (port 9094). The primary endpoint when client_broker_encryption is "TLS" or "TLS_PLAINTEXT". |
| `status.outputs.bootstrap_brokers_sasl_iam` | `string` | bootstrap_brokers_sasl_iam is a comma-separated list of broker endpoints for SASL/IAM connections (port 9098). Populated when authentication.sasl_iam_enabled is true. |
| `status.outputs.bootstrap_brokers_sasl_scram` | `string` | bootstrap_brokers_sasl_scram is a comma-separated list of broker endpoints for SASL/SCRAM connections (port 9096). Populated when authentication.sasl_scram_enabled is true. |
| `status.outputs.bootstrap_brokers_public_tls` | `string` | bootstrap_brokers_public_tls is a comma-separated list of public TLS broker endpoints. Populated when public_access_type is "SERVICE_PROVIDED_EIPS". |
| `status.outputs.bootstrap_brokers_public_sasl_iam` | `string` | bootstrap_brokers_public_sasl_iam is a comma-separated list of public SASL/IAM broker endpoints. Populated when public access and SASL/IAM authentication are both enabled. |
| `status.outputs.bootstrap_brokers_public_sasl_scram` | `string` | bootstrap_brokers_public_sasl_scram is a comma-separated list of public SASL/SCRAM broker endpoints. Populated when public access and SASL/SCRAM authentication are both enabled. |
| `status.outputs.bootstrap_brokers_vpc_connectivity_tls` | `string` | bootstrap_brokers_vpc_connectivity_tls is a comma-separated list of PrivateLink broker endpoints for mutual-TLS connections. Populated when vpc_connectivity.tls_enabled is true. |
| `status.outputs.bootstrap_brokers_vpc_connectivity_sasl_iam` | `string` | bootstrap_brokers_vpc_connectivity_sasl_iam is a comma-separated list of PrivateLink broker endpoints for SASL/IAM connections. Populated when vpc_connectivity.sasl_iam_enabled is true. |
| `status.outputs.bootstrap_brokers_vpc_connectivity_sasl_scram` | `string` | bootstrap_brokers_vpc_connectivity_sasl_scram is a comma-separated list of PrivateLink broker endpoints for SASL/SCRAM connections. Populated when vpc_connectivity.sasl_scram_enabled is true. |
| `status.outputs.zookeeper_connect_string` | `string` | zookeeper_connect_string is a comma-separated list of ZooKeeper endpoints for plaintext connections. Empty on KRaft-mode clusters (Kafka versions without ZooKeeper). |
| `status.outputs.zookeeper_connect_string_tls` | `string` | zookeeper_connect_string_tls is a comma-separated list of ZooKeeper endpoints for TLS connections. Empty on KRaft-mode clusters (Kafka versions without ZooKeeper). |
| `status.outputs.configuration_arn` | `string` | configuration_arn is the ARN of the module-managed MSK Configuration resource, if one was created from server_properties in the spec. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.subnetIds` | AwsSubnet | `status.outputs.subnet_id` |
| `spec.securityGroupIds` | AwsSecurityGroup | `status.outputs.security_group_id` |
| `spec.kmsKeyArn` | AwsKmsKey | `status.outputs.key_arn` |
| `spec.logging.cloudwatchLogs.logGroup` | AwsCloudwatchLogGroup | `status.outputs.log_group_name` |
| `spec.logging.firehose.deliveryStream` | AwsKinesisFirehose | `status.outputs.delivery_stream_name` |
| `spec.logging.s3.bucket` | AwsS3Bucket | `status.outputs.bucket_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsKinesisFirehose | `spec.mskSource.mskClusterArn` | `status.outputs.cluster_arn` |

## See Also

- [Overview](../README.md)
