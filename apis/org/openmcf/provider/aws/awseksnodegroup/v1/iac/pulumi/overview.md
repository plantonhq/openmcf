# AWS EKS Node Group Pulumi Module Overview

This document explains the architecture and implementation of the Pulumi module for AWS EKS managed node groups.

## Module Structure

The Pulumi module is organized into the following files:

```
iac/pulumi/
├── main.go              # Entry point for Pulumi program
├── Pulumi.yaml          # Pulumi project configuration
├── Makefile             # Build and deployment automation
├── debug.sh             # Debugging helper script
├── README.md            # User-facing documentation
├── examples.md          # Usage examples
└── module/
    ├── main.go          # Core resource creation logic
    ├── locals.go        # Helper functions and local variables
    └── outputs.go       # Output constant definitions
```

## Architecture Components

### 1. Entry Point (main.go)

The top-level `main.go` serves as the Pulumi program entry point. It:
- Initializes the Pulumi runtime
- Loads the stack input from the OpenMCF manifest
- Calls the module's `Resources()` function
- Handles any errors at the program level

### 2. Resource Module (module/main.go)

The `Resources()` function is the core of the module. It:

1. **Configures the AWS Provider**: 
   - Creates an AWS classic provider with credentials from the stack input
   - Falls back to default credentials if none provided (uses AWS SDK credential chain)

2. **Builds Node Group Configuration**:
   - Extracts the spec from the stack input
   - Calls `buildNodeGroupArgs()` from locals.go to construct the arguments
   - Handles all required and optional fields

3. **Creates the EKS Node Group**:
   - Uses the Pulumi AWS SDK to create an `eks.NodeGroup` resource
   - Applies the configured provider
   - Uses the resource name from the manifest metadata

4. **Exports Outputs**:
   - Exports node group name, ASG name, security group ID, and instance profile ARN
   - Aligns with the `AwsEksNodeGroupStackOutputs` protobuf message

### 3. Helper Functions (module/locals.go)

Contains reusable logic to keep the main resource creation clean:

- **`buildNodeGroupArgs()`**: Constructs the complete `eks.NodeGroupArgs` from the spec
  - Converts subnet IDs from `StringValueOrRef` to Pulumi format
  - Builds scaling configuration with min/max/desired sizes
  - Maps capacity type enum to AWS API format
  - Handles optional SSH key configuration

- **`getCapacityType()`**: Converts protobuf enum to AWS API string
  - Maps `on_demand` → `"ON_DEMAND"`
  - Maps `spot` → `"SPOT"`

- **`getDefaultDiskSize()`**: Returns recommended default disk size (100 GiB)
  - Documented rationale for default choice

### 4. Output Definitions (module/outputs.go)

Defines constant strings for Pulumi outputs:
- `OpNodeGroupName`: The managed node group name
- `OpAsgName`: The underlying Auto Scaling Group name
- `OpRemoteAccessSgId`: Security group ID for SSH access (if enabled)
- `OpInstanceProfileArn`: IAM instance profile ARN

These constants ensure consistent output naming across the module.

## Resource Flow

```
User Manifest (YAML)
       ↓
OpenMCF CLI
       ↓
Stack Input (Protobuf)
       ↓
main.go (Pulumi entry)
       ↓
module.Resources()
       ↓
buildNodeGroupArgs() ← spec.proto fields
       ↓
eks.NewNodeGroup()
       ↓
AWS EKS API
       ↓
Managed Node Group + ASG + EC2 Instances
       ↓
Stack Outputs (Protobuf)
       ↓
User (status.outputs in manifest)
```

## Key Implementation Details

### Provider Configuration

The module supports two provider configuration modes:

1. **Explicit Credentials**: When `ProviderConfig` is provided in stack input
   ```go
   provider, err = aws.NewProvider(ctx, "classic-provider", &aws.ProviderArgs{
       AccessKey: pulumi.String(awsProviderConfig.AccessKeyId),
       SecretKey: pulumi.String(awsProviderConfig.SecretAccessKey),
       Region:    pulumi.String(awsProviderConfig.GetRegion()),
       Token:     pulumi.StringPtr(awsProviderConfig.SessionToken),
   })
   ```

2. **Default Credentials**: Falls back to AWS SDK credential chain
   - Environment variables (AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY)
   - AWS credentials file (~/.aws/credentials)
   - IAM role for EC2 instances (if running on EC2)
   - ECS task role (if running in ECS)

### Foreign Key Resolution

The module handles OpenMCF's `StringValueOrRef` type for dynamic field resolution:

- **Direct Values**: When `value` is set, uses it directly
  ```yaml
  clusterName:
    value: "my-eks-cluster"
  ```

- **Resource References**: When `valueFrom` is set, OpenMCF resolves it before passing to Pulumi
  ```yaml
  clusterName:
    valueFrom:
      kind: AwsEksCluster
      name: "my-cluster"
      fieldPath: "metadata.name"
  ```

The module receives already-resolved values via `GetValue()`, so it doesn't need to handle reference resolution.

### Scaling Configuration

The scaling configuration maps directly to AWS EKS Auto Scaling Group settings:

```go
scaling := &eks.NodeGroupScalingConfigArgs{
    MinSize:     pulumi.Int(int(spec.Scaling.MinSize)),
    MaxSize:     pulumi.Int(int(spec.Scaling.MaxSize)),
    DesiredSize: pulumi.Int(int(spec.Scaling.DesiredSize)),
}
```

These values are validated by protobuf constraints (min >= 1) before reaching Pulumi.

### Optional SSH Access

SSH access is configured only when `sshKeyName` is provided:

```go
if spec.SshKeyName != "" {
    args.RemoteAccess = &eks.NodeGroupRemoteAccessArgs{
        Ec2SshKey: pulumi.String(spec.SshKeyName),
    }
}
```

This creates a security group allowing SSH access from the cluster security group.

### Capacity Type Mapping

The capacity type enum is mapped to AWS API format:

```go
func getCapacityType(capacityType awseksnodegroupv1.AwsEksNodeGroupCapacityType) pulumi.StringInput {
    if capacityType == awseksnodegroupv1.AwsEksNodeGroupCapacityType_spot {
        return pulumi.String("SPOT")
    }
    return pulumi.String("ON_DEMAND")
}
```

AWS expects uppercase with underscore: `ON_DEMAND` or `SPOT`.

### Node Labels

Kubernetes labels are passed through directly:

```go
Labels: pulumi.ToStringMap(spec.Labels)
```

These labels are applied to the node group and automatically propagate to all EC2 instances and their corresponding Kubernetes nodes.

## Validation and Testing

### Protobuf Validation

Field validation happens at the protobuf level before reaching Pulumi:
- Required fields (cluster_name, node_role_arn, subnet_ids, instance_type, scaling)
- Minimum values (min_size >= 1, max_size >= 1, desired_size >= 1)
- Array constraints (subnet_ids must have at least 2 items)
- String constraints (ssh_key_name max 255 characters, label keys/values max 63 characters)

### Component Tests

The `spec_test.go` file validates protobuf constraints:
```bash
go test -v ./apis/org/openmcf/provider/aws/awseksnodegroup/v1/
```

This ensures the protobuf validation rules are correct and comprehensive.

### Integration Testing

The module can be tested using the debug.sh script:

```bash
cd iac/pulumi/
./debug.sh ../hack/manifest.yaml
```

This runs a Pulumi preview without actually creating resources.

## Dependencies

The module depends on:

1. **Pulumi SDKs**:
   - `github.com/pulumi/pulumi/sdk/v3/go/pulumi` - Core Pulumi runtime
   - `github.com/pulumi/pulumi-aws/sdk/v7/go/aws` - AWS provider
   - `github.com/pulumi/pulumi-aws/sdk/v7/go/aws/eks` - EKS resources

2. **OpenMCF**:
   - `github.com/plantonhq/openmcf/apis/.../awseksnodegroup/v1` - Protobuf types

3. **Utilities**:
   - `github.com/pkg/errors` - Error wrapping

## Output Schema

The module exports outputs matching the `AwsEksNodeGroupStackOutputs` protobuf:

```protobuf
message AwsEksNodeGroupStackOutputs {
  string nodegroup_name = 1;
  string asg_name = 2;
  string remote_access_sg_id = 3;
  string instance_profile_arn = 4;
}
```

Currently, only `nodegroup_name` is populated; others are placeholders for future enhancement.

## Error Handling

The module uses error wrapping for context:

```go
return errors.Wrap(err, "failed to create AWS provider with custom credentials")
return errors.Wrap(err, "create EKS node group")
```

This provides a clear error chain when failures occur.

## Best Practices Implemented

1. **Separation of Concerns**: Main resource logic in main.go, helpers in locals.go
2. **Type Safety**: Strong typing via Pulumi SDK and protobuf
3. **Idempotency**: Pulumi handles resource state and ensures idempotent operations
4. **Explicit Dependencies**: Pulumi tracks dependencies automatically via resource references
5. **Minimal Privileges**: No hardcoded credentials; uses provider credential chain
6. **Documentation**: Inline comments explain implementation choices

## Future Enhancements

Potential improvements to consider:

1. **Launch Template Support**: Allow custom launch templates for advanced configuration
2. **Update Strategy**: Expose update configuration (rolling update settings)
3. **Taints**: Support node taints for workload isolation
4. **AMI Selection**: Allow custom AMI specification
5. **Node Group Timeouts**: Expose timeout configuration for creation/update/deletion
6. **Remote Access SG**: Properly export the security group ID for SSH access
7. **Complete ASG Export**: Export the full Auto Scaling Group name and details

## References

- [Pulumi AWS EKS NodeGroup](https://www.pulumi.com/registry/packages/aws/api-docs/eks/nodegroup/)
- [AWS EKS Managed Node Groups](https://docs.aws.amazon.com/eks/latest/userguide/managed-node-groups.html)
- [OpenMCF Architecture](https://github.com/plantonhq/openmcf/blob/main/architecture/deployment-component.md)

