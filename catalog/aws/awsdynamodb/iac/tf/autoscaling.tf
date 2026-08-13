# Folded capacity management: Application Auto Scaling owns the table's
# live read/write capacity on EVERY provisioned table, in one of two
# modes.
#
#   - spec.autoscaling set: the user's min/max bounds with target-tracking
#     policies (and optional scheduled adjustments).
#   - spec.autoscaling unset: pinned targets (min = max = the declared
#     provisioned_throughput values). A declared capacity change updates
#     the target, and AWS moves out-of-range capacity into the new bounds
#     by contract (live-verified 2026-08-13: re-registering a pinned
#     write target 5/5 -> 8/8 moved DescribeTable's WriteCapacityUnits to
#     8 within ~15 seconds, no table operation involved) -- so capacity
#     stays fully declarative even though the table resource ignores it
#     (table.tf).
#
# One target per dimension either way, which is what makes adding or
# removing the autoscaling block an in-place update -- the table resource
# never changes shape. The scaler's identity IS this table (one scalable
# target per table dimension), so it lives here rather than as its own
# kind. PAY_PER_REQUEST tables scale natively and get none of this.

locals {
  table_is_provisioned = var.spec.billing_mode != "PAY_PER_REQUEST"
  manage_capacity      = local.table_is_provisioned && var.spec.provisioned_throughput != null

  # Per-dimension bounds: the user's autoscaling bounds, or the pinned
  # declared capacity.
  read_min = var.spec.autoscaling != null && var.spec.autoscaling.read != null ? var.spec.autoscaling.read.min_capacity : try(var.spec.provisioned_throughput.read_capacity_units, 1)
  read_max = var.spec.autoscaling != null && var.spec.autoscaling.read != null ? var.spec.autoscaling.read.max_capacity : try(var.spec.provisioned_throughput.read_capacity_units, 1)

  write_min = var.spec.autoscaling != null && var.spec.autoscaling.write != null ? var.spec.autoscaling.write.min_capacity : try(var.spec.provisioned_throughput.write_capacity_units, 1)
  write_max = var.spec.autoscaling != null && var.spec.autoscaling.write != null ? var.spec.autoscaling.write.max_capacity : try(var.spec.provisioned_throughput.write_capacity_units, 1)
}

resource "aws_appautoscaling_target" "read" {
  count = local.manage_capacity ? 1 : 0

  service_namespace  = "dynamodb"
  resource_id        = "table/${aws_dynamodb_table.this.name}"
  scalable_dimension = "dynamodb:table:ReadCapacityUnits"
  min_capacity       = local.read_min
  max_capacity       = local.read_max
}

resource "aws_appautoscaling_target" "write" {
  count = local.manage_capacity ? 1 : 0

  service_namespace  = "dynamodb"
  resource_id        = "table/${aws_dynamodb_table.this.name}"
  scalable_dimension = "dynamodb:table:WriteCapacityUnits"
  min_capacity       = local.write_min
  max_capacity       = local.write_max
}

# Target tracking holds consumed-to-provisioned utilization near the
# target percentage -- only rendered when the user configured real
# autoscaling for the dimension (pinned targets need no policy: min = max
# leaves the scaler nothing to decide).
resource "aws_appautoscaling_policy" "read" {
  count = var.spec.autoscaling != null && var.spec.autoscaling.read != null ? 1 : 0

  name               = "${aws_dynamodb_table.this.name}-read-utilization"
  policy_type        = "TargetTrackingScaling"
  service_namespace  = aws_appautoscaling_target.read[0].service_namespace
  resource_id        = aws_appautoscaling_target.read[0].resource_id
  scalable_dimension = aws_appautoscaling_target.read[0].scalable_dimension

  target_tracking_scaling_policy_configuration {
    target_value = var.spec.autoscaling.read.target_utilization_percent

    predefined_metric_specification {
      predefined_metric_type = "DynamoDBReadCapacityUtilization"
    }

    # 0 keeps the AWS default cooldown.
    scale_in_cooldown  = var.spec.autoscaling.read.scale_in_cooldown_seconds > 0 ? var.spec.autoscaling.read.scale_in_cooldown_seconds : null
    scale_out_cooldown = var.spec.autoscaling.read.scale_out_cooldown_seconds > 0 ? var.spec.autoscaling.read.scale_out_cooldown_seconds : null
  }
}

resource "aws_appautoscaling_policy" "write" {
  count = var.spec.autoscaling != null && var.spec.autoscaling.write != null ? 1 : 0

  name               = "${aws_dynamodb_table.this.name}-write-utilization"
  policy_type        = "TargetTrackingScaling"
  service_namespace  = aws_appautoscaling_target.write[0].service_namespace
  resource_id        = aws_appautoscaling_target.write[0].resource_id
  scalable_dimension = aws_appautoscaling_target.write[0].scalable_dimension

  target_tracking_scaling_policy_configuration {
    target_value = var.spec.autoscaling.write.target_utilization_percent

    predefined_metric_specification {
      predefined_metric_type = "DynamoDBWriteCapacityUtilization"
    }

    scale_in_cooldown  = var.spec.autoscaling.write.scale_in_cooldown_seconds > 0 ? var.spec.autoscaling.write.scale_in_cooldown_seconds : null
    scale_out_cooldown = var.spec.autoscaling.write.scale_out_cooldown_seconds > 0 ? var.spec.autoscaling.write.scale_out_cooldown_seconds : null
  }
}

# Scheduled capacity adjustments, keyed by name so entries come and go
# independently. Each targets its dimension's registered scalable target
# (CEL guarantees the dimension's autoscaling config exists).
resource "aws_appautoscaling_scheduled_action" "this" {
  for_each = var.spec.autoscaling != null ? { for adjustment in var.spec.autoscaling.scheduled_adjustments : adjustment.name => adjustment } : {}

  name               = each.value.name
  service_namespace  = "dynamodb"
  resource_id        = "table/${aws_dynamodb_table.this.name}"
  scalable_dimension = each.value.dimension == "READ" ? "dynamodb:table:ReadCapacityUnits" : "dynamodb:table:WriteCapacityUnits"
  schedule           = each.value.schedule
  timezone           = each.value.timezone != "" ? each.value.timezone : null
  start_time         = each.value.start_time != "" ? each.value.start_time : null
  end_time           = each.value.end_time != "" ? each.value.end_time : null

  scalable_target_action {
    # 0 leaves the bound unchanged when the schedule fires (CEL requires
    # at least one).
    min_capacity = each.value.min_capacity > 0 ? each.value.min_capacity : null
    max_capacity = each.value.max_capacity > 0 ? each.value.max_capacity : null
  }

  depends_on = [
    aws_appautoscaling_target.read,
    aws_appautoscaling_target.write,
  ]
}
