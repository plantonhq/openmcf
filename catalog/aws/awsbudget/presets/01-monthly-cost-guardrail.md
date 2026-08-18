# Monthly Cost Guardrail

This preset creates the canonical monthly cost budget: a fixed dollar
ceiling with an early-warning alert at 80% of actual spend and a
forecast alert when AWS projects the month will breach.

## When to Use

- The first budget every account should carry — a monthly ceiling
  with humans in the loop
- Team or environment budgets when paired with a filter expression

## What You Get

- A monthly COST budget with a fixed limit (free — the alert-only
  budget class)
- Two alerts: actual spend at 80%, forecasted spend at 100%

## Customize

- Scope the budget with `metric: UnblendedCost` plus a
  `filterExpression` (services, accounts, activated tag keys, cost
  categories)
- Add SNS topics to the alerts (`subscriberSnsTopicArns`) for
  programmatic fan-out
- Escalate from alerts to enforcement with the automatic-guardrail
  preset's `actions`
