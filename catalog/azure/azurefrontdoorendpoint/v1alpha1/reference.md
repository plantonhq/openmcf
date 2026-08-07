# AzureFrontDoorEndpoint

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureFrontDoorEndpointSpec** defines the configuration for creating an
endpoint inside an Azure Front Door (Standard/Premium) profile -- the
public entry point client traffic arrives at.

Each endpoint receives a generated, globally unique hostname of the form
`{name}-{hash}.z01.azurefd.net` (the hash suffix means the NAME only has
to be unique within the profile). Routes attach to an endpoint to define
which URL patterns it serves and which origin group answers them; custom
domains CNAME onto the endpoint's generated hostname.

Endpoints are many-per-profile with independent lifecycles -- one
profile commonly fronts several applications, each behind its own
endpoint -- which is why the endpoint is a first-class kind referencing
the profile rather than a list folded into the profile's spec.

**ForceNew fields**: `profile_id`, `endpoint_name` -- both fix the
endpoint's ARM identity at creation.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFrontDoorEndpoint
metadata:
  name: test-front-door-endpoint
spec:
  profileId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cdn/profiles/test-frontdoor
  endpointName: test-web
  # Exercises the disabled state (maintenance switch).
  enabled: false
  tags:
    cost-center: platform
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.profileId` | `string \| valueFrom` | yes |  | AzureFrontDoorProfile (`status.outputs.profile_id`) |
| `spec.endpointName` | `string` | yes |  |  |
| `spec.enabled` | `bool` |  | `true` |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.profileId

`string | valueFrom` · required

The Front Door profile the endpoint lives in, by ARM ID. References
an AzureFrontDoorProfile's profile_id output so the profile and its
endpoints compose in one manifest set. Fixed at creation.

- references: AzureFrontDoorProfile (`status.outputs.profile_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorProfile, name: <that resource's name>, fieldPath: status.outputs.profile_id}} -- a bare string does not parse

### spec.endpointName

`string` · required

The endpoint's name -- unique within the profile, and the prefix of
the generated public hostname (`{name}-{hash}.z01.azurefd.net`), so
pick something recognizable in browser address bars and DNS records.

2-46 characters; letters, digits, and hyphens; must start and end
with a letter or digit.

**ForceNew**: changing the name replaces the endpoint and changes its
generated hostname -- every DNS record pointing at the old hostname
breaks.

- rule: endpoint_name must be 2-46 characters, start and end with a letter or digit, and contain only letters, digits, and hyphens
- rule: {"required":true,"string":{"minLen":"2","maxLen":"46"}}

### spec.enabled

`bool` · optional (explicit presence)

Whether the endpoint accepts traffic. Disabling stops all requests at
the edge (returns errors to clients) without deleting the endpoint or
its routes -- useful for maintenance windows and staged cutovers.
Default true.

- default: `true`

### spec.tags

`map<string, string>`

Free-form tags applied to the endpoint, merged over the
Planton-derived resource tags (organization, environment, resource
id); a user tag with the same key wins. Tags are Azure's governance
surface -- Azure Policy enforces them and Microsoft Cost Management
groups by them. Updatable in place.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFrontDoorEndpoint, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.endpoint_id` | `string` | The Azure Resource Manager ID of the endpoint -- what AzureFrontDoorRoute's endpoint_id references, and (alongside custom domain IDs) what a Front Door security policy associates a WAF policy with. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cdn/profiles/{profile}/afdEndpoints/{name} |
| `status.outputs.endpoint_name` | `string` | The endpoint's name -- unique within its profile. |
| `status.outputs.host_name` | `string` | The generated, globally unique public hostname clients connect to (`{name}-{hash}.z01.azurefd.net`). This is the CNAME target for custom-domain DNS records (e.g. an AzureDnsRecord pointing `www.example.com` at this value). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.profileId` | AzureFrontDoorProfile | `status.outputs.profile_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureFrontDoorRoute | `spec.endpointId` | `status.outputs.endpoint_id` |
| AzureFrontDoorSecurityPolicy | `spec.domainIds` | `status.outputs.endpoint_id` |

## See Also

- [Overview](../README.md)
