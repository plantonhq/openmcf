# Azure Network Watcher Flow Log -- Operational Guide

Judgment calls that matter when you run flow logs in production.

## Let the watcher stay invisible

Azure creates one Network Watcher per region (`NetworkWatcher_{region}` in `NetworkWatcherRG`) the moment the region hosts a virtual network, and allows exactly one. Leave the watcher fields unset -- the module resolves the singleton -- and resist the urge to "manage" the watcher as infrastructure: deleting `NetworkWatcherRG` is safe (Azure recreates on demand) but churns every flow log's parent. Address a watcher explicitly only if your subscription genuinely runs a self-managed one.

## Dedicate the storage account, or at least its lifecycle policy

Creating a flow log writes a storage lifecycle-management rule that OVERWRITES existing rules on the account. An account whose lifecycle policy someone hand-tuned will lose that tuning silently. The clean posture: a storage account dedicated to flow logs (or at minimum one with no hand-managed lifecycle rules), with `retentionPolicy.days` as the single retention dial.

## Version 2, unless something old objects

Schema version 1 is the provider default and exists for legacy consumers; version 2 adds flow state and byte/packet counters -- the fields capacity analysis and Traffic Analytics actually want. New flow logs should say `version: 2` explicitly.

## Scope the target to the question

A virtual-network flow log records EVERYTHING in the network -- comprehensive, and correspondingly voluminous in storage and (with Traffic Analytics) ingestion cost. A subnet or NIC target answers narrower questions at a fraction of the volume. Start from the question the records must answer; widen scope only when the question widens.

## NSG flow logs are history -- do not migrate onto them

Azure stopped accepting new NSG-targeted flow logs in June 2025 and retires the class entirely in September 2027. Validation here rejects NSG targets with that explanation. Existing NSG flow logs elsewhere in your estate should migrate to virtual-network scope on your schedule, not the deadline's.

## Traffic Analytics interval is a cost dial

The 60-minute interval (the default) suits audit and forensics; the 10-minute interval buys near-real-time visibility at roughly six times the processing cadence and correspondingly higher workspace ingestion. Pick 10 only when someone is actually watching.

## Pausing is cheaper than deleting

`enabled: false` stops collection while keeping the configuration and the already-written files. Deleting the flow log ends retention management -- the files linger until the (overwritten) lifecycle rule or your cleanup removes them. For temporary quiet, pause; delete only when the target's recording story is genuinely over.
