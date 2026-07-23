# AzureFrontDoorOrigin

One backend inside an Azure Front Door origin group: where the backend
is (host name and ports), how Front Door validates its TLS certificate,
how traffic is weighted against sibling origins, and -- on PREMIUM
profiles -- whether Front Door reaches it over Private Link instead of
the public internet.

Origins are many-per-group with independent lifecycles: a regional
stamp adds its backend to a shared group, a blue/green cutover swaps
origins one at a time, and each Private Link origin carries its own
connection-approval workflow. That is why the origin is a first-class
kind rather than a list folded into the group.

## When to Use

Use AzureFrontDoorOrigin when you need:

- **A backend in a pool** -- App Service, Container Apps, Storage,
  or any reachable hostname/IP
- **Active/passive failover or weighted traffic** -- priority tiers
  and weights within a tier (canaries, gradual cutovers)
- **Private connectivity** -- Premium-profile Private Link so the
  backend disables public access entirely

## Key Configuration

- `origin_group_id` -- the parent group, referenced from an
  AzureFrontDoorOriginGroup's output; fixed at creation
- `origin_name` -- 2-90 characters, unique within the group; ForceNew
- `host_name` -- the backend's hostname or IP
- `certificate_name_check_enabled` -- default true; keep it on
  (required with Private Link, and disabling it invites
  man-in-the-middle)
- `origin_host_header` -- unset sends the origin's own hostname, which
  is what multi-tenant Azure backends (App Service, Container Apps)
  route by
- `priority` (1-5) / `weight` (1-1000) -- failover tiers and the
  traffic split within a tier
- `private_link` -- Premium-only private connectivity; the target's
  owner approves the pending connection after deploy

## Composition

```yaml
originGroupId:
  valueFrom:
    kind: AzureFrontDoorOriginGroup
    name: api-backends
    fieldPath: status.outputs.origin_group_id
```

Routes list this origin's `origin_id` output in their `originIds` to
sequence deployment after the backend exists.

## Documentation

- [Design research](docs/README.md) -- field mapping, recorded skips
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
