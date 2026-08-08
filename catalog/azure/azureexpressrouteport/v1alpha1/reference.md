# AzureExpressRoutePort

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**Guide**: [GUIDE.md](../GUIDE.md) -- authored operational judgment for this component: conventions, trade-offs, and what pairs well with it.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureExpressRoutePort
metadata:
  name: test-express-route-port
spec:
  region: eastus
  resourceGroup:
    value: test-rg
  name: hq-port
  # The ExpressRoute DIRECT location vocabulary (narrower than circuit
  # peering locations): `az network express-route port location list`.
  peeringLocation: "Equinix-Ashburn-DC2"
  # 10 or 100 at real locations. The port bills its full monthly rate
  # from creation -- cross-connects or not.
  bandwidthInGbps: 10
  # DOT1Q: one VLAN tag per circuit (the common choice). Fixed at
  # creation, together with location and bandwidth.
  encapsulation: DOT1Q
  # Links always exist as a pair; these blocks manipulate them. They
  # start administratively disabled -- enable once the facility
  # completes the cross-connects.
  link1:
    adminEnabled: true
  link2:
    adminEnabled: true
  # One issued authorization: its ARM-generated key (sensitive)
  # surfaces in the authorization_keys output under this name.
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
| `spec.peeringLocation` | `string` | yes |  |  |
| `spec.bandwidthInGbps` | `int32` | yes |  |  |
| `spec.encapsulation` | `enum` |  |  |  |
| `spec.billingType` | `enum` |  | `METERED_DATA` |  |
| `spec.identity` | `AzureExpressRoutePortIdentity` |  |  |  |
| `spec.identity.type` | `enum` | yes |  |  |
| `spec.identity.identityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.link1` | `AzureExpressRoutePortLink` |  |  |  |
| `spec.link1.adminEnabled` | `bool` |  |  |  |
| `spec.link1.macsecCipher` | `enum` |  | `GCM_AES_128` |  |
| `spec.link1.macsecCknKeyvaultSecretId` | `string` |  |  |  |
| `spec.link1.macsecCakKeyvaultSecretId` | `string` |  |  |  |
| `spec.link1.macsecSciEnabled` | `bool` |  |  |  |
| `spec.link2` | `AzureExpressRoutePortLink` |  |  |  |
| `spec.link2.adminEnabled` | `bool` |  |  |  |
| `spec.link2.macsecCipher` | `enum` |  | `GCM_AES_128` |  |
| `spec.link2.macsecCknKeyvaultSecretId` | `string` |  |  |  |
| `spec.link2.macsecCakKeyvaultSecretId` | `string` |  |  |  |
| `spec.link2.macsecSciEnabled` | `bool` |  |  |  |
| `spec.authorizations` | `[]AzureExpressRoutePortAuthorization` |  |  |  |
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

- rule: {"required":true,"string":{"pattern":"^[^\\W_]([\\w.-]{0,78}[\\w])?$"}}

### spec.peeringLocation

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.bandwidthInGbps

`int32` · required

- rule: {"required":true,"int32":{"gte":1}}

### spec.encapsulation

`enum`

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_express_route_port_encapsulation_unspecified`
- `DOT1Q`
- `QINQ`

### spec.billingType

`enum` · optional (explicit presence)

- default: `METERED_DATA`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_express_route_port_billing_type_unspecified`
- `METERED_DATA`
- `UNLIMITED_DATA`

### spec.identity

`AzureExpressRoutePortIdentity`

- rule: identity_ids is required for USER_ASSIGNED and SYSTEM_AND_USER_ASSIGNED and must be empty for SYSTEM_ASSIGNED

### spec.identity.type

`enum` · required

- rule: {"required":true}

Allowed values (use exactly as shown):

- `azure_express_route_port_identity_type_unspecified`
- `SYSTEM_ASSIGNED`
- `USER_ASSIGNED`
- `SYSTEM_AND_USER_ASSIGNED`

### spec.identity.identityIds

`[]string | valueFrom`

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.link1

`AzureExpressRoutePortLink`

- rule: MACsec needs both keys or neither -- set macsec_ckn_keyvault_secret_id and macsec_cak_keyvault_secret_id together

### spec.link1.adminEnabled

`bool`

### spec.link1.macsecCipher

`enum` · optional (explicit presence)

- default: `GCM_AES_128`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_express_route_port_macsec_cipher_unspecified`
- `GCM_AES_128`
- `GCM_AES_256`
- `GCM_AES_XPN_128`
- `GCM_AES_XPN_256`

### spec.link1.macsecCknKeyvaultSecretId

`string`

### spec.link1.macsecCakKeyvaultSecretId

`string`

### spec.link1.macsecSciEnabled

`bool`

### spec.link2

`AzureExpressRoutePortLink`

- rule: MACsec needs both keys or neither -- set macsec_ckn_keyvault_secret_id and macsec_cak_keyvault_secret_id together

### spec.link2.adminEnabled

`bool`

### spec.link2.macsecCipher

`enum` · optional (explicit presence)

- default: `GCM_AES_128`
- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_express_route_port_macsec_cipher_unspecified`
- `GCM_AES_128`
- `GCM_AES_256`
- `GCM_AES_XPN_128`
- `GCM_AES_XPN_256`

### spec.link2.macsecCknKeyvaultSecretId

`string`

### spec.link2.macsecCakKeyvaultSecretId

`string`

### spec.link2.macsecSciEnabled

`bool`

### spec.authorizations

`[]AzureExpressRoutePortAuthorization`

### spec.authorizations[].name

`string` · required

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.tags

`map<string, string>`

## Validation Rules

- `encapsulation_required`: Choose the link encapsulation explicitly -- DOT1Q (one VLAN tag per circuit, the common choice) or QINQ (stacked tags for overlapping VLAN ranges)
- `macsec_requires_user_assigned_identity`: MACsec needs a USER_ASSIGNED (or SYSTEM_AND_USER_ASSIGNED) identity on the port -- the identity reads the CAK/CKN secrets from Key Vault
- `authorization_names_unique`: Authorization names must be unique on the port

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureExpressRoutePort, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.express_route_port_id` | `string` |  |
| `status.outputs.express_route_port_name` | `string` |  |
| `status.outputs.guid` | `string` |  |
| `status.outputs.ethertype` | `string` |  |
| `status.outputs.mtu` | `string` |  |
| `status.outputs.system_assigned_identity_principal_id` | `string` |  |
| `status.outputs.link1_id` | `string` |  |
| `status.outputs.link1_router_name` | `string` |  |
| `status.outputs.link1_interface_name` | `string` |  |
| `status.outputs.link1_patch_panel_id` | `string` |  |
| `status.outputs.link1_rack_id` | `string` |  |
| `status.outputs.link1_connector_type` | `string` |  |
| `status.outputs.link2_id` | `string` |  |
| `status.outputs.link2_router_name` | `string` |  |
| `status.outputs.link2_interface_name` | `string` |  |
| `status.outputs.link2_patch_panel_id` | `string` |  |
| `status.outputs.link2_rack_id` | `string` |  |
| `status.outputs.link2_connector_type` | `string` |  |
| `status.outputs.authorization_keys` | `map<string, string>` |  |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.identity.identityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## See Also

- [Overview](../README.md)
