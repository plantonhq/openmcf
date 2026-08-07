# AliCloud DNS Zone

Deploys a public DNS zone on Alibaba Cloud's Alidns service. Adding a domain to Alidns does not purchase or transfer the domain -- it creates a hosted zone so you can manage DNS records through Planton. After creating the zone, point your domain registrar's NS records to the Alibaba Cloud DNS servers returned in the stack outputs.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **DNS Domain** -- an `alicloud_dns_domain` resource registered in Alidns with optional resource group assignment and tags

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery.

### Alibaba Cloud Account

- **Domain ownership** -- you must own or control the domain at your registrar. After creating the zone, update the registrar's NS records to point to the Alidns nameservers from the stack outputs.
- **Domain name** -- the `domainName` field cannot be changed after creation.

## Deploy

### Console

Open the deployment store, find **AliCloud DNS Zone**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields including domain name and resource group.

### CLI

Create a manifest and apply it:

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudDnsZone
metadata:
  name: example-domain
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  domainName: example.com
  remark: Production domain for platform services
  tags:
    team: platform
```

```shell
planton apply -f alicloud-dns-zone.yaml
```

This registers example.com in Alidns. After deployment, update your registrar's NS records to the `dnsServers` from the stack outputs. A Stack Job tracks the provisioning in real time.

### InfraChart

DNS zones are standalone resources with no upstream dependencies. Downstream DNS record components reference the domain name directly (not via ValueFromRef, since the domain name is a string, not an ID).

## Key Configuration

These are the most important decisions when configuring a DNS zone. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Domain name** -- The `domainName` field specifies the domain to manage (e.g., "example.com", "sub.example.com"). This is immutable after creation.

**Resource group** -- The `resourceGroupId` field places the domain in an Alibaba Cloud resource group for access control and cost attribution. Also immutable after creation.

**Domain group** -- The `groupId` field organizes domains within the Alidns console. Useful when managing many domains.

## Outputs and Dependencies

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `domain_id` | Domain ID assigned by Alibaba Cloud | API references |
| `domain_name` | Domain name as registered | AliCloudDnsRecord domain reference |
| `dns_servers` | List of Alidns nameservers | Registrar NS record configuration |
| `group_name` | Domain group name (if assigned) | Organizational display |
| `puny_code` | Punycode representation (for internationalized domains) | IDN domain management |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Standard** -- A minimal DNS zone with just the domain name. Start from the **Standard** preset.

**Organizational** -- A DNS zone with resource group assignment, description, and cost-tracking tags. Start from the **Organizational** preset.

## Works With

- [**AliCloud DNS Record**](/cloud-catalog/ali-cloud-dns-record) -- creates DNS records within this zone
- [**AliCloud Application Load Balancer**](/cloud-catalog/ali-cloud-application-load-balancer) -- ALB dns_name is a common CNAME target for records in this zone
- [**AliCloud Network Load Balancer**](/cloud-catalog/ali-cloud-network-load-balancer) -- NLB dns_name is a common CNAME target for records in this zone
