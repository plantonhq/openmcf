locals {
  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsAppSyncApi"
    "planton.ai/resource-id"   = var.metadata.id
  }

  is_graphql = var.spec.graphql != null
  is_events  = var.spec.events != null

  # A MERGED API is declared by the merged block's presence; the
  # provider's api_type argument is derived from it.
  is_merged = local.is_graphql && var.spec.graphql.merged != null

  # The one API id every satellite hangs off, whichever arm created it.
  api_id = local.is_graphql ? aws_appsync_graphql_api.this[0].id : aws_appsync_api.this[0].api_id

  # Satellites keyed by their spec names (the for_each keys both
  # engines and the output maps share).
  datasources = { for d in var.spec.datasources : d.name => d }
  api_keys    = { for k in var.spec.api_keys : k.name => k }
  types       = local.is_graphql ? { for t in var.spec.graphql.types : t.name => t } : {}
  functions   = local.is_graphql ? { for f in var.spec.graphql.functions : f.name => f } : {}
  # "type//field" - the "//" separator is the import bridge's
  # address-key segment convention (segment 0 = type, 1 = field).
  resolvers          = local.is_graphql ? { for r in var.spec.graphql.resolvers : "${r.type}//${r.field}" => r } : {}
  channel_namespaces = local.is_events ? { for n in var.spec.events.channel_namespaces : n.name => n } : {}
  source_apis        = local.is_merged ? { for s in var.spec.graphql.merged.source_apis : s.name => s } : {}
}
