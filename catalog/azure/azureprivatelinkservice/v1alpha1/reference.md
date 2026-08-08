# AzurePrivateLinkService

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzurePrivateLinkService
metadata:
  name: test-private-link-service
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: orders-api
  # NAT addresses consumer traffic is source-NATed through. The subnet
  # must have privateLinkServiceNetworkPoliciesEnabled: false. Exactly
  # one configuration is primary.
  natIpConfigurations:
    - name: nat-1
      subnetId:
        value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/virtualNetworks/test-vnet/subnets/pls-nat
      primary: true
  # The classic shape: the service fronts a STANDARD internal load
  # balancer's frontend (exactly one of this or destinationIpAddress).
  loadBalancerFrontendIpConfigurationIds:
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/loadBalancers/orders-lb/frontendIPConfigurations/internal
  # Discoverable by one partner subscription; their connections still
  # wait for manual approval (no autoApprovalSubscriptionIds).
  visibilitySubscriptionIds:
    - "8158df85-3d6b-4d9f-8a3c-247b63cab0a8"
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.natIpConfigurations` | `[]AzurePrivateLinkServiceNatIpConfiguration` | yes |  |  |
| `spec.natIpConfigurations[].name` | `string` | yes |  |  |
| `spec.natIpConfigurations[].subnetId` | `string \| valueFrom` | yes |  | AzureSubnet (`status.outputs.subnet_id`) |
| `spec.natIpConfigurations[].privateIpAddress` | `string` |  |  |  |
| `spec.natIpConfigurations[].privateIpAddressVersion` | `string` |  | `IPv4` |  |
| `spec.natIpConfigurations[].primary` | `bool` |  |  |  |
| `spec.loadBalancerFrontendIpConfigurationIds` | `[]string \| valueFrom` |  |  | AzureLoadBalancer (`status.outputs.frontend_ip_configuration_ids`) |
| `spec.destinationIpAddress` | `string` |  |  |  |
| `spec.proxyProtocolEnabled` | `bool` |  |  |  |
| `spec.autoApprovalSubscriptionIds` | `[]string` |  |  |  |
| `spec.visibilitySubscriptionIds` | `[]string` |  |  |  |
| `spec.fqdns` | `[]string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

- rule: Private Link Service names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.natIpConfigurations

`[]AzurePrivateLinkServiceNatIpConfiguration` · required

- rule: {"repeated":{"minItems":"1","maxItems":"8"}}

### spec.natIpConfigurations[].name

`string` · required

- rule: NAT configuration names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.natIpConfigurations[].subnetId

`string | valueFrom` · required

- references: AzureSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.natIpConfigurations[].privateIpAddress

`string`

- rule: private_ip_address must be an IPv4 address from the subnet's range (or empty for dynamic assignment)

### spec.natIpConfigurations[].privateIpAddressVersion

`string` · optional (explicit presence)

- default: `IPv4`
- rule: {"string":{"in":["IPv4"]}}

### spec.natIpConfigurations[].primary

`bool`

### spec.loadBalancerFrontendIpConfigurationIds

`[]string | valueFrom`

- references: AzureLoadBalancer (`status.outputs.frontend_ip_configuration_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureLoadBalancer, name: <that resource's name>, fieldPath: status.outputs.frontend_ip_configuration_ids}} -- a bare string does not parse

### spec.destinationIpAddress

`string`

- rule: destination_ip_address must be an IPv4 address

### spec.proxyProtocolEnabled

`bool`

### spec.autoApprovalSubscriptionIds

`[]string`

- rule: {"repeated":{"items":{"string":{"uuid":true}}}}

### spec.visibilitySubscriptionIds

`[]string`

- rule: {"repeated":{"items":{"cel":[{"id":"visibility_entry_format","message":"Each visibility entry is a subscription UUID, or the single wildcard \"*\" for public discoverability","expression":"this == '*' || this.matches('^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$')"}]}}}

### spec.fqdns

`[]string`

- rule: {"repeated":{"items":{"string":{"minLen":"1"}}}}

### spec.tags

`map<string, string>`

## Validation Rules

- `exactly_one_destination`: Point the service at exactly one destination: load_balancer_frontend_ip_configuration_ids (service behind a Standard load balancer) or destination_ip_address (NAT to one fixed IP)
- `exactly_one_primary_nat_configuration`: Mark exactly one nat_ip_configuration as primary -- ARM requires a single primary NAT address
- `nat_configuration_names_unique`: NAT IP configuration names must be unique on the service

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzurePrivateLinkService, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.private_link_service_id` | `string` |  |
| `status.outputs.private_link_service_name` | `string` |  |
| `status.outputs.alias` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.natIpConfigurations[].subnetId` | AzureSubnet | `status.outputs.subnet_id` |
| `spec.loadBalancerFrontendIpConfigurationIds` | AzureLoadBalancer | `status.outputs.frontend_ip_configuration_ids` |

## See Also

- [Overview](../README.md)
