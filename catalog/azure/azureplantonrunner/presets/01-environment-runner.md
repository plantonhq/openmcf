# Environment Runner

This preset deploys the minimal runner appliance: an always-on worker in
a Container App Environment that receives deploy operations from the
control plane and executes them from inside your subscription. Three
decisions -- a resource group, an environment, and a token reference --
and everything else self-configures. The 30-second decision for getting
a runner standing on Azure.

## When to Use

- First runner in an Azure subscription
- A private AKS API server or private-endpoint database needs to be
  deployed to and operated: pick a VNet-integrated environment and the
  runner reaches that network's private endpoints
- You want zero inbound exposure: the runner app exposes no ingress at
  all and only ever dials out

## Key Configuration Choices

- **The environment decides the network boundary** -- the runner
  inherits its reach from the Container App Environment it runs in; a
  VNet-integrated environment is what makes private endpoints reachable,
  with no networking fields on the runner itself
- **Token as a managed-secret reference** -- the token authorizes the
  runner to JOIN, nothing more: on first boot the runner registers
  itself with the control plane and receives its own individually
  revocable identity, and revoking the token never touches runners it
  already admitted. On Planton the platform mints the token and writes
  it at exactly this reference before the infrastructure applies.
- **No mode or replica knobs** -- everything beyond the join (work
  queue, tunnel, API endpoints) arrives in the join response, so the
  runner self-configures on arrival; it runs as exactly one replica by
  design
- **Default sizing (0.5 vCPU / 1Gi)** -- comfortable for typical IaC
  operations; see the high-capacity preset when stacks grow large

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<runner-name>` | Name for the runner appliance | Any name you choose |
| `<resource-group-resource-name>` | Name of the AzureResourceGroup resource | Your resource group manifest's `metadata.name` |
| `<container-app-environment-resource-name>` | Name of the AzureContainerAppEnvironment resource | Your environment manifest's `metadata.name` |

The `runner-token` secret slug is yours to choose -- on Planton the
platform writes the token there automatically; elsewhere, create a token
with `planton runner token create` and store it under that slug.

## Related Presets

- `02-high-capacity` -- pinned version and larger sizing for heavy IaC workloads
