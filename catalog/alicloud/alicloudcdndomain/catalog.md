# AliCloud CDN Domain

Deploys an Alibaba Cloud CDN accelerated domain that maps a user-facing domain name to one or more origin servers. Alibaba Cloud's globally distributed edge nodes cache and serve content from the origins, reducing latency for end users. The component supports multiple origin source types, priority-based failover, HTTPS certificate configuration, and geographic scope control. After provisioning, point a CNAME record at your DNS provider to the `cname` output value for CDN acceleration to take effect. The component integrates with Planton's Provider Connections for AliCloud credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **CDN Domain** -- an `alicloud_cdn_domain_new` with the specified domain name, CDN type, geographic scope, and origin source configuration
- **Origin Sources** -- one source block per entry in the `sources` list, each with configurable type (IP, domain, OSS, common), port, priority, and weight for load distribution and failover
- **HTTPS Certificate** -- created only when `certificateConfig` is provided; configures TLS on edge nodes using an uploaded certificate, CAS-managed certificate, or free DV certificate
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- **A domain name** -- the accelerated domain (e.g., `cdn.example.com`) must be a valid domain you own. The domain name is immutable after creation.
- **At least one origin server** -- an IP address, domain name, or OSS bucket that serves the content to be cached by CDN edge nodes.
- **DNS access** -- after provisioning, you must create a CNAME record at your DNS provider pointing the accelerated domain to the `cname` output value. CDN acceleration does not take effect until this DNS change is made.
- **An SSL certificate** (optional) -- required for HTTPS. Provide a certificate from Alibaba Cloud Certificate Management Service (`cert_type: cas`), upload your own PEM-encoded certificate and key (`cert_type: upload`), or use a free DV certificate (`cert_type: free`).

## Deploy

### Console

Open the deployment store, find **AliCloud CDN Domain**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Web HTTPS** preset in the [Presets](#presets) tab to pre-populate a configuration for HTTPS-enabled web content acceleration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudCdnDomain
metadata:
  name: cdn-example
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  domainName: cdn.example.com
  cdnType: web
  scope: global
  sources:
    - type: domain
      content: origin.example.com
      port: 443
      priority: 20
```

```shell
planton apply -f cdn-domain.yaml
```

This creates a CDN domain for web content acceleration with global scope and a single domain-type origin server on port 443. HTTPS is not configured on the CDN edge. After provisioning, create a CNAME record pointing `cdn.example.com` to the `cname` output value.

## Key Configuration

These are the most important decisions when configuring a CDN domain. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**CDN type** -- Set `cdnType` to `web` (images, small files, web pages), `download` (large file downloads), or `video` (on-demand and live streaming). This determines the caching and delivery optimization strategy. The CDN type is immutable after creation.

**Geographic scope** -- Set `scope` to `domestic` (mainland China only), `overseas` (outside mainland China), or `global` (worldwide). If omitted, defaults to `domestic`. Choose `global` for applications serving users across regions.

**Origin sources** -- Configure `sources` with one or more origin servers. Each source has a `type` (ipaddr, domain, oss, common), `content` (the address), `port` (typically 80 or 443), `priority` (lower = higher priority; use 20 for primary, 30 for standby), and `weight` (for load distribution among same-priority sources). Multiple sources enable failover and load distribution.

**HTTPS** -- Configure `certificateConfig` to enable HTTPS on CDN edge nodes. Set `certType` to `cas` (Certificate Management Service), `upload` (provide PEM certificate and key), or `free` (free DV certificate from Alibaba Cloud). Set `serverCertificateStatus: on` to enable HTTPS.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `domain_name` | The accelerated domain name as registered in CDN | DNS configuration, application asset URLs |
| `cname` | CNAME assigned by Alibaba Cloud CDN | DNS CNAME record target (required for CDN to take effect) |
| `status` | Current CDN domain status (e.g., `online`, `configuring`) | Health monitoring, deployment verification |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Web content with HTTPS** -- A CDN domain for web pages, images, and small files with HTTPS enabled via a CAS-managed or free certificate, global scope, and a single domain-type origin. Start from the **Web HTTPS** preset.

**OSS static assets** -- A CDN domain backed by an Alibaba Cloud OSS bucket for serving static assets (JavaScript, CSS, images) with caching at the edge. Uses the OSS bucket's internal domain as the origin for optimal transfer speed. Start from the **OSS Static Assets** preset.

## Works With

This component operates independently and does not reference other deployment components.