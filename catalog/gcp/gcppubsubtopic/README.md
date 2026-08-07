# GcpPubSubTopic

A GcpPubSubTopic provisions a Pub/Sub topic -- the named channel to which
publishers send messages and from which subscribers consume via subscriptions.
Topics are the foundation of event-driven and streaming architectures in GCP,
decoupling producers from consumers and supporting one-to-many message delivery.

## When to Use

Use GcpPubSubTopic when you need:

- **An event bus** for asynchronous communication between microservices
- **A streaming ingestion endpoint** for real-time data pipelines (logs, metrics, events)
- **Cross-cloud data ingestion** from AWS Kinesis, AWS MSK, Azure Event Hubs, or Confluent Cloud
- **Cloud Storage ingestion** to stream object data into Pub/Sub messages
- **Customer-managed encryption** (CMEK) for topics handling sensitive or regulated data
- **Regional message storage guarantees** with configurable in-transit enforcement
- **Schema validation** to enforce message contracts at publish time (compose with `GcpPubSubSchema`)
- **Publish-side transforms** to redact, normalize, or filter messages at the topic boundary

## Prerequisites

- A GCP project (the Pub/Sub API is enabled automatically by both modules)
- Appropriate IAM permissions (`roles/pubsub.admin` or `roles/pubsub.editor`)
- For CMEK: an existing KMS key with `roles/cloudkms.cryptoKeyEncrypterDecrypter` granted to the Pub/Sub service account
- For schema validation: a `GcpPubSubSchema` resource to reference

## Quick Start

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpPubSubTopic
metadata:
  name: my-topic
spec:
  topicName: order-events
```

This creates a Pub/Sub topic with Google-managed encryption and no message
retention -- subscribers control their own retention via subscriptions.

## Configuration Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `projectId` | StringValueOrRef | No | GCP project ID (omitted = provider default project); reference `GcpProject` |
| `topicName` | string | Yes | Topic name (3-255 chars, starts with letter, no reserved `goog` prefix, immutable) |
| `kmsKeyName` | StringValueOrRef | No | CMEK encryption key; reference `GcpKmsKey` |
| `messageRetentionDuration` | string | No | Retain messages for this duration (600s-2678400s) |
| `labels` | map | No | User labels, merged beneath the platform attribution labels |
| `messageStoragePolicy` | object | No | Regional storage constraints |
| `messageStoragePolicy.allowedPersistenceRegions` | list | Yes* | GCP regions for message storage (* required when policy is set) |
| `messageStoragePolicy.enforceInTransit` | bool | No | Reject publishes from non-allowed regions |
| `schemaSettings` | object | No | Schema validation settings |
| `schemaSettings.schema` | StringValueOrRef | Yes* | Schema to validate against; reference `GcpPubSubSchema` (* required when schemaSettings is set) |
| `schemaSettings.encoding` | string | No | Message encoding: JSON or BINARY |
| `ingestionDataSourceSettings` | object | No | External data source ingestion |
| `ingestionDataSourceSettings.awsKinesis` | object | No | AWS Kinesis Data Streams ingestion (`gcpServiceAccount` references `GcpServiceAccount`) |
| `ingestionDataSourceSettings.awsMsk` | object | No | AWS MSK (Kafka) ingestion (`gcpServiceAccount` references `GcpServiceAccount`) |
| `ingestionDataSourceSettings.azureEventHubs` | object | No | Azure Event Hubs ingestion (`gcpServiceAccount` references `GcpServiceAccount`) |
| `ingestionDataSourceSettings.cloudStorage` | object | No | GCS bucket ingestion (`bucket` references `GcpGcsBucket`) |
| `ingestionDataSourceSettings.confluentCloud` | object | No | Confluent Cloud ingestion (`gcpServiceAccount` references `GcpServiceAccount`) |
| `ingestionDataSourceSettings.platformLogsSettings` | object | No | Ingestion pipeline logging |
| `messageTransforms` | list | No | Ordered JavaScript UDF pipeline applied to every published message |

## Important Notes

**Topic name is immutable.** Changing `topicName` destroys and recreates the
topic along with all its subscriptions. Choose names carefully.

**CMEK requires IAM setup.** The Pub/Sub service account
(`service-{PROJECT_NUMBER}@gcp-sa-pubsub.iam.gserviceaccount.com`) must have
`roles/cloudkms.cryptoKeyEncrypterDecrypter` on the KMS key before the topic
is created.

**Message retention is topic-level.** When `messageRetentionDuration` is set,
messages are retained regardless of subscriber acknowledgement. This enables
replay via subscription seek operations.

**Ingestion sources are mutually exclusive per topic.** Configure at most one
ingestion data source (Kinesis, MSK, Event Hubs, Cloud Storage, or Confluent
Cloud) per topic.

**Transforms run in list order.** Each JavaScript UDF receives the message and
returns the transformed message (or null/undefined to drop it); a `disabled`
transform keeps its position in the pipeline without being applied.

### Deliberately not modeled (recorded reasons)

- **`tags` (resource-manager tags)** — absent from the released `google ~> 6.x`
  line (present only on the provider's unreleased line); revisit when the
  catalog's provider line carries it.
- **`message_transforms.ai_inference` (Vertex AI transform arm)** — absent from
  the released `google ~> 6.x` line; only the JavaScript UDF arm is released.
- **`deletion_policy`** — a client-side lever that conflicts with
  Planton-managed destroy (catalog-wide skip; also absent from the released
  6.x line).
- **`schema_settings.first_revision_id`/`last_revision_id` (revision pinning)** —
  absent from the released `google ~> 6.x` line; topics validate against all
  available schema revisions.
- **Per-topic IAM (`google_pubsub_topic_iam_*`)** — resource-scoped IAM stays
  out of the catalog pending concrete pull (the additive project-level grant,
  `GcpProjectIamMember`, covers the real cases).

## Related Components

- [GcpPubSubSubscription](../gcppubsubsubscription/) -- Subscriptions that consume from this topic
- [GcpPubSubSchema](../gcppubsubschema/) -- Message contract referenced by `schemaSettings.schema`
- [GcpKmsKey](../gcpkmskey/) -- CMEK encryption key for topic encryption
- [GcpProject](../gcpproject/) -- Parent GCP project
- [GcpGcsBucket](../gcpgcsbucket/) -- Source bucket for Cloud Storage ingestion
- [GcpServiceAccount](../gcpserviceaccount/) -- Federated identity for cross-cloud ingestion

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
