# Enable the Compute Engine API so a fresh project can host the URL map.
# disable_on_destroy is false: tearing down one URL map must never disable the
# API for everything else in the project.
resource "google_project_service" "compute_api" {
  project = local.project_id
  service = "compute.googleapis.com"

  disable_dependent_services = true
  disable_on_destroy         = false
}

# A global Compute Engine URL map — the L7 routing brain of a global external
# Application Load Balancer. Host rules map request Host headers to named path
# matchers; path matchers evaluate route_rules (priority-ordered, rich matching)
# then path_rules (longest prefix), then their own default; anything unmatched
# falls through to the URL map's top-level default.
#
# name and project are immutable (ForceNew): changing either destroys and
# recreates the map, briefly breaking every target proxy referencing the old
# self_link. Routing tables, header actions, and tests update in place.
#
# Cross-field exclusivity and route_action coverage boundaries are enforced by
# the spec's CEL rules before deploy — this module maps only
# weighted_backend_services and url_rewrite inside route_action.
resource "google_compute_url_map" "this" {
  name        = local.url_map_name
  project     = local.project_id
  description = local.description

  default_service = local.default_service

  dynamic "default_url_redirect" {
    for_each = local.default_url_redirect != null ? [local.default_url_redirect] : []
    content {
      host_redirect          = default_url_redirect.value.host_redirect
      https_redirect         = default_url_redirect.value.https_redirect
      path_redirect          = default_url_redirect.value.path_redirect
      prefix_redirect        = default_url_redirect.value.prefix_redirect
      redirect_response_code = default_url_redirect.value.redirect_response_code
      strip_query            = default_url_redirect.value.strip_query
    }
  }

  dynamic "default_route_action" {
    for_each = local.default_route_action != null ? [local.default_route_action] : []
    content {
      dynamic "weighted_backend_services" {
        for_each = default_route_action.value.weighted_backend_services
        content {
          backend_service = weighted_backend_services.value.backend_service
          weight          = weighted_backend_services.value.weight
        }
      }

      dynamic "url_rewrite" {
        for_each = default_route_action.value.url_rewrite != null ? [default_route_action.value.url_rewrite] : []
        content {
          host_rewrite        = url_rewrite.value.host_rewrite
          path_prefix_rewrite = url_rewrite.value.path_prefix_rewrite
        }
      }
    }
  }

  dynamic "default_custom_error_response_policy" {
    for_each = local.default_custom_error_response_policy != null ? [local.default_custom_error_response_policy] : []
    content {
      error_service = default_custom_error_response_policy.value.error_service

      dynamic "error_response_rule" {
        for_each = default_custom_error_response_policy.value.error_response_rules
        content {
          match_response_codes   = error_response_rule.value.match_response_codes
          override_response_code = error_response_rule.value.override_response_code
          path                   = error_response_rule.value.path
        }
      }
    }
  }

  dynamic "header_action" {
    for_each = local.header_action != null ? [local.header_action] : []
    content {
      dynamic "request_headers_to_add" {
        for_each = header_action.value.request_headers_to_add
        content {
          header_name  = request_headers_to_add.value.header_name
          header_value = request_headers_to_add.value.header_value
          replace      = request_headers_to_add.value.replace
        }
      }
      request_headers_to_remove = header_action.value.request_headers_to_remove

      dynamic "response_headers_to_add" {
        for_each = header_action.value.response_headers_to_add
        content {
          header_name  = response_headers_to_add.value.header_name
          header_value = response_headers_to_add.value.header_value
          replace      = response_headers_to_add.value.replace
        }
      }
      response_headers_to_remove = header_action.value.response_headers_to_remove
    }
  }

  dynamic "host_rule" {
    for_each = local.host_rules
    content {
      hosts        = host_rule.value.hosts
      path_matcher = host_rule.value.path_matcher
      description  = host_rule.value.description
    }
  }

  dynamic "path_matcher" {
    for_each = local.path_matchers
    content {
      name            = path_matcher.value.name
      description     = path_matcher.value.description
      default_service = path_matcher.value.default_service

      dynamic "default_url_redirect" {
        for_each = path_matcher.value.default_url_redirect != null ? [path_matcher.value.default_url_redirect] : []
        content {
          host_redirect          = default_url_redirect.value.host_redirect
          https_redirect         = default_url_redirect.value.https_redirect
          path_redirect          = default_url_redirect.value.path_redirect
          prefix_redirect        = default_url_redirect.value.prefix_redirect
          redirect_response_code = default_url_redirect.value.redirect_response_code
          strip_query            = default_url_redirect.value.strip_query
        }
      }

      dynamic "default_route_action" {
        for_each = path_matcher.value.default_route_action != null ? [path_matcher.value.default_route_action] : []
        content {
          dynamic "weighted_backend_services" {
            for_each = default_route_action.value.weighted_backend_services
            content {
              backend_service = weighted_backend_services.value.backend_service
              weight          = weighted_backend_services.value.weight
            }
          }

          dynamic "url_rewrite" {
            for_each = default_route_action.value.url_rewrite != null ? [default_route_action.value.url_rewrite] : []
            content {
              host_rewrite        = url_rewrite.value.host_rewrite
              path_prefix_rewrite = url_rewrite.value.path_prefix_rewrite
            }
          }
        }
      }

      dynamic "default_custom_error_response_policy" {
        for_each = path_matcher.value.default_custom_error_response_policy != null ? [path_matcher.value.default_custom_error_response_policy] : []
        content {
          error_service = default_custom_error_response_policy.value.error_service

          dynamic "error_response_rule" {
            for_each = default_custom_error_response_policy.value.error_response_rules
            content {
              match_response_codes   = error_response_rule.value.match_response_codes
              override_response_code = error_response_rule.value.override_response_code
              path                   = error_response_rule.value.path
            }
          }
        }
      }

      dynamic "header_action" {
        for_each = path_matcher.value.header_action != null ? [path_matcher.value.header_action] : []
        content {
          dynamic "request_headers_to_add" {
            for_each = header_action.value.request_headers_to_add
            content {
              header_name  = request_headers_to_add.value.header_name
              header_value = request_headers_to_add.value.header_value
              replace      = request_headers_to_add.value.replace
            }
          }
          request_headers_to_remove = header_action.value.request_headers_to_remove

          dynamic "response_headers_to_add" {
            for_each = header_action.value.response_headers_to_add
            content {
              header_name  = response_headers_to_add.value.header_name
              header_value = response_headers_to_add.value.header_value
              replace      = response_headers_to_add.value.replace
            }
          }
          response_headers_to_remove = header_action.value.response_headers_to_remove
        }
      }

      dynamic "path_rule" {
        for_each = path_matcher.value.path_rules
        content {
          paths   = path_rule.value.paths
          service = path_rule.value.service

          dynamic "url_redirect" {
            for_each = path_rule.value.url_redirect != null ? [path_rule.value.url_redirect] : []
            content {
              host_redirect          = url_redirect.value.host_redirect
              https_redirect         = url_redirect.value.https_redirect
              path_redirect          = url_redirect.value.path_redirect
              prefix_redirect        = url_redirect.value.prefix_redirect
              redirect_response_code = url_redirect.value.redirect_response_code
              strip_query            = url_redirect.value.strip_query
            }
          }

          dynamic "route_action" {
            for_each = path_rule.value.route_action != null ? [path_rule.value.route_action] : []
            content {
              dynamic "weighted_backend_services" {
                for_each = route_action.value.weighted_backend_services
                content {
                  backend_service = weighted_backend_services.value.backend_service
                  weight          = weighted_backend_services.value.weight
                }
              }

              dynamic "url_rewrite" {
                for_each = route_action.value.url_rewrite != null ? [route_action.value.url_rewrite] : []
                content {
                  host_rewrite        = url_rewrite.value.host_rewrite
                  path_prefix_rewrite = url_rewrite.value.path_prefix_rewrite
                }
              }
            }
          }

          dynamic "custom_error_response_policy" {
            for_each = path_rule.value.custom_error_response_policy != null ? [path_rule.value.custom_error_response_policy] : []
            content {
              error_service = custom_error_response_policy.value.error_service

              dynamic "error_response_rule" {
                for_each = custom_error_response_policy.value.error_response_rules
                content {
                  match_response_codes   = error_response_rule.value.match_response_codes
                  override_response_code = error_response_rule.value.override_response_code
                  path                   = error_response_rule.value.path
                }
              }
            }
          }
        }
      }

      dynamic "route_rules" {
        for_each = path_matcher.value.route_rules
        content {
          priority = route_rules.value.priority
          service  = route_rules.value.service

          dynamic "url_redirect" {
            for_each = route_rules.value.url_redirect != null ? [route_rules.value.url_redirect] : []
            content {
              host_redirect          = url_redirect.value.host_redirect
              https_redirect         = url_redirect.value.https_redirect
              path_redirect          = url_redirect.value.path_redirect
              prefix_redirect        = url_redirect.value.prefix_redirect
              redirect_response_code = url_redirect.value.redirect_response_code
              strip_query            = url_redirect.value.strip_query
            }
          }

          dynamic "route_action" {
            for_each = route_rules.value.route_action != null ? [route_rules.value.route_action] : []
            content {
              dynamic "weighted_backend_services" {
                for_each = route_action.value.weighted_backend_services
                content {
                  backend_service = weighted_backend_services.value.backend_service
                  weight          = weighted_backend_services.value.weight
                }
              }

              dynamic "url_rewrite" {
                for_each = route_action.value.url_rewrite != null ? [route_action.value.url_rewrite] : []
                content {
                  host_rewrite          = url_rewrite.value.host_rewrite
                  path_prefix_rewrite   = url_rewrite.value.path_prefix_rewrite
                  path_template_rewrite = url_rewrite.value.path_template_rewrite
                }
              }
            }
          }

          dynamic "header_action" {
            for_each = route_rules.value.header_action != null ? [route_rules.value.header_action] : []
            content {
              dynamic "request_headers_to_add" {
                for_each = header_action.value.request_headers_to_add
                content {
                  header_name  = request_headers_to_add.value.header_name
                  header_value = request_headers_to_add.value.header_value
                  replace      = request_headers_to_add.value.replace
                }
              }
              request_headers_to_remove = header_action.value.request_headers_to_remove

              dynamic "response_headers_to_add" {
                for_each = header_action.value.response_headers_to_add
                content {
                  header_name  = response_headers_to_add.value.header_name
                  header_value = response_headers_to_add.value.header_value
                  replace      = response_headers_to_add.value.replace
                }
              }
              response_headers_to_remove = header_action.value.response_headers_to_remove
            }
          }

          dynamic "custom_error_response_policy" {
            for_each = route_rules.value.custom_error_response_policy != null ? [route_rules.value.custom_error_response_policy] : []
            content {
              error_service = custom_error_response_policy.value.error_service

              dynamic "error_response_rule" {
                for_each = custom_error_response_policy.value.error_response_rules
                content {
                  match_response_codes   = error_response_rule.value.match_response_codes
                  override_response_code = error_response_rule.value.override_response_code
                  path                   = error_response_rule.value.path
                }
              }
            }
          }

          dynamic "match_rules" {
            for_each = route_rules.value.match_rules
            content {
              prefix_match        = match_rules.value.prefix_match
              full_path_match     = match_rules.value.full_path_match
              regex_match         = match_rules.value.regex_match
              path_template_match = match_rules.value.path_template_match
              ignore_case         = match_rules.value.ignore_case

              dynamic "header_matches" {
                for_each = match_rules.value.header_matches
                content {
                  header_name   = header_matches.value.header_name
                  exact_match   = header_matches.value.exact_match
                  prefix_match  = header_matches.value.prefix_match
                  suffix_match  = header_matches.value.suffix_match
                  regex_match   = header_matches.value.regex_match
                  present_match = header_matches.value.present_match
                  invert_match  = header_matches.value.invert_match

                  dynamic "range_match" {
                    for_each = header_matches.value.range_match != null ? [header_matches.value.range_match] : []
                    content {
                      range_start = range_match.value.range_start
                      range_end   = range_match.value.range_end
                    }
                  }
                }
              }

              dynamic "query_parameter_matches" {
                for_each = match_rules.value.query_parameter_matches
                content {
                  name          = query_parameter_matches.value.name
                  exact_match   = query_parameter_matches.value.exact_match
                  present_match = query_parameter_matches.value.present_match
                  regex_match   = query_parameter_matches.value.regex_match
                }
              }

              dynamic "metadata_filters" {
                for_each = match_rules.value.metadata_filters
                content {
                  filter_match_criteria = metadata_filters.value.filter_match_criteria

                  dynamic "filter_labels" {
                    for_each = metadata_filters.value.filter_labels
                    content {
                      name  = filter_labels.value.name
                      value = filter_labels.value.value
                    }
                  }
                }
              }
            }
          }
        }
      }
    }
  }

  dynamic "test" {
    for_each = local.tests
    content {
      host                            = test.value.host
      path                            = test.value.path
      service                         = test.value.service
      description                     = test.value.description
      expected_output_url             = test.value.expected_output_url
      expected_redirect_response_code = test.value.expected_redirect_response_code

      dynamic "headers" {
        for_each = test.value.headers
        content {
          name  = headers.value.name
          value = headers.value.value
        }
      }
    }
  }

  depends_on = [google_project_service.compute_api]
}
