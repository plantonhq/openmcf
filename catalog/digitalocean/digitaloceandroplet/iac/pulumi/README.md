# DigitalOcean Droplet -- Pulumi Module

Deploys a `digitalocean:index/droplet:Droplet` from a `DigitalOceanDroplet` stack input: image and sizing, region and VPC placement, SSH keys, backups with a policy window, IPv6, the monitoring and web-console agents, volume attachments, tags, cloud-init user data, graceful shutdown, and resize behavior. Bridge SDK pin is `pulumi-digitalocean/sdk/v4 v4.49.0`.

## Module structure

- `main.go` -- Pulumi program entry point reading the stack input
- `module/main.go` -- `Resources()`: locals, provider, droplet
- `module/locals.go` -- stack-input references and the standard Planton label map
- `module/droplet.go` -- the droplet resource and stack-output exports
- `module/outputs.go` -- output key constants (the kind's outputs.proto contract)

## Outputs

Exactly the kind's stack-output contract, identical to the Terraform module: `droplet_id`, `ipv4_address`, `ipv6_address`, `ipv4_address_private`, `urn`, `vpc_uuid`.

## Behavior notes

These spec fields are real and the Terraform module wires them. This module fails the apply with `PARITY-EXCEPTION` if they are meaningfully set (proto zero values pass), until the Pulumi DigitalOcean SDK grows the matching args:

- **`spec.gpu_partition_mode`** -- no gpu_partition_mode field on Droplet
- **`spec.public_networking: false`** -- no public_networking field on Droplet (explicit `true` equals the API default and is honored by omission)

Re-evaluate each when the SDK exposes it.

- `region` is sent only when set (the zero enum value never becomes a slug); unset lets DigitalOcean choose.
- `droplet_agent` and `resize_disk` are forwarded only when present, so unset never flips a provider default.
- Tags are `spec.tags` plus the standard Planton labels rendered as `key:value` — the exact set the Terraform module applies.
- See the kind [GUIDE](../../GUIDE.md).
