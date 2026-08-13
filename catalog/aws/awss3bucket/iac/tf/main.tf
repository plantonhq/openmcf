# AWS models the bucket as a small root resource plus one satellite resource
# per behavioral setting (versioning, encryption, lifecycle, ...). This module
# mirrors that grain: each folded spec block materializes as its own provider
# resource, created only when the spec asks for it. Two satellites are ALWAYS
# created — public access block and ownership controls — because their absence
# in the spec means the secure default, and stating that default explicitly in
# state is what keeps a bucket provably private.

resource "aws_s3_bucket" "this" {
  bucket = local.bucket_name

  # force_destroy empties the bucket (all versions and delete markers) before
  # deletion — without it a non-empty bucket fails the destroy.
  force_destroy = var.spec.force_destroy

  # Object Lock can only be turned on at creation (changing it replaces the
  # bucket), which is why it lives on the root resource and not in the
  # object-lock satellite below.
  object_lock_enabled = var.spec.object_lock_enabled

  # Namespace is create-time identity like the name itself: unset keeps the
  # classic global namespace; "account-regional" scopes the name to this
  # account+region. Changing it replaces the bucket.
  bucket_namespace = var.spec.bucket_namespace != "" ? var.spec.bucket_namespace : null

  tags = local.aws_tags
}

# --- Versioning ---------------------------------------------------------------
# Only managed when the spec sets a state: an unset spec leaves the bucket
# unversioned WITHOUT creating the satellite, so a pre-existing "never
# versioned" bucket state is representable. AWS cannot return a bucket to the
# never-versioned state once enabled — flipping Enabled -> "" is rejected by
# AWS at apply time; use Suspended.
resource "aws_s3_bucket_versioning" "this" {
  count = var.spec.versioning_status != "" ? 1 : 0

  bucket = aws_s3_bucket.this.id

  versioning_configuration {
    status = var.spec.versioning_status
  }
}

# --- Default encryption -------------------------------------------------------
# Created only when the spec opts in: since January 2023 AWS itself encrypts
# every object with SSE-S3 by default, so an absent block already means
# "encrypted with AES256" without any configuration to manage.
resource "aws_s3_bucket_server_side_encryption_configuration" "this" {
  count = local.manage_encryption ? 1 : 0

  bucket = aws_s3_bucket.this.id

  rule {
    apply_server_side_encryption_by_default {
      sse_algorithm     = local.sse_algorithm
      kms_master_key_id = local.sse_kms_key_id
    }

    # Bucket keys cut SSE-KMS request costs by up to 99%; harmless under
    # AES256 where AWS ignores the flag.
    bucket_key_enabled = var.spec.encryption.bucket_key_enabled

    # Encryption types PUTs may no longer use (e.g. "SSE-C"); sent only when
    # the manifest states a posture so unset keeps AWS's own default.
    blocked_encryption_types = length(var.spec.encryption.blocked_encryption_types) > 0 ? var.spec.encryption.blocked_encryption_types : null
  }
}

# --- Public access block (always managed) --------------------------------------
resource "aws_s3_bucket_public_access_block" "this" {
  bucket = aws_s3_bucket.this.id

  block_public_acls       = local.public_access_block.block_public_acls
  block_public_policy     = local.public_access_block.block_public_policy
  ignore_public_acls      = local.public_access_block.ignore_public_acls
  restrict_public_buckets = local.public_access_block.restrict_public_buckets
}

# --- Ownership controls (always managed) ---------------------------------------
resource "aws_s3_bucket_ownership_controls" "this" {
  bucket = aws_s3_bucket.this.id

  rule {
    object_ownership = local.object_ownership
  }
}

# --- Canned ACL ----------------------------------------------------------------
# Only valid when ownership controls re-enable ACLs (CEL guarantees the
# coupling). Ordering matters: the ACL PUT fails while BucketOwnerEnforced is
# in effect, so the ownership setting must land first.
resource "aws_s3_bucket_acl" "this" {
  count = var.spec.acl != "" ? 1 : 0

  bucket = aws_s3_bucket.this.id
  acl    = var.spec.acl

  depends_on = [aws_s3_bucket_ownership_controls.this]
}

# --- Bucket policy -------------------------------------------------------------
# Applied after the public access block: a policy granting public access is
# rejected while block_public_policy is on, so when a manifest relaxes the
# guard and adds a public-read statement in the same apply, the guard change
# must land first.
resource "aws_s3_bucket_policy" "this" {
  count = local.policy != null ? 1 : 0

  bucket = aws_s3_bucket.this.id
  policy = local.policy

  depends_on = [aws_s3_bucket_public_access_block.this]
}

# --- Lifecycle configuration ---------------------------------------------------
resource "aws_s3_bucket_lifecycle_configuration" "this" {
  count = length(local.lifecycle_rules) > 0 ? 1 : 0

  bucket = aws_s3_bucket.this.id

  transition_default_minimum_object_size = var.spec.transition_default_minimum_object_size != "" ? var.spec.transition_default_minimum_object_size : null

  dynamic "rule" {
    for_each = local.lifecycle_rules
    content {
      id     = rule.value.id
      status = rule.value.status

      # AWS expresses a single-predicate filter directly on the filter block
      # and multi-predicate filters inside an `and` wrapper. The spec exposes
      # one flat filter message and this block emits whichever document shape
      # the predicate count requires.
      dynamic "filter" {
        for_each = rule.value.filter != null ? [rule.value.filter] : []
        content {
          prefix                   = filter.value.predicate_count == 1 && filter.value.prefix != "" ? filter.value.prefix : null
          object_size_greater_than = filter.value.predicate_count == 1 && filter.value.object_size_greater_than > 0 ? filter.value.object_size_greater_than : null
          object_size_less_than    = filter.value.predicate_count == 1 && filter.value.object_size_less_than > 0 ? filter.value.object_size_less_than : null

          dynamic "tag" {
            for_each = filter.value.predicate_count == 1 && length(filter.value.tags) == 1 ? filter.value.tags : {}
            content {
              key   = tag.key
              value = tag.value
            }
          }

          dynamic "and" {
            for_each = filter.value.predicate_count > 1 ? [filter.value] : []
            content {
              prefix                   = and.value.prefix != "" ? and.value.prefix : null
              tags                     = length(and.value.tags) > 0 ? and.value.tags : null
              object_size_greater_than = and.value.object_size_greater_than > 0 ? and.value.object_size_greater_than : null
              object_size_less_than    = and.value.object_size_less_than > 0 ? and.value.object_size_less_than : null
            }
          }
        }
      }

      dynamic "transition" {
        for_each = rule.value.transitions
        content {
          # days is presence-typed in the contract (null when unset), so an
          # explicit 0 — AWS's "transition on the upload day" — passes
          # through; CEL enforces exactly-one of days/date.
          days          = transition.value.days
          date          = transition.value.date != "" ? transition.value.date : null
          storage_class = transition.value.storage_class
        }
      }

      dynamic "expiration" {
        for_each = rule.value.expiration != null ? [rule.value.expiration] : []
        content {
          days                         = expiration.value.days > 0 ? expiration.value.days : null
          date                         = expiration.value.date != "" ? expiration.value.date : null
          expired_object_delete_marker = expiration.value.expired_object_delete_marker ? true : null
        }
      }

      dynamic "noncurrent_version_transition" {
        for_each = rule.value.noncurrent_version_transitions
        content {
          noncurrent_days           = noncurrent_version_transition.value.noncurrent_days
          newer_noncurrent_versions = noncurrent_version_transition.value.newer_noncurrent_versions > 0 ? noncurrent_version_transition.value.newer_noncurrent_versions : null
          storage_class             = noncurrent_version_transition.value.storage_class
        }
      }

      dynamic "noncurrent_version_expiration" {
        for_each = rule.value.noncurrent_version_expiration != null ? [rule.value.noncurrent_version_expiration] : []
        content {
          noncurrent_days           = noncurrent_version_expiration.value.noncurrent_days
          newer_noncurrent_versions = noncurrent_version_expiration.value.newer_noncurrent_versions > 0 ? noncurrent_version_expiration.value.newer_noncurrent_versions : null
        }
      }

      dynamic "abort_incomplete_multipart_upload" {
        for_each = rule.value.abort_incomplete_multipart_upload_days > 0 ? [rule.value.abort_incomplete_multipart_upload_days] : []
        content {
          days_after_initiation = abort_incomplete_multipart_upload.value
        }
      }
    }
  }

  # Lifecycle actions on noncurrent versions only make sense once versioning
  # exists; ordering avoids a transient AWS validation window.
  depends_on = [aws_s3_bucket_versioning.this]
}

# --- Replication ---------------------------------------------------------------
resource "aws_s3_bucket_replication_configuration" "this" {
  count = var.spec.replication != null ? 1 : 0

  bucket = aws_s3_bucket.this.id
  role   = var.spec.replication.role_arn

  dynamic "rule" {
    for_each = local.replication_rules
    content {
      id       = rule.value.id
      priority = rule.value.priority
      status   = rule.value.status

      # Same single-vs-and document shaping as lifecycle filters. A rule
      # without a filter covers the whole bucket; AWS still requires the
      # delete_marker_replication block on the modern (V2) schema, emitted
      # unconditionally from the spec's bool.
      dynamic "filter" {
        for_each = rule.value.filter != null ? [rule.value.filter] : []
        content {
          prefix = filter.value.predicate_count == 1 && filter.value.prefix != "" ? filter.value.prefix : null

          dynamic "tag" {
            for_each = filter.value.predicate_count == 1 && length(filter.value.tags) == 1 ? filter.value.tags : {}
            content {
              key   = tag.key
              value = tag.value
            }
          }

          dynamic "and" {
            for_each = filter.value.predicate_count > 1 ? [filter.value] : []
            content {
              prefix = and.value.prefix != "" ? and.value.prefix : null
              tags   = length(and.value.tags) > 0 ? and.value.tags : null
            }
          }
        }
      }

      delete_marker_replication {
        status = rule.value.delete_marker_replication ? "Enabled" : "Disabled"
      }

      dynamic "existing_object_replication" {
        for_each = rule.value.existing_object_replication ? [1] : []
        content {
          status = "Enabled"
        }
      }

      dynamic "source_selection_criteria" {
        for_each = (rule.value.replicate_replica_modifications || rule.value.replicate_sse_kms_encrypted_objects) ? [1] : []
        content {
          dynamic "replica_modifications" {
            for_each = rule.value.replicate_replica_modifications ? [1] : []
            content {
              status = "Enabled"
            }
          }
          dynamic "sse_kms_encrypted_objects" {
            for_each = rule.value.replicate_sse_kms_encrypted_objects ? [1] : []
            content {
              status = "Enabled"
            }
          }
        }
      }

      destination {
        bucket        = rule.value.destination.bucket_arn
        account       = rule.value.destination.account != "" ? rule.value.destination.account : null
        storage_class = rule.value.destination.storage_class != "" ? rule.value.destination.storage_class : null

        dynamic "access_control_translation" {
          for_each = rule.value.destination.change_replica_ownership_to_destination ? [1] : []
          content {
            owner = "Destination"
          }
        }

        dynamic "encryption_configuration" {
          for_each = rule.value.destination.replica_kms_key_id != "" ? [rule.value.destination.replica_kms_key_id] : []
          content {
            replica_kms_key_id = encryption_configuration.value
          }
        }

        # AWS requires metrics whenever RTC is enabled (CEL mirrors this);
        # both use the fixed 15-minute threshold AWS accepts.
        dynamic "metrics" {
          for_each = rule.value.destination.metrics_enabled ? [1] : []
          content {
            status = "Enabled"
            event_threshold {
              minutes = 15
            }
          }
        }

        dynamic "replication_time" {
          for_each = rule.value.destination.replication_time_control_enabled ? [1] : []
          content {
            status = "Enabled"
            time {
              minutes = 15
            }
          }
        }
      }
    }
  }

  # AWS rejects replication configuration until versioning is Enabled on the
  # source bucket, and the two PUTs race without an explicit edge.
  depends_on = [aws_s3_bucket_versioning.this]
}

# --- Static website hosting ----------------------------------------------------
resource "aws_s3_bucket_website_configuration" "this" {
  count = var.spec.website != null ? 1 : 0

  bucket = aws_s3_bucket.this.id

  dynamic "index_document" {
    for_each = var.spec.website.index_document_suffix != "" ? [var.spec.website.index_document_suffix] : []
    content {
      suffix = index_document.value
    }
  }

  dynamic "error_document" {
    for_each = var.spec.website.error_document_key != "" ? [var.spec.website.error_document_key] : []
    content {
      key = error_document.value
    }
  }

  dynamic "redirect_all_requests_to" {
    for_each = local.website_redirect_all != null ? [local.website_redirect_all] : []
    content {
      host_name = redirect_all_requests_to.value.host_name
      protocol  = redirect_all_requests_to.value.protocol != "" ? redirect_all_requests_to.value.protocol : null
    }
  }

  dynamic "routing_rule" {
    for_each = var.spec.website.routing_rules
    content {
      dynamic "condition" {
        for_each = routing_rule.value.condition != null ? [routing_rule.value.condition] : []
        content {
          http_error_code_returned_equals = condition.value.http_error_code_returned_equals != "" ? condition.value.http_error_code_returned_equals : null
          key_prefix_equals               = condition.value.key_prefix_equals != "" ? condition.value.key_prefix_equals : null
        }
      }
      redirect {
        host_name               = routing_rule.value.redirect.host_name != "" ? routing_rule.value.redirect.host_name : null
        http_redirect_code      = routing_rule.value.redirect.http_redirect_code != "" ? routing_rule.value.redirect.http_redirect_code : null
        protocol                = routing_rule.value.redirect.protocol != "" ? routing_rule.value.redirect.protocol : null
        replace_key_prefix_with = routing_rule.value.redirect.replace_key_prefix_with != "" ? routing_rule.value.redirect.replace_key_prefix_with : null
        replace_key_with        = routing_rule.value.redirect.replace_key_with != "" ? routing_rule.value.redirect.replace_key_with : null
      }
    }
  }
}

# --- Server access logging -----------------------------------------------------
resource "aws_s3_bucket_logging" "this" {
  count = var.spec.logging != null ? 1 : 0

  bucket = aws_s3_bucket.this.id

  target_bucket = var.spec.logging.target_bucket
  target_prefix = var.spec.logging.target_prefix

  # Partitioned key format makes access logs directly queryable via Athena
  # date partitions; the flat SimplePrefix format is the AWS default.
  dynamic "target_object_key_format" {
    for_each = var.spec.logging.partitioned_prefix_date_source != "" ? [var.spec.logging.partitioned_prefix_date_source] : []
    content {
      partitioned_prefix {
        partition_date_source = target_object_key_format.value
      }
    }
  }
}

# --- CORS ----------------------------------------------------------------------
resource "aws_s3_bucket_cors_configuration" "this" {
  count = length(var.spec.cors_rules) > 0 ? 1 : 0

  bucket = aws_s3_bucket.this.id

  dynamic "cors_rule" {
    for_each = var.spec.cors_rules
    content {
      id              = cors_rule.value.id != "" ? cors_rule.value.id : null
      allowed_methods = cors_rule.value.allowed_methods
      allowed_origins = cors_rule.value.allowed_origins
      allowed_headers = length(cors_rule.value.allowed_headers) > 0 ? cors_rule.value.allowed_headers : null
      expose_headers  = length(cors_rule.value.expose_headers) > 0 ? cors_rule.value.expose_headers : null
      max_age_seconds = cors_rule.value.max_age_seconds > 0 ? cors_rule.value.max_age_seconds : null
    }
  }
}

# --- Event notifications --------------------------------------------------------
# One notification resource carries all four arms. SQS/SNS/Lambda targets must
# already permit S3 delivery (queue/topic policy or Lambda invoke permission)
# or AWS rejects this PUT — the EventBridge arm needs no grant.
resource "aws_s3_bucket_notification" "this" {
  count = var.spec.notification != null ? 1 : 0

  bucket = aws_s3_bucket.this.id

  eventbridge = var.spec.notification.eventbridge

  dynamic "lambda_function" {
    for_each = var.spec.notification.lambda_functions
    content {
      lambda_function_arn = lambda_function.value.lambda_function_arn
      events              = lambda_function.value.events
      filter_prefix       = lambda_function.value.filter_prefix != "" ? lambda_function.value.filter_prefix : null
      filter_suffix       = lambda_function.value.filter_suffix != "" ? lambda_function.value.filter_suffix : null
    }
  }

  dynamic "queue" {
    for_each = var.spec.notification.queues
    content {
      queue_arn     = queue.value.queue_arn
      events        = queue.value.events
      filter_prefix = queue.value.filter_prefix != "" ? queue.value.filter_prefix : null
      filter_suffix = queue.value.filter_suffix != "" ? queue.value.filter_suffix : null
    }
  }

  dynamic "topic" {
    for_each = var.spec.notification.topics
    content {
      topic_arn     = topic.value.topic_arn
      events        = topic.value.events
      filter_prefix = topic.value.filter_prefix != "" ? topic.value.filter_prefix : null
      filter_suffix = topic.value.filter_suffix != "" ? topic.value.filter_suffix : null
    }
  }
}

# --- Object Lock default retention ----------------------------------------------
# The root resource's object_lock_enabled makes the bucket lock-capable; this
# satellite adds the default retention window applied to every new object.
resource "aws_s3_bucket_object_lock_configuration" "this" {
  count = var.spec.object_lock_default_retention != null ? 1 : 0

  bucket = aws_s3_bucket.this.id

  rule {
    default_retention {
      mode  = var.spec.object_lock_default_retention.mode
      days  = var.spec.object_lock_default_retention.days > 0 ? var.spec.object_lock_default_retention.days : null
      years = var.spec.object_lock_default_retention.years > 0 ? var.spec.object_lock_default_retention.years : null
    }
  }
}

# --- Transfer acceleration -------------------------------------------------------
resource "aws_s3_bucket_accelerate_configuration" "this" {
  count = var.spec.acceleration_status != "" ? 1 : 0

  bucket = aws_s3_bucket.this.id
  status = var.spec.acceleration_status
}

# --- Requester pays ---------------------------------------------------------------
resource "aws_s3_bucket_request_payment_configuration" "this" {
  count = var.spec.request_payer != "" ? 1 : 0

  bucket = aws_s3_bucket.this.id
  payer  = var.spec.request_payer
}

# --- Intelligent-Tiering archive configurations -----------------------------------
# Many-per-bucket satellite: each named configuration is its own provider
# resource keyed by name, so adding or removing one never disturbs the others.
resource "aws_s3_bucket_intelligent_tiering_configuration" "this" {
  for_each = local.intelligent_tiering_configurations

  bucket = aws_s3_bucket.this.id
  name   = each.value.name
  status = each.value.status != "" ? each.value.status : "Enabled"

  dynamic "filter" {
    for_each = (each.value.filter_prefix != "" || length(each.value.filter_tags) > 0) ? [1] : []
    content {
      prefix = each.value.filter_prefix != "" ? each.value.filter_prefix : null
      tags   = length(each.value.filter_tags) > 0 ? each.value.filter_tags : null
    }
  }

  dynamic "tiering" {
    for_each = each.value.tiers
    content {
      access_tier = tiering.value.access_tier
      days        = tiering.value.days
    }
  }
}

# --- ABAC (attribute-based access control) ---------------------------------------
# Managed only when the spec states a posture: unset leaves the bucket at
# AWS's default (disabled) without creating the satellite.
resource "aws_s3_bucket_abac" "this" {
  count = var.spec.abac_status != "" ? 1 : 0

  bucket = aws_s3_bucket.this.id

  abac_status {
    status = var.spec.abac_status
  }
}

# --- Storage-class-analysis configurations -----------------------------------------
# Many-per-bucket satellite keyed by name, mirroring intelligent tiering.
# AWS's storage-class-analysis export accepts only its V_1 schema and CSV
# format, so those provider arguments are left to their defaults and only the
# destination is user surface.
resource "aws_s3_bucket_analytics_configuration" "this" {
  for_each = local.analytics_configurations

  bucket = aws_s3_bucket.this.id
  name   = each.value.name

  dynamic "filter" {
    for_each = (each.value.filter_prefix != "" || length(each.value.filter_tags) > 0) ? [1] : []
    content {
      prefix = each.value.filter_prefix != "" ? each.value.filter_prefix : null
      tags   = length(each.value.filter_tags) > 0 ? each.value.filter_tags : null
    }
  }

  dynamic "storage_class_analysis" {
    for_each = each.value.export != null ? [each.value.export] : []
    content {
      data_export {
        destination {
          s3_bucket_destination {
            bucket_arn        = storage_class_analysis.value.bucket_arn
            bucket_account_id = storage_class_analysis.value.bucket_account_id != "" ? storage_class_analysis.value.bucket_account_id : null
            prefix            = storage_class_analysis.value.prefix != "" ? storage_class_analysis.value.prefix : null
          }
        }
      }
    }
  }
}

# --- Inventory report configurations ------------------------------------------------
# Many-per-bucket satellite keyed by name. `enabled` is always sent (AWS's
# InventoryConfiguration requires the IsEnabled member either way), derived
# from the spec's `disabled` so an unset spec matches AWS's active default.
resource "aws_s3_bucket_inventory" "this" {
  for_each = local.inventory_configurations

  bucket = aws_s3_bucket.this.id
  name   = each.value.name

  enabled                  = !each.value.disabled
  included_object_versions = each.value.included_object_versions

  optional_fields = length(each.value.optional_fields) > 0 ? each.value.optional_fields : null

  schedule {
    frequency = each.value.frequency
  }

  dynamic "filter" {
    for_each = each.value.filter_prefix != "" ? [each.value.filter_prefix] : []
    content {
      prefix = filter.value
    }
  }

  destination {
    bucket {
      bucket_arn = each.value.destination.bucket_arn
      format     = each.value.destination.format
      account_id = each.value.destination.account_id != "" ? each.value.destination.account_id : null
      prefix     = each.value.destination.prefix != "" ? each.value.destination.prefix : null

      # CEL guarantees at most one encryption arm; the block is emitted only
      # when an arm is actually chosen (an empty encryption block would send
      # an empty API struct).
      dynamic "encryption" {
        for_each = (each.value.destination.sse_s3 || each.value.destination.sse_kms_key_id != "") ? [1] : []
        content {
          dynamic "sse_kms" {
            for_each = each.value.destination.sse_kms_key_id != "" ? [each.value.destination.sse_kms_key_id] : []
            content {
              key_id = sse_kms.value
            }
          }
          dynamic "sse_s3" {
            for_each = each.value.destination.sse_s3 ? [1] : []
            content {}
          }
        }
      }
    }
  }
}

# --- CloudWatch request-metrics configurations ---------------------------------------
# Many-per-bucket satellite keyed by name. No filter block means metrics for
# every request against the bucket.
resource "aws_s3_bucket_metric" "this" {
  for_each = local.metrics_configurations

  bucket = aws_s3_bucket.this.id
  name   = each.value.name

  dynamic "filter" {
    for_each = (each.value.filter_access_point_arn != "" || each.value.filter_prefix != "" || length(each.value.filter_tags) > 0) ? [1] : []
    content {
      access_point = each.value.filter_access_point_arn != "" ? each.value.filter_access_point_arn : null
      prefix       = each.value.filter_prefix != "" ? each.value.filter_prefix : null
      tags         = length(each.value.filter_tags) > 0 ? each.value.filter_tags : null
    }
  }
}

# --- S3 Metadata (queryable metadata tables) ------------------------------------------
# The provider requires both table blocks: the journal table's expiration
# policy and the inventory table's state are stated explicitly either way.
# The destination table bucket/namespace are AWS-managed read-only state.
resource "aws_s3_bucket_metadata_configuration" "this" {
  count = var.spec.metadata_configuration != null ? 1 : 0

  bucket = aws_s3_bucket.this.id

  metadata_configuration {
    inventory_table_configuration {
      configuration_state = var.spec.metadata_configuration.inventory_table_enabled ? "ENABLED" : "DISABLED"

      dynamic "encryption_configuration" {
        for_each = var.spec.metadata_configuration.inventory_table_encryption != null ? [var.spec.metadata_configuration.inventory_table_encryption] : []
        content {
          sse_algorithm = encryption_configuration.value.sse_algorithm
          kms_key_arn   = encryption_configuration.value.kms_key_arn != "" ? encryption_configuration.value.kms_key_arn : null
        }
      }
    }

    journal_table_configuration {
      record_expiration {
        expiration = var.spec.metadata_configuration.journal_record_expiration.enabled ? "ENABLED" : "DISABLED"
        # days is only legal alongside ENABLED (CEL keeps it 0 when disabled).
        days = var.spec.metadata_configuration.journal_record_expiration.days > 0 ? var.spec.metadata_configuration.journal_record_expiration.days : null
      }

      dynamic "encryption_configuration" {
        for_each = var.spec.metadata_configuration.journal_table_encryption != null ? [var.spec.metadata_configuration.journal_table_encryption] : []
        content {
          sse_algorithm = encryption_configuration.value.sse_algorithm
          kms_key_arn   = encryption_configuration.value.kms_key_arn != "" ? encryption_configuration.value.kms_key_arn : null
        }
      }
    }
  }
}
