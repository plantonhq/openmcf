---
title: "Spot-Flexible Workers"
description: "This preset describes compute by attributes instead of naming an instance type: any current-generation x86 type with 2-8 vCPUs and 4-16 GiB of memory qualifies. Paired with an `AwsAutoScalingGroup`..."
type: "preset"
rank: "02"
presetSlug: "02-spot-flexible"
componentSlug: "launch-template"
componentTitle: "Launch Template"
provider: "aws"
icon: "package"
order: 2
---

# Spot-Flexible Workers

This preset describes compute by attributes instead of naming an instance
type: any current-generation x86 type with 2-8 vCPUs and 4-16 GiB of
memory qualifies. Paired with an `AwsAutoScalingGroup` mixed-instances
policy, that widens Spot to dozens of pools -- the diversification that
keeps a Spot fleet alive through pool interruptions.

Note the template itself stays purchase-option-neutral: put the
On-Demand/Spot split in the group's `mixedInstancesPolicy`, where base
capacity and allocation strategy live. Add `spotOptions` here only for a
template that should be Spot-only for every consumer.

## When to Use

- Interruption-tolerant workers: CI runners, queue consumers, batch
  processors, stateless web capacity above an On-Demand base
- Fleets that should adopt new instance families automatically instead of
  chasing type lists

## Key Configuration Choices

- **`memoryMib` + `vcpuCount` ranges** -- the two required dimensions;
  everything else narrows the set
- **`instanceGenerations: [current]`** -- new families qualify as AWS
  ships them; no template edit needed
- **`bareMetal: excluded`** -- keeps the pool to virtualized types with
  fast launch times
- **No `instanceType`** -- attribute selection and an exact type are
  mutually exclusive

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<fleet-name>` | Name for the launch template | Your fleet's name (e.g., `spot-workers`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `ami-<your-ami-id>` | AMI to boot from | EC2 console or SSM public parameters |
| `<security-group-resource-name>` | Name of the AwsSecurityGroup resource | Your security-group manifest's `metadata.name` |

## Common Additions

- Switch `cpuManufacturers` to `[amazon-web-services]` with an arm64 AMI
  for a Graviton fleet
- Add `spotMaxPricePercentageOverLowestPrice` (e.g., 100) to exclude
  pathologically priced pools
- Add accelerator filters (`acceleratorTypes: [gpu]`, `acceleratorNames`)
  for GPU workers

## Related Presets

- **01-web-server** -- a fixed-type On-Demand web fleet
- **03-hardened** -- stricter protection flags and a customer-managed KMS key
