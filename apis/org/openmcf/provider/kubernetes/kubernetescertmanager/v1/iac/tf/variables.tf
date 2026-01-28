variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for Kubernetes Cert-Manager deployment"
  type = object({
    # Target Kubernetes cluster
    target_cluster_name = string

    # Kubernetes namespace where cert-manager will be deployed
    namespace = optional(string, "kubernetes-cert-manager")

    # Flag to indicate if the namespace should be created
    create_namespace = optional(bool, false)

    # cert-manager version (e.g., "v1.19.1")
    kubernetes_cert_manager_version = optional(string, "v1.19.1")

    # Helm chart version
    helm_chart_version = optional(string, "v1.19.1")

    # Skip installation of self-signed issuer
    skip_install_self_signed_issuer = optional(bool, false)

    # Global ACME configuration
    acme = object({
      # ACME account email
      email = string

      # ACME server URL (defaults to Let's Encrypt production)
      server = optional(string, "https://acme-v02.api.letsencrypt.org/directory")
    })

    # List of DNS provider configurations
    dns_providers = list(object({
      # Unique name for this provider
      name = string

      # DNS zones this provider manages
      dns_zones = list(string)

      # GCP Cloud DNS provider (optional)
      gcp_cloud_dns = optional(object({
        project_id            = string
        service_account_email = string
      }))

      # AWS Route53 provider (optional)
      aws_route53 = optional(object({
        region   = string
        role_arn = string
      }))

      # Azure DNS provider (optional)
      azure_dns = optional(object({
        subscription_id = string
        resource_group  = string
        client_id       = string
      }))

      # Cloudflare provider (optional)
      cloudflare = optional(object({
        api_token = string
      }))
    }))
  })

  validation {
    condition     = length(var.spec.dns_providers) > 0
    error_message = "At least one DNS provider must be configured"
  }

  validation {
    condition     = can(regex("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$", var.spec.acme.email))
    error_message = "ACME email must be a valid email address"
  }
}
