# DigitalOcean Kubernetes Cluster -- Operational Guide

Judgment calls that matter when you run DOKS clusters.

## The version field is a creation pin, not an upgrade lever

`kubernetesVersion` decides what the cluster is CREATED at. After that, both provisioners deliberately ignore drift on it, for a hard reason: the Terraform provider destroys and recreates the entire cluster whenever the configured version is LOWER than the live one, and `autoUpgrade` routinely moves the live version ahead of your pin. Without the ignore, an auto-upgraded cluster plus a stale manifest equals an accidental cluster replacement.

Practical consequences:

- Enable `autoUpgrade` and set a `maintenancePolicy` window; patch upgrades happen there.
- Minor/major upgrades are an operational action (DigitalOcean control panel or `doctl kubernetes cluster upgrade`), not a spec edit.
- Use a full version slug (`"1.33.1-do.3"`) when you care about the exact starting point; a prefix (`"1.33"`) lets DigitalOcean pick the patch.

## Size the default pool once; grow with separate pools

Changing the default pool's `size` (or `gpuPartitionMode`) does not resize the pool — it REPLACES THE ENTIRE CLUSTER, workloads included. The inline pool is the cluster's foundation; treat it that way:

- Pick a size with headroom for system pods plus your steady-state base load.
- Add capacity classes (bigger nodes, GPU nodes, tainted dedicated pools) as separate `DigitalOceanKubernetesNodePool` resources — those resize and replace independently without touching the cluster.
- With `autoScale: true`, `nodeCount` is only the starting count; the live count drifts between `minNodes` and `maxNodes` without producing configuration diffs.

## HA is a one-way door with a price tag

`highlyAvailable: true` gives the control plane multiple replicas and a real SLA, at an extra monthly cost — and it can never be turned off again. Turn it on for production clusters whose API must survive a control-plane node loss; leave it off for everything else. Note that on newer DOKS versions DigitalOcean's own default is HA ON — this component sends an explicit false when unset, so you never get a surprise HA bill.

## Firewall the control plane, but don't lock yourself out

`controlPlaneFirewall` restricts who can reach the public Kubernetes API endpoint. Two traps:

- The allowed list must include wherever `kubectl`, CI, and the provisioner itself run. Locking the provisioner out turns every subsequent apply into a timeout.
- `enabled` is explicit, so you can stage an address list with `enabled: false` and flip it on when the list is proven.

## Network placement is create-only

`vpc`, `clusterSubnet`, `serviceSubnet`, `workerSubnetUuid`, and `isolatedWorkers` are all fixed at creation. Decide them first; retrofitting means a new cluster and a workload migration. Custom pod/service CIDRs matter when the VPC peers with networks that would collide with DigitalOcean's defaults — set them then, leave them unset otherwise.

## The kubeconfig output is a credential

`kubeconfig` is raw YAML (not base64) carrying admin credentials; write it to a file, `chmod 600` it, and point `KUBECONFIG` at it. Credentials in it expire — `kubeconfigExpireSeconds` controls the validity (0 means DigitalOcean's 7-day default); re-fetching state mints fresh ones.

## Addon toggles: unset means DigitalOcean decides

Every addon field (`routingAgent`, `corednsAutoscaler`, the GPU device plugins and DRA drivers, `rdmaSharedDevicePlugin`, `p2pOciRegistryPlugin`) is a message with one `enabled` leaf. Leaving the field out defers to DigitalOcean's own default for that addon; setting it asserts the state, on or off. The AMD and NVIDIA device plugins are each mutually exclusive with their DRA drivers — the manifest rejects both together before any provisioner runs.

## Destroy-time cleanup: read before you set it

`destroyAllAssociatedResources: true` makes destroy also delete every load balancer, volume, and volume snapshot the cluster created. That is the right call for ephemeral clusters and the wrong one anywhere volumes outlive the cluster. It only acts at destroy; it is invisible until then.

## Choosing a provisioner: the Pulumi gaps

The Pulumi bridge (v4.49.0) cannot express `sso`, `isolatedWorkers`, `workerSubnetUuid`, `gpuPartitionMode`, or any addon toggle beyond `routingAgent`. The Pulumi module fails loudly when they are set — no silent drops. If the cluster needs those surfaces today, deploy it through Terraform.

## Importing an existing cluster

Import uses the bare cluster UUID and MUTATES the remote cluster when its default pool is untagged: the provider adds a `terraform:default-node-pool` marker tag to single-pool clusters and refuses multi-pool clusters until one pool carries the tag manually. Never put that tag in `spec.tags` — the provider owns it. Also expect `registryIntegration`, `kubeconfigExpireSeconds`, and `destroyAllAssociatedResources` to stay at their configured values after import: the API never reports them back.

## What is deliberately NOT here

Additional node pools are `DigitalOceanKubernetesNodePool` resources. Load balancers and volumes appear when Kubernetes Services and PersistentVolumeClaims create them — they are Kubernetes-driven, not cluster fields. Container registries are their own resource; `registryIntegration` only wires an existing one into the cluster.
