# AzureExpressRouteCircuit

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureExpressRouteCircuit
metadata:
  name: test-express-route-circuit
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: hq-circuit
  # STANDARD reaches the geopolitical area; METERED_DATA bills outbound
  # per GB (the common production pairing).
  skuTier: STANDARD
  skuFamily: METERED_DATA
  # Service-provider mode: the trio travels together. Billing starts at
  # creation; the circuit sits NotProvisioned until the provider
  # completes the cross-connect with the service key.
  serviceProviderName: "Equinix"
  peeringLocation: "Washington DC"
  bandwidthInMbps: 50
  # One issued authorization: its ARM-generated key (sensitive) surfaces
  # in the authorization_keys output under this name.
  authorizations:
    - name: partner-team
  tags:
    purpose: hack-test
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.name` | `string` | yes |  |  |
| `spec.skuTier` | `enum` |  |  |  |
| `spec.skuFamily` | `enum` |  |  |  |
| `spec.serviceProviderName` | `string` |  |  |  |
| `spec.peeringLocation` | `string` |  |  |  |
| `spec.bandwidthInMbps` | `int32` |  |  |  |
| `spec.expressRoutePortId` | `string` |  |  |  |
| `spec.bandwidthInGbps` | `double` |  |  |  |
| `spec.rateLimitingEnabled` | `bool` |  |  |  |
| `spec.allowClassicOperations` | `bool` |  |  |  |
| `spec.authorizationKey` | `string \| valueFrom` (sensitive) |  |  |  |
| `spec.authorizations` | `[]AzureExpressRouteCircuitAuthorization` |  |  |  |
| `spec.authorizations[].name` | `string` | yes |  |  |
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

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.skuTier

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_express_route_circuit_sku_tier_unspecified`
- `BASIC`
- `LOCAL`
- `STANDARD`
- `PREMIUM`

### spec.skuFamily

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_express_route_circuit_sku_family_unspecified`
- `METERED_DATA`
- `UNLIMITED_DATA`

### spec.serviceProviderName

`string`

### spec.peeringLocation

`string`

### spec.bandwidthInMbps

`int32`

- rule: {"int32":{"gte":0}}

### spec.expressRoutePortId

`string`

### spec.bandwidthInGbps

`double`

- rule: {"double":{"gte":0}}

### spec.rateLimitingEnabled

`bool`

### spec.allowClassicOperations

`bool`

### spec.authorizationKey

`string | valueFrom` · sensitive

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.authorizations

`[]AzureExpressRouteCircuitAuthorization`

### spec.authorizations[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.tags

`map<string, string>`

## Validation Rules

- `sku_tier_required`: Choose the SKU tier explicitly -- LOCAL (metro only, no egress fees), STANDARD (geopolitical area), or PREMIUM (global reach)
- `sku_family_required`: Choose the billing family explicitly -- METERED_DATA (pay per outbound GB) or UNLIMITED_DATA (flat rate)
- `exactly_one_provisioning_mode`: Provision the circuit exactly one way: through a service provider (service_provider_name + peering_location + bandwidth_in_mbps) or on an ExpressRoute Direct port (express_route_port_id + bandwidth_in_gbps)
- `service_provider_trio_travels_together`: Service-provider mode needs all three of service_provider_name, peering_location, and bandwidth_in_mbps
- `provider_fields_only_in_provider_mode`: peering_location and bandwidth_in_mbps belong to service-provider mode -- remove them from an ExpressRoute Direct circuit
- `direct_pair_travels_together`: ExpressRoute Direct mode needs both express_route_port_id and bandwidth_in_gbps
- `gbps_only_in_direct_mode`: bandwidth_in_gbps belongs to ExpressRoute Direct mode -- service-provider circuits size bandwidth_in_mbps instead
- `authorization_names_unique`: Authorization names must be unique on the circuit

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureExpressRouteCircuit, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.express_route_circuit_id` | `string` |  |
| `status.outputs.express_route_circuit_name` | `string` |  |
| `status.outputs.service_key` | `string` |  |
| `status.outputs.service_provider_provisioning_state` | `string` |  |
| `status.outputs.authorization_keys` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureExpressRouteCircuitPeering | `spec.expressRouteCircuitName` | `status.outputs.express_route_circuit_name` |

## See Also

- [Overview](../README.md)
