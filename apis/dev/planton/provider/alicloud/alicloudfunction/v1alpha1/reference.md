# AliCloudFunction

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudFunctionSpec defines the configuration for an Alibaba Cloud Function
Compute v3 function.

Function Compute (FC) is a fully managed, event-driven serverless compute
service. FC v3 uses a service-less model where functions are top-level
resources (no alicloud_fc_service grouping required). VPC access, logging,
and IAM role are configured directly on the function.

This component wraps a single alicloud_fcv3_function resource. Triggers,
aliases, versions, and concurrency configs have independent lifecycles and
are not bundled here.

Supported runtime families:
  - Built-in: python3.12, nodejs20, java11, go1, php7.2, dotnetcore3.1
  - Custom: custom, custom.debian10, custom.debian11, custom.debian12
  - Container: custom-container (requires custom_container_config)

Provider resources:
  Terraform: alicloud_fcv3_function
  Pulumi:    fc.V3Function

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudFunction
metadata:
  name: alicloudfunction-demo
spec:
  region: cn-hangzhou
  functionName: planton-demo-func
  handler: index.handler
  runtime: python3.12
  code:
    ossBucketName: my-code-bucket
    ossObjectName: functions/demo.zip
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.functionName` | `string` | yes |  |  |
| `spec.handler` | `string` | yes |  |  |
| `spec.runtime` | `string` | yes |  |  |
| `spec.description` | `string` |  |  |  |
| `spec.cpu` | `double` |  |  |  |
| `spec.memorySize` | `int32` |  |  |  |
| `spec.timeout` | `int32` |  |  |  |
| `spec.diskSize` | `int32` |  |  |  |
| `spec.instanceConcurrency` | `int32` |  |  |  |
| `spec.code` | `AliCloudFunctionCode` |  |  |  |
| `spec.code.ossBucketName` | `string` |  |  |  |
| `spec.code.ossObjectName` | `string` |  |  |  |
| `spec.code.zipFile` | `string` |  |  |  |
| `spec.code.checksum` | `string` |  |  |  |
| `spec.role` | `string \| valueFrom` |  |  | AliCloudRamRole (`status.outputs.arn`) |
| `spec.internetAccess` | `bool` |  |  |  |
| `spec.vpcConfig` | `AliCloudFunctionVpcConfig` |  |  |  |
| `spec.vpcConfig.vpcId` | `string \| valueFrom` |  |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.vpcConfig.vswitchIds` | `[]string \| valueFrom` |  |  |  |
| `spec.vpcConfig.securityGroupId` | `string \| valueFrom` |  |  | AliCloudSecurityGroup (`status.outputs.security_group_id`) |
| `spec.logConfig` | `AliCloudFunctionLogConfig` |  |  |  |
| `spec.logConfig.project` | `string \| valueFrom` |  |  | AliCloudLogProject (`status.outputs.project_name`) |
| `spec.logConfig.logstore` | `string` |  |  |  |
| `spec.logConfig.logBeginRule` | `string` |  |  |  |
| `spec.logConfig.enableInstanceMetrics` | `bool` |  |  |  |
| `spec.logConfig.enableRequestMetrics` | `bool` |  |  |  |
| `spec.customContainerConfig` | `AliCloudFunctionCustomContainerConfig` |  |  |  |
| `spec.customContainerConfig.image` | `string` |  |  |  |
| `spec.customContainerConfig.entrypoint` | `[]string` |  |  |  |
| `spec.customContainerConfig.command` | `[]string` |  |  |  |
| `spec.customContainerConfig.port` | `int32` |  |  |  |
| `spec.customContainerConfig.healthCheckConfig` | `AliCloudFunctionHealthCheckConfig` |  |  |  |
| `spec.customContainerConfig.healthCheckConfig.initialDelaySeconds` | `int32` |  |  |  |
| `spec.customContainerConfig.healthCheckConfig.timeoutSeconds` | `int32` |  |  |  |
| `spec.customContainerConfig.healthCheckConfig.httpGetUrl` | `string` |  |  |  |
| `spec.customContainerConfig.healthCheckConfig.periodSeconds` | `int32` |  |  |  |
| `spec.customContainerConfig.healthCheckConfig.failureThreshold` | `int32` |  |  |  |
| `spec.customContainerConfig.healthCheckConfig.successThreshold` | `int32` |  |  |  |
| `spec.customRuntimeConfig` | `AliCloudFunctionCustomRuntimeConfig` |  |  |  |
| `spec.customRuntimeConfig.command` | `[]string` |  |  |  |
| `spec.customRuntimeConfig.args` | `[]string` |  |  |  |
| `spec.customRuntimeConfig.port` | `int32` |  |  |  |
| `spec.customRuntimeConfig.healthCheckConfig` | `AliCloudFunctionHealthCheckConfig` |  |  |  |
| `spec.customRuntimeConfig.healthCheckConfig.initialDelaySeconds` | `int32` |  |  |  |
| `spec.customRuntimeConfig.healthCheckConfig.timeoutSeconds` | `int32` |  |  |  |
| `spec.customRuntimeConfig.healthCheckConfig.httpGetUrl` | `string` |  |  |  |
| `spec.customRuntimeConfig.healthCheckConfig.periodSeconds` | `int32` |  |  |  |
| `spec.customRuntimeConfig.healthCheckConfig.failureThreshold` | `int32` |  |  |  |
| `spec.customRuntimeConfig.healthCheckConfig.successThreshold` | `int32` |  |  |  |
| `spec.instanceLifecycleConfig` | `AliCloudFunctionInstanceLifecycleConfig` |  |  |  |
| `spec.instanceLifecycleConfig.initializer` | `AliCloudFunctionLifecycleHook` |  |  |  |
| `spec.instanceLifecycleConfig.initializer.handler` | `string` |  |  |  |
| `spec.instanceLifecycleConfig.initializer.timeout` | `int32` |  |  |  |
| `spec.instanceLifecycleConfig.initializer.command` | `[]string` |  |  |  |
| `spec.instanceLifecycleConfig.preStop` | `AliCloudFunctionLifecycleHook` |  |  |  |
| `spec.instanceLifecycleConfig.preStop.handler` | `string` |  |  |  |
| `spec.instanceLifecycleConfig.preStop.timeout` | `int32` |  |  |  |
| `spec.instanceLifecycleConfig.preStop.command` | `[]string` |  |  |  |
| `spec.nasConfig` | `AliCloudFunctionNasConfig` |  |  |  |
| `spec.nasConfig.userId` | `int32` |  |  |  |
| `spec.nasConfig.groupId` | `int32` |  |  |  |
| `spec.nasConfig.mountPoints` | `[]AliCloudFunctionNasMountPoint` |  |  |  |
| `spec.nasConfig.mountPoints[].serverAddr` | `string` |  |  |  |
| `spec.nasConfig.mountPoints[].mountDir` | `string` |  |  |  |
| `spec.nasConfig.mountPoints[].enableTls` | `bool` |  |  |  |
| `spec.gpuConfig` | `AliCloudFunctionGpuConfig` |  |  |  |
| `spec.gpuConfig.gpuMemorySize` | `int32` |  |  |  |
| `spec.gpuConfig.gpuType` | `string` |  |  |  |
| `spec.layers` | `[]string` |  |  |  |
| `spec.environmentVariables` | `map<string, string>` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |
| `spec.resourceGroupId` | `string` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the function will be created.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.functionName

`string` · required

Function name. Must be unique within the region and account.
This field is immutable after creation (ForceNew in the provider).
1-128 characters, starting with a letter or underscore.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"128"}}

### spec.handler

`string` · required

Entry point for function invocation.
Format depends on runtime: "index.handler" (Node.js/Python),
"com.example.Main::handleRequest" (Java), "main" (Go).
Required for all runtimes including custom-container (provider enforces it).

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.runtime

`string` · required

Runtime environment for the function code.

- rule: runtime must be a valid FC v3 runtime
- rule: {"required":true}

### spec.description

`string`

Human-readable description of the function's purpose.

### spec.cpu

`double` · optional (explicit presence)

vCPU allocation for each function instance. Range: 0.05-16.
Provider computes a default based on memory_size when omitted.

- rule: cpu must be between 0.05 and 16

### spec.memorySize

`int32` · optional (explicit presence)

Memory allocation in MB for each function instance. Range: 64-32768.
Provider computes a default when omitted.

- rule: memory_size must be between 64 and 32768

### spec.timeout

`int32` · optional (explicit presence)

Maximum execution time in seconds before the function is terminated.
Range: 1-86400 (24 hours).

- rule: timeout must be between 1 and 86400

### spec.diskSize

`int32` · optional (explicit presence)

Temporary disk size in MB available to the function instance.
Minimum: 512 MB. Provider computes a default when omitted.

- rule: disk_size must be at least 512

### spec.instanceConcurrency

`int32` · optional (explicit presence)

Maximum number of concurrent requests a single instance can handle.
Range: 1-200. Higher values improve throughput but require the function
code to be safe for concurrent execution.

- rule: instance_concurrency must be between 1 and 200

### spec.code

`AliCloudFunctionCode`

Code package for the function. Not required when runtime is
"custom-container" (the container image is specified in
custom_container_config instead).

### spec.code.ossBucketName

`string`

OSS bucket name where the function code ZIP package is stored.

### spec.code.ossObjectName

`string`

OSS object key for the function code ZIP package.

### spec.code.zipFile

`string`

Base64-encoded function code ZIP package. Convenient for small functions
or testing. Mutually exclusive with oss_bucket_name/oss_object_name in
practice (the provider accepts both but only one source is meaningful).

### spec.code.checksum

`string`

CRC-64 checksum of the code package for integrity verification.

### spec.role

`string | valueFrom`

RAM role ARN that the function assumes during execution.
The role must trust the FC service principal (fc.aliyuncs.com).

- references: AliCloudRamRole (`status.outputs.arn`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudRamRole, name: <that resource's name>, fieldPath: status.outputs.arn}} -- a bare string does not parse

### spec.internetAccess

`bool` · optional (explicit presence)

Whether the function can access the public internet.
When false, outbound internet access is blocked even if the function
is not in a VPC.

### spec.vpcConfig

`AliCloudFunctionVpcConfig`

VPC configuration for functions that need to access VPC-internal
resources (databases, caches, NAS). When set, the function runs inside
the specified VPC and can reach private endpoints.

### spec.vpcConfig.vpcId

`string | valueFrom`

VPC ID. Computed by the provider from the VSwitch when omitted in TF,
but required in the Planton model for explicit cross-component wiring.

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.vpcConfig.vswitchIds

`[]string | valueFrom`

VSwitch IDs to place the function's ENIs in. Multiple VSwitches across
availability zones improve resilience.

- rule: write as {value: <literal>} or {valueFrom: {kind: <Kind>, name: <that resource's name>, fieldPath: status.outputs.<output>}} -- a bare string does not parse

### spec.vpcConfig.securityGroupId

`string | valueFrom`

Security group ID applied to the function's ENIs.

- references: AliCloudSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.logConfig

`AliCloudFunctionLogConfig`

Log Service (SLS) configuration for function invocation logs.

### spec.logConfig.project

`string | valueFrom`

SLS project name. References AliCloudLogProject.

- references: AliCloudLogProject (`status.outputs.project_name`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudLogProject, name: <that resource's name>, fieldPath: status.outputs.project_name}} -- a bare string does not parse

### spec.logConfig.logstore

`string`

SLS logstore name within the project.

### spec.logConfig.logBeginRule

`string` · optional (explicit presence)

Log begin rule for parsing log boundaries.
"None" disables log splitting. "DefaultRegex" uses the default regex.

- rule: log_begin_rule must be one of: None, DefaultRegex

### spec.logConfig.enableInstanceMetrics

`bool` · optional (explicit presence)

Enable per-instance resource usage metrics (CPU, memory).

### spec.logConfig.enableRequestMetrics

`bool` · optional (explicit presence)

Enable per-request latency and status metrics.

### spec.customContainerConfig

`AliCloudFunctionCustomContainerConfig`

Configuration for custom container runtimes (runtime="custom-container").
Specifies the container image, entrypoint, and health check.

### spec.customContainerConfig.image

`string`

Container image URI. Required when this config is set.
Example: "registry.cn-hangzhou.aliyuncs.com/my-ns/my-func:v1"

- rule: image must not be empty when custom_container_config is set

### spec.customContainerConfig.entrypoint

`[]string`

Container entrypoint override (ENTRYPOINT in Dockerfile terms).

### spec.customContainerConfig.command

`[]string`

Container command override (CMD in Dockerfile terms).

### spec.customContainerConfig.port

`int32` · optional (explicit presence)

Port the container listens on for HTTP requests.

### spec.customContainerConfig.healthCheckConfig

`AliCloudFunctionHealthCheckConfig`

Health check configuration for the container.

### spec.customContainerConfig.healthCheckConfig.initialDelaySeconds

`int32` · optional (explicit presence)

Seconds to wait after instance start before the first health check.
Range: 0-120.

- rule: initial_delay_seconds must be between 0 and 120

### spec.customContainerConfig.healthCheckConfig.timeoutSeconds

`int32` · optional (explicit presence)

Seconds to wait for a health check response before marking it failed.
Range: 0-3.

- rule: timeout_seconds must be between 0 and 3

### spec.customContainerConfig.healthCheckConfig.httpGetUrl

`string`

HTTP GET URL path for the health check (e.g., "/healthz").

### spec.customContainerConfig.healthCheckConfig.periodSeconds

`int32` · optional (explicit presence)

Seconds between consecutive health checks. Range: 0-120.

- rule: period_seconds must be between 0 and 120

### spec.customContainerConfig.healthCheckConfig.failureThreshold

`int32` · optional (explicit presence)

Number of consecutive failures before the instance is marked unhealthy.
Range: 1-120.

- rule: failure_threshold must be between 1 and 120

### spec.customContainerConfig.healthCheckConfig.successThreshold

`int32` · optional (explicit presence)

Number of consecutive successes before the instance is marked healthy.
Range: 1-120.

- rule: success_threshold must be between 1 and 120

### spec.customRuntimeConfig

`AliCloudFunctionCustomRuntimeConfig`

Configuration for custom runtimes (runtime="custom.*").
Specifies the bootstrap command, arguments, listening port, and health check.

### spec.customRuntimeConfig.command

`[]string`

Bootstrap command to start the custom runtime server.

### spec.customRuntimeConfig.args

`[]string`

Arguments passed to the bootstrap command.

### spec.customRuntimeConfig.port

`int32` · optional (explicit presence)

Port the custom runtime listens on. Range: 0-65535.

- rule: port must be between 0 and 65535

### spec.customRuntimeConfig.healthCheckConfig

`AliCloudFunctionHealthCheckConfig`

Health check configuration for the custom runtime.

### spec.customRuntimeConfig.healthCheckConfig.initialDelaySeconds

`int32` · optional (explicit presence)

Seconds to wait after instance start before the first health check.
Range: 0-120.

- rule: initial_delay_seconds must be between 0 and 120

### spec.customRuntimeConfig.healthCheckConfig.timeoutSeconds

`int32` · optional (explicit presence)

Seconds to wait for a health check response before marking it failed.
Range: 0-3.

- rule: timeout_seconds must be between 0 and 3

### spec.customRuntimeConfig.healthCheckConfig.httpGetUrl

`string`

HTTP GET URL path for the health check (e.g., "/healthz").

### spec.customRuntimeConfig.healthCheckConfig.periodSeconds

`int32` · optional (explicit presence)

Seconds between consecutive health checks. Range: 0-120.

- rule: period_seconds must be between 0 and 120

### spec.customRuntimeConfig.healthCheckConfig.failureThreshold

`int32` · optional (explicit presence)

Number of consecutive failures before the instance is marked unhealthy.
Range: 1-120.

- rule: failure_threshold must be between 1 and 120

### spec.customRuntimeConfig.healthCheckConfig.successThreshold

`int32` · optional (explicit presence)

Number of consecutive successes before the instance is marked healthy.
Range: 1-120.

- rule: success_threshold must be between 1 and 120

### spec.instanceLifecycleConfig

`AliCloudFunctionInstanceLifecycleConfig`

Lifecycle hooks for function instances. The initializer runs once when
an instance is created (warm-up logic), and pre_stop runs before an
instance is reclaimed (cleanup logic).

### spec.instanceLifecycleConfig.initializer

`AliCloudFunctionLifecycleHook`

Initializer hook. Runs once when a new instance is created, before it
receives any invocations. Use for warm-up tasks (loading models,
opening DB connections).

### spec.instanceLifecycleConfig.initializer.handler

`string`

Entry point for the hook (e.g., "index.initializer").

### spec.instanceLifecycleConfig.initializer.timeout

`int32` · optional (explicit presence)

Maximum execution time in seconds for the hook.
Initializer: 0-600. Pre-stop: 0-900.

- rule: timeout must be between 0 and 900

### spec.instanceLifecycleConfig.initializer.command

`[]string`

Bootstrap command for the hook (FC v3 supports command-based hooks
in addition to handler-based hooks).

### spec.instanceLifecycleConfig.preStop

`AliCloudFunctionLifecycleHook`

Pre-stop hook. Runs before an idle instance is reclaimed. Use for
cleanup tasks (flushing buffers, closing connections).

### spec.instanceLifecycleConfig.preStop.handler

`string`

Entry point for the hook (e.g., "index.initializer").

### spec.instanceLifecycleConfig.preStop.timeout

`int32` · optional (explicit presence)

Maximum execution time in seconds for the hook.
Initializer: 0-600. Pre-stop: 0-900.

- rule: timeout must be between 0 and 900

### spec.instanceLifecycleConfig.preStop.command

`[]string`

Bootstrap command for the hook (FC v3 supports command-based hooks
in addition to handler-based hooks).

### spec.nasConfig

`AliCloudFunctionNasConfig`

NAS file system mount configuration. Requires vpc_config to be set
(NAS mount targets are VPC-internal).

### spec.nasConfig.userId

`int32` · optional (explicit presence)

POSIX user ID for file access. Typically 0 (root) or a custom UID.

### spec.nasConfig.groupId

`int32` · optional (explicit presence)

POSIX group ID for file access.

### spec.nasConfig.mountPoints

`[]AliCloudFunctionNasMountPoint`

NAS mount points to attach to the function instance.

### spec.nasConfig.mountPoints[].serverAddr

`string`

NAS mount target address.
Format: "{file-system-id}-{mount-target-id}.{region}.nas.aliyuncs.com:/{path}"

- rule: server_addr must not be empty

### spec.nasConfig.mountPoints[].mountDir

`string`

Directory inside the function instance where the NAS path is mounted.
Must start with "/mnt/" (e.g., "/mnt/data").

- rule: mount_dir must not be empty

### spec.nasConfig.mountPoints[].enableTls

`bool` · optional (explicit presence)

Enable TLS encryption for NAS traffic.

### spec.gpuConfig

`AliCloudFunctionGpuConfig`

GPU configuration for AI/ML inference workloads.

### spec.gpuConfig.gpuMemorySize

`int32`

GPU memory in MB allocated to each function instance.

- rule: {"int32":{"gt":0}}

### spec.gpuConfig.gpuType

`string`

GPU hardware type.

- rule: gpu_type must be one of: fc.gpu.tesla.1, fc.gpu.ampere.1, fc.gpu.ada.1, g1

### spec.layers

`[]string`

Layer ARNs to attach to the function. Layers provide shared libraries
and dependencies without bundling them in the function code package.
Maximum: 5 layers.

### spec.environmentVariables

`map<string, string>`

Environment variables passed to the function at runtime.

### spec.tags

`map<string, string>`

Tags to apply to the function resource.

### spec.resourceGroupId

`string`

Alibaba Cloud resource group ID for organizational grouping (per DD05).

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudFunction, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.function_id` | `string` | The function ID assigned by Alibaba Cloud. |
| `status.outputs.function_name` | `string` | The function name (mirrors the spec input for downstream reference). |
| `status.outputs.function_arn` | `string` | The function ARN (Alibaba Cloud Resource Name). Format: acs:fc:{region}:{account-id}:functions/{function-name} Used in RAM policies and as a trigger source ARN. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.role` | AliCloudRamRole | `status.outputs.arn` |
| `spec.vpcConfig.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |
| `spec.vpcConfig.securityGroupId` | AliCloudSecurityGroup | `status.outputs.security_group_id` |
| `spec.logConfig.project` | AliCloudLogProject | `status.outputs.project_name` |

## See Also

- [Overview](./README.md)
- [Design notes](./docs/README.md)
