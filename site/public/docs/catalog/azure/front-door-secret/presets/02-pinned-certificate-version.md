---
title: "Pinned Certificate Version"
description: "This preset creates a Front Door secret pinned to ONE exact Key Vault certificate version -- Front Door keeps serving that version until the secret itself is replaced, no matter what rotation happens..."
type: "preset"
rank: "02"
presetSlug: "02-pinned-certificate-version"
componentSlug: "front-door-secret"
componentTitle: "Front Door Secret"
provider: "azure"
icon: "package"
order: 2
---

# Pinned Certificate Version

This preset creates a Front Door secret pinned to ONE exact Key Vault
certificate version -- Front Door keeps serving that version until the
secret itself is replaced, no matter what rotation happens in the
vault.

## When to Use

- Change-controlled environments where certificate rollout must be an
  explicit, auditable deployment -- never an automatic side effect of a
  vault renewal
- Clients that pin the served certificate (mobile apps, partner
  integrations with certificate allowlists), where an unannounced
  rotation is an outage

## Key Configuration Choices

- **`certificate_id` (versioned), not `versionless_id`** -- the version
  segment in the reference is what pins; compare the rotating preset
- **Rotation = replace the secret** -- the resource is immutable, so a
  new version means updating this reference (or the referenced
  certificate resource producing a new versioned id) and redeploying;
  domains follow automatically because they reference the secret by ARM
  id, which survives the replacement

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<front-door-profile-resource-name>` | The AzureFrontDoorProfile's Planton resource name | Your Front Door composition |
| `<key-vault-certificate-resource-name>` | The AzureKeyVaultCertificate holding the pinned version | Your certificate composition |
