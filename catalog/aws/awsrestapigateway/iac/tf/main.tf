# A full REST API (API Gateway v1): the API, its derived resource/
# method tree with inline integrations (or an imported OpenAPI
# document), the API-scoped satellites, and one stage with an explicit
# hash-triggered deployment.
#
# Lifecycle facts the renders below depend on:
#   - REST APIs deploy by explicit snapshot: the deployment's trigger
#     carries the definition hash from locals, so ANY definition change
#     rolls a new deployment (create_before_destroy - the stage
#     repoints, then the old deployment deletes);
#   - the resource tree renders level-by-level (a for_each resource
#     cannot reference itself); the spec caps paths at five segments;
#   - method responses are serialized upstream behind a global mutex
#     (concurrent writes on one API conflict) with a 2-minute retry;
#   - an integration response requires BOTH its method response and its
#     integration to exist - the depends_on below carries that edge;
#   - with an OpenAPI body, the provider runs CreateRestApi ->
#     PutRestApi -> a reconciliation pass that re-applies configured
#     literals the overwrite-mode import wiped - expected apply-log
#     noise, not drift;
#   - minimum_compression_size is the provider's nullable-int-as-string
#     quirk: unset means compression disabled, "0" compresses
#     everything;
#   - the standalone rest_api_policy resource owns the policy (clean
#     PATCH updates; delete resets to empty instead of touching the
#     API), so the API's own policy argument stays unset;
#   - enabling or resizing the stage cache waits up to 90 minutes
#     upstream while AWS provisions the cluster;
#   - method settings are PATCHes on the stage (a view, not a real
#     object).

resource "aws_api_gateway_rest_api" "this" {
  # metadata.name is the naming basis on both engines.
  name        = var.metadata.name
  description = var.spec.description != "" ? var.spec.description : null

  api_key_source     = var.spec.api_key_source != "" ? var.spec.api_key_source : null
  binary_media_types = length(var.spec.binary_media_types) > 0 ? var.spec.binary_media_types : null

  minimum_compression_size = var.spec.minimum_compression_size != null ? tostring(var.spec.minimum_compression_size) : null

  disable_execute_api_endpoint = var.spec.disable_execute_api_endpoint

  dynamic "endpoint_configuration" {
    for_each = var.spec.endpoint_configuration != null ? [var.spec.endpoint_configuration] : []
    content {
      types            = [endpoint_configuration.value.type]
      ip_address_type  = endpoint_configuration.value.ip_address_type != "" ? endpoint_configuration.value.ip_address_type : null
      vpc_endpoint_ids = length(endpoint_configuration.value.vpc_endpoint_ids) > 0 ? endpoint_configuration.value.vpc_endpoint_ids : null
    }
  }

  endpoint_access_mode = var.spec.endpoint_access_mode != "" ? var.spec.endpoint_access_mode : null
  security_policy      = var.spec.security_policy != "" ? var.spec.security_policy : null

  # The OpenAPI arm: the document IS the API definition.
  body              = var.spec.openapi != null ? var.spec.openapi.body : null
  fail_on_warnings  = var.spec.openapi != null ? var.spec.openapi.fail_on_warnings : null
  parameters        = var.spec.openapi != null && length(try(var.spec.openapi.parameters, {})) > 0 ? var.spec.openapi.parameters : null
  put_rest_api_mode = var.spec.openapi != null && try(var.spec.openapi.mode, "") != "" ? var.spec.openapi.mode : null

  tags = local.aws_tags
}

# The resource policy as its own object: clean PATCH updates, and
# delete resets the policy instead of touching the API.
resource "aws_api_gateway_rest_api_policy" "this" {
  count = var.spec.policy != null ? 1 : 0

  rest_api_id = aws_api_gateway_rest_api.this.id
  policy      = jsonencode(var.spec.policy)
}

# ---------------------------------------------------------------------------
# The derived resource tree, level by level (parents before children).
# ---------------------------------------------------------------------------

resource "aws_api_gateway_resource" "level1" {
  for_each = local.paths_level1

  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_rest_api.this.root_resource_id
  path_part   = element(split("/", trimprefix(each.value, "/")), 0)
}

resource "aws_api_gateway_resource" "level2" {
  for_each = local.paths_level2

  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_resource.level1["/${join("/", slice(split("/", trimprefix(each.value, "/")), 0, 1))}"].id
  path_part   = element(split("/", trimprefix(each.value, "/")), 1)
}

resource "aws_api_gateway_resource" "level3" {
  for_each = local.paths_level3

  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_resource.level2["/${join("/", slice(split("/", trimprefix(each.value, "/")), 0, 2))}"].id
  path_part   = element(split("/", trimprefix(each.value, "/")), 2)
}

resource "aws_api_gateway_resource" "level4" {
  for_each = local.paths_level4

  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_resource.level3["/${join("/", slice(split("/", trimprefix(each.value, "/")), 0, 3))}"].id
  path_part   = element(split("/", trimprefix(each.value, "/")), 3)
}

resource "aws_api_gateway_resource" "level5" {
  for_each = local.paths_level5

  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_resource.level4["/${join("/", slice(split("/", trimprefix(each.value, "/")), 0, 4))}"].id
  path_part   = element(split("/", trimprefix(each.value, "/")), 4)
}

# ---------------------------------------------------------------------------
# Named satellites referenced by routes.
# ---------------------------------------------------------------------------

resource "aws_api_gateway_model" "this" {
  for_each = local.models

  rest_api_id  = aws_api_gateway_rest_api.this.id
  name         = each.value.name
  content_type = each.value.content_type
  description  = each.value.description != "" ? each.value.description : null
  schema       = each.value.schema != "" ? each.value.schema : null
}

resource "aws_api_gateway_request_validator" "this" {
  for_each = local.request_validators

  rest_api_id                 = aws_api_gateway_rest_api.this.id
  name                        = each.value.name
  validate_request_body       = each.value.validate_request_body
  validate_request_parameters = each.value.validate_request_parameters
}

resource "aws_api_gateway_authorizer" "this" {
  for_each = local.authorizers

  rest_api_id = aws_api_gateway_rest_api.this.id
  name        = each.value.name
  type        = each.value.type

  authorizer_uri         = each.value.lambda_invoke_uri != "" ? each.value.lambda_invoke_uri : null
  authorizer_credentials = each.value.credentials_arn != "" ? each.value.credentials_arn : null
  provider_arns          = length(each.value.provider_arns) > 0 ? each.value.provider_arns : null

  identity_source                = each.value.identity_source != "" ? each.value.identity_source : null
  identity_validation_expression = each.value.identity_validation_expression != "" ? each.value.identity_validation_expression : null

  # Presence-typed: an explicit 0 disables caching; unset keeps AWS's
  # 300-second default.
  authorizer_result_ttl_in_seconds = each.value.result_ttl_seconds
}

resource "aws_api_gateway_gateway_response" "this" {
  for_each = local.gateway_responses

  rest_api_id   = aws_api_gateway_rest_api.this.id
  response_type = each.value.response_type

  status_code         = each.value.status_code != "" ? each.value.status_code : null
  response_parameters = length(each.value.response_parameters) > 0 ? each.value.response_parameters : null
  response_templates  = length(each.value.response_templates) > 0 ? each.value.response_templates : null
}

# ---------------------------------------------------------------------------
# Methods, integrations, and responses per route.
# ---------------------------------------------------------------------------

resource "aws_api_gateway_method" "this" {
  for_each = local.routes

  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = each.value.path == "/" ? aws_api_gateway_rest_api.this.root_resource_id : local.resource_id_by_path[each.value.path]
  http_method = each.value.method

  authorization        = each.value.authorization != "" ? each.value.authorization : "NONE"
  authorizer_id        = each.value.authorizer_name != "" ? aws_api_gateway_authorizer.this[each.value.authorizer_name].id : null
  authorization_scopes = length(each.value.authorization_scopes) > 0 ? each.value.authorization_scopes : null

  api_key_required = each.value.api_key_required
  operation_name   = each.value.operation_name != "" ? each.value.operation_name : null

  request_parameters = length(each.value.request_parameters) > 0 ? each.value.request_parameters : null

  # In-spec model names resolve through the created model (carrying the
  # dependency edge); the AWS built-ins pass through as literals.
  request_models = length(each.value.request_models) > 0 ? {
    for ct, m in each.value.request_models : ct => contains(keys(local.models), m) ? aws_api_gateway_model.this[m].name : m
  } : null

  request_validator_id = each.value.request_validator_name != "" ? aws_api_gateway_request_validator.this[each.value.request_validator_name].id : null
}

resource "aws_api_gateway_integration" "this" {
  for_each = local.routes

  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = each.value.path == "/" ? aws_api_gateway_rest_api.this.root_resource_id : local.resource_id_by_path[each.value.path]
  http_method = aws_api_gateway_method.this[each.key].http_method

  type = each.value.integration.type
  uri  = each.value.integration.uri != "" ? each.value.integration.uri : null

  # Lambda invocations are always POST; the spec's validation allows
  # omitting it on AWS_PROXY and the module fills it.
  integration_http_method = each.value.integration.http_method != "" ? each.value.integration.http_method : (each.value.integration.type == "AWS_PROXY" ? "POST" : null)

  credentials = each.value.integration.credentials_arn != "" ? each.value.integration.credentials_arn : null

  connection_type = each.value.integration.connection_type != "" ? each.value.integration.connection_type : null
  connection_id   = each.value.integration.vpc_link_id != "" ? each.value.integration.vpc_link_id : null

  passthrough_behavior = each.value.integration.passthrough_behavior != "" ? each.value.integration.passthrough_behavior : null
  content_handling     = each.value.integration.content_handling != "" ? each.value.integration.content_handling : null

  cache_key_parameters = length(each.value.integration.cache_key_parameters) > 0 ? each.value.integration.cache_key_parameters : null
  cache_namespace      = each.value.integration.cache_namespace != "" ? each.value.integration.cache_namespace : null

  request_parameters = length(each.value.integration.request_parameters) > 0 ? each.value.integration.request_parameters : null
  request_templates  = length(each.value.integration.request_templates) > 0 ? each.value.integration.request_templates : null

  timeout_milliseconds   = each.value.integration.timeout_milliseconds > 0 ? each.value.integration.timeout_milliseconds : null
  response_transfer_mode = each.value.integration.response_transfer_mode != "" ? each.value.integration.response_transfer_mode : null

  dynamic "tls_config" {
    for_each = each.value.integration.tls_insecure_skip_verification ? [true] : []
    content {
      insecure_skip_verification = true
    }
  }
}

resource "aws_api_gateway_method_response" "this" {
  for_each = local.route_responses

  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = each.value.route.path == "/" ? aws_api_gateway_rest_api.this.root_resource_id : local.resource_id_by_path[each.value.route.path]
  http_method = aws_api_gateway_method.this[each.value.route_key].http_method
  status_code = each.value.response.status_code

  response_models = length(each.value.response.response_models) > 0 ? {
    for ct, m in each.value.response.response_models : ct => contains(keys(local.models), m) ? aws_api_gateway_model.this[m].name : m
  } : null

  response_parameters = length(each.value.response.response_parameters) > 0 ? each.value.response.response_parameters : null
}

resource "aws_api_gateway_integration_response" "this" {
  for_each = local.route_responses

  rest_api_id = aws_api_gateway_rest_api.this.id
  resource_id = each.value.route.path == "/" ? aws_api_gateway_rest_api.this.root_resource_id : local.resource_id_by_path[each.value.route.path]
  http_method = aws_api_gateway_method.this[each.value.route_key].http_method
  status_code = aws_api_gateway_method_response.this[each.key].status_code

  selection_pattern   = each.value.response.selection_pattern != "" ? each.value.response.selection_pattern : null
  response_parameters = length(each.value.response.integration_response_parameters) > 0 ? each.value.response.integration_response_parameters : null
  response_templates  = length(each.value.response.integration_response_templates) > 0 ? each.value.response.integration_response_templates : null
  content_handling    = each.value.response.content_handling != "" ? each.value.response.content_handling : null

  # AWS requires the integration to exist too (the method response edge
  # already rides the status_code reference above).
  depends_on = [aws_api_gateway_integration.this]
}

# ---------------------------------------------------------------------------
# Documentation.
# ---------------------------------------------------------------------------

resource "aws_api_gateway_documentation_part" "this" {
  for_each = local.documentation_parts

  rest_api_id = aws_api_gateway_rest_api.this.id
  properties  = each.value.properties

  location {
    type        = each.value.location.type
    path        = each.value.location.path != "" ? each.value.location.path : null
    method      = each.value.location.method != "" ? each.value.location.method : null
    name        = each.value.location.name != "" ? each.value.location.name : null
    status_code = each.value.location.status_code != "" ? each.value.location.status_code : null
  }
}

# Publishing snapshots the parts - it must run after them.
resource "aws_api_gateway_documentation_version" "this" {
  count = var.spec.documentation != null && var.spec.documentation.published_version != null ? 1 : 0

  rest_api_id = aws_api_gateway_rest_api.this.id
  version     = var.spec.documentation.published_version.version
  description = var.spec.documentation.published_version.description != "" ? var.spec.documentation.published_version.description : null

  depends_on = [aws_api_gateway_documentation_part.this]
}

# ---------------------------------------------------------------------------
# Deployment and stage.
# ---------------------------------------------------------------------------

# The TLS client certificate presented to HTTP backends (generated
# arm; an existing certificate is referenced by ID on the stage).
resource "aws_api_gateway_client_certificate" "this" {
  count = var.spec.stage != null && var.spec.stage.client_certificate != null && var.spec.stage.client_certificate.generate ? 1 : 0

  description = var.spec.stage.client_certificate.description != "" ? var.spec.stage.client_certificate.description : null

  tags = local.aws_tags
}

resource "aws_api_gateway_deployment" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id

  triggers = {
    redeployment = local.definition_hash
  }

  # The snapshot must capture the complete definition; replacement
  # creates the new deployment before the old one deletes.
  depends_on = [
    aws_api_gateway_method.this,
    aws_api_gateway_integration.this,
    aws_api_gateway_method_response.this,
    aws_api_gateway_integration_response.this,
    aws_api_gateway_gateway_response.this,
    aws_api_gateway_rest_api_policy.this,
  ]

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "this" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  deployment_id = aws_api_gateway_deployment.this.id
  stage_name    = local.stage_name

  description = var.spec.stage != null && try(var.spec.stage.description, "") != "" ? var.spec.stage.description : null
  variables   = var.spec.stage != null && length(try(var.spec.stage.stage_variables, {})) > 0 ? var.spec.stage.stage_variables : null

  xray_tracing_enabled = var.spec.stage != null ? var.spec.stage.xray_tracing_enabled : false

  cache_cluster_enabled = var.spec.stage != null && var.spec.stage.cache_cluster != null ? var.spec.stage.cache_cluster.enabled : false
  cache_cluster_size    = var.spec.stage != null && var.spec.stage.cache_cluster != null && var.spec.stage.cache_cluster.enabled ? var.spec.stage.cache_cluster.size : null

  client_certificate_id = var.spec.stage != null && var.spec.stage.client_certificate != null ? (var.spec.stage.client_certificate.generate ? aws_api_gateway_client_certificate.this[0].id : var.spec.stage.client_certificate.existing_certificate_id) : null

  dynamic "access_log_settings" {
    for_each = var.spec.stage != null && var.spec.stage.access_log != null ? [var.spec.stage.access_log] : []
    content {
      destination_arn = access_log_settings.value.destination_arn
      format          = access_log_settings.value.format
    }
  }

  documentation_version = var.spec.stage != null && try(var.spec.stage.documentation_version, "") != "" ? var.spec.stage.documentation_version : null

  tags = local.aws_tags

  depends_on = [aws_api_gateway_documentation_version.this]
}

# Per-method overrides of logging/metrics/throttling/caching.
resource "aws_api_gateway_method_settings" "this" {
  for_each = local.method_settings

  rest_api_id = aws_api_gateway_rest_api.this.id
  stage_name  = aws_api_gateway_stage.this.stage_name
  method_path = each.value.method_path

  settings {
    metrics_enabled                            = each.value.metrics_enabled
    logging_level                              = each.value.logging_level != "" ? each.value.logging_level : null
    data_trace_enabled                         = each.value.data_trace_enabled
    throttling_burst_limit                     = each.value.throttling_burst_limit
    throttling_rate_limit                      = each.value.throttling_rate_limit
    caching_enabled                            = each.value.caching_enabled
    cache_ttl_in_seconds                       = each.value.cache_ttl_in_seconds
    cache_data_encrypted                       = each.value.cache_data_encrypted
    require_authorization_for_cache_control    = each.value.require_authorization_for_cache_control
    unauthorized_cache_control_header_strategy = each.value.unauthorized_cache_control_header_strategy != "" ? each.value.unauthorized_cache_control_header_strategy : null
  }
}
