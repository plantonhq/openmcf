##############################################
# locals.tf
#
# Computed local values for the
# KubernetesGatewayApiCrds module.
##############################################

locals {
  # Gateway API version
  version = coalesce(var.spec.version, "v1.2.1")

  # Determine if experimental channel is requested
  is_experimental = try(var.spec.install_channel.channel, "standard") == "experimental"

  # Channel name for outputs
  channel_name = local.is_experimental ? "experimental" : "standard"

  # Manifest filename based on channel
  manifest_file = local.is_experimental ? "experimental-install.yaml" : "standard-install.yaml"

  # Full URL to download CRD manifests
  manifest_url = "https://github.com/kubernetes-sigs/gateway-api/releases/download/${local.version}/${local.manifest_file}"

  # Resource labels
  labels = {
    "app.kubernetes.io/name"       = "gateway-api-crds"
    "app.kubernetes.io/instance"   = var.metadata.name
    "app.kubernetes.io/managed-by" = "openmcf"
    "app.kubernetes.io/component"  = "crds"
    "gateway-api/version"          = local.version
    "gateway-api/channel"          = local.channel_name
  }

  # Standard CRDs (always installed)
  standard_crds = [
    "gatewayclasses.gateway.networking.k8s.io",
    "gateways.gateway.networking.k8s.io",
    "httproutes.gateway.networking.k8s.io",
    "referencegrants.gateway.networking.k8s.io",
  ]

  # Experimental CRDs (only with experimental channel)
  experimental_crds = [
    "tcproutes.gateway.networking.k8s.io",
    "udproutes.gateway.networking.k8s.io",
    "tlsroutes.gateway.networking.k8s.io",
    "grpcroutes.gateway.networking.k8s.io",
  ]

  # All installed CRDs based on channel
  installed_crds = local.is_experimental ? concat(local.standard_crds, local.experimental_crds) : local.standard_crds
}
