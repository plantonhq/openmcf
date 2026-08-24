# The region's GuardDuty posture: the detector, its protection-plan
# features, finding filters, trusted/threat IP lists, findings export,
# and - for organization administrators - member management.
#
# Lifecycle facts the renders below depend on:
#   - AWS allows ONE detector per account per region; a pre-existing
#     detector (enabled by hand or by Organizations auto-enable) makes
#     create fail with "detector already exists";
#   - detector features and organization features are PATCHES onto the
#     detector (Create and Update are the same UpdateDetector call;
#     Delete is a no-op) - features removed from the spec are NOT
#     reverted by AWS, and the provider serializes feature writes per
#     detector under a global mutex, which is why every feature
#     resource below chains on the detector;
#   - the organization-configuration resource's delete is a no-op too:
#     destroying this component leaves the org posture as last
#     applied (taught in the GUIDE);
#   - the publishing destination's bucket POLICY and KMS key policy
#     must grant guardduty.amazonaws.com before create (the consumer's
#     contract on AwsS3Bucket / AwsKmsKey specs);
#   - changing TAGS on the publishing destination REPLACES it in the
#     provider (TagsSchemaForceNew) - the destination is deliberately
#     untagged here so tag sweeps never replace findings export;
#   - members and the invite accepter are the cross-account surface: a
#     member record needs the member account's root email, and the
#     accepter runs in the MEMBER account against a pending invite.

resource "aws_guardduty_detector" "this" {
  # Rendered only on an explicit choice so the module never fights the
  # provider default (enabled).
  enable = var.spec.enable == null ? true : var.spec.enable

  # Left to AWS (SIX_HOURS) unless the spec sets it: members inherit
  # the administrator's value, and an explicit send on a member
  # detector would fight the org sync forever (the idempotency gate
  # would catch exactly that).
  finding_publishing_frequency = var.spec.finding_publishing_frequency != "" ? var.spec.finding_publishing_frequency : null

  tags = local.aws_tags
}

# Protection plans, patch-keyed by feature name.
resource "aws_guardduty_detector_feature" "this" {
  for_each = local.features

  detector_id = aws_guardduty_detector.this.id
  name        = each.value.name
  status      = each.value.enabled == null ? "ENABLED" : (each.value.enabled ? "ENABLED" : "DISABLED")

  # The full sub-toggle family, always sent (see locals.tf -- AWS
  # materializes undeclared members as DISABLED, and a partial send
  # breaks post-apply plan idempotency).
  dynamic "additional_configuration" {
    for_each = local.feature_additional_configurations[each.key]
    content {
      name   = additional_configuration.value.name
      status = additional_configuration.value.status
    }
  }
}

# Finding filters, keyed by name (the for_each key both engines share).
resource "aws_guardduty_filter" "this" {
  for_each = local.filters

  detector_id = aws_guardduty_detector.this.id
  name        = each.value.name
  description = each.value.description != "" ? each.value.description : null
  action      = each.value.action
  rank        = each.value.rank

  finding_criteria {
    dynamic "criterion" {
      for_each = each.value.criteria
      content {
        field                 = criterion.value.field
        equals                = length(criterion.value.equals) > 0 ? criterion.value.equals : null
        not_equals            = length(criterion.value.not_equals) > 0 ? criterion.value.not_equals : null
        matches               = length(criterion.value.matches) > 0 ? criterion.value.matches : null
        not_matches           = length(criterion.value.not_matches) > 0 ? criterion.value.not_matches : null
        greater_than          = criterion.value.greater_than != "" ? criterion.value.greater_than : null
        greater_than_or_equal = criterion.value.greater_than_or_equal != "" ? criterion.value.greater_than_or_equal : null
        less_than             = criterion.value.less_than != "" ? criterion.value.less_than : null
        less_than_or_equal    = criterion.value.less_than_or_equal != "" ? criterion.value.less_than_or_equal : null
      }
    }
  }

  tags = local.aws_tags
}

# Trusted IP lists (AWS activates at most one per detector).
resource "aws_guardduty_ipset" "this" {
  for_each = local.ip_sets

  detector_id = aws_guardduty_detector.this.id
  name        = each.value.name
  format      = each.value.format
  location    = each.value.location
  activate    = each.value.activate

  tags = local.aws_tags
}

# Threat intel lists.
resource "aws_guardduty_threatintelset" "this" {
  for_each = local.threat_intel_sets

  detector_id = aws_guardduty_detector.this.id
  name        = each.value.name
  format      = each.value.format
  location    = each.value.location
  activate    = each.value.activate

  tags = local.aws_tags
}

# Findings export to S3 (deliberately untagged - see the header).
resource "aws_guardduty_publishing_destination" "this" {
  count = var.spec.publishing_destination != null ? 1 : 0

  detector_id      = aws_guardduty_detector.this.id
  destination_arn  = var.spec.publishing_destination.bucket_arn
  kms_key_arn      = var.spec.publishing_destination.kms_key_arn
  destination_type = "S3"
}

# ----- ADMIN side: organization administration -----

# The account-global delegation act (one per organization, performed
# from the MANAGEMENT account).
resource "aws_guardduty_organization_admin_account" "this" {
  count = var.spec.organization != null && var.spec.organization.admin_account_id != "" ? 1 : 0

  admin_account_id = var.spec.organization.admin_account_id
}

resource "aws_guardduty_organization_configuration" "this" {
  count = var.spec.organization != null ? 1 : 0

  detector_id                      = aws_guardduty_detector.this.id
  auto_enable_organization_members = var.spec.organization.auto_enable_organization_members

  # A same-apply delegation must land before org configuration.
  depends_on = [aws_guardduty_organization_admin_account.this]
}

# Organization-wide feature auto-enablement, patch-keyed by name.
resource "aws_guardduty_organization_configuration_feature" "this" {
  for_each = local.org_features

  detector_id = aws_guardduty_detector.this.id
  name        = each.value.name
  auto_enable = each.value.auto_enable

  dynamic "additional_configuration" {
    for_each = each.value.additional_configuration
    content {
      name        = additional_configuration.value.name
      auto_enable = additional_configuration.value.auto_enable
    }
  }

  depends_on = [aws_guardduty_organization_configuration.this]
}

# ----- ADMIN side: members -----

resource "aws_guardduty_member" "this" {
  for_each = local.members

  detector_id = aws_guardduty_detector.this.id
  account_id  = each.value.account_id
  email       = each.value.email

  invite                     = each.value.invite == null ? null : each.value.invite
  invitation_message         = each.value.invitation_message != "" ? each.value.invitation_message : null
  disable_email_notification = each.value.disable_email_notification
}

# Per-member protection-plan overrides, keyed "account/feature".
resource "aws_guardduty_member_detector_feature" "this" {
  for_each = local.member_features

  detector_id = aws_guardduty_detector.this.id
  account_id  = each.value.account_id
  name        = each.value.feature.name
  status      = each.value.feature.enabled == null ? "ENABLED" : (each.value.feature.enabled ? "ENABLED" : "DISABLED")

  dynamic "additional_configuration" {
    for_each = each.value.feature.additional_configuration
    content {
      name   = additional_configuration.value.name
      status = additional_configuration.value.enabled == null ? "ENABLED" : (additional_configuration.value.enabled ? "ENABLED" : "DISABLED")
    }
  }

  depends_on = [aws_guardduty_member.this]
}

# ----- MEMBER side: accept a pending invitation -----

resource "aws_guardduty_invite_accepter" "this" {
  count = var.spec.accept_invitation_from_account_id != "" ? 1 : 0

  detector_id       = aws_guardduty_detector.this.id
  master_account_id = var.spec.accept_invitation_from_account_id
}
