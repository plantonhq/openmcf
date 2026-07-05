---
title: "Application Writer Without Admin Commands"
description: "This preset defines the standard application-workload policy: full data read/write on the app's key prefix, with Redis's dangerous command category carved out."
type: "preset"
rank: "02"
presetSlug: "02-app-writer-no-admin"
componentSlug: "redis-cache-access-policy"
componentTitle: "Redis Cache Access Policy"
provider: "azure"
icon: "package"
order: 2
---

# Application Writer Without Admin Commands

This preset defines the standard application-workload policy: full data
read/write on the app's key prefix, with Redis's dangerous command
category carved out.

## When to Use

- The main application identity on a shared cache -- it reads and writes
  its own data but can never flush the cache or reconfigure the server
- Replacing "Data Contributor" grants that are broader than the
  workload needs

## Key Configuration Choices

- **`+@all -@dangerous`** -- grant-then-carve-out reads naturally: all
  commands except FLUSHALL/FLUSHDB, CONFIG, SHUTDOWN, DEBUG, and the
  rest of the dangerous category
- **`~<keyPrefix>:*`** -- even a writer stays confined to its namespace;
  a bug in one app cannot corrupt another's keys
- **Pair with an assignment** -- this policy does nothing until an
  AzureRedisCacheAccessPolicyAssignment grants it to an identity

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cache-resource-name>` | The AzureRedisCache's Planton resource name | Your cache composition |
| `<policyName>` | The policy's name (what assignments reference) | Your naming convention, e.g. `orders-writer` |
| `<keyPrefix>` | The application's key namespace | Your key-naming convention |
