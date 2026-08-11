# The log group container. All satellite resources below key off its name and
# share its lifecycle — that is why they are folded into this module instead of
# being separate kinds.
resource "aws_cloudwatch_log_group" "this" {
  name = local.resource_name

  # Retention: null leaves the retention policy unmanaged (AWS default: never
  # expire). The provider issues a separate PutRetentionPolicy call for this.
  retention_in_days = local.retention_in_days

  # KMS encryption: customer-managed key for log data at rest. Associating or
  # clearing the key is an in-place update (AssociateKmsKey/DisassociateKmsKey).
  kms_key_id = local.kms_key_id

  # Log group class is create-time (ForceNew): STANDARD, INFREQUENT_ACCESS, or
  # DELIVERY. Null lets AWS default to STANDARD.
  log_group_class = local.log_group_class

  # Deletion protection blocks every delete (including IaC destroy) until the
  # flag is cleared; applied via a separate PutLogGroupDeletionProtection call.
  deletion_protection_enabled = local.deletion_protection_enabled

  tags = local.aws_tags
}

# Metric filters: one resource per named filter (many-per-group). Each turns
# matching log events into a CloudWatch metric — the raw material for alarms
# on signals that live only in logs.
resource "aws_cloudwatch_log_metric_filter" "this" {
  for_each = { for filter in var.spec.metric_filters : filter.name => filter }

  name           = each.value.name
  log_group_name = aws_cloudwatch_log_group.this.name

  # The provider requires pattern; empty string matches every event.
  pattern = each.value.pattern

  # Tri-state pass-through (the provider attribute is Optional+Computed): only
  # an explicit false switches an existing filter back to matching raw events —
  # an omitted value keeps the filter's current setting.
  apply_on_transformed_logs = each.value.apply_on_transformed_logs

  metric_transformation {
    name      = each.value.transformation.metric_name
    namespace = each.value.transformation.metric_namespace
    value     = each.value.transformation.metric_value

    # default_value is a tri-state: null publishes nothing for non-matching
    # periods. The provider types it as a nullable-float STRING, so format the
    # number only when present. (AWS forbids combining it with dimensions —
    # the spec CEL rejects that before plan time.)
    default_value = each.value.transformation.default_value != null ? tostring(each.value.transformation.default_value) : null

    dimensions = length(each.value.transformation.dimensions) > 0 ? each.value.transformation.dimensions : null

    # The provider defaults unit to "None" when omitted.
    unit = each.value.transformation.unit != "" ? each.value.transformation.unit : null
  }
}

# Subscription filters: real-time delivery of matching events to a Kinesis
# stream, Firehose delivery stream, or Lambda function. AWS allows at most two
# per group (CEL-enforced in the spec).
resource "aws_cloudwatch_log_subscription_filter" "this" {
  for_each = { for filter in var.spec.subscription_filters : filter.name => filter }

  name           = each.value.name
  log_group_name = aws_cloudwatch_log_group.this.name

  # References arrive pre-resolved to plain ARNs by the orchestrator.
  destination_arn = each.value.destination_arn

  # The provider requires filter_pattern; empty string delivers everything.
  filter_pattern = each.value.filter_pattern

  # Required for Kinesis/Firehose destinations (CloudWatch Logs assumes this
  # role to put records); Lambda destinations authorize via a Lambda
  # permission instead, so null is correct there.
  role_arn = each.value.role_arn != "" ? each.value.role_arn : null

  # Only meaningful for Kinesis stream destinations; the provider defaults to
  # ByLogStream (per-stream ordering preserved).
  distribution = each.value.distribution != "" ? each.value.distribution : null

  # Source account/region/log enrichment for centralized log destinations.
  emit_system_fields = length(each.value.emit_system_fields) > 0 ? each.value.emit_system_fields : null

  # Tri-state pass-through (the provider attribute is Optional+Computed): only
  # an explicit false switches an existing filter back to raw events.
  apply_on_transformed_logs = each.value.apply_on_transformed_logs
}

# Data protection policy: a single group-scoped policy document (PII audit +
# masking). The Struct spec field arrives as a NESTED OBJECT in tfvars, so it
# is jsonencode'd for the provider's document argument.
resource "aws_cloudwatch_log_data_protection_policy" "this" {
  count = var.spec.data_protection_policy != null ? 1 : 0

  log_group_name  = aws_cloudwatch_log_group.this.name
  policy_document = jsonencode(var.spec.data_protection_policy)
}

# Field index policy: a single group-scoped policy listing the log fields to
# index for faster, cheaper Logs Insights queries.
resource "aws_cloudwatch_log_index_policy" "this" {
  count = var.spec.field_index_policy != null ? 1 : 0

  log_group_name  = aws_cloudwatch_log_group.this.name
  policy_document = jsonencode(var.spec.field_index_policy)
}

# Pre-created log streams: one per declared name (most streams are created at
# runtime by the writing agent — these exist for fixed-name agent configs and
# stream-scoped IAM policies). Streams die with the group.
resource "aws_cloudwatch_log_stream" "this" {
  for_each = toset(var.spec.log_streams)

  name           = each.value
  log_group_name = aws_cloudwatch_log_group.this.name
}

# Log transformer: the group's ingest-time processor pipeline (one per group,
# STANDARD class only — CEL-enforced in the spec). Each processors entry
# carries exactly ONE processor; the blocks below render whichever arm is set.
# Optional strings are sent only when set (nested defaults are not
# materialized inside repeated messages); booleans inside entries are always
# sent — PutTransformer replaces the whole pipeline, so explicit false and
# omitted false are the same write.
resource "aws_cloudwatch_log_transformer" "this" {
  count = var.spec.transformer != null ? 1 : 0

  log_group_arn = aws_cloudwatch_log_group.this.arn

  dynamic "transformer_config" {
    for_each = var.spec.transformer.processors
    content {
      dynamic "add_keys" {
        for_each = transformer_config.value.add_keys != null ? [transformer_config.value.add_keys] : []
        content {
          dynamic "entry" {
            for_each = add_keys.value.entries
            content {
              key                 = entry.value.key
              value               = entry.value.value
              overwrite_if_exists = entry.value.overwrite_if_exists
            }
          }
        }
      }

      dynamic "copy_value" {
        for_each = transformer_config.value.copy_value != null ? [transformer_config.value.copy_value] : []
        content {
          dynamic "entry" {
            for_each = copy_value.value.entries
            content {
              source              = entry.value.source
              target              = entry.value.target
              overwrite_if_exists = entry.value.overwrite_if_exists
            }
          }
        }
      }

      dynamic "csv" {
        for_each = transformer_config.value.csv != null ? [transformer_config.value.csv] : []
        content {
          columns         = length(csv.value.columns) > 0 ? csv.value.columns : null
          delimiter       = csv.value.delimiter != "" ? csv.value.delimiter : null
          quote_character = csv.value.quote_character != "" ? csv.value.quote_character : null
          source          = csv.value.source != "" ? csv.value.source : null
        }
      }

      dynamic "date_time_converter" {
        for_each = transformer_config.value.date_time_converter != null ? [transformer_config.value.date_time_converter] : []
        content {
          source          = date_time_converter.value.source
          target          = date_time_converter.value.target
          match_patterns  = date_time_converter.value.match_patterns
          locale          = date_time_converter.value.locale != "" ? date_time_converter.value.locale : null
          source_timezone = date_time_converter.value.source_timezone != "" ? date_time_converter.value.source_timezone : null
          target_format   = date_time_converter.value.target_format != "" ? date_time_converter.value.target_format : null
          target_timezone = date_time_converter.value.target_timezone != "" ? date_time_converter.value.target_timezone : null
        }
      }

      dynamic "delete_keys" {
        for_each = transformer_config.value.delete_keys != null ? [transformer_config.value.delete_keys] : []
        content {
          with_keys = delete_keys.value.with_keys
        }
      }

      dynamic "grok" {
        for_each = transformer_config.value.grok != null ? [transformer_config.value.grok] : []
        content {
          match  = grok.value.match
          source = grok.value.source != "" ? grok.value.source : null
        }
      }

      dynamic "list_to_map" {
        for_each = transformer_config.value.list_to_map != null ? [transformer_config.value.list_to_map] : []
        content {
          source            = list_to_map.value.source
          key               = list_to_map.value.key
          value_key         = list_to_map.value.value_key != "" ? list_to_map.value.value_key : null
          target            = list_to_map.value.target != "" ? list_to_map.value.target : null
          flatten           = list_to_map.value.flatten
          flattened_element = list_to_map.value.flattened_element != "" ? list_to_map.value.flattened_element : null
        }
      }

      dynamic "lower_case_string" {
        for_each = transformer_config.value.lower_case_string != null ? [transformer_config.value.lower_case_string] : []
        content {
          with_keys = lower_case_string.value.with_keys
        }
      }

      dynamic "move_keys" {
        for_each = transformer_config.value.move_keys != null ? [transformer_config.value.move_keys] : []
        content {
          dynamic "entry" {
            for_each = move_keys.value.entries
            content {
              source              = entry.value.source
              target              = entry.value.target
              overwrite_if_exists = entry.value.overwrite_if_exists
            }
          }
        }
      }

      dynamic "parse_cloudfront" {
        for_each = transformer_config.value.parse_cloudfront != null ? [transformer_config.value.parse_cloudfront] : []
        content {
          source = parse_cloudfront.value.source != "" ? parse_cloudfront.value.source : null
        }
      }

      dynamic "parse_json" {
        for_each = transformer_config.value.parse_json != null ? [transformer_config.value.parse_json] : []
        content {
          source      = parse_json.value.source != "" ? parse_json.value.source : null
          destination = parse_json.value.destination != "" ? parse_json.value.destination : null
        }
      }

      dynamic "parse_key_value" {
        for_each = transformer_config.value.parse_key_value != null ? [transformer_config.value.parse_key_value] : []
        content {
          source              = parse_key_value.value.source != "" ? parse_key_value.value.source : null
          destination         = parse_key_value.value.destination != "" ? parse_key_value.value.destination : null
          field_delimiter     = parse_key_value.value.field_delimiter != "" ? parse_key_value.value.field_delimiter : null
          key_value_delimiter = parse_key_value.value.key_value_delimiter != "" ? parse_key_value.value.key_value_delimiter : null
          key_prefix          = parse_key_value.value.key_prefix != "" ? parse_key_value.value.key_prefix : null
          non_match_value     = parse_key_value.value.non_match_value != "" ? parse_key_value.value.non_match_value : null
          overwrite_if_exists = parse_key_value.value.overwrite_if_exists
        }
      }

      dynamic "parse_postgres" {
        for_each = transformer_config.value.parse_postgres != null ? [transformer_config.value.parse_postgres] : []
        content {
          source = parse_postgres.value.source != "" ? parse_postgres.value.source : null
        }
      }

      dynamic "parse_route53" {
        for_each = transformer_config.value.parse_route53 != null ? [transformer_config.value.parse_route53] : []
        content {
          source = parse_route53.value.source != "" ? parse_route53.value.source : null
        }
      }

      dynamic "parse_to_ocsf" {
        for_each = transformer_config.value.parse_to_ocsf != null ? [transformer_config.value.parse_to_ocsf] : []
        content {
          event_source = parse_to_ocsf.value.event_source
          ocsf_version = parse_to_ocsf.value.ocsf_version
          source       = parse_to_ocsf.value.source != "" ? parse_to_ocsf.value.source : null
        }
      }

      dynamic "parse_vpc" {
        for_each = transformer_config.value.parse_vpc != null ? [transformer_config.value.parse_vpc] : []
        content {
          source = parse_vpc.value.source != "" ? parse_vpc.value.source : null
        }
      }

      dynamic "parse_waf" {
        for_each = transformer_config.value.parse_waf != null ? [transformer_config.value.parse_waf] : []
        content {
          source = parse_waf.value.source != "" ? parse_waf.value.source : null
        }
      }

      dynamic "rename_keys" {
        for_each = transformer_config.value.rename_keys != null ? [transformer_config.value.rename_keys] : []
        content {
          dynamic "entry" {
            for_each = rename_keys.value.entries
            content {
              key                 = entry.value.key
              rename_to           = entry.value.rename_to
              overwrite_if_exists = entry.value.overwrite_if_exists
            }
          }
        }
      }

      dynamic "split_string" {
        for_each = transformer_config.value.split_string != null ? [transformer_config.value.split_string] : []
        content {
          dynamic "entry" {
            for_each = split_string.value.entries
            content {
              source    = entry.value.source
              delimiter = entry.value.delimiter
            }
          }
        }
      }

      dynamic "substitute_string" {
        for_each = transformer_config.value.substitute_string != null ? [transformer_config.value.substitute_string] : []
        content {
          dynamic "entry" {
            for_each = substitute_string.value.entries
            content {
              source = entry.value.source
              from   = entry.value.from
              to     = entry.value.to
            }
          }
        }
      }

      dynamic "trim_string" {
        for_each = transformer_config.value.trim_string != null ? [transformer_config.value.trim_string] : []
        content {
          with_keys = trim_string.value.with_keys
        }
      }

      dynamic "type_converter" {
        for_each = transformer_config.value.type_converter != null ? [transformer_config.value.type_converter] : []
        content {
          dynamic "entry" {
            for_each = type_converter.value.entries
            content {
              key  = entry.value.key
              type = entry.value.type
            }
          }
        }
      }

      dynamic "upper_case_string" {
        for_each = transformer_config.value.upper_case_string != null ? [transformer_config.value.upper_case_string] : []
        content {
          with_keys = upper_case_string.value.with_keys
        }
      }
    }
  }
}
