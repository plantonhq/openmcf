# Cloudflare Workflow

A Workflow registration: the binding of a durable-execution class exported by a deployed Worker script to a named workflow in the account, with retention, step limits, and cron triggers. Instances run on the Workers platform; this resource manages only the registration.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Workflow registration** -- one `cloudflare_workflow` binding the name to the script's exported class

## Prerequisites

- **A deployed Worker script** exporting the `WorkflowEntrypoint` subclass (a `CloudflareWorker` resource)
- **A Cloudflare API token** with Account → Workers Scripts → Edit
- Production instance runs at scale ride **Workers Paid** limits; registration itself works on the free tier

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWorkflow
metadata:
  name: order-fulfillment
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  workflowName: order-fulfillment
  className: OrderFulfillment
  scriptName:
    valueFrom:
      kind: CloudflareWorker
      name: order-processor
      fieldPath: status.outputs.script_name
```

```shell
planton apply -f workflow.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `accountId` | string | The Cloudflare account. | Required, 32-hex. |
| `workflowName` | string | The workflow's identity. Create is a PUT: a name collision silently adopts the existing workflow. | Required; replaces on change. |
| `className` | string | The `WorkflowEntrypoint` subclass the script exports. | Required. |
| `scriptName` | StringValueOrRef | The deployed Worker script exporting the class. | Required; references a `CloudflareWorker`. |

### Optional Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `defaultRetention.errorRetention` | string | Cloudflare default | Retention for errored instances -- milliseconds ("5000") or a duration ("5 minutes"). |
| `defaultRetention.successRetention` | string | Cloudflare default | Retention for successful instances -- same forms. |
| `limits.steps` | int | Cloudflare default | Maximum steps per instance (>= 1). |
| `schedules[].cron` | string | none | Cron triggers (five-field, UTC). |

## Examples

### Nightly scheduled workflow with retention

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareWorkflow
metadata:
  name: nightly-reconcile
  org: acme-corp
  env: prod
spec:
  accountId: "0a1b2c3d4e5f6a7b8c9d0e1f2a3b4c5d"
  workflowName: nightly-reconcile
  className: NightlyReconcile
  scriptName:
    value: reconcile-worker
  defaultRetention:
    errorRetention: "7 days"
    successRetention: "86400000"
  limits:
    steps: 512
  schedules:
    - cron: "0 3 * * *"
```

## Destroy Semantics

Destroy is a real delete; the API keeps answering GET for the deleted workflow with an `is_deleted` marker instead of a 404. Running instances belong to the platform and are not awaited.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `workflow_name` | string | The workflow's identity -- what Worker workflow bindings reference |
| `version_id` | string | The workflow version the registration produced |

## Related Components

- [Cloudflare Worker](/docs/catalog/cloudflare/cloudflareworker) -- the script exporting the workflow class
- [Cloudflare Queue](/docs/catalog/cloudflare/cloudflarequeue) -- the sibling asynchronous-work building block
