---
title: "Where Secrets Live"
sidebar_title: "Where Secrets Live"
description: "Provider-native storage under readable, collision-proof names — your values in your store, labeled, deep-linked, and never held hostage"
icon: shield
order: 40
tags:
  - Secrets
  - Storage
  - Security
---

# Where Secrets Live

Most secrets managers store your values in their own format, in their own database. Planton does the opposite: a secret stored through Planton into a real backend lives **in your store, as itself**. This page explains exactly what lands in your vault, under what name, and why that model means Planton can never hold your secrets hostage.

## Your Store Holds the Real Value

When a secret's backend is AWS Secrets Manager, GCP Secret Manager, Azure Key Vault, HashiCorp Vault, or OpenBAO:

- A **text secret** is stored byte-for-byte. What you wrote is exactly what `gcloud secrets versions access`, the AWS console, or `vault kv get` returns.
- A **key-value secret** is stored as one canonical JSON object (Vault-family stores hold native KV fields). An External Secrets Operator extractor or a Cloud Run `secretKeyRef` reads each key directly.

There is no Planton encryption wrapper around your values in your store — the provider's own encryption, IAM, and audit apply, because the value is a native citizen of the provider. Anything that can read your store can read the value with no Planton dependency, which is precisely the point: **the day you stop using Planton, every secret is already where you need it, in the format your tools expect.**

The one deliberate exception: a single-machine local instance's built-in backend stores values envelope-encrypted in the instance's own database, with the encryption key held in the operating system's keychain — because there, Planton's database IS the store.

## Readable, Collision-Proof Names

Every managed secret has a **logical path** that encodes its full identity:

```
{prefix}/{org}/org/{slug}                  # organization-scoped
{prefix}/{org}/env/{environment}/{slug}    # environment-scoped
```

The prefix defaults to `planton` and is configurable per backend. The path is rendered into each provider's naming rules:

| Provider | Rendering | Example |
|----------|-----------|---------|
| AWS Secrets Manager | verbatim (slashes make a hierarchy in the AWS console) | `planton/acme/env/prod/db-password` |
| GCP Secret Manager | segments joined with `_` | `planton_acme_env_prod_db-password` |
| Azure Key Vault | segments joined with `--` | `planton--acme--env--prod--db-password` |
| Vault / OpenBAO | native KV path | `planton/acme/env/prod/db-password` |

Because the environment is part of the address, same-named secrets in two environments are two remote secrets — structurally, not by convention. Anyone reading your store can tell at a glance what a secret is, which organization and environment it belongs to, and that Planton manages it.

## Provenance Labels

Every managed remote secret carries labels (GCP), tags (AWS, Azure), or custom metadata (Vault-family):

| Label | Value |
|-------|-------|
| `managed-by` | `planton` |
| `org` | the organization slug |
| `scope` | `organization` or `environment` |
| `env` | the environment slug (environment scope only) |
| `secret-id` | the Planton record id |

Filter your own store by `managed-by=planton` to see exactly what Planton manages — with your own tools, no Planton API involved.

## Remote Identity on Every Secret

The moment a secret is created, its record captures the **remote identity**: the real rendered name, the logical path, and a deep link to the secret in your provider's own console. The console detail page, `planton secret describe`, and agent tools all surface it. Share the deep link with a teammate who has provider IAM access and they can read or update the value right in the provider console — out-of-band edits are first-class (see [Version History](/docs/secrets/versions)).

## What Planton's Own Database Holds

Metadata only: the secret's name, scope, backend binding, remote identity, version authorship records, and the read-audit trail. **Never the value.** A full copy of Planton's database yields no secret values for provider-backed secrets — they are in your store, under your IAM, and nowhere else.

## Related Documentation

- [Secret Backends](/docs/secrets/backends) — Configuring and verifying where secrets are stored
- [Version History](/docs/secrets/versions) — The live provider-backed timeline
- [Consuming Secrets Anywhere](/docs/secrets/integrations) — Ready-to-paste snippets for your runtime
- [Security Overview](/docs/security) — Platform-wide security architecture
