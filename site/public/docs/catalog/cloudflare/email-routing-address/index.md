---
title: "Email Routing Address"
description: "Email Routing Address deployment documentation"
icon: "package"
order: 100
componentName: "cloudflareemailroutingaddress"
---

# Email Routing Address on Cloudflare

Registers a verified destination address for Cloudflare Email Routing -- an account-scoped mailbox that routing rules and zone catch-alls forward to. Creating it sends a verification email to the mailbox; the address cannot receive forwarded mail until its owner clicks the link. Because addresses are account-scoped, you register a teammate's inbox once and reference it from any routing rule or zone catch-all in the account. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Destination Address** -- an account-scoped verified mailbox usable as a forwarding target
- **Verification Email** -- Cloudflare emails the address a confirmation link on creation; forwarding stays inert until it is verified

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Email Routing edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **Account-level access** -- destination addresses are created at the account level, so the API token must be scoped to the account.
- **Mailbox access** -- the owner of the destination mailbox must be able to click the verification link Cloudflare sends.

## Deploy

### Console

Open the deployment store, find **Email Routing Address on Cloudflare**, and click **Deploy**. The creation wizard captures the owning account and the destination email. Both are fixed at creation -- changing either replaces the address.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareEmailRoutingAddress
metadata:
  name: ops-mailbox
  org: acme-corp
  env: prod
spec:
  accountId: "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4"
  email: ops@example.com
```

```shell
planton apply -f cloudflare-email-routing-address.yaml
```

This registers `ops@example.com` as a destination. A Stack Job tracks the provisioning in real time, and Cloudflare emails the mailbox a verification link.

## Key Configuration

These are the most important decisions when configuring a destination address. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Destination Email (`email`)** -- The real mailbox forwarded mail is delivered to. Immutable -- changing it replaces the address, and the new mailbox must be re-verified.

**Account (`accountId`)** -- The owning Cloudflare account. Immutable, and must match the account of the rules and zones that forward here.

## Outputs and Dependencies

### What This Component Consumes

Nothing -- a destination address is a leaf resource.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `address_id` | The Cloudflare-assigned identifier of the address | Verification, dashboards |
| `email` | The destination email (echoed) | Referenced by a `CloudflareEmailRoutingRule` action or a `CloudflareEmailRoutingZone` catch-all (`forwardTo`) |
| `verified` | RFC3339 timestamp once verified, empty until then | Gate forwarding readiness |
| `created` | Provisioning timestamp | Auditing |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Shared team inbox** -- register `support@acme.com` once, then reference it from several routing rules across zones.

**Personal forwarding target** -- register a personal mailbox to forward a custom-domain alias to it.

## Works With

- [**Email Routing Rule on Cloudflare**](/cloud-catalog/cloudflare-email-routing-rule) -- forwards matched mail to this address
- [**Email Routing Zone on Cloudflare**](/cloud-catalog/cloudflare-email-routing-zone) -- a forwarding catch-all delivers unmatched mail to this address
