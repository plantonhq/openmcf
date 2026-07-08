locals {
  # Enums arrive as the spec enum's FULL value names; ARM wants its own
  # casing. One shared operator vocabulary serves every condition type
  # (which subset a condition accepts is spec-enforced).
  operator_map = {
    "ANY"                   = "Any"
    "EQUAL"                 = "Equal"
    "CONTAINS"              = "Contains"
    "BEGINS_WITH"           = "BeginsWith"
    "ENDS_WITH"             = "EndsWith"
    "GREATER_THAN"          = "GreaterThan"
    "GREATER_THAN_OR_EQUAL" = "GreaterThanOrEqual"
    "LESS_THAN"             = "LessThan"
    "LESS_THAN_OR_EQUAL"    = "LessThanOrEqual"
    "REG_EX"                = "RegEx"
    "WILDCARD"              = "Wildcard"
    "GEO_MATCH"             = "GeoMatch"
    "IP_MATCH"              = "IPMatch"
  }

  transform_map = {
    "LOWERCASE"    = "Lowercase"
    "UPPERCASE"    = "Uppercase"
    "TRIM"         = "Trim"
    "URL_DECODE"   = "UrlDecode"
    "URL_ENCODE"   = "UrlEncode"
    "REMOVE_NULLS" = "RemoveNulls"
  }

  behavior_on_match_map = {
    "CONTINUE" = "Continue"
    "STOP"     = "Stop"
  }

  redirect_type_map = {
    "MOVED"              = "Moved"
    "FOUND"              = "Found"
    "TEMPORARY_REDIRECT" = "TemporaryRedirect"
    "PERMANENT_REDIRECT" = "PermanentRedirect"
  }

  # The shared forwarding-protocol enum speaks TWO ARM dialects: the
  # redirect action wants Http/Https/MatchRequest while the
  # route-configuration override wants HttpOnly/HttpsOnly/MatchRequest.
  # Same semantics; each action maps its own.
  redirect_protocol_map = {
    "MATCH_REQUEST" = "MatchRequest"
    "HTTP_ONLY"     = "Http"
    "HTTPS_ONLY"    = "Https"
  }

  override_forwarding_protocol_map = {
    "MATCH_REQUEST" = "MatchRequest"
    "HTTP_ONLY"     = "HttpOnly"
    "HTTPS_ONLY"    = "HttpsOnly"
  }

  header_action_map = {
    "APPEND"    = "Append"
    "OVERWRITE" = "Overwrite"
    "DELETE"    = "Delete"
  }

  cache_behavior_map = {
    "HONOR_ORIGIN"               = "HonorOrigin"
    "OVERRIDE_ALWAYS"            = "OverrideAlways"
    "OVERRIDE_IF_ORIGIN_MISSING" = "OverrideIfOriginMissing"
    "DISABLED"                   = "Disabled"
  }

  query_string_caching_behavior_map = {
    "IGNORE_QUERY_STRING"             = "IgnoreQueryString"
    "USE_QUERY_STRING"                = "UseQueryString"
    "IGNORE_SPECIFIED_QUERY_STRINGS"  = "IgnoreSpecifiedQueryStrings"
    "INCLUDE_SPECIFIED_QUERY_STRINGS" = "IncludeSpecifiedQueryStrings"
  }

  # Rules keyed by name for for_each -- names are spec-guaranteed unique
  # within the set (each rule is its own ARM child resource).
  rules_by_name = { for rule in var.spec.rules : rule.name => rule }

  # No Azure tags: ARM does not support tags on Front Door rule sets or
  # rules, so the platform's identity tags live on the profile.
}
