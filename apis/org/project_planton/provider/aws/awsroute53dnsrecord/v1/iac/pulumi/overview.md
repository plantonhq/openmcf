# AWS Route53 DNS Record - Module Overview

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     Pulumi Module                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │   main.go    │ -> │  module/     │ -> │  AWS Route53 │      │
│  │  (entrypoint)│    │  main.go     │    │   Record     │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
│         │                   │                    │               │
│         v                   v                    v               │
│  ┌──────────────┐    ┌──────────────┐    ┌──────────────┐      │
│  │ Load Stack   │    │ AWS Provider │    │   Exports    │      │
│  │ Input (env)  │    │ (credentials)│    │  (outputs)   │      │
│  └──────────────┘    └──────────────┘    └──────────────┘      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

## Module Structure

```
iac/pulumi/
├── main.go              # Entrypoint - loads input and calls module
├── Pulumi.yaml          # Project configuration
├── Makefile             # Build automation
├── README.md            # Usage documentation
├── overview.md          # This file
└── module/
    ├── main.go          # Resource creation logic
    ├── locals.go        # Local state initialization
    └── outputs.go       # Output key constants
```

## Data Flow

1. **Input Loading**: `main.go` loads the `AwsRoute53DnsRecordStackInput` from the `STACK_INPUT` environment variable (base64-encoded YAML/JSON)

2. **Provider Setup**: Creates AWS provider with credentials from `provider_config` or uses default credentials

3. **Record Creation**: Creates the Route53 record based on spec:
   - Standard record: Uses `values` and `ttl`
   - Alias record: Uses `alias_target`
   - With routing policy: Applies weighted/latency/failover/geolocation

4. **Output Export**: Exports record FQDN, type, and routing information

## Key Design Decisions

### Single Resource Focus

Unlike `AwsRoute53Zone` which can create multiple records inline, this module focuses on a single DNS record. This enables:
- Granular lifecycle management
- Independent scaling and versioning
- Team-level autonomy (app teams manage their records)

### AWS Classic Provider

Uses `pulumi-aws` (classic) instead of `pulumi-aws-native` because:
- More mature and stable for Route53
- Better error messages
- Consistent with other Route53 modules

### Routing Policy as Oneof

The protobuf schema uses `oneof` for routing policies, ensuring only one policy type is active. The Go code handles this with type switches.

## Extension Points

### Adding New Routing Policies

1. Add policy message to `spec.proto`
2. Add case to `applyRoutingPolicy()` in `module/main.go`
3. Update documentation

### Adding New Record Types

Record types are handled via the shared `DnsRecordType` enum. To add support for specialized record type handling:
1. Add specific logic in `module/main.go`
2. Update validation in `spec.proto` if needed

## Dependencies

- `pulumi-aws` v7.x - AWS provider
- `project-planton/pkg/iac/pulumi` - Stack input loading utilities
- Generated protobuf stubs from `awsroute53dnsrecord/v1`
