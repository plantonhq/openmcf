---
title: "Cloud Run behind a Load Balancer"
description: "The most common serverless NEG: put a Cloud Run service behind a global external Application Load Balancer so it can serve a custom domain with Cloud CDN, Cloud Armor, and IAP in front of it."
type: "preset"
rank: "01"
presetSlug: "01-cloud-run-neg"
componentSlug: "region-network-endpoint-group"
componentTitle: "Region Network Endpoint Group"
provider: "gcp"
icon: "package"
order: 1
---

# Cloud Run behind a Load Balancer

The most common serverless NEG: put a Cloud Run service behind a global external Application Load Balancer so it can serve a custom domain with Cloud CDN, Cloud Armor, and IAP in front of it.

## When to Use

- A Cloud Run service that needs a custom domain on a shared load balancer (rather than the default `run.app` URL)
- Cloud Run behind Cloud CDN, Cloud Armor, or Identity-Aware Proxy
- Blue/green or canary across Cloud Run revisions (add a `tag`)

## Remix Notes

- Reference a `GcpCloudRun` resource under `cloudRun.service.valueFrom` instead of a literal name to wire the NEG to a service Planton manages.
- The NEG must be in the **same region** as the Cloud Run service.
- Reference this NEG's `self_link` from a `GcpBackendService` backend's `group`; a backend service whose backends are all serverless NEGs needs **no health check**.
- To fan a wildcard domain out to many services from one NEG, drop `service` and set a `urlMask` like `<service>.example.com`.
