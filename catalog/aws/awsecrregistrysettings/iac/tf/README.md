# AwsEcrRegistrySettings — Terraform/OpenTofu module

Manages one region's registry-level ECR configuration: `aws_ecr_registry_policy`, `aws_ecr_registry_scanning_configuration`, `aws_ecr_replication_configuration` (all count-gated singletons), `aws_ecr_pull_through_cache_rule` and `aws_ecr_repository_creation_template` (keyed by prefix), `aws_ecr_account_setting` (keyed by setting name), and `aws_ecr_pull_time_update_exclusion` (keyed by principal ARN).

Module facts worth knowing before editing:

- **Reset-not-delete arms**: scanning and replication destroy by putting the empty default back; account settings' delete is a no-op (values persist).
- **The one-value vocabularies are pinned**: WILDCARD (scanning/template filters) and PREFIX_MATCH (replication filters) are the only values AWS supports — never spec surface.
- **The account id comes from `data.aws_caller_identity`** — the registry's identity, exported as `registry_id` and composed into `registry_url`.
- **Credential/custom-role clearing is not propagated upstream** — replacing the rule is the only genuine drop (taught on the spec).
- **Nothing here is taggable at AWS** — a template's `resource_tags` are the STAMPED repositories' tags (user surface), not this module's identity tags.

Outputs mirror the Pulumi module key-for-key: `registry_id`, `registry_url`, and the three keyed maps.
