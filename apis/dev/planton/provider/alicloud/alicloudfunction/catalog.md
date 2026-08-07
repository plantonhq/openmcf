# AliCloud Function

Deploys an Alibaba Cloud Function Compute v3 function with configurable runtime, compute sizing, VPC networking, SLS logging, custom container and custom runtime support, lifecycle hooks, NAS file system mounts, and GPU acceleration. FC v3 uses a service-less model where functions are top-level resources -- VPC access, logging, IAM role, and all configuration is set directly on the function. The component integrates with Planton's Provider Connections for AliCloud credential management and supports ValueFromRef wiring to RAM roles, VPCs, security groups, and SLS log projects.

## What Gets Created

When you deploy this Cloud Resource, the IaC module provisions:

- **FC v3 Function** -- an `alicloud_fcv3_function` with the configured runtime, handler, compute sizing, and optional VPC networking, SLS logging, NAS mounts, GPU configuration, and lifecycle hooks
- **AliCloud Tags** -- resource metadata tags (organization, environment, resource kind, resource ID) merged with user-provided `tags`

## Before You Deploy

### Planton Setup

- **AliCloud Provider Connection** -- an active connection in the Connect module with credentials for the target Alibaba Cloud account. Map it as the default for your environment, or specify it explicitly when creating the Cloud Resource.
- **Planton Runner** -- required when using Runner-based credential delivery. Not needed for inline credentials.

### Alibaba Cloud Account

- A code package in an OSS bucket (ZIP file) or a container image in ACR for `custom-container` runtime. Alternatively, use the inline `zipFile` field for small functions.
- A RAM role trusting the FC service principal (`fc.aliyuncs.com`) if the function accesses other Alibaba Cloud services. Provide the role ARN directly or reference an AliCloudRamRole Cloud Resource via ValueFromRef.
- A VPC, VSwitch(es), and security group if the function needs private network access to databases, caches, or NAS mount targets. Provide IDs directly or reference AliCloudVpc, AliCloudVswitch, and AliCloudSecurityGroup Cloud Resources via ValueFromRef.
- An SLS project and logstore if logging is configured. Provide the project name directly or reference an AliCloudLogProject Cloud Resource via ValueFromRef.
- A NAS file system with mount targets in the same VPC if NAS mounts are configured.

## Deploy

### Console

Open the deployment store, find **AliCloud Function**, and click **Deploy**. The creation wizard walks you through preset selection, environment and connection configuration, and spec fields. Start from the **Event Handler** preset in the [Presets](#presets) tab to pre-populate a lightweight Python function with SLS logging.

### CLI

```yaml
apiVersion: alicloud.planton.dev/v1
kind: AliCloudFunction
metadata:
  name: hello-world
  org: acme-corp
  env: prod
spec:
  region: cn-hangzhou
  functionName: hello-world
  handler: index.handler
  runtime: python3.12
  code:
    ossBucketName: my-code-bucket
    ossObjectName: functions/hello-world.zip
```

```shell
planton apply -f function.yaml
```

This creates an FC v3 function running Python 3.12 with code from the specified OSS location. No VPC access, logging, IAM role, or compute sizing overrides are configured -- provider defaults apply.

### InfraChart

When deploying as part of a multi-resource environment, use ValueFromRef to wire the function to a RAM role, VPC, security group, and SLS log project deployed in the same InfraPipeline:

```yaml
spec:
  role:
    valueFrom:
      kind: AliCloudRamRole
      name: fc-execution-role
      fieldPath: status.outputs.arn
  vpcConfig:
    vpcId:
      valueFrom:
        kind: AliCloudVpc
        name: app-vpc
        fieldPath: status.outputs.vpc_id
    securityGroupId:
      valueFrom:
        kind: AliCloudSecurityGroup
        name: function-sg
        fieldPath: status.outputs.security_group_id
  logConfig:
    project:
      valueFrom:
        kind: AliCloudLogProject
        name: app-logs
        fieldPath: status.outputs.project_name
```

The InfraPipeline resolves the dependency graph, deploys the RAM role, VPC, security group, and log project first, then provisions the function with the resolved values.

## Key Configuration

These are the most important decisions when configuring a function. Explore the full field reference in the [API Explorer](#api-explorer) tab.

**Runtime selection** -- Choose from built-in runtimes (`python3.12`, `nodejs20`, `java11`, `go1`), custom runtimes (`custom.debian12` -- you provide an HTTP server binary), or `custom-container` (you provide a container image). The runtime determines which code deployment and configuration fields apply. `functionName` and `runtime` are immutable after creation.

**Compute sizing** -- Set `cpu` (0.05-16 vCPUs), `memorySize` (64-32768 MB), `timeout` (1-86400 seconds), and `diskSize` (minimum 512 MB). Provider computes defaults when omitted. Set `instanceConcurrency` (1-200) to control how many concurrent requests a single instance handles -- higher values improve throughput but require thread-safe code.

**VPC networking** -- Configure `vpcConfig` to place the function inside a VPC for access to private resources (databases, caches, NAS). Requires `vpcId`, one or more `vswitchIds` across availability zones, and a `securityGroupId`. Without VPC configuration, the function runs in FC's managed network with no access to VPC-internal resources.

**Lifecycle hooks** -- Configure `instanceLifecycleConfig` with an `initializer` hook (runs once at instance creation for warm-up tasks like loading ML models or opening DB connections) and a `preStop` hook (runs before instance reclaim for cleanup tasks like flushing buffers).

**GPU acceleration** -- Configure `gpuConfig` with `gpuMemorySize` and `gpuType` (e.g., `fc.gpu.ampere.1`) for AI/ML inference workloads. Requires `custom-container` or custom runtime.

## Outputs and Dependencies

### What This Component Consumes

| Dependency | Field | ValueFromRef Path |
|------------|-------|-------------------|
| **AliCloudRamRole** (optional) | `role` | `status.outputs.arn` |
| **AliCloudVpc** (optional) | `vpcConfig.vpcId` | `status.outputs.vpc_id` |
| **AliCloudSecurityGroup** (optional) | `vpcConfig.securityGroupId` | `status.outputs.security_group_id` |
| **AliCloudLogProject** (optional) | `logConfig.project` | `status.outputs.project_name` |

### What This Component Provides

After provisioning, `status.outputs` contains values that downstream Cloud Resources can consume via ValueFromRef:

| Output | Description | Common Downstream Use |
|--------|-------------|----------------------|
| `function_id` | The FC function ID assigned by Alibaba Cloud | Trigger configuration, monitoring dashboards |
| `function_name` | The function name (mirrors the spec input) | API Gateway routing, event source mappings |
| `function_arn` | The function ARN (`acs:fc:{region}:{account-id}:functions/{function-name}`) | RAM policies, trigger source ARN references |

## Common Patterns

Browse the [Presets](#presets) tab for ready-to-deploy configurations.

**Event handler** -- A lightweight Python function with SLS logging for event-driven workloads (OSS triggers, timer triggers, message queue consumers). Start from the **Event Handler** preset.

**VPC API function** -- A Node.js function with VPC access for connecting to private databases and caches, SLS logging with instance and request metrics, and environment variables for service configuration. Start from the **VPC API Function** preset.

**Custom container function** -- A containerized function with a custom image, health check configuration, lifecycle hooks for warm-up and cleanup, and optional GPU acceleration for inference workloads. Start from the **Custom Container** preset.

## Works With

- [**AliCloud RAM Role**](/cloud-catalog/ali-cloud-ram-role) -- provides the execution role the function assumes during invocation
- [**AliCloud VPC**](/cloud-catalog/ali-cloud-vpc) -- provides the VPC for private network access to databases, caches, and NAS
- [**AliCloud Security Group**](/cloud-catalog/ali-cloud-security-group) -- provides network access control for the function's VPC-attached ENIs
- [**AliCloud Log Project**](/cloud-catalog/ali-cloud-log-project) -- provides the SLS project for function invocation and metrics logging