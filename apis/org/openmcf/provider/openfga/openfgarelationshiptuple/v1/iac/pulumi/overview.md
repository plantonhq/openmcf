# OpenFGA Relationship Tuple - Pulumi Module Overview

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Pass-Through Module                       │
│                                                             │
│   ┌─────────────┐     ┌──────────────┐     ┌─────────────┐ │
│   │   main.go   │────▶│  module/     │────▶│   Exports   │ │
│   │ (entrypoint)│     │  main.go     │     │   notice    │ │
│   └─────────────┘     │  (pass-thru) │     │   only      │ │
│                       └──────────────┘     └─────────────┘ │
│                                                             │
│   ⚠️ NO RESOURCES CREATED - OpenFGA has no Pulumi provider │
└─────────────────────────────────────────────────────────────┘
```

## Why This Module Exists

OpenFGA only has a Terraform provider. There is no Pulumi provider available.

This module exists as a **placeholder** to:

1. **Maintain Consistency**: All OpenMCF deployment components have both Pulumi and Terraform modules
2. **Provide Clear Errors**: When users accidentally use Pulumi, they get clear guidance
3. **Enable Future Support**: If a Pulumi provider is created, the structure is ready

## Module Components

| File | Purpose |
|------|---------|
| `main.go` | Pulumi entrypoint - loads stack input, calls Resources() |
| `module/main.go` | Resources() - logs warnings, exports notice |
| `module/locals.go` | Placeholder locals struct (unused) |
| `module/outputs.go` | Placeholder outputs function (unused) |
| `Pulumi.yaml` | Pulumi project configuration |
| `Makefile` | Build targets with warnings |

## What Happens When Run

```
$ make up
WARNING: OpenFGA does not have a Pulumi provider.
This module is a pass-through placeholder and will not create resources.
Use: openmcf apply --manifest manifest.yaml --provisioner tofu

Outputs:
    notice: "OpenFGA Relationship Tuple was not created. No Pulumi provider available. Use Terraform/Tofu provisioner."
```

## Correct Usage

Use Terraform/Tofu instead:

```bash
openmcf apply --manifest relationship-tuple.yaml \
  --openfga-provider-config openfga-creds.yaml \
  --provisioner tofu
```

## References

- [terraform-provider-openfga](https://github.com/openfga/terraform-provider-openfga)
- [OpenFGA Documentation](https://openfga.dev/docs)
