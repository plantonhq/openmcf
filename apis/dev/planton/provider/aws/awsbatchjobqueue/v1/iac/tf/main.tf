# AWS Batch job queue.
#
# The compute-environment mapping is the queue's core: the scheduler tries
# the lowest `order` first, which is how Spot-first-overflow-to-On-Demand
# (and blue/green environment replacement) is expressed. Deleting a queue is
# disable-then-delete; the provider handles the drain wait.
resource "aws_batch_job_queue" "this" {
  # The cloud name comes from metadata.name (the catalog naming basis) --
  # set explicitly so both engines create the same queue name.
  name     = var.metadata.name
  priority = var.spec.priority
  state    = var.spec.state

  dynamic "compute_environment_order" {
    for_each = var.spec.compute_environment_order
    content {
      order               = compute_environment_order.value.order
      compute_environment = compute_environment_order.value.compute_environment
    }
  }

  # AWS quirk carried by the spec comment: once set, the scheduling policy
  # can be replaced but never removed from a live queue.
  scheduling_policy_arn = var.spec.scheduling_policy != "" ? var.spec.scheduling_policy : null

  dynamic "job_state_time_limit_action" {
    for_each = var.spec.job_state_time_limit_actions
    content {
      action           = job_state_time_limit_action.value.action
      max_time_seconds = job_state_time_limit_action.value.max_time_seconds
      reason           = job_state_time_limit_action.value.reason
      state            = job_state_time_limit_action.value.state
    }
  }

  tags = local.aws_tags
}
