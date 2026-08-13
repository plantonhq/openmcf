locals {
  # The pool's cloud name is metadata.name -- the same basis the Pulumi module
  # uses, so both engines create identically-named pools.
  resource_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key (the canonical
  # six-key identity map -- user labels never merge into cloud tags).
  aws_tags = {
    "Name"                     = local.resource_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsCognitoUserPool"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # A certificate is what makes a domain "custom" to AWS (it rides a managed
  # CloudFront distribution); a dot in the domain is the spec-level signal,
  # CEL-coupled to the certificate's presence.
  has_domain       = var.spec.domain != null && try(var.spec.domain.domain, "") != ""
  is_custom_domain = local.has_domain && strcontains(var.spec.domain.domain, ".")
}
