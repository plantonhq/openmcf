---
title: "Shared-throughput Mongo collection"
description: "A MongoDB API collection that provisions NO throughput of its own -- it draws from the parent database's shared budget (both throughput fields unset). Sharded by `userId` so it still scales out..."
type: "preset"
rank: "02"
presetSlug: "02-shared-throughput"
componentSlug: "cosmos-db-mongodb-collection"
componentTitle: "Cosmos DB MongoDB Collection"
provider: "azure"
icon: "package"
order: 2
---

# Shared-throughput Mongo collection

A MongoDB API collection that provisions NO throughput of its own --
it draws from the parent database's shared budget (both throughput
fields unset). Sharded by `userId` so it still scales out across
physical partitions even while sharing RU/s.

Use for fleets of small, similarly-sized collections where one shared
database budget is more economical than per-collection dedication --
pair it with the parent database's shared-autoscale preset.

See [`02-shared-throughput.yaml`](02-shared-throughput.yaml).
