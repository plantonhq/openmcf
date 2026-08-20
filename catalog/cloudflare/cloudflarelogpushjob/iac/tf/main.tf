# A Logpush job: continuous delivery of one Cloudflare log dataset to a
# destination the account controls. Dual scope -- exactly one of account_id
# or zone_id is set (spec validation enforces it).
#
# Two API truths this module honors:
#   - `dataset` is immutable (the provider replaces the job when it changes).
#   - `destination_conf` is NOT marked for replacement by the provider, yet
#     Cloudflare rejects changing it on an existing job (HTTP 400 at apply).
#     To repoint a job, delete and recreate it.
#
# `enabled` defaults to TRUE here even though Cloudflare's own default is
# FALSE: a declared log job is meant to ship logs. Set enabled = false
# explicitly to pause delivery.
resource "cloudflare_logpush_job" "main" {
  account_id = try(var.spec.account_id, "") != "" ? var.spec.account_id : null
  zone_id    = try(var.spec.zone_id, "") != "" ? var.spec.zone_id : null

  dataset          = var.spec.dataset
  destination_conf = var.spec.destination_conf

  name    = try(var.spec.name, "") != "" ? var.spec.name : null
  enabled = try(var.spec.enabled, null) != null ? var.spec.enabled : true
  filter  = try(var.spec.filter, "") != "" ? var.spec.filter : null
  kind    = try(var.spec.kind, "") != "" ? var.spec.kind : null

  max_upload_bytes            = try(var.spec.max_upload_bytes, null)
  max_upload_interval_seconds = try(var.spec.max_upload_interval_seconds, null)
  max_upload_records          = try(var.spec.max_upload_records, null)

  ownership_challenge = try(var.spec.ownership_challenge, "") != "" ? var.spec.ownership_challenge : null

  output_options = try(var.spec.output_options, null) != null ? {
    output_type      = try(var.spec.output_options.output_type, "") != "" ? var.spec.output_options.output_type : null
    field_names      = length(try(var.spec.output_options.field_names, [])) > 0 ? var.spec.output_options.field_names : null
    timestamp_format = try(var.spec.output_options.timestamp_format, "") != "" ? var.spec.output_options.timestamp_format : null
    sample_rate      = try(var.spec.output_options.sample_rate, null)
    batch_prefix     = try(var.spec.output_options.batch_prefix, "") != "" ? var.spec.output_options.batch_prefix : null
    batch_suffix     = try(var.spec.output_options.batch_suffix, "") != "" ? var.spec.output_options.batch_suffix : null
    record_prefix    = try(var.spec.output_options.record_prefix, "") != "" ? var.spec.output_options.record_prefix : null
    record_suffix    = try(var.spec.output_options.record_suffix, "") != "" ? var.spec.output_options.record_suffix : null
    record_delimiter = try(var.spec.output_options.record_delimiter, "") != "" ? var.spec.output_options.record_delimiter : null
    record_template  = try(var.spec.output_options.record_template, "") != "" ? var.spec.output_options.record_template : null
    field_delimiter  = try(var.spec.output_options.field_delimiter, "") != "" ? var.spec.output_options.field_delimiter : null
    merge_subrequests = try(var.spec.output_options.merge_subrequests, null)
    cve_2021_44228    = try(var.spec.output_options.cve_2021_44228, null)
  } : null
}

# The ownership-challenge issuing step, deployed only when the spec asks for
# it. ONE-SHOT at Cloudflare: the POST drops a challenge file into the
# destination and that is the whole lifecycle -- no read, no update, no
# delete (destroying this resource only forgets it), no import. The token
# inside the dropped file is fetched by the operator and fed back as
# spec.ownership_challenge -- Cloudflare deliberately routes the proof
# through storage you control, so no API can read it for you.
resource "cloudflare_logpush_ownership_challenge" "main" {
  count = var.spec.generate_ownership_challenge ? 1 : 0

  account_id = try(var.spec.account_id, "") != "" ? var.spec.account_id : null
  zone_id    = try(var.spec.zone_id, "") != "" ? var.spec.zone_id : null

  destination_conf = var.spec.destination_conf
}
