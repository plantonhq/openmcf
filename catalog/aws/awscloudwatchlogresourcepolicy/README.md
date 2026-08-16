# AwsCloudwatchLogResourcePolicy

One CloudWatch Logs resource policy: the IAM document that grants AWS services permission to write logs — Route53 query logging, EventBridge, OpenSearch slow logs, and every other service whose setup asks for a "CloudWatch Logs resource policy".

## Highlights

- **Exactly one scope per instance**: `policy_name` = a named account-wide policy (the common shape); `resource_arn` = a policy pinned to one log group. Both are identity — changing either replaces.
- **Revision-guarded**: AWS versions the policy; the modules pass the tracked revision on updates so concurrent out-of-band edits fail loudly instead of silently overwriting. Resource-scoped deletes REQUIRE the tracked revision.
- **The document is structured configuration** diffed semantically as an IAM document.
- **Untaggable at AWS** — the deliberate absence from the catalog's tag convention.

## Both Engines

Both modules render the single resource identically and export the same outputs: `policy_id` (name or ARN — the import ID; an ARN-shaped import selects the resource scope), `policy_scope`, `revision_id`.

## Chart Wiring

`resource_arn` → AwsCloudwatchLogGroup `log_group_arn`; the document's Resource entries typically reference the log groups the granted service writes to.
