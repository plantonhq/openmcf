# AwsAppSyncApi

An AWS AppSync API — a GraphQL API (SDL schema, resolvers over data sources, caching, the MERGED federation variant) XOR an Events API (real-time pub/sub over channel namespaces) — with everything AWS attaches to the API managed in-line: data sources, schema types, functions, resolvers, API keys, channel namespaces, and a custom domain.

## Highlights

- **One kind, both AppSync models**: the spec pivots on exactly one of `graphql` / `events` — the same choice AWS's own Create API flow offers. Data sources, API keys, and the custom domain sit at the top level because AWS attaches them to either model.
- **The full GraphQL surface**: multi-auth (API key, IAM, Cognito, OIDC, Lambda authorizers — one primary plus additionals), APPSYNC_JS or VTL resolvers (UNIT and PIPELINE, with name-based function joins the modules resolve to AWS function ids), individually managed schema types, per-resolver caching over the one-per-API cache singleton, query-depth/resolver-count limits, X-Ray, enhanced metrics, logging, and a WAF web ACL attached from this side (the AwsAlb pattern).
- **The full Events surface**: per-phase authorization (connect/publish/subscribe against a shared provider list), channel namespaces with inline APPSYNC_JS handlers or DIRECT data-source integrations, and namespace-level auth overrides.
- **Federation as a block**: a MERGED API declares its source APIs as references to other AwsAppSyncApi resources — schema and resolvers live on the sources, never the merged surface (AWS's contract, enforced by validation).
- **Sharp edges taught where you meet them**: GraphQL API names forbid hyphens (the explicit `api_name` field), the schema has no provider drift detection, the custom domain's ACM certificate must live in us-east-1, EventBridge data sources are replace-to-change at this provider pin, and an API key's secret is only ever shown at creation.

## Both Engines

Both modules render the pivot, data sources, and satellites identically and export the same outputs: `api_id` (import ID and the MERGED-source join key), `api_arn`, the GraphQL `graphql_url`/`realtime_url`, the Events `events_http_endpoint`/`events_realtime_endpoint`, the custom domain's `appsync_domain_name`/`domain_hosted_zone_id`, and the keyed maps (`datasource_arns`, `function_ids`, `api_key_ids`, `channel_namespace_arns`, `source_api_association_ids`, `type_formats`).

## Chart Wiring

Auth: `user_pool_id` → AwsCognitoUserPool; Lambda authorizers → AwsLambda `function_arn`. Data sources: DynamoDB `table_name` → AwsDynamodb; Lambda `function_arn` → AwsLambda; `event_bus_arn` → AwsEventBridgeBus; OpenSearch `endpoint` → AwsOpenSearchDomain; Aurora `db_cluster_identifier` → AwsRdsCluster with `aws_secret_store_arn` → AwsSecretsManagerSecret. Roles (`service_role_arn`, logging, merged execution) → AwsIamRole `role_arn`. `custom_domain.certificate_arn` → AwsCertManagerCert `cert_arn` (us-east-1). `graphql.web_acl_arn` → AwsWafWebAcl `web_acl_arn`. MERGED `source_apis[].source_api_id` → another AwsAppSyncApi's `api_id`. Applications take `graphql_url` (GraphQL) or the Events endpoints as their API base.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
