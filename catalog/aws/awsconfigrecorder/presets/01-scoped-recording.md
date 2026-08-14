# Scoped Recording

This preset records exactly the types your compliance rules evaluate —
the cost-deliberate posture, since AWS Config bills per recorded
configuration item.

## When to Use

- Starting AWS Config in a region for specific compliance rules
- Accounts where `all_supported` recording would bill on every noisy
  resource change

## What You Get

- Continuous recording of the listed types only, delivered to the
  history bucket with daily snapshots
- One year of queryable configuration history

## Customize

- Extend `resourceTypes` with the types your AwsConfigRule instances
  evaluate — the inclusion list is the bill
- The role needs the managed `AWS_ConfigRole` policy plus write access
  to the bucket; the bucket needs the `config.amazonaws.com` policy
- Raise `retentionPeriodInDays` (up to 2557) for longer audit windows

## Composing

Pair with AwsConfigRule instances scoped to the same resource types —
rules only see what the recorder captures.
