# Zero Trust App on Cloudflare

Deploys a Cloudflare Zero Trust Access Application that protects a hostname behind identity-aware access controls, with configurable email allowlists, Google Workspace group restrictions, session duration, and optional MFA enforcement. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Access Application** -- a self-hosted Access Application bound to the specified DNS zone and hostname, with configurable session duration
- **Access Policy** -- an allow or deny policy attached to the application, with email-based include rules, optional Google Workspace group includes, and an optional MFA requirement
- **Cloudflare Labels** -- resource metadata applied for organization and environment tracking

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Access edit permissions. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **An existing Cloudflare DNS zone** -- the `zoneId` is required and must correspond to a zone containing the hostname you intend to protect. Obtain it from a CloudflareDnsZone resource's `status.outputs.zone_id` or from the Cloudflare dashboard.
- **Zero Trust subscription** -- Access Applications require a Zero Trust plan enabled on the Cloudflare account.
- **DNS resolution** -- the hostname specified in `hostname` must resolve within the specified zone. Configure a DNS record (A, AAAA, or CNAME) pointing to your origin before or alongside this deployment.

## Deploy

### Console

Open the deployment store, find **Zero Trust App on Cloudflare**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Company-Wide Email** preset in the [Presets](#presets) tab to pre-populate a working configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareZeroTrustAccessApplication
metadata:
  name: internal-dashboard
  org: acme-corp
  env: prod
spec:
  applicationName: Internal Dashboard
  zoneId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  hostname: dashboard.example.com
  allowedEmails:
    - "*@acme-corp.com"
```

```shell
planton apply -f cloudflare-zero-trust-access-application.yaml
```

This creates a Zero Trust Access Application protecting `dashboard.example.com` with an allow policy granting access to all `@acme-corp.com` email addresses. MFA is not required and the session lasts 24 hours by default. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring a Zero Trust Access Application. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Policy type** -- Set `policyType` to `ALLOW` (default) to permit matching identities or `BLOCK` to deny them. Most configurations use ALLOW with a set of permitted emails or groups.

**Access control method** -- Use `allowedEmails` for email-based access (individual addresses or wildcard patterns like `*@company.com`). Use `allowedGoogleGroups` for Google Workspace group-based access with group IDs from Google Admin Console. Both can be combined in a single policy.

**Session duration** -- The `sessionDurationMinutes` field controls how long each authenticated session lasts before re-authentication is required. Defaults to 1440 minutes (24 hours). Shorten to 480 (8 hours) or less for sensitive applications.

**MFA enforcement** -- Set `requireMfa: true` to add a multi-factor authentication requirement to the access policy. Users must complete a second authentication factor before gaining access. Recommended for admin panels and sensitive internal tools.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **CloudflareDnsZone** | `zoneId` | `status.outputs.zone_id` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `application_id` | The unique Cloudflare ID of the Access Application | Access management automation, monitoring dashboards |
| `public_hostname` | The hostname protected by this Access Application | DNS verification, application endpoint documentation |
| `policy_id` | The Cloudflare ID of the Access Policy attached to this application | Policy management, audit logging references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Company-wide email access** -- Allows anyone with a company domain email to access the protected hostname. Use for internal dashboards, wikis, or tools where the entire organization should have access. Start from the **Company-Wide Email** preset.

**Team access with Google Groups and MFA** -- Restricts access to specific Google Workspace groups and requires multi-factor authentication. Use for admin panels or sensitive tools where only certain teams should have access and MFA is mandated for compliance. Start from the **Team Google Groups** preset.

## Works With

- [**DNS Zone on Cloudflare**](/cloud-catalog/cloudflare-dns-zone) -- provides the zone ID for the domain containing the protected hostname