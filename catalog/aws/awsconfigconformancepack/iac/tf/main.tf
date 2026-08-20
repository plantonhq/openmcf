# An AWS Config conformance pack at one of two scopes: this account
# (aws_config_conformance_pack) or the whole organization
# (aws_config_organization_conformance_pack). One spec, one scope
# flag, exactly one of the two resources rendered.
#
# Lifecycle facts the renders below depend on:
#   - deploying a pack REQUIRES an active Config recorder in the
#     region (AWS rejects it otherwise) - the recorder is the
#     consumer's contract (AwsConfigRecorder), never this module's;
#   - AWS never reports the template back, so template drift is
#     undetectable by design (both provider resources document it);
#   - the pack service-linked role creates the pack's rules; deleting
#     the pack deletes them (org packs unwind from member accounts,
#     which can take minutes - the provider waits);
#   - conformance packs carry NO tags at the provider (neither
#     resource has a tags argument).

resource "aws_config_conformance_pack" "this" {
  count = var.spec.organization_scope ? 0 : 1

  # metadata.name is the pack name on both engines (letters, digits,
  # hyphens, starting with a letter; account packs cap at 256 chars).
  name = var.metadata.name

  delivery_s3_bucket     = var.spec.delivery_s3_bucket != "" ? var.spec.delivery_s3_bucket : null
  delivery_s3_key_prefix = var.spec.delivery_s3_key_prefix != "" ? var.spec.delivery_s3_key_prefix : null

  # Account packs accept both template forms at once (AWS prefers the
  # S3 one); the spec CEL guarantees at least one arrives here.
  template_body   = var.spec.template_body != "" ? var.spec.template_body : null
  template_s3_uri = var.spec.template_s3_uri != "" ? var.spec.template_s3_uri : null

  dynamic "input_parameter" {
    for_each = var.spec.input_parameters
    content {
      parameter_name  = input_parameter.value.parameter_name
      parameter_value = input_parameter.value.parameter_value
    }
  }
}

resource "aws_config_organization_conformance_pack" "this" {
  count = var.spec.organization_scope ? 1 : 0

  # Organization packs cap the name at 128 chars; the delivery bucket
  # must begin with "awsconfigconforms" (AWS naming contracts enforced
  # at deploy).
  name = var.metadata.name

  delivery_s3_bucket     = var.spec.delivery_s3_bucket != "" ? var.spec.delivery_s3_bucket : null
  delivery_s3_key_prefix = var.spec.delivery_s3_key_prefix != "" ? var.spec.delivery_s3_key_prefix : null

  # Organization packs accept exactly one template form (the spec CEL
  # guarantees it).
  template_body   = var.spec.template_body != "" ? var.spec.template_body : null
  template_s3_uri = var.spec.template_s3_uri != "" ? var.spec.template_s3_uri : null

  excluded_accounts = length(var.spec.excluded_accounts) > 0 ? var.spec.excluded_accounts : null

  dynamic "input_parameter" {
    for_each = var.spec.input_parameters
    content {
      parameter_name  = input_parameter.value.parameter_name
      parameter_value = input_parameter.value.parameter_value
    }
  }
}
