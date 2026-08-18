# AwsCloudwatchSynthetics

CloudWatch Synthetics: a CANARY — a scheduled scripted probe (a Synthetics-managed Lambda running your Node.js/Python script) that checks endpoints, flows, and APIs from the outside — plus the grouping surface (owned groups and the canary's group joins).

## Highlights

- **Two independently deployable arms**: a canary instance monitors an endpoint (optionally joining groups); a groups-only instance manages shared groups many canaries join. Joins reference the group NAME, so shared groups are referenced, never fought over.
- **Code stages through S3** (the AwsLambda idiom): bucket + key + optional version; small bundles travel inline as base64 through an AwsS3ObjectSet. The runtime dictates the zip layout (`nodejs/node_modules/<file>.js` for Node.js).
- **`start_canary: false` creates the canary READY but never running** — runs, not canaries, are what Synthetics bills.
- **Contracts taught in place**: CREATE_FAILED canaries are deleted and recreated by the provider; environment variables are write-only at AWS (never put secrets there); `delete_lambda` tears the managed Lambda down with the canary.

## Both Engines

Both modules render the canary, the owned groups (name-keyed), and the joins identically, and export the same outputs: `canary_name`, `canary_arn`, `engine_arn`, `source_location_arn`, `canary_status`, and `group_arns`/`group_ids` keyed by group name.

## Chart Wiring

`execution_role_arn` → AwsIamRole, `artifact_bucket` + `code.s3_bucket` → AwsS3Bucket, `vpc_config` subnets/security groups → AwsSubnet/AwsSecurityGroup, artifact encryption key → AwsKmsKey.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
