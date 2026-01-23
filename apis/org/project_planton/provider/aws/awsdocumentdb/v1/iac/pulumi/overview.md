# AwsDocumentDb Pulumi Module Architecture

## Overview

This Pulumi module deploys an AWS DocumentDB cluster with all associated resources including subnet groups, security groups, parameter groups, and cluster instances.

## Module Structure

```
module/
├── main.go              # Controller/orchestrator
├── locals.go            # Local values and label generation
├── outputs.go           # Output constant definitions
├── cluster.go           # DocumentDB cluster resource
├── instances.go         # Cluster instance creation
├── subnet_group.go      # DB subnet group resource
├── security_group.go    # VPC security group resource
└── parameter_group.go   # Cluster parameter group resource
```

## Resource Flow

```
┌─────────────────────────────────────────────────────────────┐
│                        main.go                               │
│                   (Resources function)                       │
└──────────────────────────┬──────────────────────────────────┘
                           │
        ┌──────────────────┼──────────────────┐
        │                  │                  │
        ▼                  ▼                  ▼
┌───────────────┐  ┌───────────────┐  ┌───────────────┐
│ Security Group│  │ Subnet Group  │  │Parameter Group│
│  (optional)   │  │  (optional)   │  │  (optional)   │
└───────┬───────┘  └───────┬───────┘  └───────┬───────┘
        │                  │                  │
        └──────────────────┼──────────────────┘
                           │
                           ▼
                  ┌─────────────────┐
                  │ DocumentDB      │
                  │ Cluster         │
                  └────────┬────────┘
                           │
                           ▼
                  ┌─────────────────┐
                  │ Cluster         │
                  │ Instances       │
                  └─────────────────┘
```

## Key Design Decisions

### 1. Optional Resource Creation

Resources are created conditionally based on spec configuration:

- **Subnet Group**: Created only when `subnetIds` is provided and `dbSubnetGroupName` is not set
- **Security Group**: Created only when `allowedCidrBlocks` or `securityGroupIds` are provided
- **Parameter Group**: Created only when `clusterParameters` are provided

### 2. Credential Handling

AWS credentials are passed through `AwsDocumentDbStackInput.ProviderConfig`, not embedded in the spec. This allows for:
- Secure credential injection at deployment time
- Different credentials per environment
- IAM role assumption patterns

### 3. Instance Scaling

The module creates the specified number of instances (`instanceCount`) with identical configuration. All instances can serve read traffic, with one designated as the primary writer.

### 4. Resource Naming

Resources use the `metadata.id` from the manifest as the identifier, ensuring consistent naming across deployments.

## Outputs

The module exports the following outputs (matching `AwsDocumentDbStackOutputs`):

| Output | Description |
|--------|-------------|
| `cluster_endpoint` | Primary writer endpoint |
| `cluster_reader_endpoint` | Reader endpoint for read replicas |
| `cluster_id` | Cluster identifier |
| `cluster_arn` | Cluster ARN |
| `cluster_port` | Connection port |
| `db_subnet_group_name` | Subnet group name (if created) |
| `security_group_id` | Security group ID (if created) |
| `cluster_parameter_group_name` | Parameter group name (if created) |
| `connection_string` | MongoDB-compatible connection string template |
| `cluster_resource_id` | Internal AWS resource ID |

## Dependencies

- `github.com/pulumi/pulumi-aws/sdk/v7/go/aws` - AWS provider
- `github.com/pulumi/pulumi/sdk/v3/go/pulumi` - Pulumi SDK
- `github.com/pkg/errors` - Error wrapping
