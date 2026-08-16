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
  description = "AwsOpenSearchDomain specification"
  type = object({
    region = string
    engine_version = string
    cluster_config = object({
      instance_type = string
      instance_count = optional(number)
      dedicated_master_enabled = optional(bool, false)
      dedicated_master_type = optional(string, "")
      dedicated_master_count = optional(number, 0)
      node_options = optional(list(object({
        node_type = string
        enabled = optional(bool, false)
        instance_type = optional(string, "")
        count = optional(number, 0)
      })), [])
      zone_awareness_enabled = optional(bool, false)
      availability_zone_count = optional(number, 0)
      warm_enabled = optional(bool, false)
      warm_type = optional(string, "")
      warm_count = optional(number, 0)
      cold_storage_enabled = optional(bool, false)
      multi_az_with_standby_enabled = optional(bool, false)
    })
    ebs_options = object({
      ebs_enabled = optional(bool, false)
      volume_type = optional(string, "")
      volume_size = optional(number, 0)
      iops = optional(number, 0)
      throughput = optional(number, 0)
    })
    encrypt_at_rest_enabled = optional(bool)
    kms_key_id = optional(string, "")
    node_to_node_encryption_enabled = optional(bool)
    vpc_options = optional(object({
      subnet_ids = optional(list(string), [])
      security_group_ids = optional(list(string), [])
    }))
    domain_endpoint_options = optional(object({
      enforce_https = optional(bool)
      tls_security_policy = optional(string, "")
      custom_endpoint_enabled = optional(bool, false)
      custom_endpoint = optional(string, "")
      custom_endpoint_certificate_arn = optional(string, "")
    }))
    advanced_security_options = optional(object({
      enabled = optional(bool, false)
      internal_user_database_enabled = optional(bool, false)
      anonymous_auth_enabled = optional(bool, false)
      master_user_arn = optional(string, "")
      master_user_name = optional(string, "")
      master_user_password = optional(string, "")
      jwt_options = optional(object({
        enabled = optional(bool, false)
        jwks_url = optional(string, "")
        public_key = optional(string, "")
        roles_key = optional(string, "")
        subject_key = optional(string, "")
      }))
    }))
    cognito_options = optional(object({
      enabled = optional(bool, false)
      user_pool_id = optional(string, "")
      identity_pool_id = optional(string, "")
      role_arn = optional(string, "")
    }))
    log_publishing_options = optional(list(object({
      log_type = string
      cloudwatch_log_group_arn = string
      enabled = optional(bool)
    })), [])
    access_policies = optional(any)
    auto_tune_options = optional(object({
      desired_state = string
      maintenance_schedules = optional(list(object({
        start_at = string
        duration_hours = number
        cron_expression_for_recurrence = string
      })), [])
      rollback_on_disable = optional(string, "")
      use_off_peak_window = optional(bool, false)
    }))
    automated_snapshot_start_hour = optional(number)
    off_peak_window_options = optional(object({
      enabled = optional(bool)
      window_start_hour = optional(number)
      window_start_minute = optional(number)
    }))
    auto_software_update_enabled = optional(bool, false)
    deployment_strategy = optional(string, "")
    ip_address_type = optional(string, "")
    advanced_options = optional(map(string), {})
    aiml_options = optional(object({
      natural_language_query_generation_desired_state = optional(string, "")
      s3_vectors_engine_enabled = optional(bool, false)
      serverless_vector_acceleration_enabled = optional(bool, false)
    }))
    identity_center_options = optional(object({
      enabled_api_access = optional(bool, false)
      identity_center_instance_arn = optional(string, "")
      roles_key = optional(string, "")
      subject_key = optional(string, "")
    }))
    saml_options = optional(object({
      idp_entity_id = string
      idp_metadata_content = string
      master_backend_role = optional(string, "")
      master_user_name = optional(string, "")
      roles_key = optional(string, "")
      subject_key = optional(string, "")
      session_timeout_minutes = optional(number, 0)
    }))
    authorized_vpc_endpoint_access_accounts = optional(list(string), [])
  })
}
