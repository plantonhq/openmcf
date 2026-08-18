# AwsEcrRegistrySettings — Pulumi module

Manages one region's registry-level ECR configuration: `ecr.RegistryPolicy`, `ecr.RegistryScanningConfiguration`, `ecr.ReplicationConfiguration` (conditional singletons), `ecr.PullThroughCacheRule` and `ecr.RepositoryCreationTemplate` (one per spec entry, named by prefix), `ecr.AccountSetting` (one per configured toggle), and `ecr.PullTimeUpdateExclusion` (one per principal).

Module facts worth knowing before editing:

- **Reset-not-delete arms**: scanning and replication destroy by putting the empty default back; account settings' delete is a no-op (values persist).
- **The one-value vocabularies are pinned**: WILDCARD (scanning/template filters) and PREFIX_MATCH (replication filters) are the only values AWS supports — never spec surface.
- **The account id comes from `aws.GetCallerIdentity`** — the registry's identity, exported as `registry_id` and composed into `registry_url`.
- **Typed account settings render name/value pairs**: each configured field becomes one PutAccountSetting upsert with the matching setting name.
- **Nothing here is taggable at AWS** — a template's `ResourceTags` are the STAMPED repositories' tags (user surface), not this module's identity tags.

Outputs mirror the Terraform module key-for-key: `registry_id`, `registry_url`, and the three keyed maps.
