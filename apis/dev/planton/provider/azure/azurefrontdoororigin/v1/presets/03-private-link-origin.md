# Private Link Origin

This preset creates an origin that Front Door reaches over Azure
Private Link -- the backend never sees the public internet and can
disable public network access entirely.

## When to Use

- Locked-down architectures where the origin must only be reachable
  through Front Door (defense in depth: no direct-to-origin bypass of
  the WAF or edge policies)
- Compliance postures that forbid public endpoints on application
  backends

## Key Configuration Choices

- **PREMIUM profiles only** -- Azure rejects Private Link origins on
  Standard at apply time (the SKU lives on the profile, so the spec
  cannot check it statically)
- **`targetType` matches the backend** -- `SITES` for App Service /
  Functions, `BLOB`/`WEB` for storage, `MANAGED_ENVIRONMENTS` for
  Container Apps; omit it ONLY when the target is a Private Link
  Service fronting an internal load balancer (its ARM id is the
  attachment point)
- **`location` is the TARGET's region** -- private-link connections are
  regional even though Front Door is global
- **A manual approval completes the wiring** -- after deploy, the
  target resource's owner approves the pending connection (portal:
  Networking > Private endpoint connections); traffic flows only after
  approval

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<origin-group-resource-name>` | The AzureFrontDoorOriginGroup's Planton resource name | Your Front Door composition |
| `originName` (example value) | 2-90 chars -- rename to your convention | Your naming convention |
| `<app-name>` | The backend's hostname prefix | Your backend resource |
| `privateLink.location` (example value) | The TARGET backend's Azure region | Your backend resource |
| `privateLinkTargetId` (example value) | The backend's full ARM resource id | Your backend's id output |
