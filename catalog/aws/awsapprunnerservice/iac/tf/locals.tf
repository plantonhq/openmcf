locals {
  # The service's cloud name is the resource's metadata.name -- the same
  # basis the Pulumi module uses. ForceNew: renaming replaces the service.
  service_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  aws_tags = {
    "Name"                     = local.service_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsAppRunnerService"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Egress is VPC-routed exactly when a connector is referenced; there is no
  # module-created connector -- AwsAppRunnerVpcConnector is its own
  # composable resource shared across services.
  egress_type = var.spec.vpc_connector_arn != "" ? "VPC" : "DEFAULT"

  # The authentication block is needed for private ECR pulls (access role)
  # or any code repository (connection). Generator-flattened ref fields are
  # never null (contract default ""), so plain != "" comparisons are safe.
  needs_auth_config = (
    (var.spec.image_source != null && var.spec.image_source.access_role_arn != "") ||
    var.spec.code_source != null
  )

  # Custom domain associations keyed by domain name so entries add/remove
  # independently (a keyed set, not an ordered list -- reordering the spec
  # never touches AWS).
  custom_domains = { for d in var.spec.custom_domains : d.domain_name => d }

  # VPC Ingress Connections keyed by connection name -- the same keyed-set
  # semantics as custom_domains.
  vpc_ingress_connections = { for c in var.spec.vpc_ingress_connections : c.name => c }
}
