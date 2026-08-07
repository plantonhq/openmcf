# AWS API Gateway HTTP API (API Gateway v2).
#
# One spec materializes as up to five resource groups wired together: the API,
# its stage, deduplicated integrations, optional authorizers, and the routes
# that bind them. Locals.tf normalizes presence semantics (stage defaults,
# whole-object integration dedup, authorizer-by-name addressing).

# ---------------------------------------------------------------------------
# 1. The HTTP API
# ---------------------------------------------------------------------------

resource "aws_apigatewayv2_api" "this" {
  name          = local.api_name
  protocol_type = "HTTP"
  description   = var.spec.description != "" ? var.spec.description : null

  # Informational version label surfaced in the console and OpenAPI exports.
  version = var.spec.api_version != "" ? var.spec.api_version : null

  # When a custom domain (AwsHttpApiDomain) fronts this API, disabling the
  # default endpoint stops callers from bypassing the domain's TLS policy /
  # mTLS / WAF via https://{api-id}.execute-api...
  disable_execute_api_endpoint = var.spec.disable_execute_api_endpoint

  # ipv4 or dualstack; AWS defaults new APIs to dualstack when unset.
  ip_address_type = var.spec.ip_address_type != "" ? var.spec.ip_address_type : null

  # CORS is enforced by API Gateway itself for HTTP APIs -- configuring it
  # here means the backend never needs CORS logic.
  dynamic "cors_configuration" {
    for_each = var.spec.cors_configuration != null ? [var.spec.cors_configuration] : []
    content {
      allow_origins     = cors_configuration.value.allow_origins
      allow_methods     = cors_configuration.value.allow_methods
      allow_headers     = cors_configuration.value.allow_headers
      expose_headers    = cors_configuration.value.expose_headers
      max_age           = cors_configuration.value.max_age_seconds > 0 ? cors_configuration.value.max_age_seconds : null
      allow_credentials = cors_configuration.value.allow_credentials
    }
  }

  tags = local.aws_tags
}

# ---------------------------------------------------------------------------
# 2. The stage
# ---------------------------------------------------------------------------

resource "aws_apigatewayv2_stage" "this" {
  api_id      = aws_apigatewayv2_api.this.id
  name        = local.stage_name
  auto_deploy = local.auto_deploy
  description = local.stage_config != null && try(local.stage_config.description, "") != "" ? local.stage_config.description : null

  # Access logging streams request records to CloudWatch Logs in the
  # user-supplied format. The log group is referenced, never created here --
  # log groups are their own composable resource (AwsCloudwatchLogGroup).
  dynamic "access_log_settings" {
    for_each = try(local.stage_config.access_log, null) != null ? [local.stage_config.access_log] : []
    content {
      destination_arn = access_log_settings.value.destination_arn
      format          = access_log_settings.value.format
    }
  }

  # Stage-wide defaults. Rendered when any default is meaningful: a throttle
  # limit or detailed metrics. (data_trace_enabled and logging_level are
  # WebSocket-only knobs and are deliberately not modeled for HTTP APIs.)
  dynamic "default_route_settings" {
    for_each = (
      local.stage_config != null && (
        try(local.stage_config.default_throttle.burst_limit, 0) > 0 ||
        try(local.stage_config.default_throttle.rate_limit, 0) > 0 ||
        try(local.stage_config.detailed_metrics_enabled, false)
      )
    ) ? [1] : []
    content {
      throttling_burst_limit   = try(local.stage_config.default_throttle.burst_limit, 0) > 0 ? local.stage_config.default_throttle.burst_limit : null
      throttling_rate_limit    = try(local.stage_config.default_throttle.rate_limit, 0) > 0 ? local.stage_config.default_throttle.rate_limit : null
      detailed_metrics_enabled = try(local.stage_config.detailed_metrics_enabled, false)
    }
  }

  # Per-route overrides. Zero-valued limits inherit the stage default -- only
  # real overrides are sent.
  dynamic "route_settings" {
    for_each = local.stage_config != null ? local.stage_config.route_settings : []
    content {
      route_key                = route_settings.value.route_key
      throttling_burst_limit   = route_settings.value.throttling_burst_limit > 0 ? route_settings.value.throttling_burst_limit : null
      throttling_rate_limit    = route_settings.value.throttling_rate_limit > 0 ? route_settings.value.throttling_rate_limit : null
      detailed_metrics_enabled = route_settings.value.detailed_metrics_enabled
    }
  }

  stage_variables = local.stage_config != null ? local.stage_config.stage_variables : {}

  tags = local.aws_tags

  # Route settings reference routes by key; make sure the routes exist before
  # the stage tries to attach settings to them.
  depends_on = [aws_apigatewayv2_route.this]
}

# ---------------------------------------------------------------------------
# 3. Integrations (deduplicated)
# ---------------------------------------------------------------------------

resource "aws_apigatewayv2_integration" "this" {
  # each.value is the list of identical integration objects that hashed to
  # this key (locals.tf groups them); [0] is the canonical instance.
  for_each = local.integration_map

  api_id           = aws_apigatewayv2_api.this.id
  integration_type = each.value[0].integration_type

  # Proxy integrations carry their target here; AWS service integrations
  # (integration_subtype) express the target in request_parameters and AWS
  # rejects a URI alongside a subtype -- the spec CEL enforces the split.
  integration_uri     = each.value[0].integration_uri != "" ? each.value[0].integration_uri : null
  integration_subtype = each.value[0].integration_subtype != "" ? each.value[0].integration_subtype : null

  # Only AWS_PROXY (Lambda) and AWS service subtypes carry a payload format.
  # HTTP_PROXY rejects any PayloadFormatVersion (AWS BadRequestException on
  # 2.0); omit the field for HTTP_PROXY / VPC_LINK proxy integrations.
  payload_format_version = (
    each.value[0].integration_subtype != "" ? "1.0" :
    each.value[0].integration_type == "AWS_PROXY" ? coalesce(each.value[0].payload_format_version, "2.0") :
    null
  )

  # Lambda integrations are always invoked with POST regardless of this value.
  integration_method = each.value[0].integration_method != "" ? each.value[0].integration_method : null

  timeout_milliseconds = each.value[0].timeout_milliseconds > 0 ? each.value[0].timeout_milliseconds : null

  # Private integrations reach through a VPC link; the spec CELs guarantee
  # VPC_LINK <=> connection_id and HTTP_PROXY-only.
  connection_type = each.value[0].connection_type != "" ? each.value[0].connection_type : null
  connection_id   = each.value[0].connection_id != "" ? each.value[0].connection_id : null

  # The role API Gateway assumes to call an AWS service action (required for
  # subtype integrations by spec CEL); Lambda proxies normally rely on the
  # function's resource policy instead.
  credentials_arn = each.value[0].credentials_arn != "" ? each.value[0].credentials_arn : null

  # Parameter mappings (proxy) or service-action parameters (subtype).
  request_parameters = length(each.value[0].request_parameters) > 0 ? each.value[0].request_parameters : null

  # Response transforms keyed by backend status code.
  dynamic "response_parameters" {
    for_each = each.value[0].response_parameters
    content {
      status_code = response_parameters.value.status_code
      mappings    = response_parameters.value.mappings
    }
  }

  # SNI override for private integrations whose internal ALB serves a
  # public-domain certificate.
  dynamic "tls_config" {
    for_each = each.value[0].tls_server_name_to_verify != "" ? [1] : []
    content {
      server_name_to_verify = each.value[0].tls_server_name_to_verify
    }
  }

  description = each.value[0].description != "" ? each.value[0].description : null
}

# ---------------------------------------------------------------------------
# 4. Authorizers
# ---------------------------------------------------------------------------

resource "aws_apigatewayv2_authorizer" "this" {
  for_each = local.authorizer_map

  api_id          = aws_apigatewayv2_api.this.id
  name            = each.value.name
  authorizer_type = each.value.authorizer_type

  identity_sources = length(each.value.identity_sources) > 0 ? each.value.identity_sources : null

  # JWT: validate iss/aud claims against the configured provider.
  dynamic "jwt_configuration" {
    for_each = each.value.authorizer_type == "JWT" && each.value.jwt_configuration != null ? [each.value.jwt_configuration] : []
    content {
      issuer   = jwt_configuration.value.issuer
      audience = jwt_configuration.value.audiences
    }
  }

  # REQUEST: a Lambda decides. Simple responses ({"isAuthorized": bool}) are
  # the modern contract; payload format 2.0 pairs with them.
  authorizer_uri                    = each.value.authorizer_type == "REQUEST" && each.value.authorizer_uri != "" ? each.value.authorizer_uri : null
  authorizer_credentials_arn        = each.value.authorizer_type == "REQUEST" && each.value.authorizer_credentials_arn != "" ? each.value.authorizer_credentials_arn : null
  enable_simple_responses           = each.value.authorizer_type == "REQUEST" && each.value.enable_simple_responses ? true : null
  authorizer_payload_format_version = each.value.authorizer_type == "REQUEST" && each.value.authorizer_payload_format_version != "" ? each.value.authorizer_payload_format_version : null

  authorizer_result_ttl_in_seconds = each.value.result_ttl_seconds > 0 ? each.value.result_ttl_seconds : null
}

# ---------------------------------------------------------------------------
# 5. Routes
# ---------------------------------------------------------------------------

resource "aws_apigatewayv2_route" "this" {
  for_each = {
    for idx, route in var.spec.routes : tostring(idx) => route
  }

  api_id    = aws_apigatewayv2_api.this.id
  route_key = each.value.route_key

  # Every route targets its (deduplicated) integration.
  target = "integrations/${aws_apigatewayv2_integration.this[local.integration_keys[each.key]].id}"

  # NONE is the AWS default; only real authorization modes are sent.
  authorization_type = each.value.authorization_type != "" && each.value.authorization_type != "NONE" ? each.value.authorization_type : null

  # JWT routes bind JWT authorizers; CUSTOM routes bind REQUEST (Lambda)
  # authorizers -- the spec CEL guarantees the referenced authorizer's type
  # matches.
  authorizer_id = contains(["JWT", "CUSTOM"], each.value.authorization_type) && each.value.authorizer_name != "" ? aws_apigatewayv2_authorizer.this[each.value.authorizer_name].id : null

  authorization_scopes = each.value.authorization_type == "JWT" && length(each.value.authorization_scopes) > 0 ? each.value.authorization_scopes : null

  # Stable operationId for OpenAPI exports and generated clients.
  operation_name = each.value.operation_name != "" ? each.value.operation_name : null
}
