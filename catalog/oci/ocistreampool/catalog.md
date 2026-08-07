# Stream Pool on OCI

Deploys an Oracle Cloud Infrastructure Stream Pool -- the organizational container for OCI Streaming, a Kafka-compatible managed event-streaming service. The pool groups streams under shared Kafka compatibility settings, optional KMS encryption, and optional private networking. Streams are bundled as sub-resources within the pool and inherit the pool's encryption and endpoint configuration. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments, encryption keys, subnets, and security groups.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Stream Pool** -- the organizational container in the specified compartment with Kafka compatibility settings, optional KMS encryption, and optional private endpoint networking
- **Streams** -- one stream per entry in `streams`, each with a configurable partition count and retention period. Streams are created inside the pool and inherit its encryption and endpoint configuration.
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the pool and all streams

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the stream pool in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- For customer-managed encryption: an OCI KMS key. When omitted, Oracle-managed encryption is used.
- For private networking: a subnet and optionally network security groups. When omitted, the pool is accessible via public endpoints.

## Deploy

### Console

Open the deployment store, find **Stream Pool on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Public Kafka Compatible** preset in the [Presets](#presets) tab to pre-populate a pool with Kafka auto-topic creation, two streams, and 48-hour retention.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciStreamPool
metadata:
  name: app-events
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  kafkaSettings:
    autoCreateTopicsEnable: true
    logRetentionHours: 48
    numPartitions: 3
  streams:
    - name: events
      partitions: 5
      retentionInHours: 48
    - name: commands
      partitions: 3
      retentionInHours: 24
```

```shell
planton apply -f stream-pool.yaml
```

This creates a publicly accessible stream pool with Kafka auto-topic creation enabled, two streams (`events` with 5 partitions and `commands` with 3 partitions), and Oracle-managed encryption. No private endpoint or KMS encryption is configured.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the pool to a compartment, encryption key, and private networking deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: streaming
      fieldPath: status.outputs.compartmentId
  kmsKeyId:
    valueFrom:
      kind: OciKmsKey
      name: streaming-key
      fieldPath: status.outputs.keyId
  privateEndpointSettings:
    subnetId:
      valueFrom:
        kind: OciSubnet
        name: streaming-subnet
        fieldPath: status.outputs.subnetId
    nsgIds:
      - valueFrom:
          kind: OciSecurityGroup
          name: streaming-nsg
          fieldPath: status.outputs.networkSecurityGroupId
```

The InfraPipeline resolves the dependency graph, deploys the compartment, KMS key, subnet, and security group first, then provisions the stream pool with the resolved values.

## Key Configuration

These are the most important decisions when configuring a stream pool. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Kafka compatibility** -- Configure `kafkaSettings` to control the Kafka compatibility layer. Set `autoCreateTopicsEnable` to `true` for development (auto-creates topics on first publish) or `false` for production (prevents accidental topic creation). Set `logRetentionHours` (24-168) and `numPartitions` for auto-created topic defaults.

**Private endpoint** -- Provide `privateEndpointSettings` with a `subnetId` to make the pool accessible only from within the specified subnet. Add `nsgIds` to restrict traffic further. The entire private endpoint block is immutable after creation -- changing any field forces pool recreation.

**Stream design** -- Define streams in the `streams` list. Each stream requires a `name` and `partitions` count (minimum 1). Set `retentionInHours` between 24 and 168 (defaults to 24). Stream name, partition count, and retention are all immutable after creation.

**Encryption** -- Set `kmsKeyId` to use a customer-managed KMS key for encrypting all stream data. When omitted, Oracle-managed encryption is used. Unlike private endpoint settings, the encryption key can be changed after creation.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciKmsKey** (optional) | `kmsKeyId` | `status.outputs.keyId` |
| **OciSubnet** (optional) | `privateEndpointSettings.subnetId` | `status.outputs.subnetId` |
| **OciSecurityGroup** (optional) | `privateEndpointSettings.nsgIds` | `status.outputs.networkSecurityGroupId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `stream_pool_id` | OCID of the stream pool | IAM policy scoping, monitoring alarms, resource management |
| `endpoint_fqdn` | FQDN for accessing streams in the pool (resolves only within subnet for private pools) | Application SDK configuration, DNS records |
| `kafka_bootstrap_servers` | Kafka-compatible bootstrap server string | Kafka producer and consumer client configuration |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Public Kafka compatible** -- A publicly accessible pool with Kafka auto-topic creation, multiple streams, and 48-hour retention. Suitable for development and non-sensitive workloads where Kafka clients connect over the public endpoint. Start from the **Public Kafka Compatible** preset.

**Private encrypted** -- A production pool with KMS encryption, private endpoint networking, NSG restrictions, auto-topic creation disabled, maximum 7-day retention, and a dead-letter stream. Designed for regulated workloads requiring network isolation and customer-managed encryption. Start from the **Private Encrypted** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this stream pool
- [**KMS Key on OCI**](/cloud-catalog/oci-kms-key) -- provides the customer-managed encryption key for stream data
- [**Subnet on OCI**](/cloud-catalog/oci-subnet) -- provides the subnet for private endpoint access
- [**Network Security Group on OCI**](/cloud-catalog/oci-security-group) -- provides network security rules for the private endpoint