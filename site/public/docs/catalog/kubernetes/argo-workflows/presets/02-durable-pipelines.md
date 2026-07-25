---
title: "Durable pipelines preset"
description: "Argo Workflows with its two durability seams filled: an S3-compatible artifact store so steps pass files and archived logs outlive pods, and a Postgres archive so run history outlives the Workflow..."
type: "preset"
rank: "02"
presetSlug: "02-durable-pipelines"
componentSlug: "argo-workflows"
componentTitle: "Argo Workflows"
provider: "kubernetes"
icon: "package"
order: 2
---

# Durable pipelines preset

Argo Workflows with its two durability seams filled: an S3-compatible
artifact store so steps pass files and archived logs outlive pods, and
a Postgres archive so run history outlives the Workflow CRs
themselves. With history safe in the database, the retention policy
keeps only a working set of CRs in the cluster — the UI still shows
everything, served from the archive.

The composition is the point: the artifact endpoint and the database
host are reference fields, so this preset's literals become
`value_from` references at a KubernetesSeaweedFs and a
KubernetesPostgres resource in a real chart — the store's credential
Secret already carries the exact `accesskey`/`secretkey` pair the
chart's selectors expect, and none of those credentials ever ride this
manifest.

Change first: on a cloud store, drop `insecure`, drop the declared
Secret and set `use_ambient_credentials` with IRSA/workload identity
annotated on the runner ServiceAccount — the keyless posture. And
create the `argo_archive` database before first boot: the controller
creates its tables, never the database.

See [02-durable-pipelines.yaml](./02-durable-pipelines.yaml) for the
manifest.
