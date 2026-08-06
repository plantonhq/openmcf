# Origin Access Controls created for S3 origins that asked for one
# (s3_origin.create_origin_access_control). OAC is the modern way to serve a
# private bucket: CloudFront signs origin requests with SigV4, so the bucket
# stays fully private and only this distribution's ARN (allowed in the bucket
# policy) can read it. One OAC per requesting origin, named after the origin
# so multi-origin distributions stay legible in the console.
resource "aws_cloudfront_origin_access_control" "this" {
  for_each = local.oac_origins

  name                              = "${var.metadata.name}-${each.key}"
  description                       = "Origin access control for the ${each.key} origin of the ${var.metadata.name} distribution"
  origin_access_control_origin_type = "s3"
  signing_behavior                  = "always"
  signing_protocol                  = "sigv4"
}

resource "aws_cloudfront_distribution" "this" {
  # enabled/wait_for_deployment resolve to their annotated defaults (true) at
  # manifest load. A disabled distribution stays deployed-but-dark -- also the
  # state AWS requires before deletion, which the provider handles on destroy.
  enabled             = var.spec.enabled
  wait_for_deployment = var.spec.wait_for_deployment
  retain_on_delete    = var.spec.retain_on_delete

  aliases             = var.spec.aliases
  comment             = var.spec.comment != "" ? var.spec.comment : null
  default_root_object = var.spec.default_root_object != "" ? var.spec.default_root_object : null
  http_version        = var.spec.http_version != "" ? var.spec.http_version : null
  is_ipv6_enabled     = var.spec.is_ipv6_enabled
  price_class         = var.spec.price_class != "" ? var.spec.price_class : null
  web_acl_id          = local.web_acl_arn != "" ? local.web_acl_arn : null

  # --- Origins ------------------------------------------------------------
  # Exactly one origin-type block renders per origin (CEL enforces the arm
  # exclusivity). An origin with no arm at all is a plain public S3 REST
  # origin, which the provider expresses as an empty s3_origin_config.
  dynamic "origin" {
    for_each = local.origins_by_id
    content {
      origin_id   = origin.key
      domain_name = origin.value.domain_name
      origin_path = origin.value.origin_path != "" ? origin.value.origin_path : null

      # Zero keeps the AWS defaults (3 attempts / 10 seconds).
      connection_attempts = origin.value.connection_attempts != 0 ? origin.value.connection_attempts : null
      connection_timeout  = origin.value.connection_timeout_seconds != 0 ? origin.value.connection_timeout_seconds : null

      # The OAC attaches at the origin level: either the one this module
      # created for the origin, or an externally shared one by ID.
      origin_access_control_id = (
        try(origin.value.s3_origin.create_origin_access_control, false)
        ? aws_cloudfront_origin_access_control.this[origin.key].id
        : (
          try(origin.value.s3_origin.origin_access_control_id, "") != ""
          ? origin.value.s3_origin.origin_access_control_id
          : null
        )
      )

      dynamic "custom_header" {
        for_each = origin.value.custom_headers
        content {
          name  = custom_header.value.name
          value = custom_header.value.value
        }
      }

      dynamic "origin_shield" {
        for_each = origin.value.origin_shield != null ? [origin.value.origin_shield] : []
        content {
          enabled              = true
          origin_shield_region = origin_shield.value.origin_shield_region
        }
      }

      # Legacy OAI path: only for S3 origins already wired to an existing
      # identity -- this module never creates one (OAC supersedes it). A
      # bare s3_origin with no access arm renders the empty identity, the
      # provider's public-bucket shape.
      dynamic "s3_origin_config" {
        for_each = (
          try(origin.value.s3_origin, null) != null &&
          !try(origin.value.s3_origin.create_origin_access_control, false) &&
          try(origin.value.s3_origin.origin_access_control_id, "") == ""
        ) ? [origin.value.s3_origin] : []
        content {
          origin_access_identity = s3_origin_config.value.origin_access_identity
        }
      }

      dynamic "custom_origin_config" {
        for_each = origin.value.custom_origin != null ? [origin.value.custom_origin] : []
        content {
          origin_protocol_policy = custom_origin_config.value.protocol_policy
          http_port              = custom_origin_config.value.http_port != 0 ? custom_origin_config.value.http_port : 80
          https_port             = custom_origin_config.value.https_port != 0 ? custom_origin_config.value.https_port : 443
          # Empty keeps the safe modern floor.
          origin_ssl_protocols     = length(custom_origin_config.value.ssl_protocols) > 0 ? custom_origin_config.value.ssl_protocols : ["TLSv1.2"]
          origin_keepalive_timeout = custom_origin_config.value.keepalive_timeout_seconds != 0 ? custom_origin_config.value.keepalive_timeout_seconds : null
          origin_read_timeout      = custom_origin_config.value.read_timeout_seconds != 0 ? custom_origin_config.value.read_timeout_seconds : null
        }
      }

      dynamic "vpc_origin_config" {
        for_each = origin.value.vpc_origin != null ? [origin.value.vpc_origin] : []
        content {
          vpc_origin_id            = vpc_origin_config.value.vpc_origin_id
          origin_keepalive_timeout = vpc_origin_config.value.keepalive_timeout_seconds != 0 ? vpc_origin_config.value.keepalive_timeout_seconds : null
          origin_read_timeout      = vpc_origin_config.value.read_timeout_seconds != 0 ? vpc_origin_config.value.read_timeout_seconds : null
        }
      }
    }
  }

  # --- Origin groups (primary/failover pairs) ------------------------------
  dynamic "origin_group" {
    for_each = { for g in var.spec.origin_groups : g.origin_group_id => g }
    content {
      origin_id = origin_group.key

      failover_criteria {
        status_codes = origin_group.value.failover_status_codes
      }

      # Exactly two members, primary first (CEL enforces the pair).
      dynamic "member" {
        for_each = origin_group.value.member_origin_ids
        content {
          origin_id = member.value
        }
      }
    }
  }

  # --- Default cache behavior ----------------------------------------------
  # The same field mapping as the ordered behaviors below; CloudFront requires
  # exactly one default. Modern configurations carry a cache policy; legacy
  # ones carry forwarded_values + per-behavior TTLs (CEL keeps the two
  # generations exclusive, matching the provider's own constraint).
  default_cache_behavior {
    target_origin_id       = local.default_behavior.target_origin_id
    viewer_protocol_policy = local.default_behavior.viewer_protocol_policy
    allowed_methods        = length(local.default_behavior.allowed_methods) > 0 ? local.default_behavior.allowed_methods : ["GET", "HEAD"]
    cached_methods         = length(local.default_behavior.cached_methods) > 0 ? local.default_behavior.cached_methods : ["GET", "HEAD"]
    compress               = local.default_behavior.compress
    smooth_streaming       = local.default_behavior.smooth_streaming

    cache_policy_id            = local.default_behavior.cache_policy_id != "" ? local.default_behavior.cache_policy_id : null
    origin_request_policy_id   = local.default_behavior.origin_request_policy_id != "" ? local.default_behavior.origin_request_policy_id : null
    response_headers_policy_id = local.default_behavior.response_headers_policy_id != "" ? local.default_behavior.response_headers_policy_id : null
    field_level_encryption_id  = local.default_behavior.field_level_encryption_id != "" ? local.default_behavior.field_level_encryption_id : null
    realtime_log_config_arn    = local.default_behavior.realtime_log_config_arn != "" ? local.default_behavior.realtime_log_config_arn : null

    trusted_key_groups = length(local.default_behavior.trusted_key_group_ids) > 0 ? local.default_behavior.trusted_key_group_ids : null
    trusted_signers    = length(local.default_behavior.trusted_signers) > 0 ? local.default_behavior.trusted_signers : null

    # TTLs only apply to the legacy generation -- with a cache policy the
    # policy owns them and the provider rejects the combination.
    min_ttl     = local.default_behavior.forwarded_values != null ? local.default_behavior.min_ttl_seconds : null
    default_ttl = try(local.default_behavior.forwarded_values, null) != null && local.default_behavior.default_ttl_seconds != 0 ? local.default_behavior.default_ttl_seconds : null
    max_ttl     = try(local.default_behavior.forwarded_values, null) != null && local.default_behavior.max_ttl_seconds != 0 ? local.default_behavior.max_ttl_seconds : null

    dynamic "forwarded_values" {
      for_each = local.default_behavior.forwarded_values != null ? [local.default_behavior.forwarded_values] : []
      content {
        query_string            = forwarded_values.value.query_string
        query_string_cache_keys = forwarded_values.value.query_string_cache_keys
        headers                 = forwarded_values.value.headers
        cookies {
          forward           = forwarded_values.value.cookies_forward
          whitelisted_names = forwarded_values.value.whitelisted_cookie_names
        }
      }
    }

    dynamic "function_association" {
      for_each = local.default_behavior.function_associations
      content {
        event_type   = function_association.value.event_type
        function_arn = function_association.value.function_arn
      }
    }

    dynamic "lambda_function_association" {
      for_each = local.default_behavior.lambda_function_associations
      content {
        event_type   = lambda_function_association.value.event_type
        lambda_arn   = lambda_function_association.value.lambda_arn
        include_body = lambda_function_association.value.include_body
      }
    }

    dynamic "grpc_config" {
      for_each = local.default_behavior.grpc_enabled ? [1] : []
      content {
        enabled = true
      }
    }
  }

  # --- Ordered cache behaviors (path-matched, first match wins) ------------
  dynamic "ordered_cache_behavior" {
    for_each = var.spec.ordered_cache_behaviors
    content {
      path_pattern           = ordered_cache_behavior.value.path_pattern
      target_origin_id       = ordered_cache_behavior.value.behavior.target_origin_id
      viewer_protocol_policy = ordered_cache_behavior.value.behavior.viewer_protocol_policy
      allowed_methods        = length(ordered_cache_behavior.value.behavior.allowed_methods) > 0 ? ordered_cache_behavior.value.behavior.allowed_methods : ["GET", "HEAD"]
      cached_methods         = length(ordered_cache_behavior.value.behavior.cached_methods) > 0 ? ordered_cache_behavior.value.behavior.cached_methods : ["GET", "HEAD"]
      compress               = ordered_cache_behavior.value.behavior.compress
      smooth_streaming       = ordered_cache_behavior.value.behavior.smooth_streaming

      cache_policy_id            = ordered_cache_behavior.value.behavior.cache_policy_id != "" ? ordered_cache_behavior.value.behavior.cache_policy_id : null
      origin_request_policy_id   = ordered_cache_behavior.value.behavior.origin_request_policy_id != "" ? ordered_cache_behavior.value.behavior.origin_request_policy_id : null
      response_headers_policy_id = ordered_cache_behavior.value.behavior.response_headers_policy_id != "" ? ordered_cache_behavior.value.behavior.response_headers_policy_id : null
      field_level_encryption_id  = ordered_cache_behavior.value.behavior.field_level_encryption_id != "" ? ordered_cache_behavior.value.behavior.field_level_encryption_id : null
      realtime_log_config_arn    = ordered_cache_behavior.value.behavior.realtime_log_config_arn != "" ? ordered_cache_behavior.value.behavior.realtime_log_config_arn : null

      trusted_key_groups = length(ordered_cache_behavior.value.behavior.trusted_key_group_ids) > 0 ? ordered_cache_behavior.value.behavior.trusted_key_group_ids : null
      trusted_signers    = length(ordered_cache_behavior.value.behavior.trusted_signers) > 0 ? ordered_cache_behavior.value.behavior.trusted_signers : null

      min_ttl     = ordered_cache_behavior.value.behavior.forwarded_values != null ? ordered_cache_behavior.value.behavior.min_ttl_seconds : null
      default_ttl = try(ordered_cache_behavior.value.behavior.forwarded_values, null) != null && ordered_cache_behavior.value.behavior.default_ttl_seconds != 0 ? ordered_cache_behavior.value.behavior.default_ttl_seconds : null
      max_ttl     = try(ordered_cache_behavior.value.behavior.forwarded_values, null) != null && ordered_cache_behavior.value.behavior.max_ttl_seconds != 0 ? ordered_cache_behavior.value.behavior.max_ttl_seconds : null

      dynamic "forwarded_values" {
        for_each = ordered_cache_behavior.value.behavior.forwarded_values != null ? [ordered_cache_behavior.value.behavior.forwarded_values] : []
        content {
          query_string            = forwarded_values.value.query_string
          query_string_cache_keys = forwarded_values.value.query_string_cache_keys
          headers                 = forwarded_values.value.headers
          cookies {
            forward           = forwarded_values.value.cookies_forward
            whitelisted_names = forwarded_values.value.whitelisted_cookie_names
          }
        }
      }

      dynamic "function_association" {
        for_each = ordered_cache_behavior.value.behavior.function_associations
        content {
          event_type   = function_association.value.event_type
          function_arn = function_association.value.function_arn
        }
      }

      dynamic "lambda_function_association" {
        for_each = ordered_cache_behavior.value.behavior.lambda_function_associations
        content {
          event_type   = lambda_function_association.value.event_type
          lambda_arn   = lambda_function_association.value.lambda_arn
          include_body = lambda_function_association.value.include_body
        }
      }

      dynamic "grpc_config" {
        for_each = ordered_cache_behavior.value.behavior.grpc_enabled ? [1] : []
        content {
          enabled = true
        }
      }
    }
  }

  # --- Viewer certificate ---------------------------------------------------
  # Precedence mirrors the provider: IAM certificate, then ACM, then the
  # default *.cloudfront.net certificate. Custom certificates default to
  # SNI-only ("vip" costs dedicated IPs) and the TLSv1.2_2021 floor.
  viewer_certificate {
    cloudfront_default_certificate = local.has_custom_viewer_cert ? null : true

    acm_certificate_arn = (
      local.has_custom_viewer_cert && try(local.viewer_cert.iam_certificate_id, "") == ""
      ? try(local.viewer_cert.acm_certificate_arn, null)
      : null
    )
    iam_certificate_id = (
      local.has_custom_viewer_cert && try(local.viewer_cert.iam_certificate_id, "") != ""
      ? try(local.viewer_cert.iam_certificate_id, null)
      : null
    )

    ssl_support_method = (
      local.has_custom_viewer_cert
      ? (try(local.viewer_cert.ssl_support_method, "") != "" ? try(local.viewer_cert.ssl_support_method, null) : "sni-only")
      : null
    )
    minimum_protocol_version = (
      local.has_custom_viewer_cert
      ? (try(local.viewer_cert.minimum_protocol_version, "") != "" ? try(local.viewer_cert.minimum_protocol_version, null) : "TLSv1.2_2021")
      : null
    )
  }

  # --- Custom error responses ----------------------------------------------
  dynamic "custom_error_response" {
    for_each = var.spec.custom_error_responses
    content {
      error_code            = custom_error_response.value.error_code
      response_code         = custom_error_response.value.response_code != 0 ? custom_error_response.value.response_code : null
      response_page_path    = custom_error_response.value.response_page_path != "" ? custom_error_response.value.response_page_path : null
      error_caching_min_ttl = custom_error_response.value.error_caching_min_ttl_seconds != 0 ? custom_error_response.value.error_caching_min_ttl_seconds : null
    }
  }

  # --- Geo restriction --------------------------------------------------------
  # The restrictions block is mandatory; absent spec geo_restriction means no
  # geographic restriction.
  restrictions {
    geo_restriction {
      restriction_type = var.spec.geo_restriction != null ? var.spec.geo_restriction.restriction_type : "none"
      locations        = var.spec.geo_restriction != null ? var.spec.geo_restriction.locations : null
    }
  }

  # --- Standard access logs ----------------------------------------------------
  dynamic "logging_config" {
    for_each = var.spec.logging != null ? [var.spec.logging] : []
    content {
      bucket          = logging_config.value.bucket
      prefix          = logging_config.value.prefix != "" ? logging_config.value.prefix : null
      include_cookies = logging_config.value.include_cookies
    }
  }

  tags = local.aws_tags
}

# CloudWatch additional metrics (cache hit rate, origin latency, per-status
# error rates) -- a distribution-keyed setting materialized as its own
# provider resource, folded into the spec as one honest toggle.
resource "aws_cloudfront_monitoring_subscription" "this" {
  count = var.spec.enable_additional_metrics ? 1 : 0

  distribution_id = aws_cloudfront_distribution.this.id

  monitoring_subscription {
    realtime_metrics_subscription_config {
      realtime_metrics_subscription_status = "Enabled"
    }
  }
}
