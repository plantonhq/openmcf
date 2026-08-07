# AliCloudApplicationLoadBalancer v1

Alibaba Cloud Application Load Balancer (ALB) with bundled server groups and listeners.

## Overview

This component deploys a Layer 7 Application Load Balancer on Alibaba Cloud with:

- Multi-AZ deployment (minimum 2 zones) for high availability
- Server groups with health check and session stickiness configuration
- HTTP, HTTPS, and QUIC listeners with configurable timeouts

Server groups are created as empty targets -- backend membership is managed externally.

## Directory Structure

```
alicloudapplicationloadbalancer/
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

See [catalog.md](./catalog.md) for usage instructions and examples.

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
