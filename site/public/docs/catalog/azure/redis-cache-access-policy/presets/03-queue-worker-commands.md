---
title: "Queue Worker with Single-Command Grants"
description: "This preset shows the finest grain the ACL syntax offers: individual command grants. The identity can push, pop, and measure one queue key -- and do literally nothing else."
type: "preset"
rank: "03"
presetSlug: "03-queue-worker-commands"
componentSlug: "redis-cache-access-policy"
componentTitle: "Redis Cache Access Policy"
provider: "azure"
icon: "package"
order: 3
---

# Queue Worker with Single-Command Grants

This preset shows the finest grain the ACL syntax offers: individual
command grants. The identity can push, pop, and measure one queue key --
and do literally nothing else.

## When to Use

- Queue producers/consumers on a shared cache (list-based work queues)
- Third-party or less-trusted workloads that touch exactly one data
  structure
- Anywhere the blast radius of a leaked token must be a single key

## Key Configuration Choices

- **Single `+command` grants** -- no category is broad enough to be
  wrong; the worker gets the six list commands it uses and
  `+@connection` for AUTH/PING
- **`~<queueKey>` with no wildcard** -- one exact key, not a prefix;
  append `:*` only if the queue shards across keys
- **Producers vs consumers** -- split further if needed: producers get
  only `+lpush +rpush`, consumers only the pops

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<cache-resource-name>` | The AzureRedisCache's Planton resource name | Your cache composition |
| `<policyName>` | The policy's name (what assignments reference) | Your naming convention, e.g. `jobs-queue-worker` |
| `<queueKey>` | The exact queue key | Your key-naming convention, e.g. `jobs:pending` |
