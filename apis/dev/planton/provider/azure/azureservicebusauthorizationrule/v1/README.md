# AzureServiceBusAuthorizationRule

A SAS (shared-access-signature) authorization rule for Azure Service
Bus: a named credential with listen/send/manage rights, scoped to
exactly one of a namespace, a queue, or a topic. Its keys and connection
strings surface as sensitive outputs -- the least-privilege alternative
to the namespace's root key.

## When to Use

Use AzureServiceBusAuthorizationRule when you need:

- **Least-privilege connection strings** -- a sender holds send-only on
  its one queue; a worker holds listen-only; neither can touch anything
  else
- **Per-application credentials** -- one rule per app, so revoking or
  rotating one never breaks the others
- **The credential seam for integrations** -- systems that consume
  Service Bus by connection string reference this kind's outputs

For a zero-secret posture, skip SAS entirely: disable the namespace's
`local_auth_enabled` and grant Entra identities data-plane roles via
AzureRoleAssignment.

## Key Configuration

- `rule_name` -- unique within its scope (ForceNew; renaming
  regenerates keys)
- Exactly ONE of `namespace_id` / `queue_id` / `topic_id` -- the scope
  is the polymorphism; Azure models three identical ARM types and this
  kind dispatches to the right one
- `listen` / `send` / `manage` -- at least one; manage requires both
  others (Azure's own contract)

## Composition

```yaml
queueId:
  valueFrom:
    kind: AzureServiceBusQueue
    name: orders-queue
    fieldPath: status.outputs.queue_id
```

The geo-DR pairing consumes `status.outputs.authorization_rule_id` as
its alias credential source.

## Documentation

- [Design research](docs/README.md) -- the one-kind-three-types verdict
- [Presets](presets/) -- remixable starting points
- [Terraform module](iac/tf/README.md) / [Pulumi module](iac/pulumi/README.md)

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
