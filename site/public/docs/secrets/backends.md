---
title: "Secret Backends"
sidebar_title: "Backends"
description: "Choose where secrets are stored — from zero-config built-in storage to bring-your-own backends in your own cloud accounts"
icon: database
order: 30
tags:
  - Secrets
  - Backends
  - Storage
---

# Secret Backends

Every secret in Planton is stored in a **secret backend** — a dedicated storage system designed for sensitive data. The backend determines where your encrypted secret versions are persisted.

Out of the box, every organization gets a built-in backend powered by [OpenBAO](https://openbao.org/) (an open-source secrets management engine) with zero configuration required. Organizations that need to keep secrets in their own infrastructure can bring their own backend — AWS Secrets Manager, HashiCorp Vault, GCP Secret Manager, Azure Key Vault, or a self-hosted OpenBAO instance.

Regardless of which backend you choose, Planton encrypts every secret before it touches the backend. This is not optional. You choose **where** secrets are stored — not whether they are encrypted. See [Encryption](/docs/secrets/encryption) for how encryption works.

## The Built-In Backend

The platform backend is the default for every organization. When you create a secret without specifying a backend, it is stored here.

**What you get with zero configuration:**

- Secrets encrypted with AES-256-GCM envelope encryption before storage
- Per-organization encryption key isolation — each organization has its own key, managed by the platform
- Automatic key rotation support
- No infrastructure to deploy or maintain
- Available immediately when you create your organization

The built-in backend is a production-grade OpenBAO instance operated by Planton in high-availability mode. Your data is encrypted by Planton's encryption layer before it reaches OpenBAO, and then encrypted again by OpenBAO's own storage barrier. Two layers of encryption, zero configuration.

For most teams, the built-in backend is the right choice. It provides production-ready security without the operational overhead of managing your own secrets infrastructure.

## Bring Your Own Backend

Organizations that require secrets to remain within their own cloud accounts — for compliance, data residency, or policy reasons — can connect their own backend.

| Backend | When to Use |
|---------|-------------|
| **AWS Secrets Manager** | Your organization standardizes on AWS for secrets management, or compliance requires secrets to stay in your AWS account |
| **GCP Secret Manager** | Your organization standardizes on GCP, or you need secrets co-located with GCP workloads |
| **Azure Key Vault** | Your organization standardizes on Azure, or Azure policy requires secrets to stay in your tenant |
| **HashiCorp Vault** | You already operate a Vault cluster and want to consolidate secrets management |
| **Self-hosted OpenBAO** | You want the same technology as the built-in backend, but operated in your own infrastructure |

### Setting Up a Custom Backend

Each backend requires connection credentials so Planton can read and write secrets on your behalf. The setup varies by provider:

**AWS Secrets Manager** — Provide the AWS region, an access key ID, and a secret access key. The IAM user or role needs `secretsmanager:CreateSecret`, `secretsmanager:GetSecretValue`, `secretsmanager:PutSecretValue`, `secretsmanager:DeleteSecret`, and `secretsmanager:DescribeSecret` permissions.

**GCP Secret Manager** — Provide the GCP project ID and a service account key (JSON). The service account needs the `Secret Manager Admin` role on the project.

**Azure Key Vault** — Provide the vault URL, tenant ID, client ID, and client secret. The service principal needs `Get`, `Set`, `Delete`, and `List` permissions on the Key Vault's secrets.

**HashiCorp Vault** — Provide the Vault address, namespace (if using namespaces), a Vault token, and the KV v2 mount path. The token needs read/write access to the specified KV v2 mount.

**Self-hosted OpenBAO** — Same configuration as HashiCorp Vault (they share the same API).

Backend connection credentials are themselves encrypted and stored securely — they never appear in plaintext in Planton's database.

### The Default Backend

Each organization has a **default backend**. When you create a secret without explicitly choosing a backend, the default is used. The platform backend is the initial default. You can change the default to any backend you have configured.

Once a secret is created with a specific backend, the backend cannot be changed. To move a secret to a different backend, create a new secret on the target backend and copy the value.

<!-- SCREENSHOT: Secret backend configuration
  Page: /orgs/{org}/secrets (backend settings)
  Action: Show the backend configuration form with encryption options visible
  Focus: The encryption key source selection and CMEK configuration fields
  Alt: Secret backend configuration showing encryption key source dropdown with platform-managed and CMEK options
-->

## Choosing the Right Backend

For most organizations, the right choice is simple:

| Scenario | Backend |
|----------|---------|
| **Getting started** | Built-in (default) |
| **Enterprise compliance** | Built-in or custom |
| **Data residency requirements** | Custom backend in your region |
| **Existing Vault infrastructure** | HashiCorp Vault or OpenBAO |
| **Multi-cloud secrets consolidation** | Built-in |

Start with the defaults. The built-in backend provides production-ready security with zero configuration. Add custom backends when compliance or policy requires it.

## Related Documentation

- [Managing Secrets](/docs/secrets/managing-secrets) — Creating and managing secrets
- [Encryption](/docs/secrets/encryption) — Envelope encryption and customer-managed keys
- [Variables](/docs/variables) — Non-sensitive configuration management
- [Runner: Security Model](/docs/runner/security-model) — How secrets are protected during execution
- [Connections: Cloud Providers](/docs/connections/cloud-providers) — Connecting cloud accounts where custom backends live
