---
title: "Turnstile Widget"
description: "Turnstile Widget deployment documentation"
icon: "package"
order: 100
componentName: "cloudflareturnstilewidget"
---

# Turnstile Widget on Cloudflare

Provisions a Cloudflare Turnstile widget: a privacy-preserving CAPTCHA alternative that protects forms and endpoints from bots without the friction of traditional CAPTCHAs. A widget yields a public **site key** you embed in your page and a sensitive **secret key** your backend uses to verify tokens via the `/siteverify` endpoint. Integrates with Planton's Provider Connections for Cloudflare credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Turnstile Widget** -- a configured widget scoped to your account and domains
- **Site key and secret key** -- exported as stack outputs (the secret is sensitive)

## Before You Deploy

### Planton Setup

- **Cloudflare Provider Connection** -- an active connection in the Connect module with a Cloudflare API token that has Turnstile edit access. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline API token authentication.

### Cloudflare Account

- **An account** -- the widget is created under an existing Cloudflare account (selected from the connection's accounts).

## Deploy

### Console

Open the deployment store, find **Turnstile Widget on Cloudflare**, and click **Deploy**. The creation wizard captures the account, name, served domains, challenge mode, and region.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: cloudflare.planton.dev/v1
kind: CloudflareTurnstileWidget
metadata:
  name: signup-widget
  org: acme-corp
  env: prod
spec:
  accountId: a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4
  name: signup
  domains:
    - example.com
    - www.example.com
  mode: managed
```

```shell
planton apply -f cloudflare-turnstile-widget.yaml
```

This creates a managed-mode widget for the listed domains. A Stack Job tracks the provisioning in real time; the site and secret keys appear in `status.outputs` once issued.

## Key Configuration

These are the most important decisions when configuring a Turnstile widget. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Account (`accountId`)** -- The Cloudflare account that owns the widget. Selected from the connection's accounts.

**Domains (`domains`)** -- The hostnames allowed to serve the widget. Tokens are only issued on listed domains; use `localhost` for local development.

**Mode (`mode`)** -- `managed` (Cloudflare picks the challenge -- recommended), `non-interactive`, or `invisible`.

**Clearance Level (`clearanceLevel`)** -- Optional; the clearance granted on a Cloudflare-proxied site (`no_clearance`, `jschallenge`, `managed`, `interactive`).

**Region (`region`)** -- `world` (default) or `china`. **Immutable** -- changing it rotates the site and secret keys.

## Outputs and Dependencies

### What This Component Consumes

The widget is account-scoped; it references no other Cloud Resources.

### What This Component Provides

After provisioning, `status.outputs` contains:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `sitekey` | The public site key | Embed in the frontend Turnstile widget |
| `secret` | The sensitive secret key | Server-side token verification via `/siteverify` |
| `created_on` | Creation timestamp | Auditing |
| `modified_on` | Last-modified timestamp | Auditing |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Managed challenge** -- the recommended default; Cloudflare adapts the challenge to risk.

**Invisible widget** -- zero visible UI for low-friction flows, paired with server-side verification.

## Works With

- A **Cloudflare Worker** or backend service can reference the widget's `secret` output to verify Turnstile tokens.
