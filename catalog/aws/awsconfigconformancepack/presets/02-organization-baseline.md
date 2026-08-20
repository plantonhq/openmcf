# Organization Security Baseline

This preset deploys one S3-hosted ruleset to EVERY account in the AWS
Organization (minus named exclusions) from the management or
delegated-admin account — governance as a single apply.

## When to Use

- Org-wide security baselines (CIS/PCI-style rulesets) that every
  member account must score against
- Exempting sandbox accounts explicitly instead of silently

## What You Get

- The pack deployed into each member account, scored per account in
  the aggregated Config view
- New member accounts picked up automatically
- One template of record in S3, versioned like code

## Customize

- Organization packs accept exactly ONE template form — this preset
  uses `templateS3Uri`; swap to an inline `templateBody` if you prefer
  (never both)
- Delivery buckets at organization scope must be named with the
  `awsconfigconforms` prefix
- Every member account needs a running Config recorder for its rules
  to evaluate
