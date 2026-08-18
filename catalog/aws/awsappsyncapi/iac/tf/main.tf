# One AppSync API - the GraphQL arm XOR the Events arm - with every
# satellite AWS attaches to it managed in-line: data sources, schema
# types, functions, resolvers, the cache singleton, API keys, channel
# namespaces, the custom domain, and (for MERGED APIs) source-API
# associations.
#
# Lifecycle facts the renders below depend on:
#   - api_type (GRAPHQL/MERGED, derived from the merged block) and
#     visibility are ForceNew on the GraphQL API; the Events API's
#     name replaces on change;
#   - the schema applies through AppSync's async StartSchemaCreation
#     and the provider performs NO drift detection on it - out-of-band
#     schema edits are invisible (recorded in the import catalog);
#   - the cache is a one-per-API singleton (its AWS id IS the API id);
#     both encryption flags replace the cache on change (a cold cache,
#     not an outage); cache operations wait up to 60 minutes upstream;
#   - resolver mutations on one API serialize behind the provider's
#     per-API mutex with a 2-minute conflict retry - bulk resolver
#     changes apply one at a time;
#   - the EventBridge data source's update path silently DROPS its
#     config (an upstream defect at the pin, recorded in _inbox):
#     treat EventBridge data sources as replace-to-change by renaming
#     the entry;
#   - aws_appsync_type ignores format changes on update (a perpetual
#     diff) - replace the type entry to change its format;
#   - the domain association's create/delete wait up to 60 minutes
#     upstream; the domain itself lives in the API's region but its
#     ACM certificate must be in us-east-1 (the CloudFront class);
#   - the api key's SECRET is returned only at creation; the
#     provider's key attribute holds the key ID after any refresh, so
#     the secret is deliberately not an output - fetch it from the
#     console/CLI.

# ---------------------------------------------------------------------------
# The GraphQL arm.
# ---------------------------------------------------------------------------

resource "aws_appsync_graphql_api" "this" {
  count = local.is_graphql ? 1 : 0

  # GraphQL API names forbid hyphens, so the spec carries an explicit
  # api_name instead of metadata.name (the explicit-name-field
  # convention).
  name = var.spec.graphql.api_name

  authentication_type = var.spec.graphql.auth.type

  # MERGED is derived from the merged block - the api_type argument is
  # never spec surface.
  api_type                      = local.is_merged ? "MERGED" : null
  merged_api_execution_role_arn = local.is_merged ? var.spec.graphql.merged.execution_role_arn : null

  dynamic "user_pool_config" {
    for_each = var.spec.graphql.auth.user_pool != null ? [var.spec.graphql.auth.user_pool] : []
    content {
      user_pool_id        = user_pool_config.value.user_pool_id
      default_action      = user_pool_config.value.default_action
      app_id_client_regex = user_pool_config.value.app_id_client_regex != "" ? user_pool_config.value.app_id_client_regex : null
      aws_region          = user_pool_config.value.aws_region != "" ? user_pool_config.value.aws_region : null
    }
  }

  dynamic "openid_connect_config" {
    for_each = var.spec.graphql.auth.openid_connect != null ? [var.spec.graphql.auth.openid_connect] : []
    content {
      issuer    = openid_connect_config.value.issuer
      client_id = openid_connect_config.value.client_id != "" ? openid_connect_config.value.client_id : null
      iat_ttl   = openid_connect_config.value.iat_ttl > 0 ? openid_connect_config.value.iat_ttl : null
      auth_ttl  = openid_connect_config.value.auth_ttl > 0 ? openid_connect_config.value.auth_ttl : null
    }
  }

  dynamic "lambda_authorizer_config" {
    for_each = var.spec.graphql.auth.lambda != null ? [var.spec.graphql.auth.lambda] : []
    content {
      authorizer_uri                   = lambda_authorizer_config.value.authorizer_uri
      authorizer_result_ttl_in_seconds = lambda_authorizer_config.value.authorizer_result_ttl_in_seconds > 0 ? lambda_authorizer_config.value.authorizer_result_ttl_in_seconds : null
      identity_validation_expression   = lambda_authorizer_config.value.identity_validation_expression != "" ? lambda_authorizer_config.value.identity_validation_expression : null
    }
  }

  dynamic "additional_authentication_provider" {
    for_each = var.spec.graphql.additional_auth_providers
    content {
      authentication_type = additional_authentication_provider.value.type

      # The additional provider's user pool carries NO default_action
      # (AWS's asymmetry; the spec's CEL walls enforce it).
      dynamic "user_pool_config" {
        for_each = additional_authentication_provider.value.user_pool != null ? [additional_authentication_provider.value.user_pool] : []
        content {
          user_pool_id        = user_pool_config.value.user_pool_id
          app_id_client_regex = user_pool_config.value.app_id_client_regex != "" ? user_pool_config.value.app_id_client_regex : null
          aws_region          = user_pool_config.value.aws_region != "" ? user_pool_config.value.aws_region : null
        }
      }

      dynamic "openid_connect_config" {
        for_each = additional_authentication_provider.value.openid_connect != null ? [additional_authentication_provider.value.openid_connect] : []
        content {
          issuer    = openid_connect_config.value.issuer
          client_id = openid_connect_config.value.client_id != "" ? openid_connect_config.value.client_id : null
          iat_ttl   = openid_connect_config.value.iat_ttl > 0 ? openid_connect_config.value.iat_ttl : null
          auth_ttl  = openid_connect_config.value.auth_ttl > 0 ? openid_connect_config.value.auth_ttl : null
        }
      }

      dynamic "lambda_authorizer_config" {
        for_each = additional_authentication_provider.value.lambda != null ? [additional_authentication_provider.value.lambda] : []
        content {
          authorizer_uri                   = lambda_authorizer_config.value.authorizer_uri
          authorizer_result_ttl_in_seconds = lambda_authorizer_config.value.authorizer_result_ttl_in_seconds > 0 ? lambda_authorizer_config.value.authorizer_result_ttl_in_seconds : null
          identity_validation_expression   = lambda_authorizer_config.value.identity_validation_expression != "" ? lambda_authorizer_config.value.identity_validation_expression : null
        }
      }
    }
  }

  # MERGED APIs own no schema - it merges in from the sources (the
  # spec's CEL wall keeps the arm honest).
  schema = var.spec.graphql.schema != "" ? var.spec.graphql.schema : null

  visibility           = var.spec.graphql.visibility != "" ? var.spec.graphql.visibility : null
  introspection_config = var.spec.graphql.disable_introspection ? "DISABLED" : null
  query_depth_limit    = var.spec.graphql.query_depth_limit > 0 ? var.spec.graphql.query_depth_limit : null
  resolver_count_limit = var.spec.graphql.resolver_count_limit > 0 ? var.spec.graphql.resolver_count_limit : null
  xray_enabled         = var.spec.graphql.xray_enabled

  dynamic "log_config" {
    for_each = var.spec.graphql.log_config != null ? [var.spec.graphql.log_config] : []
    content {
      cloudwatch_logs_role_arn = log_config.value.cloudwatch_logs_role_arn
      field_log_level          = log_config.value.field_log_level
      exclude_verbose_content  = log_config.value.exclude_verbose_content
    }
  }

  dynamic "enhanced_metrics_config" {
    for_each = var.spec.graphql.enhanced_metrics != null ? [var.spec.graphql.enhanced_metrics] : []
    content {
      data_source_level_metrics_behavior = enhanced_metrics_config.value.data_source_level_metrics_behavior
      operation_level_metrics_config     = enhanced_metrics_config.value.operation_level_metrics_config
      resolver_level_metrics_behavior    = enhanced_metrics_config.value.resolver_level_metrics_behavior
    }
  }

  tags = local.aws_tags
}

# The REGIONAL web ACL attaches from the protected resource's side
# (the AwsAlb pattern - the web ACL itself never knows its
# consumers). AWS's WAF association supports GraphQL APIs only.
resource "aws_wafv2_web_acl_association" "this" {
  count = local.is_graphql && var.spec.graphql.web_acl_arn != "" ? 1 : 0

  resource_arn = aws_appsync_graphql_api.this[0].arn
  web_acl_arn  = var.spec.graphql.web_acl_arn
}

# ---------------------------------------------------------------------------
# The Events arm.
# ---------------------------------------------------------------------------

resource "aws_appsync_api" "this" {
  count = local.is_events ? 1 : 0

  # Events API names allow hyphens - metadata.name is the naming basis.
  name          = var.metadata.name
  owner_contact = var.spec.events.owner_contact != "" ? var.spec.events.owner_contact : null

  event_config {
    dynamic "auth_provider" {
      for_each = var.spec.events.auth_providers
      content {
        auth_type = auth_provider.value.type

        dynamic "cognito_config" {
          for_each = auth_provider.value.cognito != null ? [auth_provider.value.cognito] : []
          content {
            user_pool_id        = cognito_config.value.user_pool_id
            aws_region          = cognito_config.value.aws_region
            app_id_client_regex = cognito_config.value.app_id_client_regex != "" ? cognito_config.value.app_id_client_regex : null
          }
        }

        dynamic "openid_connect_config" {
          for_each = auth_provider.value.openid_connect != null ? [auth_provider.value.openid_connect] : []
          content {
            issuer    = openid_connect_config.value.issuer
            client_id = openid_connect_config.value.client_id != "" ? openid_connect_config.value.client_id : null
            iat_ttl   = openid_connect_config.value.iat_ttl > 0 ? openid_connect_config.value.iat_ttl : null
            auth_ttl  = openid_connect_config.value.auth_ttl > 0 ? openid_connect_config.value.auth_ttl : null
          }
        }

        dynamic "lambda_authorizer_config" {
          for_each = auth_provider.value.lambda != null ? [auth_provider.value.lambda] : []
          content {
            authorizer_uri                   = lambda_authorizer_config.value.authorizer_uri
            authorizer_result_ttl_in_seconds = lambda_authorizer_config.value.authorizer_result_ttl_in_seconds > 0 ? lambda_authorizer_config.value.authorizer_result_ttl_in_seconds : null
            identity_validation_expression   = lambda_authorizer_config.value.identity_validation_expression != "" ? lambda_authorizer_config.value.identity_validation_expression : null
          }
        }
      }
    }

    dynamic "connection_auth_mode" {
      for_each = var.spec.events.connection_auth_modes
      content {
        auth_type = connection_auth_mode.value
      }
    }

    dynamic "default_publish_auth_mode" {
      for_each = var.spec.events.default_publish_auth_modes
      content {
        auth_type = default_publish_auth_mode.value
      }
    }

    dynamic "default_subscribe_auth_mode" {
      for_each = var.spec.events.default_subscribe_auth_modes
      content {
        auth_type = default_subscribe_auth_mode.value
      }
    }

    dynamic "log_config" {
      for_each = var.spec.events.log_config != null ? [var.spec.events.log_config] : []
      content {
        cloudwatch_logs_role_arn = log_config.value.cloudwatch_logs_role_arn
        log_level                = log_config.value.log_level
      }
    }
  }

  tags = local.aws_tags
}

# ---------------------------------------------------------------------------
# Data sources (either arm).
# ---------------------------------------------------------------------------

resource "aws_appsync_datasource" "this" {
  for_each = local.datasources

  api_id = local.api_id
  name   = each.value.name
  type   = each.value.type

  description      = each.value.description != "" ? each.value.description : null
  service_role_arn = each.value.service_role_arn != "" ? each.value.service_role_arn : null

  dynamic "dynamodb_config" {
    for_each = each.value.dynamodb != null ? [each.value.dynamodb] : []
    content {
      table_name             = dynamodb_config.value.table_name
      region                 = dynamodb_config.value.region != "" ? dynamodb_config.value.region : null
      use_caller_credentials = dynamodb_config.value.use_caller_credentials
      versioned              = dynamodb_config.value.versioned

      dynamic "delta_sync_config" {
        for_each = dynamodb_config.value.delta_sync != null ? [dynamodb_config.value.delta_sync] : []
        content {
          delta_sync_table_name = delta_sync_config.value.delta_sync_table_name
          base_table_ttl        = delta_sync_config.value.base_table_ttl > 0 ? delta_sync_config.value.base_table_ttl : null
          delta_sync_table_ttl  = delta_sync_config.value.delta_sync_table_ttl > 0 ? delta_sync_config.value.delta_sync_table_ttl : null
        }
      }
    }
  }

  dynamic "lambda_config" {
    for_each = each.value.lambda != null ? [each.value.lambda] : []
    content {
      function_arn = lambda_config.value.function_arn
    }
  }

  dynamic "http_config" {
    for_each = each.value.http != null ? [each.value.http] : []
    content {
      endpoint = http_config.value.endpoint

      # AWS_IAM is the only authorization type - pinned here, never
      # spec surface; the sigv4 block's presence selects signing.
      dynamic "authorization_config" {
        for_each = http_config.value.sigv4 != null ? [http_config.value.sigv4] : []
        content {
          authorization_type = "AWS_IAM"
          aws_iam_config {
            signing_region       = authorization_config.value.signing_region != "" ? authorization_config.value.signing_region : null
            signing_service_name = authorization_config.value.signing_service_name != "" ? authorization_config.value.signing_service_name : null
          }
        }
      }
    }
  }

  dynamic "opensearchservice_config" {
    for_each = each.value.opensearch != null ? [each.value.opensearch] : []
    content {
      endpoint = opensearchservice_config.value.endpoint
      region   = opensearchservice_config.value.region != "" ? opensearchservice_config.value.region : null
    }
  }

  dynamic "elasticsearch_config" {
    for_each = each.value.elasticsearch != null ? [each.value.elasticsearch] : []
    content {
      endpoint = elasticsearch_config.value.endpoint
      region   = elasticsearch_config.value.region != "" ? elasticsearch_config.value.region : null
    }
  }

  dynamic "event_bridge_config" {
    for_each = each.value.eventbridge != null ? [each.value.eventbridge] : []
    content {
      event_bus_arn = event_bridge_config.value.event_bus_arn
    }
  }

  dynamic "relational_database_config" {
    for_each = each.value.relational_database != null ? [each.value.relational_database] : []
    content {
      # RDS_HTTP_ENDPOINT is the only source type - pinned here.
      source_type = "RDS_HTTP_ENDPOINT"
      http_endpoint_config {
        db_cluster_identifier = relational_database_config.value.db_cluster_identifier
        aws_secret_store_arn  = relational_database_config.value.aws_secret_store_arn
        database_name         = relational_database_config.value.database_name != "" ? relational_database_config.value.database_name : null
        schema                = relational_database_config.value.schema != "" ? relational_database_config.value.schema : null
        region                = relational_database_config.value.region != "" ? relational_database_config.value.region : null
      }
    }
  }
}

# ---------------------------------------------------------------------------
# GraphQL satellites: types, functions, resolvers, cache.
# ---------------------------------------------------------------------------

resource "aws_appsync_type" "this" {
  for_each = local.types

  api_id     = local.api_id
  definition = each.value.definition
  format     = each.value.format
}

resource "aws_appsync_function" "this" {
  for_each = local.functions

  api_id = local.api_id
  name   = each.value.name

  # In-spec data source names resolve through the created data source
  # (carrying the dependency edge); externally created data sources
  # pass through as literals.
  data_source = contains(keys(local.datasources), each.value.data_source_name) ? aws_appsync_datasource.this[each.value.data_source_name].name : each.value.data_source_name

  description = each.value.description != "" ? each.value.description : null

  code = each.value.code != "" ? each.value.code : null

  dynamic "runtime" {
    for_each = each.value.code != "" ? [true] : []
    content {
      # APPSYNC_JS is the only runtime AWS ships - pinned here.
      name            = "APPSYNC_JS"
      runtime_version = each.value.runtime_version != "" ? each.value.runtime_version : "1.0.0"
    }
  }

  # The VTL arm; functions support only the 2018-05-29 template
  # version, pinned by the provider's own default when templates are
  # used.
  request_mapping_template  = each.value.request_mapping_template != "" ? each.value.request_mapping_template : null
  response_mapping_template = each.value.response_mapping_template != "" ? each.value.response_mapping_template : null

  max_batch_size = each.value.max_batch_size > 0 ? each.value.max_batch_size : null

  dynamic "sync_config" {
    for_each = each.value.sync_config != null ? [each.value.sync_config] : []
    content {
      conflict_detection = sync_config.value.conflict_detection != "" ? sync_config.value.conflict_detection : null
      conflict_handler   = sync_config.value.conflict_handler != "" ? sync_config.value.conflict_handler : null

      dynamic "lambda_conflict_handler_config" {
        for_each = sync_config.value.lambda_conflict_handler_arn != "" ? [sync_config.value.lambda_conflict_handler_arn] : []
        content {
          lambda_conflict_handler_arn = lambda_conflict_handler_config.value
        }
      }
    }
  }
}

resource "aws_appsync_resolver" "this" {
  for_each = local.resolvers

  api_id = local.api_id
  type   = each.value.type
  field  = each.value.field

  # UNIT XOR PIPELINE, derived from which arm the entry carries (the
  # spec's CEL wall keeps it exactly one).
  kind = length(each.value.pipeline_functions) > 0 ? "PIPELINE" : "UNIT"

  data_source = each.value.data_source_name != "" ? (contains(keys(local.datasources), each.value.data_source_name) ? aws_appsync_datasource.this[each.value.data_source_name].name : each.value.data_source_name) : null

  dynamic "pipeline_config" {
    for_each = length(each.value.pipeline_functions) > 0 ? [each.value.pipeline_functions] : []
    content {
      # Pipeline entries name spec functions; the module joins names
      # to AWS function ids (externally created functions pass their
      # ids through as literals).
      functions = [for f in pipeline_config.value : contains(keys(local.functions), f) ? aws_appsync_function.this[f].function_id : f]
    }
  }

  code = each.value.code != "" ? each.value.code : null

  dynamic "runtime" {
    for_each = each.value.code != "" ? [true] : []
    content {
      name            = "APPSYNC_JS"
      runtime_version = each.value.runtime_version != "" ? each.value.runtime_version : "1.0.0"
    }
  }

  request_template  = each.value.request_template != "" ? each.value.request_template : null
  response_template = each.value.response_template != "" ? each.value.response_template : null

  max_batch_size = each.value.max_batch_size > 0 ? each.value.max_batch_size : null

  dynamic "caching_config" {
    for_each = each.value.caching != null ? [each.value.caching] : []
    content {
      caching_keys = length(caching_config.value.caching_keys) > 0 ? caching_config.value.caching_keys : null
      ttl          = caching_config.value.ttl > 0 ? caching_config.value.ttl : null
    }
  }

  dynamic "sync_config" {
    for_each = each.value.sync_config != null ? [each.value.sync_config] : []
    content {
      conflict_detection = sync_config.value.conflict_detection != "" ? sync_config.value.conflict_detection : null
      conflict_handler   = sync_config.value.conflict_handler != "" ? sync_config.value.conflict_handler : null

      dynamic "lambda_conflict_handler_config" {
        for_each = sync_config.value.lambda_conflict_handler_arn != "" ? [sync_config.value.lambda_conflict_handler_arn] : []
        content {
          lambda_conflict_handler_arn = lambda_conflict_handler_config.value
        }
      }
    }
  }
}

# The cache singleton: its AWS id IS the API id (one cache per API by
# AWS's own model). Both encryption flags replace the cache on change.
resource "aws_appsync_api_cache" "this" {
  count = local.is_graphql && var.spec.graphql.cache != null ? 1 : 0

  api_id = local.api_id

  api_caching_behavior = var.spec.graphql.cache.api_caching_behavior
  ttl                  = var.spec.graphql.cache.ttl
  type                 = var.spec.graphql.cache.type

  at_rest_encryption_enabled = var.spec.graphql.cache.at_rest_encryption_enabled
  transit_encryption_enabled = var.spec.graphql.cache.transit_encryption_enabled
}

# ---------------------------------------------------------------------------
# API keys (either arm).
# ---------------------------------------------------------------------------

resource "aws_appsync_api_key" "this" {
  for_each = local.api_keys

  api_id = local.api_id

  description = each.value.description != "" ? each.value.description : null
  expires     = each.value.expires != "" ? each.value.expires : null
}

# ---------------------------------------------------------------------------
# Channel namespaces (the Events arm).
# ---------------------------------------------------------------------------

resource "aws_appsync_channel_namespace" "this" {
  for_each = local.channel_namespaces

  api_id = local.api_id
  name   = each.value.name

  code_handlers = each.value.code_handlers != "" ? each.value.code_handlers : null

  dynamic "publish_auth_mode" {
    for_each = each.value.publish_auth_modes
    content {
      auth_type = publish_auth_mode.value
    }
  }

  dynamic "subscribe_auth_mode" {
    for_each = each.value.subscribe_auth_modes
    content {
      auth_type = subscribe_auth_mode.value
    }
  }

  dynamic "handler_configs" {
    for_each = each.value.handler_configs != null ? [each.value.handler_configs] : []
    content {
      dynamic "on_publish" {
        for_each = handler_configs.value.on_publish != null ? [handler_configs.value.on_publish] : []
        content {
          behavior = on_publish.value.behavior
          integration {
            # In-spec data source names resolve through the created
            # data source (the dependency edge).
            data_source_name = contains(keys(local.datasources), on_publish.value.data_source_name) ? aws_appsync_datasource.this[on_publish.value.data_source_name].name : on_publish.value.data_source_name

            dynamic "lambda_config" {
              for_each = on_publish.value.lambda_invoke_type != "" ? [on_publish.value.lambda_invoke_type] : []
              content {
                invoke_type = lambda_config.value
              }
            }
          }
        }
      }

      dynamic "on_subscribe" {
        for_each = handler_configs.value.on_subscribe != null ? [handler_configs.value.on_subscribe] : []
        content {
          behavior = on_subscribe.value.behavior
          integration {
            data_source_name = contains(keys(local.datasources), on_subscribe.value.data_source_name) ? aws_appsync_datasource.this[on_subscribe.value.data_source_name].name : on_subscribe.value.data_source_name

            dynamic "lambda_config" {
              for_each = on_subscribe.value.lambda_invoke_type != "" ? [on_subscribe.value.lambda_invoke_type] : []
              content {
                invoke_type = lambda_config.value
              }
            }
          }
        }
      }
    }
  }

  tags = local.aws_tags
}

# ---------------------------------------------------------------------------
# The custom domain (either arm; the certificate must be in us-east-1).
# ---------------------------------------------------------------------------

resource "aws_appsync_domain_name" "this" {
  count = var.spec.custom_domain != null ? 1 : 0

  domain_name     = var.spec.custom_domain.domain_name
  certificate_arn = var.spec.custom_domain.certificate_arn
  description     = var.spec.custom_domain.description != "" ? var.spec.custom_domain.description : null
}

# The association is 1:1 with the domain (its AWS id IS the domain
# name) and re-points in place; create/delete wait up to 60 minutes
# upstream.
resource "aws_appsync_domain_name_api_association" "this" {
  count = var.spec.custom_domain != null ? 1 : 0

  domain_name = aws_appsync_domain_name.this[0].domain_name
  api_id      = local.api_id
}

# ---------------------------------------------------------------------------
# Source-API associations (MERGED APIs).
# ---------------------------------------------------------------------------

resource "aws_appsync_source_api_association" "this" {
  for_each = local.source_apis

  merged_api_id = local.api_id
  source_api_id = each.value.source_api_id

  description = each.value.description != "" ? each.value.description : null

  # A Plugin Framework list ATTRIBUTE (not a block) upstream - assigned
  # as a value, one entry at most.
  source_api_association_config = each.value.merge_type != "" ? [{ merge_type = each.value.merge_type }] : null
}
