# Associates the cluster with an EXISTING snapshot schedule (AWS keeps
# exactly one schedule per cluster -- the schedule replaces the default
# automated-snapshot cadence). The schedule itself is an account-scoped
# resource shared by many clusters and is not created here.
resource "aws_redshift_snapshot_schedule_association" "this" {
  count = var.spec.snapshot_schedule_identifier != "" ? 1 : 0

  cluster_identifier  = aws_redshift_cluster.this.cluster_identifier
  schedule_identifier = var.spec.snapshot_schedule_identifier
}
