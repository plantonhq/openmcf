# OpenFGA Store - Pulumi Module Overview

> ⚠️ **IMPORTANT**: OpenFGA does not have a Pulumi provider. This module is a **pass-through placeholder**.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                  OpenFgaStore Component                 │
│                                                         │
│  ┌───────────────┐           ┌───────────────────────┐  │
│  │    Pulumi     │           │      Terraform        │  │
│  │  (Pass-thru)  │           │   (Real Impl)         │  │
│  │               │           │                       │  │
│  │ • Logs warn   │           │ • Creates store       │  │
│  │ • No resources│           │ • Uses openfga_store  │  │
│  │ • Exit 0      │           │ • Exports id/name     │  │
│  └───────────────┘           └───────────────────────┘  │
│         ⚠️                            ✅                 │
└─────────────────────────────────────────────────────────┘
```

## Module Structure

```
iac/pulumi/
├── main.go          # Entry point (loads input, calls module)
├── Pulumi.yaml      # Pulumi project configuration
├── Makefile         # Build/deploy targets
├── debug.sh         # Debug script
├── README.md        # This file
├── overview.md      # Architecture overview
└── module/
    ├── main.go      # Pass-through implementation
    ├── locals.go    # Local value computation
    └── outputs.go   # Output exports
```

## Flow

1. `main.go` loads stack input from context
2. Calls `module.Resources()` with the stack input
3. `module.Resources()` logs warnings and exports notice
4. No actual resources are created

## What Users Should Do

Use the Terraform module instead:

```bash
# Via CLI with --provisioner tofu flag
openmcf apply --manifest openfga-store.yaml \
  --openfga-provider-config openfga-creds.yaml \
  --provisioner tofu
```

## Future Considerations

If a Pulumi provider for OpenFGA becomes available:
1. Update `module/main.go` to create actual resources
2. Update imports to use the new provider SDK
3. Update outputs to export real values
4. Remove pass-through warnings
