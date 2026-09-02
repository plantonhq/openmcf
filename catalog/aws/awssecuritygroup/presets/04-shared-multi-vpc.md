# Shared Multi-VPC Security Group

**Use case:** One firewall definition attached by workloads in several VPCs —
instead of maintaining drifting per-VPC copies of the same rules.

The group lives in its primary VPC (`spec.vpcId`) and is shared into each
VPC listed under `spec.additionalVpcIds` (same AWS account and region).
Resources in any of those VPCs attach the group by its one
`security_group_id` output.

## What You Get

- One security group associated with every listed VPC
- Outputs: `security_group_id`, `security_group_arn`, `owner_id`, and
  `additional_vpc_association_ids` (one association id per shared VPC)

## Rules on shared groups

AWS ignores a rule that references another security group when the packet
traverses a VPC different from that group's — so keep shared-group rules on
**CIDR blocks or managed prefix lists** (as this preset does), which behave
identically in every VPC.

## When to Use

- A platform team maintains one baseline ingress/egress policy for many
  VPCs
- Workloads spread across VPCs (created by different charts) must present
  identical firewall posture

## Cost

- Security groups and VPC associations are **free**
