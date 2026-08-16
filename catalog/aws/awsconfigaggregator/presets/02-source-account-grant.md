# Source Account Grant

This preset is the OTHER side of a cross-account rollup: deployed in a
source account, it authorizes the named aggregator account+region to
collect this account's Config data. No aggregator is created here.

## When to Use

- Each member/source account of an account-list rollup (organization
  sources need no grants)
- Standing security-account topologies where one account aggregates
  many

## What You Get

- The reciprocal permission record AWS requires before cross-account
  aggregation flows
- The aggregator's pending-authorization source flips to live the
  moment this applies

## Customize

- `accountId` / `authorizedAwsRegion` name the AGGREGATOR (the
  collector), not this account
- Add more entries to authorize several aggregators from one
  deployment
- Add an `aggregation` block if this account also runs its own rollup
  — both arms are legal together
