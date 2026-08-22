---
title: "Secrets"
description: "Built-in secrets manager with envelope encryption, customer-managed keys, and just-in-time decryption — scoped to organizations and environments"
icon: key
order: 45
tags:
  - Secrets
  - Encryption
  - CMEK
  - Security
---

# Secrets

Planton includes a built-in secrets manager that gives your organization a single, secure home for sensitive values. API keys, database passwords, service account tokens, private certificates — anything that should never appear in plaintext in a log file, a database, or a configuration repository.

If you have used HashiCorp Vault, AWS Secrets Manager, or GCP Secret Manager, the model is familiar. Planton adds envelope encryption with customer-managed key support, immutable versioning, just-in-time decryption within your own infrastructure, and organization-wide and environment-specific scoping — all with zero infrastructure to deploy or maintain.

In the web console, secrets are managed from the **Secrets** tab within **Service Hub** in the sidebar.

## The Problem

In most DevOps environments, sensitive values are scattered:

- AWS access keys live in CI/CD environment variables
- Database passwords are duplicated across Kubernetes Secrets in each cluster
- API tokens are hardcoded in deployment manifests
- When someone rotates a credential, they have to track down every place it was used

This creates security risk (secrets in plaintext, in too many places), operational burden (credential rotation is manual and error-prone), and audit gaps (no record of what changed, when, or by whom).

## What Planton Offers

**Envelope encryption** — Every secret is encrypted before it touches any storage system. A unique AES-256 encryption key is generated for each secret version. That key is itself encrypted by a Key Encryption Key. Even with direct access to the storage backend, an attacker sees only ciphertext.

**Customer-managed keys** — By default, Planton manages the encryption key. Organizations that require cryptographic control can bring their own key from AWS KMS, GCP Cloud KMS, or Azure Key Vault. Revoking a customer-managed key renders all encrypted secrets permanently unreadable by anyone, including Planton.

**Just-in-time decryption** — Secrets needed for deployments are decrypted only within the [Planton Runner](/docs/runner), which runs in your infrastructure. The plaintext value exists only for the duration of the operation, only in your environment. Planton's control plane never sees the decrypted values.

**Immutable versioning** — Every change to a secret creates a new version with its own encryption key. Previous versions remain available for audit and rollback. There is no "someone changed the password but we don't know when or what it was before."

**Organization and environment scoping** — Organization secrets are available everywhere. Environment secrets are specific to a single environment (dev, staging, production). A production database password cannot be accidentally used in development.

## Secret Backends

Secrets are stored in a dedicated backend. Every organization starts with a built-in backend powered by [OpenBAO](https://openbao.org/) — ready to use with zero configuration. Organizations that need secrets to stay within their own cloud accounts can connect their own backend:

| Backend | Description |
|---------|-------------|
| **Built-in (OpenBAO)** | Platform-managed, zero configuration, per-organization key isolation. The default. |
| **AWS Secrets Manager** | Store secrets in your AWS account |
| **GCP Secret Manager** | Store secrets in your GCP project |
| **Azure Key Vault** | Store secrets in your Azure tenant |
| **HashiCorp Vault** | Use your existing Vault infrastructure |
| **Self-hosted OpenBAO** | Same technology as the built-in backend, operated by you |

See [Secret Backends](/docs/secrets/backends) for setup instructions and [Encryption](/docs/secrets/encryption) for how secrets are protected.

## How Secrets Are Used

### In CI/CD Pipelines

Services in [CI/CD](/docs/ci-cd) reference secrets by name. During the deployment pipeline, references are resolved and injected into the service's runtime environment:

```bash
# View resolved configuration for a service
planton service env-vars --env production

# Generate .env files for local development
planton service dot-env --env production
```

### In Infrastructure Deployments

Infrastructure deployments through [Infrastructure](/docs/infrastructure) can reference secrets for any input that requires sensitive values. The secret reference is stored in the resource definition; the actual value is resolved just-in-time within the Runner before the deployment executes.

### In Connections

[Connections](/docs/connections) that use inline credentials will support secret references for sensitive fields — storing the reference in the connection definition while the actual credential lives in the secrets manager, encrypted and versioned.

## Getting Started

1. **Create your first secret** — Navigate to Secrets in the web console or use `planton service secrets upsert`. The built-in backend and platform-managed encryption are ready out of the box.
2. **Reference from services** — Wire secrets into your service deployments.
3. **(Optional) Configure a custom backend** — If compliance requires secrets to stay in your own cloud account.
4. **(Optional) Enable CMEK** — If your organization requires control over the encryption key.

<!-- SCREENSHOT: Secrets Manager overview
  Page: /orgs/{org}/secrets
  Action: Show the secrets list with at least 3 secrets across different scopes
  Focus: The full page showing the secrets table with scope indicators
  Alt: Secrets Manager page showing a list of secrets with organization and environment scope indicators, descriptions, and backend assignments
-->

## Section Contents

- [Managing Secrets](/docs/secrets/managing-secrets) — Creating, versioning, and referencing sensitive values
- [Secret Backends](/docs/secrets/backends) — Where secrets are stored and how to configure backends
- [Encryption](/docs/secrets/encryption) — Envelope encryption, customer-managed keys, and the zero-trust model

## Related Documentation

- [Variables](/docs/variables) — Non-sensitive configuration management
- [Security Overview](/docs/security) — Platform-wide security architecture and how encryption fits into the broader trust model
- [Runner](/docs/runner) — The secure execution agent where secrets are decrypted
