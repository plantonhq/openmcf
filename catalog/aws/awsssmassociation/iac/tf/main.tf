# A State Manager association: the binding of an SSM document to
# targets on a schedule.
#
# Lifecycle facts the render below depends on:
#   - the document reference (name) forces replacement; every other
#     change creates a new association version in place, and the
#     provider sends the FULL argument set on update (AWS versions
#     associations whole);
#   - AWS identifies the association by a generated UUID, not a name;
#   - AWS materializes the document's declared parameter defaults into
#     the parameters map server-side, and
#     wait_for_success_timeout_seconds is a create-time wait never read
#     back - the import map declares both accordingly;
#   - the association name (association_name) is display metadata, not
#     identity.

resource "aws_ssm_association" "this" {
  name = var.spec.document_name

  association_name = var.spec.association_name != "" ? var.spec.association_name : null
  document_version = var.spec.document_version != "" ? var.spec.document_version : null

  parameters = length(var.spec.parameters) > 0 ? var.spec.parameters : null

  dynamic "targets" {
    for_each = var.spec.targets
    content {
      key    = targets.value.key
      values = targets.value.values
    }
  }

  schedule_expression         = var.spec.schedule_expression != "" ? var.spec.schedule_expression : null
  apply_only_at_cron_interval = var.spec.apply_only_at_cron_interval

  compliance_severity = var.spec.compliance_severity != "" ? var.spec.compliance_severity : null
  sync_compliance     = var.spec.sync_compliance != "" ? var.spec.sync_compliance : null

  max_concurrency = var.spec.max_concurrency != "" ? var.spec.max_concurrency : null
  max_errors      = var.spec.max_errors != "" ? var.spec.max_errors : null

  automation_target_parameter_name = var.spec.automation_target_parameter_name != "" ? var.spec.automation_target_parameter_name : null

  calendar_names = length(var.spec.calendar_names) > 0 ? var.spec.calendar_names : null

  dynamic "output_location" {
    for_each = var.spec.output_location != null ? [var.spec.output_location] : []
    content {
      s3_bucket_name = output_location.value.s3_bucket_name
      s3_key_prefix  = output_location.value.s3_key_prefix != "" ? output_location.value.s3_key_prefix : null
      s3_region      = output_location.value.s3_region != "" ? output_location.value.s3_region : null
    }
  }

  wait_for_success_timeout_seconds = var.spec.wait_for_success_timeout_seconds != 0 ? var.spec.wait_for_success_timeout_seconds : null

  tags = local.aws_tags
}
