# The DynamoDB table: key schema and indexes, capacity, streams,
# global-table replication, encryption, and recovery. Create-only in
# AWS: the table name, the primary key schema, and every local
# secondary index (LSIs can never be added or removed later).
# Everything else -- billing mode, GSIs, streams, replicas, TTL, PITR,
# SSE key, table class -- edits in place. The table composes onto its
# neighbors instead of embedding them: the KMS key, the Kinesis stream,
# and the S3 import bucket all attach by reference.
resource "aws_dynamodb_table" "this" {
  name = local.table_name

  # Empty keeps the AWS default (PROVISIONED). Spec values are the
  # exact AWS API strings, so they pass through untranslated.
  billing_mode = var.spec.billing_mode != "" ? var.spec.billing_mode : null

  # Reserved capacity applies only when the effective billing mode is
  # PROVISIONED (CEL enforces the coupling either way). These values seed
  # the table at create; live capacity is owned by Application Auto
  # Scaling from then on (see autoscaling.tf and the ignore_changes note
  # below), so declared capacity changes land through the scalable
  # targets, never through this resource.
  read_capacity  = var.spec.provisioned_throughput != null ? var.spec.provisioned_throughput.read_capacity_units : null
  write_capacity = var.spec.provisioned_throughput != null ? var.spec.provisioned_throughput.write_capacity_units : null

  # Key attributes only -- DynamoDB is schemaless beyond the keys.
  # Restore-created tables inherit schema from the source, so the list
  # may be empty (CEL enforces the coupling).
  dynamic "attribute" {
    for_each = var.spec.attribute_definitions
    content {
      name = attribute.value.name
      type = attribute.value.type
    }
  }

  hash_key  = local.table_hash_key
  range_key = local.table_range_key

  # On-demand ceilings throttle instead of billing past the cap; -1
  # explicitly removes a previously-set ceiling (0 = not configured).
  dynamic "on_demand_throughput" {
    for_each = var.spec.on_demand_throughput != null ? [var.spec.on_demand_throughput] : []
    content {
      max_read_request_units  = on_demand_throughput.value.max_read_request_units != 0 ? on_demand_throughput.value.max_read_request_units : null
      max_write_request_units = on_demand_throughput.value.max_write_request_units != 0 ? on_demand_throughput.value.max_write_request_units : null
    }
  }

  # Warm throughput only ever increases -- AWS replaces the table on a
  # decrease (the provider marks it ForceNew for that reason).
  dynamic "warm_throughput" {
    for_each = var.spec.warm_throughput != null ? [var.spec.warm_throughput] : []
    content {
      read_units_per_second  = warm_throughput.value.read_units_per_second != 0 ? warm_throughput.value.read_units_per_second : null
      write_units_per_second = warm_throughput.value.write_units_per_second != 0 ? warm_throughput.value.write_units_per_second : null
    }
  }

  # GSIs edit in place (AWS serializes to one index mutation at a
  # time). The modern key_schema shape carries multi-attribute keys
  # (1-4 HASH elements first, then 0-4 RANGE); per-index capacity
  # follows the table's billing mode (CEL enforces the coupling).
  dynamic "global_secondary_index" {
    for_each = var.spec.global_secondary_indexes
    content {
      name = global_secondary_index.value.name

      dynamic "key_schema" {
        for_each = global_secondary_index.value.key_schema
        content {
          attribute_name = key_schema.value.attribute_name
          key_type       = key_schema.value.key_type
        }
      }

      projection_type    = global_secondary_index.value.projection.type
      non_key_attributes = length(global_secondary_index.value.projection.non_key_attributes) > 0 ? global_secondary_index.value.projection.non_key_attributes : null

      read_capacity  = global_secondary_index.value.provisioned_throughput != null ? global_secondary_index.value.provisioned_throughput.read_capacity_units : null
      write_capacity = global_secondary_index.value.provisioned_throughput != null ? global_secondary_index.value.provisioned_throughput.write_capacity_units : null

      dynamic "on_demand_throughput" {
        for_each = global_secondary_index.value.on_demand_throughput != null ? [global_secondary_index.value.on_demand_throughput] : []
        content {
          max_read_request_units  = on_demand_throughput.value.max_read_request_units != 0 ? on_demand_throughput.value.max_read_request_units : null
          max_write_request_units = on_demand_throughput.value.max_write_request_units != 0 ? on_demand_throughput.value.max_write_request_units : null
        }
      }

      dynamic "warm_throughput" {
        for_each = global_secondary_index.value.warm_throughput != null ? [global_secondary_index.value.warm_throughput] : []
        content {
          read_units_per_second  = warm_throughput.value.read_units_per_second != 0 ? warm_throughput.value.read_units_per_second : null
          write_units_per_second = warm_throughput.value.write_units_per_second != 0 ? warm_throughput.value.write_units_per_second : null
        }
      }
    }
  }

  # LSIs are create-only: they can never be added or removed after the
  # table exists, and their presence permanently caps each item
  # collection at 10 GB.
  dynamic "local_secondary_index" {
    for_each = var.spec.local_secondary_indexes
    content {
      name               = local_secondary_index.value.name
      range_key          = local_secondary_index.value.range_key
      projection_type    = local_secondary_index.value.projection.type
      non_key_attributes = length(local_secondary_index.value.projection.non_key_attributes) > 0 ? local_secondary_index.value.projection.non_key_attributes : null
    }
  }

  # TTL deletes expired items free of write cost, within ~48h of the
  # epoch-seconds value in the named attribute passing.
  dynamic "ttl" {
    for_each = var.spec.ttl != null ? [var.spec.ttl] : []
    content {
      enabled        = ttl.value.enabled
      attribute_name = ttl.value.attribute_name != "" ? ttl.value.attribute_name : null
    }
  }

  # Streams carry item-level change data; global tables require the
  # NEW_AND_OLD_IMAGES view (CEL enforces it).
  stream_enabled   = var.spec.stream_enabled
  stream_view_type = var.spec.stream_view_type != "" ? var.spec.stream_view_type : null

  # Continuous backups with per-second restore granularity; 0 keeps the
  # AWS default recovery window (35 days).
  dynamic "point_in_time_recovery" {
    for_each = var.spec.point_in_time_recovery != null ? [var.spec.point_in_time_recovery] : []
    content {
      enabled                 = point_in_time_recovery.value.enabled
      recovery_period_in_days = point_in_time_recovery.value.recovery_period_in_days != 0 ? point_in_time_recovery.value.recovery_period_in_days : null
    }
  }

  # DynamoDB always encrypts; this block switches from the AWS-owned
  # key to the AWS-managed aws/dynamodb key (no KMS ARN) or a
  # customer-managed key (ARN set).
  dynamic "server_side_encryption" {
    for_each = var.spec.server_side_encryption != null ? [var.spec.server_side_encryption] : []
    content {
      enabled     = server_side_encryption.value.enabled
      kms_key_arn = server_side_encryption.value.kms_key_arn != "" ? server_side_encryption.value.kms_key_arn : null
    }
  }

  # Empty keeps the AWS default (STANDARD).
  table_class = var.spec.table_class != "" ? var.spec.table_class : null

  deletion_protection_enabled = var.spec.deletion_protection_enabled

  # Global Tables v2: each replica is an active read/write copy in
  # another region; each region encrypts independently, so the KMS key
  # is per-replica.
  dynamic "replica" {
    for_each = var.spec.replicas
    content {
      region_name                 = replica.value.region_name
      kms_key_arn                 = replica.value.kms_key_arn != "" ? replica.value.kms_key_arn : null
      point_in_time_recovery      = replica.value.point_in_time_recovery
      deletion_protection_enabled = replica.value.deletion_protection_enabled
      propagate_tags              = replica.value.propagate_tags
      consistency_mode            = replica.value.consistency_mode != "" ? replica.value.consistency_mode : null
    }
  }

  # The MRSC witness persists replicated writes for quorum but serves
  # no reads or writes; CEL pins the exact topology AWS accepts.
  dynamic "global_table_witness" {
    for_each = var.spec.global_table_witness != null ? [var.spec.global_table_witness] : []
    content {
      region_name = global_table_witness.value.region_name
    }
  }

  # Create sources are mutually exclusive (CEL): a point-in-time
  # restore by name or ARN, a backup restore, or an S3 import.
  restore_source_name      = var.spec.restore_source_name != "" ? var.spec.restore_source_name : null
  restore_source_table_arn = var.spec.restore_source_table_arn != "" ? var.spec.restore_source_table_arn : null
  restore_date_time        = var.spec.restore_date_time != "" ? var.spec.restore_date_time : null
  restore_to_latest_time   = var.spec.restore_to_latest_time ? true : null
  restore_backup_arn       = var.spec.restore_backup_arn != "" ? var.spec.restore_backup_arn : null

  # S3 import seeds a brand-new table -- billed as a one-time import
  # instead of per-item writes.
  dynamic "import_table" {
    for_each = var.spec.import_table != null ? [var.spec.import_table] : []
    content {
      input_format           = import_table.value.input_format
      input_compression_type = import_table.value.input_compression_type != "" ? import_table.value.input_compression_type : null

      s3_bucket_source {
        bucket       = import_table.value.s3_bucket
        bucket_owner = import_table.value.s3_bucket_owner != "" ? import_table.value.s3_bucket_owner : null
        key_prefix   = import_table.value.s3_key_prefix != "" ? import_table.value.s3_key_prefix : null
      }

      dynamic "input_format_options" {
        for_each = import_table.value.csv != null ? [import_table.value.csv] : []
        content {
          csv {
            delimiter   = input_format_options.value.delimiter != "" ? input_format_options.value.delimiter : null
            header_list = length(input_format_options.value.header_list) > 0 ? input_format_options.value.header_list : null
          }
        }
      }
    }
  }

  tags = local.aws_tags

  # Live capacity is Application Auto Scaling's to manage on provisioned
  # tables (autoscaling.tf registers targets in BOTH modes: user bounds
  # with tracking policies, or pinned min = max from
  # provisioned_throughput). ignore_changes must be a static literal list
  # (OpenTofu forbids a conditional expression here), so this always
  # ignores -- matching the house convention for autoscaler-managed
  # values (the Pulumi module ignores the same paths). Declared capacity
  # changes still land: the scalable targets carry them, and AWS moves
  # out-of-range capacity into the target's bounds by contract.
  lifecycle {
    ignore_changes = [read_capacity, write_capacity]
  }
}
