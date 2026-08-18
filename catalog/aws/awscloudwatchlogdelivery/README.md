# AwsCloudwatchLogDelivery

The two ways logs leave CloudWatch, as two independently deployable arms: the VENDED delivery pipeline (a delivery source wrapping one AWS resource whose service vends logs, fanning out through deliveries to S3/CloudWatch Logs/Firehose/X-Ray destinations) and the legacy CROSS-ACCOUNT Kinesis destination other accounts' subscription filters target.

## Highlights

- **The vended arm pivots on the source** (one instance = one logged resource + log type) with name-keyed deliveries; each delivery's destination is created inline or referenced by ARN (the own-XOR-existing idiom), so one-source-to-many-destinations and many-sources-to-one-destination both compose without instance collisions.
- **Destination-only instances work too** — a central logging team owns shared destinations (with their cross-account policies); producer instances reference them by ARN.
- **AWS contracts taught in place**: one delivery per (source, destination-type) pair; the CloudFront S3 prefix AWS silently prepends; the cross-account destination's access policy PERSISTS when only the policy is destroyed (the no-op-delete class) and the first create retries through IAM trust propagation.

## Both Engines

Both modules render the same six-resource surface and export the same outputs: `source_name`/`source_arn`/`source_service`, `destination_arns` keyed by name, `delivery_ids`/`delivery_arns` keyed by name, and `cross_account_destination_name`/`_arn`.

## Chart Wiring

The source's `resource_arn` takes any vended-logs producer's ARN output (a Bedrock knowledge base, a CloudFront distribution, ...); destination targets reference AwsS3Bucket/AwsCloudwatchLogGroup/AwsKinesisFirehose ARNs; the cross-account arm wires AwsIamRole + AwsKinesisStream.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
