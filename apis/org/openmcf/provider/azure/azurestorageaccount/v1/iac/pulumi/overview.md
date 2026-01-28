# Azure Storage Account Pulumi Module Architecture

## Overview

This module provides Infrastructure-as-Code for deploying Azure Storage Accounts using Pulumi with Go. It creates storage accounts with configurable replication, access tiers, network security, and blob containers.

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    AzureStorageAccountStackInput                 │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │ Target: AzureStorageAccount (manifest)                      │ │
│  │ ProviderConfig: AzureProviderConfig (credentials)           │ │
│  └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Pulumi Module                            │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐          │
│  │   locals.go  │  │    main.go   │  │  outputs.go  │          │
│  │  - Name gen  │  │  - Provider  │  │  - Export    │          │
│  │  - Tags      │  │  - Storage   │  │  - Constants │          │
│  │  - Enums     │  │  - Containers│  │              │          │
│  └──────────────┘  └──────────────┘  └──────────────┘          │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Azure Resources                             │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                    Storage Account                          │ │
│  │  - AccountKind: StorageV2 (default)                        │ │
│  │  - AccountTier: Standard/Premium                            │ │
│  │  - Replication: LRS/ZRS/GRS/GZRS/RA-GRS/RA-GZRS           │ │
│  │  - AccessTier: Hot/Cool                                     │ │
│  │  - NetworkRules: IP rules, VNet rules                       │ │
│  │  - BlobProperties: Versioning, Soft Delete                  │ │
│  └─────────────────────────────────────────────────────────────┘ │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                   Blob Containers                           │ │
│  │  - Container 1: data (private)                              │ │
│  │  - Container 2: logs (private)                              │ │
│  │  - Container N: ... (configurable access)                   │ │
│  └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

## Module Structure

```
iac/pulumi/
├── main.go           # Entrypoint - loads stack input, calls module
├── Pulumi.yaml       # Pulumi project configuration
├── Makefile          # Build automation
├── README.md         # Usage documentation
├── overview.md       # This file
└── module/
    ├── main.go       # Resource creation logic
    ├── locals.go     # Local variables, name generation, enum conversion
    └── outputs.go    # Stack output constants
```

## Data Flow

1. **Input Loading**: `main.go` loads the `AzureStorageAccountStackInput` from environment
2. **Locals Initialization**: `locals.go` computes derived values (name, tags, enum conversions)
3. **Provider Setup**: Azure provider configured with service principal credentials
4. **Resource Creation**: Storage account and containers created in sequence
5. **Output Export**: Resource identifiers and endpoints exported as stack outputs

## Key Design Decisions

### Storage Account Naming

Azure Storage Account names must be:
- 3-24 characters long
- Lowercase letters and numbers only
- Globally unique across all Azure subscriptions

The module sanitizes the metadata name by:
1. Removing dots, underscores, and hyphens
2. Converting to lowercase
3. Truncating to 24 characters maximum
4. Ensuring minimum 3 characters

### Network Security Defaults

By default, the module configures:
- Default action: Deny (secure by default)
- Azure Services bypass: Enabled (allows Azure internal services)
- No IP rules (must be explicitly configured)

This follows the principle of least privilege.

### Blob Properties Defaults

When not explicitly configured:
- Versioning: Disabled
- Blob soft delete: 7 days
- Container soft delete: 7 days

These defaults balance data protection with storage costs.

### Container Creation

Containers are created as child resources of the storage account, ensuring:
- Proper dependency ordering
- Cleanup on storage account deletion
- Consistent tagging and lifecycle management

## Extension Points

### Adding New Features

To add support for new Azure Storage features:

1. Update `spec.proto` with new fields
2. Add enum conversion functions in `locals.go`
3. Use new fields in `main.go` resource creation
4. Add corresponding outputs in `outputs.go`
5. Update stack_outputs.proto if exposing new outputs

### Multi-Region Support

The current implementation deploys to a single region. For multi-region:
- Create separate storage accounts per region
- Configure RA-GRS/RA-GZRS for read access in secondary region
- Use Azure Traffic Manager or Front Door for global routing

## Dependencies

- `github.com/pulumi/pulumi-azure/sdk/v6` - Azure Classic provider
- `github.com/pulumi/pulumi/sdk/v3` - Pulumi SDK
- `github.com/pkg/errors` - Error wrapping
- OpenMCF proto stubs
