# AwsConfigRule — Pulumi module

Manages one Config rule as one of four provider resources —
`aws:cfg/rule:Rule` (account scope) or the
`aws:cfg/organizationManagedRule:OrganizationManagedRule` /
`aws:cfg/organizationCustomRule:OrganizationCustomRule` /
`aws:cfg/organizationCustomPolicyRule:OrganizationCustomPolicyRule`
trio (organization scope) — plus the optional
`aws:cfg/remediationConfiguration:RemediationConfiguration`.

Module facts worth knowing before editing:

- **The spec CELs pre-validate the branching.** Exactly one source
  arm arrives, org-only and account-only surfaces never mix, and
  remediation only accompanies account rules — the module branches on
  presence without re-validating.
- **The source owner is derived** (managed=AWS,
  custom_lambda=CUSTOM_LAMBDA, custom_policy=CUSTOM_POLICY) and
  always sent; Guard sources also always send their
  ConfigurationItemChangeNotification trigger detail (AWS requires
  it), identically to the Terraform module.
- **Remediation consumes the rule's Name OUTPUT**, so create and
  destroy ordering ride the dependency graph.
- **Only the account-scoped rule is taggable** — the org rule
  resources carry no tags upstream.

Outputs mirror the Terraform module key-for-key: `rule_arn`,
`rule_name`, `rule_id`, `remediation_arn`.
