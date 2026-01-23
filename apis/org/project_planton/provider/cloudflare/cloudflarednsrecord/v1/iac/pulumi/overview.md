# CloudflareDnsRecord Pulumi Module Overview

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                    Project Planton CLI                           │
│                                                                  │
│  ┌──────────────────┐    ┌──────────────────────────────────┐  │
│  │   YAML Manifest   │───▶│   CloudflareDnsRecordStackInput   │  │
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
│  │   main.go     │───▶│   module/    │───▶│  Cloudflare API   │  │
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
│                    Cloudflare DNS Zone                           │
│                                                                  │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    DNS Record                             │   │
│   │  Type: A/AAAA/CNAME/MX/TXT/SRV/NS/CAA                   │   │
│   │  Name: www, @, api, etc.                                 │   │
│   │  Value: IP address, hostname, text                       │   │
│   │  Proxied: orange-cloud / grey-cloud                      │   │
│   └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

## Data Flow

1. **Input**: Project Planton CLI loads YAML manifest and converts to protobuf `CloudflareDnsRecordStackInput`
2. **Encoding**: Stack input is base64-encoded and passed via `STACK_INPUT` environment variable
3. **Entrypoint**: `main.go` decodes the stack input and calls `module.Resources()`
4. **Provider Setup**: Module initializes Cloudflare provider with credentials
5. **Resource Creation**: `dns_record.go` creates the `cloudflare_record` resource
6. **Outputs**: Record ID, hostname, type, and proxy status are exported

## Module Components

### main.go (Entrypoint)

- Entry point for Pulumi execution
- Loads `CloudflareDnsRecordStackInput` from environment
- Delegates to `module.Resources()`

### module/main.go (Controller)

- Orchestrates resource creation
- Initializes locals from stack input
- Sets up Cloudflare provider
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
- Converts protobuf enum to Cloudflare API string
- Handles optional fields (TTL, priority, comment)
- Exports outputs for downstream consumers

## Key Design Decisions

### 1. Record Type Enum

The protobuf spec uses an enum for record types to ensure type safety. The module converts these to Cloudflare API strings at runtime.

### 2. Proxied Default

The `proxied` field defaults to `false` (grey-cloud) for safety. Users must explicitly enable proxying.

### 3. TTL Handling

- TTL of 0 or 1 → Automatic (Cloudflare manages)
- TTL 60-86400 → Explicit seconds

### 4. Priority Field

Only populated for MX and SRV records. Other record types ignore this field.

### 5. Comment Support

Optional comment field (max 100 characters) for documentation purposes.

## Error Handling

- Provider setup errors are wrapped with context
- Resource creation errors include the resource type
- All errors propagate to the Pulumi engine for proper reporting

## Testing Strategy

1. **Unit Tests**: Validation tests in `spec_test.go`
2. **Preview**: `make test` runs `pulumi preview` against test manifest
3. **E2E**: Full deployment requires real Cloudflare credentials and zone

## Dependencies

- `github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare` - Cloudflare provider
- `github.com/pulumi/pulumi/sdk/v3/go/pulumi` - Pulumi SDK
- `github.com/pkg/errors` - Error wrapping
