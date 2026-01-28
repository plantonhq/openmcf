###########################
# main.tf
###########################

# Conditional namespace creation
resource "kubernetes_namespace" "external_dns" {
  count = try(var.spec.create_namespace, false) ? 1 : 0

  metadata {
    name   = local.namespace
    labels = local.final_labels
  }
}

# Data source for existing namespace
data "kubernetes_namespace" "existing" {
  count = try(var.spec.create_namespace, false) ? 0 : 1

  metadata {
    name = local.namespace
  }
}

# Create service account with cloud provider annotations
resource "kubernetes_service_account" "external_dns" {
  metadata {
    name        = local.ksa_name
    namespace   = local.namespace_name
    annotations = local.sa_annotations
    labels      = local.final_labels
  }
}

# Create secret for Cloudflare API token (only if using Cloudflare)
resource "kubernetes_secret" "cloudflare_api_token" {
  count = local.is_cloudflare ? 1 : 0

  metadata {
    name      = local.cloudflare_api_token_secret_name
    namespace = local.namespace_name
    labels    = local.final_labels
  }

  data = {
    apiKey = local.cf_api_token
  }

  type = "Opaque"
}

# Deploy ExternalDNS via Helm
resource "helm_release" "external_dns" {
  name       = local.release_name
  namespace  = local.namespace_name
  repository = local.helm_repo_url
  chart      = local.helm_chart_name
  version    = local.helm_chart_version

  atomic          = true
  cleanup_on_fail = true
  wait            = true
  wait_for_jobs   = true
  timeout         = 180

  # Service account configuration
  set {
    name  = "serviceAccount.create"
    value = "false"
  }

  set {
    name  = "serviceAccount.name"
    value = local.ksa_name
  }

  # ExternalDNS version
  set {
    name  = "image.tag"
    value = local.external_dns_version
  }

  # Provider configuration
  set {
    name  = "provider"
    value = local.provider_type
  }

  # GKE-specific configuration
  dynamic "set" {
    for_each = local.is_gke ? [1] : []
    content {
      name  = "google.project"
      value = local.gke_project_id
    }
  }

  dynamic "set" {
    for_each = local.is_gke ? [1] : []
    content {
      name  = "zoneIdFilters[0]"
      value = local.gke_dns_zone_id
    }
  }

  # EKS-specific configuration
  dynamic "set" {
    for_each = local.is_eks ? [1] : []
    content {
      name  = "zoneIdFilters[0]"
      value = local.eks_route53_zone_id
    }
  }

  # AKS-specific configuration
  dynamic "set" {
    for_each = local.is_aks && local.aks_dns_zone_id != "" ? [1] : []
    content {
      name  = "zoneIdFilters[0]"
      value = local.aks_dns_zone_id
    }
  }

  # Cloudflare-specific configuration
  dynamic "set" {
    for_each = local.is_cloudflare ? [1] : []
    content {
      name  = "sources[0]"
      value = "service"
    }
  }

  dynamic "set" {
    for_each = local.is_cloudflare ? [1] : []
    content {
      name  = "sources[1]"
      value = "ingress"
    }
  }

  dynamic "set" {
    for_each = local.is_cloudflare ? [1] : []
    content {
      name  = "sources[2]"
      value = "gateway-httproute"
    }
  }

  # Istio Gateway source - enables ExternalDNS to watch networking.istio.io/v1 Gateway resources
  dynamic "set" {
    for_each = local.is_cloudflare ? [1] : []
    content {
      name  = "sources[3]"
      value = "istio-gateway"
    }
  }

  dynamic "set" {
    for_each = local.is_cloudflare ? [1] : []
    content {
      name  = "env[0].name"
      value = "CF_API_TOKEN"
    }
  }

  dynamic "set" {
    for_each = local.is_cloudflare ? [1] : []
    content {
      name  = "env[0].valueFrom.secretKeyRef.name"
      value = local.cloudflare_api_token_secret_name
    }
  }

  dynamic "set" {
    for_each = local.is_cloudflare ? [1] : []
    content {
      name  = "env[0].valueFrom.secretKeyRef.key"
      value = "apiKey"
    }
  }

  dynamic "set" {
    for_each = local.is_cloudflare ? [1] : []
    content {
      name  = "extraArgs[0]"
      value = "--cloudflare-dns-records-per-page=5000"
    }
  }

  dynamic "set" {
    for_each = local.is_cloudflare ? [1] : []
    content {
      name  = "extraArgs[1]"
      value = "--zone-id-filter=${local.cf_dns_zone_id}"
    }
  }

  dynamic "set" {
    for_each = local.is_cloudflare && local.cf_is_proxied ? [1] : []
    content {
      name  = "extraArgs[2]"
      value = "--cloudflare-proxied"
    }
  }

  depends_on = [
    kubernetes_namespace.external_dns,
    data.kubernetes_namespace.existing,
    kubernetes_service_account.external_dns
  ]
}

