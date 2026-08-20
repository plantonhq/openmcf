# Guard Policy Rule

This preset writes a custom compliance check as a
CloudFormation-Guard policy — your own logic with ZERO compute to
operate (the Guard engine runs inside AWS Config).

## When to Use

- Organization-specific checks no managed rule covers (tagging
  standards, naming conventions, property combinations)
- When a custom Lambda evaluator would be overkill for a declarative
  check

## What You Get

- Every recorded S3 bucket evaluated against the Guard policy on
  every configuration change
- No function to deploy, patch, or pay for

## Customize

- Write your check in the
  [Guard DSL](https://docs.aws.amazon.com/cfn-guard/) against the
  recorded configuration item's shape
- Keep `scope` matched to the types the policy addresses — unscoped
  Guard rules evaluate everything the recorder captures
- Guard rules trigger on configuration changes (never on a schedule)

## Composing

Requires a running recorder in the region whose recording group
covers the scoped types.
