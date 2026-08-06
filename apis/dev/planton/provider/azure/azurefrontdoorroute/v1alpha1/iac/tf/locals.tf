locals {
  # Enums arrive as the spec enum's FULL value names; ARM wants its own
  # casing. Absent enum fields materialize the spec's documented
  # defaults here (tfvars drops zero-valued proto fields).
  protocol_map = {
    "HTTP"  = "Http"
    "HTTPS" = "Https"
  }
  supported_protocols = [for p in var.spec.supported_protocols : local.protocol_map[p]]

  forwarding_protocol_map = {
    "MATCH_REQUEST" = "MatchRequest"
    "HTTP_ONLY"     = "HttpOnly"
    "HTTPS_ONLY"    = "HttpsOnly"
  }
  forwarding_protocol = local.forwarding_protocol_map[coalesce(var.spec.forwarding_protocol, "MATCH_REQUEST")]

  query_string_caching_behavior_map = {
    "IGNORE_QUERY_STRING"             = "IgnoreQueryString"
    "USE_QUERY_STRING"                = "UseQueryString"
    "IGNORE_SPECIFIED_QUERY_STRINGS"  = "IgnoreSpecifiedQueryStrings"
    "INCLUDE_SPECIFIED_QUERY_STRINGS" = "IncludeSpecifiedQueryStrings"
  }

  # No Azure tags: ARM does not support tags on Front Door routes, so
  # the platform's identity tags live on the profile.
}
