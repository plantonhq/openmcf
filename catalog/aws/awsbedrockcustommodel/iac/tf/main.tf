# Amazon Bedrock custom model: starts a model-customization job
# (fine-tuning, continued pre-training, or distillation) that produces a
# custom model from a base foundation model and your training data.
#
# Every argument is create-time-immutable -- a customization job cannot be
# altered once started, so any change replaces the job and model. The job
# runs asynchronously in AWS after the deploy returns; job_status reports
# where it stood at read time.

resource "aws_bedrock_custom_model" "this" {
  custom_model_name = local.custom_model_name
  job_name          = local.job_name

  # The base foundation model being customized (foundation-model ARN).
  base_model_identifier = var.spec.base_model_arn

  # Sent only when set: the provider attribute is Optional+Computed and AWS
  # defaults to FINE_TUNING.
  customization_type = var.spec.customization_type != "" ? var.spec.customization_type : null

  # Per-base-model training knobs (epochCount, batchSize, learningRate, ...).
  hyperparameters = var.spec.hyperparameters

  # The role Bedrock assumes to read training data and write outputs. Must
  # trust bedrock.amazonaws.com.
  role_arn = var.spec.role_arn

  # Customer-managed key for the resulting model when referenced;
  # Bedrock-managed key otherwise.
  custom_model_kms_key_id = var.spec.custom_model_kms_key_arn != "" ? var.spec.custom_model_kms_key_arn : null

  training_data_config {
    s3_uri = var.spec.training_data_s3_uri
  }

  output_data_config {
    s3_uri = var.spec.output_data_s3_uri
  }

  # Up to 10 validation datasets -- Bedrock reports per-dataset metrics on
  # the finished job.
  dynamic "validation_data_config" {
    for_each = length(var.spec.validation_data_s3_uris) > 0 ? [true] : []
    content {
      dynamic "validator" {
        for_each = var.spec.validation_data_s3_uris
        content {
          s3_uri = validator.value
        }
      }
    }
  }

  # Run the job's data access inside the caller's VPC (both members
  # required together -- CEL-enforced).
  dynamic "vpc_config" {
    for_each = var.spec.vpc_config != null ? [var.spec.vpc_config] : []
    content {
      subnet_ids         = vpc_config.value.subnet_ids
      security_group_ids = vpc_config.value.security_group_ids
    }
  }

  tags = local.aws_tags
}
