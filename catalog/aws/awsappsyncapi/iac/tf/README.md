# AwsAppSyncApi — OpenTofu module

Manages one AppSync API — `aws_appsync_graphql_api` XOR `aws_appsync_api` (Events) — with its satellites: `aws_appsync_datasource`, `aws_appsync_type`, `aws_appsync_function`, `aws_appsync_resolver`, `aws_appsync_api_cache`, `aws_appsync_api_key`, `aws_appsync_channel_namespace`, `aws_appsync_domain_name` + its API association, `aws_appsync_source_api_association` (MERGED), and `aws_wafv2_web_acl_association`.

Module facts worth knowing before editing:

- **variables.tf is generator-owned** — regenerate with `planton tofu generate-variables AwsAppSyncApi`, never hand-edit.
- **api_type is derived**: the merged block's presence renders `MERGED`; there is no spec field for it. Same idiom for resolver `kind` (pipeline_functions non-empty → `PIPELINE`).
- **Name joins carry the edges**: functions and channel handlers reference data sources by spec name through the created resource; resolvers' `pipeline_functions` join names to AWS function ids. Externally created objects pass through as literals.
- **Resolver instances key `"type//field"`** — the `//` separator is the import bridge's address-key segment convention.
- **One-value vocabularies are pinned**: APPSYNC_JS (the only runtime), AWS_IAM (the only HTTP auth type), RDS_HTTP_ENDPOINT (the only relational source type).
- **The EventBridge data source is replace-to-change** — the provider's update path drops its config upstream (recorded in the import catalog and _inbox); rename the entry instead of editing it.
- **The api key secret is never an output** — the provider degrades the `key` attribute to the key ID on read; the `api_key_ids` map is honest about carrying IDs.

Outputs mirror the Pulumi module key-for-key: `api_id` (import ID), `api_arn`, `graphql_url`, `realtime_url`, `events_http_endpoint`, `events_realtime_endpoint`, `appsync_domain_name`, `domain_hosted_zone_id`, `datasource_arns`, `function_ids`, `api_key_ids`, `channel_namespace_arns`, `source_api_association_ids`, `type_formats`.
