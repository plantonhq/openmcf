# Azure DNS Record Pulumi Module Architecture

## Overview

This module creates individual DNS records in Azure DNS Zones using the Pulumi Azure provider.

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    AzureDnsRecordStackInput                 │
│  ┌───────────────────────┐  ┌────────────────────────────┐  │
│  │   AzureDnsRecord      │  │   AzureProviderConfig      │  │
│  │   - metadata          │  │   - client_id              │  │
│  │   - spec              │  │   - client_secret          │  │
│  │     - resource_group  │  │   - subscription_id        │  │
│  │     - zone_name       │  │   - tenant_id              │  │
│  │     - record_type     │  │                            │  │
│  │     - name            │  │                            │  │
│  │     - values          │  │                            │  │
│  │     - ttl_seconds     │  │                            │  │
│  └───────────────────────┘  └────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Pulumi Module                          │
│  ┌────────────────┐  ┌────────────────┐  ┌──────────────┐  │
│  │   main.go      │  │   locals.go    │  │  outputs.go  │  │
│  │   - Provider   │  │   - Tags       │  │  - OpRecordId│  │
│  │   - Resources  │  │   - TTL        │  │  - OpFqdn    │  │
│  │   - Exports    │  │   - MxPriority │  │              │  │
│  └────────────────┘  └────────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              Azure DNS Resources (conditional)              │
│  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐           │
│  │A Record │ │AAAA Rec │ │CNAME Rec│ │MX Record│ ...       │
│  └─────────┘ └─────────┘ └─────────┘ └─────────┘           │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Stack Outputs                          │
│  ┌────────────────────┐  ┌─────────────────────────────┐   │
│  │   record_id        │  │   fqdn                      │   │
│  │   (ARM resource ID)│  │   (fully qualified domain)  │   │
│  └────────────────────┘  └─────────────────────────────┘   │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

1. **Input Loading**: `stackinput.LoadStackInput()` deserializes the `AzureDnsRecordStackInput` protobuf from the `STACK_INPUT` environment variable

2. **Locals Initialization**: `initializeLocals()` extracts:
   - Zone name (from literal or reference)
   - Resource group
   - Record name
   - TTL with default fallback
   - MX priority with default fallback
   - Azure resource tags

3. **Provider Setup**: Creates Azure provider with credentials from `ProviderConfig`

4. **Record Creation**: Based on `record_type`, creates the appropriate DNS record:
   - A → `dns.NewARecord`
   - AAAA → `dns.NewAaaaRecord`
   - CNAME → `dns.NewCNameRecord`
   - MX → `dns.NewMxRecord`
   - TXT → `dns.NewTxtRecord`
   - NS → `dns.NewNsRecord`
   - SRV → `dns.NewSrvRecord`
   - CAA → `dns.NewCaaRecord`
   - PTR → `dns.NewPtrRecord`

5. **Output Export**: Exports record ID and FQDN for downstream consumption

## Design Decisions

### Single Record Per Deployment

Unlike the `AzureDnsZone` module which supports embedded records, this module creates exactly one record per deployment. This enables:
- Independent lifecycle management
- Fine-grained access control
- Clear resource ownership

### Type Switch vs Polymorphism

The module uses a switch statement on `record_type` rather than interface-based polymorphism. This is simpler and more explicit for the limited set of DNS record types.

### Tag Propagation

All records inherit tags from the resource metadata, enabling:
- Cost allocation
- Resource grouping
- Compliance tracking
