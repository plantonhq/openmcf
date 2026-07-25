---
title: "SigNoz on your own ClickHouse"
description: "The composition posture: SigNoz runs the observability product, a KubernetesClickHouse (with its KubernetesAltinityOperator) runs the database — each with its own lifecycle, sizing and operations...."
type: "preset"
rank: "03"
presetSlug: "03-external-clickhouse"
componentSlug: "signoz"
componentTitle: "SigNoz"
provider: "kubernetes"
icon: "package"
order: 3
---

# SigNoz on your own ClickHouse

The composition posture: SigNoz runs the observability product, a
KubernetesClickHouse (with its KubernetesAltinityOperator) runs the
database — each with its own lifecycle, sizing and operations. Every
connection detail here is a reference to that resource's outputs: the
client Service, the logical cluster name, and the auth Secret carrying
the declared user's password. Nothing ClickHouse-related installs with
SigNoz, and no password appears anywhere in this manifest.

The one contract to honor on the ClickHouse side: declare the `signoz`
user with `access_management` and cluster-wide DDL grants — SigNoz
creates and migrates its own databases (signoz_traces, signoz_metrics,
signoz_logs and companions) and runs `ON CLUSTER` DDL against the
cluster name it is given.

**When to use:** the database deserves its own operations — independent
scaling, cold storage tiers, deep profile/quota tuning, or one ClickHouse
serving more than SigNoz.

**When to move on:** this IS the destination posture; the bundled presets
are the on-ramp.
