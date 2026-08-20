# A device posture rule: a health check WARP evaluates on enrolled devices,
# which Access and Gateway policies can then require. A plain CRUD resource
# (real create/update/delete; only the account forces replacement).
#
# The input tree carries every check family's parameters; unset fields are
# never sent, so each rule's payload holds exactly the fields its type reads.
resource "cloudflare_zero_trust_device_posture_rule" "main" {
  account_id = var.spec.account_id
  name       = var.spec.name
  type       = var.spec.type

  description = try(var.spec.description, "") != "" ? var.spec.description : null
  expiration  = try(var.spec.expiration, "") != "" ? var.spec.expiration : null
  schedule    = try(var.spec.schedule, "") != "" ? var.spec.schedule : null

  match = length(try(var.spec.match, [])) > 0 ? [
    for row in var.spec.match : { platform = row.platform }
  ] : null

  input = try(var.spec.input, null) != null ? {
    operating_system   = try(var.spec.input.operating_system, "") != "" ? var.spec.input.operating_system : null
    path               = try(var.spec.input.path, "") != "" ? var.spec.input.path : null
    exists             = try(var.spec.input.exists, null)
    sha256             = try(var.spec.input.sha256, "") != "" ? var.spec.input.sha256 : null
    thumbprint         = try(var.spec.input.thumbprint, "") != "" ? var.spec.input.thumbprint : null
    id                 = try(var.spec.input.id, "") != "" ? var.spec.input.id : null
    domain             = try(var.spec.input.domain, "") != "" ? var.spec.input.domain : null
    operator           = try(var.spec.input.operator, "") != "" ? var.spec.input.operator : null
    version            = try(var.spec.input.version, "") != "" ? var.spec.input.version : null
    os_distro_name     = try(var.spec.input.os_distro_name, "") != "" ? var.spec.input.os_distro_name : null
    os_distro_revision = try(var.spec.input.os_distro_revision, "") != "" ? var.spec.input.os_distro_revision : null
    os_version_extra   = try(var.spec.input.os_version_extra, "") != "" ? var.spec.input.os_version_extra : null
    enabled            = try(var.spec.input.enabled, null)
    check_disks        = length(try(var.spec.input.check_disks, [])) > 0 ? var.spec.input.check_disks : null
    require_all        = try(var.spec.input.require_all, null)
    certificate_id     = try(var.spec.input.certificate_id, "") != "" ? var.spec.input.certificate_id : null
    cn                 = try(var.spec.input.cn, "") != "" ? var.spec.input.cn : null
    check_private_key  = try(var.spec.input.check_private_key, null)
    extended_key_usage = length(try(var.spec.input.extended_key_usage, [])) > 0 ? var.spec.input.extended_key_usage : null

    locations = try(var.spec.input.locations, null) != null ? {
      paths        = length(try(var.spec.input.locations.paths, [])) > 0 ? var.spec.input.locations.paths : null
      trust_stores = length(try(var.spec.input.locations.trust_stores, [])) > 0 ? var.spec.input.locations.trust_stores : null
    } : null

    subject_alternative_names = length(try(var.spec.input.subject_alternative_names, [])) > 0 ? var.spec.input.subject_alternative_names : null
    update_window_days        = try(var.spec.input.update_window_days, null)
    compliance_status         = try(var.spec.input.compliance_status, "") != "" ? var.spec.input.compliance_status : null
    connection_id             = try(var.spec.input.connection_id, "") != "" ? var.spec.input.connection_id : null
    last_seen                 = try(var.spec.input.last_seen, "") != "" ? var.spec.input.last_seen : null
    os                        = try(var.spec.input.os, "") != "" ? var.spec.input.os : null
    overall                   = try(var.spec.input.overall, "") != "" ? var.spec.input.overall : null
    sensor_config             = try(var.spec.input.sensor_config, "") != "" ? var.spec.input.sensor_config : null
    state                     = try(var.spec.input.state, "") != "" ? var.spec.input.state : null
    version_operator          = try(var.spec.input.version_operator, "") != "" ? var.spec.input.version_operator : null
    auth_state                = length(try(var.spec.input.auth_state, [])) > 0 ? var.spec.input.auth_state : null
    count_operator            = try(var.spec.input.count_operator, "") != "" ? var.spec.input.count_operator : null
    issue_count               = try(var.spec.input.issue_count, "") != "" ? var.spec.input.issue_count : null
    eid_last_seen             = try(var.spec.input.eid_last_seen, "") != "" ? var.spec.input.eid_last_seen : null
    risk_level                = try(var.spec.input.risk_level, "") != "" ? var.spec.input.risk_level : null
    score_operator            = try(var.spec.input.score_operator, "") != "" ? var.spec.input.score_operator : null
    total_score               = try(var.spec.input.total_score, null)
    active_threats            = try(var.spec.input.active_threats, null)
    infected                  = try(var.spec.input.infected, null)
    is_active                 = try(var.spec.input.is_active, null)
    network_status            = try(var.spec.input.network_status, "") != "" ? var.spec.input.network_status : null
    operational_state         = try(var.spec.input.operational_state, "") != "" ? var.spec.input.operational_state : null
    score                     = try(var.spec.input.score, null)
  } : null
}
