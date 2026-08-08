# Regional NAT Gateway

A regional NAT gateway: one AWS-managed gateway that spans every availability zone of a VPC, replacing the classic per-AZ gateway fleet. AWS places capacity in each zone, allocates and rotates the Elastic IPs, and keeps egress zone-local -- one resource instead of one gateway, one EIP, and one route-table edit per AZ.

## When to Use

- Multi-AZ VPCs where running (and paying for) one zonal NAT gateway per AZ is operational overhead you want AWS to absorb
- Architectures that previously funneled all AZs through a single zonal gateway to save cost, accepting cross-AZ data charges and a zonal blast radius -- regional mode removes that trade-off
- New VPCs where you want NAT egress that automatically follows AZ expansion

## Key Configuration Choices

- **Regional placement** (`availabilityMode: regional`) -- the gateway references the whole VPC (`vpcId`), never a subnet. Zonal-mode inputs (`subnetId`, `allocationId`, private-IP knobs) do not apply and are rejected at validation.
- **Auto mode** -- with no `availabilityZoneAddresses`, AWS chooses the zones and manages the Elastic IPs itself. Pin zones (and, optionally, your own EIP allocations per zone) by listing `availabilityZoneAddresses` entries instead; switching a live gateway between auto and manual replaces it.
- **Compose by reference** -- `vpcId` resolves from an `AwsVpc`, so the gateway composes into a chart without hardcoding ids.
- **IAM note** -- creating a regional gateway requires the `ec2:DescribeAvailabilityZones` permission on the deploying role, on top of the usual NAT gateway permissions.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<vpc-name>` | The `AwsVpc` resource the gateway spans | Your VPC manifest's `metadata.name`, or the VPC console for the vpc-id to use as a literal `value:` instead |

## After Deploying

Point private subnets' default routes at the gateway (`target_type: nat_gateway` with the gateway's `nat_gateway_id` output) exactly as with a zonal gateway -- routing is unchanged; only the gateway's placement model differs.
