# AzureVirtualHubConnection

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.name` | `string` | yes |  |  |
| `spec.virtualHubId` | `string \| valueFrom` | yes |  | AzureVirtualHub (`status.outputs.virtual_hub_id`) |
| `spec.remoteVirtualNetworkId` | `string \| valueFrom` | yes |  | AzureVirtualNetwork (`status.outputs.virtual_network_id`) |
| `spec.internetSecurityEnabled` | `bool` |  |  |  |
| `spec.routing` | `AzureVirtualHubConnectionRouting` |  |  |  |
| `spec.routing.associatedRouteTableId` | `string \| valueFrom` |  |  | AzureVirtualHub (`status.outputs.default_route_table_id`) |
| `spec.routing.inboundRouteMapId` | `string \| valueFrom` |  |  | AzureVirtualHub (`status.outputs.route_map_ids`) |
| `spec.routing.outboundRouteMapId` | `string \| valueFrom` |  |  | AzureVirtualHub (`status.outputs.route_map_ids`) |
| `spec.routing.propagatedRouteTable` | `AzureVirtualHubConnectionPropagatedRouteTable` |  |  |  |
| `spec.routing.propagatedRouteTable.labels` | `[]string` |  |  |  |
| `spec.routing.propagatedRouteTable.routeTableIds` | `[]string \| valueFrom` |  |  | AzureVirtualHub (`status.outputs.default_route_table_id`) |
| `spec.routing.staticVnetRoutes` | `[]AzureVirtualHubConnectionStaticVnetRoute` |  |  |  |
| `spec.routing.staticVnetRoutes[].name` | `string` | yes |  |  |
| `spec.routing.staticVnetRoutes[].addressPrefixes` | `[]string` | yes |  |  |
| `spec.routing.staticVnetRoutes[].nextHopIpAddress` | `string` | yes |  |  |
| `spec.routing.staticVnetLocalRouteOverrideCriteria` | `enum` |  | `CONTAINS` |  |
| `spec.routing.staticVnetPropagateStaticRoutesEnabled` | `bool` |  | `true` |  |

## Field Details

### spec.name

`string` · required

- rule: {"required":true,"string":{"pattern":"^[0-9a-zA-Z][-_.0-9a-zA-Z]{0,78}[_0-9a-zA-Z]$"}}

### spec.virtualHubId

`string | valueFrom` · required

- references: AzureVirtualHub (`status.outputs.virtual_hub_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.virtual_hub_id}} -- a bare string does not parse

### spec.remoteVirtualNetworkId

`string | valueFrom` · required

- references: AzureVirtualNetwork (`status.outputs.virtual_network_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualNetwork, name: <that resource's name>, fieldPath: status.outputs.virtual_network_id}} -- a bare string does not parse

### spec.internetSecurityEnabled

`bool`

### spec.routing

`AzureVirtualHubConnectionRouting`

- rule: Give the routing block something to do: an associated route table, a propagated_route_table, or static_vnet_routes -- or omit routing entirely for ARM's defaults

### spec.routing.associatedRouteTableId

`string | valueFrom`

- references: AzureVirtualHub (`status.outputs.default_route_table_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.default_route_table_id}} -- a bare string does not parse

### spec.routing.inboundRouteMapId

`string | valueFrom`

- references: AzureVirtualHub (`status.outputs.route_map_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.route_map_ids}} -- a bare string does not parse

### spec.routing.outboundRouteMapId

`string | valueFrom`

- references: AzureVirtualHub (`status.outputs.route_map_ids`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.route_map_ids}} -- a bare string does not parse

### spec.routing.propagatedRouteTable

`AzureVirtualHubConnectionPropagatedRouteTable`

- rule: Give propagation a target: labels, route_table_ids, or both

### spec.routing.propagatedRouteTable.labels

`[]string`

### spec.routing.propagatedRouteTable.routeTableIds

`[]string | valueFrom`

- references: AzureVirtualHub (`status.outputs.default_route_table_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHub, name: <that resource's name>, fieldPath: status.outputs.default_route_table_id}} -- a bare string does not parse

### spec.routing.staticVnetRoutes

`[]AzureVirtualHubConnectionStaticVnetRoute`

### spec.routing.staticVnetRoutes[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.routing.staticVnetRoutes[].addressPrefixes

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"pattern":"^([0-9]{1,3}\\.){3}[0-9]{1,3}/([0-9]|[1-2][0-9]|3[0-2])$"}}}}

### spec.routing.staticVnetRoutes[].nextHopIpAddress

`string` · required

- rule: {"required":true,"string":{"pattern":"^([0-9]{1,3}\\.){3}[0-9]{1,3}$"}}

### spec.routing.staticVnetLocalRouteOverrideCriteria

`enum` · optional (explicit presence)

- default: `CONTAINS`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_hub_connection_static_vnet_local_route_override_criteria_unspecified`
- `CONTAINS`
- `EQUAL`

### spec.routing.staticVnetPropagateStaticRoutesEnabled

`bool` · optional (explicit presence)

- default: `true`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureVirtualHubConnection, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.virtual_hub_connection_id` | `string` |  |
| `status.outputs.virtual_hub_connection_name` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.virtualHubId` | AzureVirtualHub | `status.outputs.virtual_hub_id` |
| `spec.remoteVirtualNetworkId` | AzureVirtualNetwork | `status.outputs.virtual_network_id` |
| `spec.routing.associatedRouteTableId` | AzureVirtualHub | `status.outputs.default_route_table_id` |
| `spec.routing.inboundRouteMapId` | AzureVirtualHub | `status.outputs.route_map_ids` |
| `spec.routing.outboundRouteMapId` | AzureVirtualHub | `status.outputs.route_map_ids` |
| `spec.routing.propagatedRouteTable.routeTableIds` | AzureVirtualHub | `status.outputs.default_route_table_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureVirtualHub | `spec.bgpConnections[].virtualNetworkConnectionId` | `status.outputs.virtual_hub_connection_id` |

## See Also

- [Overview](../README.md)
