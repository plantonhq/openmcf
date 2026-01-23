# DigitalOcean DNS Record - Module Architecture

## Overview

This Pulumi module creates a single DNS record in DigitalOcean's DNS service.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    Stack Input (JSON)                        │
│  - target: DigitalOceanDnsRecord manifest                   │
│  - providerConfig: DigitalOcean credentials                 │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    main.go (Entrypoint)                      │
│  1. Load stack input from environment                        │
│  2. Call module.Resources()                                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                 module/main.go (Resources)                   │
│  1. Initialize locals                                        │
│  2. Create DigitalOcean provider                            │
│  3. Create DNS record                                        │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│              module/dns_record.go (dnsRecord)               │
│  1. Build DnsRecordArgs from spec                           │
│  2. Add type-specific fields (MX, SRV, CAA)                 │
│  3. Create digitalocean.DnsRecord                           │
│  4. Export stack outputs                                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                  DigitalOcean API                            │
│  - Create DNS record in specified domain                     │
└─────────────────────────────────────────────────────────────┘
```

## Module Files

| File | Purpose |
|------|---------|
| `main.go` | Pulumi entrypoint, loads stack input |
| `module/main.go` | Resources function, orchestrates creation |
| `module/locals.go` | Initializes local values from input |
| `module/outputs.go` | Defines output constants |
| `module/dns_record.go` | DNS record creation logic |

## Key Design Decisions

### 1. Single Record Per Manifest

Unlike the DNS Zone component which can create multiple records, this component creates exactly one DNS record per manifest. This provides:
- Simpler lifecycle management
- Clearer resource boundaries
- Easier troubleshooting

### 2. Type-Specific Field Handling

Record types have different required fields:
- **MX**: priority
- **SRV**: priority, weight, port
- **CAA**: flags, tag

The module conditionally adds these fields based on record type.

### 3. TTL Default

If `ttl_seconds` is not specified, defaults to 1800 (30 minutes), matching DigitalOcean's default behavior.

### 4. Output Construction

The `hostname` output is constructed from the record name:
- Root records (`@`): Returns just the domain
- Subdomains: Returns `{name}.{domain}`

## Dependencies

- `pulumi/pulumi-digitalocean/sdk/v4`: DigitalOcean Pulumi provider
- `pkg/iac/pulumi/pulumimodule`: Shared Pulumi utilities
