# Workload Account

This preset creates the standard production member account: placed in
a workloads OU, root-only billing access, and the billing/security
contacts AWS should actually call.

## When to Use

- Every long-lived workload account in a governed organization
- Estates that route AWS's billing and security notifications to the
  right teams instead of the root inbox

## What You Get

- A member account inside the referenced OU (guardrails attached
  there govern it from birth)
- `iamUserAccessToBilling: DENY` — only the root user sees billing
- Billing and security alternate contacts on file from day one
- The default delete contract (removal, not closure) — destroying the
  resource never destroys a production account's data

## Customize

- Use a plus-addressed root email (`aws+workloads-prod@example.com`) —
  emails are unique across ALL of AWS, forever-ish
- Add the `operations` contact and `primaryContact` for complete
  records
- `roleName` names the bootstrap role AWS pre-creates (default
  `OrganizationAccountAccessRole`) — write-once, decide before the
  first deploy
