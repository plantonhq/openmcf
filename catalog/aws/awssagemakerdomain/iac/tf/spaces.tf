# ---------------------------------------------------------------------------
# Spaces — the folded aws_sagemaker_space satellites, one per spec.spaces
# entry, keyed by space name. A space's settings tree is DELIBERATELY
# different from the domain's default_space_settings (AWS uses distinct
# types: SpaceSettings vs DefaultSpaceSettings): it adds app_type and
# per-space Code Editor settings, requires a resource spec on each
# configured app, carries a timeout-only idle dial, and mounts existing EFS
# file systems directly by id.
#
# ownership_settings and space_sharing_settings travel together (provider
# RequiredWith, CEL-enforced) and are never sent on update — the provider
# does not support changing them after create.
# ---------------------------------------------------------------------------

resource "aws_sagemaker_space" "this" {
  for_each = { for s in var.spec.spaces : s.space_name => s }

  # A space's ownership references its owner profile by NAME (a plain
  # string, no resource linkage), so the engine sees no implicit edge --
  # without this, a space racing its owner profile fails at create.
  depends_on = [aws_sagemaker_user_profile.this]

  domain_id  = aws_sagemaker_domain.this.id
  space_name = each.value.space_name

  space_display_name = each.value.display_name != "" ? each.value.display_name : null

  dynamic "ownership_settings" {
    for_each = each.value.ownership_settings != null ? [each.value.ownership_settings] : []
    content {
      owner_user_profile_name = ownership_settings.value.owner_user_profile_name
    }
  }

  dynamic "space_sharing_settings" {
    for_each = each.value.space_sharing_settings != null ? [each.value.space_sharing_settings] : []
    content {
      sharing_type = space_sharing_settings.value.sharing_type
    }
  }

  dynamic "space_settings" {
    for_each = each.value.space_settings != null ? [each.value.space_settings] : []
    content {
      app_type = space_settings.value.app_type

      # --- JupyterLab (resource spec required on a space, CEL-enforced) ---
      dynamic "jupyter_lab_app_settings" {
        for_each = space_settings.value.jupyter_lab_app_settings != null ? [space_settings.value.jupyter_lab_app_settings] : []
        content {
          default_resource_spec {
            instance_type                 = jupyter_lab_app_settings.value.default_resource_spec.instance_type != "" ? jupyter_lab_app_settings.value.default_resource_spec.instance_type : null
            lifecycle_config_arn          = jupyter_lab_app_settings.value.default_resource_spec.lifecycle_config_arn != "" ? jupyter_lab_app_settings.value.default_resource_spec.lifecycle_config_arn : null
            sagemaker_image_arn           = jupyter_lab_app_settings.value.default_resource_spec.sagemaker_image_arn != "" ? jupyter_lab_app_settings.value.default_resource_spec.sagemaker_image_arn : null
            sagemaker_image_version_alias = jupyter_lab_app_settings.value.default_resource_spec.sagemaker_image_version_alias != "" ? jupyter_lab_app_settings.value.default_resource_spec.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = jupyter_lab_app_settings.value.default_resource_spec.sagemaker_image_version_arn != "" ? jupyter_lab_app_settings.value.default_resource_spec.sagemaker_image_version_arn : null
          }

          dynamic "code_repository" {
            for_each = jupyter_lab_app_settings.value.code_repositories
            content {
              repository_url = code_repository.value.repository_url
            }
          }

          # The space idle dial carries only the timeout — no
          # lifecycle_management switch and no min/max guardrails (those
          # exist only on the domain/user plane). The module forwards it
          # verbatim; legality is resolved server-side — CreateSpace 400s
          # ("Idle Shutdown is disabled for this space") when idle shutdown
          # resolves DISABLED for the space through the domain/owner-profile
          # inheritance chain (live-verified 2026-08-13). The spec's
          # space_idle_requires_owner_lifecycle_* CELs catch the in-manifest
          # contradiction before the API does.
          dynamic "app_lifecycle_management" {
            for_each = jupyter_lab_app_settings.value.idle_settings != null ? [jupyter_lab_app_settings.value.idle_settings] : []
            content {
              idle_settings {
                idle_timeout_in_minutes = app_lifecycle_management.value.idle_timeout_in_minutes
              }
            }
          }
        }
      }

      # --- Code Editor (the space form) ---
      dynamic "code_editor_app_settings" {
        for_each = space_settings.value.code_editor_app_settings != null ? [space_settings.value.code_editor_app_settings] : []
        content {
          default_resource_spec {
            instance_type                 = code_editor_app_settings.value.default_resource_spec.instance_type != "" ? code_editor_app_settings.value.default_resource_spec.instance_type : null
            lifecycle_config_arn          = code_editor_app_settings.value.default_resource_spec.lifecycle_config_arn != "" ? code_editor_app_settings.value.default_resource_spec.lifecycle_config_arn : null
            sagemaker_image_arn           = code_editor_app_settings.value.default_resource_spec.sagemaker_image_arn != "" ? code_editor_app_settings.value.default_resource_spec.sagemaker_image_arn : null
            sagemaker_image_version_alias = code_editor_app_settings.value.default_resource_spec.sagemaker_image_version_alias != "" ? code_editor_app_settings.value.default_resource_spec.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = code_editor_app_settings.value.default_resource_spec.sagemaker_image_version_arn != "" ? code_editor_app_settings.value.default_resource_spec.sagemaker_image_version_arn : null
          }

          dynamic "app_lifecycle_management" {
            for_each = code_editor_app_settings.value.idle_settings != null ? [code_editor_app_settings.value.idle_settings] : []
            content {
              idle_settings {
                idle_timeout_in_minutes = app_lifecycle_management.value.idle_timeout_in_minutes
              }
            }
          }
        }
      }

      # --- classic Jupyter Server (same shape as the domain baseline;
      # default_resource_spec required on a space, CEL-enforced) ---
      dynamic "jupyter_server_app_settings" {
        for_each = space_settings.value.jupyter_server_app_settings != null ? [space_settings.value.jupyter_server_app_settings] : []
        content {
          lifecycle_config_arns = length(jupyter_server_app_settings.value.lifecycle_config_arns) > 0 ? jupyter_server_app_settings.value.lifecycle_config_arns : null

          default_resource_spec {
            instance_type                 = jupyter_server_app_settings.value.default_resource_spec.instance_type != "" ? jupyter_server_app_settings.value.default_resource_spec.instance_type : null
            lifecycle_config_arn          = jupyter_server_app_settings.value.default_resource_spec.lifecycle_config_arn != "" ? jupyter_server_app_settings.value.default_resource_spec.lifecycle_config_arn : null
            sagemaker_image_arn           = jupyter_server_app_settings.value.default_resource_spec.sagemaker_image_arn != "" ? jupyter_server_app_settings.value.default_resource_spec.sagemaker_image_arn : null
            sagemaker_image_version_alias = jupyter_server_app_settings.value.default_resource_spec.sagemaker_image_version_alias != "" ? jupyter_server_app_settings.value.default_resource_spec.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = jupyter_server_app_settings.value.default_resource_spec.sagemaker_image_version_arn != "" ? jupyter_server_app_settings.value.default_resource_spec.sagemaker_image_version_arn : null
          }

          dynamic "code_repository" {
            for_each = jupyter_server_app_settings.value.code_repositories
            content {
              repository_url = code_repository.value.repository_url
            }
          }
        }
      }

      # --- KernelGateway (same shape as the domain baseline) ---
      dynamic "kernel_gateway_app_settings" {
        for_each = space_settings.value.kernel_gateway_app_settings != null ? [space_settings.value.kernel_gateway_app_settings] : []
        content {
          lifecycle_config_arns = length(kernel_gateway_app_settings.value.lifecycle_config_arns) > 0 ? kernel_gateway_app_settings.value.lifecycle_config_arns : null

          default_resource_spec {
            instance_type                 = kernel_gateway_app_settings.value.default_resource_spec.instance_type != "" ? kernel_gateway_app_settings.value.default_resource_spec.instance_type : null
            lifecycle_config_arn          = kernel_gateway_app_settings.value.default_resource_spec.lifecycle_config_arn != "" ? kernel_gateway_app_settings.value.default_resource_spec.lifecycle_config_arn : null
            sagemaker_image_arn           = kernel_gateway_app_settings.value.default_resource_spec.sagemaker_image_arn != "" ? kernel_gateway_app_settings.value.default_resource_spec.sagemaker_image_arn : null
            sagemaker_image_version_alias = kernel_gateway_app_settings.value.default_resource_spec.sagemaker_image_version_alias != "" ? kernel_gateway_app_settings.value.default_resource_spec.sagemaker_image_version_alias : null
            sagemaker_image_version_arn   = kernel_gateway_app_settings.value.default_resource_spec.sagemaker_image_version_arn != "" ? kernel_gateway_app_settings.value.default_resource_spec.sagemaker_image_version_arn : null
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

      # --- mounted EFS file systems (by id — the space form has no
      # per-mount path, unlike the domain baseline's config) ---
      dynamic "custom_file_system" {
        for_each = space_settings.value.custom_file_systems
        content {
          efs_file_system {
            file_system_id = custom_file_system.value.file_system_id
          }
        }
      }

      # --- the space's EBS volume (a single concrete size) ---
      dynamic "space_storage_settings" {
        for_each = space_settings.value.space_storage_settings != null ? [space_settings.value.space_storage_settings] : []
        content {
          ebs_storage_settings {
            ebs_volume_size_in_gb = space_storage_settings.value.ebs_volume_size_in_gb
          }
        }
      }
    }
  }

  tags = local.tags
}
