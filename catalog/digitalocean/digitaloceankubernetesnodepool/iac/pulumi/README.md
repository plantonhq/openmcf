# DigitalOcean Kubernetes Node Pool -- Pulumi Module

Deploys a `digitalocean:index/kubernetesNodePool:KubernetesNodePool` from a `DigitalOceanKubernetesNodePool` stack input: owning cluster, Droplet size, fixed or autoscaled node count, Kubernetes labels and taints, and DigitalOcean tags. Bridge SDK pin is `pulumi-digitalocean/sdk/v4 v4.49.0`.

## Module structure

- `main.go` -- Pulumi program entry point reading the stack input
- `module/main.go` -- `Resources()`: locals, provider, node pool
- `module/locals.go` -- stack-input references and the standard Planton label map
- `module/node_pool.go` -- the node-pool resource and stack-output exports
- `module/outputs.go` -- output key constants (the kind's outputs.proto contract)

## Outputs

Exactly the kind's stack-output contract, identical to the Terraform module: `node_pool_id`, `cluster_id`, `node_ids`, `droplet_ids`.

`node_ids` are the DOKS node object UUIDs (`nodes[*].id`); `droplet_ids` are the backing Droplet ids. Both engines export the same pair.

## Behavior notes

These spec fields are real and the Terraform module wires them. This module fails the apply with `PARITY-EXCEPTION` if they are meaningfully set (proto zero values pass), until the Pulumi DigitalOcean SDK grows the matching args:

- **`spec.gpu_partition_mode`** -- no `gpu_partition_mode` field on KubernetesNodePool

Re-evaluate when the SDK exposes it.

- `cluster` resolves to the owning DOKS cluster's UUID.
- Autoscaling bounds (`autoScale` / `minNodes` / `maxNodes`) are sent only when `autoScale` is true.
- Kubernetes node labels are user labels over the standard Planton identity labels — the exact map the Terraform module applies.
- Tags are `spec.tags` plus the standard Planton labels rendered as `key:value` — the exact set the Terraform module applies.
- See the kind [GUIDE](../../GUIDE.md).
