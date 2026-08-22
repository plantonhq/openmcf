# AWS Budget

A spend or usage threshold AWS evaluates continuously — with alerts
when actual or forecasted spend crosses it, and optional automatic
guardrails (restrict IAM, attach an SCP, stop instances) when it
breaches.

## What Gets Managed

- The budget: cost or usage, its reset period, and exactly one funding
  shape — a fixed limit, per-period planned limits, or an
  auto-adjusting limit AWS recomputes from history or forecast.
- What the budget watches: a filter tree over services, regions,
  accounts, tags, and cost categories (or the legacy name/values
  filters).
- Alerts: actual/forecasted thresholds notifying email addresses and
  SNS topics.
- Actions: staged or automatic responses when a threshold breaches —
  apply a deny policy to IAM principals, attach an SCP, or stop
  EC2/RDS instances.

## Before You Deploy

### Planton Setup

- An organization and environment.
- An AWS provider connection with Budgets permissions.

### AWS Account

- Budgets is account-global and free (the first two action-enabled
  budgets are free; more bill pennies).
- Actions need an execution role trusting `budgets.amazonaws.com`
  ([AWS IAM Role](/cloud-catalog/aws-iam-role)); tag-based filters
  need cost-allocation tags activated in Billing.

## Deploy

### Console

Create the resource from the AWS catalog, pick the budget type and
limit, add alerts, and deploy.

### CLI

```bash
planton apply -f budget.yaml
```

## After Deploy

- The budget evaluates continuously; alerts fire as thresholds cross.
- Outputs publish the budget name/ARN, the owning account, and each
  action's AWS-generated ID.
