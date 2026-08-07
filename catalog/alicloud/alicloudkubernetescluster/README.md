# AliCloudKubernetesCluster v1

Alibaba Cloud ACK Managed Kubernetes cluster.

## Overview

This component deploys an ACK (Alibaba Cloud Container Service for Kubernetes) managed cluster with:

- Fully managed control plane with standard or professional SLA
- Flannel or Terway CNI support for pod networking
- RRSA (RAM Roles for Service Accounts) for pod-level IAM
- Control plane and audit logging to SLS
- Maintenance windows and automatic version upgrades
- Kubernetes Secrets encryption via KMS

Worker nodes are managed separately through AliCloudKubernetesNodePool components.

## Directory Structure

```
alicloudkubernetescluster/
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
