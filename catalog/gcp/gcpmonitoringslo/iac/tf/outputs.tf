# Both outputs derive from the SLO's server-assigned resource name
# (projects/{p}/services/{s}/serviceLevelObjectives/{id}) so they are
# correct on every service arm — including an existing service in the
# provider's ambient project, where no service resource exists in this
# module to read a name from.
output "slo_name" {
  description = "Resource name of the SLO (projects/{p}/services/{s}/serviceLevelObjectives/{id}) — the burn-rate alert handle"
  value       = google_monitoring_slo.this.name
}

output "service_name" {
  description = "Resource name of the Monitoring service the SLO measures (projects/{p}/services/{s})"
  value       = split("/serviceLevelObjectives/", google_monitoring_slo.this.name)[0]
}
