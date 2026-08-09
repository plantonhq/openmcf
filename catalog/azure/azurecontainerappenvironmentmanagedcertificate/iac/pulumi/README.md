# AzureContainerAppEnvironmentManagedCertificate - Pulumi Module

Pulumi implementation for the
AzureContainerAppEnvironmentManagedCertificate component.

## Architecture

```
containerapp.EnvironmentManagedCertificate (one Azure-managed certificate)
```

## Key Design Decisions

- **The validation enum matches Azure's wire values verbatim**
  (`HTTP` / `CNAME`); unspecified deploys `HTTP`, sent explicitly so
  both engines send identical request bodies.
- **Create blocks on domain-validation proof** -- the `asuid` TXT record
  and the CNAME (or HTTP routing) must exist before deploy; Azure polls
  validation for up to ~30 minutes.
- **The issued certificate attaches to the matching custom-domain
  binding asynchronously** -- the binding module tolerates that drift by
  design; this module only owns issuance.
- **PARITY-EXCEPTION on tag shape** versus the Terraform module
  (documented in both engines) -- output-neutral.

## Provider

Built via the shared `pulumiazureprovider.Get` builder -- static client
secret, keyless web identity, or ambient chain, resolved from the stack
input. Never construct the provider inline.
