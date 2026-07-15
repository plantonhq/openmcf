locals {
  # The API's cloud name is the resource's metadata.name -- the same basis
  # the Pulumi module uses.
  api_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  aws_tags = {
    "Name"                     = local.api_name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsHttpApiGateway"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Stage defaults: name "$default" when unset, and auto-deploy ON when the
  # spec does not say otherwise -- a declarative spec should be self-applying,
  # so only an explicit auto_deploy=false turns it off. The Pulumi module
  # implements the identical presence rule.
  stage_config = var.spec.stage
  stage_name   = local.stage_config != null && try(local.stage_config.name, "") != "" ? local.stage_config.name : "$default"
  auto_deploy  = local.stage_config != null ? coalesce(local.stage_config.auto_deploy, true) : true

  # Integration deduplication: routes whose integration blocks are IDENTICAL
  # across every field share one API Gateway integration resource; any
  # difference (even a single request-parameter mapping) yields a separate
  # integration. Keying on a hash of the whole object keeps the dedup honest
  # as the integration surface grows -- a partial key (type:uri:payload) would
  # silently merge integrations that differ in the newer fields. The Pulumi
  # module dedups with the same whole-object rule.
  integration_keys = {
    for idx, route in var.spec.routes : tostring(idx) => md5(jsonencode(route.integration))
  }
  integration_map = {
    for idx, route in var.spec.routes : local.integration_keys[tostring(idx)] => route.integration...
  }

  # Authorizer map: authorizers are addressed by their unique name (validated
  # in the spec), which is also how routes bind to them.
  authorizer_map = {
    for auth in var.spec.authorizers : auth.name => auth
  }
}
