locals {
  # The custom model name is metadata.name -- the naming basis both engines
  # share. AWS allows 1-63 characters.
  custom_model_name = var.metadata.name

  # The customization job name defaults to metadata.name. Job names are
  # unique per account FOREVER (AWS never reuses them, even after delete),
  # so re-running a customization needs an explicit spec.job_name.
  job_name = var.spec.job_name != "" ? var.spec.job_name : var.metadata.name

  # Resource-identity tags match the Pulumi module key-for-key.
  aws_tags = {
    "Name"                     = var.metadata.name
    "planton.ai/resource"      = "true"
    "planton.ai/organization"  = var.metadata.org
    "planton.ai/environment"   = var.metadata.env
    "planton.ai/resource-kind" = "AwsBedrockCustomModel"
    "planton.ai/resource-id"   = var.metadata.id
  }
}
