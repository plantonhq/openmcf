# AwsConfigRule — Terraform/OpenTofu module

Manages one Config rule as one of four provider resources —
`aws_config_config_rule` (account scope) or the
`aws_config_organization_managed_rule` /
`aws_config_organization_custom_rule` /
`aws_config_organization_custom_policy_rule` trio (organization
scope) — plus the optional `aws_config_remediation_configuration`.

Module facts worth knowing before editing:

- **The spec CELs pre-validate the branching.** Exactly one source
  arm arrives, org-only and account-only surfaces never mix, and
  remediation only accompanies account rules — the module branches on
  presence without re-validating.
- **The source owner is derived** (managed=AWS,
  custom_lambda=CUSTOM_LAMBDA, custom_policy=CUSTOM_POLICY) and
  always sent; Guard sources also always send their
  ConfigurationItemChangeNotification trigger detail (AWS requires
  it), identically to the Pulumi module.
- **Remediation references the rule RESOURCE** (not the variable), so
  create and destroy ordering ride the graph.
- **Only the account-scoped rule is taggable** — the org rule
  resources carry no tags upstream.

Outputs mirror the Pulumi module key-for-key: `rule_arn`, `rule_name`,
`rule_id`, `remediation_arn`.
