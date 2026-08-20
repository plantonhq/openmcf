# Enforced Sandbox Budget

This preset creates a tag-scoped budget that does not just alert — at
100% of the limit it AUTOMATICALLY applies AWS's deny-all policy to
the sandbox users group, freezing further provisioning until spend
resets or a human reverses it.

## When to Use

- Sandbox and experiment accounts where overspend must stop itself
- Any environment where a hard spend ceiling beats a paged human

## What You Get

- A monthly budget scoped by the `environment: sandbox` tag through
  the modern filter expression
- An 80% early-warning alert, then an automatic IAM freeze at 100% —
  AWS reverses the policy when spend drops back under the threshold

## Customize

- Change `approvalModel: MANUAL` to stage the freeze for human
  approval instead of executing it
- Swap the arm: `scpActionDefinition` attaches an SCP to organization
  targets; `ssmActionDefinition` stops tagged EC2/RDS instances
- The execution role must trust budgets.amazonaws.com and carry
  AWS's `AWSBudgetsActionsWithAWSResourceControlAccess` policy
