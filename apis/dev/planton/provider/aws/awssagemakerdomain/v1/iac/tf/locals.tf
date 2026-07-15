locals {
  # The domain's cloud name is the resource's metadata.name. AWS constrains it
  # to 1-63 characters of [0-9A-Za-z-] and makes it create-time-immutable
  # (changing the name replaces the domain), which is why the name is not spec
  # surface. Same basis as the Pulumi module.
  domain_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  # With spec.tag_propagation = "ENABLED", AWS copies these tags onto the
  # apps, spaces, and user profiles created inside the domain.
  tags = {
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsSagemakerDomain"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Block-presence switches. The spec enables each optional feature by the
  # presence of its message (or, for flattened single-field wrappers, by the
  # presence of the scalar), so the module's job is translating presence into
  # the provider's nested blocks -- never an always-emitted block with zero
  # values, which would pin AWS defaults and create phantom diffs.
  dus = var.spec.default_user_settings

  # domain_settings is one provider block carrying five independent
  # domain-administration dials; the spec models them as top-level fields
  # (everything in a Domain spec is a "domain setting" -- the wrapper adds no
  # information for manifest authors). Emit the block when any dial is set.
  has_domain_settings = (
    length(var.spec.domain_security_group_ids) > 0 ||
    var.spec.docker_settings != null ||
    var.spec.execution_role_identity_config != null ||
    var.spec.r_studio_server_pro_domain_settings != null ||
    var.spec.trusted_identity_propagation_status != null
  )
}
