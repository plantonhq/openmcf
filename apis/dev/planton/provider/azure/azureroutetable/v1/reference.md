# AzureRouteTable

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureRouteTableSpec** defines the configuration for creating an Azure
route table: a set of user-defined routes (UDRs) that override Azure's
default system routing for the subnets it is attached to.

The classic uses are forced tunneling (send 0.0.0.0/0 to a firewall or
on-premises gateway instead of the internet), steering traffic through a
network virtual appliance (NVA), and black-holing unwanted prefixes. One
route table is typically shared by many subnets -- it has its own
lifecycle, and editing its routes changes routing for every subnet
attached to it at once.

Routes are folded inside this resource because Azure applies them only as
part of the table (a route has no life of its own). The subnet-side
attachment is expressed on AzureSubnet, matching Azure's model where a
subnet declares which route table it uses.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureRouteTable
metadata:
  name: test-route-table
  org: test-org
  env: dev
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: test-rt
  routes:
    - name: default-via-firewall
      addressPrefix: "0.0.0.0/0"
      nextHopType: VIRTUAL_APPLIANCE
      nextHopInIpAddress:
        value: "10.0.1.4"
    - name: blackhole-rfc1918
      addressPrefix: "192.168.0.0/16"
      nextHopType: NONE
  bgpRoutePropagationEnabled: false
  tags:
    cost-center: platform
    owner: network-team
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.routes` | `[]AzureRouteTableRoute` |  |  |  |
| `spec.routes[].name` | `string` | yes |  |  |
| `spec.routes[].addressPrefix` | `string` | yes |  |  |
| `spec.routes[].nextHopType` | `enum` | yes |  |  |
| `spec.routes[].nextHopInIpAddress` | `string \| valueFrom` |  |  | AzureFirewall (`status.outputs.private_ip_address`) |
| `spec.bgpRoutePropagationEnabled` | `bool` |  | `true` |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

The Azure region where the route table will be created, e.g. "eastus",
"westeurope". Must match the region of the virtual networks whose
subnets attach it; changing the region replaces the table.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.resourceGroup

`string | valueFrom` · required

The Azure resource group the route table will be created in.
Can be a literal resource-group name or a reference to an
AzureResourceGroup's name output.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.name

`string` · required

The name of the route table. Must be unique within the resource group;
1-80 characters (alphanumerics, underscores, periods, and hyphens; must
start with a letter or number and end with a letter, number, or
underscore). Changing the name replaces the table -- detaching it from
every subnet until the replacement is re-attached -- so name it durably,
after the routing policy it implements ("egress-via-firewall").

- rule: Route table names start with a letter or number, end with a letter, number, or underscore, and may contain alphanumerics, underscores, periods, and hyphens
- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.routes

`[]AzureRouteTableRoute`

The user-defined routes in this table. Each route steers traffic bound
for its address_prefix to its next hop, overriding Azure's system route
for that prefix (the most specific prefix wins; among equals, a
user-defined route beats a system route). An empty list is a valid --
and common -- starting point: attach the empty table to subnets first,
then add routes as the topology grows. Routes update in place.

- rule: next_hop_in_ip_address is required when next_hop_type is VIRTUAL_APPLIANCE and must be omitted for every other next hop type

### spec.routes[].name

`string` · required

The route's name, unique within the table; 1-80 characters
(alphanumerics, underscores, periods, and hyphens). Name it after what
the route does ("default-via-firewall", "blackhole-rfc1918").

- rule: {"required":true,"string":{"minLen":"1","maxLen":"80"}}

### spec.routes[].addressPrefix

`string` · required

The destination this route applies to: a CIDR block ("10.1.0.0/16",
"0.0.0.0/0" for the default route) or an Azure service tag
("ApiManagement", "AzureBackup") to route a whole Azure service's
prefix set at once.

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.routes[].nextHopType

`enum` · required

Where matching traffic is sent. VIRTUAL_APPLIANCE requires
next_hop_in_ip_address (the appliance's IP); the other hop types
forbid it.

- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_route_table_next_hop_type_unspecified` -- Not specified -- invalid; every route must declare its next hop.
- `VIRTUAL_NETWORK_GATEWAY` -- Route traffic to the virtual network gateway (VPN/ExpressRoute) -- the forced-tunneling hop for sending traffic on-premises.
- `VNET_LOCAL` -- Route traffic within the virtual network's own address space (restore or refine intra-VNet routing).
- `INTERNET` -- Route traffic directly to the internet.
- `VIRTUAL_APPLIANCE` -- Route traffic to a network virtual appliance (firewall, router) at next_hop_in_ip_address.
- `NONE` -- Drop the traffic (black-hole route).

### spec.routes[].nextHopInIpAddress

`string | valueFrom`

The private IP address packets are forwarded to -- the network virtual
appliance's (firewall's, router's) address. Required exactly when
next_hop_type is VIRTUAL_APPLIANCE; meaningless (and rejected) for
every other hop type. Can be a literal IP or a reference to the
appliance's address output; defaults to an AzureFirewall's
private_ip_address -- the hub-spoke seam where spoke tables send
egress through the hub firewall without hand-copying its address.
For a non-firewall appliance (an NVA VM's NIC), reference its address
output explicitly or pass the literal IP.

- references: AzureFirewall (`status.outputs.private_ip_address`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFirewall, name: <that resource's name>, fieldPath: status.outputs.private_ip_address}} -- a bare string does not parse

### spec.bgpRoutePropagationEnabled

`bool` · optional (explicit presence)

Whether routes learned from on-premises via BGP (over ExpressRoute or
VPN gateways) propagate into the subnets attached to this table. Azure
defaults to true. Disable for forced-tunneling designs where learned
routes could bypass the user-defined ones -- the standard hardening for
firewall egress subnets.

- default: `true`

### spec.tags

`map<string, string>`

Free-form tags applied to the route table, merged over the
Planton-derived resource tags (organization, environment, resource id);
a user tag with the same key wins. Tags are Azure's governance surface
-- Azure Policy enforces them and Microsoft Cost Management groups by
them -- so carry your org's ownership/cost-center conventions here.
Updatable in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureRouteTable, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.route_table_id` | `string` | The Azure Resource Manager ID of the route table. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Network/routeTables/{name} This is the primary output referenced by downstream resources via StringValueOrRef. |
| `status.outputs.route_table_name` | `string` | The name of the route table. Echoed from the spec for convenience. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.routes[].nextHopInIpAddress` | AzureFirewall | `status.outputs.private_ip_address` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureSubnet | `spec.routeTableId` | `status.outputs.route_table_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
