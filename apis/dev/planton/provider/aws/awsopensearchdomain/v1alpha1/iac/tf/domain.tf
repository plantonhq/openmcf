# ---------------------------------------------------------------------------
# OpenSearch Service domain
#
# Create-time (ForceNew) surfaces: the domain name, the at-rest KMS key,
# adding/removing vpc_options, and disabling either encryption toggle or FGAC
# once enabled. Everything else -- topology, storage, endpoint policy,
# Auto-Tune, log publishing -- updates in place, though topology changes
# usually run as a blue/green deployment behind the endpoint
# (deployment_strategy tunes how that migration provisions capacity).
# ---------------------------------------------------------------------------

resource "aws_opensearch_domain" "this" {
  domain_name    = local.domain_name
  engine_version = var.spec.engine_version

  # ---------------------------------------------------------------------------
  # Cluster configuration
  # ---------------------------------------------------------------------------

  cluster_config {
    instance_type  = var.spec.cluster_config.instance_type
    instance_count = coalesce(var.spec.cluster_config.instance_count, 1)

    # Dedicated master nodes
    dedicated_master_enabled = var.spec.cluster_config.dedicated_master_enabled
    dedicated_master_type    = var.spec.cluster_config.dedicated_master_enabled && var.spec.cluster_config.dedicated_master_type != "" ? var.spec.cluster_config.dedicated_master_type : null
    dedicated_master_count   = var.spec.cluster_config.dedicated_master_enabled && var.spec.cluster_config.dedicated_master_count > 0 ? var.spec.cluster_config.dedicated_master_count : null

    # Coordinator node pools: request routing, query fan-out, and response
    # aggregation offloaded from the data nodes.
    dynamic "node_options" {
      for_each = var.spec.cluster_config.node_options
      content {
        node_type = node_options.value.node_type
        node_config {
          enabled = node_options.value.enabled
          type    = node_options.value.instance_type != "" ? node_options.value.instance_type : null
          count   = node_options.value.count > 0 ? node_options.value.count : null
        }
      }
    }

    # Zone awareness
    zone_awareness_enabled = var.spec.cluster_config.zone_awareness_enabled

    dynamic "zone_awareness_config" {
      for_each = var.spec.cluster_config.zone_awareness_enabled && var.spec.cluster_config.availability_zone_count > 0 ? [1] : []
      content {
        availability_zone_count = var.spec.cluster_config.availability_zone_count
      }
    }

    # UltraWarm storage
    warm_enabled = var.spec.cluster_config.warm_enabled
    warm_type    = var.spec.cluster_config.warm_enabled ? var.spec.cluster_config.warm_type : null
    warm_count   = var.spec.cluster_config.warm_enabled ? var.spec.cluster_config.warm_count : null

    # Cold storage (S3-backed; requires UltraWarm)
    dynamic "cold_storage_options" {
      for_each = var.spec.cluster_config.cold_storage_enabled ? [1] : []
      content {
        enabled = true
      }
    }

    # Multi-AZ with standby (99.99% SLA posture)
    multi_az_with_standby_enabled = var.spec.cluster_config.multi_az_with_standby_enabled
  }

  # ---------------------------------------------------------------------------
  # EBS storage
  # ---------------------------------------------------------------------------

  ebs_options {
    ebs_enabled = var.spec.ebs_options.ebs_enabled
    volume_type = var.spec.ebs_options.ebs_enabled && var.spec.ebs_options.volume_type != "" ? var.spec.ebs_options.volume_type : null
    volume_size = var.spec.ebs_options.ebs_enabled && var.spec.ebs_options.volume_size > 0 ? var.spec.ebs_options.volume_size : null
    iops        = var.spec.ebs_options.ebs_enabled && var.spec.ebs_options.iops > 0 ? var.spec.ebs_options.iops : null
    throughput  = var.spec.ebs_options.ebs_enabled && var.spec.ebs_options.throughput > 0 ? var.spec.ebs_options.throughput : null
  }

  # ---------------------------------------------------------------------------
  # Encryption -- both toggles are one-way (cannot be disabled once on), and
  # the KMS key is fixed at creation.
  # ---------------------------------------------------------------------------

  encrypt_at_rest {
    enabled    = var.spec.encrypt_at_rest_enabled
    kms_key_id = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null
  }

  node_to_node_encryption {
    enabled = var.spec.node_to_node_encryption_enabled
  }

  # ---------------------------------------------------------------------------
  # VPC (conditional; ForceNew)
  # ---------------------------------------------------------------------------

  dynamic "vpc_options" {
    for_each = local.has_vpc_options ? [1] : []
    content {
      subnet_ids         = var.spec.vpc_options.subnet_ids
      security_group_ids = length(var.spec.vpc_options.security_group_ids) > 0 ? var.spec.vpc_options.security_group_ids : null
    }
  }

  # ---------------------------------------------------------------------------
  # Domain endpoint options
  # ---------------------------------------------------------------------------

  dynamic "domain_endpoint_options" {
    for_each = var.spec.domain_endpoint_options != null ? [var.spec.domain_endpoint_options] : []
    content {
      enforce_https                   = domain_endpoint_options.value.enforce_https
      tls_security_policy             = domain_endpoint_options.value.tls_security_policy != "" ? domain_endpoint_options.value.tls_security_policy : null
      custom_endpoint_enabled         = domain_endpoint_options.value.custom_endpoint_enabled
      custom_endpoint                 = domain_endpoint_options.value.custom_endpoint_enabled ? domain_endpoint_options.value.custom_endpoint : null
      custom_endpoint_certificate_arn = domain_endpoint_options.value.custom_endpoint_enabled && domain_endpoint_options.value.custom_endpoint_certificate_arn != "" ? domain_endpoint_options.value.custom_endpoint_certificate_arn : null
    }
  }

  # ---------------------------------------------------------------------------
  # Fine-grained access control (FGAC). Emitted only when enabled -- FGAC is
  # one-way in AWS. JWT bearer auth and anonymous auth ride inside it.
  # ---------------------------------------------------------------------------

  dynamic "advanced_security_options" {
    for_each = local.has_advanced_security ? [var.spec.advanced_security_options] : []
    content {
      enabled                        = true
      internal_user_database_enabled = advanced_security_options.value.internal_user_database_enabled
      anonymous_auth_enabled         = advanced_security_options.value.anonymous_auth_enabled ? true : null

      dynamic "master_user_options" {
        for_each = advanced_security_options.value.master_user_arn != "" || advanced_security_options.value.master_user_name != "" ? [1] : []
        content {
          master_user_arn      = advanced_security_options.value.master_user_arn != "" ? advanced_security_options.value.master_user_arn : null
          master_user_name     = advanced_security_options.value.master_user_name != "" ? advanced_security_options.value.master_user_name : null
          master_user_password = advanced_security_options.value.master_user_password != "" ? advanced_security_options.value.master_user_password : null
        }
      }

      dynamic "jwt_options" {
        for_each = advanced_security_options.value.jwt_options != null ? [advanced_security_options.value.jwt_options] : []
        content {
          enabled     = jwt_options.value.enabled
          jwks_url    = jwt_options.value.jwks_url != "" ? jwt_options.value.jwks_url : null
          public_key  = jwt_options.value.public_key != "" ? jwt_options.value.public_key : null
          roles_key   = jwt_options.value.roles_key != "" ? jwt_options.value.roles_key : null
          subject_key = jwt_options.value.subject_key != "" ? jwt_options.value.subject_key : null
        }
      }
    }
  }

  # ---------------------------------------------------------------------------
  # Cognito authentication for OpenSearch Dashboards
  # ---------------------------------------------------------------------------

  dynamic "cognito_options" {
    for_each = local.has_cognito ? [var.spec.cognito_options] : []
    content {
      enabled          = true
      user_pool_id     = cognito_options.value.user_pool_id
      identity_pool_id = cognito_options.value.identity_pool_id
      role_arn         = cognito_options.value.role_arn
    }
  }

  # ---------------------------------------------------------------------------
  # Log publishing (conditional, repeated)
  # ---------------------------------------------------------------------------

  dynamic "log_publishing_options" {
    for_each = var.spec.log_publishing_options
    content {
      log_type                 = log_publishing_options.value.log_type
      cloudwatch_log_group_arn = log_publishing_options.value.cloudwatch_log_group_arn
      enabled                  = log_publishing_options.value.enabled
    }
  }

  # ---------------------------------------------------------------------------
  # Access policies
  # ---------------------------------------------------------------------------

  access_policies = local.access_policies_json

  # ---------------------------------------------------------------------------
  # Auto-Tune. Non-disruptive JVM tuning applies immediately; blue/green
  # optimizations wait for a maintenance schedule or the off-peak window.
  # Not supported on t2/t3 (burstable) instance types.
  # ---------------------------------------------------------------------------

  dynamic "auto_tune_options" {
    for_each = var.spec.auto_tune_options != null ? [var.spec.auto_tune_options] : []
    content {
      desired_state       = auto_tune_options.value.desired_state
      rollback_on_disable = auto_tune_options.value.rollback_on_disable != "" ? auto_tune_options.value.rollback_on_disable : null
      use_off_peak_window = auto_tune_options.value.use_off_peak_window

      dynamic "maintenance_schedule" {
        for_each = auto_tune_options.value.maintenance_schedules
        content {
          start_at                       = maintenance_schedule.value.start_at
          cron_expression_for_recurrence = maintenance_schedule.value.cron_expression_for_recurrence
          duration {
            # HOURS is the only duration unit the AWS API supports.
            unit  = "HOURS"
            value = maintenance_schedule.value.duration_hours
          }
        }
      }
    }
  }

  # ---------------------------------------------------------------------------
  # Snapshots, off-peak window, software updates
  # ---------------------------------------------------------------------------

  dynamic "snapshot_options" {
    for_each = var.spec.automated_snapshot_start_hour != null ? [1] : []
    content {
      automated_snapshot_start_hour = var.spec.automated_snapshot_start_hour
    }
  }

  dynamic "off_peak_window_options" {
    for_each = var.spec.off_peak_window_options != null ? [var.spec.off_peak_window_options] : []
    content {
      enabled = off_peak_window_options.value.enabled

      dynamic "off_peak_window" {
        for_each = off_peak_window_options.value.window_start_hour != null ? [1] : []
        content {
          window_start_time {
            hours   = off_peak_window_options.value.window_start_hour
            minutes = coalesce(off_peak_window_options.value.window_start_minute, 0)
          }
        }
      }
    }
  }

  software_update_options {
    auto_software_update_enabled = var.spec.auto_software_update_enabled
  }

  # ---------------------------------------------------------------------------
  # Blue/green deployment strategy for config changes that require one
  # ---------------------------------------------------------------------------

  dynamic "deployment_strategy_options" {
    for_each = var.spec.deployment_strategy != "" ? [1] : []
    content {
      deployment_strategy = var.spec.deployment_strategy
    }
  }

  # ---------------------------------------------------------------------------
  # IP and advanced options
  # ---------------------------------------------------------------------------

  # One-way: dualstack -> ipv4 replaces the domain.
  ip_address_type = var.spec.ip_address_type != "" ? var.spec.ip_address_type : null

  advanced_options = length(var.spec.advanced_options) > 0 ? var.spec.advanced_options : null

  # ---------------------------------------------------------------------------
  # AI/ML capabilities
  # ---------------------------------------------------------------------------

  dynamic "aiml_options" {
    for_each = var.spec.aiml_options != null ? [var.spec.aiml_options] : []
    content {
      dynamic "natural_language_query_generation_options" {
        for_each = aiml_options.value.natural_language_query_generation_desired_state != "" ? [1] : []
        content {
          desired_state = aiml_options.value.natural_language_query_generation_desired_state
        }
      }

      dynamic "s3_vectors_engine" {
        for_each = aiml_options.value.s3_vectors_engine_enabled ? [1] : []
        content {
          enabled = true
        }
      }

      dynamic "serverless_vector_acceleration" {
        for_each = aiml_options.value.serverless_vector_acceleration_enabled ? [1] : []
        content {
          enabled = true
        }
      }
    }
  }

  # ---------------------------------------------------------------------------
  # IAM Identity Center
  # ---------------------------------------------------------------------------

  dynamic "identity_center_options" {
    for_each = var.spec.identity_center_options != null ? [var.spec.identity_center_options] : []
    content {
      enabled_api_access           = identity_center_options.value.enabled_api_access
      identity_center_instance_arn = identity_center_options.value.identity_center_instance_arn != "" ? identity_center_options.value.identity_center_instance_arn : null
      roles_key                    = identity_center_options.value.roles_key != "" ? identity_center_options.value.roles_key : null
      subject_key                  = identity_center_options.value.subject_key != "" ? identity_center_options.value.subject_key : null
    }
  }

  # ---------------------------------------------------------------------------
  # Tags
  # ---------------------------------------------------------------------------

  tags = local.aws_tags
}
