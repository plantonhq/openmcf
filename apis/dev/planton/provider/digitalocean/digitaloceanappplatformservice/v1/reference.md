# DigitalOceanAppPlatformService

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `digital-ocean.planton.dev/v1`

DigitalOceanAppPlatformServiceSpec defines the specification required to deploy a containerized service or application on DigitalOcean App Platform.
It focuses on essential fields (following the 80/20 principle) such as the service's source (either a git repository or a container image from DigitalOcean Container Registry), resource sizing, scaling, and optional custom domain configuration.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.serviceName` | `string` | yes |  |  |
| `spec.region` | `enum` | yes |  |  |
| `spec.serviceType` | `enum` | yes |  |  |
| `spec.gitSource` | `DigitalOceanAppPlatformGitSource` |  |  |  |
| `spec.gitSource.repoUrl` | `string` | yes |  |  |
| `spec.gitSource.branch` | `string` | yes |  |  |
| `spec.gitSource.buildCommand` | `string` |  |  |  |
| `spec.gitSource.runCommand` | `string` |  |  |  |
| `spec.imageSource` | `DigitalOceanAppPlatformRegistrySource` |  |  |  |
| `spec.imageSource.registry` | `string \| valueFrom` | yes |  | DigitalOceanContainerRegistry (`status.outputs.server_url`) |
| `spec.imageSource.repository` | `string` | yes |  |  |
| `spec.imageSource.tag` | `string` | yes |  |  |
| `spec.instanceSizeSlug` | `enum` | yes | `basic-xxs` |  |
| `spec.instanceCount` | `uint32` |  | `1` |  |
| `spec.enableAutoscale` | `bool` |  |  |  |
| `spec.minInstanceCount` | `uint32` |  |  |  |
| `spec.maxInstanceCount` | `uint32` |  |  |  |
| `spec.env` | `map<string, string>` |  |  |  |
| `spec.customDomain` | `string \| valueFrom` |  |  | DigitalOceanDnsZone (`spec.domain_name`) |

## Field Details

### spec.serviceName

`string` · required

name of the app (must be unique within the user's DigitalOcean account).
Constraints: should be DNS-friendly (e.g., lowercase alphanumeric and hyphens), maximum 63 characters.

- rule: {"required":true,"string":{"maxLen":"63","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?$"}}

### spec.region

`enum` · required

region in which to deploy the app (DigitalOcean data center region slug).

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

### spec.serviceType

`enum` · required

type of service being deployed (e.g., a web service that receives HTTP traffic, a background worker, or a one-off job).

- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_app_platform_service_type_unspecified`
- `web_service`
- `worker`
- `job`

### spec.gitSource

`DigitalOceanAppPlatformGitSource`

git repository source configuration (for App Platform to build and deploy from source code).

### spec.gitSource.repoUrl

`string` · required

repo_url is the URL of the git repository (HTTPS or git) containing the source code.

- rule: {"required":true}

### spec.gitSource.branch

`string` · required

branch specifies the git branch to deploy from.

- rule: {"required":true}

### spec.gitSource.buildCommand

`string`

build_command optionally overrides the default build command for the app.
Example: "npm run build". If not provided, DigitalOcean will auto-detect build steps or use defaults.

### spec.gitSource.runCommand

`string`

run_command optionally overrides the start command for the app.
Example: "npm start". If not provided, defaults are inferred from the build or Dockerfile.

### spec.imageSource

`DigitalOceanAppPlatformRegistrySource`

container image source configuration (deploy a pre-built image, typically from DigitalOcean Container Registry).

### spec.imageSource.registry

`string | valueFrom` · required

registry is a reference to a DigitalOceanContainerRegistry resource that hosts the image.
This typically provides the registry URL and ensures credentials are available for pulling the image.

- references: DigitalOceanContainerRegistry (`status.outputs.server_url`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanContainerRegistry, name: <that resource's name>, fieldPath: status.outputs.server_url}} -- a bare string does not parse

### spec.imageSource.repository

`string` · required

repository is the name of the repository in the registry containing the image.
For example, "myapp/backend".

- rule: {"required":true}

### spec.imageSource.tag

`string` · required

tag is the image tag to deploy.
For example, "latest" or a specific version like "v1.0.0".

- rule: {"required":true}

### spec.instanceSizeSlug

`enum` · required

instance_size_slug specifies the instance size (plan) for this service.
Determines the CPU/memory allocated per instance. Common options include "basic-xxs", "basic-xs", "basic-s", "basic-m", and professional tiers.
Default (if not specified by user): "basic-xxs".

- default: `basic-xxs`
- rule: {"required":true}

Allowed values (use exactly as shown):

- `digital_ocean_app_platform_instance_size_unspecified`
- `basic_xxs`
- `basic_xs`
- `basic_s`
- `basic_m`
- `basic_l`
- `professional_xs`
- `professional_s`
- `professional_m`
- `professional_l`
- `professional_xl`

### spec.instanceCount

`uint32`

instance_count is the number of instances (containers) to run for this service.
Default: 1.

- default: `1`

### spec.enableAutoscale

`bool`

enable_autoscale controls whether to use auto-scaling for this service.
If true, the service will automatically scale between the specified min and max instance counts based on load.
Default: false (manual scaling).

### spec.minInstanceCount

`uint32`

min_instance_count specifies the minimum number of instances to run when auto-scaling is enabled.
Required if enable_autoscale = true.

### spec.maxInstanceCount

`uint32`

max_instance_count specifies the maximum number of instances to run when auto-scaling is enabled.
Required if enable_autoscale = true.

### spec.env

`map<string, string>`

env is a map of environment variables to set in the app's runtime environment.
Keys are variable names and values are their corresponding values.

### spec.customDomain

`string | valueFrom`

custom_domain is an optional custom domain to use for the app, in addition to the default ondigitalocean.app hostname.
Provide a reference to a DigitalOceanDnsZone resource (typically its domain name). The system will create the necessary DNS records.

- references: DigitalOceanDnsZone (`spec.domain_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: DigitalOceanDnsZone, name: <that resource's name>, fieldPath: spec.domain_name}} -- a bare string does not parse

## Outputs

Reference an output from another manifest as `valueFrom: {kind: DigitalOceanAppPlatformService, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.app_id` | `string` | app_id is the unique identifier of the app (DigitalOcean App Platform application ID). |
| `status.outputs.default_hostname` | `string` | default_hostname is the default hostname assigned to the app (usually ending in "ondigitalocean.app"). |
| `status.outputs.live_url` | `string` | live_url is the publicly accessible URL (including protocol) of the deployed service. This may be the same as the default hostname with "https://" prefix, or a custom domain if one was configured. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.imageSource.registry` | DigitalOceanContainerRegistry | `status.outputs.server_url` |
| `spec.customDomain` | DigitalOceanDnsZone | `spec.domain_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
