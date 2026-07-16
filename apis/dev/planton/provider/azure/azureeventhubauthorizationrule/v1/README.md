# AzureEventHubAuthorizationRule

A SAS (shared-access-signature) authorization rule for Azure Event Hubs:
a named credential with listen/send/manage rights, scoped to exactly one
of a namespace or a single event hub. Its keys and connection strings
surface as sensitive outputs -- the least-privilege alternative to the
namespace's root key.

## When to Use

Use AzureEventHubAuthorizationRule when you need:

- **Least-privilege connection strings** -- a producer holds send-only
  on its one hub; a consumer holds listen-only; neither can touch
  anything else in the namespace
- **Per-application credentials** -- one rule per app, so revoking or
  rotating one never breaks the others
- **The credential seam for integrations** -- systems that consume Event
  Hubs by connection string reference this kind's outputs

For a zero-secret posture, skip SAS entirely: disable the namespace's
`local_authentication_enabled` and grant Entra identities data-plane
roles (Azure Event Hubs Data Owner/Sender/Receiver) via
AzureRoleAssignment.

## Key Configuration

- `rule_name` -- unique within its scope, up to 60 characters (ForceNew;
  renaming regenerates keys); `RootManageSharedAccessKey` is reserved
  for the namespace's built-in root rule, whose keys already surface as
  AzureEventHubNamespace outputs
- Exactly ONE of `namespace_id` / `event_hub_id` -- the scope is the
  polymorphism; Azure models two identical ARM types and this kind
  dispatches to the right one (fixed at creation)
- `listen` / `send` / `manage` -- at least one; manage requires both
  others (Azure's own contract)

## Composition

```yaml
eventHubId:
  valueFrom:
    kind: AzureEventHub
    name: telemetry-stream
    fieldPath: status.outputs.event_hub_id
```

Applications consume `status.outputs.primary_connection_string`
(hub-scoped rules append `EntityPath={hub}` so the credential is
ready-to-use as-is).

## Documentation

- [Design research](docs/README.md) -- the one-kind-two-types verdict
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)
