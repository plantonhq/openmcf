# CivoDnsRecord Pulumi Module Overview

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Project Planton CLI                           │
│                                                                  │
│  ┌──────────────────┐    ┌──────────────────────────────────┐  │
│  │   YAML Manifest   │───▶│   CivoDnsRecordStackInput         │  │
│  └──────────────────┘    │   (protobuf)                       │  │
│                          └──────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    │ base64 encoded
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Pulumi Module                               │
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────────┐  │
│  │   main.go     │───▶│   module/    │───▶│     Civo API      │  │
│  │  (entrypoint) │    │  main.go     │    │                   │  │
│  └──────────────┘    └──────────────┘    └──────────────────┘  │
│                              │                                   │
│                              ▼                                   │
│                      ┌──────────────┐                           │
│                      │ dns_record.go │                           │
│                      │  (resource)   │                           │
│                      └──────────────┘                           │
└─────────────────────────────────────────────────────────────────┘
                                    │
                                    ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Civo DNS Zone                               │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    DNS Record                             │   │
│   │  Type: A/AAAA/CNAME/MX/TXT/SRV/NS                       │   │
│   │  Name: www, @, api, etc.                                 │   │
│   │  Value: IP address, hostname, text                       │   │
│   │  TTL: 60-86400 seconds                                   │   │
│   └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

1. **Input**: Project Planton CLI loads YAML manifest and converts to protobuf `CivoDnsRecordStackInput`
2. **Encoding**: Stack input is base64-encoded and passed via `STACK_INPUT` environment variable
3. **Entrypoint**: `main.go` decodes the stack input and calls `module.Resources()`
4. **Provider Setup**: Module initializes Civo provider with credentials
5. **Resource Creation**: `dns_record.go` creates the `civo_dns_domain_record` resource
6. **Outputs**: Record ID, hostname, type, and account ID are exported

## Module Components

### main.go (Entrypoint)

- Entry point for Pulumi execution
- Loads `CivoDnsRecordStackInput` from environment
- Delegates to `module.Resources()`

### module/main.go (Controller)

- Orchestrates resource creation
- Initializes locals from stack input
- Sets up Civo provider
- Calls resource creation functions

### module/locals.go

- Defines `Locals` struct for shared data
- Extracts commonly used values from stack input
- Provides convenient access to spec and metadata

### module/outputs.go

- Defines output constant names
- Used for consistent export key naming

### module/dns_record.go

- Contains the actual resource creation logic
- Converts protobuf enum to Civo API string
- Handles optional fields (TTL, priority)
- Exports outputs for downstream consumers

## Key Design Decisions

### 1. Record Type Enum

The protobuf spec uses an enum for record types to ensure type safety. The module converts these to Civo API strings at runtime.

### 2. TTL Default

The `ttl` field defaults to 3600 (1 hour) if not specified. Valid range is 60-86400 seconds.

### 3. Priority Field

Only populated for MX and SRV records. Other record types ignore this field.

### 4. Zone Reference

Records reference a zone by its ID (`zone_id`), which must be obtained from a CivoDnsZone resource or the Civo dashboard.

## Error Handling

- Provider setup errors are wrapped with context
- Resource creation errors include the resource type
- All errors propagate to the Pulumi engine for proper reporting

## Testing Strategy

1. **Unit Tests**: Validation tests in `spec_test.go`
2. **Preview**: `make test` runs `pulumi preview` against test manifest
3. **E2E**: Full deployment requires real Civo credentials and zone

## Dependencies

- `github.com/pulumi/pulumi-civo/sdk/v2/go/civo` - Civo provider
- `github.com/pulumi/pulumi/sdk/v3/go/pulumi` - Pulumi SDK
- `github.com/pkg/errors` - Error wrapping
