# Deny-Leave SCP

This preset creates the foundational guardrail: no member account can
remove itself from the organization, attached at the root so it binds
everything.

## When to Use

- Every governed organization, on day one — without it any member
  account's root user can walk out of all guardrails with one API call
- The first SCP after enabling SERVICE_CONTROL_POLICY on the
  organization

## What You Get

- A deny on `organizations:LeaveOrganization` for the whole tree (SCPs
  never bind the management account — it can still remove accounts
  deliberately)
- A safe first SCP: it forbids one administrative action and cannot
  break workloads

## Customize

- Wire `targetId` to your organization's `root_id` output (a
  `valueFrom` reference to the AwsOrganization resource) instead of a
  literal
- Common companions in the same document: deny `account:CloseAccount`,
  deny disabling CloudTrail (`cloudtrail:StopLogging`,
  `cloudtrail:DeleteTrail`)
- Requires SERVICE_CONTROL_POLICY in the organization's
  `enabledPolicyTypes` — attachments fail (policy creation succeeds)
  until it is
