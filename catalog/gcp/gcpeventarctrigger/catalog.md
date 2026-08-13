# GCP Eventarc Trigger

Creates an Eventarc trigger — the routing rule "when THIS event happens, call THAT service": events matching the criteria (a Pub/Sub message, a Cloud Storage object change, an audit-log entry, a SaaS partner event) are delivered as CloudEvents to a Cloud Run service, a GKE service, a Workflow execution, or a private HTTP endpoint.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Eventarc trigger** -- an `eventarc.Trigger` with the configured criteria, destination, identity, and transport
- **Partner channel** (when `partnerChannel` is set) -- an `eventarc.Channel` the trigger is wired to, with its one-time activation token exported
- **Google channel config** (when `googleChannelCryptoKey` is set) -- the per-project-per-location CMEK singleton
- **Eventarc API enablement** -- `eventarc.googleapis.com` enabled in the target project (never disabled on destroy)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project.
- **Planton Runner** -- required when using Runner-based credential delivery.

### GCP Project

- **A GCP project** to host the trigger (directly or via a GcpProject reference).
- **IAM**: the deploying identity needs `roles/eventarc.admin` (or narrower); the trigger's `serviceAccount` needs `roles/eventarc.eventReceiver`, plus `roles/run.invoker` for authenticated Cloud Run destinations.

## Deploy

### Console

Open the deployment store, find **GCP Eventarc Trigger**, and click **Deploy**. Start from the **Pub/Sub to Cloud Run** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpEventarcTrigger
metadata:
  name: order-events
  org: acme-corp
  env: prod
spec:
  location: us-central1
  matchingCriteria:
    - attribute: type
      value: google.cloud.pubsub.topic.v1.messagePublished
  destination:
    cloudRunService:
      service:
        value: order-processor
  serviceAccount:
    value: eventarc-invoker@my-project.iam.gserviceaccount.com
```

```shell
planton apply -f trigger.yaml
```

### InfraChart

The event-routing backbone in one chart: a GcpPubSubTopic, this trigger consuming it as transport, and the GcpCloudRun destination — publish and the service is invoked.

## Key Configuration

**matchingCriteria** -- ALL criteria must match. Every trigger MUST filter the `type` attribute (the [event catalog](https://cloud.google.com/eventarc/docs/reference/supported-events) lists types and their filterable attributes). `match-path-pattern` is the only non-exact operator (Storage object names, audit-log resourceName).

**destination** -- exactly one arm: `cloudRunService` (the workhorse), `gke` (needs one-time `gcloud eventarc gke-destinations init` per project), `workflow` (references a GcpWorkflow), or `httpEndpoint` (private endpoints via a VPC network attachment).

**serviceAccount** -- the trigger's identity: `roles/eventarc.eventReceiver` always; `roles/run.invoker` for authenticated Cloud Run destinations; REQUIRED for audit-log triggers.

**transportPubsubTopic** -- messagePublished triggers only: consume an EXISTING topic instead of letting Eventarc mint one. The topic is never deleted with the trigger.

**partnerChannel** -- receive events from a SaaS partner: the module creates the channel and wires the trigger; hand the `partner_channel_activation_token` output to the partner to complete the handshake.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |
| **GcpCloudRun** (destination arm) | `destination.cloudRunService.service` | `status.outputs.service_name` |
| **GcpGkeCluster** (destination arm) | `destination.gke.cluster` | `status.outputs.name` |
| **GcpWorkflow** (destination arm) | `destination.workflow` | `status.outputs.workflow_id` |
| **GcpPubSubTopic** (optional) | `transportPubsubTopic` | `status.outputs.topic_id` |
| **GcpServiceAccount** (optional) | `serviceAccount` | `status.outputs.email` |
| **GcpKmsKey** (optional) | `partnerChannel.cryptoKey`, `googleChannelCryptoKey` | `status.outputs.key_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `trigger_name` | The trigger name | gcloud commands, debugging |
| `trigger_id` | Full trigger resource name | Monitoring filters, the canonical API handle |
| `partner_channel_activation_token` | One-time partner handshake token (sensitive) | Handed to the SaaS partner's console |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Pub/Sub to Cloud Run** -- the workhorse: a message published to a topic invokes a service. Start from the **Pub/Sub to Cloud Run** preset.

**Audit-log to Workflow** -- react to control-plane changes (a bucket created, an IAM grant) with an orchestrated response. Start from the **Audit Log to Workflow** preset.

## Works With

- [**GCP Cloud Run**](/cloud-catalog/gcp-cloud-run) -- the most common destination
- [**GCP Workflow**](/cloud-catalog/gcp-workflow) -- orchestrated event handling
- [**GCP GKE Cluster**](/cloud-catalog/gcp-gke-cluster) -- GKE service destinations
- [**GCP Pub/Sub Topic**](/cloud-catalog/gcp-pub-sub-topic) -- bring-your-own transport
- [**GCP Service Account**](/cloud-catalog/gcp-service-account) -- the trigger's identity
- [**GCP KMS Key**](/cloud-catalog/gcp-kms-key) -- channel CMEK
