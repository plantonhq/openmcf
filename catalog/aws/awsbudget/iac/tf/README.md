# AwsBudget — Terraform/OpenTofu module

Manages one Budgets budget (`aws_budgets_budget`) and its folded
actions (`aws_budgets_budget_action`, one per `spec.actions` entry).

Module facts worth knowing before editing:

- **The name is `spec.budget_name`, never `metadata.name`** — budget
  names legally carry spaces and mixed case metadata.name cannot.
- **Exactly one funding shape renders** (spec-enforced): the fixed
  limit maps to `limit_amount`/`limit_unit`, planned limits to
  `planned_limit` blocks, auto-adjustment to `auto_adjust_data`.
- **The filter generations are mutually exclusive** (spec-enforced):
  legacy `cost_filter`/`cost_types` vs the modern `metrics` +
  `filter_expression` pair. The spec's singular `metric` renders as
  the provider's single-element `metrics` list.
- **The filter tree is LEVELED (root → node → leaf)** — exactly the
  nesting AWS accepts, so the dynamic blocks unroll it 1:1 with no
  depth checks (neither side can express an illegal tree; the Pulumi
  module walks the same levels).
- **Actions render with `for_each` keyed by `spec.actions[].name`** —
  the state key and the `action_ids` output-map key. Definition arms
  match `action_type` (spec-enforced).
- **`account_id` renders on budget AND actions when set** — a
  member-account budget's actions must live in the same account.

Outputs mirror the Pulumi module key-for-key: `budget_name`,
`budget_arn`, `account_id`, `action_ids` (map keyed by action name).
