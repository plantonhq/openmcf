# The Fargate profile composes onto its neighbors instead of embedding
# them: the cluster attaches by reference, the pod execution role is a
# referenced AwsIamRole that carries its own policies (this module never
# modifies a role it merely references), and the subnets are referenced
# AwsSubnet nodes -- private only, which AWS enforces at create time.
#
# The ENTIRE profile is create-only in AWS: any change to name, cluster,
# role, subnets, or selectors replaces the profile. AWS also serializes
# profile operations per cluster (one create or delete at a time) -- the
# provider simply waits; no ordering configuration is needed here.
resource "aws_eks_fargate_profile" "this" {
  fargate_profile_name   = local.fargate_profile_name
  cluster_name           = var.spec.cluster_name
  pod_execution_role_arn = var.spec.pod_execution_role_arn
  subnet_ids             = var.spec.subnet_ids

  # A pod runs on Fargate when it matches ANY selector; within one
  # selector the namespace AND every label must match.
  dynamic "selector" {
    for_each = var.spec.selectors
    content {
      namespace = selector.value.namespace
      labels    = length(selector.value.labels) > 0 ? selector.value.labels : null
    }
  }

  tags = local.aws_tags
}
