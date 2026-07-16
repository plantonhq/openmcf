# Enable the Dataproc API — the control plane that owns the policy.
# disable_on_destroy is false: tearing down one policy must never disable
# the API for everything else in the project.
resource "google_project_service" "dataproc_api" {
  project = local.project_id
  service = "dataproc.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# The Dataproc autoscaling policy. A shareable resource: one policy can
# govern many clusters (each attaches it by reference), so scaling
# behavior is tuned in one place. Policy contents are mutable — updating
# re-tunes every attached cluster — but the API refuses to delete a
# policy while any cluster references it.
resource "google_dataproc_autoscaling_policy" "policy" {
  policy_id = local.policy_id
  location  = local.location
  project   = local.project_id

  # Primary workers are the stable base (they carry HDFS DataNodes);
  # min_instances 0 accepts the API's default of 2.
  worker_config {
    max_instances = var.spec.worker_config.max_instances
    min_instances = var.spec.worker_config.min_instances > 0 ? var.spec.worker_config.min_instances : null
    weight        = var.spec.worker_config.weight > 0 ? var.spec.worker_config.weight : null
  }

  # Secondary workers are the cost-optimized burst arm — no HDFS data,
  # so the autoscaler can add and remove them aggressively.
  dynamic "secondary_worker_config" {
    for_each = var.spec.secondary_worker_config != null ? [var.spec.secondary_worker_config] : []
    content {
      max_instances = secondary_worker_config.value.max_instances > 0 ? secondary_worker_config.value.max_instances : null
      min_instances = secondary_worker_config.value.min_instances > 0 ? secondary_worker_config.value.min_instances : null
      weight        = secondary_worker_config.value.weight > 0 ? secondary_worker_config.value.weight : null
    }
  }

  # The YARN memory-based algorithm. The scale factors express what
  # fraction of pending/available memory the autoscaler acts on per
  # cooldown period — 1.0 is maximally aggressive, small values smooth
  # scaling at the cost of reaction speed.
  basic_algorithm {
    cooldown_period = var.spec.basic_algorithm.cooldown_period != "" ? var.spec.basic_algorithm.cooldown_period : null

    yarn_config {
      graceful_decommission_timeout  = var.spec.basic_algorithm.yarn_config.graceful_decommission_timeout
      scale_up_factor                = var.spec.basic_algorithm.yarn_config.scale_up_factor
      scale_down_factor              = var.spec.basic_algorithm.yarn_config.scale_down_factor
      scale_up_min_worker_fraction   = var.spec.basic_algorithm.yarn_config.scale_up_min_worker_fraction > 0 ? var.spec.basic_algorithm.yarn_config.scale_up_min_worker_fraction : null
      scale_down_min_worker_fraction = var.spec.basic_algorithm.yarn_config.scale_down_min_worker_fraction > 0 ? var.spec.basic_algorithm.yarn_config.scale_down_min_worker_fraction : null
    }
  }

  depends_on = [
    google_project_service.dataproc_api,
  ]
}
