# On-Demand Simple Table

This preset creates a DynamoDB table with on-demand (pay-per-request) billing and a simple partition key. On-demand pricing automatically scales to handle any traffic level without capacity planning and costs nothing while idle. Point-in-time recovery and deletion protection are enabled for production safety. This is the 30-second default for most DynamoDB use cases.

## When to Use

- Key-value stores, session stores, user profiles, or any workload with a simple primary key
- Applications with unpredictable or variable traffic where capacity planning is impractical
- New tables where traffic patterns are not yet established

## Key Configuration Choices

- **On-demand billing** (`billingMode: PAY_PER_REQUEST`) -- No capacity planning; pay only for reads/writes consumed; scales instantly. Add an `onDemandThroughput` ceiling later if you want a hard spend guardrail.
- **String partition key** (`keySchema: [{attributeName: id, keyType: HASH}]`) -- Simple HASH key on `id`; sufficient for most key-value access patterns
- **Point-in-time recovery** (`pointInTimeRecovery.enabled: true`) -- Continuous backups enabling restoration to any second in the recovery window (the AWS default keeps 35 days; set `recoveryPeriodInDays` to shorten it)
- **Deletion protection** (`deletionProtectionEnabled: true`) -- Prevents accidental table deletion

## Placeholders to Replace

- `<aws-region>` -- The AWS region for the table (e.g. `us-west-2`)

The table uses a generic `id` partition key. Rename the `id` attribute and adjust the type (`S` for string, `N` for number, `B` for binary) based on your data model.

## Related Presets

- **02-provisioned-production** -- Use instead for sustained, predictable workloads where provisioned capacity is more cost-effective
- **03-global-table** -- Use instead for multi-region active-active deployments
