---
title: "Certificate"
description: "Certificate deployment documentation"
icon: "package"
order: 100
componentName: "digitaloceancertificate"
---

# Certificate on DigitalOcean

Deploys an SSL/TLS certificate on DigitalOcean, supporting both free auto-renewing Let's Encrypt certificates and user-provided custom certificates. The certificate name can be referenced by DigitalOcean Load Balancers for HTTPS termination. Integrates with Planton's Provider Connections for DigitalOcean credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DigitalOcean Certificate** -- a `digitalocean_certificate` resource with the specified name and type (Let's Encrypt or custom)
- **Let's Encrypt Certificate** -- created only when `type` is `lets_encrypt`; includes the specified domains with automatic ACME validation and renewal managed by DigitalOcean
- **Custom Certificate** -- created only when `type` is `custom`; uploads the provided PEM-encoded leaf certificate, private key, and optional intermediate chain with `create_before_destroy` lifecycle for zero-downtime rotation
- **DigitalOcean Tags** -- applied when `tags` are specified for organizational grouping

## Before You Deploy

### Planton Setup

- **DigitalOcean Provider Connection** -- an active connection in the Connect module with a DigitalOcean API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### DigitalOcean Account

- **For Let's Encrypt certificates** -- the domain(s) must be managed by DigitalOcean DNS so that ACME validation can succeed. Wildcard certificates (e.g., `*.example.com`) require DNS validation.
- **For custom certificates** -- have the PEM-encoded leaf certificate, private key, and optional intermediate chain ready.

## Deploy

### Console

Open the deployment store, find **Certificate on DigitalOcean**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Let's Encrypt Certificate** preset in the [Presets](#presets) tab for a free, auto-renewing certificate.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: digital-ocean.planton.dev/v1
kind: DigitalOceanCertificate
metadata:
  name: my-cert
  org: acme-corp
  env: prod
spec:
  certificateName: my-cert
  type: lets_encrypt
  letsEncrypt:
    domains:
      - example.com
      - www.example.com
```

```shell
planton apply -f certificate.yaml
```

This creates a free Let's Encrypt certificate covering `example.com` and `www.example.com` with automatic renewal enabled. Reference the certificate name in load balancer configurations for HTTPS termination.

## Key Configuration

These are the most important decisions when configuring a DigitalOcean certificate. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Certificate type** -- Set `type` to `lets_encrypt` for a free, auto-renewing certificate managed by DigitalOcean, or `custom` when providing your own PEM-encoded certificate from an enterprise CA or purchased provider. Only one source branch (`letsEncrypt` or `custom`) can be specified.

**Let's Encrypt domains** -- The `letsEncrypt.domains` field accepts multiple FQDNs and wildcard domains (e.g., `*.example.com`). All domains must have DNS managed by DigitalOcean for ACME validation. Set `disableAutoRenew` to `true` only if you need manual control over renewal timing.

**Custom certificate materials** -- When using `custom` type, provide `leafCertificate` (server cert), `privateKey` (matching key), and optionally `certificateChain` (intermediate certs). Custom certificates have no automatic renewal -- plan for manual rotation before expiry.

**Certificate naming** -- Use a stable, descriptive `certificateName` since load balancers reference certificates by name. This ensures IaC state survives Let's Encrypt renewals without breaking load balancer configuration.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `certificate_id` | Unique identifier (UUID) of the certificate in DigitalOcean | Load balancer HTTPS listener configuration, API operations |
| `expiry_rfc3339` | Expiration timestamp in RFC 3339 format | Certificate rotation monitoring, alerting dashboards |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Let's Encrypt certificate** -- Free, auto-renewing certificate for one or more domains managed by DigitalOcean DNS. Covers the common case of HTTPS for public websites and APIs. Start from the **Let's Encrypt Certificate** preset.

**Custom certificate** -- User-provided PEM certificate for enterprise CAs, purchased certificates, or EV certificates not available through Let's Encrypt. Requires manual rotation before expiry. Start from the **Custom Certificate** preset.

## Works With

This component operates independently and does not reference other deployment components.