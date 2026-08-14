# AwsCloudTrail — Pulumi module

Manages a CloudTrail trail (`aws:cloudtrail/trail:Trail`) and the
optional organization delegated-administrator registration
(`aws:cloudtrail/organizationDelegatedAdminAccount:OrganizationDelegatedAdminAccount`).

Module facts worth knowing before editing:

- **The bucket policy is validated at create.** AWS rejects the trail
  without the `cloudtrail.amazonaws.com` policy on the delivery
  bucket — the policy is the consumer's contract on AwsS3Bucket,
  never rendered here.
- **One selector style renders.** The spec CEL guarantees classic and
  advanced selectors never arrive together; the module renders
  whichever is present.
- **The CloudWatch group ARN is normalized.** AWS expects the `:*`
  suffix form; the module appends it when the referenced value lacks
  it, so both engines send the identical ARN.
- **The delegated-admin registration is presence-gated** on the spec
  field and account-global — it has no structural edge to the trail
  resource.

Outputs mirror the Terraform module key-for-key: `trail_arn`,
`home_region`, `sns_topic_arn`.
