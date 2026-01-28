# Pulumi Module Architecture: AWS Lambda

## Overview

This Pulumi module provisions AWS Lambda functions through a declarative, protobuf-defined specification. Lambda functions are serverless compute units that run code in response to events without managing servers, supporting both zip-based deployments (from S3) and container image deployments (from ECR).

The module follows OpenMCF's standard pattern: **input transformation → resource provisioning → output extraction**, with special handling for Lambda's dual packaging models and VPC networking complexity.

## Module Structure

```
iac/pulumi/
├── main.go              # Pulumi entrypoint (loads stack input, calls module)
├── Pulumi.yaml          # Pulumi project metadata
├── Makefile             # Build/test helpers
├── debug.sh             # Delve debugging wrapper
└── module/
    ├── main.go          # Orchestration logic (provider setup, output export)
    ├── locals.go        # Input transformation and tag construction
    ├── function.go      # Lambda function and log group resource implementation
    └── outputs.go       # Output key constants
```

### File Responsibilities

#### `main.go` (entrypoint)
- Unmarshals `AwsLambdaStackInput` from Pulumi stack configuration
- Delegates to `module.Resources()` for actual provisioning
- Minimal logic—purely a thin entrypoint wrapper

**Standard Pulumi Pattern**: This file is identical across all AWS components, only varying in the protobuf type imported.

#### `module/main.go` (orchestrator)
**Key Function:** `Resources(ctx, stackInput)`

Responsibilities:
1. **Locals Initialization**: Transform stack input into typed locals struct
2. **Provider Configuration**: Handle two provider scenarios:
   - **Default**: Create provider with ambient AWS credentials (IAM role, environment variables)
   - **Explicit**: Create provider with credentials from `stackInput.ProviderConfig` (access key, secret key, session token)
3. **Resource Creation**: Invoke `lambdaFunction()` with initialized locals and provider
4. **Error Propagation**: Wrap errors with context for debugging

**Design Decision**: Provider handling supports both CI/CD environments (IRSA, instance profiles) and local development (explicit credentials).

#### `module/locals.go` (input transformer)
**Key Function:** `initializeLocals(ctx, stackInput)`

Transforms the protobuf `AwsLambdaStackInput` into a strongly-typed `Locals` struct and constructs AWS resource tags.

**Tag Construction**: Every Lambda function receives six tags:
- `resource=true`: Marks this as a OpenMCF managed resource
- `organization=<org>`: Organization ID from metadata
- `environment=<env>`: Environment (dev/staging/prod) from metadata
- `resource-kind=AwsLambda`: CloudResourceKind enum string
- `resource-id=<id>`: Unique resource identifier from metadata
- `name=<name>`: Function name from metadata

**Why Tags Matter**: These tags enable:
- Cost allocation reporting by org/env/resource-kind (Lambda costs can be significant)
- Policy enforcement (e.g., "only functions tagged env=prod can access production secrets")
- Resource discovery and inventory management
- CloudWatch Logs Insights filtering by resource metadata

#### `module/function.go` (resource implementation)
**Key Function:** `lambdaFunction(ctx, locals, provider)`

This is the core implementation file containing all Lambda function provisioning logic. It translates the protobuf `AwsLambdaSpec` into Pulumi's Lambda resources.

**Implementation Pattern**: The module follows a **multi-phase construction** approach:

1. **Phase 1: Architecture Resolution**
   - Map protobuf enum `Architecture` to AWS strings ("x86_64", "arm64")
   - Default to x86_64 if not specified

2. **Phase 2: Environment Variables**
   - Convert protobuf map to Pulumi StringMap
   - Handle empty environment gracefully

3. **Phase 3: Layer ARN Resolution**
   - Extract layer ARNs from `StringValueOrRef` (supports literals and references)
   - Build layer ARN array

4. **Phase 4: VPC Configuration**
   - Extract subnet IDs and security group IDs from `StringValueOrRef`
   - Only create VPC config if subnets or security groups are provided

5. **Phase 5: CloudWatch Log Group**
   - Create dedicated log group with standard name pattern (`/aws/lambda/{functionName}`)
   - Set 30-day retention by default
   - Apply resource tags

6. **Phase 6: Lambda Function Creation**
   - Build `FunctionArgs` based on code source type (S3 vs container image)
   - Set memory, timeout, KMS key, layers, VPC config
   - Handle conditional fields based on code source type

7. **Phase 7: Output Exports**
   - Export function ARN, name, log group name, role ARN, and layer ARNs

#### `module/outputs.go` (constants)
Defines string constants for output keys, preventing typos and providing a single source of truth.

**Exported Outputs**:
- `function_arn`: Full ARN of the Lambda function
- `function_name`: Name of the function (for API calls, event source mappings)
- `log_group_name`: CloudWatch Logs group name (for log queries, metric filters)
- `role_arn`: Execution role ARN (for policy validation)
- `layer_arns`: Attached layer ARNs (for version tracking)

## Data Flow Diagram

```
┌──────────────────────────────────┐
│ AwsLambdaStackInput (protobuf)   │
│  ├─ target: AwsLambda            │
│  │   ├─ metadata                 │
│  │   └─ spec: AwsLambdaSpec      │
│  │       ├─ function_name        │
│  │       ├─ code_source_type     │
│  │       ├─ runtime/handler OR   │
│  │       │   image_uri           │
│  │       ├─ role_arn             │
│  │       ├─ memory_mb            │
│  │       ├─ timeout_seconds      │
│  │       ├─ environment          │
│  │       ├─ vpc config           │
│  │       └─ layers               │
│  └─ provider_config (optional)   │
└──────────────┬───────────────────┘
               │
               ▼
      ┌────────────────┐
      │initializeLocals│
      └────────┬────────┘
               │ Creates:
               │ - Locals.AwsLambda
               │ - Locals.AwsTags (6 tags)
               ▼
     ┌──────────────────────┐
     │ AWS Provider Setup   │
     │  ├─ Ambient creds OR │
     │  └─ Explicit creds   │
     └──────────┬───────────┘
                │
                ▼
       ┌──────────────────┐
       │ lambdaFunction() │
       └────────┬─────────┘
                │
                ├─ Map architecture enum:
                │    X86_64 → "x86_64"
                │    ARM64 → "arm64"
                │
                ├─ Extract environment variables
                │
                ├─ Resolve layer ARNs (StringValueOrRef)
                │
                ├─ Resolve VPC config (StringValueOrRef)
                │
                ├─ Create CloudWatch Log Group:
                │    cloudwatch.NewLogGroup(ctx, "log-group", &LogGroupArgs{
                │      Name: "/aws/lambda/{functionName}",
                │      RetentionInDays: 30,
                │      Tags: awsTags
                │    })
                │
                ├─ Build Lambda FunctionArgs based on code source:
                │    
                │    IF code_source_type == S3:
                │      - Set Runtime (nodejs18.x, python3.11, etc.)
                │      - Set Handler (index.handler, etc.)
                │      - Set S3Bucket, S3Key, S3ObjectVersion
                │      - PackageType = "Zip" (default)
                │    
                │    IF code_source_type == IMAGE:
                │      - Set ImageUri (ECR image URI)
                │      - Set PackageType = "Image"
                │      - Runtime/Handler ignored (defined in image)
                │
                ├─ Set common configuration:
                │    - Description
                │    - Role ARN
                │    - Architectures (x86_64 or arm64)
                │    - MemorySize, Timeout
                │    - Environment (if provided)
                │    - KmsKeyArn (if provided)
                │    - Layers (if provided)
                │    - VpcConfig (if subnets/SGs provided)
                │    - Tags
                │
                ├─ Create Lambda Function:
                │    lambda.NewFunction(ctx, functionName, args, provider)
                │
                └─ Export outputs:
                     ctx.Export("function_arn", fn.Arn)
                     ctx.Export("function_name", fn.Name)
                     ctx.Export("log_group_name", logGroup.Name)
                     ctx.Export("role_arn", roleArn)
                     ctx.Export("layer_arns", layerArns)
                │
                ▼
       ┌─────────────────────────────┐
       │  AWS Lambda & CloudWatch    │
       │   ├─ Lambda Function         │
       │   │   ├─ Code (S3 or Image)  │
       │   │   ├─ Execution Role      │
       │   │   ├─ VPC ENIs (optional) │
       │   │   └─ Environment Config  │
       │   └─ CloudWatch Log Group    │
       │       ├─ 30-day retention    │
       │       └─ Tagged for tracking │
       └─────────────────────────────┘
```

## Resource Relationships

```
AwsLambdaSpec
  │
  ├─ code_source_type ────────────┐
  │    (S3 or IMAGE)               │
  │                                │
  ├─ runtime/handler OR ───────┐  │
  │   image_uri                 │  │
  │                             │  │
  ├─ role_arn ─────────────────┼──┼─ Lambda Function
  │    (StringValueOrRef)       │  │   ├─ Name
  │                             │  │   ├─ ARN
  ├─ memory_mb ────────────────┼──┼───┤   ├─ Code Source (S3 or Image)
  ├─ timeout_seconds ──────────┼──┼───┤   ├─ Runtime (if S3)
  ├─ reserved_concurrency ─────┼──┼───┤   ├─ Handler (if S3)
  │                             │  │   ├─ Role ARN
  ├─ environment ──────────────┼──┼───┤   ├─ Memory/Timeout
  │                             │  │   ├─ Architecture (x86_64/arm64)
  ├─ subnets ──────────────────┼──┼───┤   ├─ Environment Variables
  │    (StringValueOrRef[])    │  │   │   │   └─ Optional KMS encryption
  ├─ security_groups ──────────┼──┼───┤   │
  │    (StringValueOrRef[])    │  │   ├─ VPC Config (optional)
  │                             │  │   │   ├─ Subnets
  ├─ architecture ─────────────┼──┼───┤   │   └─ Security Groups
  │                             │  │   │
  ├─ layer_arns ───────────────┼──┼───┤   ├─ Layers (up to 5)
  │    (StringValueOrRef[])    │  │   │
  │                             │  │   └─ Tags (org/env/resource metadata)
  └─ kms_key_arn ──────────────┘  │
       (StringValueOrRef)         │
                                  │
                                  ├─ CloudWatch Log Group
                                  │   ├─ Name: /aws/lambda/{functionName}
                                  │   ├─ Retention: 30 days
                                  │   └─ Tags
                                  │
                                  └─ Outputs:
                                      ├─ function_arn
                                      ├─ function_name
                                      ├─ log_group_name
                                      ├─ role_arn
                                      └─ layer_arns
```

### Critical Relationships

**Code Source Type → Configuration**:
- **S3 Package**: Requires `runtime`, `handler`, `s3.bucket`, `s3.key`; optional `s3.object_version`
- **Container Image**: Requires `image_uri`; runtime/handler defined in Dockerfile and ignored here
- **Package Type**: Automatically set to "Zip" (default) for S3 or "Image" for container

**VPC Configuration → ENI Creation**:
- If subnets or security groups provided, Lambda creates ENIs (Elastic Network Interfaces)
- ENI creation adds 10-30 seconds to cold start time
- Requires IAM permissions: `ec2:CreateNetworkInterface`, `ec2:DescribeNetworkInterfaces`, `ec2:DeleteNetworkInterface`
- Functions with VPC config can access VPC resources but NOT the internet unless routed through NAT

**Layers → Dependency Ordering**:
- Layers are extracted in order (up to 5 layers)
- Layer order matters: later layers can override files from earlier layers
- Layers share the `/opt` directory in the Lambda execution environment

**Role ARN → Execution Permissions**:
- Execution role must allow `lambda.amazonaws.com` to assume it (trust policy)
- Role must have permissions for: CloudWatch Logs, VPC ENI creation (if VPC), and any services the function accesses
- The module exports the role ARN for validation

## Key Design Decisions

### 1. Automatic CloudWatch Log Group Creation

**Decision**: Always create a dedicated CloudWatch log group with a standard name pattern.

**Implementation**:
```go
logGroupName := pulumi.String("/aws/lambda/" + functionName)
logGroup, err = cloudwatch.NewLogGroup(ctx, "log-group", &cloudwatch.LogGroupArgs{
    Name:            logGroupName,
    RetentionInDays: pulumi.Int(30),
    Tags:            pulumi.ToStringMap(locals.AwsTags),
})
```

**Rationale**:
- **Predictable Naming**: Standard AWS pattern `/aws/lambda/{functionName}`
- **Retention Control**: Default 30 days prevents unbounded log storage costs
- **Resource Tagging**: Log group tagged with same metadata as function for cost tracking
- **Explicit Management**: Without explicit creation, Lambda creates the log group automatically but without tags or retention

**Alternative Approach**: Let Lambda auto-create log groups. This module creates them explicitly for better control.

### 2. Code Source Type Branching

**Decision**: Use conditional logic based on `code_source_type` enum to set different FunctionArgs fields.

**Implementation**:
```go
if spec.CodeSourceType == awslambdav1.CodeSourceType_CODE_SOURCE_TYPE_S3 {
    args.Runtime = pulumi.String(spec.Runtime)
    args.Handler = pulumi.String(spec.Handler)
    args.S3Bucket = pulumi.String(spec.S3.Bucket)
    args.S3Key = pulumi.String(spec.S3.Key)
}

if spec.CodeSourceType == awslambdav1.CodeSourceType_CODE_SOURCE_TYPE_IMAGE {
    args.ImageUri = pulumi.String(spec.ImageUri)
    args.PackageType = pulumi.String("Image")
}
```

**Rationale**:
- **API Compatibility**: AWS Lambda API has mutually exclusive field requirements for zip vs image
- **Clear Intent**: Branching makes it obvious which code path is taken
- **Validation Upstream**: Protobuf validation ensures required fields are present for each code source type

**Trade-off**: More verbose than reflection-based mapping, but explicit branching catches API changes at compile time.

### 3. StringValueOrRef Pattern for Cross-Resource References

**Decision**: Fields like `role_arn`, `subnets`, `security_groups`, `layer_arns`, and `kms_key_arn` use `StringValueOrRef` to support both literal values and references to other resources.

**Implementation**:
```go
// Direct value
spec.RoleArn.GetValue() → "arn:aws:iam::123456789012:role/lambda-exec"

// Or reference to another resource
spec.RoleArn → references AwsIamRole.status.outputs.role_arn
```

**Rationale**:
- **Flexibility**: Users can provide ARNs directly or reference other OpenMCF resources
- **Dependency Management**: References create implicit dependencies between resources
- **Simplified Workflow**: No need to manually extract and copy ARNs between resources

**Helper Function**: `valuefrom.ToStringArray()` resolves `StringValueOrRef[]` to plain string arrays.

### 4. Conditional VPC Configuration

**Decision**: Only create VPC config if subnets or security groups are provided.

**Implementation**:
```go
if len(subnetIds) > 0 || len(sgIds) > 0 {
    args.VpcConfig = &awslambda.FunctionVpcConfigArgs{
        SubnetIds:        pulumi.ToStringArray(subnetIds),
        SecurityGroupIds: pulumi.ToStringArray(sgIds),
    }
}
```

**Rationale**:
- **VPC is Optional**: Not all Lambda functions need VPC access
- **Cold Start Impact**: VPC-attached functions have slower cold starts (10-30s) due to ENI creation
- **Zero-Value Handling**: Don't create VPC config with empty arrays (AWS rejects it)

**Best Practice**: Only use VPC when Lambda needs to access VPC resources (RDS, ElastiCache, private APIs).

### 5. Fixed 30-Day Log Retention

**Decision**: Hard-code CloudWatch log retention to 30 days.

**Rationale**:
- **Cost Control**: Lambda generates significant logs; unbounded retention can be expensive
- **Compliance Balance**: 30 days covers most debugging and audit requirements
- **Simplicity**: One less configuration field for users to decide

**Trade-off**: Users needing longer retention must manually update the log group. Future enhancement could add `log_retention_days` to spec.

## Lambda-Specific Implementation Details

### Architecture Enum Mapping

The module maps protobuf architecture enum to AWS strings.

**Enum Definitions** (spec.proto):
```protobuf
enum Architecture {
  ARCHITECTURE_UNSPECIFIED = 0;
  X86_64 = 1;
  ARM64 = 2;
}
```

**Mapping**:
```go
var arch pulumi.StringInput = pulumi.String("x86_64")  // Default
if spec.Architecture == awslambdav1.Architecture_ARM64 {
    arch = pulumi.String("arm64")
}
```

**AWS Strings**:
- `X86_64 (1)` → `"x86_64"` - Intel/AMD processors, widest runtime support
- `ARM64 (2)` → `"arm64"` - AWS Graviton2 processors, 20% cheaper, better price/performance

**Default**: x86_64 for maximum compatibility.

### S3 Code Package Handling

For S3-based deployments, three fields are required: bucket, key, and optionally object version.

**Implementation**:
```go
if spec.S3.Bucket != "" && spec.S3.Key != "" {
    args.S3Bucket = pulumi.String(spec.S3.Bucket)
    args.S3Key = pulumi.String(spec.S3.Key)
    if spec.S3.ObjectVersion != "" {
        args.S3ObjectVersion = pulumi.String(spec.S3.ObjectVersion)
    }
}
```

**S3 Object Version**: Optional but recommended for reproducible deployments. Without it, Lambda uses the latest version, which can change unexpectedly.

### Container Image URI Handling

For container-based deployments, only the image URI is required.

**Implementation**:
```go
if spec.ImageUri != "" {
    args.ImageUri = pulumi.String(spec.ImageUri)
    args.PackageType = pulumi.String("Image")
}
```

**Image URI Format**: Must be an ECR image in the same region or a public ECR image: `123456789012.dkr.ecr.us-east-1.amazonaws.com/my-app:v1.0.0`

**Package Type Override**: Explicitly set to "Image" (default is "Zip").

### VPC Configuration and ENI Lifecycle

When VPC config is provided, Lambda creates ENIs in the specified subnets.

**Critical Requirements**:
1. **Subnet IDs**: At least one subnet (recommended: 2+ for HA across AZs)
2. **Security Groups**: At least one security group with outbound rules
3. **IAM Permissions**: Execution role must have `AWSLambdaVPCAccessExecutionRole` policy

**Implementation**:
```go
subnetIds := valuefrom.ToStringArray(spec.Subnets)
sgIds := valuefrom.ToStringArray(spec.SecurityGroups)

if len(subnetIds) > 0 || len(sgIds) > 0 {
    args.VpcConfig = &awslambda.FunctionVpcConfigArgs{
        SubnetIds:        pulumi.ToStringArray(subnetIds),
        SecurityGroupIds: pulumi.ToStringArray(sgIds),
    }
}
```

**ENI Behavior**:
- Lambda creates one ENI per subnet-security group combination
- ENIs are reused across invocations (Hyperplane architecture)
- First invocation after deployment or update is slow (ENI creation)
- Subsequent invocations within the same hour reuse ENIs (fast)

### Layer ARN Resolution

Layers provide shared code and dependencies without including them in every deployment package.

**Implementation**:
```go
layerArns := valuefrom.ToStringArray(spec.LayerArns)
if len(layerArns) > 0 {
    args.Layers = pulumi.ToStringArray(layerArns)
}
```

**Layer Ordering**: The order in the array matters. If multiple layers contain the same file, later layers override earlier ones.

**Layer Path**: All layers are extracted to `/opt` in the Lambda execution environment. Example: layer with `python/lib/requests/` → available at `/opt/python/lib/requests/`

### KMS Key for Environment Variable Encryption

Lambda automatically encrypts environment variables at rest using AWS-managed keys. Custom KMS keys provide additional control.

**Implementation**:
```go
if spec.KmsKeyArn.GetValue() != "" {
    args.KmsKeyArn = pulumi.String(spec.KmsKeyArn.GetValue())
}
```

**Why Custom KMS Keys?**:
- Cross-account encryption key access
- Audit who decrypts environment variables via CloudTrail
- Key policy controls (restrict decrypt to specific roles)
- Compliance requirements for customer-managed keys

**Cost Impact**: KMS API calls for encrypt/decrypt environment variables (minimal unless function scales to thousands of invocations).

### Memory and Timeout Defaults

The module only sets memory and timeout if explicitly provided (non-zero).

**Implementation**:
```go
if spec.MemoryMb != 0 {
    args.MemorySize = pulumi.Int(int(spec.MemoryMb))
}
if spec.TimeoutSeconds != 0 {
    args.Timeout = pulumi.Int(int(spec.TimeoutSeconds))
}
```

**AWS Defaults**: If omitted, AWS uses 128 MB memory and 3 seconds timeout.

**Rationale**: Allow Pulumi to use AWS defaults rather than hard-coding them. Users explicitly set values only when needed.

## Error Handling Philosophy

The module follows a **fail-fast** approach:

1. **Validation at Protobuf Level**: The `spec.proto` validation rules (including CEL validations) catch configuration errors before Pulumi runs
2. **AWS API Errors Propagate**: If AWS rejects a configuration, the error propagates immediately (no silent fallbacks)
3. **Wrapped Errors**: All errors include context (`errors.Wrap()`) for easier debugging

**Rationale**: Infrastructure as code demands predictability. Silent defaults or error recovery can mask misconfiguration and create surprise behavior.

## Common Pitfalls and Gotchas

### Pitfall 1: VPC Cold Start Latency
**Symptom**: First function invocation after deployment takes 10-30 seconds.

**Cause**: Lambda creating ENIs in VPC subnets.

**Solution**: Expected behavior. Use provisioned concurrency if cold starts are unacceptable, or avoid VPC if the function doesn't need private resource access.

### Pitfall 2: Missing VPC IAM Permissions
**Symptom**: Lambda function fails to create with error about network interfaces.

**Cause**: Execution role lacks VPC-related permissions.

**Solution**: Attach `AWSLambdaVPCAccessExecutionRole` managed policy to the execution role.

### Pitfall 3: Runtime Ignored for Container Images
**Symptom**: Set `runtime: python3.11` but function runs a different runtime.

**Cause**: When `code_source_type: IMAGE`, runtime is defined by the Dockerfile, not the spec.

**Solution**: Remove `runtime` and `handler` fields when using container images, or understand they're ignored.

### Pitfall 4: Layer Limit Exceeded
**Symptom**: AWS API error: `The total unzipped size of the function and all layers cannot exceed 250 MB`

**Cause**: Too many layers or layers with large dependencies.

**Solution**: Consolidate layers, remove unused dependencies, or bundle more code directly in the deployment package.

### Pitfall 5: S3 Bucket Region Mismatch
**Symptom**: Lambda function fails to find S3 code package.

**Cause**: S3 bucket is in a different region than the Lambda function.

**Solution**: S3 bucket and Lambda must be in the same region. Use S3 replication or copy packages to the target region.

### Pitfall 6: Missing IAM Role Permissions
**Symptom**: Lambda function fails at runtime with access denied errors.

**Cause**: Execution role doesn't have permissions for services the function accesses (DynamoDB, S3, etc.).

**Solution**: Update the IAM role to grant required permissions. Use `role_arn` reference to an `AwsIamRole` resource for centralized management.

## Testing and Debugging

### Debugging with Delve

The `debug.sh` script enables step-through debugging:

1. Uncomment the binary option in `Pulumi.yaml`:
   ```yaml
   runtime:
     options:
       binary: ./debug.sh
   ```
2. Run Pulumi CLI commands normally
3. The debug script launches Delve, allowing breakpoints in any module file

**Use Case**: Debugging VPC config resolution, code source branching, or AWS API errors.

### Manual Testing with Sample Manifest

Use the sample manifest in `iac/hack/manifest.yaml` for local testing:

```bash
cd iac/pulumi
pulumi stack init dev
pulumi config set aws:region us-east-1
# Set AWS credentials via environment variables or AWS_PROFILE
pulumi up
```

**Verification**:
```bash
# Get outputs
pulumi stack output function_arn
pulumi stack output function_name
pulumi stack output log_group_name

# Test function invocation
aws lambda invoke \
  --function-name $(pulumi stack output function_name) \
  --payload '{"key":"value"}' \
  response.json

# View logs
aws logs tail /aws/lambda/$(pulumi stack output function_name) --follow
```

## Performance Considerations

### Resource Creation Time

Lambda function creation time varies by configuration:
- **Basic function (no VPC)**: 5-10 seconds
- **Function with VPC**: 30-60 seconds (ENI creation time)
- **Function with many layers**: +2-5 seconds per layer
- **Function with container image**: +5-10 seconds (image pull)

**Total**: Expect 10-60 seconds depending on complexity.

### Cold Start Implications

VPC-attached functions have significantly slower cold starts:
- **No VPC**: 100-1000ms
- **With VPC**: 1-10+ seconds (ENI attachment)

**Optimization Strategies**:
- Use provisioned concurrency for latency-critical functions
- Keep deployment packages small (<10 MB unzipped)
- Minimize layer count and size
- Consider ARM64 architecture (slightly faster cold starts)

### Pulumi State Tracking

The module creates two Pulumi resources:
1. One `cloudwatch.LogGroup` resource
2. One `lambda.Function` resource

**State Size Impact**: Minimal for individual functions. Large-scale deployments (100+ functions) should monitor state file size.

## Lambda-Specific Best Practices

### Memory Sizing Strategy

Start with minimum memory (128 MB) and increase based on:
- **CloudWatch Logs**: Check "Max Memory Used" metric
- **Duration vs Cost**: Higher memory = faster execution = lower duration cost (up to a point)
- **Sweet Spot**: Often 512-1024 MB balances performance and cost

**Power Tuning**: Use AWS Lambda Power Tuning tool to find optimal memory for your workload.

### Timeout Configuration

Set timeout slightly higher than worst-case expected runtime:
- **API handlers**: 3-30 seconds
- **Data processing**: 60-300 seconds
- **Long-running jobs**: 300-900 seconds (max)

**Warning**: Timeout = max billable duration even if function finishes early. Don't set unnecessarily high.

### Environment Variable Encryption

**Without Custom KMS Key**: Environment variables encrypted with AWS-managed key (free, automatic).

**With Custom KMS Key**: Function needs decrypt permission at cold start to read environment variables.

**IAM Policy Required**:
```json
{
  "Effect": "Allow",
  "Action": "kms:Decrypt",
  "Resource": "arn:aws:kms:us-east-1:123456789012:key/abc-123"
}
```

### Reserved Concurrency vs Unreserved Pool

**Omit or set to -1**: Function uses unreserved account concurrency (default 1000, adjustable)

**Set to positive integer**: Reserve dedicated concurrency
- Guarantees availability under load
- Prevents this function from consuming all account concurrency
- Reduces available concurrency for other functions

**Set to 0**: Throttles all invocations (useful for emergency shutoff)

## Cost Optimization

Lambda pricing has three components:
1. **Requests**: $0.20 per 1M requests
2. **Duration**: $0.0000166667 per GB-second (varies by region)
3. **Provisioned Concurrency**: $0.015 per GB-hour (if used)

**Optimization Strategies**:

**Reduce Memory**:
- Benchmark function with different memory sizes
- AWS Lambda Power Tuning tool automates this
- Lower memory = lower cost per invocation

**Reduce Duration**:
- Optimize code (faster execution = lower cost)
- Consider ARM64 architecture (20% cheaper for same performance)
- Remove unused dependencies from deployment package

**Efficient Packaging**:
- Keep zip packages small (<10 MB)
- Use layers for shared dependencies
- Exclude dev dependencies (test frameworks, build tools)

**Right-Size Timeout**:
- Don't set 900s timeout for a 5s function
- Timeout caps maximum cost per invocation

**Monitor with CloudWatch**:
- Use Cost Explorer to track Lambda costs by function (via tags)
- Set up billing alarms for unexpected cost spikes

## Conclusion

This Pulumi module demonstrates OpenMCF's philosophy: **support both simple and complex use cases through a single, flexible API**.

The architecture accommodates:
- **Simple functions**: S3 zip package, no VPC, minimal config
- **Complex functions**: Container images, VPC networking, layers, encrypted environment, custom memory/timeout

Key design principles:
- **Explicit branching** for S3 vs container image code
- **Automatic log group creation** with retention and tagging
- **Conditional VPC configuration** to avoid cold start penalties when unnecessary
- **StringValueOrRef pattern** for cross-resource references
- **Fixed retention period** for cost predictability

For teams familiar with Terraform or CloudFormation, the mapping should feel natural. For teams new to Lambda, the research documentation (`docs/README.md`) provides critical context about serverless patterns, cold starts, and cost optimization.

The result is a production-ready module that provisions Lambda functions with proper IAM roles, VPC networking (when needed), log management, and comprehensive output exports—handling the complexity of Lambda's dual packaging models while keeping the API surface clean.

