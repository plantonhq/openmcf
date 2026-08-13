# GCP Pub/Sub Topic

Deploys a Pub/Sub topic in a GCP project with configurable message retention, regional storage policies, schema validation, CMEK encryption, and data ingestion from external sources (AWS Kinesis, AWS MSK, Azure Event Hubs, Cloud Storage, Confluent Cloud). The component integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects, KMS keys, and GCS buckets.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Pub/Sub Topic** -- a named topic resource in the specified GCP project for publishing and distributing messages to subscribers
- **Message Storage Policy** -- created only when `messageStoragePolicy` is configured; restricts which GCP regions may persist messages, with optional in-transit enforcement that rejects publishes from non-allowed regions
- **Schema Validation** -- created only when `schemaSettings` is configured; validates all published messages against a Pub/Sub schema in JSON or BINARY encoding
- **Ingestion Pipeline** -- created only when `ingestionDataSourceSettings` is configured; ingests data from AWS Kinesis, AWS MSK, Azure Event Hubs, Cloud Storage, or Confluent Cloud into the topic
- **CMEK Encryption** -- created only when `kmsKeyName` is provided; encrypts messages at rest using a customer-managed Cloud KMS key instead of Google-managed keys
- **GCP Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the topic will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **Cloud Pub/Sub API** enabled in the target project.
- **Cloud KMS CryptoKey Encrypter/Decrypter role** granted to the Pub/Sub service account on the KMS key (if using CMEK encryption).
- **GCS bucket** accessible to the Pub/Sub service account (if using Cloud Storage ingestion).

## Deploy

### Console

Open the deployment store, find **GCP Pub/Sub Topic**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Basic Topic** preset in the [Presets](#presets) tab to pre-populate a minimal configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpPubSubTopic
metadata:
  name: order-events
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  topicName: order-events
```

```shell
planton apply -f pubsub-topic.yaml
```

This creates a topic with Google-managed encryption, no message retention, no regional storage constraints, and no schema validation. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the topic to a GCP project and KMS key deployed in the same InfraPipeline:

```yaml
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: production-project
      fieldPath: status.outputs.project_id
  kmsKeyName:
    valueFrom:
      kind: GcpKmsKey
      name: pubsub-encryption-key
      fieldPath: status.outputs.key_id
```

The InfraPipeline resolves the dependency graph, deploys the project and KMS key first, then provisions the topic with CMEK encryption.

## Key Configuration

These are the most important decisions when configuring a Pub/Sub topic. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Topic name** -- The `topicName` field is immutable after creation. Changing it requires destroying and recreating the topic along with all attached subscriptions. Choose a descriptive, stable name.

**Message retention** -- Set `messageRetentionDuration` (e.g., `"604800s"` for 7 days) to retain published messages at the topic level independently of subscription-level retention. Enables any subscription to seek to a timestamp within the retention window. Range: 600s to 2678400s.

**Regional storage policy** -- Configure `messageStoragePolicy.allowedPersistenceRegions` to restrict where messages are stored. Set `enforceInTransit: true` to additionally reject publish calls from non-allowed regions. Required for data residency compliance.

**Schema validation** -- Set `schemaSettings.schema` to a GcpPubSubSchema reference (its `schema_id` output) or a fully qualified literal path to validate all published messages. Invalid messages are rejected at publish time. Encoding must be `JSON` or `BINARY`. If the referenced schema is later deleted, the topic validates against the `_deleted-schema_` sentinel and publishes fail -- destroy topics before their schema.

**Data ingestion** -- Configure `ingestionDataSourceSettings` to stream data from external sources (AWS Kinesis, AWS MSK, Azure Event Hubs, Cloud Storage, or Confluent Cloud) directly into the topic without writing custom publisher code.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpKmsKey** (optional) | `kmsKeyName` | `status.outputs.key_id` |
| **GcpPubSubSchema** (optional) | `schemaSettings.schema` | `status.outputs.schema_id` |
| **GcpGcsBucket** (optional) | `ingestionDataSourceSettings.cloudStorage.bucket` | `status.outputs.bucket_id` |
| **GcpServiceAccount** (optional) | ingestion `gcpServiceAccount` fields | `status.outputs.email` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `topic_id` | Fully qualified topic ID (`projects/{project}/topics/{name}`) | Subscription topic references, Cloud Function triggers, Cloud Scheduler targets |
| `topic_name` | Short topic name | Display, logging, monitoring filters |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Basic topic** -- Minimal topic with Google-managed encryption and no retention or storage policy. Suitable for development and simple event routing. Start from the **Basic Topic** preset.

**Regional encrypted** -- CMEK-encrypted topic with regional storage constraints and in-transit enforcement. Suitable for regulated workloads with data residency requirements. Start from the **Regional Encrypted** preset.

**Message retention** -- Topic with 7-day message retention enabling subscription-level seek and replay. Suitable for event sourcing and audit trail patterns where consumers may need to reprocess historical messages. Start from the **Message Retention** preset.

**Cloud Storage ingestion** -- Topic configured to ingest data from a GCS bucket with glob-based filtering and text format parsing. Suitable for batch-to-stream bridge patterns where files landing in a bucket should trigger downstream event processing. Start from the **Cloud Storage Ingestion** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the topic is created
- [**GCP Pub/Sub Schema**](/cloud-catalog/gcp-pub-sub-schema) -- provides the message contract enforced at publish time
- [**GCP Pub/Sub Subscription**](/cloud-catalog/gcp-pub-sub-subscription) -- consumes this topic's stream (pull, push, BigQuery, or Cloud Storage delivery)
- [**GCP KMS Key**](/cloud-catalog/gcp-kms-key) -- provides the CMEK encryption key for messages at rest
- [**GCP GCS Bucket**](/cloud-catalog/gcp-gcs-bucket) -- provides the source bucket for Cloud Storage ingestion