# An SSM Parameter Store entry: one named configuration value
# (String, StringList, or SecureString).
#
# Lifecycle facts the render below depends on:
#   - the parameter's name is spec.parameter_name (an explicit field -
#     hierarchical names carry slashes metadata.name cannot), and
#     changing it forces replacement;
#   - the spec's plain `value` arm renders as insecure_value (readable
#     in plans - the arm's whole point) and the secret `secure_value`
#     arm renders as the provider's sensitive value argument; the spec
#     guarantees exactly one arm is set;
#   - overwrite renders only when TRUE: the provider's unset behavior
#     (fail on a pre-existing foreign name at create, overwrite own
#     updates) is the safe default, and an explicit false would break
#     the provider's own update path;
#   - an Advanced -> Standard tier downgrade forces replacement (AWS
#     forbids it in place), and Intelligent-Tiering is resolved
#     server-side to Standard or Advanced per write.

resource "aws_ssm_parameter" "this" {
  name = var.spec.parameter_name
  type = var.spec.type

  value          = var.spec.secure_value != "" ? var.spec.secure_value : null
  insecure_value = var.spec.value != "" ? var.spec.value : null

  description     = var.spec.description != "" ? var.spec.description : null
  allowed_pattern = var.spec.allowed_pattern != "" ? var.spec.allowed_pattern : null
  tier            = var.spec.tier != "" ? var.spec.tier : null
  key_id          = var.spec.key_id != "" ? var.spec.key_id : null
  data_type       = var.spec.data_type != "" ? var.spec.data_type : null

  overwrite = var.spec.overwrite ? true : null

  tags = local.aws_tags
}
