# AliCloudEipAddress

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudEipAddressSpec defines the configuration for an Alibaba Cloud
Elastic IP Address (EIP).

An EIP is a static, public IPv4 address that can be associated with ECS
instances, NAT gateways, SLB/ALB/NLB load balancers, and VPN gateways.
Unlike an instance's auto-assigned public IP, an EIP persists independently
of the resource lifecycle -- it can be released from one resource and
re-associated with another without changing the address.

This component creates a standalone EIP. Association with a target resource
is handled by the downstream component (e.g., AliCloudNatGateway includes
an eip_id field for EIP association).

Provider resources:
  Terraform: alicloud_eip_address
  Pulumi:    ecs.EipAddress

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudEipAddress
metadata:
  name: alicloudeipaddress-demo
spec:
  region: cn-hangzhou
  addressName: demo-eip
  description: Demo EIP for local testing
  bandwidth: 10
  tags:
    purpose: demo
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.addressName` | `string` |  |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.bandwidth` | `int32` |  | `5` |  |
| `spec.internetChargeType` | `string` |  | `PayByTraffic` |  |
| `spec.isp` | `string` |  | `BGP` |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the EIP will be allocated.
The EIP can only be associated with resources in the same region.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.addressName

`string`

EIP display name. 1-128 characters; must start with a letter or Chinese
character; can contain digits, underscores, periods, and hyphens.
Maps to the provider field `address_name`.

- rule: {"string":{"maxLen":"128"}}

### spec.description

`string`

Human-readable description of the EIP's intended use.
2-256 characters; cannot start with http:// or https://.

### spec.bandwidth

`int32` · optional (explicit presence)

Maximum outbound bandwidth in Mbps.
When internet_charge_type is "PayByTraffic", this is the bandwidth ceiling
(you pay per GB transferred, not per Mbps reserved).
When internet_charge_type is "PayByBandwidth", this is the guaranteed
bandwidth (you pay for the full allocation regardless of usage).
Default: 5

- default: `5`
- rule: {"int32":{"lte":1000,"gte":1}}

### spec.internetChargeType

`string` · optional (explicit presence)

Metering method for the EIP.

"PayByTraffic" -- billed per GB of outbound traffic. Best for bursty or
  unpredictable workloads. The bandwidth field acts as a ceiling.
"PayByBandwidth" -- billed for the reserved bandwidth regardless of usage.
  Best for steady, high-throughput workloads.

This field is immutable after creation (ForceNew in the provider).
Default: "PayByTraffic"

- default: `PayByTraffic`
- rule: internet_charge_type must be one of: PayByTraffic, PayByBandwidth

### spec.isp

`string` · optional (explicit presence)

Internet Service Provider line type for the EIP.

"BGP" -- multi-line BGP (default, recommended for most regions).
"BGP_PRO" -- premium BGP with optimized routing for China mainland.
"ChinaTelecom", "ChinaUnicom", "ChinaMobile" -- single-carrier lines.
"ChinaTelecom_L2", "ChinaUnicom_L2", "ChinaMobile_L2" -- L2 single-carrier.
"BGP_FinanceCloud" -- BGP for Chinese finance cloud regions.
"BGP_International" -- international BGP (outside mainland China).

Availability varies by region. BGP is available in all regions.
This field is immutable after creation (ForceNew in the provider).
Default: "BGP"

- default: `BGP`
- rule: isp must be one of: BGP, BGP_PRO, ChinaTelecom, ChinaUnicom, ChinaMobile, ChinaTelecom_L2, ChinaUnicom_L2, ChinaMobile_L2, BGP_FinanceCloud, BGP_International

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).
If omitted, the EIP is placed in the account's default resource group.

### spec.tags

`map<string, string>`

Tags to apply to the EIP.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudEipAddress, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.eip_id` | `string` | The EIP allocation ID assigned by Alibaba Cloud. Referenced by downstream components (NatGateway, ALB, VPN, etc.) via StringValueOrRef for EIP association. |
| `status.outputs.ip_address` | `string` | The public IPv4 address allocated to this EIP. This is the internet-routable address that external clients connect to. |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AliCloudNatGateway | `spec.eipId` | `status.outputs.eip_id` |
| AliCloudNetworkLoadBalancer | `spec.zoneMappings[].allocationId` | `status.outputs.eip_id` |

## See Also

- [Overview](../README.md)
