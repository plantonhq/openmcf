<p align="center">
  <img src="logo.svg" alt="AWS Config Aggregator" width="80"/>
</p>

# AWS Config Aggregator

Manage [AWS Config aggregation](https://docs.aws.amazon.com/config/latest/developerguide/aggregate-data.html)
— the cross-account, cross-region rollup of Config data into one
queryable view, covering BOTH sides of the relationship.

## What Gets Managed

- **The aggregator** (`spec.aggregation`; `metadata.name` is the
  aggregator name): what it collects from — an explicit account list
  or the whole AWS Organization, exactly one — across listed regions
  or all of them. It references no Config recorder: aggregation works
  in an account with zero recorders.
- **The reciprocal grants** (`spec.authorizations`): deployed in each
  SOURCE account, each grant names the aggregator account+region
  allowed to collect from it. Organization-sourced aggregators need
  no grants.

A same-account rollup needs only the aggregation arm. A cross-account
topology deploys this component twice: once in the aggregator account
(aggregation arm) and once per source account (authorizations arm).
An account playing both roles declares both.

Destroying this component **deletes whichever arms it manages** —
aggregated data disappears from the view; source accounts' Config
data is untouched. Aggregators and grants are free.

See [v1alpha1/reference.md](v1alpha1/reference.md) for the full field
reference generated from the spec proto.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
