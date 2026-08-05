# OciDynamicRoutingGateway

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciDynamicRoutingGatewaySpec defines the specification for an Oracle Cloud
Infrastructure Dynamic Routing Gateway (DRG) and its sub-resources.

A DRG is a virtual router that enables connectivity between a VCN and
networks outside the VCN's region: other VCNs (peering), on-premises
networks (Site-to-Site VPN via IPSec tunnels, FastConnect via virtual
circuits), and cross-region VCNs (remote peering connections). In a
hub-and-spoke topology the DRG acts as the hub.

This component bundles the DRG with its attachments, route tables,
route distributions, distribution statements, and static route rules
into a single deployment unit. Sub-resources reference each other by
display_name rather than OCID, making the YAML experience clean and
self-contained.

Excluded resources:
  - DrgAttachmentManagement: manages auto-created attachments for IPSec
    tunnels, virtual circuits, and remote peering connections. Those
    attachments are owned by the respective network resources.
  - DrgAttachmentsList: read-only data source.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.displayName` | `string` |  |  |  |
| `spec.attachments` | `[]DrgAttachment` |  |  |  |
| `spec.attachments[].displayName` | `string` | yes |  |  |
| `spec.attachments[].networkDetails` | `NetworkDetails` | yes |  |  |
| `spec.attachments[].networkDetails.type` | `enum` |  |  |  |
| `spec.attachments[].networkDetails.id` | `string \| valueFrom` | yes |  |  |
| `spec.attachments[].networkDetails.routeTableId` | `string` |  |  |  |
| `spec.attachments[].networkDetails.vcnRouteType` | `enum` |  |  |  |
| `spec.attachments[].drgRouteTableName` | `string` |  |  |  |
| `spec.attachments[].exportDrgRouteDistributionName` | `string` |  |  |  |
| `spec.routeTables` | `[]DrgRouteTable` |  |  |  |
| `spec.routeTables[].displayName` | `string` | yes |  |  |
| `spec.routeTables[].importDrgRouteDistributionName` | `string` |  |  |  |
| `spec.routeTables[].isEcmpEnabled` | `bool` |  |  |  |
| `spec.routeTables[].staticRouteRules` | `[]StaticRouteRule` |  |  |  |
| `spec.routeTables[].staticRouteRules[].destination` | `string` | yes |  |  |
| `spec.routeTables[].staticRouteRules[].nextHopAttachmentName` | `string` | yes |  |  |
| `spec.routeDistributions` | `[]DrgRouteDistribution` |  |  |  |
| `spec.routeDistributions[].displayName` | `string` | yes |  |  |
| `spec.routeDistributions[].distributionType` | `enum` |  |  |  |
| `spec.routeDistributions[].statements` | `[]DistributionStatement` |  |  |  |
| `spec.routeDistributions[].statements[].priority` | `int32` |  |  |  |
| `spec.routeDistributions[].statements[].matchCriteria` | `MatchCriteria` | yes |  |  |
| `spec.routeDistributions[].statements[].matchCriteria.matchType` | `enum` |  |  |  |
| `spec.routeDistributions[].statements[].matchCriteria.attachmentType` | `string` |  |  |  |
| `spec.routeDistributions[].statements[].matchCriteria.drgAttachmentName` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the DRG will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.displayName

`string`

Human-readable name for the DRG shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.attachments

`[]DrgAttachment`

Network attachments that connect VCNs or other network resources to this
DRG. Each attachment can optionally reference a custom route table and/or
a custom export route distribution defined within this component.

### spec.attachments[].displayName

`string` · required

Unique name for this attachment within the DRG component.
Used by other sub-resources (route rules, distribution statements)
to reference this attachment by name.

- rule: {"string":{"minLen":"1"}}

### spec.attachments[].networkDetails

`NetworkDetails` · required

Details of the network being attached.

- rule: {"required":true}

### spec.attachments[].networkDetails.type

`enum`

Type of network resource being attached.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `network_type_unspecified`
- `vcn`
- `ipsec_tunnel`
- `remote_peering_connection`
- `virtual_circuit`
- `loopback`

### spec.attachments[].networkDetails.id

`string | valueFrom` · required

OCID of the network resource (VCN, IPSec connection, virtual circuit,
or remote peering connection). For VCN attachments this is the VCN
OCID; for IPSec it is the IPSec connection OCID, etc.

Uses StringValueOrRef without a default_kind because the referenced
resource type depends on the attachment type.

- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.attachments[].networkDetails.routeTableId

`string`

OCID of the VCN route table to use for ingress routing (transit
routing). Only applicable for VCN attachments. When set, traffic
entering the VCN through this DRG attachment is routed according
to this VCN route table rather than the VCN's default route table.

### spec.attachments[].networkDetails.vcnRouteType

`enum`

Controls whether VCN CIDRs or individual subnet CIDRs are imported
into the DRG route table. Only applicable for VCN attachments.

Allowed values (use exactly as shown):

- `vcn_route_type_unspecified`
- `vcn_cidrs` -- Import the VCN's CIDR blocks as aggregate routes.
- `subnet_cidrs` -- Import individual subnet CIDR blocks for finer-grained routing.

### spec.attachments[].drgRouteTableName

`string`

Name of a DRG route table defined in this component's route_tables
list. When set, the attachment uses this custom route table instead
of the default route table for its network type.

### spec.attachments[].exportDrgRouteDistributionName

`string`

Name of a route distribution defined in this component's
route_distributions list. When set, the attachment uses this
distribution for exporting routes instead of the DRG's default
export distribution.

### spec.routeTables

`[]DrgRouteTable`

Custom DRG route tables for controlling traffic routing within the DRG.
OCI automatically creates default route tables per network type; these
are additional, user-defined tables. Route tables can import routes
from a distribution and contain static route rules.

### spec.routeTables[].displayName

`string` · required

Unique name for this route table within the DRG component.
Attachments reference route tables by this name.

- rule: {"string":{"minLen":"1"}}

### spec.routeTables[].importDrgRouteDistributionName

`string`

Name of a route distribution defined in this component's
route_distributions list. When set, routes from matching
attachments are automatically imported into this route table.

### spec.routeTables[].isEcmpEnabled

`bool`

When true, enables Equal-Cost Multi-Path (ECMP) routing across
multiple IPSec tunnels or virtual circuits. Traffic is distributed
across paths with equal-cost routes.

### spec.routeTables[].staticRouteRules

`[]StaticRouteRule`

Static route rules for this route table. Static routes take
precedence over dynamically imported routes.

### spec.routeTables[].staticRouteRules[].destination

`string` · required

Destination CIDR block (IPv4 or IPv6) for this route.
Example: "10.0.0.0/8", "192.168.1.0/24".

- rule: {"string":{"minLen":"1"}}

### spec.routeTables[].staticRouteRules[].nextHopAttachmentName

`string` · required

Name of the DRG attachment (defined in this component's attachments
list) that serves as the next hop for traffic matching this route.

- rule: {"string":{"minLen":"1"}}

### spec.routeDistributions

`[]DrgRouteDistribution`

Custom route distributions for controlling which routes are advertised
to DRG route tables (import) or to DRG attachments (export). OCI
automatically creates a default export distribution; these are
additional, user-defined distributions.

### spec.routeDistributions[].displayName

`string` · required

Unique name for this distribution within the DRG component.
Route tables and attachments reference distributions by this name.

- rule: {"string":{"minLen":"1"}}

### spec.routeDistributions[].distributionType

`enum`

Whether this distribution controls route import or export.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `distribution_type_unspecified`
- `import_routes` -- Import distribution: controls which routes are imported into a DRG route table from its associated attachments.
- `export_routes` -- Export distribution: controls which routes are exported from the DRG to an attachment's connected network.

### spec.routeDistributions[].statements

`[]DistributionStatement`

Prioritized statements that define which routes are accepted.
Statements are evaluated in priority order (lowest number first);
the first matching statement determines the action.

### spec.routeDistributions[].statements[].priority

`int32`

Priority of this statement (1-65535). Lower numbers are evaluated
first. Priorities must be unique within a distribution.

- rule: {"int32":{"lte":65535,"gte":1}}

### spec.routeDistributions[].statements[].matchCriteria

`MatchCriteria` · required

Criteria for matching routes. Determines which routes this
statement applies to.

- rule: {"required":true}

### spec.routeDistributions[].statements[].matchCriteria.matchType

`enum`

How to match routes for this statement.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `match_type_unspecified`
- `match_all` -- Match all routes regardless of attachment.
- `drg_attachment_type` -- Match routes from attachments of a specific network type.
- `drg_attachment_id` -- Match routes from a specific DRG attachment.

### spec.routeDistributions[].statements[].matchCriteria.attachmentType

`string`

Network attachment type to match against. Required when match_type
is drg_attachment_type. Accepted values: "VCN", "IPSEC_TUNNEL",
"VIRTUAL_CIRCUIT", "REMOTE_PEERING_CONNECTION".

### spec.routeDistributions[].statements[].matchCriteria.drgAttachmentName

`string`

Name of a DRG attachment (defined in this component's attachments
list) to match against. Required when match_type is
drg_attachment_id.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciDynamicRoutingGateway, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.drg_id` | `string` | OCID of the DRG. |
| `status.outputs.default_export_drg_route_distribution_id` | `string` | OCID of the default export route distribution that OCI automatically creates for the DRG. Useful for configuring external DRG attachments (e.g., IPSec tunnels managed outside this component). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
