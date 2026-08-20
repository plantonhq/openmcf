---
title: "Service Accounts for Automated Workflows"
date: 2026-03-02
category: feature
tags:
  - platform
  - console
  - cli
excerpt: "Create service accounts with API keys to authenticate CI/CD pipelines, runners, and custom automation against the Planton platform."
author:
  - name: Swarup Donepudi
    title: Founder
---

You can now create service accounts to give your automation a proper identity on the platform. Service accounts carry their own API keys, so CI/CD pipelines, runners, and scripts can authenticate without borrowing a human user's credentials. Every action taken by a service account is tracked under its own identity, giving you a clear audit trail for machine-driven operations.

## Managing Service Accounts from the Console

A new **Service Accounts** tab is available under **Org Settings**. If you haven't created any yet, you'll see an educational overview explaining what service accounts are, how they relate to runners, and when you'd want to create one manually.

To create a service account, click the **Create** button and provide a name and optional description. The platform provisions the backing identity and access controls automatically.

From the detail page for any service account, you can create API keys, view existing keys by fingerprint, revoke specific keys, or delete the entire service account. Deleting a service account cascades — all its API keys are revoked and its access is cleaned up across the platform.

When you create a new API key, the raw key is shown exactly once in a dedicated dialog with a copy button. Once you close the dialog, the key cannot be retrieved again. This matches the security pattern used by GCP, AWS, and other cloud providers.

## Managing Service Accounts from the CLI

The new `planton service-account` command group (alias: `sa`) gives you the same lifecycle control from the terminal:

```bash
# Create a service account
planton service-account create --name "deploy-bot" --org my-org

# Create an API key for it
planton service-account key create sa-xxxx

# List all service accounts in your org
planton service-account list --org my-org

# Revoke a specific key
planton service-account key revoke sa-xxxx --key-id ak-yyyy
```

To authenticate the CLI as a service account, use the API key directly:

```bash
planton auth login --api-key pck_...
```

This works with any API key — whether it belongs to a human user or a service account. Once logged in, `planton auth who` shows your identity type and, for service accounts, the owning account name and organization. The `planton auth list` command also distinguishes between OAuth sessions, user API keys, and service account API keys across your configured environments.

## Runners Get Service Accounts Automatically

When you generate credentials for a runner, the platform now provisions a service account behind the scenes and includes its API key in the credentials file. The runner uses this identity to resolve secrets and configuration values at runtime — for example, when an IaC module references a connection secret, the runner authenticates to the control plane with its service account API key to fetch the actual value.

You don't need to create or manage runner service accounts manually. They're provisioned on first credential generation and persist across credential regeneration cycles. If you regenerate credentials, the old API key is revoked and a new one is issued automatically.

## Why This Matters

- **No more shared credentials** — give each pipeline, tool, or runner its own identity with individually revocable API keys
- **Auditable automation** — every API call made by a service account is tracked under its own machine identity, separate from human users
- **Familiar model** — if you've used GCP service accounts or AWS IAM roles for workload identity, this works the same way
- **Non-interactive auth** — `planton auth login --api-key` provides a proper headless login path for CI/CD environments where browser-based OAuth isn't possible
