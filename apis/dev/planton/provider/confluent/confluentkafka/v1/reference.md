# ConfluentKafka

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `confluent.planton.dev/v1`

confluent-kafka spec

## Example

```yaml
---
# Example ConfluentKafka manifest for testing deployments
# This file demonstrates a Standard production cluster configuration

apiVersion: confluent.planton.dev/v1
kind: ConfluentKafka
metadata:
  name: test-kafka-cluster
  labels:
    environment: test
    team: platform
    purpose: testing
spec:
  # Cloud provider (AWS, AZURE, or GCP)
  cloud: AWS
  
  # Cloud-specific region
  region: us-east-2
  
  # Availability configuration
  # SINGLE_ZONE: Development/testing, no SLA
  # MULTI_ZONE: Production, 99.99% SLA
  availability: MULTI_ZONE
  
  # Confluent Cloud environment ID
  # Replace with your actual environment ID
  environment_id: env-test-abc123
  
  # Cluster type (BASIC, STANDARD, ENTERPRISE, DEDICATED)
  # Default: STANDARD if not specified
  cluster_type: STANDARD
  
  # Optional: Custom display name for Confluent Cloud UI
  # If not specified, defaults to metadata.name
  display_name: Test Kafka Cluster
  
  # Optional: Dedicated cluster configuration
  # Required when cluster_type is DEDICATED
  # dedicated_config:
  #   cku: 2  # Minimum 1 CKU
  
  # Optional: Network configuration for private networking
  # Only available for ENTERPRISE and DEDICATED cluster types
  # network_config:
  #   network_id: n-abc123  # Pre-created network resource ID
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.cloud` | `string` | yes |  |  |
| `spec.region` | `string` | yes |  |  |
| `spec.availability` | `string` | yes |  |  |
| `spec.environmentId` | `string` | yes |  |  |
| `spec.clusterType` | `string` |  |  |  |
| `spec.dedicatedConfig` | `ConfluentKafkaDedicatedConfig` |  |  |  |
| `spec.dedicatedConfig.cku` | `int32` | yes |  |  |
| `spec.networkConfig` | `ConfluentKafkaNetworkConfig` |  |  |  |
| `spec.networkConfig.networkId` | `string` | yes |  |  |
| `spec.displayName` | `string` |  |  |  |

## Field Details

### spec.cloud

`string` · required

cloud provider
https://www.pulumi.com/registry/packages/confluentcloud/api-docs/kafkacluster/#cloud_yaml

- rule: {"required":true,"string":{"in":["AWS","AZURE","GCP"]}}

### spec.region

`string` · required

region is the cloud-specific region where the cluster will be deployed (e.g., us-east-2, us-central1, eastus)
https://www.pulumi.com/registry/packages/confluentcloud/api-docs/kafkacluster/#region_yaml

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.availability

`string` · required

availability determines high availability configuration
SINGLE_ZONE: Development/testing, no SLA
MULTI_ZONE: Production, 99.99% SLA (required for Standard/Dedicated)
LOW and HIGH are legacy values for Basic clusters
https://www.pulumi.com/registry/packages/confluentcloud/api-docs/kafkacluster/#availability_yaml

- rule: {"required":true,"string":{"in":["SINGLE_ZONE","MULTI_ZONE","LOW","HIGH"]}}

### spec.environmentId

`string` · required

environment_id is the ID of the Confluent Cloud environment (parent container for clusters)
https://www.pulumi.com/registry/packages/confluentcloud/api-docs/kafkacluster/#environment_yaml

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.clusterType

`string`

cluster_type determines the deployment type and capabilities
BASIC: Multi-tenant, development/testing, single-zone only, public internet only
STANDARD: Multi-tenant, production, elastic scaling, public internet only
ENTERPRISE: Multi-tenant, production, elastic scaling, supports private networking
DEDICATED: Single-tenant, production, provisioned capacity (CKU), supports private networking
Default: STANDARD (if not specified)

- rule: {"ignore":"IGNORE_IF_ZERO_VALUE","string":{"in":["BASIC","STANDARD","ENTERPRISE","DEDICATED"]}}

### spec.dedicatedConfig

`ConfluentKafkaDedicatedConfig`

dedicated_config is required only when cluster_type is DEDICATED
configures provisioned capacity (CKU) for dedicated clusters

### spec.dedicatedConfig.cku

`int32` · required

cku (Confluent Kafka Units) is the provisioned capacity for dedicated clusters
Minimum: 1 CKU, can be scaled up/down but not to zero
https://www.pulumi.com/registry/packages/confluentcloud/api-docs/kafkacluster/#cku_yaml

- rule: {"required":true,"int32":{"gte":1}}

### spec.networkConfig

`ConfluentKafkaNetworkConfig`

network_config configures private networking (PrivateLink, VNet Peering, Private Service Connect)
Only available for ENTERPRISE and DEDICATED cluster types
Optional: If not specified, cluster uses public internet access

### spec.networkConfig.networkId

`string` · required

network_id is the ID of the Confluent Cloud network resource
Must be pre-created in the same environment
https://docs.confluent.io/cloud/current/networking/overview.html

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.displayName

`string`

display_name is the human-readable name shown in Confluent Cloud UI
Optional: If not specified, defaults to metadata.name

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ConfluentKafka, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.id` | `string` | The provider-assigned unique ID for this managed resource. https://www.pulumi.com/registry/packages/confluentcloud/api-docs/kafkacluster/#id_yaml |
| `status.outputs.bootstrap_endpoint` | `string` | The bootstrap endpoint used by Kafka clients to connect to the Kafka cluster. (e.g., SASL_SSL://pkc-00000.us-central1.gcp.confluent.cloud:9092). https://www.pulumi.com/registry/packages/confluentcloud/api-docs/kafkacluster/#bootstrapendpoint_yaml |
| `status.outputs.crn` | `string` | The Confluent Resource Name of the Kafka cluster, for example, crn://confluent.cloud/organization=1111aaaa-11aa-11aa-11aa-111111aaaaaa/environment=env-abc123/cloud-cluster=lkc-abc123. https://www.pulumi.com/registry/packages/confluentcloud/api-docs/kafkacluster/#rbaccrn_yaml |
| `status.outputs.rest_endpoint` | `string` | The REST endpoint of the Kafka cluster (e.g., https://pkc-00000.us-central1.gcp.confluent.cloud:443). https://www.pulumi.com/registry/packages/confluentcloud/api-docs/kafkacluster/#restendpoint_yaml |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
