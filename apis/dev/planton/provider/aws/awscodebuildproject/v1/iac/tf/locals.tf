locals {
  # The project's cloud name is the resource's metadata.name. CodeBuild
  # project names are create-time-immutable (changing the name replaces the
  # project), which is why the name is not spec surface. Same basis as the
  # Pulumi module.
  project_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  tags = {
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsCodeBuildProject"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # Lambda compute ignores build/queued timeouts (AWS caps Lambda builds
  # itself, and the provider diff-suppresses the arguments for these types).
  # The spec's CEL already rejects EXPLICIT timeouts on Lambda environments;
  # this guard additionally keeps the spec-level defaults (60/480) from being
  # sent, so the plan matches what AWS actually stores.
  is_lambda_env = contains(["LINUX_LAMBDA_CONTAINER", "ARM_LAMBDA_CONTAINER"], var.spec.environment.type)

  # Block-presence switches. The spec enables each optional feature by the
  # presence of its message; the module translates presence into the
  # provider's nested blocks.
  has_webhook    = var.spec.webhook != null
  has_vpc_config = var.spec.vpc_config != null
  # An absent cache block and an explicit NO_CACHE deploy identically; the
  # cache type is a presence-carrying optional, so null means the NO_CACHE
  # default.
  has_cache           = var.spec.cache != null && var.spec.cache.type != null && var.spec.cache.type != "NO_CACHE"
  has_logs_config     = var.spec.logs_config != null
  has_batch_config    = var.spec.build_batch_config != null
  has_resource_policy = var.spec.resource_policy != null

  # The resource policy arrives as a structured JSON document (the tfvars
  # layer passes nested objects, never pre-serialized strings).
  resource_policy_json = local.has_resource_policy ? jsonencode(var.spec.resource_policy) : null
}
