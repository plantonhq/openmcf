---
title: "Managing Secrets"
description: "Create, scope, reference, reveal, and audit sensitive values — with environment-level governance and a read trail that answers who read what, when"
icon: key
order: 20
tags:
  - Secrets
  - Configuration
  - Audit
---

# Managing Secrets

A secret in Planton is a named container for sensitive data — API keys, database passwords, service account tokens, private certificates. Secrets are scoped to an organization or a single environment, versioned through a [live provider-backed timeline](/docs/secrets/versions), stored [provider-native in your own backend](/docs/secrets/where-secrets-live), and read-audited.

## Scoping: Organization vs Environment

**Organization secrets** are available across all environments. Use these for values that do not change between environments — a third-party API key, a shared monitoring token.

**Environment secrets** belong to a single environment — the production database password, the staging service account. Environment scoping is real governance, not a naming convention: a teammate granted access to only the staging environment can manage staging's secrets and only staging's, and marking production as protected protects production's secrets exactly as it protects production's infrastructure.

## Creating a Secret

### Using the Web Console

Settings → Secrets → **Create Secret** walks a short wizard:

1. **Scope** — organization-wide, or pick the environment.
2. **Destination** — which [backend](/docs/secrets/backends) stores the value (the organization default is preselected).
3. **Data** — a single text value, or key-value pairs. Creating the container without a value is legitimate — seed the value later or out-of-band.
4. **Review** — shows the logical path the secret will live under before anything is created.

The success screen shows the secret's **real remote name** with a copy button, a deep link to your provider's console, and the full [integration snippet catalog](/docs/secrets/integrations).

### Using the CLI

```bash
# Create or update a text secret's value
planton secret set stripe-api-key 'sk_live_...'

# Key-value secrets take KEY=VALUE pairs
planton secret set cloudflare r2.access-key-id=... r2.secret-access-key=...

# Environment-scoped secrets take the environment flag
planton secret set db-password 'the-value' --env production

# List, inspect, read
planton secret list
planton secret describe db-password    # includes the remote identity
planton secret get db-password -o plain
```

## Referencing Secrets

Anything on the platform that consumes configuration accepts the `$secret/` reference grammar:

```
$secret/<slug>                    # org-scoped text secret
$secret/<slug>/<key>              # one key of an org-scoped key-value secret
$secret/@<env>/<slug>             # environment-scoped text secret
$secret/@<env>/<slug>/<key>       # one key of an environment-scoped key-value secret
```

References live in service manifests, infrastructure resource definitions, and connection fields; the value never does. At deployment time the [Runner](/docs/runner) resolves references just-in-time inside your infrastructure — "latest" being the [backend's latest](/docs/secrets/versions) — uses the value, and discards it. Planton's own database never stores a secret value for provider-backed secrets.

For local development the CLI mirrors the same resolution:

```bash
planton service env run     # run with resolved configuration
planton service env pull    # write .env files (git-ignore verified, 0600)
planton service env check   # per-reference resolution report
```

## The Read Story

Every read of a secret's **value** — a console reveal, a CLI read, a deploy-time resolution — writes one immutable audit entry in the same breath as serving the value. Each entry records who (the person, the API key, or the exact platform work such as a specific stack job), through which surface, when, and **exactly which version was served**. A read whose record cannot be written is refused.

The feed rides the same permission as reading the value itself:

- **Console** — the Recent Access card on the secret's detail page.
- **CLI** — `planton secret access <slug>` (alias `reads`).

Entries carry coordinates and actors, never values.

## Permissions

Reading a secret's **metadata** and reading its **value** are separate permissions, so a teammate can see that a secret exists without being able to reveal it. API keys can be clamped to entries-only. Agent surfaces can write secret values but never read them back.

## Related Documentation

- [Secret Backends](/docs/secrets/backends) — Where secrets are stored
- [Where Secrets Live](/docs/secrets/where-secrets-live) — Naming, labels, and remote identity
- [Version History](/docs/secrets/versions) — The live timeline, out-of-band edits, rollback
- [Consuming Secrets Anywhere](/docs/secrets/integrations) — Ready-to-paste snippets
- [Variables](/docs/variables) — Non-sensitive configuration management
