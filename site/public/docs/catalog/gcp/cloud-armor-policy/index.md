---
title: "Cloud Armor Policy"
description: "Cloud Armor Policy deployment documentation"
icon: "package"
order: 100
componentName: "gcpcloudarmorpolicy"
---

# GCP Cloud Armor Policy

Deploys a Cloud Armor security policy with configurable rules for IP allowlisting/denylisting, rate limiting, ban escalation, OWASP WAF protection, and Layer 7 DDoS defense. The policy attaches to HTTP(S) load balancers, Cloud CDN backends, or internal Traffic Director services. Integrates with Planton's Provider Connections for GCP credential management and supports ValueFromRef wiring to GCP projects.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Security Policy** -- a `compute.SecurityPolicy` in the specified GCP project, configured with the chosen policy type, rules, and advanced options
- **Security Rules** -- one rule per entry in `rules`, each with a priority, action (allow, deny, throttle, rate_based_ban, redirect), match condition (IP ranges or CEL expression), and optional rate limiting, redirect, header injection, or WAF exclusion configuration
- **Adaptive Protection** -- created only when `adaptiveProtectionConfig` is present; enables automatic Layer 7 DDoS detection and alerting
- **Advanced Options** -- created only when `advancedOptionsConfig` is present; configures JSON body parsing (with optional custom content types), logging verbosity, and client IP resolution headers
- **Default Rule** -- if no rule at priority 2147483647 is provided, the IaC module auto-adds a default "allow all" rule
- **GCP Labels** -- resource metadata labels (resource name, kind, organization, environment) applied automatically for tracking and governance

## Before You Deploy

### Planton Setup

- **GCP Provider Connection** -- an active connection in the Connect module with credentials for the target GCP project. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or browser OAuth authentication modes.

### GCP Project

- **A GCP project** where the security policy will be created. Provide the project ID directly or reference a GcpProject Cloud Resource via ValueFromRef.
- **Compute Engine API** (`compute.googleapis.com`) enabled in the target project.
- **An HTTP(S) load balancer** or backend service to attach the policy to (configured outside this Cloud Resource via the backend service's `securityPolicy` field).

## Deploy

### Console

Open the deployment store, find **GCP Cloud Armor Policy**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Basic IP Allowlist** preset in the [Presets](#presets) tab to pre-populate a deny-by-default policy with IP-based allow rules.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: gcp.planton.dev/v1alpha1
kind: GcpCloudArmorPolicy
metadata:
  name: api-protection
  org: acme-corp
  env: prod
spec:
  projectId:
    value: "acme-prod-12345"
  policyName: api-protection
  type: CLOUD_ARMOR
  rules:
    - action: allow
      priority: 1000
      match:
        versionedExpr: SRC_IPS_V1
        srcIpRanges: ["10.0.0.0/8"]
    - action: "deny(403)"
      priority: 2147483647
      match:
        versionedExpr: SRC_IPS_V1
        srcIpRanges: ["*"]
```

```shell
planton apply -f gcp-cloud-armor-policy.yaml
```

This creates a CLOUD_ARMOR policy that allows traffic from RFC 1918 ranges and denies everything else with 403. The policy must be attached to a backend service separately. A Stack Job tracks the provisioning in real time.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the policy to a GCP project deployed in the same InfraPipeline:

```yaml
spec:
  projectId:
    valueFrom:
      kind: GcpProject
      name: production-project
      fieldPath: status.outputs.project_id
```

The InfraPipeline resolves the dependency graph, deploys the project first, then provisions the Cloud Armor policy with the resolved project ID.

## Key Configuration

These are the most important decisions when configuring a Cloud Armor policy. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Policy type** -- Set `type` to CLOUD_ARMOR (default) for backend security policies with full WAF, rate limiting, and header injection. Use CLOUD_ARMOR_EDGE for CDN and backend bucket protection (IP/geo rules only). Use CLOUD_ARMOR_INTERNAL_SERVICE for internal Traffic Director services. Type is immutable after creation.

**Rule priority and evaluation** -- Rules are evaluated from lowest priority number (highest precedence) to highest. Priority 2147483647 is reserved for the default rule. Plan priority numbering with gaps (e.g., 1000, 2000, 3000) to allow inserting rules later without renumbering.

**Rate limiting and ban escalation** -- Use `throttle` action for simple rate limiting or `rate_based_ban` for two-tier protection. Rate limiting uses `rateLimitThreshold` to cap requests per interval; ban escalation adds `banThreshold` and `banDurationSec` to fully block persistent abusers. Set `enforceOnKey` to `IP` for per-source limiting.

**WAF rules and exclusions** -- Match against preconfigured OWASP rule sets (e.g., `sqli-v33-stable`, `xss-v33-stable`) using CEL expressions. Add `preconfiguredWafConfig` exclusions for request fields that trigger false positives -- critical for production APIs that accept user-generated content.

**Adaptive Protection** -- Enable `adaptiveProtectionConfig.enableLayer7DdosDefense` for automatic anomaly detection. STANDARD visibility is available to all Cloud Armor users. PREMIUM requires Managed Protection Plus.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **GcpProject** | `projectId` | `status.outputs.project_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `policy_id` | Fully qualified resource ID (`projects/{p}/global/securityPolicies/{name}`) | Backend service security policy references |
| `policy_name` | Name of the security policy in GCP | Audit logs, monitoring dashboards |
| `policy_self_link` | Self-link URI of the security policy | Attaching to backend services, load balancers, CDN configurations |
| `fingerprint` | Server-computed fingerprint for optimistic concurrency | Out-of-band policy updates |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Basic IP allowlist** -- Deny-by-default policy that allows traffic only from specified CIDR ranges (corporate networks, VPNs). All other traffic receives 403 Forbidden. Suitable for internal dashboards and admin APIs. Start from the **Basic IP Allowlist** preset.

**Rate limiting for APIs** -- Per-IP rate limiting with ban escalation. Throttles traffic beyond 100 requests per minute, then bans persistent abusers exceeding 500 requests over 5 minutes for 1 hour. Start from the **Rate Limiting API** preset.

**WAF OWASP protection** -- OWASP WAF rules blocking SQL injection and XSS attacks, with adaptive Layer 7 DDoS protection, JSON body parsing, and verbose logging. Suitable for internet-facing web applications. Start from the **WAF OWASP Protection** preset.

## Works With

- [**GCP Project**](/cloud-catalog/gcp-project) -- provides the GCP project where the security policy is created