# Cloudflare waiting room event: a scheduled window during which the event's
# override values temporarily replace the room's settings. Every override is
# null-means-inherit -- unset fields are never sent, so the room's value stays
# in charge during the window. Time ordering (start >= 1 min before end,
# prequeue >= 5 min before start) is validated by the spec's CEL up front; the
# API enforces the same rules with opaque errors.
resource "cloudflare_waiting_room_event" "main" {
  zone_id         = var.spec.zone_id
  waiting_room_id = var.spec.waiting_room_id

  name             = var.spec.name
  event_start_time = var.spec.event_start_time
  event_end_time   = var.spec.event_end_time

  prequeue_start_time    = var.spec.prequeue_start_time != "" ? var.spec.prequeue_start_time : null
  shuffle_at_event_start = var.spec.shuffle_at_event_start
  description            = var.spec.description != "" ? var.spec.description : null
  suspended              = var.spec.suspended

  custom_page_html        = var.spec.custom_page_html != "" ? var.spec.custom_page_html : null
  disable_session_renewal = var.spec.disable_session_renewal
  new_users_per_minute    = var.spec.new_users_per_minute
  total_active_users      = var.spec.total_active_users
  queueing_method         = var.spec.queueing_method
  session_duration        = var.spec.session_duration
  turnstile_action        = var.spec.turnstile_action
  turnstile_mode          = var.spec.turnstile_mode
}
