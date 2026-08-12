# Azure Event Grid Topic -- Operational Guide

Judgment calls that matter when you run custom topics in production.

## Pick the input schema like a public API version

The schema is create-only, and a replaced topic means a new endpoint hostname and repointed publishers -- so choose deliberately on day one. Default to CloudEvents 1.0 for anything new (the vendor-neutral standard every modern handler understands); keep Event Grid schema only when existing tooling already speaks it; reserve custom schema for retrofitting a producer whose JSON cannot change.

## The name is a region-wide claim

Topic names share the region's public DNS namespace with every other Azure customer. Namespace them like you would a storage account (`{org}-{app}-events`, not `orders`) -- both to avoid conflicts and to keep the hostname self-describing in publisher configs and firewall logs.

## Publish auth: keys to start, Entra to finish

Keys are the fast path and rotate cleanly (two keys, move-regenerate-move). For production, prefer Entra ID (grant publishers the EventGrid Data Sender role on the topic) and then set `local_auth_enabled: false` -- the keys stop working, and leaked-key risk drops to zero. Flip it only after every publisher is on Entra; the switch is immediate.

## No subscribers means silent drops, not queues

A topic does not store events: anything published with no matching subscription is evaluated and discarded. That makes "deploy topic, verify publish, wire subscriptions later" a trap in production cutovers -- create at least the dead-letter-backed subscription first, then cut publishers over.

## The IP firewall gates only the public path

Inbound IP rules apply while public network access is enabled; disabling public access ignores them and leaves private endpoints as the only publish path. For the common locked-down-but-public shape (CI publishers with known egress), keep public access on and list the egress CIDRs -- the locked-down preset is that shape.

## Identity is for delivery, not publishing

The topic's managed identity authenticates OUTBOUND delivery (to identity-secured Event Hubs, queues, storage dead-letter destinations) -- it has nothing to do with who may publish. Enable system-assigned identity when subscriptions will target locked-down destinations, and grant it there before creating those subscriptions.
