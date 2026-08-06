# Third-Party PrivateLink Service

This preset connects a VPC privately to a vendor's PrivateLink service
(a SaaS database, observability platform, or partner API published
behind their NLB) -- the vendor's endpoint-service name goes in, and
traffic to them never leaves the AWS network.

## When to Use

- Consuming a SaaS vendor's PrivateLink offering (most data and
  observability vendors publish one)
- Cross-account service access inside one organization without VPC
  peering or Transit Gateway
- Replacing IP allowlists with a private, identity-scoped connection

## Key Configuration Choices

- **The vendor's `vpce-svc-*` name** -- from their onboarding docs.
  After deploy, the `state` output reads `pendingAcceptance` until the
  vendor accepts the connection (same-account services can set
  `autoAccept: true` instead).
- **No `privateDnsEnabled`** -- third-party services have no AWS public
  DNS name to override. Clients use the endpoint's own name from the
  `dns_name` output, or a private Route53 alias you create onto it
  (`dns_name` + `hosted_zone_id` outputs are the alias target pair).
  Vendors that verified a private DNS domain are the exception -- flip
  it on if their docs say so.
- **Dedicated security group** -- allow the vendor's service port from
  the clients' CIDR, nothing else.

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<endpoint-resource-name>` | Name for this endpoint resource (e.g. `vendor-privatelink`) | Your naming convention |
| `<aws-region>` | AWS region code (e.g., `us-west-2`) -- must match the vendor service's region | Your deployment region |
| `<vpc-resource-name>` | Name of the AwsVpc resource | Your VPC manifest's `metadata.name` |
| `<private-subnet-a-resource-name>` / `<private-subnet-b-resource-name>` | Names of the AwsSubnet resources (different AZs) | Your subnet manifests' `metadata.name` |
| `<endpoint-sg-resource-name>` | Name of the AwsSecurityGroup for the endpoint ENIs | Your security-group manifest's `metadata.name` |

## Common Additions

- An `AwsRoute53DnsRecord` alias onto the `dns_name` /
  `hosted_zone_id` outputs for a friendly internal hostname
- `serviceRegion` when the vendor publishes in another region
  (cross-region interface endpoints)

## Related Presets

- **01-s3-gateway** -- the free route-based endpoint for S3
- **02-interface-endpoint** -- private access to AWS's own services
