locals {
  # One alert resource per spec row. Keys carry the row index so two rows
  # may share a display name without colliding (the zone-records pattern);
  # reordering rows does churn addresses, which is harmless for the
  # seconds-fast alert objects.
  alerts = { for idx, alert in coalesce(var.spec.alerts, []) : "${idx}-${alert.alert_name}" => alert }
}
