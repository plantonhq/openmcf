# Amazon Bedrock guardrail: content-safety policy families (content
# filters, denied topics, word filters, sensitive-information handling,
# contextual grounding) evaluated on model inputs and outputs, plus
# immutable published versions for production pinning.
#
# The module owns the guardrail's mutable DRAFT definition and one
# aws_bedrock_guardrail_version per spec.versions entry.

# Deploy-time identity, consumed only when the spec carries the portable
# geography-qualified guardrail-profile id ("us.guardrail.v1:0"): the AWS
# API accepts that id directly, but the provider's schema types
# guardrail_profile_identifier as an ARN (fwtypes.ARNType -- stricter
# than the API, live-caught 2026-08-13), so the module composes the
# account-scoped profile ARN here rather than forcing every manifest to
# embed its account id.
data "aws_caller_identity" "current" {
  count = local.compose_cross_region_arn ? 1 : 0
}

data "aws_partition" "current" {
  count = local.compose_cross_region_arn ? 1 : 0
}

resource "aws_bedrock_guardrail" "this" {
  # Create-time naming basis; doubles as the Name tag. metadata.name on
  # both engines.
  name = local.guardrail_name

  # Description is sent only when set: the provider attribute is
  # Optional+Computed, so sending "" would fight AWS's normalization.
  description = var.spec.description != "" ? var.spec.description : null

  # Required by AWS for every guardrail: what the caller sees when the
  # guardrail intervenes on input/output.
  blocked_input_messaging   = var.spec.blocked_input_messaging
  blocked_outputs_messaging = var.spec.blocked_outputs_messaging

  # Customer-managed key when referenced; Bedrock-managed key otherwise.
  kms_key_arn = var.spec.kms_key_arn != "" ? var.spec.kms_key_arn : null

  # -------------------------------------------------------------------
  # Content filters
  # -------------------------------------------------------------------
  dynamic "content_policy_config" {
    for_each = local.has_content_policy ? [var.spec.content_policy] : []
    content {
      # The safeguard tier rides an optional single-entry list attribute.
      # Omitted entirely when unset -- the provider treats the absent list
      # as "AWS default" (CLASSIC) and pins whatever AWS materializes.
      tier_config = content_policy_config.value.tier != "" ? [{ tier_name = content_policy_config.value.tier }] : null

      dynamic "filters_config" {
        for_each = content_policy_config.value.filters
        content {
          type            = filters_config.value.type
          input_strength  = filters_config.value.input_strength
          output_strength = filters_config.value.output_strength
          # Action/enabled arms are send-when-set: AWS defaults actions to
          # BLOCK and enabled to true; explicit values (including false)
          # are always transmitted so disablement is expressible.
          input_action      = filters_config.value.input_action != "" ? filters_config.value.input_action : null
          output_action     = filters_config.value.output_action != "" ? filters_config.value.output_action : null
          input_enabled     = filters_config.value.input_enabled
          output_enabled    = filters_config.value.output_enabled
          input_modalities  = length(filters_config.value.input_modalities) > 0 ? filters_config.value.input_modalities : null
          output_modalities = length(filters_config.value.output_modalities) > 0 ? filters_config.value.output_modalities : null
        }
      }
    }
  }

  # -------------------------------------------------------------------
  # Denied topics
  # -------------------------------------------------------------------
  dynamic "topic_policy_config" {
    for_each = local.has_topic_policy ? [var.spec.topic_policy] : []
    content {
      tier_config = topic_policy_config.value.tier != "" ? [{ tier_name = topic_policy_config.value.tier }] : null

      dynamic "topics_config" {
        for_each = topic_policy_config.value.topics
        content {
          name       = topics_config.value.name
          definition = topics_config.value.definition
          examples   = length(topics_config.value.examples) > 0 ? topics_config.value.examples : null
          # DENY is the only topic type AWS defines -- the modules own the
          # constant so the spec never asks for a one-value field.
          type = "DENY"
        }
      }
    }
  }

  # -------------------------------------------------------------------
  # Word filters
  # -------------------------------------------------------------------
  dynamic "word_policy_config" {
    for_each = local.has_word_policy ? [var.spec.word_policy] : []
    content {
      # The AWS-managed profanity list -- PROFANITY is the only managed
      # list type AWS defines; presence of spec.profanity_filter enables
      # it and the modules own the type constant.
      dynamic "managed_word_lists_config" {
        for_each = word_policy_config.value.profanity_filter != null ? [word_policy_config.value.profanity_filter] : []
        content {
          type           = "PROFANITY"
          input_action   = managed_word_lists_config.value.input_action != "" ? managed_word_lists_config.value.input_action : null
          output_action  = managed_word_lists_config.value.output_action != "" ? managed_word_lists_config.value.output_action : null
          input_enabled  = managed_word_lists_config.value.input_enabled
          output_enabled = managed_word_lists_config.value.output_enabled
        }
      }

      dynamic "words_config" {
        for_each = word_policy_config.value.custom_words
        content {
          text           = words_config.value.text
          input_action   = words_config.value.input_action != "" ? words_config.value.input_action : null
          output_action  = words_config.value.output_action != "" ? words_config.value.output_action : null
          input_enabled  = words_config.value.input_enabled
          output_enabled = words_config.value.output_enabled
        }
      }
    }
  }

  # -------------------------------------------------------------------
  # Sensitive information (PII + regexes)
  # -------------------------------------------------------------------
  dynamic "sensitive_information_policy_config" {
    for_each = local.has_sensitive_information ? [var.spec.sensitive_information_policy] : []
    content {
      dynamic "pii_entities_config" {
        for_each = sensitive_information_policy_config.value.pii_entities
        content {
          type   = pii_entities_config.value.type
          action = pii_entities_config.value.action
          # Per-direction overrides are Optional+Computed at the provider:
          # AWS materializes them from `action` when omitted, and once set
          # they never revert to AWS-derived (taught on the spec fields).
          input_action   = pii_entities_config.value.input_action != "" ? pii_entities_config.value.input_action : null
          output_action  = pii_entities_config.value.output_action != "" ? pii_entities_config.value.output_action : null
          input_enabled  = pii_entities_config.value.input_enabled
          output_enabled = pii_entities_config.value.output_enabled
        }
      }

      dynamic "regexes_config" {
        for_each = sensitive_information_policy_config.value.regexes
        content {
          name           = regexes_config.value.name
          pattern        = regexes_config.value.pattern
          description    = regexes_config.value.description != "" ? regexes_config.value.description : null
          action         = regexes_config.value.action
          input_action   = regexes_config.value.input_action != "" ? regexes_config.value.input_action : null
          output_action  = regexes_config.value.output_action != "" ? regexes_config.value.output_action : null
          input_enabled  = regexes_config.value.input_enabled
          output_enabled = regexes_config.value.output_enabled
        }
      }
    }
  }

  # -------------------------------------------------------------------
  # Contextual grounding
  # -------------------------------------------------------------------
  dynamic "contextual_grounding_policy_config" {
    for_each = local.has_contextual_grounding_policy ? [var.spec.contextual_grounding_policy] : []
    content {
      dynamic "filters_config" {
        for_each = contextual_grounding_policy_config.value.filters
        content {
          type      = filters_config.value.type
          threshold = filters_config.value.threshold
        }
      }
    }
  }

  # -------------------------------------------------------------------
  # Cross-region inference profile. The spec accepts the portable
  # geography-qualified profile id ("us.guardrail.v1:0") or the full
  # profile ARN; the provider's schema demands an ARN (stricter than the
  # AWS API, which resolves the plain id -- live-verified 2026-08-13), so
  # the id shape is composed into the caller's account-scoped ARN
  # (local.cross_region_identifier). Required by AWS whenever any policy
  # family uses the STANDARD tier (CEL-enforced at manifest time).
  # -------------------------------------------------------------------
  dynamic "cross_region_config" {
    for_each = local.has_cross_region ? [local.cross_region_identifier] : []
    content {
      guardrail_profile_identifier = cross_region_config.value
    }
  }

  tags = local.aws_tags
}

# Published versions -- one immutable numbered version per spec.versions
# entry, keyed by the entry's stable name (the logical name "published"
# also disambiguates the import map's version placeholder from the
# guardrail's own DRAFT).
resource "aws_bedrock_guardrail_version" "published" {
  for_each = local.versions

  guardrail_arn = aws_bedrock_guardrail.this.guardrail_arn
  description   = each.value.description != "" ? each.value.description : null

  # Keep the published version in AWS when the entry (or the whole
  # guardrail) is removed from management.
  skip_destroy = each.value.keep_on_delete

  # The draft must carry the intended definition before publishing.
  depends_on = [aws_bedrock_guardrail.this]
}
