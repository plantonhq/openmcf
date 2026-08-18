# AwsAppSyncApi — Pulumi module

Manages one AppSync API — `appsync.GraphQLApi` XOR `appsync.Api` (Events) — with its satellites: `appsync.DataSource`, `appsync.Type`, `appsync.Function`, `appsync.Resolver`, `appsync.ApiCache`, `appsync.ApiKey`, `appsync.ChannelNamespace`, `appsync.DomainName` + `appsync.DomainNameApiAssociation`, `appsync.SourceApiAssociation` (MERGED), and `wafv2.WebAclAssociation`.

Module facts worth knowing before editing:

- **The arm decides the pivot**: `graphql.go` builds the GraphQL API (plus the WAF association and cache singleton), `events.go` the Events API; both return the shared `createdApi` handle every satellite hangs off (`datasources.go`, `satellites.go`).
- **api_type is derived** from the merged block's presence; resolver `kind` from pipeline_functions. Neither is spec surface.
- **Name joins carry the edges**: in-spec data source names resolve through the created `DataSource` (its Name output), and pipeline entries resolve to created functions' `FunctionId` outputs — externally created objects pass through as literals.
- **One-value vocabularies are pinned**: APPSYNC_JS, AWS_IAM (HTTP signing), RDS_HTTP_ENDPOINT.
- **The EventBridge data source is replace-to-change** at the pin (upstream update defect) — rename the entry.
- **`type_formats` is an import-derivation echo map** (the type's composite import ID carries its format, which no resource output echoes).

Outputs mirror the Terraform module key-for-key: `api_id` (import ID), `api_arn`, `graphql_url`, `realtime_url`, `events_http_endpoint`, `events_realtime_endpoint`, `appsync_domain_name`, `domain_hosted_zone_id`, `datasource_arns`, `function_ids`, `api_key_ids`, `channel_namespace_arns`, `source_api_association_ids`, `type_formats`.
