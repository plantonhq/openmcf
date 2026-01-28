##############################################
# main.tf
#
# Main orchestration file for KubernetesGhaRunnerScaleSetController
# deployment using Terraform.
#
# This module deploys the GitHub Actions Runner Scale Set Controller
# using the official Helm chart.
#
# Infrastructure Components:
#  1. Kubernetes Namespace (if create_namespace is true)
#  2. Helm Release for the controller
#
# The controller then manages:
#  - AutoScalingRunnerSet CRDs
#  - AutoScalingListener CRDs
#  - EphemeralRunner CRDs
#  - EphemeralRunnerSet CRDs
#
# After deployment, users can create AutoScalingRunnerSet resources
# to deploy actual GitHub Actions runners.
##############################################

# Create namespace if requested
resource "kubernetes_namespace" "controller" {
  count = var.create_namespace ? 1 : 0

  metadata {
    name   = var.namespace
    labels = local.labels
  }
}

# Deploy the controller via Helm
# For OCI charts, the full URL must be passed as the chart parameter
# (repository doesn't work with OCI registries in Terraform helm_release)
resource "helm_release" "controller" {
  name             = local.release_name
  namespace        = var.namespace
  create_namespace = false # We handle namespace creation ourselves

  chart   = local.chart_oci
  version = var.helm_chart_version

  values = [yamlencode(local.helm_values_final)]

  depends_on = [kubernetes_namespace.controller]
}

