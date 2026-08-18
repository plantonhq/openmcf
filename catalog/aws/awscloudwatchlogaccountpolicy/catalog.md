# AWS CloudWatch Logs Account Policy

Set a logging rule once, account-wide: mask PII in every log group, forward everything to one destination, index the fields your queries filter on, or transform events at ingest — without touching each log group.

## What Gets Managed

- One account-level policy per (name, type) pair: its type (data protection, subscription filter, field index, transformer, or metric extraction), its type-specific document, and the selection criteria narrowing which log groups it covers.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with `logs:PutAccountPolicy` / `logs:DescribeAccountPolicies` / `logs:DeleteAccountPolicy`.

### Blast Radius

Account policies apply to EVERY log group in the region unless `selection_criteria` narrows them. Start narrowed (a prefix selection), widen deliberately.

## After You Deploy

- The policy applies to matching log groups immediately (and to ones created later).
- AWS allows a bounded number of policies per type per account (one for most types) — plan names accordingly.

## Common Changes

- Document edits apply in place; name, type, and selection criteria replace the policy.
