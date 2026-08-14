---
title: "Development Basic Search"
description: "The development shape: a single-replica `basic` service -- the cheapest dedicated tier, with defaults everywhere else (API-key auth, public endpoint, one partition)."
type: "preset"
rank: "02"
presetSlug: "02-dev-basic-search"
componentSlug: "ai-search-service"
componentTitle: "AI Search Service"
provider: "azure"
icon: "package"
order: 2
---

# Development Basic Search

The development shape: a single-replica `basic` service -- the
cheapest dedicated tier, with defaults everywhere else (API-key auth,
public endpoint, one partition).

## When to Use

- Development and test environments
- Small indexes (basic caps at 3 partitions / 3 replicas)
- Trying the service before sizing production

## Key Configuration Choices

- `sku: basic` -- dedicated (unlike `free`, which is one
  shared-cluster service per subscription); upgrades IN PLACE to
  standard/standard2/standard3 when the workload grows.
- Counts left unset -- both default to 1; basic caps both at 3.
- Everything else stays on provider defaults: API keys on, public
  endpoint open, no semantic ranking.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-resource-group-name>` | The resource group to create the service in | Portal -> Resource groups |
| `acme-search-dev` | Your globally-unique service name | It becomes {name}.search.windows.net |

## Related Presets

- `01-production-search` -- the production posture this grows into.
