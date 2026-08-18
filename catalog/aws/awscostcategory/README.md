<p align="center">
  <img src="logo.svg" alt="AWS Cost Category" width="80"/>
</p>

# AWS Cost Category

Manage a [Cost Explorer cost category](https://docs.aws.amazon.com/awsaccountbilling/latest/aboutv2/manage-cost-categories.html)
— a named dimension YOU define that groups every line of spend into
exactly one of its values, everywhere Cost Explorer reports.

## What Gets Managed

- **The category** (`spec.categoryName` is the name — an explicit
  field because category names legally carry spaces `metadata.name`
  cannot): its ordered **rules** (first match wins) — REGULAR rules
  assign a value you name to the spend a Cost Explorer expression
  selects (dimensions, tags, and other cost categories composed with
  and/or/not — the leveled shape is exactly the nesting AWS accepts);
  INHERITED_VALUE rules take the value from a dimension itself (the
  account name, a tag's value — one rule fanning out to N values).
- **defaultValue** for unmatched spend (otherwise "Uncategorized"),
  the **effectiveStart** month, and **splitChargeRules** re-allocating
  one value's costs across targets (proportionally, evenly, or by
  fixed percentages) for shared-cost chargeback.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
