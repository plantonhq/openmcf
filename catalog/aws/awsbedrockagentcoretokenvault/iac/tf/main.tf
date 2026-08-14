# AgentCore token vault encryption for one region -- a settings
# singleton: AWS provisions ONE default vault per account+region and
# this module sets which KMS key encrypts it. metadata.name never
# reaches AWS.
#
# Lifecycle facts the render below depends on:
#   - an unset token_vault_id targets AWS's one default vault (the
#     module sends "default", mirroring the provider's own default);
#     changing the id replaces the configuration onto another vault;
#   - key_type/kms_key_arn pairing is CEL-enforced upstream of this
#     module (CustomerManagedKey requires the ARN, ServiceManagedKey
#     forbids it), so the conditional below never sends a stray ARN;
#   - destroy is a NO-OP at AWS: the last-applied key setting REMAINS
#     in effect. Reverting to AWS-managed encryption is an APPLY with
#     key_type ServiceManagedKey, never a destroy.

resource "aws_bedrockagentcore_token_vault_cmk" "this" {
  token_vault_id = var.spec.token_vault_id != "" ? var.spec.token_vault_id : "default"

  kms_configuration {
    key_type    = var.spec.key_type
    kms_key_arn = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null
  }
}
