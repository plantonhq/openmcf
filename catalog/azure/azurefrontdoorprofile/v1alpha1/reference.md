# AzureFrontDoorProfile

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureFrontDoorProfileSpec** defines the configuration for creating an
Azure Front Door (Standard/Premium) profile -- the top-level container
for a global content-delivery and application-acceleration deployment on
Microsoft's edge network.

The profile is deliberately just the container: it owns the SKU tier,
the origin response timeout, the managed identity, access-log scrubbing,
and tags. The delivery surface composes from first-class resources that
reference this profile:
- AzureFrontDoorEndpoint -- the public entry hostname (*.azurefd.net)
- AzureFrontDoorOriginGroup -- a load-balanced backend pool
- AzureFrontDoorOrigin -- one backend inside an origin group
- AzureFrontDoorRoute -- connects an endpoint to an origin group by
  URL pattern

This mirrors Azure's own resource model (each is a separate ARM child
resource with an independent lifecycle) and keeps regional stamps
composable: a new region can add its origin to a shared origin group
without touching the profile or any other region's resources.

**SKU tiers**:
- STANDARD: global HTTP load balancing, SSL offloading, caching,
  compression, URL routing. 99.99% SLA.
- PREMIUM: everything in STANDARD plus Private Link to origins,
  managed WAF rule sets, and Bot Manager.

**Azure Front Door is a global resource** -- it has no region; Azure
deploys it across all edge locations worldwide. It still lives inside a
resource group for ARM organization.

**ForceNew fields**: `profile_name`, `sku`. Additionally, Azure rejects
a PREMIUM -> STANDARD change outright (upgrades recreate; downgrades are
not supported at all).

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFrontDoorProfile
metadata:
  name: test-front-door-profile
spec:
  resourceGroup:
    value: test-rg
  profileName: test-frontdoor
  # Exercises the Premium tier (Private Link origins and managed WAF
  # rules are Premium-only; the sku is ForceNew and cannot downgrade).
  sku: PREMIUM
  responseTimeoutSeconds: 60
  # Exercises the identity seam: a system-assigned identity for keyless
  # Key Vault certificate access.
  identity:
    type: SYSTEM_ASSIGNED
  # Exercises log scrubbing (presence enables it).
  logScrubbingVariables:
    - REQUEST_IP_ADDRESS
    - REQUEST_URI
  tags:
    cost-center: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.resourceGroup` | `string \| valueFrom` | yes |  | AzureResourceGroup (`status.outputs.resource_group_name`) |
| `spec.profileName` | `string` | yes |  |  |
| `spec.sku` | `enum` |  |  |  |
| `spec.responseTimeoutSeconds` | `int32` |  | `120` |  |
| `spec.identity` | `AzureFrontDoorProfileIdentity` |  |  |  |
| `spec.identity.type` | `enum` |  |  |  |
| `spec.identity.userAssignedIdentityIds` | `[]string \| valueFrom` |  |  | AzureUserAssignedIdentity (`status.outputs.identity_id`) |
| `spec.logScrubbingVariables` | `[]enum` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.resourceGroup

`string | valueFrom` · required

The Azure Resource Group the Front Door profile is created in.
Front Door is global (no region), but every ARM resource belongs to a
resource group for organization, RBAC scoping, and lifecycle grouping.

- references: AzureResourceGroup (`status.outputs.resource_group_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureResourceGroup, name: <that resource's name>, fieldPath: status.outputs.resource_group_name}} -- a bare string does not parse

### spec.profileName

`string` · required

The profile's name -- unique within the resource group. This is the
ARM identity every child resource (endpoint, origin group, rule set,
custom domain, secret, security policy) is nested under.

2-90 characters; letters, digits, and hyphens; must start and end
with a letter or digit.

**ForceNew**: changing the name replaces the profile AND everything
nested under it.

- rule: profile_name must be 2-90 characters, start and end with a letter or digit, and contain only letters, digits, and hyphens
- rule: {"required":true,"string":{"minLen":"2","maxLen":"90"}}

### spec.sku

`enum`

The pricing/capability tier. Unspecified deploys STANDARD -- the
right answer unless you need Private Link to origins or the managed
WAF rule sets, which are PREMIUM-only.

**ForceNew** -- and Azure additionally refuses a PREMIUM -> STANDARD
downgrade even as a replace, so choose PREMIUM deliberately.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_front_door_profile_sku_unspecified` -- Not specified -- deploys STANDARD, the production default.
- `STANDARD` -- Global load balancing, SSL offloading, caching, compression, and URL-based routing. No Private Link to origins, no managed WAF rules.
- `PREMIUM` -- Everything in STANDARD plus Private Link to origins, managed WAF rule sets (Microsoft_DefaultRuleSet, Bot Manager), and JS challenge/CAPTCHA custom-rule actions.

### spec.responseTimeoutSeconds

`int32` · optional (explicit presence)

How long Front Door waits for the origin's response before returning
a 504 to the client, in seconds (16-240, default 120). Raise it for
slow APIs and large downloads; lower it when fast failover matters
more than slow origins completing.

- default: `120`
- rule: {"int32":{"lte":240,"gte":16}}

### spec.identity

`AzureFrontDoorProfileIdentity`

The profile's managed identity. Front Door uses it to read
customer-managed TLS certificates from Key Vault (via
AzureFrontDoorSecret) without an access-policy secret -- assign one
when custom domains will carry bring-your-own certificates.

- rule: user_assigned_identity_ids is required with USER_ASSIGNED or SYSTEM_AND_USER_ASSIGNED, and must be empty with SYSTEM_ASSIGNED

### spec.identity.type

`enum`

The identity model: SYSTEM_ASSIGNED (Azure creates and rotates a
service principal bound to the profile's lifecycle), USER_ASSIGNED
(bring identities from user_assigned_identity_ids, shareable across
resources), or SYSTEM_AND_USER_ASSIGNED (both).

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_profile_identity_type_unspecified` -- Not specified -- invalid; choose an explicit identity model.
- `SYSTEM_ASSIGNED` -- Azure creates a service principal bound to the profile's lifecycle.
- `USER_ASSIGNED` -- Bring your own AzureUserAssignedIdentity entries -- shareable across resources and grantable before the profile exists.
- `SYSTEM_AND_USER_ASSIGNED` -- Both a system-assigned principal and user-assigned identities.

### spec.identity.userAssignedIdentityIds

`[]string | valueFrom`

The user-assigned identities to attach -- required when (and only
meaningful when) type includes USER_ASSIGNED. Each entry references
an AzureUserAssignedIdentity's ARM id.

- references: AzureUserAssignedIdentity (`status.outputs.identity_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureUserAssignedIdentity, name: <that resource's name>, fieldPath: status.outputs.identity_id}} -- a bare string does not parse

### spec.logScrubbingVariables

`[]enum`

Scrub (mask) sensitive request data out of Front Door's access logs
before they are written. Each entry names one request part to scrub;
listing at least one entry enables scrubbing, an empty list leaves it
disabled. Azure scrubs ALL values of the selected part (the service
supports only the match-everything operator on profiles).

- rule: {"repeated":{"maxItems":"3","unique":true,"items":{"enum":{"definedOnly":true,"notIn":[0]}}}}

Allowed values (use exactly as shown):

- `azure_front_door_profile_log_scrubbing_variable_unspecified` -- Not specified -- invalid; name the request part to scrub.
- `QUERY_STRING_ARG_NAMES` -- Mask every query-string argument name/value in logged request URLs.
- `REQUEST_IP_ADDRESS` -- Mask the client IP address in access-log entries.
- `REQUEST_URI` -- Mask the request URI (path and query) in access-log entries.

### spec.tags

`map<string, string>`

Free-form tags applied to the profile, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins. Tags are Azure's governance
surface -- Azure Policy enforces them and Microsoft Cost Management
groups by them. Updatable in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFrontDoorProfile, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.profile_id` | `string` | The Azure Resource Manager ID of the Front Door profile -- what AzureFrontDoorEndpoint and AzureFrontDoorOriginGroup reference as their parent. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cdn/profiles/{name} |
| `status.outputs.profile_name` | `string` | The profile's name -- the namespace every child resource is nested under in ARM. |
| `status.outputs.resource_guid` | `string` | The Front Door service's own GUID for this profile (distinct from the ARM resource ID). Azure asks for it when validating traffic ownership -- e.g. the apex-domain "afdverify" DNS record. |
| `status.outputs.identity_principal_id` | `string` | The object (principal) ID of the profile's system-assigned managed identity -- the principal to grant Key Vault access to for bring-your-own TLS certificates. Empty when the identity block is absent or user-assigned only. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.resourceGroup` | AzureResourceGroup | `status.outputs.resource_group_name` |
| `spec.identity.userAssignedIdentityIds` | AzureUserAssignedIdentity | `status.outputs.identity_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureFrontDoorCustomDomain | `spec.profileId` | `status.outputs.profile_id` |
| AzureFrontDoorEndpoint | `spec.profileId` | `status.outputs.profile_id` |
| AzureFrontDoorOriginGroup | `spec.profileId` | `status.outputs.profile_id` |
| AzureFrontDoorRuleSet | `spec.profileId` | `status.outputs.profile_id` |
| AzureFrontDoorSecret | `spec.profileId` | `status.outputs.profile_id` |
| AzureFrontDoorSecurityPolicy | `spec.profileId` | `status.outputs.profile_id` |
| AzureFunctionApp | `spec.siteConfig.ipRestrictions[].headers.xAzureFdid` | `status.outputs.resource_guid` |
| AzureFunctionApp | `spec.siteConfig.scmIpRestrictions[].headers.xAzureFdid` | `status.outputs.resource_guid` |
| AzureFunctionAppFlexConsumption | `spec.siteConfig.ipRestrictions[].headers.xAzureFdid` | `status.outputs.resource_guid` |
| AzureFunctionAppFlexConsumption | `spec.siteConfig.scmIpRestrictions[].headers.xAzureFdid` | `status.outputs.resource_guid` |
| AzureLinuxWebApp | `spec.siteConfig.ipRestrictions[].headers.xAzureFdid` | `status.outputs.resource_guid` |
| AzureLinuxWebApp | `spec.siteConfig.scmIpRestrictions[].headers.xAzureFdid` | `status.outputs.resource_guid` |

## See Also

- [Overview](../README.md)
