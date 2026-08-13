resource "aws_s3_object" "objects" {
  for_each = local.content_objects_map

  bucket = local.bucket_name
  key    = each.value.key

  # Content source: the spec guarantees exactly one of content/content_base64
  # is set (proto oneof + CEL); the unset arm arrives as "" and null here so
  # the provider sees only the chosen one.
  content        = each.value.content != "" ? each.value.content : null
  content_base64 = each.value.content_base64 != "" ? each.value.content_base64 : null

  # Sent explicitly, resolving the spec default (application/octet-stream)
  # here rather than leaving it to the provider: an omitted Content-Type would
  # make S3 store its own default (binary/octet-stream), diverging from what
  # the manifest declares and from the Pulumi engine.
  content_type = coalesce(each.value.content_type, "application/octet-stream")

  # HTTP presentation headers. Empty strings become null so unset stays
  # indistinguishable from the AWS defaults (no phantom headers stored).
  cache_control       = each.value.cache_control != "" ? each.value.cache_control : null
  content_encoding    = each.value.content_encoding != "" ? each.value.content_encoding : null
  content_disposition = each.value.content_disposition != "" ? each.value.content_disposition : null
  content_language    = each.value.content_language != "" ? each.value.content_language : null

  # User metadata (x-amz-meta-*). Keys are spec-validated lowercase, matching
  # what S3 stores, so reads never drift from the manifest.
  metadata = length(each.value.metadata) > 0 ? each.value.metadata : null

  # Website redirect: inert unless the bucket has static website hosting.
  website_redirect = each.value.website_redirect != "" ? each.value.website_redirect : null

  # Storage placement. Unset means STANDARD (AWS's default) — pass null so
  # the provider computes rather than pins a value.
  storage_class = each.value.storage_class != "" ? each.value.storage_class : null

  # Per-object encryption OVERRIDE. Unset inherits the bucket's default
  # encryption, which is where uniform posture belongs. A kms_key alone
  # implies aws:kms (the provider sets ServerSideEncryption when SSEKMSKeyId
  # is sent); the spec CEL rejects the contradictory kms_key + AES256 pair.
  server_side_encryption = each.value.server_side_encryption != "" ? each.value.server_side_encryption : null
  kms_key_id             = each.value.kms_key != "" ? each.value.kms_key : null
  bucket_key_enabled     = each.value.bucket_key_enabled

  # Upload-integrity checksum, stored alongside the object.
  checksum_algorithm = each.value.checksum_algorithm != "" ? each.value.checksum_algorithm : null

  # Object Lock (requires an Object Lock-enabled bucket; the spec CEL
  # guarantees mode and retain-until arrive as a pair).
  object_lock_mode              = each.value.object_lock_mode != "" ? each.value.object_lock_mode : null
  object_lock_retain_until_date = each.value.object_lock_retain_until_date != "" ? each.value.object_lock_retain_until_date : null
  object_lock_legal_hold_status = each.value.object_lock_legal_hold_status != "" ? each.value.object_lock_legal_hold_status : null

  # Canned ACL — only valid on buckets whose ownership controls permit ACLs;
  # modern (BucketOwnerEnforced) buckets reject it at apply time.
  acl = each.value.acl != "" ? each.value.acl : null

  # The GOVERNANCE-retention delete bypass (x-amz-bypass-governance-retention).
  # Only valid on Object Lock-enabled buckets — S3 rejects the flag on regular
  # buckets, failing the destroy. Versioned-bucket cleanup needs no flag:
  # deleting an object always removes all of its versions.
  force_destroy = each.value.force_destroy

  # Merge order: identity tags < set-level tags < object-level tags, so an
  # object can specialize but never lose its resource-identity attribution.
  tags = merge(local.set_tags, try(each.value.tags, {}))
}

resource "aws_s3_object_copy" "objects" {
  for_each = local.copy_objects_map

  # A copy may name another object of this same set as its source; the source
  # string carries no dependency the graph can see, so copies always wait for
  # the set's content objects.
  depends_on = [aws_s3_object.objects]

  bucket = local.bucket_name
  key    = each.value.key

  # "bucket/key", or "arn.../object/key" for access-point sources (composed in
  # locals.tf). The provider URL-escapes the whole string, which is also why a
  # ?versionId= selector cannot work at this pin — copies take the source's
  # current version.
  source = local.copy_sources[each.key]

  # Metadata directive: REPLACE writes this entry's metadata/headers to the
  # copy; the omitted default (COPY) preserves everything the source object
  # carried. AWS silently ignores header values sent under COPY, so headers
  # are gated here exactly as the spec CEL gates them — content_type included,
  # because the manifest loader materializes its default on every object and
  # an ungated send would echo-drift against the source's real content type.
  metadata_directive  = each.value.copy_from.replace_metadata ? "REPLACE" : null
  content_type        = each.value.copy_from.replace_metadata ? coalesce(each.value.content_type, "application/octet-stream") : null
  cache_control       = each.value.copy_from.replace_metadata && each.value.cache_control != "" ? each.value.cache_control : null
  content_encoding    = each.value.copy_from.replace_metadata && each.value.content_encoding != "" ? each.value.content_encoding : null
  content_disposition = each.value.copy_from.replace_metadata && each.value.content_disposition != "" ? each.value.content_disposition : null
  content_language    = each.value.copy_from.replace_metadata && each.value.content_language != "" ? each.value.content_language : null
  metadata            = each.value.copy_from.replace_metadata && length(each.value.metadata) > 0 ? each.value.metadata : null
  website_redirect    = each.value.copy_from.replace_metadata && each.value.website_redirect != "" ? each.value.website_redirect : null

  # Copy-time preconditions, evaluated by S3 against the SOURCE object; a
  # failed precondition fails the deploy (the promotion guard).
  copy_if_match            = each.value.copy_from.copy_if_match != "" ? each.value.copy_from.copy_if_match : null
  copy_if_none_match       = each.value.copy_from.copy_if_none_match != "" ? each.value.copy_from.copy_if_none_match : null
  copy_if_modified_since   = each.value.copy_from.copy_if_modified_since != "" ? each.value.copy_from.copy_if_modified_since : null
  copy_if_unmodified_since = each.value.copy_from.copy_if_unmodified_since != "" ? each.value.copy_from.copy_if_unmodified_since : null

  # The Expires header on the copied object (only the copy resource offers it).
  expires = each.value.copy_from.expires != "" ? each.value.copy_from.expires : null

  # Requester Pays acknowledgment for the SOURCE bucket.
  request_payer = each.value.copy_from.request_payer != "" ? each.value.copy_from.request_payer : null

  # Destination placement — the same surface content-sourced objects get.
  storage_class          = each.value.storage_class != "" ? each.value.storage_class : null
  server_side_encryption = each.value.server_side_encryption != "" ? each.value.server_side_encryption : null
  kms_key_id             = each.value.kms_key != "" ? each.value.kms_key : null
  bucket_key_enabled     = each.value.bucket_key_enabled
  checksum_algorithm     = each.value.checksum_algorithm != "" ? each.value.checksum_algorithm : null

  # Object Lock on the copy (requires an Object Lock-enabled destination
  # bucket; the spec CEL guarantees mode and retain-until arrive as a pair).
  object_lock_mode              = each.value.object_lock_mode != "" ? each.value.object_lock_mode : null
  object_lock_retain_until_date = each.value.object_lock_retain_until_date != "" ? each.value.object_lock_retain_until_date : null
  object_lock_legal_hold_status = each.value.object_lock_legal_hold_status != "" ? each.value.object_lock_legal_hold_status : null

  # Canned ACL — same bucket-ownership constraint as content objects.
  acl = each.value.acl != "" ? each.value.acl : null

  # The GOVERNANCE-retention delete bypass (see the content resource's note).
  force_destroy = each.value.force_destroy

  # Identity tags always reach the copy: the tag-set is sent with the REPLACE
  # tagging directive because the provider never derives the directive itself —
  # tags sent under the AWS-default COPY directive are silently ignored.
  tagging_directive = "REPLACE"
  tags              = merge(local.set_tags, try(each.value.tags, {}))
}
