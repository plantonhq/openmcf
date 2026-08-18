# AwsRoute53ResolverQueryLog

A Resolver query logging configuration — the record of every DNS query VPCs make through Route 53 Resolver — with its VPC associations managed in-line. The destination is a CloudWatch log group, an S3 bucket, or a Kinesis Data Firehose stream.

## Highlights

- **This is RESOLVER logging, not zone logging**: it records everything associated VPCs ask (cached answers, forwards, firewall verdicts included) — a different surface from the hosted-zone query logging on AwsRoute53Zone, which sees only what Route 53 answers for one public zone.
- **Immutable by contract, taught up front**: name and destination are both ForceNew — the configuration replaces on change and existing log data stays in the destination.
- **The async-failure class is the verifier's job**: an association against an unwritable destination FAILS after a clean apply — the E2E verifier asserts ACTIVE, never mere existence.

## Both Engines

Both modules render the configuration and its associations identically and export the same outputs: `query_log_config_id` (import ID), `query_log_config_arn`, `share_status`, plus the `association_ids` map keyed by the resolved VPC id.

## Chart Wiring

`destination_arn` → AwsCloudwatchLogGroup `log_group_arn` (or any S3 bucket / Firehose stream ARN); `vpc_ids` → AwsVpc `vpc_id`. Pair with AwsRoute53ResolverFirewall to see which queries its rules fired on.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
