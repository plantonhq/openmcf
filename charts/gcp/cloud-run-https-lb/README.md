# GCP Cloud Run HTTPS Load Balancer

A production HTTPS front door for Cloud Run: a global external load balancer with a static IPv4 address, a Google-managed TLS certificate, and modern EXTERNAL_MANAGED routing — the composition GCP itself recommends when a Cloud Run service outgrows domain mappings. Optional Cloud Armor protection and Cloud DNS publishing ride the same graph. The chart deploys a hello service by default so its defaults produce a working stack; point it at your own image, or front a Cloud Run service you already run.

## What it deploys

| Resource | Kind | Purpose |
|----------|------|---------|
| Cloud Run service | `GcpCloudRun` | The workload (optional — `deployService: false` fronts an existing service instead) |
| Static global IP | `GcpGlobalAddress` | The stable address DNS points at |
| Managed certificate | `GcpManagedSslCertificate` | Google-issued and auto-renewed TLS for the domain |
| Serverless NEG | `GcpRegionNetworkEndpointGroup` | Bridges Cloud Run into the load balancer |
| Cloud Armor policy | `GcpCloudArmorPolicy` | Optional (`cloudArmorEnabled`) — Layer-7 DDoS defense on the backend |
| Backend service | `GcpBackendService` | Fronts the NEG; where Cloud Armor attaches |
| URL map | `GcpUrlMap` | The routing table (default-service; grows path routing later) |
| Target HTTPS proxy | `GcpTargetHttpsProxy` | TLS termination with the managed certificate |
| Global forwarding rule | `GcpGlobalForwardingRule` | The public 443 entry point on the static IP |
| DNS A record | `GcpDnsRecord` | Optional (`dnsEnabled`) — publishes the domain at the IP in Cloud DNS |

## Architecture

```
                      Internet
                         │
                 https://app.example.com
                         │
        ┌────────────────▼────────────────┐
        │  GcpGlobalForwardingRule :443   │   Layer 1 — the front door
        │  (static IP: GcpGlobalAddress)  │
        └────────────────┬────────────────┘
        ┌────────────────▼────────────────┐
        │  GcpTargetHttpsProxy            │   Layer 2 — TLS termination
        │  (GcpManagedSslCertificate)     │
        └────────────────┬────────────────┘
        ┌────────────────▼────────────────┐
        │  GcpUrlMap → GcpBackendService  │   Layer 3 — routing (+ Cloud
        │  (optional GcpCloudArmorPolicy) │   Armor when enabled)
        └────────────────┬────────────────┘
        ┌────────────────▼────────────────┐
        │  Serverless NEG → GcpCloudRun   │   Layer 4 — the workload
        │  (ingress: INTERNAL_LB only)    │
        └─────────────────────────────────┘
```

The chart-deployed service accepts traffic ONLY through the load balancer (`ingress: INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER`), so the default `run.app` URL cannot bypass the domain, the certificate, or Cloud Armor.

## Parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `appName` | `web-app` | Base name for every resource (env-prefixed) |
| `region` | `us-central1` | Region of the Cloud Run service and its NEG |
| `domain` | `app.example.com` | The domain the LB serves; the certificate is issued for exactly this name |
| `deployService` | `true` | Deploy the chart's own Cloud Run service; `false` fronts an existing one |
| `existingServiceName` | `my-existing-service` | The existing service to front when `deployService: false` |
| `containerImage` | Google hello image | Image for the chart-deployed service |
| `cloudArmorEnabled` | `false` | Create and attach a Cloud Armor policy (L7 DDoS defense) |
| `dnsEnabled` | `false` | Publish the domain's A record in Cloud DNS |
| `dnsZoneName` | `my-dns-zone` | The `GcpDnsZone` resource hosting the domain (when `dnsEnabled`) |

## After deployment

1. **Point DNS at the load balancer.** With `dnsEnabled: true` the A record is already published in your Cloud DNS zone; otherwise read the static IP from the `GcpGlobalAddress` output (`address`) and create the A record at your DNS host.
2. **Wait for the certificate.** The Google-managed certificate provisions only after the domain resolves to the IP — typically within 15 minutes of DNS propagating. Until then, HTTPS for the domain returns an SSL error while the LB itself is healthy.
3. **Verify end to end**: `curl -I https://<domain>` should answer from Cloud Run (check the `server` header). The `run.app` URL of the chart-deployed service refuses external traffic by design.
4. **Bring-your-own service?** Set that service's ingress to `INGRESS_TRAFFIC_INTERNAL_LOAD_BALANCER` yourself — the chart never mutates a service it does not own — or leave it open if you want both entry points during a migration.

## Day-2 notes

- **Path routing and traffic management grow on the URL map** — weighted canary splits, mirrors, CORS, CDN cache policies, and fault injection are all `GcpUrlMap` surface; nothing else in the stack changes.
- **Harden Cloud Armor incrementally**: the chart ships the policy open (GCP's default allow) with Layer-7 DDoS defense on. Add allow/deny rules on the policy resource; an explicit rule set must include a priority-`2147483647` default rule.
- **Adding a second domain** means adding it to the certificate's `domains` (certificate replacement, zero downtime) and a host rule on the URL map.
- **Scale-appropriate alternative**: a single service on a single domain with no Cloud Armor/CDN needs may not need this stack at all — `GcpCloudRunDomainMapping` serves that case without a load balancer, and both can coexist during a migration.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
