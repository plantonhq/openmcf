# AwsCloudwatchLogDelivery — Terraform/OpenTofu module

Manages the vended-log pipeline (`aws_cloudwatch_log_delivery_source`, `aws_cloudwatch_log_delivery_destination` + policy, `aws_cloudwatch_log_delivery`) and/or the legacy cross-account destination (`aws_cloudwatch_log_destination` + policy).

Module facts worth knowing before editing:

- **The vended source is per (resource, log_type)** — name, log type, and resource ARN all replace on change; one source fans out to name-keyed deliveries.
- **Each delivery's destination is own-XOR-existing** (spec-guaranteed): owned destinations resolve to their resource's ARN with an implicit dependency; external destinations pass the literal ARN through.
- **AWS allows ONE delivery per (source, destination-type) pair** — one to S3 plus one to Firehose, never two to S3.
- **CloudFront suffix paths**: AWS prepends `AWSLogs/{account-id}/CloudFront/` server-side; the provider strips it on reads — configure only your own path segment.
- **The `destinations` variable arrives untyped (`any`)** — its entries carry a free-form policy Struct, so locals normalize each entry (the inline-policies class).
- **The cross-account destination's policy delete is a NO-OP at AWS** — the policy persists; only destroying the destination is real. Its first create retries up to two minutes for IAM trust propagation.
- **Tags land on the taggable four** (source, destinations, deliveries, cross-account destination); AWS does not return them on Get for the vended family.

Outputs mirror the Pulumi module key-for-key: `source_name`/`source_arn`/`source_service`, `destination_arns` (keyed by name), `delivery_ids`/`delivery_arns` (keyed by name), `cross_account_destination_name`/`_arn`.
