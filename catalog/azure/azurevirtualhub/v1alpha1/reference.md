# AzureVirtualHub

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.virtualWanId` | `string \| valueFrom` | yes |  | AzureVirtualWan (`status.outputs.virtual_wan_id`) |
| `spec.addressPrefix` | `string` | yes |  |  |
| `spec.sku` | `enum` |  | `STANDARD` |  |
| `spec.hubRoutingPreference` | `enum` |  | `EXPRESS_ROUTE` |  |
| `spec.branchToBranchTrafficEnabled` | `bool` |  |  |  |
| `spec.virtualRouterAutoScaleMinCapacity` | `int32` |  | `2` |  |
| `spec.routes` | `[]AzureVirtualHubRoute` |  |  |  |
| `spec.routes[].addressPrefixes` | `[]string` | yes |  |  |
| `spec.routes[].nextHopIpAddress` | `string` | yes |  |  |
| `spec.routeTables` | `[]AzureVirtualHubRouteTable` |  |  |  |
| `spec.routeTables[].name` | `string` | yes |  |  |
| `spec.routeTables[].labels` | `[]string` |  |  |  |
| `spec.routeTables[].routes` | `[]AzureVirtualHubRouteTableRoute` |  |  |  |
| `spec.routeTables[].routes[].name` | `string` | yes |  |  |
| `spec.routeTables[].routes[].destinationsType` | `enum` |  |  |  |
| `spec.routeTables[].routes[].destinations` | `[]string` | yes |  |  |
| `spec.routeTables[].routes[].nextHop` | `string \| valueFrom` | yes |  | AzureFirewall (`status.outputs.firewall_id`) |
| `spec.routeMaps` | `[]AzureVirtualHubRouteMap` |  |  |  |
| `spec.routeMaps[].name` | `string` | yes |  |  |
| `spec.routeMaps[].rules` | `[]AzureVirtualHubRouteMapRule` |  |  |  |
| `spec.routeMaps[].rules[].name` | `string` | yes |  |  |
| `spec.routeMaps[].rules[].matchCriteria` | `[]AzureVirtualHubRouteMapMatchCriterion` |  |  |  |
| `spec.routeMaps[].rules[].matchCriteria[].matchCondition` | `enum` |  |  |  |
| `spec.routeMaps[].rules[].matchCriteria[].asPath` | `[]string` |  |  |  |
| `spec.routeMaps[].rules[].matchCriteria[].community` | `[]string` |  |  |  |
| `spec.routeMaps[].rules[].matchCriteria[].routePrefix` | `[]string` |  |  |  |
| `spec.routeMaps[].rules[].actions` | `[]AzureVirtualHubRouteMapAction` |  |  |  |
| `spec.routeMaps[].rules[].actions[].type` | `enum` |  |  |  |
| `spec.routeMaps[].rules[].actions[].parameters` | `[]AzureVirtualHubRouteMapActionParameter` |  |  |  |
| `spec.routeMaps[].rules[].actions[].parameters[].asPath` | `[]string` |  |  |  |
| `spec.routeMaps[].rules[].actions[].parameters[].community` | `[]string` |  |  |  |
| `spec.routeMaps[].rules[].actions[].parameters[].routePrefix` | `[]string` |  |  |  |
| `spec.routeMaps[].rules[].nextStepIfMatched` | `enum` |  |  |  |
| `spec.bgpConnections` | `[]AzureVirtualHubBgpConnection` |  |  |  |
| `spec.bgpConnections[].name` | `string` | yes |  |  |
| `spec.bgpConnections[].peerAsn` | `int64` |  |  |  |
| `spec.bgpConnections[].peerIp` | `string` | yes |  |  |
| `spec.bgpConnections[].virtualNetworkConnectionId` | `string \| valueFrom` |  |  | AzureVirtualHubConnection (`status.outputs.virtual_hub_connection_id`) |
| `spec.routingIntent` | `AzureVirtualHubRoutingIntent` |  |  |  |
| `spec.routingIntent.name` | `string` | yes |  |  |
| `spec.routingIntent.routingPolicies` | `[]AzureVirtualHubRoutingPolicy` | yes |  |  |
| `spec.routingIntent.routingPolicies[].name` | `string` | yes |  |  |
| `spec.routingIntent.routingPolicies[].destinations` | `[]enum` | yes |  |  |
| `spec.routingIntent.routingPolicies[].nextHop` | `string \| valueFrom` | yes |  | AzureFirewall (`status.outputs.firewall_id`) |
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

- rule: {"required":true,"string":{"minLen":"1","maxLen":"256"}}

### spec.virtualWanId

`string | valueFrom` · required

- references: AzureVirtualWan (`status.outputs.virtual_wan_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualWan, name: <that resource's name>, fieldPath: status.outputs.virtual_wan_id}} -- a bare string does not parse

### spec.addressPrefix

`string` · required

- rule: {"required":true,"string":{"pattern":"^([0-9]{1,3}\\.){3}[0-9]{1,3}/([0-9]|[1-2][0-9]|3[0-2])$"}}

### spec.sku

`enum` · optional (explicit presence)

- default: `STANDARD`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_hub_sku_unspecified`
- `BASIC`
- `STANDARD`

### spec.hubRoutingPreference

`enum` · optional (explicit presence)

- default: `EXPRESS_ROUTE`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_hub_routing_preference_unspecified`
- `EXPRESS_ROUTE`
- `VPN_GATEWAY`
- `AS_PATH`

### spec.branchToBranchTrafficEnabled

`bool`

### spec.virtualRouterAutoScaleMinCapacity

`int32` · optional (explicit presence)

- default: `2`
- rule: {"int32":{"gte":2}}

### spec.routes

`[]AzureVirtualHubRoute`

### spec.routes[].addressPrefixes

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"pattern":"^([0-9]{1,3}\\.){3}[0-9]{1,3}/([0-9]|[1-2][0-9]|3[0-2])$"}}}}

### spec.routes[].nextHopIpAddress

`string` · required

- rule: {"required":true,"string":{"pattern":"^([0-9]{1,3}\\.){3}[0-9]{1,3}$"}}

### spec.routeTables

`[]AzureVirtualHubRouteTable`

- rule: Route names must be unique within the route table

### spec.routeTables[].name

`string` · required

- rule: {"required":true,"string":{"pattern":"^[^<>%&:?/+]+$"}}

### spec.routeTables[].labels

`[]string`

### spec.routeTables[].routes

`[]AzureVirtualHubRouteTableRoute`

- rule: Choose how destinations are interpreted: CIDR (prefixes, the common choice), RESOURCE_ID, or SERVICE

### spec.routeTables[].routes[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.routeTables[].routes[].destinationsType

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_hub_route_destinations_type_unspecified`
- `CIDR`
- `RESOURCE_ID`
- `SERVICE`

### spec.routeTables[].routes[].destinations

`[]string` · required

- rule: {"repeated":{"minItems":"1","items":{"string":{"minLen":"1"}}}}

### spec.routeTables[].routes[].nextHop

`string | valueFrom` · required

- references: AzureFirewall (`status.outputs.firewall_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFirewall, name: <that resource's name>, fieldPath: status.outputs.firewall_id}} -- a bare string does not parse

### spec.routeMaps

`[]AzureVirtualHubRouteMap`

- rule: Rule names must be unique within the route map

### spec.routeMaps[].name

`string` · required

- rule: {"required":true,"string":{"pattern":"^[a-zA-Z0-9]([a-zA-Z0-9_.-]{0,78}[a-zA-Z_0-9])?$"}}

### spec.routeMaps[].rules

`[]AzureVirtualHubRouteMapRule`

### spec.routeMaps[].rules[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.routeMaps[].rules[].matchCriteria

`[]AzureVirtualHubRouteMapMatchCriterion`

- rule: Choose how the criterion compares: CONTAINS, EQUALS, NOT_CONTAINS, or NOT_EQUALS
- rule: Give the criterion something to match: as_path, community, or route_prefix

### spec.routeMaps[].rules[].matchCriteria[].matchCondition

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_hub_route_map_match_condition_unspecified`
- `CONTAINS`
- `EQUALS`
- `NOT_CONTAINS`
- `NOT_EQUALS`

### spec.routeMaps[].rules[].matchCriteria[].asPath

`[]string`

### spec.routeMaps[].rules[].matchCriteria[].community

`[]string`

### spec.routeMaps[].rules[].matchCriteria[].routePrefix

`[]string`

### spec.routeMaps[].rules[].actions

`[]AzureVirtualHubRouteMapAction`

- rule: Choose the action: ADD, REMOVE, or REPLACE (transform attributes), or DROP (discard the route)
- rule: ADD, REMOVE, and REPLACE actions need at least one parameter (the attribute values to transform); only DROP takes none

### spec.routeMaps[].rules[].actions[].type

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_hub_route_map_action_type_unspecified`
- `ADD`
- `DROP`
- `REMOVE`
- `REPLACE`

### spec.routeMaps[].rules[].actions[].parameters

`[]AzureVirtualHubRouteMapActionParameter`

### spec.routeMaps[].rules[].actions[].parameters[].asPath

`[]string`

### spec.routeMaps[].rules[].actions[].parameters[].community

`[]string`

### spec.routeMaps[].rules[].actions[].parameters[].routePrefix

`[]string`

### spec.routeMaps[].rules[].nextStepIfMatched

`enum` · optional (explicit presence)

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_virtual_hub_route_map_next_step_unspecified`
- `CONTINUE`
- `TERMINATE`

### spec.bgpConnections

`[]AzureVirtualHubBgpConnection`

### spec.bgpConnections[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.bgpConnections[].peerAsn

`int64`

- rule: {"int64":{"gte":"0"}}

### spec.bgpConnections[].peerIp

`string` · required

- rule: {"required":true,"string":{"pattern":"^([0-9]{1,3}\\.){3}[0-9]{1,3}$"}}

### spec.bgpConnections[].virtualNetworkConnectionId

`string | valueFrom`

- references: AzureVirtualHubConnection (`status.outputs.virtual_hub_connection_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureVirtualHubConnection, name: <that resource's name>, fieldPath: status.outputs.virtual_hub_connection_id}} -- a bare string does not parse

### spec.routingIntent

`AzureVirtualHubRoutingIntent`

- rule: Routing policy names must be unique within the routing intent

### spec.routingIntent.name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.routingIntent.routingPolicies

`[]AzureVirtualHubRoutingPolicy` · required

- rule: {"repeated":{"minItems":"1"}}

### spec.routingIntent.routingPolicies[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.routingIntent.routingPolicies[].destinations

`[]enum` · required

- rule: {"repeated":{"minItems":"1","unique":true,"items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_virtual_hub_routing_policy_destination_unspecified`
- `INTERNET`
- `PRIVATE_TRAFFIC`

### spec.routingIntent.routingPolicies[].nextHop

`string | valueFrom` · required

- references: AzureFirewall (`status.outputs.firewall_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFirewall, name: <that resource's name>, fieldPath: status.outputs.firewall_id}} -- a bare string does not parse

### spec.tags

`map<string, string>`

## Validation Rules

- `route_table_names_unique`: Route table names must be unique on the hub
- `route_map_names_unique`: Route map names must be unique on the hub
- `bgp_connection_names_unique`: BGP connection names must be unique on the hub

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureVirtualHub, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.virtual_hub_id` | `string` |  |
| `status.outputs.virtual_hub_name` | `string` |  |
| `status.outputs.default_route_table_id` | `string` |  |
| `status.outputs.virtual_router_asn` | `int64` |  |
| `status.outputs.virtual_router_ips` | `[]string` |  |
| `status.outputs.route_table_ids` | `map<string, string>` |  |
| `status.outputs.route_map_ids` | `map<string, string>` |  |
| `status.outputs.bgp_connection_ids` | `map<string, string>` |  |
| `status.outputs.routing_intent_id` | `string` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.virtualWanId` | AzureVirtualWan | `status.outputs.virtual_wan_id` |
| `spec.routeTables[].routes[].nextHop` | AzureFirewall | `status.outputs.firewall_id` |
| `spec.bgpConnections[].virtualNetworkConnectionId` | AzureVirtualHubConnection | `status.outputs.virtual_hub_connection_id` |
| `spec.routingIntent.routingPolicies[].nextHop` | AzureFirewall | `status.outputs.firewall_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureExpressRouteGateway | `spec.virtualHubId` | `status.outputs.virtual_hub_id` |
| AzureExpressRouteGateway | `spec.connections[].routing.associatedRouteTableId` | `status.outputs.default_route_table_id` |
| AzureExpressRouteGateway | `spec.connections[].routing.inboundRouteMapId` | `status.outputs.route_map_ids` |
| AzureExpressRouteGateway | `spec.connections[].routing.outboundRouteMapId` | `status.outputs.route_map_ids` |
| AzureExpressRouteGateway | `spec.connections[].routing.propagatedRouteTable.routeTableIds` | `status.outputs.default_route_table_id` |
| AzureVirtualHubConnection | `spec.virtualHubId` | `status.outputs.virtual_hub_id` |
| AzureVirtualHubConnection | `spec.routing.associatedRouteTableId` | `status.outputs.default_route_table_id` |
| AzureVirtualHubConnection | `spec.routing.inboundRouteMapId` | `status.outputs.route_map_ids` |
| AzureVirtualHubConnection | `spec.routing.outboundRouteMapId` | `status.outputs.route_map_ids` |
| AzureVirtualHubConnection | `spec.routing.propagatedRouteTable.routeTableIds` | `status.outputs.default_route_table_id` |

## See Also

- [Overview](../README.md)
