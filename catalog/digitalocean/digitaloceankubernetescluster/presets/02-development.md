# Development Kubernetes Cluster

This preset creates a minimal DigitalOcean Kubernetes cluster for development and testing: a non-HA control plane and a fixed two-node pool of smaller instances, keeping costs low while still providing VPC isolation. Everything else stays at DigitalOcean's defaults -- surge upgrades on, no control-plane firewall, DigitalOcean-assigned subnets.

## When to Use

- Development, staging, or CI/CD environments
- Learning and experimentation with Kubernetes
- Short-lived clusters for feature-branch testing

## Key Configuration Choices

- **Non-HA control plane** -- `highlyAvailable` omitted (an explicit false is sent). Sufficient for non-critical workloads and avoids the HA surcharge.
- **Fixed-size node pool** -- 2 nodes with no autoscaling for predictable cost. Changing the pool's `size` later replaces the entire cluster; prefer adding a separate `DigitalOceanKubernetesNodePool` when the workload grows.
- **Smaller instances** (`size: s-2vcpu-4gb`) -- half the resources of the production preset.
- **No control-plane firewall** -- omitted for developer convenience; the API server accepts connections from anywhere. Add a `controlPlaneFirewall` block if that is not acceptable.
- **No auto-upgrade or maintenance policy** -- manual control over upgrades in dev environments.
- **VPC reference** (`vpc.valueFrom`) -- references a `DigitalOceanVpc` resource named `my-vpc` and resolves to its exported `vpc_id`; rename it to your VPC resource, or replace the block with `value: <uuid>` for an unmanaged VPC.

## Related Presets

- **01-production-ha** -- Use instead for production workloads requiring HA, autoscaling, and API server security
