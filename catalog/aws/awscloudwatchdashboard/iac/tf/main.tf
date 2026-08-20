# One CloudWatch dashboard: a named widget canvas.
#
# Lifecycle facts the render below depends on:
#   - the dashboard's name is spec.dashboard_name (an explicit field -
#     dashboard names carry uppercase metadata.name cannot), and
#     changing it replaces the dashboard;
#   - create and update are the same AWS call (PutDashboard is a pure
#     upsert), so every body change applies in place;
#   - dashboards are untaggable at AWS - no tags argument exists on the
#     resource.

resource "aws_cloudwatch_dashboard" "this" {
  dashboard_name = var.spec.dashboard_name
  dashboard_body = local.dashboard_body
}
