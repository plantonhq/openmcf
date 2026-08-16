# Managed Rule with Remediation

This preset deploys an AWS-authored compliance check with a
human-triggered SSM fix attached — the zero-code path to "detect and
repair".

## When to Use

- Standard compliance checks AWS already maintains (the managed-rule
  catalog covers hundreds)
- When non-compliant resources should be repairable from the Config
  console with one click

## What You Get

- Versioning compliance over every recorded S3 bucket
- The `AWS-ConfigureS3BucketVersioning` SSM document wired as the
  fix, throttled (10% concurrency) and circuit-broken (stop at 50%
  failures)

## Customize

- Swap `ruleIdentifier` for any managed rule; keep `scope` matched to
  what the rule evaluates
- `<remediation-role-arn>`: a role SSM Automation can assume with
  permissions for the document's actions
- Flip `automatic: true` ONLY after the manual fix has proven itself —
  AWS then requires the retry contract the preset already carries

## Composing

Requires a running recorder in the region (the scoped-recording
preset of AwsConfigRecorder records `AWS::S3::Bucket` already).
