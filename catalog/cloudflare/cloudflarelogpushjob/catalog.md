# Cloudflare Logpush Job

A Logpush job: continuous delivery of one Cloudflare log dataset to storage or an analytics endpoint you own. Real create, update, delete.

## What Gets Created

When you deploy this resource, the IaC module provisions:

- **Logpush job** -- one `cloudflare_logpush_job`
- **Ownership challenge** -- one `cloudflare_logpush_ownership_challenge`, only when `generateOwnershipChallenge` is set

## Prerequisites

- **A Logpush entitlement** -- most datasets are Enterprise-gated
- **A destination you control** (R2, S3-compatible, GCS, Splunk, Datadog, HTTPS) with write credentials
- **A Cloudflare API token** with Logs → Write on the chosen scope

## Quick Start

```yaml
apiVersion: cloudflare.planton.dev/v1alpha1
kind: CloudflareLogpushJob
metadata:
  name: http-logs-to-r2
  org: acme-corp
  env: prod
spec:
  zoneId:
    value: "023e105f4ecef8ad9ca31a8372d0c353"
  dataset: http_requests
  destinationConf:
    value: "r2://logs/{DATE}?account-id=...&access-key-id=...&secret-access-key=..."
  name: http-logs-to-r2
```

```shell
planton apply -f logpush-job.yaml
```

## Configuration Reference

### Required Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `dataset` | string | Which log dataset ships. | One of 35 provider values (enum-walled); immutable. |
| `destinationConf` | string | Destination URI with embedded credentials. | Required, sensitive; rejected on update by the API. |
| `accountId` / `zoneId` | string | The scope. | Exactly one must be set. |

### Optional Fields

| Field | Type | Description | Validation |
|-------|------|-------------|------------|
| `enabled` | bool | Whether logs ship. | Defaults to true here (Cloudflare's default is false). |
| `filter` | string | JSON filter expression. | Free-form Cloudflare filter grammar. |
| `kind` | string | Job variant. | Empty (batching) or `edge` (Instant Logs). |
| `maxUploadBytes` | int64 | Batch size ceiling. | 0, or 5000000–1000000000. |
| `maxUploadIntervalSeconds` | int64 | Batch interval ceiling. | 0, or 30–300. |
| `maxUploadRecords` | int64 | Batch record ceiling. | 0, or 1000–1000000. |
| `outputOptions` | object | Record shaping. | `outputType` ndjson/csv; `timestampFormat` one of five; `sampleRate` 0–1. |
| `ownershipChallenge` | string | Destination ownership token. | Sensitive, write-only. |
| `generateOwnershipChallenge` | bool | Issue the challenge file. | One-shot; never deleted at Cloudflare. |

## Destroy Semantics

Real delete for the job -- delivery stops immediately, and already-shipped objects remain in the destination. The ownership challenge is a one-shot record Cloudflare keeps forever; destroying it only drops it from state.

## Stack Outputs

| Output | Type | Description |
|--------|------|-------------|
| `jobId` | string | The Cloudflare-assigned numeric job id, in string form |
| `accountId` | string | The account scope, when account-scoped |
| `zoneId` | string | The zone scope, when zone-scoped |
| `ownershipChallengeFilename` | string | Challenge file name (challenge arm only) |
| `ownershipChallengeMessage` | string | Message accompanying the challenge |
| `ownershipChallengeValid` | bool | Whether Cloudflare found the destination valid |

## Related Components

- [Cloudflare DNS Zone](/docs/catalog/cloudflare/cloudflarednszone) -- the scope for zone datasets
- [Cloudflare R2 Bucket](/docs/catalog/cloudflare/cloudflarer2bucket) -- the usual same-account destination
- [Cloudflare Notification Policy](/docs/catalog/cloudflare/cloudflarenotificationpolicy) -- alert when a job starts failing
