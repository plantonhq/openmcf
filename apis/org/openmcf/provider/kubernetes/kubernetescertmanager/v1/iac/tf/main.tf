# Kubernetes Cert-Manager Terraform Module

# Conditionally create namespace for cert-manager
resource "kubernetes_namespace" "cert_manager" {
  count = var.spec.create_namespace ? 1 : 0
  
  metadata {
    name = local.namespace
  }
}

# Look up existing namespace when not creating
data "kubernetes_namespace" "cert_manager" {
  count = var.spec.create_namespace ? 0 : 1
  
  metadata {
    name = local.namespace
  }
}

# Create ServiceAccount with workload identity annotations
resource "kubernetes_service_account" "cert_manager" {
  metadata {
    name        = local.ksa_name
    namespace   = local.namespace_name
    annotations = local.sa_annotations
  }
}

# Deploy cert-manager Helm chart
resource "helm_release" "cert_manager" {
  name       = local.helm_chart_name
  repository = local.helm_chart_repo
  chart      = local.helm_chart_name
  version    = local.helm_chart_version
  namespace  = local.namespace_name

  # Wait for resources to be ready
  wait          = true
  wait_for_jobs = true
  timeout       = 180

  # Enable atomic rollback on failure
  atomic          = true
  cleanup_on_fail = true

  # Helm values
  values = [yamlencode({
    installCRDs = true
    serviceAccount = {
      create = false
      name   = local.ksa_name
    }
    extraArgs = [
      "--dns01-recursive-nameservers-only",
      "--dns01-recursive-nameservers=1.1.1.1:53,8.8.8.8:53"
    ]
  })]

  depends_on = [kubernetes_service_account.cert_manager]
}

# Create Kubernetes Secrets for Cloudflare providers
# Secret names use metadata.name prefix for uniqueness when multiple instances share a namespace
resource "kubernetes_secret" "cloudflare" {
  for_each = { for provider in local.cloudflare_providers : provider.name => provider }

  metadata {
    name      = "${var.metadata.name}-${each.key}-credentials"
    namespace = local.namespace_name
  }

  data = {
    "api-token" = each.value.cloudflare.api_token
  }
}

# Create ClusterIssuer resources (one per domain)
resource "kubernetes_manifest" "cluster_issuer" {
  for_each = { for issuer in local.cluster_issuers : issuer.domain => issuer }

  manifest = {
    apiVersion = "cert-manager.io/v1"
    kind       = "ClusterIssuer"

    metadata = {
      name = each.value.domain
    }

    spec = {
      acme = {
        email  = each.value.acme_email
        server = each.value.acme_server

        privateKeySecretRef = {
          name = "letsencrypt-${each.value.domain}-account-key"
        }

        solvers = [
          # GCP Cloud DNS solver
          each.value.gcp_cloud_dns != null ? {
            dns01 = {
              cloudDNS = {
                project = each.value.gcp_cloud_dns.project_id
              }
            }
          } :
          # AWS Route53 solver
          each.value.aws_route53 != null ? {
            dns01 = {
              route53 = {
                region = each.value.aws_route53.region
              }
            }
          } :
          # Azure DNS solver
          each.value.azure_dns != null ? {
            dns01 = {
              azureDNS = {
                subscriptionID    = each.value.azure_dns.subscription_id
                resourceGroupName = each.value.azure_dns.resource_group
              }
            }
          } :
          # Cloudflare solver
          # Uses computed secret name with metadata.name prefix for uniqueness
          each.value.cloudflare != null ? {
            dns01 = {
              cloudflare = {
                apiTokenSecretRef = {
                  name = "${var.metadata.name}-${each.value.provider_name}-credentials"
                  key  = "api-token"
                }
              }
            }
          } : null
        ]
      }
    }
  }

  depends_on = [
    helm_release.cert_manager,
    kubernetes_secret.cloudflare
  ]
}

