# AwsElasticIp

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `aws.planton.dev/v1alpha1`

AwsElasticIpSpec defines the desired configuration for an AWS Elastic IP address.

An Elastic IP (EIP) is a static, public IPv4 address allocated from Amazon's
address pool (or optionally from a Bring-Your-Own-IP pool). Unlike ephemeral
public IPs that change when an instance stops, an EIP persists until explicitly
released — making it the foundation for stable public endpoints.

Common consumers of Elastic IPs:
- Network Load Balancers (static IP per subnet via allocation_id)
- NAT Gateways (stable outbound IP for private subnets)
- EC2 instances (persistent public IP across stop/start cycles)

For the 95%+ use case, no spec fields are required — simply create the resource
to allocate a VPC EIP from Amazon's pool. The optional fields below support
advanced scenarios such as Bring-Your-Own-IP (BYOIP) and AWS Local/Wavelength
zone deployments.

All optional fields are ForceNew in the underlying provider: changing any value
requires replacing the EIP (new allocation). Treat allocated EIPs as immutable.

Credentials, region, and deployment workflow live outside this spec in stack inputs.

## Example

```yaml
apiVersion: aws.planton.dev/v1alpha1
kind: AwsElasticIp
metadata:
  name: static-egress-ip-demo
spec:
  region: us-west-2
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.publicIpv4Pool` | `string` |  |  |  |
| `spec.address` | `string` |  |  |  |
| `spec.networkBorderGroup` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

The AWS region where the resource will be created.
Example: "us-west-2", "eu-west-1"

- rule: {"string":{"minLen":"1"}}

### spec.publicIpv4Pool

`string`

EC2 IPv4 address pool identifier to allocate from. Set this to a BYOIP pool
ID when you want the EIP to come from your own registered IP address range
instead of Amazon's pool. When omitted, Amazon allocates from its own pool.

This field is ForceNew: changing it requires replacing the EIP.

### spec.address

`string`

Request a specific IP address from the BYOIP pool identified by
`public_ipv4_pool`. The address must belong to the specified pool.
Only meaningful when `public_ipv4_pool` is set.

This field is ForceNew: changing it requires replacing the EIP.

### spec.networkBorderGroup

`string`

Network border group that controls the location scope of this EIP. Used for
allocating EIPs in AWS Local Zones or Wavelength zones. When omitted, the
EIP is scoped to the Region.

Example values: "us-east-1", "us-east-1-wl1-bos-wlz-1" (Wavelength),
"us-west-2-lax-1a" (Local Zone).

This field is ForceNew: changing it requires replacing the EIP.

## Validation Rules

- `address_requires_byoip_pool`: address requires public_ipv4_pool to be set (specific IPs can only be requested from a BYOIP pool)

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AwsElasticIp, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.allocation_id` | `string` | The allocation ID of the Elastic IP (e.g., "eipalloc-0123456789abcdef0"). This is the primary identifier used to reference the EIP in other AWS resources such as NLB subnet mappings, NAT Gateways, and EIP associations. |
| `status.outputs.public_ip` | `string` | The public IPv4 address assigned to this Elastic IP. |
| `status.outputs.arn` | `string` | The Amazon Resource Name (ARN) of the Elastic IP. Used for IAM policies and resource-level permissions. |
| `status.outputs.public_dns` | `string` | The public DNS hostname associated with the Elastic IP (e.g., "ec2-1-2-3-4.compute-1.amazonaws.com"). |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AwsNatGateway | `spec.allocationId` | `status.outputs.allocation_id` |
| AwsNatGateway | `spec.secondaryAllocationIds` | `status.outputs.allocation_id` |
| AwsNatGateway | `spec.availabilityZoneAddresses[].allocationIds` | `status.outputs.allocation_id` |
| AwsNlb | `spec.subnetMappings[].allocationId` | `status.outputs.allocation_id` |
| AwsRedshiftCluster | `spec.elasticIp` | `status.outputs.public_ip` |

## See Also

- [Overview](../README.md)
