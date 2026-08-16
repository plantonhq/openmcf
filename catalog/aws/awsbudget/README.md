<p align="center">
  <img src="logo.svg" alt="AWS Budget" width="80"/>
</p>

# AWS Budget

Manage an [AWS Budgets budget](https://docs.aws.amazon.com/cost-management/latest/userguide/budgets-managing-costs.html)
— a spend or usage threshold AWS evaluates continuously — with its
alert notifications and folded budget ACTIONS (automatic or staged
responses when a threshold breaches).

## What Gets Managed

- **The budget** (`spec.budgetName` is the name — an explicit field
  because budget names legally carry spaces `metadata.name` cannot):
  its **type** (COST, USAGE, or the RI/Savings Plans coverage and
  utilization trackers), **time unit**, and exactly ONE funding shape —
  a fixed **limit**, per-period **plannedLimits**, or **autoAdjust**
  (AWS recomputes the limit from history or forecast).
- **Filtering**, one generation per budget: the modern **metric** +
  **filterExpression** tree (dimensions, tag keys, and cost categories
  composed with and/or/not — the leveled shape is exactly the nesting
  AWS accepts) or the legacy **costFilters**/**costTypes** pairs.
- **Notifications**: actual/forecasted threshold alerts to email
  addresses and SNS topics.
- **Actions** — the folded satellites: each entry applies a
  restrictive IAM policy, attaches an SCP, or stops EC2/RDS instances
  via SSM when its threshold breaches, automatically or staged for
  MANUAL approval, through a budgets-trusting execution role.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
