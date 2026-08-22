# DigitalOcean Function -- Pulumi Module

Deploys a `digitalocean:index/app:App` with a single functions component from a `DigitalOceanFunction` spec. There is no standalone Functions resource.

The app name is `metadata.name`. The functions component name is `spec.functionName`. Source (`git` / `github` / `gitlab` / `bitbucket`) is actually set on the component — a Functions deploy with an empty source cannot build.

Runtime, memory, timeout, and schedules are read by App Platform from `project.yml` inside `sourceDirectory`. They are not Pulumi args.

## Prerequisites

- Pulumi CLI 3.x
- Go 1.21+
- `DIGITALOCEAN_TOKEN`

## Outputs

| Output | Description |
|--------|-------------|
| `function_id` | App UUID that hosts the functions component |
| `https_endpoint` | Public HTTPS URL |
| `default_hostname` | Default `ondigitalocean.app` hostname |

## Usage

```go
package main

import (
    digitaloceanfunctionv1alpha1 "github.com/plantonhq/planton/catalog/digitalocean/digitaloceanfunction/v1alpha1"
    "github.com/plantonhq/planton/catalog/digitalocean/digitaloceanfunction/iac/pulumi/module"
    "github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
    pulumi.Run(func(ctx *pulumi.Context) error {
        return module.Resources(ctx, stackInput)
    })
}
```

The Planton runner supplies `stackInput`. See the kind [README](../../README.md) and [GUIDE](../../GUIDE.md).
