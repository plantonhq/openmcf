---
title: "Human Operator Grant"
description: "This preset grants full data access (including admin commands) to a human user or an Entra group -- the break-glass and on-call path that replaces sharing the access keys."
type: "preset"
rank: "03"
presetSlug: "03-human-operator-grant"
componentSlug: "redis-cache-access-policy-assignment"
componentTitle: "Redis Cache Access Policy Assignment"
provider: "azure"
icon: "package"
order: 3
---

# Human Operator Grant

This preset grants full data access (including admin commands) to a
human user or an Entra group -- the break-glass and on-call path that
replaces sharing the access keys.

## When to Use

- On-call engineers who need redis-cli access without a shared secret
- Replacing every "the key is in the team vault" workflow: access is
  personal, auditable, and revocable per user
- Prefer granting a GROUP over individual users -- membership changes
  then need no deploys

## Key Configuration Choices

- **"Data Owner"** -- includes the dangerous/admin commands; humans
  debugging production sometimes need them, applications never do
- **Literal `objectId`** -- users and groups are not Planton resources;
  find the id with `az ad user show` / `az ad group show`
- **Connecting as a human** -- `redis-cli` with username = the object id
  (or the alias) and password = a token from
  `az account get-access-token --scope https://redis.azure.com/.default`

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cache-resource-name>` | The AzureRedisCache's Planton resource name | Your cache composition |
| `<assignmentName>` | The grant's name, unique within the cache | Convention: `oncall-data-owner` |
| `<objectId>` | The user's or group's Entra object id | `az ad user show --id <upn> --query id` or `az ad group show --group <name> --query id` |
| `<userOrGroupName>` | A readable label (doubles as an alternative Redis username) | The user's UPN or the group name |
