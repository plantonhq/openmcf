# AzureApplicationSecurityGroup -- Design Research

## The Resource

An Azure Application Security Group (`Microsoft.Network/applicationSecurityGroups`)
is a named grouping of network interfaces that NSG security rules reference
by name instead of by IP prefix, enabling address-independent
micro-segmentation. The component maps onto `azurerm_application_security_group`
(azurerm v4.x, `internal/services/network/application_security_group_resource.go`),
parity-verified against pulumi-azure v6 (`network.ApplicationSecurityGroup`).

## Field Mapping (azurerm → spec)

| azurerm | spec | Notes |
|---|---|---|
| `name` | `name` | Required, ForceNew |
| `location` | `region` | Required, ForceNew |
| `resource_group_name` | `resource_group` | FK → AzureResourceGroup |
| `tags` | `tags` | User tags merged over Planton-derived tags; **only updatable field** |
| `id` (computed) | `application_security_group_id` output | Join key for member and rule references |

The azurerm resource has exactly these four arguments -- the ASG is
deliberately member-less. This is the full v4.80 surface, not a subset.

## Decomposition Decisions

- **The group holds no members.** Azure models ASG membership from the
  member side: `network_interface_application_security_group_association`,
  a network interface's `application_security_group_ids`, a VM scale set IP
  configuration's ASG ids, and NSG rules'
  `source`/`destination_application_security_group_ids`. Planton mirrors
  this exactly -- the ASG kind is the empty anchor; every consumer
  references `application_security_group_id`. Modeling membership on the
  group would invert Azure's own grain and fight every consumer's lifecycle.
- **First-class, not folded.** The group is created once and referenced by
  many members and rules, each with an independent lifecycle -- the textbook
  split criterion. Folding it into an NSG or NIC would hide the
  segmentation anchor from the resource graph.

## Composition Seams

- `AzureNetworkInterface.application_security_group_ids` → this kind's
  `application_security_group_id` (NIC-side membership).
- `AzureVirtualMachineScaleSet` instance IP configuration ASG ids → this
  kind (scale-set membership).
- `AzureNetworkSecurityGroup` rule
  `source`/`destination_application_security_group_ids` → this kind (the
  rule targets the group).

## Lifecycle

- `name` and `region` are ForceNew: renaming or moving the group replaces
  it, orphaning every rule and NIC reference until they are re-pointed.
- `tags` updates in place (the provider's Update only touches tags).
- Create/Read/Update/Delete are simple CRUD against
  `Microsoft.Network/applicationSecurityGroups`; no polling beyond the
  standard create-or-update.

## Parity Notes

- Both engines create a single `ApplicationSecurityGroup`/
  `azurerm_application_security_group` with identical name, location,
  resource group, and merged tags.
- Tag-shape divergence (lowered enum string vs snake-case literal;
  `resource_id` omission when `metadata.id` is empty) is the family-wide
  `PARITY-EXCEPTION` documented in both modules -- output-neutral, since
  tags never feed stack outputs.
