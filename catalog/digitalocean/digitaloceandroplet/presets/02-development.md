# Development Droplet

This preset creates the smallest real DigitalOcean Droplet for development and testing: name, size, and image only. It deliberately omits region and VPC -- DigitalOcean picks a region with available capacity and places the Droplet in that region's default VPC -- keeping cost and configuration at the floor.

## When to Use

- Development, staging, or testing environments
- Short-lived instances for CI/CD build agents
- Prototyping and experimentation

## Key Configuration Choices

- **Minimal sizing** (`size: s-1vcpu-1gb`) -- the smallest general-purpose Droplet, sufficient for light dev workloads.
- **No region** -- omitted, so DigitalOcean chooses one with available capacity. Pin a region only when the Droplet must sit next to other resources.
- **No VPC** -- omitted, so the Droplet lands in the chosen region's default VPC (the resulting UUID is exported as the `vpc_uuid` output). Reference a `DigitalOceanVpc` when dev must mirror production networking.
- **No backups** -- `enableBackups` omitted (defaults to `false`). Dev environments are ephemeral and can be recreated.
- **Development tag** -- enables separate firewall rules for dev environments.

## Placeholders to Replace

- `metadata.name` / `dropletName` -- your Droplet's name. Everything else works as-is.

## Related Presets

- **01-production** -- use instead for production workloads: SSH keys, backups with a policy window, VPC isolation, monitoring, and larger sizing.
