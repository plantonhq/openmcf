# DDoS-Protected Edge Network

This preset hardens a network that fronts public IPs: an attached DDoS
Protection Plan (always-on traffic monitoring with adaptive mitigation for
every public IP in the network) and virtual network encryption for
VM-to-VM traffic over the Azure backbone.

The DDoS Protection Plan is a separate, org-level, separately billed Azure
resource -- one plan protects up to 100 networks across subscriptions, so
it is typically created once by a platform team and attached everywhere by
ARM ID. ARM keeps attachment and activation distinct: `enable: false`
keeps the plan attached with protection off, ready to re-activate without
re-attaching.

## When to Use

- Networks hosting internet-facing workloads (Application Gateway, Load
  Balancer, public IPs)
- Environments with compliance requirements for traffic encryption in
  transit inside the network
- Any network where volumetric-attack mitigation is worth the plan's cost

## Key Configuration Choices

- **`enable: true`** -- protection is active; flip to `false` to keep the
  attachment during cost reviews without losing the wiring
- **`encryption: ALLOW_UNENCRYPTED`** -- the only enforcement mode ARM
  currently accepts: VM sizes that support encryption encrypt, others flow
  in plaintext (`DROP_UNENCRYPTED` is defined by the API but not yet
  generally available)

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<resource-group-name>` | The resource group to create the network in | The resource group's `status.outputs.resource_group_name` |
| `<ddos-protection-plan-arm-id>` | The shared plan's full ARM ID | Azure portal → DDoS protection plans → Properties → Resource ID |
| `<cost-center>` | Your org's cost-attribution tag value | Your tagging convention |
