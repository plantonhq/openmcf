# Amazon SageMaker model package group (the model registry's unit of
# organization) with its folded resource policy.
#
# Lifecycle facts the renders below depend on:
#   - the group is immutable except tags (even the description is
#     ForceNew upstream - a description edit replaces the group);
#   - the policy is an idempotent upsert (PutModelPackageGroupPolicy)
#     that updates in place; removing spec.resource_policy deletes the
#     policy resource, leaving the group open to its own account only;
#   - model package VERSIONS register into the group imperatively
#     (training pipelines, SDK) - never through this module.

resource "aws_sagemaker_model_package_group" "this" {
  # The component's name IS the group name.
  model_package_group_name = local.group_name

  model_package_group_description = var.spec.description != "" ? var.spec.description : null

  tags = local.aws_tags
}

# Cross-account sharing policy - present exactly when the spec carries
# one.
resource "aws_sagemaker_model_package_group_policy" "this" {
  count = local.has_policy ? 1 : 0

  model_package_group_name = aws_sagemaker_model_package_group.this.model_package_group_name

  resource_policy = jsonencode(var.spec.resource_policy)
}
