# Target HTTPS Proxy on Google Cloud

Deploys a global Compute Engine target HTTPS proxy — the TLS-termination node of a global external Application Load Balancer. It binds a global forwarding rule (the VIP) to a URL map (the routing brain) and owns the client-facing handshake: which certificates are presented, which TLS policy constrains ciphers and versions, whether QUIC (HTTP/3) is negotiated, and whether TLS 1.3 0-RTT early data is accepted. Certificates attach through exactly one of three mechanisms — the classic compute-certificate list, Certificate Manager certificates (cross-region internal ALB), or an SNI-scale certificate map. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to projects, URL maps, certificates, and SSL policies.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Compute Engine Target HTTPS Proxy (global)** -- bound to the configured URL map, certificate mechanism, SSL policy, and QUIC/early-data posture

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the proxy will be created.
- **Compute Engine API** (`compute.googleapis.com`) enabled in the target project.
- **A URL map** (GcpUrlMap) and **at least one certificate** (GcpManagedSslCertificate or GcpSslCertificate) covering every host the map routes.

## Deploy

### Console

Open the deployment store, find **Target HTTPS Proxy on Google Cloud**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Managed Cert Frontend** preset in the [Presets](#presets) tab.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpTargetHttpsProxy
metadata:
  name: web-https-proxy
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  urlMap:
    value: "https://www.googleapis.com/compute/v1/projects/acme-prod-12345/global/urlMaps/prod-web-map"
  sslCertificates:
    - value: "https://www.googleapis.com/compute/v1/projects/acme-prod-12345/global/sslCertificates/acme-cert"
```

```shell
planton apply -f target-https-proxy.yaml
```

This creates the standard serving frontend: TLS terminated with the attached certificate, requests routed by the URL map.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the chain:

```yaml
spec:
  urlMap:
    valueFrom:
      kind: GcpUrlMap
      name: prod-web-map
      fieldPath: status.outputs.self_link
  sslCertificates:
    - valueFrom:
        kind: GcpManagedSslCertificate
        name: acme-cert
        fieldPath: status.outputs.self_link
```

The InfraPipeline resolves the dependency graph — certificate and URL map first, then this proxy — and a downstream GcpGlobalForwardingRule references this proxy's `self_link` as its target (it is the forwarding rule's DEFAULT target kind).

## Key Configuration

These are the most important decisions when configuring a target HTTPS proxy. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Certificate mechanism** -- Exactly one of three (GCP rejects combinations): `sslCertificates` (up to 15 — Google-managed and self-managed compute certificates share one collection), `certificateManagerCertificates` (only the cross-region INTERNAL_MANAGED ALB honors them), or `certificateMap` (SNI-scale selection for many domains — SaaS custom domains). All mutable: rotation is attach-before-detach, in place.

**URL map** -- Required; mutable in place (setUrlMap) — repointing a live frontend at a new routing table causes zero downtime.

**TLS posture** -- `sslPolicy` (reference a GcpSslPolicy; empty = GCP's permissive default, min TLS 1.0), `serverTlsPolicy` (the mTLS lever — can demand client certificates; the ONLY TLS mechanism on Traffic Director), `quicOverride` (HTTP/3; unset lets Google decide), and `tlsEarlyData` (0-RTT; replayable by design — STRICT is the furthest most APIs should go; IMMUTABLE).

**Keep-alive / proxy bind** -- Same immutable dials as the HTTP sibling: 5-1200s keep-alive (EXTERNAL_MANAGED only) and the Traffic Director mesh bind.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |
| **GcpUrlMap** | `urlMap` | `status.outputs.self_link` |
| **GcpManagedSslCertificate** | `sslCertificates[]` | `status.outputs.self_link` |
| **GcpSslCertificate** | `sslCertificates[]` (explicit kind) | `status.outputs.self_link` |
| **GcpCertManagerCert** | `certificateManagerCertificates[]` | `status.outputs.certificate_name` |
| **GcpSslPolicy** | `sslPolicy` | `status.outputs.self_link` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `self_link` | Self-link URI of the proxy | GcpGlobalForwardingRule `target` (the default target kind) |
| `proxy_name` | Name as it exists in GCP | Audit, fleet inventory |
| `proxy_id` | Server-assigned numeric ID | Diagnostics |
| `fingerprint` | Optimistic-concurrency token | Out-of-band gcloud updates |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Managed cert frontend** -- The standard shape: a Google-managed certificate + the application URL map. Start from the **Managed Cert Frontend** preset.

**Certificate map SaaS** -- SNI-scale certificate selection for many customer domains. Start from the **Certificate Map SaaS** preset.

**mTLS server TLS policy** -- Demanding and validating client certificates. Start from the **MTLS Server TLS Policy** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the proxy is created
- [**GCP URL Map**](/cloud-catalog/gcp-url-map) -- the routing table for decrypted requests
- [**GCP Managed SSL Certificate**](/cloud-catalog/gcp-managed-ssl-certificate) -- the auto-renewing certificates this proxy presents
- [**GCP SSL Policy**](/cloud-catalog/gcp-ssl-policy) -- the TLS versions/ciphers floor
- [**GCP Global Forwarding Rule**](/cloud-catalog/gcp-global-forwarding-rule) -- consumes this proxy's `self_link` as its target
