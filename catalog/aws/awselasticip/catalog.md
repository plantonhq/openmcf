# AWS Elastic IP

Deploys a static public IPv4 address from Amazon's address pool or from a Bring-Your-Own-IP (BYOIP) pool, with optional network border group scoping for Local Zones and Wavelength zones. The Elastic IP integrates with Planton's Provider Connections for AWS credential management.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **Elastic IP** -- a VPC-domain static public IPv4 address allocated from Amazon's pool (default) or from a BYOIP pool when `publicIpv4Pool` is specified
- **AWS Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) applied to the Elastic IP

## Before You Deploy

### Planton Setup

- **AWS Provider Connection** -- an active connection in the Connect module with credentials for the target AWS account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials or cross-account trust authentication modes.

### AWS Account

- **Elastic IP quota** -- AWS accounts have a default limit of 5 Elastic IPs per region. Request a quota increase before allocating additional addresses. Each EIP incurs hourly charges when not associated with a running resource.
- **BYOIP pool** (optional) -- if using a Bring-Your-Own-IP pool, the pool must be provisioned and advertised in the target region before referencing it in `publicIpv4Pool`.
- **Network border group** (optional) -- required when allocating an EIP in an AWS Local Zone or Wavelength zone. Verify the target zone is enabled in your account.

## Deploy

### Console

Open the deployment store, find **AWS Elastic IP**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Standard EIP** preset in the [Presets](#presets) tab to pre-populate a default configuration.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: aws.planton.dev/v1
kind: AwsElasticIp
metadata:
  name: nlb-frontend
  org: acme-corp
  env: prod
spec:
  region: us-west-2
```

```shell
planton apply -f elastic-ip.yaml
```

This allocates a standard VPC Elastic IP from Amazon's address pool in us-west-2. No BYOIP or network border group is configured. The allocated IP persists until the Cloud Resource is destroyed -- it does not change across Stack Job runs. A Stack Job tracks the provisioning in real time.

## Key Configuration

These are the most important decisions when configuring an Elastic IP. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**BYOIP pool** -- Set `publicIpv4Pool` to a BYOIP pool ID when you need the Elastic IP to come from your own registered IP address range. Use `address` to request a specific IP from that pool. Both fields are ForceNew -- changing either requires replacing the EIP with a new allocation.

**Network border group** -- Set `networkBorderGroup` to scope the EIP to a specific AWS Local Zone or Wavelength zone instead of the default regional scope. Required when associating the EIP with resources in edge locations. This field is also ForceNew.

**Immutability** -- All optional spec fields (`publicIpv4Pool`, `address`, `networkBorderGroup`) trigger resource replacement when changed. Treat an allocated Elastic IP as immutable -- any configuration change allocates a new IP address.

## Outputs and Dependencies

### What This Component Consumes

This component has no foreign key dependencies.

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `allocation_id` | EIP allocation identifier (e.g., eipalloc-0123...) | Network Load Balancer subnet mappings, NAT Gateway attachment |
| `public_ip` | Static public IPv4 address | DNS A records, firewall allowlists, application configuration |
| `arn` | Amazon Resource Name of the EIP | IAM policies, resource-level permissions |
| `public_dns` | Public DNS hostname (e.g., ec2-1-2-3-4.compute-1.amazonaws.com) | DNS CNAME records, reverse-lookup verification |

The `allocation_id` is the primary output consumed by downstream resources. Network Load Balancers reference it in `subnetMappings` to bind a static IP per subnet. NAT Gateways reference it to provide a stable outbound IP for private subnets.

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard EIP** -- An Elastic IP from Amazon's pool with no special configuration. The most common pattern for NLB frontends, NAT Gateways, and EC2 instances needing stable public addresses. Start from the **Standard EIP** preset.

**BYOIP pool allocation** -- An Elastic IP from a customer-owned IP range, useful for maintaining IP reputation (e.g., email servers) or meeting contractual requirements for specific IP addresses. Start from the **BYOIP Pool** preset.

## Works With

This component operates independently and does not reference other deployment components.
