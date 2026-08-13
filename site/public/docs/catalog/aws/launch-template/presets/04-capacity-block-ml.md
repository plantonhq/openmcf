---
title: "Capacity Block ML Training Fleet"
description: "This preset launches into a pre-purchased Capacity Block -- AWS's reservation product for GPU capacity, bought for a defined window at a known price. `marketType: capacity-block` tells EC2 this is..."
type: "preset"
rank: "04"
presetSlug: "04-capacity-block-ml"
componentSlug: "launch-template"
componentTitle: "Launch Template"
provider: "aws"
icon: "package"
order: 4
---

# Capacity Block ML Training Fleet

This preset launches into a pre-purchased Capacity Block -- AWS's
reservation product for GPU capacity, bought for a defined window at a
known price. `marketType: capacity-block` tells EC2 this is Capacity Block
consumption, and `capacityReservation.capacityReservationId` names the
block; a launch outside the block's window, or without the target, is
rejected rather than silently billed On-Demand.

The dataset volume hydrates from its snapshot at the maximum paid rate
(300 MiB/s) so training does not stall on lazy block loads during the
paid reservation window -- exactly when slow I/O costs the most.

## When to Use

- Model training or fine-tuning runs scheduled inside a purchased
  Capacity Block window
- Any GPU workload where "could not get capacity" is worse than paying
  for a reservation

## Key Configuration Choices

- **`marketType: capacity-block`** -- the purchase market; requires the
  reservation target below
- **`capacityReservation.capacityReservationId`** -- the Capacity Block
  to launch into; blocks are targeted, never open
- **`volumeInitializationRateMibps: 300`** -- pre-warm the dataset volume
  at the fastest paid rate; only meaningful with `snapshotId`
- **`httpTokens: required`** -- IMDSv2 stays enforced on training fleets
  too

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<training-fleet-name>` | Name for the launch template | Your fleet's name (e.g., `llm-finetune`) |
| `<aws-region>` | AWS region code (e.g., `us-east-1`) | Your deployment region |
| `ami-<your-ami-id>` | GPU AMI to boot from (e.g., a Deep Learning AMI) | EC2 console or SSM public parameters |
| `<gpu-instance-type>` | The instance type the block reserves (e.g., `p5.48xlarge`) | Your Capacity Block's reservation details |
| `cr-<your-capacity-block-id>` | The Capacity Block reservation ID | EC2 console, Capacity Reservations |
| `snap-<your-dataset-snapshot-id>` | Snapshot holding the training dataset | EC2 console, Snapshots |
| `<security-group-resource-name>` | Name of the AwsSecurityGroup resource | Your security-group manifest's `metadata.name` |
| `<instance-profile-resource-name>` | Name of the AwsIamInstanceProfile resource | Your instance-profile manifest's `metadata.name` |

## Common Additions

- `placement.groupName` with a cluster placement group for multi-node
  training interconnect
- `networkInterfaces` with `enaSrd.enabled: true` for lower tail latency
  between nodes that support ENA Express
- `licenseConfigurationArns` when the training stack carries BYOL
  licensing

## Related Presets

- **02-spot-flexible** -- the opposite cost posture: interruption-tolerant
  attribute-based Spot
- **03-hardened** -- stricter protection flags and a customer-managed KMS key
