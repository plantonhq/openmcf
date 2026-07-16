---
title: "Preset: Authorized View"
description: "Use this preset to expose a filtered or aggregated slice of sensitive data without granting readers any access to the raw tables. The view lives in a reader-facing dataset; the source dataset..."
type: "preset"
rank: "02"
presetSlug: "02-authorized-view"
componentSlug: "bigquery-table"
componentTitle: "BigQuery Table"
provider: "gcp"
icon: "package"
order: 2
---

# Preset: Authorized View

## When to Use

Use this preset to expose a filtered or aggregated slice of sensitive data
without granting readers any access to the raw tables. The view lives in a
reader-facing dataset; the source dataset authorizes the VIEW itself (not
its readers) to read the underlying data.

## What It Creates

- A logical view defined entirely by its SQL query
- No schema, partitioning, or clustering (the query defines the shape)
- Deletion protection OFF so definition changes stay frictionless

## The Authorized-View Pattern (two halves)

1. **This preset** creates the view in a reporting dataset that analysts can
   read.
2. **The source dataset** lists the view in its `access` entries:

```yaml
# On the source GcpBigQueryDataset:
access:
  - view:
      projectId: my-gcp-project
      datasetId: reporting_views
      tableId: revenue_summary
```

The view then reads the source on behalf of its own readers — no role on
the raw data required, and updating the view's query requires re-granting.

## How to Use

1. Replace the `datasetId` ref's `name` with your reader-facing dataset
2. Replace `tableId` and the `view.query` SQL
3. Add the matching authorized-view access entry on the source dataset
