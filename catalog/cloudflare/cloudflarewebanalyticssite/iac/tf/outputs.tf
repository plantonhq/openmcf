output "site_tag" {
  description = "The Cloudflare-assigned site tag (the site's identity in every RUM API path)"
  value       = cloudflare_web_analytics_site.main.site_tag
}

output "site_token" {
  description = "The site's measurement token, embedded by the JavaScript beacon (treated as a credential in outputs)"
  value       = cloudflare_web_analytics_site.main.site_token
  sensitive   = true
}

output "snippet" {
  description = "The ready-to-embed JavaScript snippet (carries the site token, so secret-marked like it)"
  value       = cloudflare_web_analytics_site.main.snippet
  sensitive   = true
}

output "ruleset_id" {
  description = "The site's ruleset ID (the parent object the include/exclude rules live under)"
  value       = cloudflare_web_analytics_site.main.ruleset.id
}
