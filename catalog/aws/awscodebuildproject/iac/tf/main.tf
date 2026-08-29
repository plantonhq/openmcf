# CodeBuild project -- the build configuration unit. The project itself is a
# metadata-only control-plane resource: creating it provisions nothing until
# a build actually starts, so create/update/delete are near-instant. The only
# operational wait AWS imposes is IAM eventual consistency on a freshly
# created service role, which the provider absorbs with a bounded retry.

resource "aws_codebuild_project" "this" {
  name         = local.project_name
  description  = var.spec.description != "" ? var.spec.description : null
  service_role = var.spec.service_role

  # Lambda environments ignore timeouts entirely (AWS caps them itself);
  # sending the spec defaults would create a permanent plan diff there.
  build_timeout  = local.is_lambda_env ? null : var.spec.build_timeout
  queued_timeout = local.is_lambda_env ? null : var.spec.queued_timeout

  concurrent_build_limit = var.spec.concurrent_build_limit > 0 ? var.spec.concurrent_build_limit : null

  # Additional automatic retries after a failed build (not total attempts).
  auto_retry_limit = var.spec.auto_retry_limit > 0 ? var.spec.auto_retry_limit : null

  badge_enabled  = var.spec.badge_enabled
  encryption_key = var.spec.encryption_key != "" ? var.spec.encryption_key : null
  source_version = var.spec.source_version != "" ? var.spec.source_version : null

  # Public visibility is managed by a separate AWS API call under the hood
  # (UpdateProjectVisibility); the provider sequences it after project
  # create/update. The resource access role is what CodeBuild uses to read
  # the logs/artifacts it re-exposes publicly.
  project_visibility   = var.spec.project_visibility
  resource_access_role = var.spec.resource_access_role != "" ? var.spec.resource_access_role : null

  tags = local.tags

  # --- Primary source -------------------------------------------------------

  source {
    type            = var.spec.source.type
    location        = var.spec.source.location != "" ? var.spec.source.location : null
    buildspec       = var.spec.source.buildspec != "" ? var.spec.source.buildspec : null
    git_clone_depth = var.spec.source.git_clone_depth > 0 ? var.spec.source.git_clone_depth : null
    insecure_ssl    = var.spec.source.insecure_ssl

    report_build_status = var.spec.source.report_build_status

    dynamic "git_submodules_config" {
      for_each = var.spec.source.git_submodules_config != null ? [var.spec.source.git_submodules_config] : []
      content {
        fetch_submodules = git_submodules_config.value.fetch_submodules
      }
    }

    dynamic "build_status_config" {
      for_each = var.spec.source.build_status_config != null ? [var.spec.source.build_status_config] : []
      content {
        context    = build_status_config.value.context != "" ? build_status_config.value.context : null
        target_url = build_status_config.value.target_url != "" ? build_status_config.value.target_url : null
      }
    }

    # Pins the authorization for this source (CodeConnections connection or
    # Secrets Manager token), overriding the account-level source credential.
    dynamic "auth" {
      for_each = var.spec.source.auth != null ? [var.spec.source.auth] : []
      content {
        type     = auth.value.type
        resource = auth.value.resource
      }
    }
  }

  # --- Secondary sources (identifier-addressed extra checkouts) ------------

  dynamic "secondary_sources" {
    for_each = var.spec.secondary_sources
    content {
      source_identifier = secondary_sources.value.source_identifier
      type              = secondary_sources.value.type
      location          = secondary_sources.value.location != "" ? secondary_sources.value.location : null
      buildspec         = secondary_sources.value.buildspec != "" ? secondary_sources.value.buildspec : null
      git_clone_depth   = secondary_sources.value.git_clone_depth > 0 ? secondary_sources.value.git_clone_depth : null
      insecure_ssl      = secondary_sources.value.insecure_ssl

      report_build_status = secondary_sources.value.report_build_status

      dynamic "git_submodules_config" {
        for_each = secondary_sources.value.git_submodules_config != null ? [secondary_sources.value.git_submodules_config] : []
        content {
          fetch_submodules = git_submodules_config.value.fetch_submodules
        }
      }

      dynamic "build_status_config" {
        for_each = secondary_sources.value.build_status_config != null ? [secondary_sources.value.build_status_config] : []
        content {
          context    = build_status_config.value.context != "" ? build_status_config.value.context : null
          target_url = build_status_config.value.target_url != "" ? build_status_config.value.target_url : null
        }
      }

      dynamic "auth" {
        for_each = secondary_sources.value.auth != null ? [secondary_sources.value.auth] : []
        content {
          type     = auth.value.type
          resource = auth.value.resource
        }
      }
    }
  }

  # Per-secondary-source version pins (branch/tag/commit per identifier).
  dynamic "secondary_source_version" {
    for_each = var.spec.secondary_source_versions
    content {
      source_identifier = secondary_source_version.value.source_identifier
      source_version    = secondary_source_version.value.source_version
    }
  }

  # --- Environment ----------------------------------------------------------

  environment {
    type                        = var.spec.environment.type
    compute_type                = var.spec.environment.compute_type
    image                       = var.spec.environment.image
    certificate                 = var.spec.environment.certificate != "" ? var.spec.environment.certificate : null
    privileged_mode             = var.spec.environment.privileged_mode
    image_pull_credentials_type = var.spec.environment.image_pull_credentials_type
    # Kernel selection is Linux container/EC2 surface (the spec's CEL scopes
    # it); omitted, AWS chooses (the argument is Optional+Computed).
    host_kernel = var.spec.environment.host_kernel != "" ? var.spec.environment.host_kernel : null

    dynamic "environment_variable" {
      for_each = var.spec.environment.environment_variables
      content {
        name  = environment_variable.value.name
        value = environment_variable.value.value
        type  = environment_variable.value.type
      }
    }

    dynamic "registry_credential" {
      for_each = var.spec.environment.registry_credential != null ? [var.spec.environment.registry_credential] : []
      content {
        credential          = registry_credential.value.credential
        credential_provider = registry_credential.value.credential_provider
      }
    }

    # Persistent, dedicated Docker server: layer state survives across
    # builds, unlike the per-build daemon privileged_mode provides.
    dynamic "docker_server" {
      for_each = var.spec.environment.docker_server != null ? [var.spec.environment.docker_server] : []
      content {
        compute_type       = docker_server.value.compute_type
        security_group_ids = length(docker_server.value.security_group_ids) > 0 ? docker_server.value.security_group_ids : null
      }
    }

    # Reserved-capacity fleet membership -- pre-provisioned, always-warm
    # build machines. The fleet is a shared account-level resource; the
    # project only references its ARN.
    dynamic "fleet" {
      for_each = var.spec.environment.fleet_arn != "" ? [var.spec.environment.fleet_arn] : []
      content {
        fleet_arn = fleet.value
      }
    }
  }

  # --- Primary artifacts ----------------------------------------------------

  artifacts {
    type                   = var.spec.artifacts.type
    artifact_identifier    = var.spec.artifacts.artifact_identifier != "" ? var.spec.artifacts.artifact_identifier : null
    location               = var.spec.artifacts.location != "" ? var.spec.artifacts.location : null
    name                   = var.spec.artifacts.name != "" ? var.spec.artifacts.name : null
    path                   = var.spec.artifacts.path != "" ? var.spec.artifacts.path : null
    packaging              = var.spec.artifacts.packaging != "" ? var.spec.artifacts.packaging : null
    namespace_type         = var.spec.artifacts.namespace_type != "" ? var.spec.artifacts.namespace_type : null
    encryption_disabled    = var.spec.artifacts.encryption_disabled
    override_artifact_name = var.spec.artifacts.override_artifact_name
    bucket_owner_access    = var.spec.artifacts.bucket_owner_access != "" ? var.spec.artifacts.bucket_owner_access : null
  }

  # --- Secondary artifacts (identifier-addressed extra outputs) ------------

  dynamic "secondary_artifacts" {
    for_each = var.spec.secondary_artifacts
    content {
      artifact_identifier    = secondary_artifacts.value.artifact_identifier
      type                   = secondary_artifacts.value.type
      location               = secondary_artifacts.value.location != "" ? secondary_artifacts.value.location : null
      name                   = secondary_artifacts.value.name != "" ? secondary_artifacts.value.name : null
      path                   = secondary_artifacts.value.path != "" ? secondary_artifacts.value.path : null
      packaging              = secondary_artifacts.value.packaging != "" ? secondary_artifacts.value.packaging : null
      namespace_type         = secondary_artifacts.value.namespace_type != "" ? secondary_artifacts.value.namespace_type : null
      encryption_disabled    = secondary_artifacts.value.encryption_disabled
      override_artifact_name = secondary_artifacts.value.override_artifact_name
      bucket_owner_access    = secondary_artifacts.value.bucket_owner_access != "" ? secondary_artifacts.value.bucket_owner_access : null
    }
  }

  # --- Cache ------------------------------------------------------------------

  dynamic "cache" {
    for_each = local.has_cache ? [var.spec.cache] : []
    content {
      type            = cache.value.type
      location        = cache.value.location != "" ? cache.value.location : null
      modes           = cache.value.type == "LOCAL" ? cache.value.modes : null
      cache_namespace = cache.value.cache_namespace != "" ? cache.value.cache_namespace : null
    }
  }

  # --- Logs -------------------------------------------------------------------

  dynamic "logs_config" {
    for_each = local.has_logs_config ? [var.spec.logs_config] : []
    content {
      dynamic "cloudwatch_logs" {
        for_each = logs_config.value.cloudwatch_logs != null ? [logs_config.value.cloudwatch_logs] : []
        content {
          status      = cloudwatch_logs.value.status
          group_name  = cloudwatch_logs.value.group_name != "" ? cloudwatch_logs.value.group_name : null
          stream_name = cloudwatch_logs.value.stream_name != "" ? cloudwatch_logs.value.stream_name : null
        }
      }
      dynamic "s3_logs" {
        for_each = logs_config.value.s3_logs != null ? [logs_config.value.s3_logs] : []
        content {
          status = s3_logs.value.status
          # AWS stores S3 build logs only under a prefix -- the provider
          # takes one "bucket/prefix" string, composed here from the
          # spec's two halves.
          location            = s3_logs.value.bucket != "" ? "${s3_logs.value.bucket}/${s3_logs.value.prefix}" : null
          encryption_disabled = s3_logs.value.encryption_disabled
          bucket_owner_access = s3_logs.value.bucket_owner_access != "" ? s3_logs.value.bucket_owner_access : null
        }
      }
    }
  }

  # --- VPC placement ----------------------------------------------------------

  dynamic "vpc_config" {
    for_each = local.has_vpc_config ? [var.spec.vpc_config] : []
    content {
      vpc_id             = var.spec.vpc_config.vpc_id
      subnets            = var.spec.vpc_config.subnet_ids
      security_group_ids = var.spec.vpc_config.security_group_ids
    }
  }

  # --- EFS mounts (shared caches that outlive individual builds) -------------

  dynamic "file_system_locations" {
    for_each = var.spec.file_system_locations
    content {
      type          = file_system_locations.value.type
      identifier    = file_system_locations.value.identifier
      location      = file_system_locations.value.location
      mount_point   = file_system_locations.value.mount_point
      mount_options = file_system_locations.value.mount_options != "" ? file_system_locations.value.mount_options : null
    }
  }

  # --- Batch builds -----------------------------------------------------------

  dynamic "build_batch_config" {
    for_each = local.has_batch_config ? [var.spec.build_batch_config] : []
    content {
      service_role      = build_batch_config.value.service_role
      combine_artifacts = build_batch_config.value.combine_artifacts
      timeout_in_mins   = build_batch_config.value.timeout_in_mins > 0 ? build_batch_config.value.timeout_in_mins : null

      dynamic "restrictions" {
        for_each = build_batch_config.value.restrictions != null ? [build_batch_config.value.restrictions] : []
        content {
          compute_types_allowed  = length(restrictions.value.compute_types_allowed) > 0 ? restrictions.value.compute_types_allowed : null
          maximum_builds_allowed = restrictions.value.maximum_builds_allowed > 0 ? restrictions.value.maximum_builds_allowed : null
        }
      }
    }
  }
}

# --- Webhook (folded 1:1 satellite) ------------------------------------------
#
# The webhook registers the project with its source provider so repository
# events trigger builds. With manual_creation, AWS mints the payload URL and
# HMAC secret WITHOUT registering anything -- the operator wires the
# repository webhook by hand from this module's outputs.

resource "aws_codebuild_webhook" "this" {
  count        = local.has_webhook ? 1 : 0
  project_name = aws_codebuild_project.this.name

  build_type      = var.spec.webhook.build_type != "" ? var.spec.webhook.build_type : null
  manual_creation = var.spec.webhook.manual_creation ? true : null

  dynamic "filter_group" {
    for_each = var.spec.webhook.filter_groups
    content {
      dynamic "filter" {
        for_each = filter_group.value.filters
        content {
          type                    = filter.value.type
          pattern                 = filter.value.pattern
          exclude_matched_pattern = filter.value.exclude_matched_pattern
        }
      }
    }
  }

  # Organization/group-scoped webhooks fire for every repository in scope
  # (runner projects, org-wide CI) instead of a single repository.
  dynamic "scope_configuration" {
    for_each = var.spec.webhook.scope_configuration != null ? [var.spec.webhook.scope_configuration] : []
    content {
      name   = scope_configuration.value.name
      scope  = scope_configuration.value.scope
      domain = scope_configuration.value.domain != "" ? scope_configuration.value.domain : null
    }
  }

  # Comment-approval gate for PR-triggered builds -- protects CI secrets
  # from untrusted fork code.
  dynamic "pull_request_build_policy" {
    for_each = var.spec.webhook.pull_request_build_policy != null ? [var.spec.webhook.pull_request_build_policy] : []
    content {
      requires_comment_approval = pull_request_build_policy.value.requires_comment_approval
      approver_roles            = length(pull_request_build_policy.value.approver_roles) > 0 ? pull_request_build_policy.value.approver_roles : null
    }
  }
}

# --- Resource policy (folded 1:1 satellite) -----------------------------------
#
# A resource-based IAM policy on the project itself -- the cross-account
# access mechanism (e.g., a central CI account starting builds here). One
# document per project, keyed by the project ARN, so it folds rather than
# being its own kind.

resource "aws_codebuild_resource_policy" "this" {
  count = local.has_resource_policy ? 1 : 0

  resource_arn = aws_codebuild_project.this.arn
  policy       = local.resource_policy_json
}
