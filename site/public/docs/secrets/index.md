---
title: "Secrets"
description: "Secrets that live in your own vault — readable names, live version history, full visibility in your cloud console, and copyable integration snippets for wherever the value needs to go"
icon: key
order: 45
tags:
  - Secrets
  - Security
---

# Secrets

Planton includes a built-in secrets manager that gives your organization a single home for sensitive values — API keys, database passwords, service account tokens, private certificates — with one promise most platforms cannot make: **you choose the vault, and Planton never takes your secrets hostage.**

A secret stored through Planton into your AWS, GCP, Azure, Vault, or OpenBAO is a first-class citizen of that store. It carries a name anyone can read at a glance, labels you can filter on, and a one-click deep link to the provider's own console. A teammate can update the value right there, and Planton shows that edit the moment anyone looks. External Secrets Operator, Cloud Run, ECS — anything that reads your store natively — reads exactly what Planton wrote, from day one, with no Planton dependency.

In the web console, secrets are managed under your organization's **Settings → Secrets**.

## The Problem

In most DevOps environments, sensitive values are scattered:

- AWS access keys live in CI/CD environment variables
- Database passwords are duplicated across Kubernetes Secrets in each cluster
- API tokens are hardcoded in deployment manifests
- When someone rotates a credential, they have to track down every place it was used

And most secrets managers add a second problem while solving the first: they store your values in a proprietary format, in their own database — so every consumer of the value now depends on their API, and leaving means migrating every secret.

## What Planton Offers

**Your store, your values** — Real backends hold the real value: text secrets byte-for-byte, key-value secrets as one clean JSON object. Your own tools read them natively. See [Where Secrets Live](/docs/secrets/where-secrets-live).

**Readable, collision-proof naming** — Every secret's remote name encodes the organization, scope, and environment (`planton/acme/env/prod/db-password`, rendered per provider's naming rules). Same-named secrets in two environments are two secrets, structurally.

**Live version history** — The provider is the source of truth. A version added directly in the cloud console appears in Planton's timeline immediately, and it is what the next deployment resolves. See [Version History](/docs/secrets/versions).

**Just-in-time resolution** — Deployments carry references, never values. References become plaintext on the [Planton Runner](/docs/runner) inside your infrastructure, at the moment of use. Planton's own database never stores a secret value.

**Read auditing** — Every read of a secret's value is recorded the moment it happens: who read it, through which surface, and exactly which version was served.

**Integration snippets** — Every secret's detail page offers ready-to-paste snippets for External Secrets Operator, Kubernetes CSI, Cloud Run, ECS, the provider's own CLI, and Planton's own manifests. See [Consuming Secrets Anywhere](/docs/secrets/integrations).

**Organization and environment scoping with real governance** — Environment secrets are governed by environment-level access control: a teammate scoped to staging manages staging's secrets and only staging's, and a protected production environment protects production's secrets too.

## Secret Backends

Secrets are stored in a dedicated backend. Every organization starts with a zero-configuration default; organizations that want secrets in their own accounts connect their own:

| Backend | Where the value lives |
|---------|-----------------------|
| **Built-in default** | Hosted organizations get a Planton-operated OpenBAO vault (native KV values, per-organization isolation); a local desktop instance gets a built-in store, envelope-encrypted in the instance's own database |
| **AWS Secrets Manager** | Your AWS account — native values, tags, and hierarchy |
| **GCP Secret Manager** | Your GCP project — native values and labels |
| **Azure Key Vault** | Your Azure tenant — native values and tags |
| **HashiCorp Vault / OpenBAO** | Your vault, native KV v2 entries |

See [Secret Backends](/docs/secrets/backends) for setup, ambient authentication, and the verify-before-save health probe.

## How Secrets Are Used

Services, infrastructure deployments, and connections reference secrets by name with the `$secret/` grammar — `$secret/stripe-api-key`, or `$secret/@production/db-password` for an environment-scoped secret. References resolve at deployment time on the Runner; the CLI mirrors the same resolution for local development:

```bash
# View resolved configuration for a service
planton service env run

# Generate .env files for local development
planton service env pull
```

## Getting Started

1. **Create your first secret** — Settings → Secrets in the web console, or `planton secret set my-api-key 'the-value'`. The default backend is ready out of the box.
2. **See where it lives** — the secret's detail page shows its real remote name, its logical path, and a deep link to your provider's console.
3. **Wire it in** — copy the integration snippet for wherever the value needs to go, or reference it from service and infrastructure manifests with `$secret/...`.
4. **(Optional) Bring your own backend** — if compliance requires secrets to stay in your own account, connect your store and verify it before saving.

## Section Contents

- [Managing Secrets](/docs/secrets/managing-secrets) — Creating, referencing, revealing, and auditing sensitive values
- [Secret Backends](/docs/secrets/backends) — Where secrets are stored and how to configure and verify backends
- [Where Secrets Live](/docs/secrets/where-secrets-live) — Naming, labels, remote identity, and the no-hostage storage model
- [Version History](/docs/secrets/versions) — The live provider-backed timeline, out-of-band edits, and rollback
- [Consuming Secrets Anywhere](/docs/secrets/integrations) — ESO, Kubernetes CSI, Cloud Run, ECS, provider CLIs, and Planton manifests

## Related Documentation

- [Variables](/docs/variables) — Non-sensitive configuration management
- [Security Overview](/docs/security) — Platform-wide security architecture
- [Runner](/docs/runner) — The secure execution agent where references become values
