# Production HA Kubernetes Cluster

This preset creates a production-grade DigitalOcean Kubernetes (DOKS) cluster with a highly available control plane, an autoscaling default node pool, automatic patch upgrades in a Sunday-night maintenance window, DigitalOcean Container Registry integration, and a control-plane firewall restricting API server access. This is the recommended starting point for any production Kubernetes workload on DigitalOcean.

## When to Use

- Production applications requiring high availability and automatic recovery
- Workloads that need node autoscaling based on scheduling pressure
- Teams using DigitalOcean Container Registry for private images
- Environments requiring restricted Kubernetes API access

## Key Configuration Choices

- **HA control plane** (`highlyAvailable: true`) -- multiple control-plane replicas eliminate the single point of failure for the Kubernetes API. HA is one-way: it cannot be turned off once enabled, and it carries an additional monthly cost.
- **Auto-upgrade with a pinned window** (`autoUpgrade: true` + `maintenancePolicy`) -- patch releases apply automatically during Sunday 03:00 UTC. Adjust the day and start time to your low-traffic hours.
- **Registry integration** (`registryIntegration: true`) -- pulls from the account's DigitalOcean Container Registry work cluster-wide without imagePullSecrets.
- **Control-plane firewall** (`controlPlaneFirewall`) -- restricts `kubectl` and API access to the listed addresses. Replace `203.0.113.0/24` with your management network's CIDR; plain IPs are also accepted.
- **Autoscaling node pool** (`autoScale: true`, 2-5 nodes of `s-4vcpu-8gb`) -- scales with pod scheduling pressure. Changing the pool's `size` later replaces the entire cluster, so size it deliberately; add growth with separate `DigitalOceanKubernetesNodePool` resources instead.
- **VPC reference** (`vpc.valueFrom`) -- references a `DigitalOceanVpc` resource named `my-vpc` and resolves to its exported `vpc_id`; rename it to your VPC resource, or replace the block with `value: <uuid>` for an unmanaged VPC.

## Related Presets

- **02-development** -- Use instead for dev/test clusters where HA and autoscaling are unnecessary
