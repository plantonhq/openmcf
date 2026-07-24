---
title: "Dev single node preset"
description: "The smallest declarable OpenSearch that actually serves: one pool, one node carrying every role, a small PVC, and operator-generated TLS on both the transport and HTTP layers. For developers and CI..."
type: "preset"
rank: "01"
presetSlug: "01-dev-single-node"
componentSlug: "opensearch"
componentTitle: "OpenSearch"
provider: "kubernetes"
icon: "package"
order: 1
---

# Dev single node preset

The smallest declarable OpenSearch that actually serves: one pool,
one node carrying every role, a small PVC, and operator-generated TLS
on both the transport and HTTP layers. For developers and CI who need
the real OpenSearch API surface without production ceremony.

The trade-offs are total on the durability side — one node means no
replica can be placed anywhere, so every pod restart is a brief
outage and losing the volume loses the data. And be clear-eyed about
the credentials: TLS here is real, but the bootstrapped admin login
in the `<name>-admin-password` Secret is the image's well-known demo
credential set, not a generated password. That is acceptable only
inside a private cluster with nothing sensitive in the indices. Do
not put this preset anywhere reachable from outside, and do not grow
it into production by just raising replicas — graduate to the
production preset's split pools and its security guidance instead.

The first thing to change is `version` (pin the line your clients
test against) and the pool's memory/heap pair if your documents are
large — keep heap at about half the container memory.

See [01-dev-single-node.yaml](./01-dev-single-node.yaml) for the
manifest.
