# DigitalOceanDroplet

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1alpha1`

DigitalOceanDropletSpec defines the user configuration for a DigitalOcean Droplet (VM).

## Example

```yaml
apiVersion: digital-ocean.planton.dev/v1alpha1
kind: DigitalOceanDroplet
metadata:
  name: first-droplet                        # Kubernetes object name
spec:
  dropletName: first-droplet                         # droplet hostname
  region: blr1                             # NYC3 | SFO3 | FRA1 etc.
  size: s-2vcpu-4gb                        # enum value
  image: ubuntu-22-04-x64                  # official Ubuntu image slug
  vpc:
    value: b5648f9e-a28a-4760-bb87-b2fad07ae295                           # UUID or ref to DigitalOceanVpc
  enableIpv6: false                         # optional
  enableBackups: false                     # optional
  disableMonitoring: false                 # keep DO monitoring agent
#  volumeIds:                               # optional block‑volume attachment(s)
#    - value: 93a7a5b4-62ce-11f0-b9db-0a58ac1466b2                      # UUID or ref to DigitalOceanVolume
  tags:                                    # unique,  ≤ .64 .chars each
    - planton
  userData: |                              # cloud‑init (≤32 .KiB)
    #cloud-config
    package_update: true
    runcmd:
      - apt-get install -y nginx
  timezone: utc                            # utc (default) | local
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.dropletName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.size` | `string` | yes |  |  |
| `spec.image` | `string` | yes |  |  |
| `spec.vpc` | `string \| valueFrom` | yes |  | DigitalOceanVpc (`status.outputs.vpc_id`) |
| `spec.enableIpv6` | `bool` |  |  |  |
| `spec.enableBackups` | `bool` |  |  |  |
| `spec.disableMonitoring` | `bool` |  |  |  |
| `spec.volumeIds` | `[]string \| valueFrom` |  |  | DigitalOceanVolume (`status.outputs.volume_id`) |
| `spec.tags` | `[]string` |  |  |  |
| `spec.userData` | `string` |  |  |  |
| `spec.timezone` | `enum` |  | `UTC` |  |

## Field Details

### spec.dropletName

`string` · required

droplet hostname (DNS-compatible, <=63 chars)

- rule: {"required":true,"string":{"maxLen":"63","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"}}

### spec.region

`enum` · required

region slug (datacenter location for the droplet)

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_region_unspecified` -- 0: default / unspecified region
- `nyc3` -- new york 3
- `sfo3` -- san francisco 3
- `fra1` -- frankfurt 1
- `sgp1` -- singapore 1
- `lon1` -- london 1
- `tor1` -- toronto 1
- `blr1` -- bangalore 1
- `ams3` -- amsterdam 3

### spec.size

`string` · required

Droplet size slug, e.g. "s-2vcpu-4gb" or "g-8vcpu-32gb".
Valid values: must match the regexp "^[a-z0-9]+(-[a-z0-9]+)+$" and
must be accepted by the DigitalOcean /v2/sizes API at creation time.

- rule: {"required":true,"string":{"pattern":"^[a-z0-9]+(-[a-z0-9]+)+$"}}

### spec.image

`string` · required

image slug for the droplet base image (e.g. "ubuntu-22-04-x64")

- rule: {"required":true,"string":{"pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"}}

### spec.vpc

`string | valueFrom` · required

target vpc network uuid for the droplet

- references: DigitalOceanVpc (`status.outputs.vpc_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.enableIpv6

`bool`

enable IPv6 networking (disabled by default)

### spec.enableBackups

`bool`

enable automated backups (disabled by default)

### spec.disableMonitoring

`bool`

disable digitalocean monitoring agent (monitoring on by default)

### spec.volumeIds

`[]string | valueFrom`

block storage volumes to attach (must reside in same region)

- references: DigitalOceanVolume (`status.outputs.volume_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanVolume, name: <that resource's name>, fieldPath: status.outputs.volume_id}} -- a bare string does not parse

### spec.tags

`[]string`

tags to apply to the droplet (must be unique)

- rule: {"repeated":{"unique":true}}

### spec.userData

`string`

cloud-init user data script (<=32 KiB)

- rule: {"string":{"maxBytes":"32768"}}

### spec.timezone

`enum` · optional (explicit presence)

timezone setting for the droplet's clock (default: UTC)

- default: `UTC`

Allowed values (use exactly as shown):

- `utc`
- `local`

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanDroplet, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.droplet_id` | `string` | droplet unique identifier (DigitalOcean ID) |
| `status.outputs.ipv4_address` | `string` | primary IPv4 address (public if available, otherwise private) |
| `status.outputs.ipv6_address` | `string` | IPv6 address (if IPv6 was enabled) |
| `status.outputs.image_id` | `int64` | image ID of the droplet’s base image |
| `status.outputs.vpc_uuid` | `string` | VPC network UUID in which the droplet resides |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpc` | DigitalOceanVpc | `status.outputs.vpc_id` |
| `spec.volumeIds` | DigitalOceanVolume | `status.outputs.volume_id` |

## Referenced By

Fields on other kinds that can point at this resource:

| Kind | Field | Reads |
|---|---|---|
| DigitalOceanLoadBalancer | `spec.dropletIds` | `status.outputs.droplet_id` |

## See Also

- [Overview](../README.md)
