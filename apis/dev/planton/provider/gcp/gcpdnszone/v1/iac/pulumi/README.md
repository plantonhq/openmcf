# Pulumi Module: GcpDnsZone

Go Pulumi program provisioning `gcp.dns.ManagedZone` with API enablement.

## Resources

- `projects.Service` — enables `dns.googleapis.com`
- `dns.ManagedZone` — the managed zone

## Outputs

Exported as stack outputs: `zone_id`, `zone_name`, `nameservers`.

## Notes

- Zone-only module; records belong in GcpDnsRecord
- No IAM bindings or record sets created here
