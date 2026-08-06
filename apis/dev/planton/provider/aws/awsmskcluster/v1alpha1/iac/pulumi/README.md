# AwsMskCluster — Pulumi IaC Module

Pulumi module for provisioning AWS MSK (Managed Streaming for Apache Kafka) clusters using the Planton `AwsMskClusterSpec`.

## Overview

This module creates:
- An MSK Cluster with configurable brokers, encryption, authentication, connectivity, logging, and monitoring. The referenced `securityGroupIds` (required, ≥1) attach directly to the broker network interfaces; ingress rules live on those first-class security-group nodes, never on a module-managed shadow group.
- An inline MSK Configuration from `serverProperties` (conditional — only when the map is non-empty).
- SCRAM secret associations (one per ARN in `scramSecretArns`) and a cluster policy (from `clusterPolicy`) — folded satellites in `module/satellites.go`.

## Usage

### As a Pulumi program

The module is designed to be invoked from the entry point in `main.go`, which loads an `AwsMskClusterStackInput` and calls `module.Resources()`:

```go
package main

import (
    awsmskclusterv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsmskcluster/v1alpha1"
    "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsmskcluster/v1alpha1/iac/pulumi/module"
    "github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        stackInput := &awsmskclusterv1.AwsMskClusterStackInput{}
        if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
            return err
        }
        return module.Resources(ctx, stackInput)
    })
}
```

### Stack Input

The stack input is an `AwsMskClusterStackInput` protobuf message containing:
- `target` — the `AwsMskCluster` resource (metadata + spec).
- `provider_config` — optional AWS credentials (region, access key, secret key, session token).

### Outputs

The module exports 17 stack outputs (see `module/outputs.go` for keys). Access them via `pulumi stack output`:

```bash
pulumi stack output cluster_arn
pulumi stack output bootstrap_brokers_sasl_iam
pulumi stack output zookeeper_connect_string_tls
```

## File Structure

| File | Purpose |
|------|---------|
| `main.go` | Entry point — loads stack input, runs Pulumi program |
| `module/main.go` | Orchestrator — resource creation flow + output exports |
| `module/locals.go` | Locals initialization (labels, resolved target) |
| `module/configuration.go` | Inline MSK Configuration from server_properties |
| `module/cluster.go` | MSK Cluster resource creation |
| `module/satellites.go` | SCRAM secret associations + cluster policy |
| `module/outputs.go` | Output key constants |

## Prerequisites

- Go 1.21+
- Pulumi CLI v3+
- AWS credentials (ambient or via stack input)
- `pulumi-aws` plugin v7

## Related

- [Spec reference](../../README.md)
- [Module architecture](./overview.md)
- [Examples](../../examples.md)
