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
  description = "AwsSagemakerDomain specification"
  type = object({
    region = string
    auth_mode = string
    vpc_id = string
    subnet_ids = list(string)
    kms_key_id = optional(string, "")
    app_network_access_type = optional(string)
    app_security_group_management = optional(string)
    tag_propagation = optional(string)
    home_efs_retention_policy = optional(string)
    default_user_settings = object({
      execution_role_arn = string
      security_group_ids = optional(list(string), [])
      default_landing_uri = optional(string, "")
      studio_web_portal = optional(string)
      auto_mount_home_efs = optional(string)
      jupyter_lab_app_settings = optional(object({
        default_resource_spec = optional(object({
          instance_type = optional(string, "")
          lifecycle_config_arn = optional(string, "")
          sagemaker_image_arn = optional(string, "")
          sagemaker_image_version_alias = optional(string, "")
          sagemaker_image_version_arn = optional(string, "")
        }))
        lifecycle_config_arns = optional(list(string), [])
        built_in_lifecycle_config_arn = optional(string, "")
        custom_images = optional(list(object({
          app_image_config_name = string
          image_name = string
          image_version_number = optional(number)
        })), [])
        code_repositories = optional(list(object({
          repository_url = string
        })), [])
        idle_settings = optional(object({
          lifecycle_management = optional(string)
          idle_timeout_in_minutes = number
          min_idle_timeout_in_minutes = number
          max_idle_timeout_in_minutes = number
        }))
        emr_settings = optional(object({
          assumable_role_arns = optional(list(string), [])
          execution_role_arns = optional(list(string), [])
        }))
      }))
      jupyter_server_app_settings = optional(object({
        default_resource_spec = optional(object({
          instance_type = optional(string, "")
          lifecycle_config_arn = optional(string, "")
          sagemaker_image_arn = optional(string, "")
          sagemaker_image_version_alias = optional(string, "")
          sagemaker_image_version_arn = optional(string, "")
        }))
        lifecycle_config_arns = optional(list(string), [])
        code_repositories = optional(list(object({
          repository_url = string
        })), [])
      }))
      kernel_gateway_app_settings = optional(object({
        default_resource_spec = optional(object({
          instance_type = optional(string, "")
          lifecycle_config_arn = optional(string, "")
          sagemaker_image_arn = optional(string, "")
          sagemaker_image_version_alias = optional(string, "")
          sagemaker_image_version_arn = optional(string, "")
        }))
        lifecycle_config_arns = optional(list(string), [])
        custom_images = optional(list(object({
          app_image_config_name = string
          image_name = string
          image_version_number = optional(number)
        })), [])
      }))
      code_editor_app_settings = optional(object({
        default_resource_spec = optional(object({
          instance_type = optional(string, "")
          lifecycle_config_arn = optional(string, "")
          sagemaker_image_arn = optional(string, "")
          sagemaker_image_version_alias = optional(string, "")
          sagemaker_image_version_arn = optional(string, "")
        }))
        lifecycle_config_arns = optional(list(string), [])
        built_in_lifecycle_config_arn = optional(string, "")
        custom_images = optional(list(object({
          app_image_config_name = string
          image_name = string
          image_version_number = optional(number)
        })), [])
        idle_settings = optional(object({
          lifecycle_management = optional(string)
          idle_timeout_in_minutes = number
          min_idle_timeout_in_minutes = number
          max_idle_timeout_in_minutes = number
        }))
      }))
      tensor_board_app_settings = optional(object({
        default_resource_spec = optional(object({
          instance_type = optional(string, "")
          lifecycle_config_arn = optional(string, "")
          sagemaker_image_arn = optional(string, "")
          sagemaker_image_version_alias = optional(string, "")
          sagemaker_image_version_arn = optional(string, "")
        }))
      }))
      r_session_app_settings = optional(object({
        default_resource_spec = optional(object({
          instance_type = optional(string, "")
          lifecycle_config_arn = optional(string, "")
          sagemaker_image_arn = optional(string, "")
          sagemaker_image_version_alias = optional(string, "")
          sagemaker_image_version_arn = optional(string, "")
        }))
        custom_images = optional(list(object({
          app_image_config_name = string
          image_name = string
          image_version_number = optional(number)
        })), [])
      }))
      r_studio_server_pro_app_settings = optional(object({
        access_status = optional(string, "")
        user_group = optional(string, "")
      }))
      canvas_app_settings = optional(object({
        direct_deploy_status = optional(string)
        emr_serverless_settings = optional(object({
          execution_role_arn = optional(string, "")
          status = optional(string, "")
        }))
        generative_ai_bedrock_role_arn = optional(string, "")
        identity_provider_oauth_settings = optional(list(object({
          data_source_name = optional(string, "")
          secret_arn = string
          status = optional(string, "")
        })), [])
        kendra_settings_status = optional(string)
        model_register_settings = optional(object({
          cross_account_model_register_role_arn = optional(string, "")
          status = optional(string, "")
        }))
        time_series_forecasting_settings = optional(object({
          amazon_forecast_role_arn = optional(string, "")
          status = optional(string, "")
        }))
        workspace_settings = optional(object({
          s3_artifact_path = optional(string, "")
          s3_kms_key_id = optional(string, "")
        }))
      }))
      sharing_settings = optional(object({
        notebook_output_option = optional(string)
        s3_kms_key_id = optional(string, "")
        s3_output_path = optional(string, "")
      }))
      space_storage_settings = optional(object({
        default_ebs_volume_size_in_gb = number
        maximum_ebs_volume_size_in_gb = number
      }))
      custom_file_system_configs = optional(list(object({
        efs_file_system_config = object({
          file_system_id = string
          file_system_path = string
        })
      })), [])
      custom_posix_user_config = optional(object({
        uid = number
        gid = number
      }))
      studio_web_portal_settings = optional(object({
        hidden_app_types = optional(list(string), [])
        hidden_instance_types = optional(list(string), [])
        hidden_ml_tools = optional(list(string), [])
      }))
    })
    default_space_settings = optional(object({
      execution_role_arn = string
      security_group_ids = optional(list(string), [])
      jupyter_lab_app_settings = optional(object({
        default_resource_spec = optional(object({
          instance_type = optional(string, "")
          lifecycle_config_arn = optional(string, "")
          sagemaker_image_arn = optional(string, "")
          sagemaker_image_version_alias = optional(string, "")
          sagemaker_image_version_arn = optional(string, "")
        }))
        lifecycle_config_arns = optional(list(string), [])
        built_in_lifecycle_config_arn = optional(string, "")
        custom_images = optional(list(object({
          app_image_config_name = string
          image_name = string
          image_version_number = optional(number)
        })), [])
        code_repositories = optional(list(object({
          repository_url = string
        })), [])
        idle_settings = optional(object({
          lifecycle_management = optional(string)
          idle_timeout_in_minutes = number
          min_idle_timeout_in_minutes = number
          max_idle_timeout_in_minutes = number
        }))
        emr_settings = optional(object({
          assumable_role_arns = optional(list(string), [])
          execution_role_arns = optional(list(string), [])
        }))
      }))
      jupyter_server_app_settings = optional(object({
        default_resource_spec = optional(object({
          instance_type = optional(string, "")
          lifecycle_config_arn = optional(string, "")
          sagemaker_image_arn = optional(string, "")
          sagemaker_image_version_alias = optional(string, "")
          sagemaker_image_version_arn = optional(string, "")
        }))
        lifecycle_config_arns = optional(list(string), [])
        code_repositories = optional(list(object({
          repository_url = string
        })), [])
      }))
      kernel_gateway_app_settings = optional(object({
        default_resource_spec = optional(object({
          instance_type = optional(string, "")
          lifecycle_config_arn = optional(string, "")
          sagemaker_image_arn = optional(string, "")
          sagemaker_image_version_alias = optional(string, "")
          sagemaker_image_version_arn = optional(string, "")
        }))
        lifecycle_config_arns = optional(list(string), [])
        custom_images = optional(list(object({
          app_image_config_name = string
          image_name = string
          image_version_number = optional(number)
        })), [])
      }))
      space_storage_settings = optional(object({
        default_ebs_volume_size_in_gb = number
        maximum_ebs_volume_size_in_gb = number
      }))
      custom_file_system_configs = optional(list(object({
        efs_file_system_config = object({
          file_system_id = string
          file_system_path = string
        })
      })), [])
      custom_posix_user_config = optional(object({
        uid = number
        gid = number
      }))
    }))
    domain_security_group_ids = optional(list(string), [])
    docker_settings = optional(object({
      enable_docker_access = optional(string, "")
      vpc_only_trusted_accounts = optional(list(string), [])
    }))
    execution_role_identity_config = optional(string)
    r_studio_server_pro_domain_settings = optional(object({
      domain_execution_role_arn = string
      r_studio_connect_url = optional(string, "")
      r_studio_package_manager_url = optional(string, "")
      default_resource_spec = optional(object({
        instance_type = optional(string, "")
        lifecycle_config_arn = optional(string, "")
        sagemaker_image_arn = optional(string, "")
        sagemaker_image_version_alias = optional(string, "")
        sagemaker_image_version_arn = optional(string, "")
      }))
    }))
    trusted_identity_propagation_status = optional(string)
    user_profiles = optional(list(object({
      user_profile_name = string
      single_sign_on_user_identifier = optional(string, "")
      single_sign_on_user_value = optional(string, "")
      user_settings = optional(object({
        execution_role_arn = string
        security_group_ids = optional(list(string), [])
        default_landing_uri = optional(string, "")
        studio_web_portal = optional(string)
        auto_mount_home_efs = optional(string)
        jupyter_lab_app_settings = optional(object({
          default_resource_spec = optional(object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          }))
          lifecycle_config_arns = optional(list(string), [])
          built_in_lifecycle_config_arn = optional(string, "")
          custom_images = optional(list(object({
            app_image_config_name = string
            image_name = string
            image_version_number = optional(number)
          })), [])
          code_repositories = optional(list(object({
            repository_url = string
          })), [])
          idle_settings = optional(object({
            lifecycle_management = optional(string)
            idle_timeout_in_minutes = number
            min_idle_timeout_in_minutes = number
            max_idle_timeout_in_minutes = number
          }))
          emr_settings = optional(object({
            assumable_role_arns = optional(list(string), [])
            execution_role_arns = optional(list(string), [])
          }))
        }))
        jupyter_server_app_settings = optional(object({
          default_resource_spec = optional(object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          }))
          lifecycle_config_arns = optional(list(string), [])
          code_repositories = optional(list(object({
            repository_url = string
          })), [])
        }))
        kernel_gateway_app_settings = optional(object({
          default_resource_spec = optional(object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          }))
          lifecycle_config_arns = optional(list(string), [])
          custom_images = optional(list(object({
            app_image_config_name = string
            image_name = string
            image_version_number = optional(number)
          })), [])
        }))
        code_editor_app_settings = optional(object({
          default_resource_spec = optional(object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          }))
          lifecycle_config_arns = optional(list(string), [])
          built_in_lifecycle_config_arn = optional(string, "")
          custom_images = optional(list(object({
            app_image_config_name = string
            image_name = string
            image_version_number = optional(number)
          })), [])
          idle_settings = optional(object({
            lifecycle_management = optional(string)
            idle_timeout_in_minutes = number
            min_idle_timeout_in_minutes = number
            max_idle_timeout_in_minutes = number
          }))
        }))
        tensor_board_app_settings = optional(object({
          default_resource_spec = optional(object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          }))
        }))
        r_session_app_settings = optional(object({
          default_resource_spec = optional(object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          }))
          custom_images = optional(list(object({
            app_image_config_name = string
            image_name = string
            image_version_number = optional(number)
          })), [])
        }))
        r_studio_server_pro_app_settings = optional(object({
          access_status = optional(string, "")
          user_group = optional(string, "")
        }))
        canvas_app_settings = optional(object({
          direct_deploy_status = optional(string)
          emr_serverless_settings = optional(object({
            execution_role_arn = optional(string, "")
            status = optional(string, "")
          }))
          generative_ai_bedrock_role_arn = optional(string, "")
          identity_provider_oauth_settings = optional(list(object({
            data_source_name = optional(string, "")
            secret_arn = string
            status = optional(string, "")
          })), [])
          kendra_settings_status = optional(string)
          model_register_settings = optional(object({
            cross_account_model_register_role_arn = optional(string, "")
            status = optional(string, "")
          }))
          time_series_forecasting_settings = optional(object({
            amazon_forecast_role_arn = optional(string, "")
            status = optional(string, "")
          }))
          workspace_settings = optional(object({
            s3_artifact_path = optional(string, "")
            s3_kms_key_id = optional(string, "")
          }))
        }))
        sharing_settings = optional(object({
          notebook_output_option = optional(string)
          s3_kms_key_id = optional(string, "")
          s3_output_path = optional(string, "")
        }))
        space_storage_settings = optional(object({
          default_ebs_volume_size_in_gb = number
          maximum_ebs_volume_size_in_gb = number
        }))
        custom_file_system_configs = optional(list(object({
          efs_file_system_config = object({
            file_system_id = string
            file_system_path = string
          })
        })), [])
        custom_posix_user_config = optional(object({
          uid = number
          gid = number
        }))
        studio_web_portal_settings = optional(object({
          hidden_app_types = optional(list(string), [])
          hidden_instance_types = optional(list(string), [])
          hidden_ml_tools = optional(list(string), [])
        }))
      }))
    })), [])
    spaces = optional(list(object({
      space_name = string
      display_name = optional(string, "")
      ownership_settings = optional(object({
        owner_user_profile_name = string
      }))
      space_sharing_settings = optional(object({
        sharing_type = string
      }))
      space_settings = optional(object({
        app_type = optional(string)
        jupyter_lab_app_settings = optional(object({
          default_resource_spec = object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          })
          code_repositories = optional(list(object({
            repository_url = string
          })), [])
          idle_settings = optional(object({
            idle_timeout_in_minutes = optional(number)
          }))
        }))
        code_editor_app_settings = optional(object({
          default_resource_spec = object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          })
          idle_settings = optional(object({
            idle_timeout_in_minutes = optional(number)
          }))
        }))
        jupyter_server_app_settings = optional(object({
          default_resource_spec = optional(object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          }))
          lifecycle_config_arns = optional(list(string), [])
          code_repositories = optional(list(object({
            repository_url = string
          })), [])
        }))
        kernel_gateway_app_settings = optional(object({
          default_resource_spec = optional(object({
            instance_type = optional(string, "")
            lifecycle_config_arn = optional(string, "")
            sagemaker_image_arn = optional(string, "")
            sagemaker_image_version_alias = optional(string, "")
            sagemaker_image_version_arn = optional(string, "")
          }))
          lifecycle_config_arns = optional(list(string), [])
          custom_images = optional(list(object({
            app_image_config_name = string
            image_name = string
            image_version_number = optional(number)
          })), [])
        }))
        custom_file_systems = optional(list(object({
          file_system_id = string
        })), [])
        space_storage_settings = optional(object({
          ebs_volume_size_in_gb = number
        }))
      }))
    })), [])
  })
}
