---
title: "Capacity Block ML Training Node"
description: "This preset launches a GPU instance into a pre-purchased EC2 Capacity Block -- reserved GPU capacity for a defined time window at a committed price, the way AWS sells scarce accelerators (P5/P4,..."
type: "preset"
rank: "05"
presetSlug: "05-capacity-block-ml"
componentSlug: "ec2-instance"
componentTitle: "EC2 Instance"
provider: "aws"
icon: "package"
order: 5
---

# Capacity Block ML Training Node

This preset launches a GPU instance into a pre-purchased EC2 Capacity Block -- reserved GPU capacity for a defined time window at a committed price, the way AWS sells scarce accelerators (P5/P4, Trainium) for ML training runs. The instance targets the block's reservation directly; when the block's window ends, AWS reclaims the capacity.

## When to Use

- Scheduled ML training runs on scarce GPU capacity purchased ahead of time
- Fine-tuning or experimentation windows where On-Demand GPU capacity is unavailable in-region
- Any workload whose capacity was bought as a Capacity Block and must land inside it

## Key Configuration Choices

- **Capacity-block market** (`marketType: capacity-block`) -- Launches into the purchased block instead of the On-Demand pool; the reservation target is required for this market and validation enforces it
- **Direct reservation target** (`capacityReservation.capacityReservationId`) -- Names the block's reservation; a resource-group ARN works when blocks are pooled
- **Instance type must match the block** (`instanceType`) -- Capacity Blocks are purchased per instance type and AZ; a mismatch fails the launch
- **High-throughput gp3 root** (`throughputMibps: 1000`) -- Checkpoint writes are the classic training bottleneck; gp3 goes to 2000 MiB/s when needed
- **At-creation volume tags** (`volumeTags`) -- Uniform tags on every volume at creation satisfy ABAC/SCP cost-governance policies without a post-creation tagging window

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<aws-region>` | Region where the Capacity Block was purchased | AWS region list |
| `cr-0123456789abcdef0` | The Capacity Block's reservation id | EC2 console > Capacity Reservations |
| `ami-0123456789abcdef0` | GPU AMI (Deep Learning AMI or your own) | AWS EC2 AMI catalog |
| `<private-subnet-id>` | Subnet in the block's availability zone | `AwsSubnet` status outputs |
| `<security-group-id>` | Security group for the training node | `AwsSecurityGroup` status outputs |
| `<instance-profile-name>` | IAM instance profile NAME (S3 dataset/checkpoint access) | `AwsIamInstanceProfile` status outputs |

## Related Presets

- **03-spot-worker** -- The other discounted market: interruption-tolerant, no reservation
- **01-ssm-managed** -- The On-Demand baseline
