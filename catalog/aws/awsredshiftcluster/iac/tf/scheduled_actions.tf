# Scheduled actions pause, resume, or resize this cluster on a cron/at
# schedule -- the nights-and-weekends cost lever. Each entry renders one
# aws_redshift_scheduled_action keyed by its (account-unique) name; the
# spec's exactly-one-arm CEL guarantees a single target arm per entry.
#
# The IAM role's trust policy must allow scheduler.redshift.amazonaws.com
# to sts:AssumeRole -- AWS validates the trust at create, and the
# provider retries briefly while a freshly created trust propagates.
resource "aws_redshift_scheduled_action" "this" {
  for_each = { for a in var.spec.scheduled_actions : a.name => a }

  name        = each.value.name
  description = each.value.description != "" ? each.value.description : null

  # The spec's zero value (disabled=false) is AWS's default: enabled.
  enable = !each.value.disabled

  schedule   = each.value.schedule
  start_time = each.value.start_time != "" ? each.value.start_time : null
  end_time   = each.value.end_time != "" ? each.value.end_time : null

  iam_role = each.value.iam_role_arn

  target_action {
    dynamic "pause_cluster" {
      for_each = each.value.pause_cluster ? [1] : []
      content {
        cluster_identifier = aws_redshift_cluster.this.cluster_identifier
      }
    }

    dynamic "resume_cluster" {
      for_each = each.value.resume_cluster ? [1] : []
      content {
        cluster_identifier = aws_redshift_cluster.this.cluster_identifier
      }
    }

    dynamic "resize_cluster" {
      for_each = each.value.resize_cluster != null ? [each.value.resize_cluster] : []
      content {
        cluster_identifier = aws_redshift_cluster.this.cluster_identifier
        classic            = resize_cluster.value.classic

        # Unset members keep the cluster's current topology.
        cluster_type    = resize_cluster.value.cluster_type != "" ? resize_cluster.value.cluster_type : null
        node_type       = resize_cluster.value.node_type != "" ? resize_cluster.value.node_type : null
        number_of_nodes = resize_cluster.value.number_of_nodes != 0 ? resize_cluster.value.number_of_nodes : null
      }
    }
  }
}
