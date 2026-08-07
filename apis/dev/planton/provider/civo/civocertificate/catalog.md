# Certificate on Civo

Deploys a TLS certificate on Civo Cloud, supporting free Let's Encrypt certificates with automatic renewal and custom user-provided certificates with PEM-encoded keys and chains. Integrates with Planton's Provider Connections for Civo credential management. Note: the Civo Terraform/Pulumi provider does not yet expose a certificate resource, so the IaC module validates and stores the specification while certificates are managed via the Civo dashboard or API until provider support is added.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Civo Certificate** -- a TLS certificate associated with the Civo account, either auto-managed via Let's Encrypt or uploaded as a custom PEM-encoded certificate
- **Auto-Renewal Configuration** -- created only when using Let's Encrypt type with auto-renewal enabled (default); renews the certificate before the 90-day expiration

## Before You Deploy

### Planton Setup

- **Civo Provider Connection** -- an active connection in the Connect module with a Civo API token. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Civo Account

- **A Civo account** with DNS configured for the target domain(s) when using Let's Encrypt. Domain validation requires DNS records to be resolvable via Civo's nameservers.
- **PEM-encoded certificate files** when using custom type -- the leaf certificate, private key, and optionally the intermediate certificate chain.

## Deploy

### Console

Open the deployment store, find **Certificate on Civo**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Let's Encrypt** preset in the [Presets](#presets) tab for a free, auto-renewing wildcard certificate.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: civo.planton.dev/v1
kind: CivoCertificate
metadata:
  name: app-cert
  org: acme-corp
  env: prod
spec:
  certificateName: app-cert
  type: letsEncrypt
  letsEncrypt:
    domains:
      - "example.com"
      - "*.example.com"
```

```shell
planton apply -f civo-certificate.yaml
```

This requests a Let's Encrypt wildcard certificate covering the apex domain and all subdomains with automatic renewal enabled. No custom PEM files are needed. A Stack Job tracks the operation in real time.

## Key Configuration

These are the most important decisions when configuring a Civo certificate. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Certificate type** -- Set `type` to `letsEncrypt` for free, browser-trusted certificates with automatic renewal, or `custom` to upload your own PEM-encoded certificate and private key. Let's Encrypt certificates expire every 90 days but renew automatically unless `disableAutoRenew` is set to `true`.

**Domain coverage** -- For Let's Encrypt, specify one or more domains in the `domains` list. Use wildcard entries (e.g., `"*.example.com"`) to cover all subdomains with a single certificate. Include both the apex and wildcard to cover both `example.com` and `*.example.com`.

**Auto-renewal** -- Enabled by default for Let's Encrypt certificates. Set `disableAutoRenew` to `true` only if you need manual control over certificate renewal timing -- for example, in staging environments where you test certificate rotation procedures.

**Custom certificate chain** -- When using a custom certificate, provide `leafCertificate` and `privateKey` as PEM-encoded strings. Include `certificateChain` for intermediate CA certificates that clients may need to validate the full trust path.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `certificate_id` | Unique identifier of the certificate in Civo | Load balancer TLS configuration, Civo API references |
| `expiry_rfc3339` | Certificate expiration timestamp in RFC 3339 format | Monitoring dashboards, renewal alerting |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Let's Encrypt wildcard** -- a free wildcard certificate covering the apex domain and all subdomains with automatic renewal. Covers the majority of web hosting and API scenarios without managing certificate files manually. Start from the **Let's Encrypt** preset.

## Works With

This component operates independently and does not reference other deployment components.