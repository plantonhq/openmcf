# DigitalOcean Database Replica -- Operational Guide

What experience with this component teaches that the field reference cannot.

## Why region and size are required here

The upstream provider marks both optional ("inherit the primary's"), but reads them back unconditionally without computing them -- so an omitted value diffs on the next apply, and because region is create-only, that diff schedules a FULL REPLICA REPLACEMENT. This spec requires both fields: writing the primary's region and size explicitly is the same outcome with none of the landmine. When the primary resizes, revisit the replica's size too (it must stay >= the primary's).

## Tags replace the replica

Replica tags are create-only upstream. Changing the tag list REPLACES the replica: DigitalOcean seeds a fresh one from the primary and the old endpoint dies. No primary data is at risk, but read consumers see the endpoint churn and the reseed takes cluster-create time. Decide the tag set at birth (the modules add the standard resource-metadata tags automatically -- account for them when planning any future retag).

## Creation is as slow as a cluster

A replica seeds from the primary's backup chain; DigitalOcean even retries creation through 412 responses while a young primary's FIRST backup completes. Budget the same ~5 minutes you budget for a cluster, and expect brand-new primaries to add delay.

## Cross-region replicas and VPCs

A cross-region replica joins the REPLICA region's VPC (the primary's VPC does not span regions) -- wire `vpc` to a VPC in `region`. Both are create-only; moving a replica between regions or networks is a replace.

## What a replica is NOT

Not automatic failover -- DigitalOcean offers manual console promotion, not managed HA (the primary's `node_count` standbys are the HA story). Not a backup -- replicas follow deletes and corruption faithfully; the primary's backups are the recovery path. Not writable -- writes go to the primary, always.

## Resize choreography

`size` and `storage_size_mib` change IN PLACE (the one update path), waited to "online". Grow storage with size when the slug's default would shrink below the current allocation -- storage can never decrease.

## What is deliberately NOT here

Promote-to-primary (exists in DigitalOcean's API but is not bridged by the provider -- a recorded absence to re-evaluate), replica-level users/databases (the primary owns them), and replica-level firewall rules (trusted sources are governed at the cluster family level; whether the primary's rule set covers replica endpoints is a live-verification item).
