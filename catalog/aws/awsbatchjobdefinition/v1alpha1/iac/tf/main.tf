# AWS Batch job definition.
#
# Every meaningful change registers a NEW revision (revisions are immutable
# in AWS); with deregister_on_new_revision at its default the previous
# revision is deregistered so exactly one ACTIVE revision tracks this
# resource. Only single-container ECS-based jobs are modeled; multinode
# (nodeProperties), multi-container ECS (ecsProperties), and Batch-on-EKS
# (eksProperties) are separate long-tail arms.
resource "aws_batch_job_definition" "this" {
  # The cloud name comes from metadata.name (the catalog naming basis) --
  # revisions register under this name in both engines.
  name = var.metadata.name
  type = "container"

  container_properties = jsonencode(local.container_properties)

  platform_capabilities = length(var.spec.platform_capabilities) > 0 ? var.spec.platform_capabilities : null
  parameters            = length(var.spec.parameters) > 0 ? var.spec.parameters : null
  scheduling_priority   = var.spec.scheduling_priority > 0 ? var.spec.scheduling_priority : null
  propagate_tags        = var.spec.propagate_tags
  # Platform-defaulted to true; an explicit false keeps every historical
  # revision ACTIVE for out-of-band consumers.
  deregister_on_new_revision = var.spec.deregister_on_new_revision

  dynamic "retry_strategy" {
    for_each = var.spec.retry_strategy != null ? [var.spec.retry_strategy] : []
    content {
      attempts = retry_strategy.value.attempts > 0 ? retry_strategy.value.attempts : null

      dynamic "evaluate_on_exit" {
        for_each = retry_strategy.value.evaluate_on_exit
        content {
          action           = evaluate_on_exit.value.action
          on_exit_code     = evaluate_on_exit.value.on_exit_code != "" ? evaluate_on_exit.value.on_exit_code : null
          on_reason        = evaluate_on_exit.value.on_reason != "" ? evaluate_on_exit.value.on_reason : null
          on_status_reason = evaluate_on_exit.value.on_status_reason != "" ? evaluate_on_exit.value.on_status_reason : null
        }
      }
    }
  }

  dynamic "timeout" {
    for_each = var.spec.timeout != null ? [var.spec.timeout] : []
    content {
      attempt_duration_seconds = timeout.value.attempt_duration_seconds > 0 ? timeout.value.attempt_duration_seconds : null
    }
  }

  tags = local.aws_tags
}
