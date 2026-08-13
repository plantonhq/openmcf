---
title: "User-Assigned Identity"
description: "This preset creates a project that acts as a pre-created user-assigned identity instead of a system-assigned one -- the shape for teams whose data grants (storage containers, search indexes, Key..."
type: "preset"
rank: "02"
presetSlug: "02-user-assigned-identity"
componentSlug: "cognitive-account-project"
componentTitle: "Cognitive Account Project"
provider: "azure"
icon: "package"
order: 2
---

# User-Assigned Identity

This preset creates a project that acts as a pre-created user-assigned identity instead of a system-assigned one -- the shape for teams whose data grants (storage containers, search indexes, Key Vault secrets) must exist BEFORE the project so agents work on first run.

## When to Use

- Grants are managed by a separate platform/security team and pre-provisioned
- One identity is shared by related workloads and audited as a unit
- Recreating the project must not change the principal your grants bind to

## Key Configuration Choices

- **USER_ASSIGNED with explicit identity IDs** -- reference `AzureUserAssignedIdentity` resources so the grants compose in the same chart
- **Stable principal across replaces** -- unlike a system identity, the principal survives project recreation
- **The grants target the identity, not the account** -- agents act as the PROJECT's identity

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-cognitive-account-id>` | ARM ID of the parent account (kind AIServices, project management enabled) | `AzureCognitiveAccount` status outputs (`cognitive_account_id`), or reference it with valueFrom |
| `<your-user-assigned-identity-id>` | ARM ID of the pre-granted identity | `AzureUserAssignedIdentity` status outputs (`identity_id`), or reference it with valueFrom |
