# AliCloudSaeApplication

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `alicloud.planton.dev/v1alpha1`

AliCloudSaeApplicationSpec defines the configuration for an Alibaba Cloud
Serverless App Engine (SAE) application.

SAE is a fully managed, container-based serverless compute platform that
combines the simplicity of PaaS with the flexibility of containers. It
supports deploying applications as container images, JAR packages, WAR
packages, or Python/PHP ZIP archives. SAE handles provisioning, scaling,
load balancing, and log collection automatically.

CPU and memory are specified as discrete tiers (millicores and MB),
not arbitrary values. Each instance runs a single copy of the application
and can be scaled horizontally via the replicas field.

Provider resources:
  Terraform: alicloud_sae_application
  Pulumi:    sae.Application

## Example

```yaml
apiVersion: alicloud.planton.dev/v1alpha1
kind: AliCloudSaeApplication
metadata:
  name: alicloudsaeapplication-demo
spec:
  region: cn-hangzhou
  appName: sae-demo
  packageType: Image
  replicas: 1
  cpu: 500
  memory: 1024
  imageUrl: registry.cn-hangzhou.aliyuncs.com/my-ns/demo:latest
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.region` | `string` | yes |  |  |
| `spec.appName` | `string` | yes |  |  |
| `spec.appDescription` | `string` |  |  |  |
| `spec.packageType` | `string` | yes |  |  |
| `spec.replicas` | `int32` | yes |  |  |
| `spec.cpu` | `int32` | yes |  |  |
| `spec.memory` | `int32` | yes |  |  |
| `spec.vpcId` | `string \| valueFrom` |  |  | AliCloudVpc (`status.outputs.vpc_id`) |
| `spec.vswitchId` | `string \| valueFrom` |  |  | AliCloudVswitch (`status.outputs.vswitch_id`) |
| `spec.securityGroupId` | `string \| valueFrom` |  |  | AliCloudSecurityGroup (`status.outputs.security_group_id`) |
| `spec.namespaceId` | `string` |  |  |  |
| `spec.imageUrl` | `string` |  |  |  |
| `spec.packageUrl` | `string` |  |  |  |
| `spec.packageVersion` | `string` |  |  |  |
| `spec.command` | `string` |  |  |  |
| `spec.commandArgs` | `[]string` |  |  |  |
| `spec.envs` | `map<string, string>` |  |  |  |
| `spec.jdk` | `string` |  |  |  |
| `spec.jarStartOptions` | `string` |  |  |  |
| `spec.jarStartArgs` | `string` |  |  |  |
| `spec.programmingLanguage` | `string` |  |  |  |
| `spec.timezone` | `string` |  |  |  |
| `spec.terminationGracePeriodSeconds` | `int32` |  |  |  |
| `spec.minReadyInstances` | `int32` |  |  |  |
| `spec.acrInstanceId` | `string` |  |  |  |
| `spec.liveness` | `AliCloudSaeApplicationHealthCheck` |  |  |  |
| `spec.liveness.httpGet` | `AliCloudSaeApplicationHttpGetAction` |  |  |  |
| `spec.liveness.httpGet.path` | `string` |  |  |  |
| `spec.liveness.httpGet.port` | `int32` |  |  |  |
| `spec.liveness.tcpSocket` | `AliCloudSaeApplicationTcpSocketAction` |  |  |  |
| `spec.liveness.tcpSocket.port` | `int32` |  |  |  |
| `spec.liveness.exec` | `AliCloudSaeApplicationExecAction` |  |  |  |
| `spec.liveness.exec.command` | `string` |  |  |  |
| `spec.liveness.initialDelaySeconds` | `int32` |  |  |  |
| `spec.liveness.periodSeconds` | `int32` |  |  |  |
| `spec.liveness.timeoutSeconds` | `int32` |  |  |  |
| `spec.liveness.failureThreshold` | `int32` |  |  |  |
| `spec.liveness.successThreshold` | `int32` |  |  |  |
| `spec.readiness` | `AliCloudSaeApplicationHealthCheck` |  |  |  |
| `spec.readiness.httpGet` | `AliCloudSaeApplicationHttpGetAction` |  |  |  |
| `spec.readiness.httpGet.path` | `string` |  |  |  |
| `spec.readiness.httpGet.port` | `int32` |  |  |  |
| `spec.readiness.tcpSocket` | `AliCloudSaeApplicationTcpSocketAction` |  |  |  |
| `spec.readiness.tcpSocket.port` | `int32` |  |  |  |
| `spec.readiness.exec` | `AliCloudSaeApplicationExecAction` |  |  |  |
| `spec.readiness.exec.command` | `string` |  |  |  |
| `spec.readiness.initialDelaySeconds` | `int32` |  |  |  |
| `spec.readiness.periodSeconds` | `int32` |  |  |  |
| `spec.readiness.timeoutSeconds` | `int32` |  |  |  |
| `spec.readiness.failureThreshold` | `int32` |  |  |  |
| `spec.readiness.successThreshold` | `int32` |  |  |  |
| `spec.customHostAliases` | `[]AliCloudSaeApplicationCustomHostAlias` |  |  |  |
| `spec.customHostAliases[].hostName` | `string` |  |  |  |
| `spec.customHostAliases[].ip` | `string` |  |  |  |
| `spec.updateStrategy` | `AliCloudSaeApplicationUpdateStrategy` |  |  |  |
| `spec.updateStrategy.type` | `string` |  |  |  |
| `spec.updateStrategy.batchUpdate` | `AliCloudSaeApplicationBatchUpdate` |  |  |  |
| `spec.updateStrategy.batchUpdate.batch` | `int32` |  |  |  |
| `spec.updateStrategy.batchUpdate.batchWaitTime` | `int32` |  |  |  |
| `spec.updateStrategy.batchUpdate.releaseType` | `string` |  |  |  |
| `spec.slsConfigs` | `string` |  |  |  |
| `spec.tags` | `map<string, string>` |  |  |  |

## Field Details

### spec.region

`string` · required

Alibaba Cloud region where the SAE application will be created.
Examples: "cn-hangzhou", "cn-shanghai", "us-west-1", "ap-southeast-1".

- rule: {"required":true,"string":{"minLen":"1"}}

### spec.appName

`string` · required

Application name. Must start with a letter, and can contain letters,
digits, and dashes. Maximum 36 characters.
This field is immutable after creation (ForceNew in the provider).

- rule: {"required":true,"string":{"minLen":"1","maxLen":"36"}}

### spec.appDescription

`string`

Human-readable description of the application. Maximum 1024 characters.

- rule: {"string":{"maxLen":"1024"}}

### spec.packageType

`string` · required

Application package type. Determines how the application code is
delivered and which runtime fields apply.
This field is immutable after creation (ForceNew in the provider).

- rule: package_type must be one of: Image, FatJar, War, PythonZip, PhpZip
- rule: {"required":true}

### spec.replicas

`int32` · required

Number of application instances to run. Minimum 1.

- rule: {"required":true,"int32":{"gte":1}}

### spec.cpu

`int32` · required

CPU allocation per instance in millicores.

- rule: cpu must be one of: 500, 1000, 2000, 4000, 8000, 16000, 32000
- rule: {"required":true}

### spec.memory

`int32` · required

Memory allocation per instance in MB.

- rule: memory must be one of: 1024, 2048, 4096, 8192, 12288, 16384, 24576, 32768, 65536, 131072
- rule: {"required":true}

### spec.vpcId

`string | valueFrom`

VPC ID for VPC-based deployment. When set, the application runs inside
the specified VPC. Optional — SAE provides managed networking by default.
This field is immutable after creation (ForceNew in the provider).

- references: AliCloudVpc (`status.outputs.vpc_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVpc, name: <that resource's name>, fieldPath: status.outputs.vpc_id}} -- a bare string does not parse

### spec.vswitchId

`string | valueFrom`

VSwitch ID for VPC-based deployment. Determines the subnet and
availability zone.

- references: AliCloudVswitch (`status.outputs.vswitch_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudVswitch, name: <that resource's name>, fieldPath: status.outputs.vswitch_id}} -- a bare string does not parse

### spec.securityGroupId

`string | valueFrom`

Security group ID for VPC-based deployment. Controls inbound/outbound
traffic rules for the application instances.

- references: AliCloudSecurityGroup (`status.outputs.security_group_id`)
- rule: write as {value: <literal>} or {valueFrom: {kind: AliCloudSecurityGroup, name: <that resource's name>, fieldPath: status.outputs.security_group_id}} -- a bare string does not parse

### spec.namespaceId

`string`

SAE namespace ID. Format: "{region}:{short_id}". Applications in the
same namespace share configuration items (ConfigMaps) and can discover
each other via built-in service registration.
If omitted, the application is placed in the default namespace.
This field is immutable after creation (ForceNew in the provider).

### spec.imageUrl

`string`

Container image URL. Required when package_type is "Image".
Example: "registry.cn-hangzhou.aliyuncs.com/my-ns/my-app:v1"

### spec.packageUrl

`string`

Deployment package URL (OSS or HTTP). Required when package_type is
"FatJar", "War", "PythonZip", or "PhpZip".

### spec.packageVersion

`string`

Deployment package version identifier. Recommended for non-Image
deployments to track which version is deployed.

### spec.command

`string`

Container start command override (ENTRYPOINT equivalent).

### spec.commandArgs

`[]string`

Container start command arguments (CMD equivalent).
Maps to the provider's command_args_v2 field.

### spec.envs

`map<string, string>`

Environment variables passed to the application at runtime.
The IaC modules convert this map to the JSON array format that
the SAE API expects: [{"name":"K","value":"V"},...].

### spec.jdk

`string`

JDK version for Java applications (package_type FatJar or War).
Examples: "Open JDK 8", "Open JDK 11", "Open JDK 17", "Dragonwell 8",
"Dragonwell 11", "Dragonwell 17".

### spec.jarStartOptions

`string`

JVM startup options for FatJar applications.
Example: "-Xms512m -Xmx1024m -Dserver.port=8080"

### spec.jarStartArgs

`string`

Application startup arguments for FatJar applications.

### spec.programmingLanguage

`string` · optional (explicit presence)

Programming language. Affects which runtime features are available.
This field is immutable after creation (ForceNew in the provider).

- rule: programming_language must be one of: java, php, other

### spec.timezone

`string`

Application timezone. Affects log timestamps and scheduled tasks.

### spec.terminationGracePeriodSeconds

`int32` · optional (explicit presence)

Graceful shutdown timeout in seconds. The application receives SIGTERM
and has this many seconds to clean up before SIGKILL.
Range: 1-60. Provider default: 30.

- rule: termination_grace_period_seconds must be between 1 and 60

### spec.minReadyInstances

`int32` · optional (explicit presence)

Minimum number of available instances during a rolling deployment.
Ensures service availability while new instances are starting.

### spec.acrInstanceId

`string`

ACR Enterprise Edition instance ID for pulling images from a private
container registry. Only needed when image_url points to an ACR EE
instance (not the default ACR Personal Edition).

### spec.liveness

`AliCloudSaeApplicationHealthCheck`

Liveness probe configuration. If the liveness check fails repeatedly,
SAE restarts the application instance.

### spec.liveness.httpGet

`AliCloudSaeApplicationHttpGetAction`

HTTP GET probe. SAE sends an HTTP GET request to the specified path
and port, and considers the check successful if the response status
is 200-399.

### spec.liveness.httpGet.path

`string`

URL path for the HTTP GET request (e.g., "/healthz").

### spec.liveness.httpGet.port

`int32`

Port to send the HTTP GET request to.

### spec.liveness.tcpSocket

`AliCloudSaeApplicationTcpSocketAction`

TCP socket probe. SAE attempts a TCP connection to the specified port
and considers the check successful if the connection is established.

### spec.liveness.tcpSocket.port

`int32`

Port to attempt a TCP connection to.

### spec.liveness.exec

`AliCloudSaeApplicationExecAction`

Exec probe. SAE executes the specified command inside the container
and considers the check successful if the command exits with code 0.

### spec.liveness.exec.command

`string`

Command to execute inside the container. The check succeeds if the
command exits with code 0.

### spec.liveness.initialDelaySeconds

`int32` · optional (explicit presence)

Seconds to wait after the container starts before initiating the
first health check.

### spec.liveness.periodSeconds

`int32` · optional (explicit presence)

Seconds between consecutive health checks.

### spec.liveness.timeoutSeconds

`int32` · optional (explicit presence)

Seconds to wait for a health check response before marking it failed.

### spec.liveness.failureThreshold

`int32` · optional (explicit presence)

Number of consecutive failures before considering the container
unhealthy (liveness) or not ready (readiness).

### spec.liveness.successThreshold

`int32` · optional (explicit presence)

Number of consecutive successes before considering the container
healthy or ready again after a failure.

### spec.readiness

`AliCloudSaeApplicationHealthCheck`

Readiness probe configuration. Until the readiness check passes, SAE
does not route traffic to the instance.

### spec.readiness.httpGet

`AliCloudSaeApplicationHttpGetAction`

HTTP GET probe. SAE sends an HTTP GET request to the specified path
and port, and considers the check successful if the response status
is 200-399.

### spec.readiness.httpGet.path

`string`

URL path for the HTTP GET request (e.g., "/healthz").

### spec.readiness.httpGet.port

`int32`

Port to send the HTTP GET request to.

### spec.readiness.tcpSocket

`AliCloudSaeApplicationTcpSocketAction`

TCP socket probe. SAE attempts a TCP connection to the specified port
and considers the check successful if the connection is established.

### spec.readiness.tcpSocket.port

`int32`

Port to attempt a TCP connection to.

### spec.readiness.exec

`AliCloudSaeApplicationExecAction`

Exec probe. SAE executes the specified command inside the container
and considers the check successful if the command exits with code 0.

### spec.readiness.exec.command

`string`

Command to execute inside the container. The check succeeds if the
command exits with code 0.

### spec.readiness.initialDelaySeconds

`int32` · optional (explicit presence)

Seconds to wait after the container starts before initiating the
first health check.

### spec.readiness.periodSeconds

`int32` · optional (explicit presence)

Seconds between consecutive health checks.

### spec.readiness.timeoutSeconds

`int32` · optional (explicit presence)

Seconds to wait for a health check response before marking it failed.

### spec.readiness.failureThreshold

`int32` · optional (explicit presence)

Number of consecutive failures before considering the container
unhealthy (liveness) or not ready (readiness).

### spec.readiness.successThreshold

`int32` · optional (explicit presence)

Number of consecutive successes before considering the container
healthy or ready again after a failure.

### spec.customHostAliases

`[]AliCloudSaeApplicationCustomHostAlias`

Custom hostname-to-IP mappings injected into the container's /etc/hosts.
Useful for resolving internal service names without DNS.

### spec.customHostAliases[].hostName

`string`

Hostname to map.

### spec.customHostAliases[].ip

`string`

IP address the hostname resolves to.

### spec.updateStrategy

`AliCloudSaeApplicationUpdateStrategy`

Deployment update strategy controlling how new versions are rolled out.

### spec.updateStrategy.type

`string` · optional (explicit presence)

Release policy type.
"BatchUpdate": release in sequential batches.
"GrayBatchUpdate": canary-style phased release.

- rule: type must be one of: BatchUpdate, GrayBatchUpdate

### spec.updateStrategy.batchUpdate

`AliCloudSaeApplicationBatchUpdate`

Batch update configuration.

### spec.updateStrategy.batchUpdate.batch

`int32` · optional (explicit presence)

Number of batches to split the release into.

### spec.updateStrategy.batchUpdate.batchWaitTime

`int32` · optional (explicit presence)

Seconds to wait between batches.

### spec.updateStrategy.batchUpdate.releaseType

`string` · optional (explicit presence)

Release processing method.
"auto": automatically proceed to the next batch.
"manual": pause between batches for manual approval.

- rule: release_type must be one of: auto, manual

### spec.slsConfigs

`string`

SLS log collection configuration as a JSON string.
Format: [{"logDir":"/path","logType":"stdout"},...]
When set, SAE collects application logs and ships them to SLS.

### spec.tags

`map<string, string>`

Tags to apply to the SAE application resource.

## Outputs

Reference an output from another manifest as `valueFrom: {kind: AliCloudSaeApplication, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.app_id` | `string` | The SAE application ID assigned by Alibaba Cloud. |
| `status.outputs.app_name` | `string` | The application name (mirrors the spec input for downstream reference). |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.vpcId` | AliCloudVpc | `status.outputs.vpc_id` |
| `spec.vswitchId` | AliCloudVswitch | `status.outputs.vswitch_id` |
| `spec.securityGroupId` | AliCloudSecurityGroup | `status.outputs.security_group_id` |

## See Also

- [Overview](../README.md)
