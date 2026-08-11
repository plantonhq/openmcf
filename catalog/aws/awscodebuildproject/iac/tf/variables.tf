variable "metadata" {
  description = "Cloud resource metadata"
  type = object({
    name = string
    id = optional(string, "")
    org = optional(string, "")
    env = optional(string, "")
    labels = optional(map(string), {})
    annotations = optional(map(string), {})
    tags = optional(list(string), [])
  })
}

variable "spec" {
  description = "AwsCodeBuildProject specification"
  type = object({
    region = string
    source = object({
      type = string
      location = optional(string, "")
      buildspec = optional(string, "")
      git_clone_depth = optional(number, 0)
      git_submodules_config = optional(object({
        fetch_submodules = optional(bool, false)
      }))
      insecure_ssl = optional(bool, false)
      report_build_status = optional(bool, false)
      build_status_config = optional(object({
        context = optional(string, "")
        target_url = optional(string, "")
      }))
      auth = optional(object({
        type = string
        resource = string
      }))
      source_identifier = optional(string, "")
    })
    secondary_sources = optional(list(object({
      type = string
      location = optional(string, "")
      buildspec = optional(string, "")
      git_clone_depth = optional(number, 0)
      git_submodules_config = optional(object({
        fetch_submodules = optional(bool, false)
      }))
      insecure_ssl = optional(bool, false)
      report_build_status = optional(bool, false)
      build_status_config = optional(object({
        context = optional(string, "")
        target_url = optional(string, "")
      }))
      auth = optional(object({
        type = string
        resource = string
      }))
      source_identifier = optional(string, "")
    })), [])
    secondary_source_versions = optional(list(object({
      source_identifier = string
      source_version = string
    })), [])
    environment = object({
      type = string
      compute_type = string
      image = string
      certificate = optional(string, "")
      privileged_mode = optional(bool, false)
      image_pull_credentials_type = optional(string)
      environment_variables = optional(list(object({
        name = string
        value = string
        type = optional(string)
      })), [])
      registry_credential = optional(object({
        credential = string
        credential_provider = string
      }))
      docker_server = optional(object({
        compute_type = string
        security_group_ids = optional(list(string), [])
      }))
      fleet_arn = optional(string, "")
      host_kernel = optional(string, "")
    })
    artifacts = object({
      type = string
      location = optional(string, "")
      name = optional(string, "")
      path = optional(string, "")
      packaging = optional(string, "")
      namespace_type = optional(string, "")
      encryption_disabled = optional(bool, false)
      override_artifact_name = optional(bool, false)
      bucket_owner_access = optional(string, "")
      artifact_identifier = optional(string, "")
    })
    secondary_artifacts = optional(list(object({
      type = string
      location = optional(string, "")
      name = optional(string, "")
      path = optional(string, "")
      packaging = optional(string, "")
      namespace_type = optional(string, "")
      encryption_disabled = optional(bool, false)
      override_artifact_name = optional(bool, false)
      bucket_owner_access = optional(string, "")
      artifact_identifier = optional(string, "")
    })), [])
    service_role = string
    description = optional(string, "")
    encryption_key = optional(string, "")
    build_timeout = optional(number)
    queued_timeout = optional(number)
    concurrent_build_limit = optional(number, 0)
    auto_retry_limit = optional(number, 0)
    badge_enabled = optional(bool, false)
    source_version = optional(string, "")
    cache = optional(object({
      type = optional(string)
      location = optional(string, "")
      modes = optional(list(string), [])
      cache_namespace = optional(string, "")
    }))
    logs_config = optional(object({
      cloudwatch_logs = optional(object({
        status = optional(string)
        group_name = optional(string, "")
        stream_name = optional(string, "")
      }))
      s3_logs = optional(object({
        status = optional(string)
        location = optional(string, "")
        encryption_disabled = optional(bool, false)
        bucket_owner_access = optional(string, "")
      }))
    }))
    vpc_config = optional(object({
      vpc_id = string
      subnet_ids = list(string)
      security_group_ids = list(string)
    }))
    file_system_locations = optional(list(object({
      identifier = string
      location = string
      mount_point = string
      mount_options = optional(string, "")
      type = optional(string)
    })), [])
    build_batch_config = optional(object({
      service_role = string
      combine_artifacts = optional(bool, false)
      timeout_in_mins = optional(number, 0)
      restrictions = optional(object({
        compute_types_allowed = optional(list(string), [])
        maximum_builds_allowed = optional(number, 0)
      }))
    }))
    project_visibility = optional(string)
    resource_access_role = optional(string, "")
    resource_policy = optional(any)
    webhook = optional(object({
      build_type = optional(string, "")
      manual_creation = optional(bool, false)
      filter_groups = optional(list(object({
        filters = list(object({
          type = string
          pattern = string
          exclude_matched_pattern = optional(bool, false)
        }))
      })), [])
      scope_configuration = optional(object({
        name = string
        scope = string
        domain = optional(string, "")
      }))
      pull_request_build_policy = optional(object({
        requires_comment_approval = string
        approver_roles = optional(list(string), [])
      }))
    }))
  })
}
