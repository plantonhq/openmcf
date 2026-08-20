output "job_id" {
  description = "The Cloudflare-assigned numeric job ID, in string form"
  value       = tostring(cloudflare_logpush_job.main.id)
}

output "account_id" {
  description = "The account the job lives in (account-scoped jobs)"
  value       = try(var.spec.account_id, "")
}

output "zone_id" {
  description = "The zone the job lives in (zone-scoped jobs)"
  value       = try(var.spec.zone_id, "")
}

output "ownership_challenge_filename" {
  description = "Name of the challenge file Cloudflare dropped into the destination (when the challenge arm ran); fetch its content and set spec.ownership_challenge"
  value       = var.spec.generate_ownership_challenge ? cloudflare_logpush_ownership_challenge.main[0].filename : ""
}

output "ownership_challenge_message" {
  description = "Message accompanying the issued challenge, when Cloudflare returns one"
  value       = var.spec.generate_ownership_challenge ? cloudflare_logpush_ownership_challenge.main[0].message : ""
}

output "ownership_challenge_valid" {
  description = "Whether Cloudflare reported the destination valid when issuing the challenge"
  value       = var.spec.generate_ownership_challenge ? cloudflare_logpush_ownership_challenge.main[0].valid : false
}
