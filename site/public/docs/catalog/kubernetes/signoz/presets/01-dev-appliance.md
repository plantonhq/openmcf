---
title: "Dev appliance"
description: "The smallest honest SigNoz: the whole platform — UI, ingestion collector, schema migrator and the bundled ClickHouse — from a four-line manifest. The one thing that is NOT upstream-default here is..."
type: "preset"
rank: "01"
presetSlug: "01-dev-appliance"
componentSlug: "signoz"
componentTitle: "SigNoz"
provider: "kubernetes"
icon: "package"
order: 1
---

# Dev appliance

The smallest honest SigNoz: the whole platform — UI, ingestion collector,
schema migrator and the bundled ClickHouse — from a four-line manifest.
The one thing that is NOT upstream-default here is the database password:
the module generates a unique credential per install and exports it
through the `<name>-clickhouse-auth` Secret, so the chart's
publicly-documented default password never reaches a cluster.

**When to use:** local development, demos, a team's first look at
OpenTelemetry-native observability with one UI for traces, metrics and
logs.

**When to move on:** size the database and turn on alerting email with
`02-production-bundled`, or point SigNoz at a ClickHouse you operate
yourself with `03-external-clickhouse`.
