locals {
  # Health-probe enums arrive as the spec enum's FULL value names; ARM
  # wants its own casing. request_type absent means HEAD (tfvars drops
  # zero-valued proto fields; the module materializes the documented
  # default).
  health_probe_protocol_map = {
    "HTTP"  = "Http"
    "HTTPS" = "Https"
  }
  health_probe_request_type_map = {
    "HEAD" = "HEAD"
    "GET"  = "GET"
  }

  # No Azure tags: ARM does not support tags on Front Door origin
  # groups, so the platform's identity tags live on the profile.
}
