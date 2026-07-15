---
title: "Global Table (Multi-Region Active-Active)"
description: "This preset creates a DynamoDB Global Tables v2 deployment: the table in its home region plus one full read/write replica in a second region. Applications write to whichever region is closest and..."
type: "preset"
rank: "03"
presetSlug: "03-global-table"
componentSlug: "dynamodb"
componentTitle: "DynamoDB"
provider: "aws"
icon: "package"
order: 3
---

# Global Table (Multi-Region Active-Active)

This preset creates a DynamoDB Global Tables v2 deployment: the table in its home region plus one full read/write replica in a second region. Applications write to whichever region is closest and DynamoDB replicates changes both ways (typically within a second). Streams with the `NEW_AND_OLD_IMAGES` view are enabled because global tables require them.

## When to Use

- Multi-region applications that need low-latency reads AND writes in every region
- Disaster-recovery postures where a whole-region outage must not lose the database
- Global user bases where a single-region table forces cross-continent round trips

## Key Configuration Choices

- **On-demand billing** (`billingMode: PAY_PER_REQUEST`) -- The recommended pairing for global tables: each region scales independently with its own traffic
- **Streams with NEW_AND_OLD_IMAGES** -- Required by global tables; the replication machinery consumes the change stream
- **One replica** (`replicas`) -- Add more entries for more regions; adding or removing an entry adds or removes that region's replica in place
- **Per-replica safety** -- Point-in-time recovery and deletion protection are set on the replica independently of the home table; `propagateTags: true` keeps replica tags in sync
- **Eventual consistency** (the default) -- Last-writer-wins across regions. For Multi-Region Strong Consistency set `consistencyMode: STRONG` on exactly two replicas, or on one replica plus a `globalTableWitness` region

## Placeholders to Replace

- `<aws-region>` -- The table's home region (e.g. `us-west-2`)
- `<replica-region>` -- The replica's region (e.g. `eu-west-1`); must differ from the home region

## Related Presets

- **01-on-demand-simple** -- Use instead for single-region workloads
- **02-provisioned-production** -- Use instead for sustained single-region traffic on reserved capacity
