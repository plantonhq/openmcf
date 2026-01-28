# Infrastructure as Code (IaC) Overview

This directory contains the Infrastructure as Code (IaC) implementations for the **GcpCertManagerCert** resource. The implementations are available in both Pulumi (Go) and Terraform, allowing you to choose the IaC tool that best fits your workflow.

## Directory Structure

```
iac/
├── overview.md          # This file
├── pulumi/              # Pulumi implementation (Go)
│   ├── main.go          # Pulumi entry point
│   ├── module/          # Core Pulumi logic
│   │   ├── main.go
│   │   ├── cert_manager_cert.go
│   │   ├── locals.go
│   │   └── outputs.go
│   ├── Pulumi.yaml      # Pulumi project configuration
│   ├── Makefile         # Helper commands
│   ├── debug.sh         # Debug script
│   ├── README.md        # Pulumi-specific documentation
│   └── examples.md      # Pulumi usage examples
└── tf/                  # Terraform implementation
    ├── main.tf          # Certificate resources
    ├── variables.tf     # Input variables
    ├── outputs.tf       # Output values
    ├── locals.tf        # Local variables
    ├── provider.tf      # Provider configuration
    ├── README.md        # Terraform-specific documentation
    └── examples.md      # Terraform usage examples
```

## Supported Backends

### Pulumi (Go)

The Pulumi implementation uses the Go programming language and the official GCP Pulumi provider. It provides:

- **Type Safety**: Strong typing through Go's type system
- **Programmatic Control**: Full programming capabilities for complex logic
- **State Management**: Pulumi handles state management automatically
- **Rich SDK**: Access to the full Pulumi ecosystem

See [pulumi/README.md](pulumi/README.md) for detailed Pulumi-specific documentation.

### Terraform

The Terraform implementation uses HCL (HashiCorp Configuration Language) and the official GCP Terraform provider. It provides:

- **Declarative Syntax**: Clear, readable infrastructure definitions
- **Wide Adoption**: Industry-standard IaC tool
- **Module System**: Reusable infrastructure components
- **State Management**: Robust state handling

See [tf/README.md](tf/README.md) for detailed Terraform-specific documentation.

## Resource Types

Both implementations support creating two types of GCP certificates:

### Certificate Manager (MANAGED)

- Uses Google Certificate Manager
- Modern, feature-rich certificate management
- Supports certificate maps and advanced features
- Automatic DNS authorization and validation
- Works with various GCP services

### Load Balancer Certificates (LOAD_BALANCER)

- Uses Google-managed SSL certificates
- Optimized for GCP load balancers
- Classic certificate management approach
- Automatic provisioning when attached to load balancers

## How It Works

Both implementations follow the same general flow:

1. **Certificate Creation**:
   - For MANAGED type: Create Certificate Manager certificate with DNS authorizations
   - For LOAD_BALANCER type: Create Google-managed SSL certificate

2. **DNS Validation**:
   - For MANAGED type: Create DNS authorization records in Cloud DNS
   - For LOAD_BALANCER type: Validation happens automatically when attached to LB

3. **Domain Verification**:
   - GCP verifies domain ownership through DNS records
   - Certificate becomes active after successful validation

4. **Output Values**:
   - Certificate ID
   - Certificate name
   - Primary domain name
   - Certificate status

## Integration with OpenMCF CLI

While you can use these implementations directly, they are designed to be used through the OpenMCF CLI, which:

- Validates your manifest against the Protobuf schema
- Converts YAML manifests to stack inputs
- Manages provider credentials
- Handles deployment lifecycle

## Choosing Between Pulumi and Terraform

**Use Pulumi if:**
- You prefer programming languages over declarative syntax
- You need complex logic or conditionals
- You want strong type checking
- You're comfortable with Go

**Use Terraform if:**
- You prefer declarative configuration
- Your team is already using Terraform
- You want a widely adopted standard
- You value extensive community modules

Both implementations produce identical infrastructure, so choose based on your preferences and existing tooling.

## Next Steps

- For Pulumi: See [pulumi/README.md](pulumi/README.md)
- For Terraform: See [tf/README.md](tf/README.md)
- For examples: See the main [examples.md](../examples.md) or implementation-specific examples

