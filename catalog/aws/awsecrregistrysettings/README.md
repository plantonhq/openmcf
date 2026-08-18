# AwsEcrRegistrySettings

The REGISTRY-level ECR configuration for one region — the registry policy, scanning, replication, pull-through cache rules, repository creation templates, account settings, and pull-time update exclusions. Individual repositories are the AwsEcrRepo kind; this kind governs what all of them share.

## Highlights

- **A settings singleton, honestly framed**: AWS keeps exactly one private registry per account+region — deploy at most one instance per region, and know that destroy semantics differ per arm (the policy and keyed collections delete; scanning and replication RESET to empty defaults; account settings PERSIST at last-applied values — each taught on its arm).
- **Typed account settings, not magic strings**: the provider's name/value pairs became three typed fields (`basic_scan_type_version`, `blob_mounting`, `registry_policy_scope`) — the console can offer real choices and validation catches invalid pairings at manifest time.
- **The template guard the provider forgot**: mutability exclusion filters require an *_WITH_EXCLUSION mode — enforced here at validation (the provider's repository resource has this guard; its template resource does not, and the API failure is late and cryptic).

## Both Engines

Both modules render every arm identically and export the same outputs: `registry_id` (the account id — the singletons' import ID), `registry_url` (the pull URL base), plus keyed maps for cache rules, creation templates, and pull-time exclusions.

## Chart Wiring

`pull_through_cache_rules[].credential_arn` → AwsSecretsManagerSecret `secret_arn`; `custom_role_arn` fields → AwsIamRole `role_arn`; template `encryption.kms_key` → AwsKmsKey `key_arn`; `pull_time_update_exclusions` → AwsIamRole `role_arn`. Pair with AwsEcrRepo kinds — this kind sets the registry posture their images live under.
