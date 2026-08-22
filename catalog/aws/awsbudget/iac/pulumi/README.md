# AwsBudget — Pulumi module (Go)

Manages one Budgets budget (`budgets.Budget`) and its folded actions
(`budgets.BudgetAction`, one per `spec.actions` entry).

Module facts worth knowing before editing:

- **The name is `spec.budget_name`, never `metadata.name`** — budget
  names legally carry spaces and mixed case metadata.name cannot.
- **Exactly one funding shape renders** (spec-enforced): the fixed
  limit maps to `LimitAmount`/`LimitUnit`, planned limits to
  `PlannedLimits`, auto-adjustment to `AutoAdjustData`.
- **The filter generations are mutually exclusive** (spec-enforced):
  legacy `cost_filters`/`cost_types` vs the modern `metric` +
  `filter_expression` pair. The provider's `metrics` is a
  single-element list the SDK flattens to one string — the spec's
  singular `metric` maps 1:1.
- **The filter tree is LEVELED (root → node → leaf)** — exactly the
  nesting AWS accepts, so `filter_expression.go` is a 1:1 typed walk
  with no depth checks (the SDK generates a distinct type per tree
  path; neither side can express an illegal tree).
- **Actions key by `spec.actions[].name`** — the logical name and the
  `action_ids` output-map key. Definition arms match `action_type`
  (spec-enforced); references (execution role, policy ARNs,
  principals, instance IDs) arrive resolved.
- **`AccountId` renders on budget AND actions when set** — a
  member-account budget's actions must live in the same account.

Outputs mirror the Terraform module key-for-key: `budget_name`,
`budget_arn`, `account_id`, `action_ids` (map keyed by action name).
