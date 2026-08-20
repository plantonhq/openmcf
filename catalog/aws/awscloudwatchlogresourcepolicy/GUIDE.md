# AwsCloudwatchLogResourcePolicy — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## Ten per region — consolidate

The account-scope quota (10) fills fast when every stack ships its own grant. Prefer one policy per service class (`route53-query-logging`, `eventbridge-delivery`) with Resource patterns wide enough for the class, owned by one instance.

## The revision guard is your friend

Updates send the last-seen revision; if someone edited the policy in the console since, the apply fails with a revision mismatch instead of overwriting their change. Re-plan (refreshing state) and apply again to take ownership of the merged truth.

## Resource scope needs its revision to delete

AWS refuses a resource-scoped delete without the current revision ID — which lives in state. Never hand-clear state for this kind; a lost revision means deleting the policy via CLI (`aws logs delete-resource-policy`) before the next apply.

## Scope choice

Account scope (name) is right for service grants that cover a path pattern. Resource scope (one group's ARN) is for a grant that must live and die with a single group — rare, and it makes the group's ARN the policy's identity.
