##############################################
# locals.tf
#
# Computed values and transformations for the
# KubernetesGhaRunnerScaleSetController module.
##############################################

locals {
  # Release configuration - use resource name for consistency with other components
  release_name = var.metadata.name
  # For OCI charts, the full URL must be passed as the chart parameter
  # (repository doesn't work with OCI registries in Terraform helm_release)
  chart_oci = "oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller"

  # Standard labels
  labels = merge(
    {
      "openmcf.org/resource"      = "true"
      "openmcf.org/resource-name" = var.metadata.name
      "openmcf.org/resource-kind" = "KubernetesGhaRunnerScaleSetController"
    },
    var.metadata.id != "" ? { "openmcf.org/resource-id" = var.metadata.id } : {},
    var.metadata.org != "" ? { "openmcf.org/organization" = var.metadata.org } : {},
    var.metadata.env != "" ? { "openmcf.org/environment" = var.metadata.env } : {}
  )

  # Metrics enabled check
  metrics_enabled = var.metrics != null && var.metrics.controller_manager_addr != ""

  # Build Helm values
  helm_values = {
    replicaCount = var.replica_count
    labels       = local.labels

    resources = {
      requests = {
        cpu    = try(var.container.resources.requests.cpu, "100m")
        memory = try(var.container.resources.requests.memory, "128Mi")
      }
      limits = {
        cpu    = try(var.container.resources.limits.cpu, "500m")
        memory = try(var.container.resources.limits.memory, "512Mi")
      }
    }

    flags = merge(
      {
        logLevel                      = var.flags.log_level
        logFormat                     = var.flags.log_format
        runnerMaxConcurrentReconciles = var.flags.runner_max_concurrent_reconciles
        updateStrategy                = var.flags.update_strategy
      },
      var.flags.watch_single_namespace != "" ? { watchSingleNamespace = var.flags.watch_single_namespace } : {},
      length(var.flags.exclude_label_propagation_prefixes) > 0 ? { excludeLabelPropagationPrefixes = var.flags.exclude_label_propagation_prefixes } : {},
      var.flags.k8s_client_rate_limiter_qps > 0 ? { k8sClientRateLimiterQPS = var.flags.k8s_client_rate_limiter_qps } : {},
      var.flags.k8s_client_rate_limiter_burst > 0 ? { k8sClientRateLimiterBurst = var.flags.k8s_client_rate_limiter_burst } : {}
    )

    priorityClassName = var.priority_class_name

    imagePullSecrets = [for s in var.image_pull_secrets : { name = s }]
  }

  # Add image configuration if provided
  helm_values_with_image = var.container.image.repository != "" ? merge(local.helm_values, {
    image = {
      repository = var.container.image.repository
      tag        = var.container.image.tag
      pullPolicy = var.container.image.pull_policy
    }
  }) : local.helm_values

  # Add metrics configuration if enabled
  helm_values_final = local.metrics_enabled ? merge(local.helm_values_with_image, {
    metrics = {
      controllerManagerAddr = var.metrics.controller_manager_addr
      listenerAddr          = var.metrics.listener_addr
      listenerEndpoint      = var.metrics.listener_endpoint
    }
  }) : local.helm_values_with_image
}

