# AwsCostCategory — Terraform/OpenTofu module

Manages one Cost Explorer cost category (`aws_ce_cost_category`).

Module facts worth knowing before editing:

- **The name is `spec.category_name`, never `metadata.name`** —
  category names legally carry spaces metadata.name cannot; changing
  it replaces the category.
- **`rule_version` is module-pinned** to `CostCategoryExpression.v1`,
  the only value the AWS API accepts — never a spec knob (recorded as
  a manifest exclusion).
- **Rules are ORDERED** (first match wins) and render exactly in spec
  order; each is either a REGULAR expression rule or an
  INHERITED_VALUE dimension rule (spec-enforced shapes).
- **The rule expression is LEVELED (root → node → leaf)** — exactly
  the nesting AWS accepts, so the dynamic blocks unroll it 1:1 with no
  depth checks (the inner `rule` dynamic uses an explicit `expr`
  iterator to avoid shadowing the outer rule; the Pulumi module walks
  the same levels).
- **FIXED split-charge rules carry allocation percentages**
  (spec-enforced).

Outputs mirror the Pulumi module key-for-key: `category_arn`,
`category_name`, `effective_start`, `effective_end`.
