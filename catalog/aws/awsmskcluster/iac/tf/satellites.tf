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
# PrivateLink access (kafka:CreateVpcConnection and friends). The spec carries
# the policy as a structured document; it is serialized to JSON here.
resource "aws_msk_cluster_policy" "this" {
  count = var.spec.cluster_policy != null ? 1 : 0

  cluster_arn = aws_msk_cluster.this.arn
  policy      = jsonencode(var.spec.cluster_policy)
}

# Declared Kafka topics, managed through the MSK topic API -- no Kafka client
# or bootstrap connectivity needed. Keyed by topic name so adding or removing
# one topic never churns its neighbors. Topic deletion requires
# delete.topic.enable=true on the cluster (MSK's default).
resource "aws_msk_topic" "this" {
  for_each = { for t in var.spec.topics : t.name => t }

  cluster_arn        = aws_msk_cluster.this.arn
  name               = each.value.name
  partition_count    = each.value.partition_count
  replication_factor = each.value.replication_factor

  # The provider takes topic configs as a JSON document; entries with no
  # overrides omit the argument entirely.
  configs = length(each.value.configs) > 0 ? jsonencode(each.value.configs) : null
}
