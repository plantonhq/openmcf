# AliCloudDnsZone

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1`

AliCloudDnsZoneSpec defines the configuration for an Alibaba Cloud DNS
domain managed by the Alidns service.

A DNS domain is the top-level container for DNS records. Registering a domain
in Alidns does not purchase or transfer the domain -- it adds it to the
Alidns hosted zone so that you can create DNS records against it. After
adding the domain, point your registrar's nameserver (NS) records to the
Alibaba Cloud DNS servers returned in the stack outputs.

This component creates a single Alidns domain. DNS records within the domain
are managed by the AliCloudDnsRecord component.

## Example

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudDnsZone
metadata:
  name: aliclouddnszone-demo
spec:
  region: cn-hangzhou
  domainName: demo.example.com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.domainName` | `string` | yes |  |  |
| `spec.groupId` | `string` |  |  |  |
| `spec.remark` | `string` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region for provider initialization.
While Alidns is a global service, the provider requires a region.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.domainName

`string` · required

The domain name to manage in Alidns.
Must be a valid domain (e.g., "example.com", "sub.example.com").
Cannot be changed after creation (ForceNew).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253"}}

### spec.groupId

`string`

Alidns domain group ID for organizational grouping.
Groups help organize large numbers of domains in the Alidns console.
If omitted, the domain is placed in the default group.

### spec.remark

`string`

Remark or description for the domain.
Visible in the Alidns console; useful for noting the domain's purpose.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for access control and cost attribution.
If omitted, the domain is placed in the account's default resource group.
Cannot be changed after creation (ForceNew).

### spec.tags

`map<string, string>`

Tags to apply to the domain.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudDnsZone, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.domain_id` | `string` | The domain ID assigned by Alibaba Cloud. |
| `status.outputs.domain_name` | `string` | The domain name as registered in Alidns. |
| `status.outputs.dns_servers` | `[]string` | The DNS server names assigned by Alibaba Cloud. Point your domain registrar's NS records to these servers for Alidns to serve DNS queries for the domain. |
| `status.outputs.group_name` | `string` | The domain group name (computed from the group_id). Empty when the domain is in the default group. |
| `status.outputs.puny_code` | `string` | Punycode representation of the domain name. Populated when the domain contains internationalized (non-ASCII) characters. |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
