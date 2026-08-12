# Amazon SES (SESv2) email identity.
#
# The identity resource carries the identity string, its default
# configuration set, and DKIM signing; the custom MAIL FROM domain,
# feedback forwarding, and authorization policies are separate AWS
# sub-resources keyed by the identity, materialized below so each updates
# independently of the identity itself.
resource "aws_sesv2_email_identity" "this" {
  # The identity string IS the AWS identifier (the domain or address being
  # verified) -- deliberately from the spec, not metadata.name, because
  # the identity must be the exact DNS name mail is sent from.
  email_identity = var.spec.email_identity

  # The default configuration set for every message sent from this
  # identity; attachable/swappable in place. The ref-flattened string
  # carries the contract default "" when absent -- a plain != ""
  # comparison, never coalesce(x, ""), which errors in HCL when every
  # argument is empty.
  configuration_set_name = var.spec.configuration_set != "" ? var.spec.configuration_set : null

  # DKIM signing. Two arms behind one block (spec-level CELs keep them
  # exclusive): Easy DKIM rotates to next_signing_key_length, BYODKIM
  # brings the key/selector pair. An absent block accepts AWS's Easy DKIM
  # default (RSA_2048_BIT) -- SES always DKIM-enables new identities.
  dynamic "dkim_signing_attributes" {
    for_each = local.dkim != null ? [local.dkim] : []
    content {
      next_signing_key_length    = local.byodkim ? null : (dkim_signing_attributes.value.next_signing_key_length != "" ? dkim_signing_attributes.value.next_signing_key_length : null)
      domain_signing_private_key = local.byodkim ? dkim_signing_attributes.value.domain_signing_private_key : null
      domain_signing_selector    = local.byodkim ? dkim_signing_attributes.value.domain_signing_selector : null
    }
  }

  tags = local.aws_tags
}

# Custom MAIL FROM domain -- a PUT-style satellite on the identity (one
# per identity), created only when the manifest configures it. Deleting
# the block reverts the identity to the amazonses.com envelope sender.
resource "aws_sesv2_email_identity_mail_from_attributes" "this" {
  count = var.spec.mail_from != null ? 1 : 0

  email_identity   = aws_sesv2_email_identity.this.email_identity
  mail_from_domain = var.spec.mail_from.mail_from_domain

  # Contract-defaulted to USE_DEFAULT_VALUE (AWS's own default): a missing
  # MX record degrades alignment but never silently stops mail unless the
  # manifest explicitly opts into REJECT_MESSAGE.
  behavior_on_mx_failure = coalesce(var.spec.mail_from.behavior_on_mx_failure, "USE_DEFAULT_VALUE")
}

# Bounce/complaint email forwarding -- materialized only when the manifest
# takes an explicit position. Absent, the setting is UNMANAGED: a fresh
# identity gets AWS's default (forwarding ON), but SES retains the
# last-written value per identity name across even identity deletion
# (live-verified 2026-08-12), so a previously-managed name keeps its old
# position. Note the provider's DESTROY of this resource writes
# forwarding=false (PutEmailIdentityFeedbackAttributes with the unset/zero
# value), which then persists for the name -- removing an explicit position
# is a one-way write to false until a new explicit true. Set false once
# event destinations or SNS feedback carry bounces.
resource "aws_sesv2_email_identity_feedback_attributes" "this" {
  count = var.spec.email_forwarding_enabled != null ? 1 : 0

  email_identity           = aws_sesv2_email_identity.this.email_identity
  email_forwarding_enabled = var.spec.email_forwarding_enabled
}

# Authorization policies -- one AWS sub-resource per named entry (the
# cross-account sending grants).
resource "aws_sesv2_email_identity_policy" "this" {
  for_each = local.policies

  email_identity = aws_sesv2_email_identity.this.email_identity
  policy_name    = each.key

  # The policy Struct arrives from the tfvars layer as a nested object --
  # encode it to the JSON document the provider expects.
  policy = jsonencode(each.value.policy)
}
