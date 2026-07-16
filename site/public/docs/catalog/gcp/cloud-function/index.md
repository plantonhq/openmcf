---
title: "Cloud Function"
description: "Cloud Function deployment documentation"
icon: "package"
order: 100
componentName: "gcpcloudfunction"
---

# GCP Cloud Function

Source-based serverless compute on Cloud Run: ship a zip, Cloud Build containerizes it, Cloud Run serves it — HTTP endpoints or CloudEvent consumers via Eventarc.

**Enum:** 602 · **ID prefix:** `cldfunc` · **Provider:** GCP · **API:** `gcp.planton.dev/v1`

## At a Glance

| | |
|---|---|
| **Creates** | `google_cloudfunctions2_function` (+ public-invoker IAM member when requested) |
| **Triggers** | HTTPS, or 125+ CloudEvent sources via Eventarc |
| **Composes with** | GcpGcsBucket (source), GcpServerlessVpcConnector (private egress), GcpPubSubTopic (events), GcpServiceAccount (identity), GcpKmsKey (CMEK) |
| **Engines** | Terraform (~> 6.0) and Pulumi |

## When to Use

- **Webhooks and small HTTP APIs** — public or IAM-gated HTTPS endpoints from source
- **Event processors** — Pub/Sub, Storage, Firestore events through Eventarc
- **Database-backed glue** — private egress through a Serverless VPC Access connector, credentials from Secret Manager
- **LB-fronted functions** — `ALLOW_INTERNAL_AND_GCLB` + a serverless NEG puts the function behind Cloud Armor / CDN

## Quick Example

```yaml
apiVersion: gcp.planton.dev/v1
kind: GcpCloudFunction
metadata:
  name: hello-api
spec:
  region: us-central1
  buildConfig:
    runtime: nodejs22
    entryPoint: helloHttp
    source:
      storageSource:
        bucket:
          value: my-source-bucket
        object: functions/hello-v1.zip
  serviceConfig:
    allowUnauthenticated: true
```

## Key Fields

- `buildConfig` — runtime, entry point, GCS or repo source, build identity, base-image update policy
- `serviceConfig` — memory/CPU/timeout/concurrency, env + Secret Manager references (env vars and volume files), VPC connector egress, ingress, scaling, traffic pinning
- `trigger` — HTTP (default) or an Eventarc event type with filters and retry policy

## Outputs

`functionId`, `name`, `functionUrl`, `uri`, `cloudRunServiceId`, `serviceAccountEmail`, `eventarcTriggerId`, `state`, `environment`, `updateTime`

## See Also

- [README](README.md) — full configuration reference
- [GcpServerlessVpcConnector](/docs/catalog/gcp/serverless-vpc-access-connector) — the private-egress bridge
- Presets: [HTTP API](presets/01-http-api.yaml), [Pub/Sub processor](presets/02-pubsub-event.yaml), [private VPC egress](presets/03-private-vpc-egress.yaml)
