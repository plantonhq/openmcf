locals {
  # The pipeline's cloud name is the resource's metadata.name. Pipeline names
  # are create-time-immutable (changing the name replaces the pipeline),
  # which is why the name is not spec surface. Same basis as the Pulumi
  # module.
  pipeline_name = var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key. Identity
  # tagging is the only tagging surface this module manages; user-defined
  # custom tags are a platform-wide concern, not per-kind spec surface.
  tags = {
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsCodePipeline"
    "planton.ai/resource-id"   = var.metadata.id
  }

  # AWS models the artifact store two ways: exactly one store without a
  # region (single-region pipeline) or one store per region (cross-region
  # pipeline). The spec's CEL enforces the shape; the module only needs to
  # know which arm it is rendering.
  is_single_region = length(var.spec.artifact_stores) == 1 && var.spec.artifact_stores[0].region == ""

  # The spec defaults to V2/SUPERSEDED, while the PROVIDER defaults
  # pipeline_type to V1 -- an omitted value must deploy the same pipeline an
  # explicit V2 would, so the spec defaults are applied here, never left to
  # the provider.
  pipeline_type  = var.spec.pipeline_type != null && var.spec.pipeline_type != "" ? var.spec.pipeline_type : "V2"
  execution_mode = var.spec.execution_mode != null && var.spec.execution_mode != "" ? var.spec.execution_mode : "SUPERSEDED"
}
