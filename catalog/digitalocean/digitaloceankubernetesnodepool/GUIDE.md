# DigitalOcean Kubernetes Node Pool -- Operational Guide

Live-earned judgment for operating additional DOKS worker pools. The [README](README.md) covers what the component models; this covers how it behaves.

## The default pool is not this kind

Every DOKS cluster carries an inline default pool that belongs to the `DigitalOceanKubernetesCluster` resource. DigitalOcean marks it with the `terraform:default-node-pool` tag and refuses to import it as a standalone pool. Grow a cluster with this kind; never try to adopt the default pool here.

## Size changes replace the pool

`size` (and `gpuPartitionMode`) are ForceNew: changing them tears the pool's nodes down and creates new ones. Workloads on the pool reschedule during the replacement. The safe pattern for a capacity-class change is blue/green at the pool level: create the new pool, drain workloads onto it (labels + selectors), then delete the old one.

## Autoscaling drift is by design

With `autoScale: true`, the live node count moves between `minNodes` and `maxNodes` without touching your manifest -- the provider suppresses the diff by comparing against the pool's actual count. `nodeCount` only seeds the pool. The spec enforces the bounds' coherence early (`minNodes >= 1`, `maxNodes >= minNodes`); DigitalOcean would reject incoherent bounds only at apply time.

## Labels and taints travel with the pool, not the nodes

DOKS reapplies the pool's labels and taints to every node it creates -- including replacements after autoscaler scale-ups and node recycles. Imperative `kubectl label node`/`kubectl taint node` edits on individual nodes are lost on the next node rotation; put them here instead.

## Tags versus labels

`tags` are DigitalOcean-side: they group the pool's Droplets for billing attribution and can be targeted by DigitalOcean Cloud Firewalls. `labels` are Kubernetes-side: they drive scheduling. The provider silently filters DOKS's own machinery tags (`k8s:*`, `terraform:*`) out of state -- never author tags with those prefixes.

## Importing an existing pool

The import id is the plain pool UUID (`doctl kubernetes cluster node-pool list <cluster-id>`); the provider recovers the owning cluster by scanning the account. Two expected non-defects on a blind round-trip: the module's Planton identity labels/tags appear as additions on the first plan (they were not on the manually created pool), and a default pool import is refused outright.

## GPU pools

`gpuPartitionMode` accepts the two AMD partition tokens (`AMD_PARTITION_MODE_SPX_NPS1`, `AMD_PARTITION_MODE_DPX_NPS2`) and only makes sense on AMD GPU size slugs. It is a Terraform-only arm today: the Pulumi DigitalOcean SDK v4.49.0 has no field for it, and the Pulumi provisioner fails loudly rather than silently dropping it. Re-evaluate when the SDK catches up.

## What is deliberately NOT here

Per-node customization (DOKS nodes are cattle -- the pool is the unit), Kubernetes-side objects (deployments, tolerations -- author them in your workload manifests), and the cluster's own settings (version, addons, maintenance -- the `DigitalOceanKubernetesCluster` kind owns those).
