# AwsCostCategory — Pulumi module (Go)

Manages one Cost Explorer cost category
(`costexplorer.CostCategory`).

Module facts worth knowing before editing:

- **The name is `spec.category_name`, never `metadata.name`** —
  category names legally carry spaces metadata.name cannot; changing
  it replaces the category.
- **`RuleVersion` is module-pinned** to `CostCategoryExpression.v1`,
  the only value the AWS API accepts — never a spec knob (recorded as
  a manifest exclusion).
- **Rules are ORDERED** (first match wins) and render exactly in spec
  order; each is either a REGULAR expression rule or an
  INHERITED_VALUE dimension rule (spec-enforced shapes).
- **The rule expression is LEVELED (root → node → leaf)** — exactly
  the nesting AWS accepts, so `expression.go` is a 1:1 typed walk
  with no depth checks (the SDK generates a distinct type per tree
  path; neither side can express an illegal tree).
- **FIXED split-charge rules carry allocation percentages**
  (spec-enforced).

Outputs mirror the Terraform module key-for-key: `category_arn`,
`category_name`, `effective_start`, `effective_end`.
