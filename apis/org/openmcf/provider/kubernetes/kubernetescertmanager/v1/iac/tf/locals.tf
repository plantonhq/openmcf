# Local values for Cert-Manager deployment

locals {
  # Namespace (use from spec or default)
  namespace = var.spec.namespace

  # Namespace reference (either created or looked up)
  namespace_name = var.spec.create_namespace ? kubernetes_namespace.cert_manager[0].metadata[0].name : data.kubernetes_namespace.cert_manager[0].metadata[0].name

  # Helm chart constants
  helm_chart_name    = "cert-manager"
  helm_chart_repo    = "https://charts.jetstack.io"
  helm_chart_version = var.spec.helm_chart_version

  # Computed resource names to avoid conflicts when multiple instances share a namespace
  # Format: {metadata.name}-{purpose}
  # Users can prefix metadata.name with component type if needed (e.g., "cert-manager-prod")
  # ServiceAccount name uses metadata.name for uniqueness
  ksa_name = var.metadata.name

  # Build ServiceAccount annotations for workload identity
  # Multiple providers may configure different annotations
  sa_annotations = merge(
    { for provider in var.spec.dns_providers :
      "iam.gke.io/gcp-service-account" => provider.gcp_cloud_dns.service_account_email
      if provider.gcp_cloud_dns != null
    },
    { for provider in var.spec.dns_providers :
      "eks.amazonaws.com/role-arn" => provider.aws_route53.role_arn
      if provider.aws_route53 != null
    },
    { for provider in var.spec.dns_providers :
      "azure.workload.identity/client-id" => provider.azure_dns.client_id
      if provider.azure_dns != null
    }
  )

  # Extract Cloudflare providers for secret creation
  cloudflare_providers = [
    for provider in var.spec.dns_providers :
    provider if provider.cloudflare != null
  ]

  # Build ClusterIssuer list (one per domain across all providers)
  cluster_issuers = flatten([
    for provider in var.spec.dns_providers : [
      for zone in provider.dns_zones : {
        domain        = zone
        provider_name = provider.name
        acme_email    = var.spec.acme.email
        acme_server   = var.spec.acme.server

        # Provider-specific configuration
        gcp_cloud_dns = provider.gcp_cloud_dns
        aws_route53   = provider.aws_route53
        azure_dns     = provider.azure_dns
        cloudflare    = provider.cloudflare
      }
    ]
  ])

  # Outputs for stack_outputs
  release_name = "cert-manager"

  # Solver identity (first non-null identity found)
  solver_identity = try(
    [for provider in var.spec.dns_providers :
      provider.gcp_cloud_dns.service_account_email
      if provider.gcp_cloud_dns != null
    ][0],
    try(
      [for provider in var.spec.dns_providers :
        provider.aws_route53.role_arn
        if provider.aws_route53 != null
      ][0],
      try(
        [for provider in var.spec.dns_providers :
          provider.azure_dns.client_id
          if provider.azure_dns != null
        ][0],
        ""
      )
    )
  )
}

locals {
  # Cloudflare secret name (first cloudflare provider's secret)
  # Uses metadata.name prefix for uniqueness when multiple instances share a namespace
  cloudflare_secret_name = length(local.cloudflare_providers) > 0 ? "${var.metadata.name}-${local.cloudflare_providers[0].name}-credentials" : ""
}

