# App Runner observability configuration (version).
#
# AWS versions this resource, but the provider's update path is TAGS-ONLY
# and does not replace on trace changes: adding, removing, or changing the
# trace block on an existing configuration plans an in-place update that
# silently changes nothing at AWS (an upstream provider gap -- it briefly
# made trace settings ForceNew in 2023 and reverted the same day). To
# actually change tracing posture, register a NEW configuration (a new
# metadata.name) and repoint referencing services. With X-Ray as the only
# vendor AWS supports, the block's practical states are present or absent
# at creation time.
resource "aws_apprunner_observability_configuration" "this" {
  observability_configuration_name = local.resource_name

  # The trace block is emitted only when the spec configures tracing; a
  # configuration without it is valid but inert. X-Ray is the only vendor
  # App Runner supports today -- applications must emit spans through the
  # AWS Distro for OpenTelemetry SDK for the collector to forward anything.
  dynamic "trace_configuration" {
    for_each = var.spec.trace_configuration != null ? [var.spec.trace_configuration] : []
    content {
      vendor = trace_configuration.value.vendor
    }
  }

  tags = local.aws_tags
}
