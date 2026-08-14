# AwsRestApiUsagePlan — Terraform/OpenTofu module

Deploys an API Gateway usage plan (`aws_api_gateway_usage_plan`) plus
the API keys it admits (`aws_api_gateway_api_key`) and the
attachments (`aws_api_gateway_usage_plan_key`).

Module facts worth knowing before editing:

- **Keys are created and attached here.** A key with no plan is an
  orphan; the module always attaches every `api_keys` entry.
- **Key values are sensitive and not outputs.** IDs and ARNs are
  exported; the value is read from AWS when distributing to consumers.
- **Quota and throttle are independent optional blocks.** Omit either
  to skip that ceiling.
- **Method throttles live on the `api_stages` entry** they apply to
  (`path` is `{resource}/{METHOD}`).

Outputs mirror the Pulumi module key-for-key: `usage_plan_id`,
`usage_plan_arn`, `api_key_ids`, `api_key_arns`.
