# Amazon SageMaker Domain.
#
# One provider resource carries the whole surface. Three provider behaviors
# shape the wiring below:
#
#   - Optional nested blocks are presence-driven: an absent block is the
#     feature's AWS-default state, so each spec message maps to a `dynamic`
#     block guarded by its presence -- never an always-emitted block with
#     zero values, which would pin defaults and create phantom diffs.
#   - Presence-carrying scalars (proto3 optionals arriving as null) pass
#     through as-is: null means "not sent", so AWS applies its own default
#     and the plan stays clean. Plain strings use "" as the absent sentinel
#     and are converted to null here.
#   - The same `default_resource_spec` shape repeats across every app type
#     (JupyterLab, classic Jupyter Server, KernelGateway, Code Editor,
#     TensorBoard, RSession, RStudio). The provider models each as its own
#     block type, so each app-settings block below carries its own copy.
resource "aws_sagemaker_domain" "this" {
  domain_name = local.domain_name

  # Identity plane: IAM vs IAM Identity Center (SSO). ForceNew.
  auth_mode = var.spec.auth_mode

  # Network plane. VPC + subnets are ForceNew; the access type (public vs
  # VpcOnly) is updatable in place.
  vpc_id                  = var.spec.vpc_id
  subnet_ids              = var.spec.subnet_ids
  app_network_access_type = var.spec.app_network_access_type

  # Only honored by AWS when RStudio is configured (the spec enforces the
  # pairing, so a silent no-op cannot reach this argument).
  app_security_group_management = var.spec.app_security_group_management

  # EFS home-directory encryption key. ForceNew.
  kms_key_id = var.spec.kms_key_id != "" ? var.spec.kms_key_id : null

  # Whether the domain's tags propagate to apps/spaces/user profiles created
  # inside it. Null defers to AWS's default (DISABLED).
  tag_propagation = var.spec.tag_propagation

  tags = local.tags

  # What happens to the domain's auto-created EFS file system on destroy.
  # AWS's default is Retain, which leaves a billing orphan behind -- the spec
  # surfaces the decision so ephemeral domains can opt into Delete.
  dynamic "retention_policy" {
    for_each = var.spec.home_efs_retention_policy != null ? [var.spec.home_efs_retention_policy] : []
    content {
      home_efs_file_system = retention_policy.value
    }
  }

  # --- domain_settings (domain-scoped administration) ---

  dynamic "domain_settings" {
    for_each = local.has_domain_settings ? [1] : []
    content {
      security_group_ids             = length(var.spec.domain_security_group_ids) > 0 ? var.spec.domain_security_group_ids : null
      execution_role_identity_config = var.spec.execution_role_identity_config

      dynamic "docker_settings" {
        for_each = var.spec.docker_settings != null ? [var.spec.docker_settings] : []
        content {
          enable_docker_access      = docker_settings.value.enable_docker_access != "" ? docker_settings.value.enable_docker_access : null
          vpc_only_trusted_accounts = length(docker_settings.value.vpc_only_trusted_accounts) > 0 ? docker_settings.value.vpc_only_trusted_accounts : null
        }
      }

      # RStudio (Posit) Workbench activation: configuring the domain execution
      # role is what turns the RStudio app plane on for the domain.
      dynamic "r_studio_server_pro_domain_settings" {
        for_each = var.spec.r_studio_server_pro_domain_settings != null ? [var.spec.r_studio_server_pro_domain_settings] : []
        content {
          domain_execution_role_arn    = r_studio_server_pro_domain_settings.value.domain_execution_role_arn
          r_studio_connect_url         = r_studio_server_pro_domain_settings.value.r_studio_connect_url != "" ? r_studio_server_pro_domain_settings.value.r_studio_connect_url : null
          r_studio_package_manager_url = r_studio_server_pro_domain_settings.value.r_studio_package_manager_url != "" ? r_studio_server_pro_domain_settings.value.r_studio_package_manager_url : null

          dynamic "default_resource_spec" {
            for_each = r_studio_server_pro_domain_settings.value.default_resource_spec != null ? [r_studio_server_pro_domain_settings.value.default_resource_spec] : []
            content {
              instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
              lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
              sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
              sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
              sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
            }
          }
        }
      }

      # Trusted identity propagation: the spec's CEL guarantees ENABLED only
      # ever reaches AWS on an SSO domain (AWS rejects it under IAM auth).
      dynamic "trusted_identity_propagation_settings" {
        for_each = var.spec.trusted_identity_propagation_status != null ? [var.spec.trusted_identity_propagation_status] : []
        content {
          status = trusted_identity_propagation_settings.value
        }
      }
    }
  }

  # --- default_user_settings (the per-user baseline) ---

  default_user_settings {
    execution_role      = local.dus.execution_role_arn
    security_groups     = length(local.dus.security_group_ids) > 0 ? local.dus.security_group_ids : null
    default_landing_uri = local.dus.default_landing_uri != "" ? local.dus.default_landing_uri : null
    studio_web_portal   = local.dus.studio_web_portal
    auto_mount_home_efs = local.dus.auto_mount_home_efs

    # --- JupyterLab (the primary Studio IDE) ---

    dynamic "jupyter_lab_app_settings" {
      for_each = local.dus.jupyter_lab_app_settings != null ? [local.dus.jupyter_lab_app_settings] : []
      content {
        lifecycle_config_arns         = length(jupyter_lab_app_settings.value.lifecycle_config_arns) > 0 ? jupyter_lab_app_settings.value.lifecycle_config_arns : null
        built_in_lifecycle_config_arn = jupyter_lab_app_settings.value.built_in_lifecycle_config_arn != "" ? jupyter_lab_app_settings.value.built_in_lifecycle_config_arn : null

        dynamic "default_resource_spec" {
          for_each = jupyter_lab_app_settings.value.default_resource_spec != null ? [jupyter_lab_app_settings.value.default_resource_spec] : []
          content {
            instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
            lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
            sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
            sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
          }
        }

        dynamic "custom_image" {
          for_each = jupyter_lab_app_settings.value.custom_images
          content {
            app_image_config_name = custom_image.value.app_image_config_name
            image_name            = custom_image.value.image_name
            image_version_number  = custom_image.value.image_version_number
          }
        }

        dynamic "code_repository" {
          for_each = jupyter_lab_app_settings.value.code_repositories
          content {
            repository_url = code_repository.value.repository_url
          }
        }

        # The spec folds idle_settings directly under the app settings and
        # makes block presence the enable switch; the provider nests them
        # inside a single-purpose app_lifecycle_management wrapper, both
        # reconstructed here. lifecycle_management defaults to ENABLED when
        # the block is present; an explicit DISABLED keeps the timeouts as
        # published guardrails without enforcing auto-shutdown. All three timeouts
        # are required by the live API whenever the block is sent (absent
        # members transmit as 0 and AWS rejects them), so they pass through
        # unconditionally.
        dynamic "app_lifecycle_management" {
          for_each = jupyter_lab_app_settings.value.idle_settings != null ? [jupyter_lab_app_settings.value.idle_settings] : []
          content {
            idle_settings {
              lifecycle_management        = coalesce(app_lifecycle_management.value.lifecycle_management, "ENABLED")
              idle_timeout_in_minutes     = app_lifecycle_management.value.idle_timeout_in_minutes
              min_idle_timeout_in_minutes = app_lifecycle_management.value.min_idle_timeout_in_minutes
              max_idle_timeout_in_minutes = app_lifecycle_management.value.max_idle_timeout_in_minutes
            }
          }
        }

        dynamic "emr_settings" {
          for_each = jupyter_lab_app_settings.value.emr_settings != null ? [jupyter_lab_app_settings.value.emr_settings] : []
          content {
            assumable_role_arns = length(emr_settings.value.assumable_role_arns) > 0 ? emr_settings.value.assumable_role_arns : null
            execution_role_arns = length(emr_settings.value.execution_role_arns) > 0 ? emr_settings.value.execution_role_arns : null
          }
        }
      }
    }

    # --- classic Jupyter Server (Studio Classic) ---

    dynamic "jupyter_server_app_settings" {
      for_each = local.dus.jupyter_server_app_settings != null ? [local.dus.jupyter_server_app_settings] : []
      content {
        lifecycle_config_arns = length(jupyter_server_app_settings.value.lifecycle_config_arns) > 0 ? jupyter_server_app_settings.value.lifecycle_config_arns : null

        dynamic "default_resource_spec" {
          for_each = jupyter_server_app_settings.value.default_resource_spec != null ? [jupyter_server_app_settings.value.default_resource_spec] : []
          content {
            instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
            lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
            sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
            sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
          }
        }

        dynamic "code_repository" {
          for_each = jupyter_server_app_settings.value.code_repositories
          content {
            repository_url = code_repository.value.repository_url
          }
        }
      }
    }

    # --- KernelGateway (bring-your-own-image kernels) ---

    dynamic "kernel_gateway_app_settings" {
      for_each = local.dus.kernel_gateway_app_settings != null ? [local.dus.kernel_gateway_app_settings] : []
      content {
        lifecycle_config_arns = length(kernel_gateway_app_settings.value.lifecycle_config_arns) > 0 ? kernel_gateway_app_settings.value.lifecycle_config_arns : null

        dynamic "default_resource_spec" {
          for_each = kernel_gateway_app_settings.value.default_resource_spec != null ? [kernel_gateway_app_settings.value.default_resource_spec] : []
          content {
            instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
            lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
            sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
            sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
          }
        }

        dynamic "custom_image" {
          for_each = kernel_gateway_app_settings.value.custom_images
          content {
            app_image_config_name = custom_image.value.app_image_config_name
            image_name            = custom_image.value.image_name
            image_version_number  = custom_image.value.image_version_number
          }
        }
      }
    }

    # --- Code Editor (SageMaker's VS Code / Code-OSS IDE) ---

    dynamic "code_editor_app_settings" {
      for_each = local.dus.code_editor_app_settings != null ? [local.dus.code_editor_app_settings] : []
      content {
        lifecycle_config_arns         = length(code_editor_app_settings.value.lifecycle_config_arns) > 0 ? code_editor_app_settings.value.lifecycle_config_arns : null
        built_in_lifecycle_config_arn = code_editor_app_settings.value.built_in_lifecycle_config_arn != "" ? code_editor_app_settings.value.built_in_lifecycle_config_arn : null

        dynamic "default_resource_spec" {
          for_each = code_editor_app_settings.value.default_resource_spec != null ? [code_editor_app_settings.value.default_resource_spec] : []
          content {
            instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
            lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
            sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
            sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
          }
        }

        dynamic "custom_image" {
          for_each = code_editor_app_settings.value.custom_images
          content {
            app_image_config_name = custom_image.value.app_image_config_name
            image_name            = custom_image.value.image_name
            image_version_number  = custom_image.value.image_version_number
          }
        }

        dynamic "app_lifecycle_management" {
          for_each = code_editor_app_settings.value.idle_settings != null ? [code_editor_app_settings.value.idle_settings] : []
          content {
            idle_settings {
              lifecycle_management        = coalesce(app_lifecycle_management.value.lifecycle_management, "ENABLED")
              idle_timeout_in_minutes     = app_lifecycle_management.value.idle_timeout_in_minutes
              min_idle_timeout_in_minutes = app_lifecycle_management.value.min_idle_timeout_in_minutes
              max_idle_timeout_in_minutes = app_lifecycle_management.value.max_idle_timeout_in_minutes
            }
          }
        }
      }
    }

    # --- TensorBoard ---

    dynamic "tensor_board_app_settings" {
      for_each = local.dus.tensor_board_app_settings != null ? [local.dus.tensor_board_app_settings] : []
      content {
        dynamic "default_resource_spec" {
          for_each = tensor_board_app_settings.value.default_resource_spec != null ? [tensor_board_app_settings.value.default_resource_spec] : []
          content {
            instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
            lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
            sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
            sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
          }
        }
      }
    }

    # --- RSession (R kernels behind RStudio) ---

    dynamic "r_session_app_settings" {
      for_each = local.dus.r_session_app_settings != null ? [local.dus.r_session_app_settings] : []
      content {
        dynamic "default_resource_spec" {
          for_each = r_session_app_settings.value.default_resource_spec != null ? [r_session_app_settings.value.default_resource_spec] : []
          content {
            instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
            lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
            sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
            sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
          }
        }

        dynamic "custom_image" {
          for_each = r_session_app_settings.value.custom_images
          content {
            app_image_config_name = custom_image.value.app_image_config_name
            image_name            = custom_image.value.image_name
            image_version_number  = custom_image.value.image_version_number
          }
        }
      }
    }

    # --- RStudio Server Pro (per-user access) ---

    dynamic "r_studio_server_pro_app_settings" {
      for_each = local.dus.r_studio_server_pro_app_settings != null ? [local.dus.r_studio_server_pro_app_settings] : []
      content {
        access_status = r_studio_server_pro_app_settings.value.access_status != "" ? r_studio_server_pro_app_settings.value.access_status : null
        user_group    = r_studio_server_pro_app_settings.value.user_group != "" ? r_studio_server_pro_app_settings.value.user_group : null
      }
    }

    # --- Canvas (no-code ML) ---
    #
    # The spec flattens the provider's two single-field wrappers
    # (direct_deploy_settings.status and kendra_settings.status) to plain
    # status scalars; the wrappers are reconstructed here.

    dynamic "canvas_app_settings" {
      for_each = local.dus.canvas_app_settings != null ? [local.dus.canvas_app_settings] : []
      content {
        dynamic "direct_deploy_settings" {
          for_each = canvas_app_settings.value.direct_deploy_status != null ? [canvas_app_settings.value.direct_deploy_status] : []
          content {
            status = direct_deploy_settings.value
          }
        }

        dynamic "emr_serverless_settings" {
          for_each = canvas_app_settings.value.emr_serverless_settings != null ? [canvas_app_settings.value.emr_serverless_settings] : []
          content {
            execution_role_arn = emr_serverless_settings.value.execution_role_arn != "" ? emr_serverless_settings.value.execution_role_arn : null
            status             = emr_serverless_settings.value.status != "" ? emr_serverless_settings.value.status : null
          }
        }

        # Setting the Bedrock role is what enables Canvas generative AI; the
        # provider wraps the single role in a generative_ai_settings block.
        dynamic "generative_ai_settings" {
          for_each = canvas_app_settings.value.generative_ai_bedrock_role_arn != "" ? [canvas_app_settings.value.generative_ai_bedrock_role_arn] : []
          content {
            amazon_bedrock_role_arn = generative_ai_settings.value
          }
        }

        dynamic "identity_provider_oauth_settings" {
          for_each = canvas_app_settings.value.identity_provider_oauth_settings
          content {
            data_source_name = identity_provider_oauth_settings.value.data_source_name != "" ? identity_provider_oauth_settings.value.data_source_name : null
            secret_arn       = identity_provider_oauth_settings.value.secret_arn
            status           = identity_provider_oauth_settings.value.status != "" ? identity_provider_oauth_settings.value.status : null
          }
        }

        dynamic "kendra_settings" {
          for_each = canvas_app_settings.value.kendra_settings_status != null ? [canvas_app_settings.value.kendra_settings_status] : []
          content {
            status = kendra_settings.value
          }
        }

        dynamic "model_register_settings" {
          for_each = canvas_app_settings.value.model_register_settings != null ? [canvas_app_settings.value.model_register_settings] : []
          content {
            cross_account_model_register_role_arn = model_register_settings.value.cross_account_model_register_role_arn != "" ? model_register_settings.value.cross_account_model_register_role_arn : null
            status                                = model_register_settings.value.status != "" ? model_register_settings.value.status : null
          }
        }

        dynamic "time_series_forecasting_settings" {
          for_each = canvas_app_settings.value.time_series_forecasting_settings != null ? [canvas_app_settings.value.time_series_forecasting_settings] : []
          content {
            amazon_forecast_role_arn = time_series_forecasting_settings.value.amazon_forecast_role_arn != "" ? time_series_forecasting_settings.value.amazon_forecast_role_arn : null
            status                   = time_series_forecasting_settings.value.status != "" ? time_series_forecasting_settings.value.status : null
          }
        }

        dynamic "workspace_settings" {
          for_each = canvas_app_settings.value.workspace_settings != null ? [canvas_app_settings.value.workspace_settings] : []
          content {
            s3_artifact_path = workspace_settings.value.s3_artifact_path != "" ? workspace_settings.value.s3_artifact_path : null
            s3_kms_key_id    = workspace_settings.value.s3_kms_key_id != "" ? workspace_settings.value.s3_kms_key_id : null
          }
        }
      }
    }

    # --- sharing, storage, custom mounts, POSIX identity, UI governance ---

    dynamic "sharing_settings" {
      for_each = local.dus.sharing_settings != null ? [local.dus.sharing_settings] : []
      content {
        notebook_output_option = sharing_settings.value.notebook_output_option
        s3_kms_key_id          = sharing_settings.value.s3_kms_key_id != "" ? sharing_settings.value.s3_kms_key_id : null
        s3_output_path         = sharing_settings.value.s3_output_path != "" ? sharing_settings.value.s3_output_path : null
      }
    }

    # The spec flattens the provider's single-purpose default_ebs_storage_settings
    # wrapper; it is reconstructed here.
    dynamic "space_storage_settings" {
      for_each = local.dus.space_storage_settings != null ? [local.dus.space_storage_settings] : []
      content {
        default_ebs_storage_settings {
          default_ebs_volume_size_in_gb = space_storage_settings.value.default_ebs_volume_size_in_gb
          maximum_ebs_volume_size_in_gb = space_storage_settings.value.maximum_ebs_volume_size_in_gb
        }
      }
    }

    dynamic "custom_file_system_config" {
      for_each = local.dus.custom_file_system_configs
      content {
        efs_file_system_config {
          file_system_id   = custom_file_system_config.value.efs_file_system_config.file_system_id
          file_system_path = custom_file_system_config.value.efs_file_system_config.file_system_path
        }
      }
    }

    dynamic "custom_posix_user_config" {
      for_each = local.dus.custom_posix_user_config != null ? [local.dus.custom_posix_user_config] : []
      content {
        uid = custom_posix_user_config.value.uid
        gid = custom_posix_user_config.value.gid
      }
    }

    dynamic "studio_web_portal_settings" {
      for_each = local.dus.studio_web_portal_settings != null ? [local.dus.studio_web_portal_settings] : []
      content {
        hidden_app_types      = length(studio_web_portal_settings.value.hidden_app_types) > 0 ? studio_web_portal_settings.value.hidden_app_types : null
        hidden_instance_types = length(studio_web_portal_settings.value.hidden_instance_types) > 0 ? studio_web_portal_settings.value.hidden_instance_types : null
        hidden_ml_tools       = length(studio_web_portal_settings.value.hidden_ml_tools) > 0 ? studio_web_portal_settings.value.hidden_ml_tools : null
      }
    }
  }

  # --- default_space_settings (the shared-space baseline) ---

  dynamic "default_space_settings" {
    for_each = var.spec.default_space_settings != null ? [var.spec.default_space_settings] : []
    content {
      execution_role  = default_space_settings.value.execution_role_arn
      security_groups = length(default_space_settings.value.security_group_ids) > 0 ? default_space_settings.value.security_group_ids : null

      dynamic "jupyter_lab_app_settings" {
        for_each = default_space_settings.value.jupyter_lab_app_settings != null ? [default_space_settings.value.jupyter_lab_app_settings] : []
        content {
          lifecycle_config_arns         = length(jupyter_lab_app_settings.value.lifecycle_config_arns) > 0 ? jupyter_lab_app_settings.value.lifecycle_config_arns : null
          built_in_lifecycle_config_arn = jupyter_lab_app_settings.value.built_in_lifecycle_config_arn != "" ? jupyter_lab_app_settings.value.built_in_lifecycle_config_arn : null

          dynamic "default_resource_spec" {
            for_each = jupyter_lab_app_settings.value.default_resource_spec != null ? [jupyter_lab_app_settings.value.default_resource_spec] : []
            content {
              instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
              lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
              sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
              sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
              sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
            }
          }

          dynamic "custom_image" {
            for_each = jupyter_lab_app_settings.value.custom_images
            content {
              app_image_config_name = custom_image.value.app_image_config_name
              image_name            = custom_image.value.image_name
              image_version_number  = custom_image.value.image_version_number
            }
          }

          dynamic "code_repository" {
            for_each = jupyter_lab_app_settings.value.code_repositories
            content {
              repository_url = code_repository.value.repository_url
            }
          }

          dynamic "app_lifecycle_management" {
            for_each = jupyter_lab_app_settings.value.idle_settings != null ? [jupyter_lab_app_settings.value.idle_settings] : []
            content {
              idle_settings {
                lifecycle_management        = coalesce(app_lifecycle_management.value.lifecycle_management, "ENABLED")
                idle_timeout_in_minutes     = app_lifecycle_management.value.idle_timeout_in_minutes
                min_idle_timeout_in_minutes = app_lifecycle_management.value.min_idle_timeout_in_minutes
                max_idle_timeout_in_minutes = app_lifecycle_management.value.max_idle_timeout_in_minutes
              }
            }
          }

          dynamic "emr_settings" {
            for_each = jupyter_lab_app_settings.value.emr_settings != null ? [jupyter_lab_app_settings.value.emr_settings] : []
            content {
              assumable_role_arns = length(emr_settings.value.assumable_role_arns) > 0 ? emr_settings.value.assumable_role_arns : null
              execution_role_arns = length(emr_settings.value.execution_role_arns) > 0 ? emr_settings.value.execution_role_arns : null
            }
          }
        }
      }

      dynamic "jupyter_server_app_settings" {
        for_each = default_space_settings.value.jupyter_server_app_settings != null ? [default_space_settings.value.jupyter_server_app_settings] : []
        content {
          lifecycle_config_arns = length(jupyter_server_app_settings.value.lifecycle_config_arns) > 0 ? jupyter_server_app_settings.value.lifecycle_config_arns : null

          dynamic "default_resource_spec" {
            for_each = jupyter_server_app_settings.value.default_resource_spec != null ? [jupyter_server_app_settings.value.default_resource_spec] : []
            content {
              instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
              lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
              sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
              sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
              sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
            }
          }

          dynamic "code_repository" {
            for_each = jupyter_server_app_settings.value.code_repositories
            content {
              repository_url = code_repository.value.repository_url
            }
          }
        }
      }

      dynamic "kernel_gateway_app_settings" {
        for_each = default_space_settings.value.kernel_gateway_app_settings != null ? [default_space_settings.value.kernel_gateway_app_settings] : []
        content {
          lifecycle_config_arns = length(kernel_gateway_app_settings.value.lifecycle_config_arns) > 0 ? kernel_gateway_app_settings.value.lifecycle_config_arns : null

          dynamic "default_resource_spec" {
            for_each = kernel_gateway_app_settings.value.default_resource_spec != null ? [kernel_gateway_app_settings.value.default_resource_spec] : []
            content {
              instance_type                 = default_resource_spec.value.instance_type != "" ? default_resource_spec.value.instance_type : null
              lifecycle_config_arn          = default_resource_spec.value.lifecycle_config_arn != "" ? default_resource_spec.value.lifecycle_config_arn : null
              sagemaker_image_arn           = default_resource_spec.value.sagemaker_image_arn != "" ? default_resource_spec.value.sagemaker_image_arn : null
              sagemaker_image_version_alias = default_resource_spec.value.sagemaker_image_version_alias != "" ? default_resource_spec.value.sagemaker_image_version_alias : null
              sagemaker_image_version_arn   = default_resource_spec.value.sagemaker_image_version_arn != "" ? default_resource_spec.value.sagemaker_image_version_arn : null
            }
          }

          dynamic "custom_image" {
            for_each = kernel_gateway_app_settings.value.custom_images
            content {
              app_image_config_name = custom_image.value.app_image_config_name
              image_name            = custom_image.value.image_name
              image_version_number  = custom_image.value.image_version_number
            }
          }
        }
      }

      dynamic "space_storage_settings" {
        for_each = default_space_settings.value.space_storage_settings != null ? [default_space_settings.value.space_storage_settings] : []
        content {
          default_ebs_storage_settings {
            default_ebs_volume_size_in_gb = space_storage_settings.value.default_ebs_volume_size_in_gb
            maximum_ebs_volume_size_in_gb = space_storage_settings.value.maximum_ebs_volume_size_in_gb
          }
        }
      }

      dynamic "custom_file_system_config" {
        for_each = default_space_settings.value.custom_file_system_configs
        content {
          efs_file_system_config {
            file_system_id   = custom_file_system_config.value.efs_file_system_config.file_system_id
            file_system_path = custom_file_system_config.value.efs_file_system_config.file_system_path
          }
        }
      }

      dynamic "custom_posix_user_config" {
        for_each = default_space_settings.value.custom_posix_user_config != null ? [default_space_settings.value.custom_posix_user_config] : []
        content {
          uid = custom_posix_user_config.value.uid
          gid = custom_posix_user_config.value.gid
        }
      }
    }
  }
}
