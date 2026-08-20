# AwsConfigAggregator — Pulumi module

Manages the aggregation arms: the aggregator
(`aws:cfg/configurationAggregator:ConfigurationAggregator`,
presence-gated on `spec.aggregation`) and the reciprocal
source-account grants
(`aws:cfg/aggregateAuthorization:AggregateAuthorization`, one per
`spec.authorizations` entry).

Module facts worth knowing before editing:

- **Grants are keyed by identity.** Each grant resource is named
  `grant-{account_id}:{authorized_aws_region}` — the provider's own
  import ID — so reordering the spec list never churns resources; the
  `authorization_arns` output echoes the same-keyed map.
- **Exactly one source shape renders** (the spec CEL guarantees it);
  empty region lists are omitted so both engines send identical
  payloads.
- **The provider's replacement rule is left alone.** A source block
  APPEARING on an existing aggregator replaces it; content changes
  update in place — provider behavior, nothing encoded here.
- **The deprecated `region` alias is never rendered** —
  `authorized_aws_region` is the surviving argument.

Outputs mirror the Terraform module key-for-key: `aggregator_name`,
`aggregator_arn`, `authorization_arns` (empty string / empty map on
arms the instance does not manage).
