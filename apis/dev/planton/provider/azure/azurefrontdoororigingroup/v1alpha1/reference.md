# AzureFrontDoorOriginGroup

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `azure.planton.dev/v1alpha1`

**AzureFrontDoorOriginGroupSpec** defines the configuration for creating
an origin group inside an Azure Front Door (Standard/Premium) profile --
the load-balanced pool of backends a route sends traffic to.

The group carries the pool-level behavior: how Front Door samples origin
health, how it weighs latency when picking an origin, whether clients
stick to one origin (session affinity), and how traffic shifts back to
an origin that just recovered. The backends themselves are first-class
AzureFrontDoorOrigin resources referencing this group -- so a regional
stamp can add its backend to a shared group without touching the group
or any other region's origins.

**ForceNew fields**: `profile_id`, `origin_group_name` -- both fix the
group's ARM identity at creation.

## Example

```yaml
apiVersion: azure.planton.dev/v1alpha1
kind: AzureFrontDoorOriginGroup
metadata:
  name: test-front-door-origin-group
spec:
  profileId:
    value: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/test-rg/providers/Microsoft.Cdn/profiles/test-frontdoor
  originGroupName: api-backends
  # Exercises explicit load-balancing dials.
  loadBalancing:
    sampleSize: 6
    successfulSamplesRequired: 4
    additionalLatencyInMilliseconds: 200
  # Exercises the health probe (absence would disable probing).
  healthProbe:
    protocol: HTTPS
    intervalInSeconds: 30
    requestType: GET
    path: /api/healthz
  # Exercises stateless-API affinity off + a fast traffic-restore ramp.
  sessionAffinityEnabled: false
  restoreTrafficTimeToHealedOrNewEndpointInMinutes: 5
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.profileId` | `string \| valueFrom` | yes |  | AzureFrontDoorProfile (`status.outputs.profile_id`) |
| `spec.originGroupName` | `string` | yes |  |  |
| `spec.loadBalancing` | `AzureFrontDoorOriginGroupLoadBalancing` |  |  |  |
| `spec.loadBalancing.sampleSize` | `int32` |  | `4` |  |
| `spec.loadBalancing.successfulSamplesRequired` | `int32` |  | `3` |  |
| `spec.loadBalancing.additionalLatencyInMilliseconds` | `int32` |  | `50` |  |
| `spec.healthProbe` | `AzureFrontDoorOriginGroupHealthProbe` |  |  |  |
| `spec.healthProbe.protocol` | `enum` |  |  |  |
| `spec.healthProbe.intervalInSeconds` | `int32` | yes |  |  |
| `spec.healthProbe.requestType` | `enum` |  |  |  |
| `spec.healthProbe.path` | `string` |  | `/` |  |
| `spec.sessionAffinityEnabled` | `bool` |  | `true` |  |
| `spec.restoreTrafficTimeToHealedOrNewEndpointInMinutes` | `int32` |  | `10` |  |

## Field Details

### spec.profileId

`string | valueFrom` · required

The Front Door profile the origin group lives in, by ARM ID.
References an AzureFrontDoorProfile's profile_id output so the
profile and its origin groups compose in one manifest set. Fixed at
creation.

- references: AzureFrontDoorProfile (`status.outputs.profile_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: AzureFrontDoorProfile, name: <that resource's name>, fieldPath: status.outputs.profile_id}} -- a bare string does not parse

### spec.originGroupName

`string` · required

The origin group's name -- unique within the profile. Routes and
origins reference the group by its ARM ID (see the origin_group_id
output), so the name is mostly a human-facing label in the portal.

2-90 characters; letters, digits, and hyphens; must start and end
with a letter or digit.

**ForceNew**: changing the name replaces the group AND every origin
nested under it.

- rule: origin_group_name must be 2-90 characters, start and end with a letter or digit, and contain only letters, digits, and hyphens
- rule: {"required":true,"string":{"minLen":"2","maxLen":"90"}}

### spec.loadBalancing

`AzureFrontDoorOriginGroupLoadBalancing`

How Front Door distributes traffic across the group's origins.
Azure requires load-balancing settings on every origin group; leave
this unset to deploy Azure's defaults (sample size 4, 3 successful
samples required, 50 ms additional latency).

### spec.loadBalancing.sampleSize

`int32` · optional (explicit presence)

How many recent health-probe samples to consider per origin, 0-255.
Default 4 (Azure's default).

- default: `4`
- rule: {"int32":{"lte":255,"gte":0}}

### spec.loadBalancing.successfulSamplesRequired

`int32` · optional (explicit presence)

How many of those samples must have succeeded for the origin to
count as healthy, 0-255. Default 3 (Azure's default). Lower it
relative to sample_size to tolerate flaky probes; raise it to eject
an origin faster.

- default: `3`
- rule: {"int32":{"lte":255,"gte":0}}

### spec.loadBalancing.additionalLatencyInMilliseconds

`int32` · optional (explicit presence)

The latency window in milliseconds, 0-1000: origins within this many
milliseconds of the fastest origin are treated as equally fast and
share traffic by weight. Small values pin traffic to the closest
origin; larger values spread it across geographically dispersed
backends. Default 50 (Azure's default).

- default: `50`
- rule: {"int32":{"lte":1000,"gte":0}}

### spec.healthProbe

`AzureFrontDoorOriginGroupHealthProbe`

Periodic health probing of the group's origins. When set, Front Door
probes each origin at the configured interval and takes unhealthy
origins out of rotation until they recover; when unset, probing is
disabled and every origin is assumed healthy. Skip probes for
single-origin groups (they only add origin load -- traffic has
nowhere else to go); configure them whenever the group has two or
more origins.

### spec.healthProbe.protocol

`enum`

The protocol probes are sent over. Probe over HTTPS when origins
serve TLS (the common case) so the probe exercises the same path as
real traffic.

- rule: {"enum":{"definedOnly":true,"notIn":[0]}}

Allowed values (use exactly as shown):

- `azure_front_door_origin_group_health_probe_protocol_unspecified` -- Not specified -- invalid; probes must declare Http or Https.
- `HTTP` -- Probe over plain HTTP (the origin's http_port).
- `HTTPS` -- Probe over TLS (the origin's https_port) -- exercises certificate validity as part of health.

### spec.healthProbe.intervalInSeconds

`int32` · required

Seconds between probes to each origin, 1-255. Shorter intervals
detect failures faster but multiply probe load across Front Door's
edge fleet -- every edge location probes every origin, so an
aggressive interval can be a meaningful fraction of a small origin's
traffic. 30-120 s suits most production workloads.

- rule: {"required":true,"int32":{"lte":255,"gte":1}}

### spec.healthProbe.requestType

`enum`

The HTTP method probes use. HEAD (the default) is cheapest and
right for most origins; switch to GET when the health endpoint only
answers GET or the origin needs to compute a body to prove health.

- rule: {"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `azure_front_door_origin_group_health_probe_request_type_unspecified` -- Not specified -- deploys HEAD, Azure's default.
- `HEAD` -- HEAD request -- no response body; cheapest probe.
- `GET` -- GET request -- use when the health endpoint needs to produce a body or only answers GET.

### spec.healthProbe.path

`string` · optional (explicit presence)

The URL path probed on each origin. Default "/". Point it at a
dedicated health endpoint (e.g. "/healthz") so probes exercise the
application, not just the web server.

- default: `/`
- rule: path must start with '/'

### spec.sessionAffinityEnabled

`bool` · optional (explicit presence)

Route subsequent requests from the same client to the same origin
(cookie-based affinity). Keep enabled for stateful backends
(in-memory sessions); disable for stateless APIs so traffic spreads
evenly. Default true (Azure's default).

- default: `true`

### spec.restoreTrafficTimeToHealedOrNewEndpointInMinutes

`int32` · optional (explicit presence)

How many minutes to gradually shift traffic to an origin that just
became healthy (or was just added), 0-50. A ramp avoids
cold-starting a recovered backend with its full traffic share.
0 shifts immediately. Default 10 (Azure's default).

- default: `10`
- rule: {"int32":{"lte":50,"gte":0}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AzureFrontDoorOriginGroup, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.origin_group_id` | `string` | The Azure Resource Manager ID of the origin group -- what AzureFrontDoorOrigin's origin_group_id (parent) and AzureFrontDoorRoute's origin_group_id (destination) reference. Format: /subscriptions/{sub}/resourceGroups/{rg}/providers/Microsoft.Cdn/profiles/{profile}/originGroups/{name} |
| `status.outputs.origin_group_name` | `string` | The origin group's name -- unique within its profile. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.profileId` | AzureFrontDoorProfile | `status.outputs.profile_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| AzureFrontDoorOrigin | `spec.originGroupId` | `status.outputs.origin_group_id` |
| AzureFrontDoorRoute | `spec.originGroupId` | `status.outputs.origin_group_id` |
| AzureFrontDoorRuleSet | `spec.rules[].actions.routeConfigurationOverride.originGroupId` | `status.outputs.origin_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
