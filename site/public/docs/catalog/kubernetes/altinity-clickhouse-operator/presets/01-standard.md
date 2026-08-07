---
title: "Standard preset"
description: "The one-per-cluster engine install: the Altinity operator watching every namespace, with real operator credentials and modest container sizing. Declare it once; every KubernetesClickHouse resource in..."
type: "preset"
rank: "01"
presetSlug: "01-standard"
componentSlug: "altinity-clickhouse-operator"
componentTitle: "Altinity ClickHouse Operator"
provider: "kubernetes"
icon: "package"
order: 1
---

# Standard preset

The one-per-cluster engine install: the Altinity operator watching
every namespace, with real operator credentials and modest container
sizing. Declare it once; every KubernetesClickHouse resource in the
cluster is reconciled by this single install.

Two defaults are deliberately overridden. The chart's own watch scope
is the operator's namespace only — almost never what a platform wants,
so this preset watches everything (`.*` — entries are regular
expressions; narrow the list on multi-tenant clusters where teams
must not see each other's databases). And the operator's credentials
— the login it uses on every ClickHouse cluster it manages — default
to a publicly documented username/password pair upstream; this preset
forces the conversation by shipping a placeholder you must replace.

CRD lifecycle needs no decision here: the chart installs the four
CRDs on first install, the built-in hook upgrades them on chart
upgrades, and uninstalling the operator never deletes them — running
ClickHouse clusters and their data survive an operator removal.

The first thing to change is the password; the second, on shared
clusters, is the watch list.

See [01-standard.yaml](./01-standard.yaml) for the manifest.
