# Amazon SageMaker AI image with its folded versions.
#
# Lifecycle facts the renders below depend on:
#   - the provider sleeps ~1 minute before CreateImage (IAM propagation)
#     - every create is at least a minute;
#   - version numbers are AWS-assigned, monotonic, never reused: a
#     changed base_image REPLACES the version under a NEW number (the
#     old number stays retired). Entries are keyed by position -
#     append-only, taught on the spec;
#   - image_version carries NO tags upstream (by provider design);
#   - AWS serializes version creation per image (the provider holds a
#     mutex) - versions attach one at a time.

resource "aws_sagemaker_image" "this" {
  # The component's name IS the image name.
  image_name = local.image_name

  role_arn = var.spec.role_arn

  display_name = var.spec.display_name != "" ? var.spec.display_name : null
  description  = var.spec.description != "" ? var.spec.description : null

  tags = local.aws_tags
}

resource "aws_sagemaker_image_version" "this" {
  for_each = local.versions

  image_name = aws_sagemaker_image.this.image_name

  # The version's identity - changing it replaces the version under a
  # new AWS-assigned number.
  base_image = each.value.base_image

  aliases = length(each.value.aliases) > 0 ? each.value.aliases : null

  horovod          = each.value.horovod
  job_type         = each.value.job_type != "" ? each.value.job_type : null
  ml_framework     = each.value.ml_framework != "" ? each.value.ml_framework : null
  processor        = each.value.processor != "" ? each.value.processor : null
  programming_lang = each.value.programming_lang != "" ? each.value.programming_lang : null
  release_notes    = each.value.release_notes != "" ? each.value.release_notes : null
  vendor_guidance  = each.value.vendor_guidance != "" ? each.value.vendor_guidance : null
}
