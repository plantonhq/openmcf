# Stack outputs — exactly the DigitalOceanUptimeCheckStackOutputs contract,
# identical across both provisioners. Alert row ids are not outputs; the
# composed alert rows import as "{check_id},{alert_id}" with the alert id
# found via the API or console.

output "check_id" {
  description = "UUID of the uptime check (the API identity, and the import id)"
  value       = digitalocean_uptime_check.check.id
}
