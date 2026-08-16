# AwsCloudwatchSynthetics — Pulumi module (Go)

Manages one optional canary, owned groups, and the canary's group joins (`synthetics.Canary`, `synthetics.Group`, `synthetics.GroupAssociation`).

Module facts worth knowing before editing:

- **The canary's name is `metadata.name`** (the canary charset is lowercase letters, digits, hyphens, underscores); renames replace the canary.
- **The artifact location is assembled here** — the spec's bucket + prefix become the provider's one `s3://bucket/prefix` string.
- **`StartCanary: false` creates the canary READY but never running** — no run costs until started; the provider stops/starts around updates.
- **CREATE_FAILED canaries are deleted and recreated by the provider** — AWS offers no other repair.
- **`RunConfig.EnvironmentVariables` are write-only at AWS** — reads never return them; never put secrets there.
- **Joins reference the group by NAME** — owned groups get a DependsOn; external names resolve as-is. The association is create/delete only and untaggable.

Outputs mirror the Terraform module key-for-key: `canary_name`, `canary_arn`, `engine_arn`, `source_location_arn`, `canary_status`, `group_arns`/`group_ids` (keyed by group name).
