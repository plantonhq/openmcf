# A customer-owned SSM document: a reusable definition of the actions
# Systems Manager performs (Command, Automation, Session, ...).
#
# Lifecycle facts the render below depends on:
#   - the document name is metadata.name on both engines and changing
#     it forces replacement;
#   - updating the content creates a NEW document version and the
#     provider promotes it to the default version; schema-1.x documents
#     only update when the content itself changes (an AWS rule);
#   - permissions is the provider's flat share map - the spec's
#     share_with_account_ids renders as {type: "Share", account_ids:
#     "<comma-joined>"} (AWS applies changes in batches of 20);
#   - attachment metadata is never read back by any SSM API - the
#     import map declares attachment_sources config-only.

resource "aws_ssm_document" "this" {
  name = local.document_name

  content       = var.spec.content
  document_type = var.spec.document_type

  document_format = var.spec.document_format != "" ? var.spec.document_format : null
  target_type     = var.spec.target_type != "" ? var.spec.target_type : null
  version_name    = var.spec.version_name != "" ? var.spec.version_name : null

  dynamic "attachments_source" {
    for_each = var.spec.attachment_sources
    content {
      key    = attachments_source.value.key
      name   = attachments_source.value.name != "" ? attachments_source.value.name : null
      values = attachments_source.value.values
    }
  }

  permissions = length(var.spec.share_with_account_ids) > 0 ? {
    type        = "Share"
    account_ids = join(",", var.spec.share_with_account_ids)
  } : null

  tags = local.aws_tags
}
