# S3 Best Practices Pack

This preset creates a small, readable pack of three managed S3 rules —
no public reads, no public writes, versioning on — deployed and scored
as one unit in this account.

## When to Use

- The first conformance pack in an account (small enough to read end
  to end)
- Baseline S3 hygiene scored as one compliance number

## What You Get

- Three AWS-managed rules, created and owned by the pack (prefixed
  with the pack name)
- One conformance score in the Config console instead of three
  scattered rules
- One-step teardown: deleting the pack deletes its rules

## Customize

- Requires a running Config recorder in the region — deploy an
  AwsConfigRecorder first
- Add rules to the template (or switch to `templateS3Uri` for AWS's
  published operational-best-practices templates)
- Set `organizationScope: true` from the management account to score
  every member account
