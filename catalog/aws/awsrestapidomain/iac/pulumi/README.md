# AwsRestApiDomain — Pulumi module (Go)

Deploys an API Gateway v1 custom domain (`apigateway.DomainName`) plus
base-path mappings and optional PRIVATE access associations.

Module facts worth knowing before editing:

- **Certificate argument is endpoint-type-specific.** REGIONAL /
  PRIVATE use the regional certificate argument; EDGE uses the
  CloudFront certificate argument. The spec's single `certificate_arn`
  field is wired to the matching SDK argument.
- **Base-path mapping IDs are composite** (`domain/basePath`). The
  empty base path is keyed `(none)` in outputs (AWS's own sentinel).
- **`routing_mode` is the v1 knob** that arbitrates between base-path
  mappings and the v2 routing-rule surface (which stays on
  AwsHttpApiDomain).
- **DNS is not created here.** Alias targets and zone IDs are outputs.

Outputs mirror the Terraform module key-for-key: `domain_name`,
`domain_name_arn`, `domain_name_id`, regional and CloudFront targets
and zone IDs, plus the mapping and access-association maps.
