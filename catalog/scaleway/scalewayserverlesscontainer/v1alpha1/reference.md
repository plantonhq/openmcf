# ScalewayServerlessContainer

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `scaleway.planton.dev/v1alpha1`

ScalewayServerlessContainerSpec defines the specification for a
Scaleway serverless container deployment.

This is a composite resource that creates a container namespace,
the container itself, and optional cron triggers. The namespace is
an internal implementation detail -- users interact with the
container as a single resource.

**Scaleway serverless container model:**
Scaleway organizes containers into namespaces. Each namespace can
hold multiple containers, but in the Planton model we create one
namespace per container for clean lifecycle management and isolation.
Environment variables and secrets are set on the container (not the
namespace) for simplicity and clarity.

**Key difference from ScalewayServerlessFunction:**
Containers deploy pre-built Docker images (from any OCI registry)
instead of source code with a runtime. They expose a listening port,
support health checks, and offer fine-grained CPU/memory/scaling
controls suited for long-running HTTP services.

**Composition pattern:** Mid-tier resource (DAG Layer 2).
Upstream: `private_network_id` references ScalewayPrivateNetwork.
          `image.registry_endpoint` references ScalewayContainerRegistry.
Downstream: `domain_name` output for ScalewayDnsRecord.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.image` | `ScalewayServerlessContainerImage` | yes |  |  |
| `spec.image.registryEndpoint` | `string \| valueFrom` | yes |  | ScalewayContainerRegistry (`status.outputs.endpoint`) |
| `spec.image.name` | `string` | yes |  |  |
| `spec.image.tag` | `string` | yes |  |  |
| `spec.registrySha256` | `string` |  |  |  |
| `spec.port` | `uint32` |  | `8080` |  |
| `spec.privacy` | `enum` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.memoryLimitMb` | `uint32` |  | `256` |  |
| `spec.cpuLimit` | `uint32` |  |  |  |
| `spec.minScale` | `uint32` |  |  |  |
| `spec.maxScale` | `uint32` |  | `20` |  |
| `spec.timeoutSeconds` | `uint32` |  | `300` |  |
| `spec.httpOption` | `enum` |  |  |  |
| `spec.protocol` | `enum` |  |  |  |
| `spec.commands` | `[]string` |  |  |  |
| `spec.args` | `[]string` |  |  |  |
| `spec.env` | `ScalewayServerlessContainerEnv` |  |  |  |
| `spec.env.variables` | `[]ScalewayServerlessContainerEnvVar` |  |  |  |
| `spec.env.variables[].name` | `string` | yes |  |  |
| `spec.env.variables[].value` | `string` | yes |  |  |
| `spec.env.secrets` | `[]ScalewayServerlessContainerEnvVar` |  |  |  |
| `spec.env.secrets[].name` | `string` | yes |  |  |
| `spec.env.secrets[].value` | `string` | yes |  |  |
| `spec.privateNetworkId` | `string \| valueFrom` |  |  | ScalewayPrivateNetwork (`status.outputs.private_network_id`) |
| `spec.sandbox` | `string` |  |  |  |
| `spec.healthCheck` | `ScalewayServerlessContainerHealthCheck` |  |  |  |
| `spec.healthCheck.path` | `string` | yes |  |  |
| `spec.healthCheck.failureThreshold` | `uint32` |  |  |  |
| `spec.healthCheck.intervalSeconds` | `uint32` |  |  |  |
| `spec.scalingOption` | `ScalewayServerlessContainerScalingOption` |  |  |  |
| `spec.scalingOption.concurrentRequestsThreshold` | `uint32` |  |  |  |
| `spec.scalingOption.cpuUsageThreshold` | `uint32` |  |  |  |
| `spec.scalingOption.memoryUsageThreshold` | `uint32` |  |  |  |
| `spec.localStorageLimitMb` | `uint32` |  |  |  |
| `spec.deploy` | `bool` |  | `true` |  |
| `spec.cronTriggers` | `[]ScalewayServerlessContainerCronTrigger` |  |  |  |
| `spec.cronTriggers[].name` | `string` |  |  |  |
| `spec.cronTriggers[].schedule` | `string` | yes |  |  |
| `spec.cronTriggers[].args` | `string` | yes |  |  |

## Field Details

### spec.region

`string` · required

region where the container namespace and container are deployed.

Scaleway serverless containers are regional resources.
Valid regions: "fr-par", "nl-ams", "pl-waw".

- rule: {"string":{"minLen":"1"}}

### spec.image

`ScalewayServerlessContainerImage` · required

image specifies the container image to deploy.

The image is defined as a structured message with three fields:
registry endpoint, image name, and tag. The IaC module composes
these into the full image URL: "{registry_endpoint}/{name}:{tag}".

The `registry_endpoint` field supports `StringValueOrRef`, enabling
infra charts to create a DAG edge from ScalewayContainerRegistry
to this container. For external registries (Docker Hub, GHCR),
use a plain `value`.

- rule: {"required":true}

### spec.image.registryEndpoint

`string | valueFrom` · required

registry_endpoint is the base URL of the container registry.

For Scaleway Container Registry, use `valueFrom` to reference
the endpoint output of a ScalewayContainerRegistry resource.
For external registries, provide the URL as a plain `value`.

Format examples:
  "rg.fr-par.scw.cloud/my-namespace" (Scaleway)
  "docker.io/library" (Docker Hub)
  "ghcr.io/my-org" (GitHub Container Registry)

- references: ScalewayContainerRegistry (`status.outputs.endpoint`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayContainerRegistry, name: <that resource's name>, fieldPath: status.outputs.endpoint}} -- a bare string does not parse

### spec.image.name

`string` · required

name is the image name within the registry.

This is the repository path relative to the registry endpoint.
Examples: "my-app", "backend/api", "nginx".

- rule: {"string":{"minLen":"1"}}

### spec.image.tag

`string` · required

tag is the image tag to deploy.

Examples: "latest", "v1.2.3", "sha-abc1234", "production".
Use immutable tags (semantic versions, git SHAs) in production
for reproducible deployments.

- rule: {"string":{"minLen":"1"}}

### spec.registrySha256

`string`

registry_sha256 is a deployment trigger string.

When this value changes, the IaC module triggers a redeployment
of the container. This can be any string -- typically the SHA256
digest of the container image, but it could be a CI build number,
git commit SHA, or any value that changes when redeployment is
desired.

Analogous to `zip_hash` on ScalewayServerlessFunction.
Leave empty to skip change-detection-based redeployment.

### spec.port

`uint32`

port is the listening port exposed by the container.

The container must listen on this port for incoming HTTP requests.
Scaleway routes traffic to this port. Most containers use 8080
(the default) but any valid port is accepted.

- default: `8080`

### spec.privacy

`enum` · required

privacy controls how the container endpoint is authenticated.

- rule: privacy must be specified
- rule: {"required":true,"enum":{"definedOnly":true}}

Allowed values (use exactly as shown):

- `scaleway_serverless_container_privacy_unspecified` -- Unspecified (invalid).
- `privacy_public` -- Public: the container endpoint is publicly accessible without authentication.
- `privacy_private` -- Private: the container endpoint requires a valid authentication token. Tokens are managed separately via the Scaleway console or API (not bundled in this resource).

### spec.description

`string`

description of the container. Optional, human-readable.

### spec.memoryLimitMb

`uint32`

memory_limit_mb is the memory allocated to each container instance
in megabytes.

Determines the memory ceiling for the container. Common values:
128, 256, 512, 1024, 2048, 4096. Defaults to 256 MB.

- default: `256`

### spec.cpuLimit

`uint32`

cpu_limit is the vCPU allocated to each container instance in
milliCPU units (1000 = 1 vCPU).

When set to 0 or omitted, Scaleway auto-allocates CPU proportional
to the memory limit. Set explicitly for workloads with specific
CPU requirements. Common values: 70, 140, 280, 560, 1120.

### spec.minScale

`uint32`

min_scale is the minimum number of always-running container instances.

Set to 0 (default) for scale-to-zero behavior -- the container
stops when idle and incurs no compute charges.
Set to 1+ for always-warm instances (eliminates cold starts but
incurs continuous billing).

### spec.maxScale

`uint32`

max_scale is the maximum number of concurrent container instances.

The container auto-scales based on incoming workload but never
exceeds this limit. Defaults to 20 if not specified.

- default: `20`

### spec.timeoutSeconds

`uint32`

timeout_seconds is the maximum time in seconds a single request
can spend being processed before Scaleway terminates it.

Defaults to 300 seconds (5 minutes) if not specified.

- default: `300`

### spec.httpOption

`enum`

http_option controls HTTP/HTTPS behavior for the container endpoint.

Defaults to "enabled" (both HTTP and HTTPS allowed) if not specified.

Allowed values (use exactly as shown):

- `scaleway_serverless_container_http_option_unspecified` -- Unspecified -- defaults to "enabled" behavior.
- `enabled` -- Enabled: both HTTP and HTTPS requests are accepted.
- `redirected` -- Redirected: HTTP requests are automatically redirected to HTTPS.

### spec.protocol

`enum`

protocol selects the communication protocol between the Scaleway
gateway and the container.

Defaults to HTTP/1.1 if not specified. Use h2c for gRPC services
or HTTP/2 backends.

Allowed values (use exactly as shown):

- `scaleway_serverless_container_protocol_unspecified` -- Unspecified -- defaults to HTTP/1.1.
- `http1` -- HTTP/1.1: standard HTTP protocol. Suitable for most web services.
- `h2c` -- h2c: HTTP/2 over cleartext. Required for gRPC services and beneficial for multiplexed HTTP/2 backends.

### spec.commands

`[]string`

commands overrides the default CMD from the container image.

This is equivalent to Docker's CMD instruction. When set, it
replaces the image's default command entirely.

Example: ["node", "server.js"]

### spec.args

`[]string`

args overrides the default ENTRYPOINT arguments from the container
image.

These are passed as arguments to the command specified in the
`commands` field (or the image's default ENTRYPOINT).

Example: ["--port", "8080", "--workers", "4"]

### spec.env

`ScalewayServerlessContainerEnv`

env groups environment variables and secrets for the container.

Variables are non-secret and visible in logs/dashboards.
Secrets are encrypted at rest and masked in the Scaleway console.

Modeled as repeated name-value messages (Kubernetes-style) rather
than maps to preserve sort order and enable future `valueFrom`
extension. This is a component-local message (NOT shared with
ScalewayServerlessFunction).

### spec.env.variables

`[]ScalewayServerlessContainerEnvVar`

variables are non-secret environment variables.

These are visible in the Scaleway console and may appear in
container logs. Do not put sensitive values here -- use `secrets`.

### spec.env.variables[].name

`string` · required

name is the environment variable name (e.g., "DATABASE_URL").

- rule: {"string":{"minLen":"1"}}

### spec.env.variables[].value

`string` · required

value is the environment variable value.

- rule: {"string":{"minLen":"1"}}

### spec.env.secrets

`[]ScalewayServerlessContainerEnvVar`

secrets are encrypted environment variables.

Scaleway stores these encrypted at rest and masks them in the
console. Use for database URLs, API keys, tokens, and other
sensitive configuration.

### spec.env.secrets[].name

`string` · required

name is the environment variable name (e.g., "DATABASE_URL").

- rule: {"string":{"minLen":"1"}}

### spec.env.secrets[].value

`string` · required

value is the environment variable value.

- rule: {"string":{"minLen":"1"}}

### spec.privateNetworkId

`string | valueFrom`

private_network_id optionally connects the container to a Scaleway
Private Network for VPC-internal communication.

When connected, the container can reach resources on the Private
Network (databases, Redis clusters, other services) without
traversing the public internet.

Leave unset for containers that only need public internet access.

- references: ScalewayPrivateNetwork (`status.outputs.private_network_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: ScalewayPrivateNetwork, name: <that resource's name>, fieldPath: status.outputs.private_network_id}} -- a bare string does not parse

### spec.sandbox

`string`

sandbox selects the execution environment for the container.

Available sandboxes depend on the Scaleway platform version.
Common values: "v1" (standard), "v2" (enhanced security).
Leave empty to use the platform default.

### spec.healthCheck

`ScalewayServerlessContainerHealthCheck`

health_check configures HTTP health checking for the container.

When configured, Scaleway periodically sends HTTP requests to the
specified path and considers the container unhealthy after the
configured number of consecutive failures.

Leave unset to rely on Scaleway's default liveness detection
(TCP port check).

### spec.healthCheck.path

`string` · required

path is the HTTP path to probe for health checks.

The health check sends an HTTP GET to this path on the
container's configured port. The container should return
a 2xx status code when healthy.

Common values: "/health", "/healthz", "/ready", "/".

- rule: {"string":{"minLen":"1"}}

### spec.healthCheck.failureThreshold

`uint32`

failure_threshold is the number of consecutive health check
failures before the container instance is considered unhealthy.

A higher threshold tolerates transient failures but takes longer
to detect actual problems. Typical values: 2-5.

### spec.healthCheck.intervalSeconds

`uint32`

interval_seconds is the period between health check probes
in seconds.

Shorter intervals detect failures faster but generate more
probe traffic. Typical values: 10-60.

### spec.scalingOption

`ScalewayServerlessContainerScalingOption`

scaling_option configures autoscaling thresholds for the container.

These thresholds control when Scaleway adds or removes container
instances. Only one threshold should typically be set -- Scaleway
uses the first threshold that is breached.

Leave unset to use Scaleway's default scaling behavior.

### spec.scalingOption.concurrentRequestsThreshold

`uint32`

concurrent_requests_threshold is the target number of concurrent
requests per container instance before scaling up.

When the average concurrent requests across instances exceeds
this value, Scaleway adds instances (up to max_scale).
Set to 0 or omit to disable request-based scaling.

Typical values: 10-100 depending on request latency.

### spec.scalingOption.cpuUsageThreshold

`uint32`

cpu_usage_threshold is the target CPU usage percentage per
container instance before scaling up.

When average CPU usage exceeds this percentage, Scaleway adds
instances. Set to 0 or omit to disable CPU-based scaling.

Typical values: 50-80.

### spec.scalingOption.memoryUsageThreshold

`uint32`

memory_usage_threshold is the target memory usage percentage per
container instance before scaling up.

When average memory usage exceeds this percentage, Scaleway adds
instances. Set to 0 or omit to disable memory-based scaling.

Typical values: 60-90.

### spec.localStorageLimitMb

`uint32`

local_storage_limit_mb is the local (ephemeral) storage available
to each container instance in megabytes.

This storage is lost when the container instance is stopped or
replaced. Use for temporary files, caches, or build artifacts.
Leave at 0 or unset to use the platform default.

### spec.deploy

`bool`

deploy controls whether the container is deployed (started) after
creation or update.

When true (default), the IaC module triggers deployment immediately
after provisioning. Set to false to create the container
infrastructure without starting it -- useful for pre-provisioning
before a release.

- default: `true`

### spec.cronTriggers

`[]ScalewayServerlessContainerCronTrigger`

cron_triggers defines optional scheduled triggers for the container.

Each trigger creates a `scaleway_container_cron` resource that
invokes the container on the specified schedule with the given
JSON arguments.

Common patterns:
  Hourly cleanup:    schedule = "0 * * * *"
  Nightly backup:    schedule = "0 2 * * *"
  Every 5 minutes:   schedule = "*/5 * * * *"

### spec.cronTriggers[].name

`string`

name is an optional human-readable identifier for the trigger.

If not provided, Scaleway auto-generates a name. Providing a name
makes triggers easier to identify in the console and logs.

### spec.cronTriggers[].schedule

`string` · required

schedule is a UNIX CRON expression defining when the container
is invoked.

Examples:
  "0 * * * *"    -- every hour at minute 0
  "0 2 * * *"    -- daily at 2:00 AM
  "*/5 * * * *"  -- every 5 minutes
  "0 0 * * 0"    -- weekly on Sunday at midnight

- rule: {"string":{"minLen":"1"}}

### spec.cronTriggers[].args

`string` · required

args is a JSON string passed to the container's event object on
each scheduled invocation.

Must be valid JSON. Use "{}" for no arguments.

Example: '{"task": "cleanup", "batch_size": 100}'

- rule: {"string":{"minLen":"1"}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: ScalewayServerlessContainer, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.container_id` | `string` | The unique identifier of the deployed serverless container. Format: Scaleway-assigned UUID. Used for API operations, Scaleway CLI commands, and Terraform import. |
| `status.outputs.namespace_id` | `string` | The unique identifier of the container namespace. Useful for managing additional resources that reference the namespace (external cron triggers, tokens, or additional containers not managed by this resource). |
| `status.outputs.domain_name` | `string` | The native Scaleway domain name for invoking the container. This is the HTTPS endpoint automatically assigned by Scaleway. Downstream ScalewayDnsRecord resources can create CNAME records pointing to this domain for custom domain routing. Example: "mycontainer-abc123.containers.fnc.fr-par.scw.cloud" |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.image.registryEndpoint` | ScalewayContainerRegistry | `status.outputs.endpoint` |
| `spec.privateNetworkId` | ScalewayPrivateNetwork | `status.outputs.private_network_id` |

## See Also

- [Overview](../README.md)
