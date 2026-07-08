# App Runner observability configuration (version).
#
# AWS versions this resource: the trace settings are create-time immutable,
# so any change destroys this revision and registers the NEXT revision under
# the same configuration name. Referencing services pick up the new
# revision-carrying ARN through the resource graph on their next deployment.
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
