# Enable the Cloud Run Admin API before creating the job so a fresh project
# works first try. disable_on_destroy=false: turning an API off on teardown
# is a project-wide blast radius no single resource should own.
resource "google_project_service" "run_api" {
  project = local.project_id
  service = "run.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Cloud Run job: a run-to-completion workload definition. Each manual or
# scheduled trigger creates an execution that stamps out task_count tasks with
# up to parallelism running concurrently. This resource owns the definition;
# individual executions are separate API objects.
resource "google_cloud_run_v2_job" "main" {
  name         = local.job_name
  location     = var.spec.region
  project      = local.project_id
  launch_stage = local.launch_stage

  labels      = local.final_labels
  annotations = length(var.spec.annotations) > 0 ? var.spec.annotations : null

  deletion_protection = var.spec.deletion_protection

  dynamic "binary_authorization" {
    for_each = var.spec.binary_authorization != null ? [var.spec.binary_authorization] : []
    content {
      use_default              = binary_authorization.value.use_default ? true : null
      policy                   = binary_authorization.value.policy != "" ? binary_authorization.value.policy : null
      breakglass_justification = binary_authorization.value.breakglass_justification != "" ? binary_authorization.value.breakglass_justification : null
    }
  }

  template {
    task_count   = var.spec.task_count
    parallelism = var.spec.parallelism

    # The inner template describes one task: containers, volumes, identity,
    # networking, and per-task limits.
    template {
      service_account       = local.service_account
      execution_environment = local.execution_environment
      encryption_key        = local.encryption_key
      timeout               = local.timeout
      max_retries           = var.spec.template.max_retries
      gpu_zonal_redundancy_disabled = var.spec.gpu_zonal_redundancy_disabled ? true : null

      dynamic "node_selector" {
        for_each = var.spec.template.node_selector != null ? [var.spec.template.node_selector] : []
        content {
          accelerator = node_selector.value.accelerator
        }
      }

      dynamic "vpc_access" {
        for_each = local.has_vpc_access ? [1] : []
        content {
          connector = local.vpc_connector
          egress    = local.vpc_egress

          dynamic "network_interfaces" {
            for_each = local.vpc_interfaces
            content {
              network    = network_interfaces.value.network != "" ? network_interfaces.value.network : null
              subnetwork = network_interfaces.value.subnetwork != "" ? network_interfaces.value.subnetwork : null
              tags       = length(network_interfaces.value.tags) > 0 ? network_interfaces.value.tags : null
            }
          }
        }
      }

      dynamic "volumes" {
        for_each = var.spec.template.volumes
        content {
          name = volumes.value.name

          dynamic "cloud_sql_instance" {
            for_each = volumes.value.cloud_sql_instance != null ? [volumes.value.cloud_sql_instance] : []
            content {
              instances = cloud_sql_instance.value.instances
            }
          }

          dynamic "secret" {
            for_each = volumes.value.secret != null ? [volumes.value.secret] : []
            content {
              secret       = secret.value.secret
              default_mode = secret.value.default_mode

              dynamic "items" {
                for_each = secret.value.items
                content {
                  path    = items.value.path
                  version = items.value.version != "" ? items.value.version : null
                  mode    = items.value.mode
                }
              }
            }
          }

          dynamic "empty_dir" {
            for_each = volumes.value.empty_dir != null ? [volumes.value.empty_dir] : []
            content {
              medium     = empty_dir.value.medium != "" ? empty_dir.value.medium : null
              size_limit = empty_dir.value.size_limit != "" ? empty_dir.value.size_limit : null
            }
          }

          dynamic "gcs" {
            for_each = volumes.value.gcs != null ? [volumes.value.gcs] : []
            content {
              bucket    = gcs.value.bucket
              read_only = gcs.value.read_only
            }
          }

          dynamic "nfs" {
            for_each = volumes.value.nfs != null ? [volumes.value.nfs] : []
            content {
              server    = nfs.value.server
              path      = nfs.value.path
              read_only = nfs.value.read_only
            }
          }
        }
      }

      dynamic "containers" {
        for_each = var.spec.template.containers
        content {
          name        = containers.value.name != "" ? containers.value.name : null
          image       = containers.value.image
          command     = length(containers.value.command) > 0 ? containers.value.command : null
          args        = length(containers.value.args) > 0 ? containers.value.args : null
          working_dir = containers.value.working_dir != "" ? containers.value.working_dir : null
          depends_on  = length(containers.value.depends_on) > 0 ? containers.value.depends_on : null

          dynamic "env" {
            for_each = containers.value.env
            content {
              name  = env.value.name
              value = env.value.value_from_secret == null ? env.value.value : null

              dynamic "value_source" {
                for_each = env.value.value_from_secret != null ? [env.value.value_from_secret] : []
                content {
                  secret_key_ref {
                    secret  = value_source.value.secret
                    version = value_source.value.version != "" ? value_source.value.version : null
                  }
                }
              }
            }
          }

          dynamic "resources" {
            for_each = containers.value.resources != null ? [containers.value.resources] : []
            content {
              limits = (resources.value.cpu != "" || resources.value.memory != "") ? merge(
                resources.value.cpu != "" ? { cpu = resources.value.cpu } : {},
                resources.value.memory != "" ? { memory = resources.value.memory } : {},
              ) : null
            }
          }

          dynamic "volume_mounts" {
            for_each = containers.value.volume_mounts
            content {
              name       = volume_mounts.value.name
              mount_path = volume_mounts.value.mount_path
            }
          }
        }
      }
    }
  }

  depends_on = [google_project_service.run_api]
}
