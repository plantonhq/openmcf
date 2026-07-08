---
title: "Human Operator Grant"
description: "This preset grants a human user -- or an Entra group, covering a whole on-call rotation with one assignment -- data-plane access to a Managed Redis instance. Personal, auditable access with no shared..."
type: "preset"
rank: "02"
presetSlug: "02-human-operator"
componentSlug: "managed-redis-access-policy-assignment"
componentTitle: "Managed Redis Access Policy Assignment"
provider: "azure"
icon: "package"
order: 2
---

# Human Operator Grant

This preset grants a human user -- or an Entra group, covering a whole
on-call rotation with one assignment -- data-plane access to a Managed
Redis instance. Personal, auditable access with no shared key in a team
vault.

## When to Use

- On-call engineers who need redis-cli access during incidents
- Debugging sessions against staging or production caches
- Replacing a shared access key passed around in a password manager

## Key Configuration Choices

- **A literal object ID** -- find a user's object ID in Microsoft Entra
  ID > Users, or a group's under Groups; the connection username IS the
  object ID and every command is attributable to it
- **Prefer granting a group** -- membership changes then never touch
  infrastructure manifests
- **Revocation is deletion** -- removing this resource revokes access
  immediately; nothing to rotate

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<managed-redis-resource-name>` | The AzureManagedRedis being granted on | Your cache manifest |
| `<user-or-group-object-id>` | The Entra user or group object ID (a GUID) | Microsoft Entra ID > Users / Groups |
