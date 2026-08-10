---
title: "Workflow"
description: "Workflow deployment documentation"
icon: "package"
order: 100
componentName: "gcpworkflow"
---

# GCP Workflow

Creates a Cloud Workflows workflow — a serverless orchestrator that executes a sequence of steps (HTTP calls, GCP connector calls, conditionals, retries) defined in YAML or JSON. The glue between services: an Eventarc trigger fires an execution, the workflow calls Cloud Run / BigQuery / Pub/Sub in order, handles errors, and records the result.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Workflow** -- a `workflows.Workflow` with the configured source, service account, CMEK, logging levels, and env vars
- **Workflows API enablement** -- `workflows.googleapis.com` enabled in the target project (never disabled on destroy)

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project.
- **Planton Runner** -- required when using Runner-based credential delivery.

### GCP Project

- **A GCP project** to host the workflow (directly or via a GcpProject reference).
- **IAM**: the deploying identity needs `roles/workflows.editor` or broader, plus `iam.serviceAccounts.actAs` on the workflow's service account when one is set.

## Deploy

### Console

Open the deployment store, find **GCP Workflow**, and click **Deploy**. Start from the **HTTP Orchestrator** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpWorkflow
metadata:
  name: order-orchestrator
  org: acme-corp
  env: prod
spec:
  region: us-central1
  sourceContents: |
    main:
      params: [input]
      steps:
        - charge:
            call: http.post
            args:
              url: https://payments.internal/charge
              auth:
                type: OIDC
            result: charge_result
        - done:
            return: ${charge_result.body}
  serviceAccount:
    value: order-workflow@my-project.iam.gserviceaccount.com
```

```shell
planton apply -f workflow.yaml
```

### InfraChart

The event-driven backbone in one chart: a GcpPubSubTopic, a GcpEventarcTrigger whose `destination.workflow` references this workflow, and the workflow itself — publish a message and the orchestration runs.

## Key Configuration

**sourceContents** -- the [workflow definition](https://cloud.google.com/workflows/docs/reference/syntax) executed on each run (≤128KB — the API cap, enforced at manifest time). Every change deploys a NEW revision; running executions finish on the revision they started with.

**serviceAccount** -- the identity the workflow's HTTP and connector calls carry. Grant a dedicated account only the roles the steps need — the default compute account is fine for experiments, wrong for production.

**callLogLevel / executionHistoryLevel** -- how much call and step detail lands in Cloud Logging and execution history. `LOG_ALL_CALLS` + `EXECUTION_HISTORY_DETAILED` while developing; scale back for cost in steady state.

**userEnvVars** -- values visible via `sys.get_env()` (≤20 entries; keys must not start `GOOGLE` or `WORKFLOWS`). An env-var change deploys a new revision.

**deletionProtection** -- ON by default: destroy fails until you explicitly set `false`. The workflow's execution HISTORY dies with the workflow — see the guide.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** (optional) | `projectId` | `status.outputs.project_id` |
| **GcpServiceAccount** (optional) | `serviceAccount` | `status.outputs.email` |
| **GcpKmsKey** (optional) | `cryptoKey` | `status.outputs.key_id` |

### What This Component Provides

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `workflow_id` | Full resource name | A GcpEventarcTrigger's `destination.workflow`; a GcpEventarcMessageBus pipeline's `destination.workflow` |
| `workflow_name` | Short name | Display, gcloud commands |
| `revision_id` | Deployed revision | Deploy tracking — compare across applies |
| `state` | `ACTIVE` / `UNAVAILABLE` | Health checks |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**HTTP orchestrator** -- call services in order with per-step auth and retries. Start from the **HTTP Orchestrator** preset.

**Event-driven step function** -- an Eventarc trigger starts an execution per event; the workflow fans out the work. Start from the **Event Handler** preset.

## Works With

- [**GCP Eventarc Trigger**](/cloud-catalog/gcp-eventarc-trigger) -- starts an execution per matching event
- [**GCP Eventarc Message Bus**](/cloud-catalog/gcp-eventarc-message-bus) -- pipelines deliver messages into executions
- [**GCP Service Account**](/cloud-catalog/gcp-service-account) -- the identity the steps run as
- [**GCP KMS Key**](/cloud-catalog/gcp-kms-key) -- CMEK for workflow and execution data
- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project
