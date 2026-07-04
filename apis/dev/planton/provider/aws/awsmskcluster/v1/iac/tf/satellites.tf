# Cluster-scoped settings AWS models as standalone resources but that are
# honestly part of the cluster's own configuration -- each is keyed by the
# cluster ARN, owned by exactly one cluster, and referenced by nothing else.

# SASL/SCRAM credential associations: one association per Secrets Manager
# secret, materialized per-ARN so adding or removing a secret updates in
# place instead of churning a batch list.
resource "aws_msk_single_scram_secret_association" "this" {
  for_each = toset(var.spec.scram_secret_arns)

  cluster_arn = aws_msk_cluster.this.arn
  secret_arn  = each.value
}

# A resource-based IAM policy on the cluster -- the grant behind cross-account
# PrivateLink access (kafka:CreateVpcConnection and friends).
resource "aws_msk_cluster_policy" "this" {
  count = var.spec.cluster_policy != "" ? 1 : 0

  cluster_arn = aws_msk_cluster.this.arn
  policy      = var.spec.cluster_policy
}
