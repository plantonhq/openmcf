variable "metadata" {
  description = "Metadata for the resource, including name and labels"
  type = object({
    name    = string,
    id      = optional(string),
    org     = optional(string),
    env     = optional(string),
    labels  = optional(map(string)),
    tags    = optional(list(string)),
    version = optional(object({ id = string, message = string }))
  })
}

variable "spec" {
  description = "Specification for the GCP Compute Engine global URL map"
  type = object({
    # The GCP project that owns the URL map. The CLI's tfvars converter
    # resolves StringValueOrRef fields to their literal string before the
    # module runs, so this arrives as a plain string.
    # If empty, the provider's default project is used (see locals.tf).
    project_id = optional(string, "")

    # Name of the URL map in GCP (RFC1035). Empty defaults to metadata.name
    # (see locals.tf). Immutable (ForceNew).
    url_map_name = optional(string, "")

    description = optional(string, "")

    # Exactly one default target is required (enforced by the spec's CEL
    # upstream). StringValueOrRef fields arrive as plain strings.
    default_service = optional(string, "")

    default_url_redirect = optional(object({
      host_redirect          = optional(string, "")
      https_redirect         = optional(bool, false)
      path_redirect          = optional(string, "")
      prefix_redirect        = optional(string, "")
      redirect_response_code = optional(string, "")
      strip_query            = optional(bool, false)
    }))

    default_route_action = optional(object({
      weighted_backend_services = optional(list(object({
        backend_service = string
        weight          = optional(number, 0)
      })), [])
      url_rewrite = optional(object({
        host_rewrite        = optional(string, "")
        path_prefix_rewrite = optional(string, "")
      }))
    }))

    default_custom_error_response_policy = optional(object({
      error_service = optional(string, "")
      error_response_rules = optional(list(object({
        match_response_codes   = list(string)
        override_response_code = optional(number)
        path                   = string
      })), [])
    }))

    header_action = optional(object({
      request_headers_to_add = optional(list(object({
        header_name  = string
        header_value = string
        replace      = optional(bool, false)
      })), [])
      request_headers_to_remove = optional(list(string), [])
      response_headers_to_add = optional(list(object({
        header_name  = string
        header_value = string
        replace      = optional(bool, false)
      })), [])
      response_headers_to_remove = optional(list(string), [])
    }))

    host_rules = optional(list(object({
      hosts        = list(string)
      path_matcher = string
      description  = optional(string, "")
    })), [])

    path_matchers = optional(list(object({
      name            = string
      description     = optional(string, "")
      default_service = optional(string, "")
      default_url_redirect = optional(object({
        host_redirect          = optional(string, "")
        https_redirect         = optional(bool, false)
        path_redirect          = optional(string, "")
        prefix_redirect        = optional(string, "")
        redirect_response_code = optional(string, "")
        strip_query            = optional(bool, false)
      }))
      default_route_action = optional(object({
        weighted_backend_services = optional(list(object({
          backend_service = string
          weight          = optional(number, 0)
        })), [])
        url_rewrite = optional(object({
          host_rewrite        = optional(string, "")
          path_prefix_rewrite = optional(string, "")
        }))
      }))
      default_custom_error_response_policy = optional(object({
        error_service = optional(string, "")
        error_response_rules = optional(list(object({
          match_response_codes   = list(string)
          override_response_code = optional(number)
          path                   = string
        })), [])
      }))
      header_action = optional(object({
        request_headers_to_add = optional(list(object({
          header_name  = string
          header_value = string
          replace      = optional(bool, false)
        })), [])
        request_headers_to_remove = optional(list(string), [])
        response_headers_to_add = optional(list(object({
          header_name  = string
          header_value = string
          replace      = optional(bool, false)
        })), [])
        response_headers_to_remove = optional(list(string), [])
      }))
      path_rules = optional(list(object({
        paths   = list(string)
        service = optional(string, "")
        url_redirect = optional(object({
          host_redirect          = optional(string, "")
          https_redirect         = optional(bool, false)
          path_redirect          = optional(string, "")
          prefix_redirect        = optional(string, "")
          redirect_response_code = optional(string, "")
          strip_query            = optional(bool, false)
        }))
        route_action = optional(object({
          weighted_backend_services = optional(list(object({
            backend_service = string
            weight          = optional(number, 0)
          })), [])
          url_rewrite = optional(object({
            host_rewrite        = optional(string, "")
            path_prefix_rewrite = optional(string, "")
          }))
        }))
        custom_error_response_policy = optional(object({
          error_service = optional(string, "")
          error_response_rules = optional(list(object({
            match_response_codes   = list(string)
            override_response_code = optional(number)
            path                   = string
          })), [])
        }))
      })), [])
      route_rules = optional(list(object({
        priority = number
        service  = optional(string, "")
        url_redirect = optional(object({
          host_redirect          = optional(string, "")
          https_redirect         = optional(bool, false)
          path_redirect          = optional(string, "")
          prefix_redirect        = optional(string, "")
          redirect_response_code = optional(string, "")
          strip_query            = optional(bool, false)
        }))
        route_action = optional(object({
          weighted_backend_services = optional(list(object({
            backend_service = string
            weight          = optional(number, 0)
          })), [])
          url_rewrite = optional(object({
            host_rewrite          = optional(string, "")
            path_prefix_rewrite   = optional(string, "")
            path_template_rewrite = optional(string, "")
          }))
        }))
        header_action = optional(object({
          request_headers_to_add = optional(list(object({
            header_name  = string
            header_value = string
            replace      = optional(bool, false)
          })), [])
          request_headers_to_remove = optional(list(string), [])
          response_headers_to_add = optional(list(object({
            header_name  = string
            header_value = string
            replace      = optional(bool, false)
          })), [])
          response_headers_to_remove = optional(list(string), [])
        }))
        custom_error_response_policy = optional(object({
          error_service = optional(string, "")
          error_response_rules = optional(list(object({
            match_response_codes   = list(string)
            override_response_code = optional(number)
            path                   = string
          })), [])
        }))
        match_rules = list(object({
          prefix_match        = optional(string, "")
          full_path_match     = optional(string, "")
          regex_match         = optional(string, "")
          path_template_match = optional(string, "")
          ignore_case         = optional(bool, false)
          header_matches = optional(list(object({
            header_name   = string
            exact_match   = optional(string, "")
            prefix_match  = optional(string, "")
            suffix_match  = optional(string, "")
            regex_match   = optional(string, "")
            present_match = optional(bool, false)
            invert_match  = optional(bool, false)
            range_match = optional(object({
              range_start = optional(number, 0)
              range_end   = optional(number, 0)
            }))
          })), [])
          query_parameter_matches = optional(list(object({
            name          = string
            exact_match   = optional(string, "")
            present_match = optional(bool, false)
            regex_match   = optional(string, "")
          })), [])
          metadata_filters = optional(list(object({
            filter_match_criteria = string
            filter_labels = list(object({
              name  = string
              value = string
            }))
          })), [])
        }))
      })), [])
    })), [])

    tests = optional(list(object({
      host                            = string
      path                            = string
      service                         = optional(string, "")
      description                     = optional(string, "")
      expected_output_url             = optional(string, "")
      expected_redirect_response_code = optional(number)
      headers = optional(list(object({
        name  = string
        value = string
      })), [])
    })), [])
  })
}
