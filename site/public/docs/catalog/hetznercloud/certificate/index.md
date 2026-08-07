---
title: "Certificate"
description: "Certificate deployment documentation"
icon: "package"
order: 100
componentName: "hetznercloudcertificate"
---

# Hetzner Cloud Certificate

Provisions a TLS certificate on Hetzner Cloud for use with load balancer HTTPS services. Supports two mutually exclusive modes: uploaded certificates where you provide your own PEM-encoded certificate chain and private key, or managed certificates where Hetzner Cloud automatically obtains and renews a Let's Encrypt certificate for specified domain names.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions one of:

- **Uploaded Certificate** -- an `hcloud_uploaded_certificate` resource storing your PEM-encoded certificate chain and private key
- **Managed Certificate** -- an `hcloud_managed_certificate` resource that automatically obtains and renews a Let's Encrypt certificate for the specified domains

## Before You Deploy

### Hetzner Cloud Account

- **A Hetzner Cloud account** with an active project and API token.

### For Uploaded Certificates

- **A PEM-encoded TLS certificate chain** (server cert + intermediate CAs).
- **The corresponding PEM-encoded private key** (sensitive material).

### For Managed Certificates

- **Domain names** with DNS records pointing to a Hetzner Cloud load balancer so the ACME HTTP-01 challenge can succeed.
- **A load balancer** with an HTTPS service configured to reference this certificate.

## Deploy

### Console

Open the deployment store, find **Hetzner Cloud Certificate**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and the choice between uploaded and managed certificate types.

### CLI

**Managed certificate (Let's Encrypt):**

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudCertificate
metadata:
  name: web-cert
  org: acme-corp
  env: prod
spec:
  managed:
    domainNames:
      - "example.com"
      - "www.example.com"
```

**Uploaded certificate:**

```yaml
apiVersion: hetzner-cloud.planton.dev/v1
kind: HetznerCloudCertificate
metadata:
  name: custom-cert
  org: acme-corp
  env: prod
spec:
  uploaded:
    certificate: |
      -----BEGIN CERTIFICATE-----
      ...
      -----END CERTIFICATE-----
    privateKey: |
      -----BEGIN PRIVATE KEY-----
      ...
      -----END PRIVATE KEY-----
```

```shell
planton apply -f hetznercloud-certificate.yaml
```

A Stack Job tracks the provisioning in real time. Reference the certificate in HetznerCloudLoadBalancer HTTPS services via `certificateIds`.

## Key Configuration

These are the most important decisions when configuring a certificate. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Certificate type** -- Exactly one of `uploaded` or `managed` must be set. Uploaded certificates require you to provide and renew the certificate yourself. Managed certificates use Let's Encrypt for automatic issuance and renewal but require DNS records pointing to a Hetzner Cloud load balancer.

**Domain names (managed)** -- The `managed.domainNames` field lists the domains covered by the certificate. Hetzner Cloud issues a single SAN certificate covering all listed domains. Changing this list forces replacement.

**Certificate chain (uploaded)** -- The `uploaded.certificate` field accepts the PEM-encoded certificate chain (server cert first, root last). Changing forces replacement.

**Private key (uploaded)** -- The `uploaded.privateKey` field accepts the PEM-encoded private key. This is sensitive material handled as a secret by the IaC modules. Changing forces replacement.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `certificate_id` | Hetzner Cloud numeric ID of the certificate | HetznerCloudLoadBalancer HTTPS service `certificateIds` |
| `type` | Certificate type ("uploaded" or "managed") | Display and monitoring |
| `fingerprint` | SHA256 fingerprint of the certificate | Certificate verification |
| `not_valid_before` | Certificate validity start (ISO-8601) | Expiration monitoring |
| `not_valid_after` | Certificate validity end (ISO-8601) | Expiration monitoring |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Let's Encrypt auto-renewal** -- Create a managed certificate with domain names and reference it from a load balancer HTTPS service. Hetzner Cloud handles issuance and renewal automatically.

**Custom CA certificate** -- Upload a certificate from a corporate CA or a paid certificate authority for domains requiring extended validation or specific trust chains.

## Works With

- [**Hetzner Cloud Load Balancer**](/cloud-catalog/hetznercloud-load-balancer) -- HTTPS services reference this certificate for TLS termination
