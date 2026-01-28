# AWS Route53 DNS Record Pulumi Module Architecture

## Module Structure

```
pulumi/
├── main.go          # Entrypoint - initializes stack and calls Resources
├── Pulumi.yaml      # Project configuration
├── Makefile         # Build automation
├── README.md        # Usage documentation
├── overview.md      # This file
└── module/
    ├── main.go      # Core resource creation logic
    ├── locals.go    # Local variable initialization
    └── outputs.go   # Output constant definitions
```

## Resource Flow

```
AwsRoute53DnsRecordStackInput
         │
         ▼
   initializeLocals()
         │
         ▼
    Extract zone_id from StringValueOrRef
         │
         ├──────────────────────────────────┐
         │                                  │
         ▼                                  ▼
   Standard Record                    Alias Record
   (values + ttl)              (dns_name + zone_id from StringValueOrRef)
         │                                  │
         └──────────┬───────────────────────┘
                    │
                    ▼
         Apply Routing Policy
         (weighted/latency/failover/geo)
                    │
                    ▼
      aws_route53_record resource
                    │
                    ▼
            Export outputs
```

## Key Implementation Details

### StringValueOrRef Handling

The `StringValueOrRef` type allows fields to be either literal values or references to other resources:

```go
// Extract zone_id from StringValueOrRef
zoneId := ""
if spec.ZoneId != nil {
    zoneId = spec.ZoneId.GetValue()
}
```

The CLI resolves `value_from` references before invoking Pulumi, so the module always receives the resolved value via `GetValue()`.

### Alias Detection

Alias records are detected by checking if `alias_target.dns_name` has a non-empty value:

```go
isAlias := spec.AliasTarget != nil &&
    spec.AliasTarget.DnsName != nil &&
    spec.AliasTarget.DnsName.GetValue() != ""
```

### Routing Policy Application

Routing policies are applied via a type switch on the oneof policy field:

```go
switch p := policy.Policy.(type) {
case *AwsRoute53RoutingPolicy_Weighted:
    // Apply weighted routing
case *AwsRoute53RoutingPolicy_Latency:
    // Apply latency routing
case *AwsRoute53RoutingPolicy_Failover:
    // Apply failover routing
case *AwsRoute53RoutingPolicy_Geolocation:
    // Apply geolocation routing
}
```

## Outputs

| Constant | Value | Description |
|----------|-------|-------------|
| `OpFqdn` | `fqdn` | Fully qualified domain name |
| `OpRecordType` | `record_type` | DNS record type |
| `OpZoneId` | `zone_id` | Route53 zone ID |
| `OpIsAlias` | `is_alias` | Boolean - is alias record |
| `OpSetIdentifier` | `set_identifier` | Routing set identifier |
