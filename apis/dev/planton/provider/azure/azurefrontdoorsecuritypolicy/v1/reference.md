# AzureFrontDoorSecurityPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1`

**AzureFrontDoorSecurityPolicySpec** defines the configuration for
creating a security policy inside an Azure Front Door
(Standard/Premium) profile -- the association that attaches a Front
Door WAF policy (AzureFrontDoorFirewallPolicy) to the hostnames the
profile serves. A WAF policy enforces nothing until a security policy
associates it; this kind is the enforcement seam.

The `domain_ids` list names the hostnames the WAF protects. Each
entry is either an AzureFrontDoorEndpoint's ARM ID (protect the
endpoint's generated *.azurefd.net hostname) or an
AzureFrontDoorCustomDomain's ARM ID (protect that custom hostname) --
Azure accepts both resource types interchangeably here.

**Path scope**: the WAF applies to ALL paths (`/*`) on every
associated domain -- Azure's security policies accept no other
pattern, so it is not configurable. Scope enforcement by choosing
WHICH domains to associate, not which paths.

**Domain caps ride the profile's sku**: a STANDARD profile allows up
to 100 domains per association, PREMIUM up to 500 (checked at deploy
time -- the cap lives on the profile, not this resource). The WAF
policy's own sku must also MATCH the profile's sku; Azure rejects a
mismatched pairing.

**ForceNew fields**: `profile_id`, `security_policy_name`, and
`firewall_policy_id` -- swapping the WAF policy replaces the security
policy (a fast, metadata-only operation). The domain list updates in
place.

## Example

```yaml
apiVersion: azure.planton.dev/v1
kind: AzureFrontDoorSecurityPolicy
metadata:
  name: test-front-door-security-policy
spec:
  profileId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cdn/profiles/test-frontdoor
  securityPolicyName: edge-waf-attach
  firewallPolicyId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Network/frontDoorWebApplicationFirewallPolicies/testedgewaf
  # A mixed list: the endpoint's default *.azurefd.net hostname plus a
  # custom domain -- exercises both accepted ID shapes.
  domainIds:
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cdn/profiles/test-frontdoor/afdEndpoints/web
    - value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cdn/profiles/test-frontdoor/customDomains/www-example-com
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.profileId` | `string \| valueFrom` | yes |  | AzureFrontDoorProfile (`status.outputs.profile_id`) |
| `spec.securityPolicyName` | `string` | yes |  |  |
| `spec.firewallPolicyId` | `string \| valueFrom` | yes |  | AzureFrontDoorFirewallPolicy (`status.outputs.firewall_policy_id`) |
| `spec.domainIds` | `[]string \| valueFrom` | yes |  | AzureFrontDoorEndpoint (`status.outputs.endpoint_id`) |

## Field Details

### spec.profileId

`string | valueFrom` · required

The Front Door profile the security policy lives in, by ARM ID.
References an AzureFrontDoorProfile's profile_id output so the
profile and its security policies compose in one manifest set.
Fixed at creation.

- references: AzureFrontDoorProfile (`status.outputs.profile_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorProfile, name: <that resource's name>, fieldPath: status.outputs.profile_id}} -- a bare string does not parse

### spec.securityPolicyName

`string` · required

The security policy's name -- unique within the profile. Must
begin and end with a letter or digit and may contain only letters,
digits, and hyphens.

**ForceNew**: changing the name replaces the security policy (a
brief enforcement gap while the old association tears down).

- rule: security_policy_name must begin and end with a letter or digit and may contain only letters, digits, and hyphens
- rule: {"required":true}

### spec.firewallPolicyId

`string | valueFrom` · required

The Front Door WAF policy to enforce, by ARM ID. References an
AzureFrontDoorFirewallPolicy's firewall_policy_id output. The WAF
policy's sku must match the profile's sku (Azure rejects the
pairing otherwise, at deploy time).

**ForceNew**: pointing at a different WAF policy replaces the
security policy.

- references: AzureFrontDoorFirewallPolicy (`status.outputs.firewall_policy_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorFirewallPolicy, name: <that resource's name>, fieldPath: status.outputs.firewall_policy_id}} -- a bare string does not parse

### spec.domainIds

`[]string | valueFrom` · required

The hostnames the WAF protects, each by ARM ID: an
AzureFrontDoorEndpoint's endpoint_id (the generated *.azurefd.net
hostname) or an AzureFrontDoorCustomDomain's custom_domain_id --
Azure accepts both types here, and a list can mix them. At least
one; at most 500 (a STANDARD profile caps the list at 100 -- the
cap rides the profile's sku and is checked at deploy time).
Updatable in place: adding a domain extends protection without
touching the others.

- references: AzureFrontDoorEndpoint (`status.outputs.endpoint_id`)
- rule: {"repeated":{"minItems":"1","maxItems":"500"}}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorEndpoint, name: <that resource's name>, fieldPath: status.outputs.endpoint_id}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFrontDoorSecurityPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.security_policy_id` | `string` | The Azure Resource Manager ID of the security policy. Nothing composes on a security policy (it is itself the association), so this exists for operational addressing -- diagnostics, RBAC scoping, and ARM reads. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cdn/profiles/{profile}/securityPolicies/{name} |
| `status.outputs.security_policy_name` | `string` | The security policy's name -- unique within its profile. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.profileId` | AzureFrontDoorProfile | `status.outputs.profile_id` |
| `spec.firewallPolicyId` | AzureFrontDoorFirewallPolicy | `status.outputs.firewall_policy_id` |
| `spec.domainIds` | AzureFrontDoorEndpoint | `status.outputs.endpoint_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
