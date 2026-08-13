# AWS Batch job definition.
#
# Every meaningful change registers a NEW revision (revisions are immutable
# in AWS); with deregister_on_new_revision at its default the previous
# revision is deregistered so exactly one ACTIVE revision tracks this
# resource. Two workload arms are modeled, both of AWS type "container":
# single-container ECS-based jobs (containerProperties, a JSON document
# built in locals.tf) and Batch-on-EKS pod jobs (eksProperties, the typed
# block below). The spec guarantees exactly one arm is set. Multinode
# (nodeProperties, type "multinode") and multi-container ECS
# (ecsProperties) remain unmodeled long-tail shapes.
resource "aws_batch_job_definition" "this" {
  # The cloud name comes from metadata.name (the catalog naming basis) --
  # revisions register under this name in both engines.
  name = var.metadata.name
  # EKS pod jobs are ALSO type "container" in AWS -- the arm is selected by
  # which properties document is present, not by the type.
  type = "container"

  container_properties = var.spec.container != null ? jsonencode(local.container_properties) : null

  # Batch-on-EKS pod jobs. Field-for-field send parity with the Pulumi
  # module: presence-typed spec fields (host_network, security context
  # run_as_user/run_as_group/allow_privilege_escalation) pass through as
  # null when unset so AWS's own defaults apply; plain bools are always
  # sent (state-pinned) when their enclosing block renders.
  dynamic "eks_properties" {
    for_each = var.spec.eks != null ? [var.spec.eks] : []
    content {
      pod_properties {
        # Unset means AWS's default, which is TRUE for Batch pods.
        host_network            = eks_properties.value.host_network
        dns_policy              = eks_properties.value.dns_policy != "" ? eks_properties.value.dns_policy : null
        service_account_name    = eks_properties.value.service_account_name != "" ? eks_properties.value.service_account_name : null
        share_process_namespace = eks_properties.value.share_process_namespace

        dynamic "containers" {
          for_each = eks_properties.value.containers
          content {
            image             = containers.value.image
            name              = containers.value.name != "" ? containers.value.name : null
            command           = length(containers.value.command) > 0 ? containers.value.command : null
            args              = length(containers.value.args) > 0 ? containers.value.args : null
            image_pull_policy = containers.value.image_pull_policy != "" ? containers.value.image_pull_policy : null

            # Sorted for determinism (the provider stores env as a set).
            dynamic "env" {
              for_each = { for name in sort(keys(containers.value.env)) : name => containers.value.env[name] }
              content {
                name  = env.key
                value = env.value
              }
            }

            dynamic "resources" {
              for_each = containers.value.resources != null ? [containers.value.resources] : []
              content {
                limits   = length(resources.value.limits) > 0 ? resources.value.limits : null
                requests = length(resources.value.requests) > 0 ? resources.value.requests : null
              }
            }

            dynamic "security_context" {
              for_each = containers.value.security_context != null ? [containers.value.security_context] : []
              content {
                run_as_user                = security_context.value.run_as_user
                run_as_group               = security_context.value.run_as_group
                run_as_non_root            = security_context.value.run_as_non_root
                allow_privilege_escalation = security_context.value.allow_privilege_escalation
                privileged                 = security_context.value.privileged
                read_only_root_file_system = security_context.value.read_only_root_file_system
              }
            }

            dynamic "volume_mounts" {
              for_each = containers.value.volume_mounts
              content {
                name       = volume_mounts.value.name
                mount_path = volume_mounts.value.mount_path
                read_only  = volume_mounts.value.read_only
              }
            }
          }
        }

        dynamic "init_containers" {
          for_each = eks_properties.value.init_containers
          content {
            image             = init_containers.value.image
            name              = init_containers.value.name != "" ? init_containers.value.name : null
            command           = length(init_containers.value.command) > 0 ? init_containers.value.command : null
            args              = length(init_containers.value.args) > 0 ? init_containers.value.args : null
            image_pull_policy = init_containers.value.image_pull_policy != "" ? init_containers.value.image_pull_policy : null

            dynamic "env" {
              for_each = { for name in sort(keys(init_containers.value.env)) : name => init_containers.value.env[name] }
              content {
                name  = env.key
                value = env.value
              }
            }

            dynamic "resources" {
              for_each = init_containers.value.resources != null ? [init_containers.value.resources] : []
              content {
                limits   = length(resources.value.limits) > 0 ? resources.value.limits : null
                requests = length(resources.value.requests) > 0 ? resources.value.requests : null
              }
            }

            dynamic "security_context" {
              for_each = init_containers.value.security_context != null ? [init_containers.value.security_context] : []
              content {
                run_as_user                = security_context.value.run_as_user
                run_as_group               = security_context.value.run_as_group
                run_as_non_root            = security_context.value.run_as_non_root
                allow_privilege_escalation = security_context.value.allow_privilege_escalation
                privileged                 = security_context.value.privileged
                read_only_root_file_system = security_context.value.read_only_root_file_system
              }
            }

            dynamic "volume_mounts" {
              for_each = init_containers.value.volume_mounts
              content {
                name       = volume_mounts.value.name
                mount_path = volume_mounts.value.mount_path
                read_only  = volume_mounts.value.read_only
              }
            }
          }
        }

        # The provider's block is singular-named; the spec carries plain
        # secret names.
        dynamic "image_pull_secret" {
          for_each = eks_properties.value.image_pull_secret_names
          content {
            name = image_pull_secret.value
          }
        }

        dynamic "metadata" {
          for_each = length(eks_properties.value.pod_labels) > 0 ? [1] : []
          content {
            labels = eks_properties.value.pod_labels
          }
        }

        dynamic "volumes" {
          for_each = eks_properties.value.volumes
          content {
            name = volumes.value.name

            dynamic "empty_dir" {
              for_each = volumes.value.empty_dir != null ? [volumes.value.empty_dir] : []
              content {
                # Unset medium means node-backed storage (AWS default "").
                medium     = empty_dir.value.medium != "" ? empty_dir.value.medium : null
                size_limit = empty_dir.value.size_limit
              }
            }

            dynamic "host_path" {
              for_each = volumes.value.host_path != "" ? [volumes.value.host_path] : []
              content {
                path = host_path.value
              }
            }

            dynamic "secret" {
              for_each = volumes.value.secret != null ? [volumes.value.secret] : []
              content {
                secret_name = secret.value.secret_name
                optional    = secret.value.optional
              }
            }
          }
        }
      }
    }
  }

  platform_capabilities = length(var.spec.platform_capabilities) > 0 ? var.spec.platform_capabilities : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  scheduling_priority   = var.spec.scheduling_priority > 0 ? var.spec.scheduling_priority : null
  propagate_tags        = var.spec.propagate_tags
  # Platform-defaulted to true; an explicit false keeps every historical
  # revision ACTIVE for out-of-band consumers.
  deregister_on_new_revision = var.spec.deregister_on_new_revision

  dynamic "retry_strategy" {
    for_each = var.spec.retry_strategy != null ? [var.spec.retry_strategy] : []
    content {
      attempts = retry_strategy.value.attempts > 0 ? retry_strategy.value.attempts : null

      dynamic "evaluate_on_exit" {
        for_each = retry_strategy.value.evaluate_on_exit
        content {
          action           = evaluate_on_exit.value.action
          on_exit_code     = evaluate_on_exit.value.on_exit_code != "" ? evaluate_on_exit.value.on_exit_code : null
          on_reason        = evaluate_on_exit.value.on_reason != "" ? evaluate_on_exit.value.on_reason : null
          on_status_reason = evaluate_on_exit.value.on_status_reason != "" ? evaluate_on_exit.value.on_status_reason : null
        }
      }
    }
  }

  dynamic "timeout" {
    for_each = var.spec.timeout != null ? [var.spec.timeout] : []
    content {
      attempt_duration_seconds = timeout.value.attempt_duration_seconds > 0 ? timeout.value.attempt_duration_seconds : null
    }
  }

  tags = local.aws_tags
}
