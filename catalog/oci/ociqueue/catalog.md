# Queue on OCI

Deploys an Oracle Cloud Infrastructure Queue -- a fully managed, serverless message queue for asynchronous communication between decoupled services. The queue supports at-least-once delivery, configurable dead-letter routing, optional KMS encryption, large messages up to 512 KB, and consumer groups for partitioned consumption. The component integrates with Planton's Provider Connections for OCI credential management and supports ValueFromRef wiring to compartments and encryption keys.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Queue** -- a serverless message queue in the specified compartment with configurable retention, visibility timeout, polling timeout, dead-letter routing, and optional capabilities (large messages, consumer groups)
- **Freeform Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the queue

## Before You Deploy

### Planton Setup

- **OCI Provider Connection** -- an active connection in the Connect module with credentials for the target OCI tenancy. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### OCI Tenancy

- A compartment to place the queue in. Provide the compartment OCID directly or reference an OciCompartment Cloud Resource via ValueFromRef.
- For customer-managed encryption: an OCI KMS key. When omitted, Oracle-managed encryption is used.

## Deploy

### Console

Open the deployment store, find **Queue on OCI**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard with DLQ** preset in the [Presets](#presets) tab to pre-populate a queue with 7-day retention and dead-letter routing.

### CLI

```yaml
apiVersion: oci.planton.dev/v1
kind: OciQueue
metadata:
  name: order-events
  org: acme-corp
  env: prod
spec:
  compartmentId:
    value: "ocid1.compartment.oc1..example"
  retentionInSeconds: 604800
  visibilityInSeconds: 30
  timeoutInSeconds: 30
  deadLetterQueueDeliveryCount: 5
```

```shell
planton apply -f queue.yaml
```

This creates a queue with 7-day retention, 30-second visibility timeout, and dead-letter routing after 5 failed delivery attempts. Oracle-managed encryption is used. Large messages and consumer groups are not enabled.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the queue to a compartment and encryption key deployed in the same InfraPipeline:

```yaml
spec:
  compartmentId:
    valueFrom:
      kind: OciCompartment
      name: messaging
      fieldPath: status.outputs.compartmentId
  customEncryptionKeyId:
    valueFrom:
      kind: OciKmsKey
      name: messaging-key
      fieldPath: status.outputs.keyId
```

The InfraPipeline resolves the dependency graph, deploys the compartment and KMS key first, then provisions the queue with the resolved values.

## Key Configuration

These are the most important decisions when configuring a queue. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Retention period** -- Set `retentionInSeconds` to control how long unconsumed messages remain in the queue before automatic deletion. Defaults to 604800 (7 days). This field is immutable after creation -- changing it forces queue recreation. Choose a retention window that covers your longest expected consumer outage.

**Dead-letter routing** -- Set `deadLetterQueueDeliveryCount` to the number of failed delivery attempts before a message is moved to the dead-letter queue. Set to 0 to disable dead-letter routing entirely. The consumer group config can override this value per consumer group.

**Large messages** -- Set `isLargeMessagesEnabled` to `true` to accept messages up to 512 KB. This adds the LARGE_MESSAGES capability to the queue. When disabled, standard OCI message size limits apply.

**Consumer groups** -- Provide `consumerGroupConfig` to enable the CONSUMER_GROUPS capability for partitioned consumption. The primary consumer group is created with the queue. Set `isPrimaryEnabled`, `primaryDeadLetterQueueDeliveryCount`, and `primaryDisplayName` to control its behavior.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **OciCompartment** | `compartmentId` | `status.outputs.compartmentId` |
| **OciKmsKey** (optional) | `customEncryptionKeyId` | `status.outputs.keyId` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `queue_id` | OCID of the queue | Monitoring alarms, IAM policy scoping, application configuration |
| `messages_endpoint` | Endpoint URL for publishing and consuming messages | Application SDK clients, OCI Functions triggers |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard with DLQ** -- A production queue with 7-day retention, 30-second visibility timeout, and dead-letter routing after 5 failed attempts. Covers most asynchronous messaging use cases. Start from the **Standard with DLQ** preset.

**Encrypted large messages** -- A security-hardened queue with customer-managed KMS encryption, large message support (512 KB), and consumer groups enabled. Designed for regulated workloads with large payloads. Start from the **Encrypted Large Messages** preset.

## Works With

- [**Compartment on OCI**](/cloud-catalog/oci-compartment) -- provides the compartment that scopes this queue
- [**KMS Key on OCI**](/cloud-catalog/oci-kms-key) -- provides the customer-managed encryption key for message content