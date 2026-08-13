# Amazon Bedrock AgentCore memory: a managed short-term (raw session
# events) and long-term (strategy-extracted records) store with folded
# strategy satellites.
#
# Lifecycle facts the renders below depend on:
#   - strategy writes go through the parent memory's update API -- AWS
#     serializes them per memory and the provider holds a per-memory
#     lock, so strategy operations can take tens of minutes (the
#     provider's own timeout is 45m);
#   - the deprecated strategy-level memory_execution_role_arn and
#     namespaces arguments are never sent (excluded-deprecated; the
#     memory-level role and namespace_templates are the living
#     surfaces);
#   - MEMORY_RECORDS is the only stream content type AWS defines -- the
#     module owns the constant.

resource "aws_bedrockagentcore_memory" "this" {
  # AWS's memory-name charset (letter first, then letters/digits/_) is
  # stricter than metadata.name conventions, so the name is an explicit
  # spec field. Changing it replaces the memory.
  name = var.spec.memory_name

  # Required by AWS: days raw session events survive (the short-term
  # window).
  event_expiry_duration = var.spec.event_expiry_days

  description = var.spec.description != "" ? var.spec.description : null

  # Changing the key replaces the memory (provider-enforced).
  encryption_key_arn = var.spec.encryption_key_arn != "" ? var.spec.encryption_key_arn : null

  memory_execution_role_arn = var.spec.execution_role_arn != "" ? var.spec.execution_role_arn : null

  # Metadata keys indexed for filtered retrieval (1-10; changing the set
  # replaces the memory, provider-enforced).
  dynamic "indexed_key" {
    for_each = var.spec.indexed_keys
    content {
      key  = indexed_key.value.key
      type = indexed_key.value.type
    }
  }

  # Stream long-term records to Kinesis as they are written.
  dynamic "stream_delivery_resources" {
    for_each = local.has_kinesis ? [var.spec.kinesis_delivery] : []
    content {
      resource {
        kinesis {
          data_stream_arn = stream_delivery_resources.value.data_stream_arn
          content_configuration {
            # MEMORY_RECORDS is the only content type AWS defines.
            type  = "MEMORY_RECORDS"
            level = stream_delivery_resources.value.content_level != "" ? stream_delivery_resources.value.content_level : null
          }
        }
      }
    }
  }

  tags = local.aws_tags
}

# Long-term extraction pipelines. Built-in types run AWS-managed
# extraction; CUSTOM carries prompt/model overrides (exactly paired,
# spec-validated).
resource "aws_bedrockagentcore_memory_strategy" "this" {
  for_each = local.strategies

  memory_id = aws_bedrockagentcore_memory.this.id
  name      = each.value.name

  # Changing the type replaces the strategy (provider-enforced).
  type = each.value.type

  description = each.value.description != "" ? each.value.description : null

  # Required by the provider: its deprecated `namespaces` twin and this
  # field are an exactly-one pair, and the living surface is this one.
  namespace_templates = each.value.namespace_templates

  # Prompt/model overrides -- present exactly when type is CUSTOM
  # (spec-validated).
  dynamic "configuration" {
    for_each = each.value.custom != null ? [each.value.custom] : []
    content {
      type = configuration.value.type

      dynamic "extraction" {
        for_each = configuration.value.extraction != null ? [configuration.value.extraction] : []
        content {
          append_to_prompt = extraction.value.append_to_prompt
          model_id         = extraction.value.model_id
        }
      }

      dynamic "consolidation" {
        for_each = configuration.value.consolidation != null ? [configuration.value.consolidation] : []
        content {
          append_to_prompt = consolidation.value.append_to_prompt
          model_id         = consolidation.value.model_id
        }
      }

      dynamic "reflection" {
        for_each = configuration.value.reflection != null ? [configuration.value.reflection] : []
        content {
          append_to_prompt    = reflection.value.append_to_prompt
          model_id            = reflection.value.model_id
          namespace_templates = reflection.value.namespace_templates
        }
      }
    }
  }

  # EPISODIC reflection namespaces -- only legal on EPISODIC strategies
  # (spec-validated).
  dynamic "reflection_configuration" {
    for_each = length(each.value.reflection_namespace_templates) > 0 ? [each.value.reflection_namespace_templates] : []
    content {
      namespace_templates = reflection_configuration.value
    }
  }
}
