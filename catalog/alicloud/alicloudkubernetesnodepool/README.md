# AliCloudKubernetesNodePool v1

Alibaba Cloud ACK Kubernetes node pool.

## Overview

This component deploys a node pool within an ACK Managed Kubernetes cluster with:

- Flexible instance type selection with multi-type support for availability
- Auto-scaling via cluster auto-scaler integration
- Managed node pool lifecycle (auto-repair, auto-upgrade, vulnerability patching)
- Spot instance support for cost optimization
- System and data disk configuration with ESSD and encryption support
- Kubernetes labels and taints for workload scheduling

A node pool is the unit of worker node management in ACK. Each pool has its own
instance types, scaling policy, and node configuration. Multiple pools can be
attached to a single cluster to support heterogeneous workloads.

## Directory Structure

```
alicloudkubernetesnodepool/
├── README.md              # This file
├── catalog.md             # User-facing documentation
├── logo.svg               # Component logo
├── presets/               # Ready-to-use configuration manifests
├── e2e/                   # Validated example manifest + live scenarios
├── iac/
│   ├── pulumi/            # Pulumi Go module
│   │   ├── main.go        # Entrypoint
│   │   └── module/        # Module implementation
│   └── tf/                # Terraform HCL module
└── v1alpha1/              # The versioned contract
    ├── api.proto          # KRM resource definition (apiVersion/kind/metadata/spec/status)
    ├── spec.proto         # Specification (all configurable fields)
    ├── input.proto        # IaC module input
    ├── outputs.proto      # IaC module outputs
    ├── spec_test.go       # Validation tests
    └── reference.md       # Generated field reference
```

## Quick Start

See e2e/manifest.yaml for YAML manifest examples.

## Development

```bash
# Build and test
go build ./...
go test -v ./...
go vet ./...

# Terraform validation
cd iac/tf/
terraform init
terraform validate
```

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
