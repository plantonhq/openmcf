# Amazon Bedrock Provisioned Throughput: dedicated model serving capacity
# purchased in model units -- the required serving path for fine-tuned
# custom models.
#
# COST: this resource bills from the moment it is created. Without a
# commitment it bills hourly and deletes any time; with a OneMonth/
# SixMonths commitment it bills for the FULL term and cannot be deleted
# until the term lapses (the provider refuses the destroy server-side).
# Every argument except tags is create-time-immutable.

resource "aws_bedrock_provisioned_model_throughput" "this" {
  provisioned_model_name = local.provisioned_model_name

  # The model capacity is bought for -- a custom model's output ARN (the
  # default reference wiring) or a foundation-model ARN where AWS allows
  # direct provisioning.
  model_arn = var.spec.model_arn

  # A model unit's throughput (tokens/minute) is model-specific; account
  # quotas for NO-COMMITMENT units default low (often 0-2).
  model_units = var.spec.model_units

  # Omitted entirely for no commitment (hourly billing) -- the provider
  # treats the absent argument as the no-commitment purchase.
  commitment_duration = var.spec.commitment_duration != "" ? var.spec.commitment_duration : null

  tags = local.aws_tags
}
