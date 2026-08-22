# Developers Group

This preset creates the everyday engineering group: PowerUserAccess
for building, with an inline deny guard keeping billing and cost
configuration out of reach.

## When to Use

- Engineering teams that provision infrastructure but must not touch
  billing, budgets, or cost tooling
- Any group where one broad managed policy needs a narrow local
  carve-out

## What You Get

- PowerUserAccess (everything except IAM/organization management)
  via AWS's maintained managed policy
- An inline `deny-billing` policy that lives and dies with the group
  — explicit deny beats any allow

## Customize

- Deny wins: extend the inline document's Action list for other
  carve-outs (regions, services) rather than forking the managed
  policy
- Managed-policy attachments reconcile individually — add your own
  AwsIamPolicy references (`status.outputs.policy_arn`) alongside the
  AWS-managed one
- Membership is authoritative: the declared users ARE the group;
  additions made in the console disappear on the next apply
