# GcpDnsRecord Pulumi Module Architecture

## Overview

This document describes the architecture and design decisions of the GcpDnsRecord Pulumi module.

## Module Structure

```
pulumi/
├── main.go              # Entry point - loads stack input and calls module
├── Pulumi.yaml          # Pulumi project configuration
├── Makefile             # Build automation
├── README.md            # Usage documentation
├── debug.sh             # Local debugging script
├── overview.md          # This file
└── module/
    ├── main.go          # Resource creation logic
    ├── locals.go        # Input transformation
    └── outputs.go       # Output constants
```

## Data Flow

```
┌─────────────────┐     ┌──────────────┐     ┌─────────────────┐
│ YAML Manifest   │────▶│  Stack Input │────▶│    Locals       │
│ (GcpDnsRecord)  │     │  (protobuf)  │     │ (extracted)     │
└─────────────────┘     └──────────────┘     └────────┬────────┘
                                                      │
                                                      ▼
┌─────────────────┐     ┌──────────────┐     ┌─────────────────┐
│    Outputs      │◀────│  DNS Record  │◀────│  GCP Provider   │
│ (exported)      │     │  (created)   │     │ (configured)    │
└─────────────────┘     └──────────────┘     └─────────────────┘
```

## Component Responsibilities

### main.go (entrypoint)

- Entry point for Pulumi runtime
- Loads `GcpDnsRecordStackInput` from environment
- Invokes module.Resources()

### module/main.go

- Creates GCP provider with credentials
- Creates the `dns.RecordSet` resource
- Exports outputs for the stack

### module/locals.go

- Transforms proto messages into usable values
- Extracts project ID from StringValueOrRef
- Handles default values (TTL)

### module/outputs.go

- Defines output key constants
- Maps to `stack_outputs.proto` fields

## Resource Creation

The module creates a single GCP resource:

### google_dns_record_set

```
┌─────────────────────────────────────┐
│       google_dns_record_set         │
├─────────────────────────────────────┤
│ project      = spec.projectId       │
│ managed_zone = spec.managedZone     │
│ name         = spec.name            │
│ type         = spec.recordType      │
│ ttl          = spec.ttlSeconds      │
│ rrdatas      = spec.values          │
└─────────────────────────────────────┘
```

## Error Handling

Errors are wrapped with context using `pkg/errors`:

```go
if err != nil {
    return errors.Wrap(err, "failed to create DNS record")
}
```

This provides:
- Stack traces for debugging
- Context about where errors occurred
- Clear error messages for users

## Provider Configuration

The GCP provider is configured using credentials from `ProviderConfig`:

```go
gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
```

This abstracts credential handling and supports:
- Service account JSON keys
- Application Default Credentials
- Workload Identity

## Design Decisions

### 1. Single Resource Focus

Unlike GcpDnsZone which can create multiple records inline, this component manages exactly one record set. This enables:
- Independent lifecycle management
- Fine-grained access control
- Team-specific record ownership

### 2. No Zone Creation

The module assumes the managed zone already exists. This separation of concerns means:
- Records can be added to zones managed elsewhere
- No risk of accidental zone deletion
- Simpler resource management

### 3. Full Record Type Support

All DNS record types are supported through the shared `DnsRecordType` enum, ensuring consistency with other DNS components.

## Testing

To test locally:

1. Set up credentials:
```bash
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/sa.json
```

2. Create a test stack:
```bash
pulumi stack init test
```

3. Run with manifest:
```bash
export PLANTON_CLOUD_RESOURCE_MANIFEST=$(cat ../hack/manifest.yaml)
pulumi up
```
