output "gateway_id" {
  description = "The gateway's id (URL slug) -- the segment clients put in the gateway endpoint URL"
  value       = cloudflare_ai_gateway.main.id
}

output "dynamic_route_ids" {
  description = "The id of each managed dynamic route, keyed by route name"
  value       = { for name, route in cloudflare_ai_gateway_dynamic_routing.routes : name => route.id }
}
