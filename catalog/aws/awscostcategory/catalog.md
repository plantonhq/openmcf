# AWS Cost Category

Your own dimension over AWS spend: rules that group every cost line
into named values ("Cost Center" = platform/data/mobile), visible
everywhere Cost Explorer reports — with optional split-charge rules
for shared-cost chargeback.

## What Gets Managed

- The category and its ordered rules — expression rules (services,
  accounts, tags, other categories) and inherited-value rules (the
  value IS the account name or a tag's value).
- The default value unmatched spend lands in.
- Split-charge rules: re-allocate a shared value's costs across
  targets proportionally, evenly, or by fixed percentages.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Cost Explorer permissions.

### AWS Account

- Cost categories are free; tag-based rules need cost-allocation tags
  activated in Billing.

## Deploy

### Console

Create the resource from the AWS catalog, define the rules in
priority order, and deploy.

### CLI

```bash
planton apply -f cost-category.yaml
```

## After Deploy

- AWS categorizes the current and previous month retroactively;
  values appear as a Cost Explorer grouping dimension.
- Outputs publish the category's ARN, name, and effective window.
