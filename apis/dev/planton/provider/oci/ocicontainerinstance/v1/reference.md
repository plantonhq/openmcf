# OciContainerInstance

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `oci.planton.dev/v1`

OciContainerInstanceSpec defines the specification for an Oracle Cloud
Infrastructure Container Instance.

A container instance is OCI's serverless container service (analogous to
AWS Fargate or Azure Container Instances). It runs one or more containers
in a pod-like construct that shares networking (VNICs) and volumes, without
requiring you to manage the underlying compute infrastructure.

Key characteristics:
  - Multiple containers per instance share the same network namespace
  - Volumes (emptydir, configfile) can be mounted across containers
  - Health checks (HTTP, TCP) enable automated container lifecycle management
  - Image pull secrets support both basic auth and OCI Vault-based credentials
  - Shapes are always flex (CI.Standard.E4.Flex, CI.Standard.E3.Flex)

The `state` lifecycle field (ACTIVE/INACTIVE) is intentionally omitted.
Planton resources are always deployed to their active state; delete the
resource to decommission it.

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.compartmentId` | `string \| valueFrom` | yes |  | OciCompartment (`status.outputs.compartment_id`) |
| `spec.availabilityDomain` | `string` | yes |  |  |
| `spec.displayName` | `string` |  |  |  |
| `spec.shape` | `string` | yes |  |  |
| `spec.shapeConfig` | `ShapeConfig` | yes |  |  |
| `spec.shapeConfig.ocpus` | `float` |  |  |  |
| `spec.shapeConfig.memoryInGbs` | `float` |  |  |  |
| `spec.containers` | `[]Container` | yes |  |  |
| `spec.containers[].imageUrl` | `string` | yes |  |  |
| `spec.containers[].displayName` | `string` |  |  |  |
| `spec.containers[].command` | `[]string` |  |  |  |
| `spec.containers[].arguments` | `[]string` |  |  |  |
| `spec.containers[].environmentVariables` | `map<string, string>` |  |  |  |
| `spec.containers[].workingDirectory` | `string` |  |  |  |
| `spec.containers[].isResourcePrincipalDisabled` | `bool` |  |  |  |
| `spec.containers[].resourceConfig` | `ContainerResourceConfig` |  |  |  |
| `spec.containers[].resourceConfig.memoryLimitInGbs` | `float` |  |  |  |
| `spec.containers[].resourceConfig.vcpusLimit` | `float` |  |  |  |
| `spec.containers[].healthChecks` | `[]HealthCheck` |  |  |  |
| `spec.containers[].healthChecks[].healthCheckType` | `enum` |  |  |  |
| `spec.containers[].healthChecks[].port` | `int32` |  |  |  |
| `spec.containers[].healthChecks[].name` | `string` |  |  |  |
| `spec.containers[].healthChecks[].path` | `string` |  |  |  |
| `spec.containers[].healthChecks[].failureAction` | `enum` |  |  |  |
| `spec.containers[].healthChecks[].failureThreshold` | `int32` |  |  |  |
| `spec.containers[].healthChecks[].successThreshold` | `int32` |  |  |  |
| `spec.containers[].healthChecks[].initialDelayInSeconds` | `int32` |  |  |  |
| `spec.containers[].healthChecks[].intervalInSeconds` | `int32` |  |  |  |
| `spec.containers[].healthChecks[].timeoutInSeconds` | `int32` |  |  |  |
| `spec.containers[].healthChecks[].headers` | `[]HealthCheckHeader` |  |  |  |
| `spec.containers[].healthChecks[].headers[].name` | `string` |  |  |  |
| `spec.containers[].healthChecks[].headers[].value` | `string` |  |  |  |
| `spec.containers[].securityContext` | `SecurityContext` |  |  |  |
| `spec.containers[].securityContext.isNonRootUserCheckEnabled` | `bool` |  |  |  |
| `spec.containers[].securityContext.isRootFileSystemReadonly` | `bool` |  |  |  |
| `spec.containers[].securityContext.runAsUser` | `int32` |  |  |  |
| `spec.containers[].securityContext.runAsGroup` | `int32` |  |  |  |
| `spec.containers[].securityContext.capabilities` | `Capabilities` |  |  |  |
| `spec.containers[].securityContext.capabilities.addCapabilities` | `[]string` |  |  |  |
| `spec.containers[].securityContext.capabilities.dropCapabilities` | `[]string` |  |  |  |
| `spec.containers[].volumeMounts` | `[]VolumeMount` |  |  |  |
| `spec.containers[].volumeMounts[].mountPath` | `string` | yes |  |  |
| `spec.containers[].volumeMounts[].volumeName` | `string` | yes |  |  |
| `spec.containers[].volumeMounts[].isReadOnly` | `bool` |  |  |  |
| `spec.containers[].volumeMounts[].partition` | `int32` |  |  |  |
| `spec.containers[].volumeMounts[].subPath` | `string` |  |  |  |
| `spec.vnics` | `[]Vnic` | yes |  |  |
| `spec.vnics[].subnetId` | `string \| valueFrom` | yes |  | OciSubnet (`status.outputs.subnet_id`) |
| `spec.vnics[].displayName` | `string` |  |  |  |
| `spec.vnics[].hostnameLabel` | `string` |  |  |  |
| `spec.vnics[].isPublicIpAssigned` | `bool` |  |  |  |
| `spec.vnics[].nsgIds` | `[]string \| valueFrom` |  |  | OciSecurityGroup (`status.outputs.network_security_group_id`) |
| `spec.vnics[].privateIp` | `string` |  |  |  |
| `spec.vnics[].skipSourceDestCheck` | `bool` |  |  |  |
| `spec.containerRestartPolicy` | `enum` |  |  |  |
| `spec.faultDomain` | `string` |  |  |  |
| `spec.gracefulShutdownTimeoutInSeconds` | `int64` |  |  |  |
| `spec.dnsConfig` | `DnsConfig` |  |  |  |
| `spec.dnsConfig.nameservers` | `[]string` |  |  |  |
| `spec.dnsConfig.options` | `[]string` |  |  |  |
| `spec.dnsConfig.searches` | `[]string` |  |  |  |
| `spec.imagePullSecrets` | `[]ImagePullSecret` |  |  |  |
| `spec.imagePullSecrets[].registryEndpoint` | `string` | yes |  |  |
| `spec.imagePullSecrets[].secretType` | `enum` |  |  |  |
| `spec.imagePullSecrets[].username` | `string` |  |  |  |
| `spec.imagePullSecrets[].password` | `string` (sensitive) |  |  |  |
| `spec.imagePullSecrets[].secretId` | `string \| valueFrom` |  |  |  |
| `spec.volumes` | `[]Volume` |  |  |  |
| `spec.volumes[].name` | `string` | yes |  |  |
| `spec.volumes[].volumeType` | `enum` |  |  |  |
| `spec.volumes[].backingStore` | `string` |  |  |  |
| `spec.volumes[].configs` | `[]VolumeConfig` |  |  |  |
| `spec.volumes[].configs[].data` | `string` |  |  |  |
| `spec.volumes[].configs[].fileName` | `string` |  |  |  |
| `spec.volumes[].configs[].path` | `string` |  |  |  |

## Field Details

### spec.compartmentId

`string | valueFrom` · required

OCID of the compartment where the container instance will be created.

- references: OciCompartment (`status.outputs.compartment_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciCompartment, name: <that resource's name>, fieldPath: status.outputs.compartment_id}} -- a bare string does not parse

### spec.availabilityDomain

`string` · required

Availability domain where the container instance runs.
Example: "Uocm:PHX-AD-1".

- rule: {"string":{"minLen":"1"}}

### spec.displayName

`string`

Human-readable name for the container instance shown in the OCI Console.
Falls back to metadata.name if not provided.

### spec.shape

`string` · required

Compute shape for the container instance.
Example: "CI.Standard.E4.Flex", "CI.Standard.E3.Flex".

- rule: {"string":{"minLen":"1"}}

### spec.shapeConfig

`ShapeConfig` · required

CPU and memory allocation for the entire container instance.
Individual containers can set resource limits within this envelope.

- rule: {"required":true}

### spec.shapeConfig.ocpus

`float`

Number of OCPUs allocated to the container instance.
Example: 1.0, 2.0, 4.0.

- rule: {"float":{"gt":0}}

### spec.shapeConfig.memoryInGbs

`float`

Memory in gigabytes. When omitted, OCI assigns a default based
on the OCPU count (typically 1 GB per OCPU minimum).

### spec.containers

`[]Container` · required

Containers to run on this instance. At least one container is required.
Multiple containers share the same network namespace and can communicate
over localhost.

- rule: {"repeated":{"minItems":"1"}}

### spec.containers[].imageUrl

`string` · required

Container image URL.
Example: "docker.io/library/nginx:latest", "ghcr.io/org/app:v1.2".
Default registry is docker.io/library if not specified.

- rule: {"string":{"minLen":"1"}}

### spec.containers[].displayName

`string`

Human-readable name for the container.

### spec.containers[].command

`[]string`

Overrides the image's ENTRYPOINT. Each element is a separate argument.

### spec.containers[].arguments

`[]string`

Arguments passed to the ENTRYPOINT process.
Total size of all arguments combined must be <= 64 KB.

### spec.containers[].environmentVariables

`map<string, string>`

Environment variables injected into the container.
Total size of all names + values combined must be <= 64 KB.

### spec.containers[].workingDirectory

`string`

Working directory for the container's entrypoint process.

### spec.containers[].isResourcePrincipalDisabled

`bool`

When true, disables OCI resource principal access for this container.
Resource principal (v2.2) is enabled by default.

### spec.containers[].resourceConfig

`ContainerResourceConfig`

CPU and memory limits for this container. When omitted, the container
can use all resources available to the instance.

### spec.containers[].resourceConfig.memoryLimitInGbs

`float`

Maximum memory in gigabytes the container can consume.
When omitted, the container can use all available instance memory.

### spec.containers[].resourceConfig.vcpusLimit

`float`

Maximum logical CPUs the container can consume.
1 OCPU = 2 logical CPUs. Values can be fractional (e.g., 0.5).
When omitted, the container can use all available instance CPUs.

### spec.containers[].healthChecks

`[]HealthCheck`

Health checks for monitoring container readiness. Supports HTTP and
TCP probe types with configurable thresholds and intervals.

### spec.containers[].healthChecks[].healthCheckType

`enum`

Protocol for the health check.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `health_check_type_unspecified`
- `http`
- `tcp`

### spec.containers[].healthChecks[].port

`int32`

Port to probe.

- rule: {"int32":{"gt":0}}

### spec.containers[].healthChecks[].name

`string`

Optional name for the health check, unique within the container instance.

### spec.containers[].healthChecks[].path

`string`

URL path for HTTP health checks. Required when health_check_type is http.
Example: "/healthz".

### spec.containers[].healthChecks[].failureAction

`enum`

Action to take when the health check fails.
When omitted, defaults to KILL.

Allowed values (use exactly as shown):

- `failure_action_unspecified`
- `kill`
- `none`

### spec.containers[].healthChecks[].failureThreshold

`int32`

Consecutive failures required to consider the container unhealthy.

### spec.containers[].healthChecks[].successThreshold

`int32`

Consecutive successes required to consider the container healthy again.

### spec.containers[].healthChecks[].initialDelayInSeconds

`int32`

Seconds to wait after container start before running the first check.

### spec.containers[].healthChecks[].intervalInSeconds

`int32`

Seconds between consecutive health checks.

### spec.containers[].healthChecks[].timeoutInSeconds

`int32`

Seconds to wait for a health check response before considering it failed.

### spec.containers[].healthChecks[].headers

`[]HealthCheckHeader`

Custom HTTP headers sent with HTTP health checks.

### spec.containers[].healthChecks[].headers[].name

`string`

### spec.containers[].healthChecks[].headers[].value

`string`

### spec.containers[].securityContext

`SecurityContext`

Linux security settings for the container process.

### spec.containers[].securityContext.isNonRootUserCheckEnabled

`bool`

When true, validates at runtime that the container does not run as
UID 0. Fails the container start if the image runs as root.

### spec.containers[].securityContext.isRootFileSystemReadonly

`bool`

When true, the container's root filesystem is mounted read-only.

### spec.containers[].securityContext.runAsUser

`int32`

User ID (UID) for the container's entrypoint process.
Defaults to the UID specified in the container image.

### spec.containers[].securityContext.runAsGroup

`int32`

Group ID (GID) for the container's entrypoint process.
When specified, run_as_user should also be provided.

### spec.containers[].securityContext.capabilities

`Capabilities`

Linux capabilities to add or drop from the container process.

### spec.containers[].securityContext.capabilities.addCapabilities

`[]string`

Capabilities to add to the container process.
Example: ["NET_ADMIN", "SYS_TIME"]

### spec.containers[].securityContext.capabilities.dropCapabilities

`[]string`

Capabilities to drop from the container process.
Example: ["ALL"]

### spec.containers[].volumeMounts

`[]VolumeMount`

Volumes to mount into this container's filesystem.

### spec.containers[].volumeMounts[].mountPath

`string` · required

Path inside the container where the volume is mounted.
Example: "/data", "/etc/config".

- rule: {"string":{"minLen":"1"}}

### spec.containers[].volumeMounts[].volumeName

`string` · required

Name of the volume to mount. Must match a volume defined in the
instance-level volumes list.

- rule: {"string":{"minLen":"1"}}

### spec.containers[].volumeMounts[].isReadOnly

`bool`

When true, the volume is mounted read-only.

### spec.containers[].volumeMounts[].partition

`int32`

If the volume has partitions, the partition number to mount.

### spec.containers[].volumeMounts[].subPath

`string`

Sub-path within the volume to mount instead of the volume root.

### spec.vnics

`[]Vnic` · required

Virtual network interface cards providing network connectivity.
At least one VNIC is required. Each VNIC is attached to a subnet.

- rule: {"repeated":{"minItems":"1"}}

### spec.vnics[].subnetId

`string | valueFrom` · required

OCID of the subnet in which to create the VNIC.

- references: OciSubnet (`status.outputs.subnet_id`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSubnet, name: <that resource's name>, fieldPath: status.outputs.subnet_id}} -- a bare string does not parse

### spec.vnics[].displayName

`string`

Human-readable name for the VNIC.

### spec.vnics[].hostnameLabel

`string`

Hostname label for the VNIC's primary private IP in subnet DNS.

### spec.vnics[].isPublicIpAssigned

`bool` · optional (explicit presence)

Whether to assign a public IP to the VNIC.
When omitted, uses the subnet's default public IP assignment setting.

### spec.vnics[].nsgIds

`[]string | valueFrom`

OCIDs of network security groups to add this VNIC to.

- references: OciSecurityGroup (`status.outputs.network_security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: OciSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.network_security_group_id}} -- a bare string does not parse

### spec.vnics[].privateIp

`string`

Static private IP address within the subnet's CIDR.
When omitted, OCI assigns one automatically.

### spec.vnics[].skipSourceDestCheck

`bool`

When true, disables source/destination checking on the VNIC.
Required for NAT instances or virtual routers.

### spec.containerRestartPolicy

`enum`

Restart policy applied to all containers in this instance.
When omitted, defaults to ALWAYS.

Allowed values (use exactly as shown):

- `restart_policy_unspecified`
- `always`
- `never`
- `on_failure`

### spec.faultDomain

`string`

Fault domain within the availability domain.
When omitted, OCI selects a fault domain automatically.
Example: "FAULT-DOMAIN-1".

### spec.gracefulShutdownTimeoutInSeconds

`int64`

Seconds to wait for containers to gracefully terminate before
forcefully stopping them. Applies when the instance is stopped or deleted.

### spec.dnsConfig

`DnsConfig`

DNS resolver configuration for containers. When omitted, containers
inherit DNS settings from the subnet's DHCP options.

### spec.dnsConfig.nameservers

`[]string`

IP addresses of DNS name servers (IPv4 or IPv6).
When omitted, uses nameservers from the subnet's DHCP options.

### spec.dnsConfig.options

`[]string`

Resolver options in resolv.conf format.
Example: ["ndots:5", "edns0"]

### spec.dnsConfig.searches

`[]string`

Search domains for unqualified hostname lookups.
When omitted, uses searches from the subnet's DHCP options.

### spec.imagePullSecrets

`[]ImagePullSecret`

Credentials for pulling container images from private registries.

### spec.imagePullSecrets[].registryEndpoint

`string` · required

Registry endpoint URL.
Example: "ghcr.io", "docker.io", "us-ashburn-1.ocir.io".

- rule: {"string":{"minLen":"1"}}

### spec.imagePullSecrets[].secretType

`enum`

Authentication method for the registry.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `secret_type_unspecified`
- `basic`
- `vault`

### spec.imagePullSecrets[].username

`string`

Username for basic authentication. Required when secret_type is basic.
Must be base64-encoded.

### spec.imagePullSecrets[].password

`string` · sensitive

Password for basic authentication. Required when secret_type is basic.
Must be base64-encoded.

### spec.imagePullSecrets[].secretId

`string | valueFrom`

OCID of an OCI Vault secret containing registry credentials.
Required when secret_type is vault.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.volumes

`[]Volume`

Volumes accessible to containers via volume mounts. A container
instance supports up to 32 volumes.

### spec.volumes[].name

`string` · required

Unique name for this volume within the container instance.
Containers reference this name in their volume_mounts.

- rule: {"string":{"minLen":"1"}}

### spec.volumes[].volumeType

`enum`

Storage backing for the volume.

- rule: {"enum":{"notIn":[0]}}

Allowed values (use exactly as shown):

- `volume_type_unspecified`
- `emptydir`
- `configfile`

### spec.volumes[].backingStore

`string`

Backing store for emptydir volumes. Options: "EPHEMERAL_STORAGE"
(disk-backed) or "MEMORY" (tmpfs). Only applicable when volume_type
is emptydir.

### spec.volumes[].configs

`[]VolumeConfig`

Config file entries for configfile volumes. Each entry becomes a file
in the volume. Only applicable when volume_type is configfile.

### spec.volumes[].configs[].data

`string`

Base64-encoded contents of the file. Decoded to plain text at mount time.

### spec.volumes[].configs[].fileName

`string`

Name of the file within the volume. Must be unique across the volume.

### spec.volumes[].configs[].path

`string`

Optional relative path within the volume mount directory.
When omitted, the file is placed at the volume mount root.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: OciContainerInstance, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.container_instance_id` | `string` | OCID of the container instance. |
| `status.outputs.container_ids` | `string` | Comma-separated OCIDs of the individual containers within the instance. Useful for operational tasks (viewing logs, exec into container). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.compartmentId` | OciCompartment | `status.outputs.compartment_id` |
| `spec.vnics[].subnetId` | OciSubnet | `status.outputs.subnet_id` |
| `spec.vnics[].nsgIds` | OciSecurityGroup | `status.outputs.network_security_group_id` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
