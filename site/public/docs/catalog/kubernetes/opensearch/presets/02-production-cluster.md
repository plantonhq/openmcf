---
title: "Production cluster preset"
description: "The production topology: three dedicated cluster-manager nodes (the coordination quorum, isolated from query load), three data/ingest nodes on fast persistent storage, PodDisruptionBudgets limiting..."
type: "preset"
rank: "02"
presetSlug: "02-production-cluster"
componentSlug: "opensearch"
componentTitle: "OpenSearch"
provider: "kubernetes"
icon: "package"
order: 2
---

# Production cluster preset

The production topology: three dedicated cluster-manager nodes (the
coordination quorum, isolated from query load), three data/ingest
nodes on fast persistent storage, PodDisruptionBudgets limiting both
pools to one node down at a time, drain-before-stop for rolling
operations, TLS on both layers, and the Dashboards console with two
replicas. Pools scale independently — add data nodes without touching
the quorum.

Do not use this for development (the dev preset is one node and a
fraction of the resources) or as-is for anything internet-facing.
The one thing this preset cannot do for you is credentials: the
operator bootstraps the security plugin with the image's well-known
demo admin credentials — a generated-random password does not exist
at this release. Bring your own `security.config` (your
internal_users.yml, admin certificate, and admin credentials
Secrets), or treat rotating the admin password through the security
API as step one after install, before the first index is created.

Change first: the two `storage_class` placeholders (a literal class
name, or a valueFrom reference to a KubernetesStorageClass), then the
data pool's `disk_size`/memory/heap to match your index volume —
keeping heap at about half the container memory.

See [02-production-cluster.yaml](./02-production-cluster.yaml) for
the manifest.
