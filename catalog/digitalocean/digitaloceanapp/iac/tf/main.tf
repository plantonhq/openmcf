resource "digitalocean_app" "main" {
  project_id = var.spec.project_id != "" ? var.spec.project_id : null

  spec {
    name                          = var.spec.app_name
    region                        = local.region
    disable_edge_cache            = var.spec.disable_edge_cache
    disable_email_obfuscation     = var.spec.disable_email_obfuscation
    enhanced_threat_control_enabled = var.spec.enhanced_threat_control_enabled
    features                      = length(var.spec.features) > 0 ? var.spec.features : null

    dynamic "env" {
      for_each = var.spec.envs
      content {
        key   = env.value.key
        value = env.value.secret != "" ? env.value.secret : env.value.plaintext
        type  = env.value.secret != "" ? "SECRET" : "GENERAL"
        scope = (env.value.scope == "" || endswith(env.value.scope, "_unspecified")) ? "RUN_AND_BUILD_TIME" : upper(env.value.scope)
      }
    }

    dynamic "service" {
      for_each = var.spec.services
      content {
        name               = service.value.name
        source_dir         = service.value.source_dir != "" ? service.value.source_dir : null
        environment_slug   = service.value.environment_slug != "" ? service.value.environment_slug : null
        dockerfile_path    = service.value.dockerfile_path != "" ? service.value.dockerfile_path : null
        build_command      = service.value.build_command != "" ? service.value.build_command : null
        run_command        = service.value.run_command != "" ? service.value.run_command : null
        instance_size_slug = service.value.instance_size_slug != "" ? service.value.instance_size_slug : null
        instance_count     = service.value.autoscaling != null ? null : (service.value.instance_count > 0 ? service.value.instance_count : null)
        http_port          = service.value.http_port
        internal_ports     = length(service.value.internal_ports) > 0 ? service.value.internal_ports : null

        dynamic "git" {
          for_each = service.value.git != null ? [service.value.git] : []
          content {
            repo_clone_url = git.value.repo_clone_url
            branch         = git.value.branch
          }
        }
        dynamic "github" {
          for_each = service.value.github != null ? [service.value.github] : []
          content {
            repo           = github.value.repo
            branch         = github.value.branch
            deploy_on_push = github.value.deploy_on_push
          }
        }
        dynamic "gitlab" {
          for_each = service.value.gitlab != null ? [service.value.gitlab] : []
          content {
            repo           = gitlab.value.repo
            branch         = gitlab.value.branch
            deploy_on_push = gitlab.value.deploy_on_push
          }
        }
        dynamic "bitbucket" {
          for_each = service.value.bitbucket != null ? [service.value.bitbucket] : []
          content {
            repo           = bitbucket.value.repo
            branch         = bitbucket.value.branch
            deploy_on_push = bitbucket.value.deploy_on_push
          }
        }
        dynamic "image" {
          for_each = service.value.image != null ? [service.value.image] : []
          content {
            registry_type         = upper(image.value.registry_type)
            registry              = image.value.registry != "" ? image.value.registry : null
            repository            = image.value.repository
            tag                   = image.value.tag != "" ? image.value.tag : null
            digest                = image.value.digest != "" ? image.value.digest : null
            registry_credentials  = image.value.registry_credentials != "" ? image.value.registry_credentials : null
            dynamic "deploy_on_push" {
              for_each = image.value.deploy_on_push ? [1] : []
              content {
                enabled = true
              }
            }
          }
        }
        dynamic "health_check" {
          for_each = service.value.health_check != null ? [service.value.health_check] : []
          content {
            port                  = health_check.value.port
            http_path             = health_check.value.http_path != "" ? health_check.value.http_path : null
            initial_delay_seconds = health_check.value.initial_delay_seconds
            period_seconds        = health_check.value.period_seconds
            timeout_seconds       = health_check.value.timeout_seconds
            success_threshold     = health_check.value.success_threshold
            failure_threshold     = health_check.value.failure_threshold
          }
        }
        dynamic "liveness_health_check" {
          for_each = service.value.liveness_health_check != null ? [service.value.liveness_health_check] : []
          content {
            port                  = liveness_health_check.value.port
            http_path             = liveness_health_check.value.http_path != "" ? liveness_health_check.value.http_path : null
            initial_delay_seconds = liveness_health_check.value.initial_delay_seconds
            period_seconds        = liveness_health_check.value.period_seconds
            timeout_seconds       = liveness_health_check.value.timeout_seconds
            success_threshold     = liveness_health_check.value.success_threshold
            failure_threshold     = liveness_health_check.value.failure_threshold
          }
        }
        dynamic "autoscaling" {
          for_each = service.value.autoscaling != null ? [service.value.autoscaling] : []
          content {
            min_instance_count = autoscaling.value.min_instance_count
            max_instance_count = autoscaling.value.max_instance_count
            metrics {
              cpu {
                percent = autoscaling.value.cpu_percent
              }
            }
          }
        }
        dynamic "termination" {
          for_each = service.value.termination != null ? [service.value.termination] : []
          content {
            grace_period_seconds = termination.value.grace_period_seconds
            drain_seconds        = termination.value.drain_seconds
          }
        }
        dynamic "env" {
          for_each = service.value.envs
          content {
            key   = env.value.key
            value = env.value.secret != "" ? env.value.secret : env.value.plaintext
            type  = env.value.secret != "" ? "SECRET" : "GENERAL"
            scope = (env.value.scope == "" || endswith(env.value.scope, "_unspecified")) ? "RUN_AND_BUILD_TIME" : upper(env.value.scope)
          }
        }
        dynamic "alert" {
          for_each = service.value.alerts
          content {
            rule     = upper(alert.value.rule)
            operator = upper(alert.value.operator)
            window   = upper(alert.value.window)
            value    = alert.value.value
            disabled = alert.value.disabled
            dynamic "destinations" {
              for_each = alert.value.destinations != null ? [alert.value.destinations] : []
              content {
                emails = length(destinations.value.emails) > 0 ? destinations.value.emails : null
                dynamic "slack_webhooks" {
                  for_each = destinations.value.slack_webhooks
                  content {
                    channel = slack_webhooks.value.channel
                    url     = slack_webhooks.value.url
                  }
                }
              }
            }
          }
        }
        dynamic "log_destination" {
          for_each = service.value.log_destinations
          content {
            name = log_destination.value.name
            dynamic "papertrail" {
              for_each = log_destination.value.papertrail != null ? [log_destination.value.papertrail] : []
              content { endpoint = papertrail.value.endpoint }
            }
            dynamic "datadog" {
              for_each = log_destination.value.datadog != null ? [log_destination.value.datadog] : []
              content {
                api_key  = datadog.value.api_key
                endpoint = datadog.value.endpoint != "" ? datadog.value.endpoint : null
              }
            }
            dynamic "logtail" {
              for_each = log_destination.value.logtail != null ? [log_destination.value.logtail] : []
              content { token = logtail.value.token }
            }
            dynamic "open_search" {
              for_each = log_destination.value.open_search != null ? [log_destination.value.open_search] : []
              content {
                endpoint     = open_search.value.endpoint != "" ? open_search.value.endpoint : null
                index_name   = open_search.value.index_name != "" ? open_search.value.index_name : null
                cluster_name = open_search.value.cluster_name != "" ? open_search.value.cluster_name : null
                basic_auth {
                  user     = try(open_search.value.basic_auth.user, null)
                  password = try(open_search.value.basic_auth.password, null)
                }
              }
            }
          }
        }
      }
    }

    dynamic "worker" {
      for_each = var.spec.workers
      content {
        name               = worker.value.name
        source_dir         = worker.value.source_dir != "" ? worker.value.source_dir : null
        environment_slug   = worker.value.environment_slug != "" ? worker.value.environment_slug : null
        dockerfile_path    = worker.value.dockerfile_path != "" ? worker.value.dockerfile_path : null
        build_command      = worker.value.build_command != "" ? worker.value.build_command : null
        run_command        = worker.value.run_command != "" ? worker.value.run_command : null
        instance_size_slug = worker.value.instance_size_slug != "" ? worker.value.instance_size_slug : null
        instance_count     = worker.value.autoscaling != null ? null : (worker.value.instance_count > 0 ? worker.value.instance_count : null)

        dynamic "git" {
          for_each = worker.value.git != null ? [worker.value.git] : []
          content {
            repo_clone_url = git.value.repo_clone_url
            branch         = git.value.branch
          }
        }
        dynamic "github" {
          for_each = worker.value.github != null ? [worker.value.github] : []
          content {
            repo           = github.value.repo
            branch         = github.value.branch
            deploy_on_push = github.value.deploy_on_push
          }
        }
        dynamic "gitlab" {
          for_each = worker.value.gitlab != null ? [worker.value.gitlab] : []
          content {
            repo           = gitlab.value.repo
            branch         = gitlab.value.branch
            deploy_on_push = gitlab.value.deploy_on_push
          }
        }
        dynamic "bitbucket" {
          for_each = worker.value.bitbucket != null ? [worker.value.bitbucket] : []
          content {
            repo           = bitbucket.value.repo
            branch         = bitbucket.value.branch
            deploy_on_push = bitbucket.value.deploy_on_push
          }
        }
        dynamic "image" {
          for_each = worker.value.image != null ? [worker.value.image] : []
          content {
            registry_type        = upper(image.value.registry_type)
            registry             = image.value.registry != "" ? image.value.registry : null
            repository           = image.value.repository
            tag                  = image.value.tag != "" ? image.value.tag : null
            digest               = image.value.digest != "" ? image.value.digest : null
            registry_credentials = image.value.registry_credentials != "" ? image.value.registry_credentials : null
            dynamic "deploy_on_push" {
              for_each = image.value.deploy_on_push ? [1] : []
              content { enabled = true }
            }
          }
        }
        dynamic "liveness_health_check" {
          for_each = worker.value.liveness_health_check != null ? [worker.value.liveness_health_check] : []
          content {
            port                  = liveness_health_check.value.port
            http_path             = liveness_health_check.value.http_path != "" ? liveness_health_check.value.http_path : null
            initial_delay_seconds = liveness_health_check.value.initial_delay_seconds
            period_seconds        = liveness_health_check.value.period_seconds
            timeout_seconds       = liveness_health_check.value.timeout_seconds
            success_threshold     = liveness_health_check.value.success_threshold
            failure_threshold     = liveness_health_check.value.failure_threshold
          }
        }
        dynamic "autoscaling" {
          for_each = worker.value.autoscaling != null ? [worker.value.autoscaling] : []
          content {
            min_instance_count = autoscaling.value.min_instance_count
            max_instance_count = autoscaling.value.max_instance_count
            metrics {
              cpu { percent = autoscaling.value.cpu_percent }
            }
          }
        }
        dynamic "termination" {
          for_each = worker.value.termination != null ? [worker.value.termination] : []
          content {
            grace_period_seconds = termination.value.grace_period_seconds
          }
        }
        dynamic "env" {
          for_each = worker.value.envs
          content {
            key   = env.value.key
            value = env.value.secret != "" ? env.value.secret : env.value.plaintext
            type  = env.value.secret != "" ? "SECRET" : "GENERAL"
            scope = (env.value.scope == "" || endswith(env.value.scope, "_unspecified")) ? "RUN_AND_BUILD_TIME" : upper(env.value.scope)
          }
        }
        dynamic "alert" {
          for_each = worker.value.alerts
          content {
            rule     = upper(alert.value.rule)
            operator = upper(alert.value.operator)
            window   = upper(alert.value.window)
            value    = alert.value.value
            disabled = alert.value.disabled
            dynamic "destinations" {
              for_each = alert.value.destinations != null ? [alert.value.destinations] : []
              content {
                emails = length(destinations.value.emails) > 0 ? destinations.value.emails : null
                dynamic "slack_webhooks" {
                  for_each = destinations.value.slack_webhooks
                  content {
                    channel = slack_webhooks.value.channel
                    url     = slack_webhooks.value.url
                  }
                }
              }
            }
          }
        }
        dynamic "log_destination" {
          for_each = worker.value.log_destinations
          content {
            name = log_destination.value.name
            dynamic "papertrail" {
              for_each = log_destination.value.papertrail != null ? [log_destination.value.papertrail] : []
              content { endpoint = papertrail.value.endpoint }
            }
            dynamic "datadog" {
              for_each = log_destination.value.datadog != null ? [log_destination.value.datadog] : []
              content {
                api_key  = datadog.value.api_key
                endpoint = datadog.value.endpoint != "" ? datadog.value.endpoint : null
              }
            }
            dynamic "logtail" {
              for_each = log_destination.value.logtail != null ? [log_destination.value.logtail] : []
              content { token = logtail.value.token }
            }
            dynamic "open_search" {
              for_each = log_destination.value.open_search != null ? [log_destination.value.open_search] : []
              content {
                endpoint     = open_search.value.endpoint != "" ? open_search.value.endpoint : null
                index_name   = open_search.value.index_name != "" ? open_search.value.index_name : null
                cluster_name = open_search.value.cluster_name != "" ? open_search.value.cluster_name : null
                basic_auth {
                  user     = try(open_search.value.basic_auth.user, null)
                  password = try(open_search.value.basic_auth.password, null)
                }
              }
            }
          }
        }
      }
    }

    dynamic "job" {
      for_each = var.spec.jobs
      content {
        name               = job.value.name
        source_dir         = job.value.source_dir != "" ? job.value.source_dir : null
        environment_slug   = job.value.environment_slug != "" ? job.value.environment_slug : null
        dockerfile_path    = job.value.dockerfile_path != "" ? job.value.dockerfile_path : null
        build_command      = job.value.build_command != "" ? job.value.build_command : null
        run_command        = job.value.run_command != "" ? job.value.run_command : null
        instance_size_slug = job.value.instance_size_slug != "" ? job.value.instance_size_slug : null
        instance_count     = job.value.instance_count > 0 ? job.value.instance_count : null
        kind               = (job.value.kind == "" || endswith(job.value.kind, "_unspecified")) ? null : upper(job.value.kind)

        dynamic "git" {
          for_each = job.value.git != null ? [job.value.git] : []
          content {
            repo_clone_url = git.value.repo_clone_url
            branch         = git.value.branch
          }
        }
        dynamic "github" {
          for_each = job.value.github != null ? [job.value.github] : []
          content {
            repo           = github.value.repo
            branch         = github.value.branch
            deploy_on_push = github.value.deploy_on_push
          }
        }
        dynamic "gitlab" {
          for_each = job.value.gitlab != null ? [job.value.gitlab] : []
          content {
            repo           = gitlab.value.repo
            branch         = gitlab.value.branch
            deploy_on_push = gitlab.value.deploy_on_push
          }
        }
        dynamic "bitbucket" {
          for_each = job.value.bitbucket != null ? [job.value.bitbucket] : []
          content {
            repo           = bitbucket.value.repo
            branch         = bitbucket.value.branch
            deploy_on_push = bitbucket.value.deploy_on_push
          }
        }
        dynamic "image" {
          for_each = job.value.image != null ? [job.value.image] : []
          content {
            registry_type        = upper(image.value.registry_type)
            registry             = image.value.registry != "" ? image.value.registry : null
            repository           = image.value.repository
            tag                  = image.value.tag != "" ? image.value.tag : null
            digest               = image.value.digest != "" ? image.value.digest : null
            registry_credentials = image.value.registry_credentials != "" ? image.value.registry_credentials : null
            dynamic "deploy_on_push" {
              for_each = image.value.deploy_on_push ? [1] : []
              content { enabled = true }
            }
          }
        }
        dynamic "termination" {
          for_each = job.value.termination != null ? [job.value.termination] : []
          content {
            grace_period_seconds = termination.value.grace_period_seconds
          }
        }
        dynamic "env" {
          for_each = job.value.envs
          content {
            key   = env.value.key
            value = env.value.secret != "" ? env.value.secret : env.value.plaintext
            type  = env.value.secret != "" ? "SECRET" : "GENERAL"
            scope = (env.value.scope == "" || endswith(env.value.scope, "_unspecified")) ? "RUN_AND_BUILD_TIME" : upper(env.value.scope)
          }
        }
        dynamic "alert" {
          for_each = job.value.alerts
          content {
            rule     = upper(alert.value.rule)
            operator = upper(alert.value.operator)
            window   = upper(alert.value.window)
            value    = alert.value.value
            disabled = alert.value.disabled
            dynamic "destinations" {
              for_each = alert.value.destinations != null ? [alert.value.destinations] : []
              content {
                emails = length(destinations.value.emails) > 0 ? destinations.value.emails : null
                dynamic "slack_webhooks" {
                  for_each = destinations.value.slack_webhooks
                  content {
                    channel = slack_webhooks.value.channel
                    url     = slack_webhooks.value.url
                  }
                }
              }
            }
          }
        }
        dynamic "log_destination" {
          for_each = job.value.log_destinations
          content {
            name = log_destination.value.name
            dynamic "papertrail" {
              for_each = log_destination.value.papertrail != null ? [log_destination.value.papertrail] : []
              content { endpoint = papertrail.value.endpoint }
            }
            dynamic "datadog" {
              for_each = log_destination.value.datadog != null ? [log_destination.value.datadog] : []
              content {
                api_key  = datadog.value.api_key
                endpoint = datadog.value.endpoint != "" ? datadog.value.endpoint : null
              }
            }
            dynamic "logtail" {
              for_each = log_destination.value.logtail != null ? [log_destination.value.logtail] : []
              content { token = logtail.value.token }
            }
            dynamic "open_search" {
              for_each = log_destination.value.open_search != null ? [log_destination.value.open_search] : []
              content {
                endpoint     = open_search.value.endpoint != "" ? open_search.value.endpoint : null
                index_name   = open_search.value.index_name != "" ? open_search.value.index_name : null
                cluster_name = open_search.value.cluster_name != "" ? open_search.value.cluster_name : null
                basic_auth {
                  user     = try(open_search.value.basic_auth.user, null)
                  password = try(open_search.value.basic_auth.password, null)
                }
              }
            }
          }
        }
      }
    }

    dynamic "static_site" {
      for_each = var.spec.static_sites
      content {
        name              = static_site.value.name
        source_dir        = static_site.value.source_dir != "" ? static_site.value.source_dir : null
        environment_slug  = static_site.value.environment_slug != "" ? static_site.value.environment_slug : null
        dockerfile_path   = static_site.value.dockerfile_path != "" ? static_site.value.dockerfile_path : null
        build_command     = static_site.value.build_command != "" ? static_site.value.build_command : null
        output_dir        = static_site.value.output_dir != "" ? static_site.value.output_dir : null
        index_document    = static_site.value.index_document != "" ? static_site.value.index_document : null
        error_document    = static_site.value.error_document != "" ? static_site.value.error_document : null
        catchall_document = static_site.value.catchall_document != "" ? static_site.value.catchall_document : null

        dynamic "git" {
          for_each = static_site.value.git != null ? [static_site.value.git] : []
          content {
            repo_clone_url = git.value.repo_clone_url
            branch         = git.value.branch
          }
        }
        dynamic "github" {
          for_each = static_site.value.github != null ? [static_site.value.github] : []
          content {
            repo           = github.value.repo
            branch         = github.value.branch
            deploy_on_push = github.value.deploy_on_push
          }
        }
        dynamic "gitlab" {
          for_each = static_site.value.gitlab != null ? [static_site.value.gitlab] : []
          content {
            repo           = gitlab.value.repo
            branch         = gitlab.value.branch
            deploy_on_push = gitlab.value.deploy_on_push
          }
        }
        dynamic "bitbucket" {
          for_each = static_site.value.bitbucket != null ? [static_site.value.bitbucket] : []
          content {
            repo           = bitbucket.value.repo
            branch         = bitbucket.value.branch
            deploy_on_push = bitbucket.value.deploy_on_push
          }
        }
        dynamic "env" {
          for_each = static_site.value.envs
          content {
            key   = env.value.key
            value = env.value.secret != "" ? env.value.secret : env.value.plaintext
            type  = env.value.secret != "" ? "SECRET" : "GENERAL"
            scope = (env.value.scope == "" || endswith(env.value.scope, "_unspecified")) ? "RUN_AND_BUILD_TIME" : upper(env.value.scope)
          }
        }
      }
    }

    dynamic "function" {
      for_each = var.spec.functions
      content {
        name       = function.value.name
        source_dir = function.value.source_dir != "" ? function.value.source_dir : null

        dynamic "git" {
          for_each = function.value.git != null ? [function.value.git] : []
          content {
            repo_clone_url = git.value.repo_clone_url
            branch         = git.value.branch
          }
        }
        dynamic "github" {
          for_each = function.value.github != null ? [function.value.github] : []
          content {
            repo           = github.value.repo
            branch         = github.value.branch
            deploy_on_push = github.value.deploy_on_push
          }
        }
        dynamic "gitlab" {
          for_each = function.value.gitlab != null ? [function.value.gitlab] : []
          content {
            repo           = gitlab.value.repo
            branch         = gitlab.value.branch
            deploy_on_push = gitlab.value.deploy_on_push
          }
        }
        dynamic "bitbucket" {
          for_each = function.value.bitbucket != null ? [function.value.bitbucket] : []
          content {
            repo           = bitbucket.value.repo
            branch         = bitbucket.value.branch
            deploy_on_push = bitbucket.value.deploy_on_push
          }
        }
        dynamic "env" {
          for_each = function.value.envs
          content {
            key   = env.value.key
            value = env.value.secret != "" ? env.value.secret : env.value.plaintext
            type  = env.value.secret != "" ? "SECRET" : "GENERAL"
            scope = (env.value.scope == "" || endswith(env.value.scope, "_unspecified")) ? "RUN_AND_BUILD_TIME" : upper(env.value.scope)
          }
        }
        dynamic "alert" {
          for_each = function.value.alerts
          content {
            rule     = upper(alert.value.rule)
            operator = upper(alert.value.operator)
            window   = upper(alert.value.window)
            value    = alert.value.value
            disabled = alert.value.disabled
            dynamic "destinations" {
              for_each = alert.value.destinations != null ? [alert.value.destinations] : []
              content {
                emails = length(destinations.value.emails) > 0 ? destinations.value.emails : null
                dynamic "slack_webhooks" {
                  for_each = destinations.value.slack_webhooks
                  content {
                    channel = slack_webhooks.value.channel
                    url     = slack_webhooks.value.url
                  }
                }
              }
            }
          }
        }
        dynamic "log_destination" {
          for_each = function.value.log_destinations
          content {
            name = log_destination.value.name
            dynamic "papertrail" {
              for_each = log_destination.value.papertrail != null ? [log_destination.value.papertrail] : []
              content { endpoint = papertrail.value.endpoint }
            }
            dynamic "datadog" {
              for_each = log_destination.value.datadog != null ? [log_destination.value.datadog] : []
              content {
                api_key  = datadog.value.api_key
                endpoint = datadog.value.endpoint != "" ? datadog.value.endpoint : null
              }
            }
            dynamic "logtail" {
              for_each = log_destination.value.logtail != null ? [log_destination.value.logtail] : []
              content { token = logtail.value.token }
            }
            dynamic "open_search" {
              for_each = log_destination.value.open_search != null ? [log_destination.value.open_search] : []
              content {
                endpoint     = open_search.value.endpoint != "" ? open_search.value.endpoint : null
                index_name   = open_search.value.index_name != "" ? open_search.value.index_name : null
                cluster_name = open_search.value.cluster_name != "" ? open_search.value.cluster_name : null
                basic_auth {
                  user     = try(open_search.value.basic_auth.user, null)
                  password = try(open_search.value.basic_auth.password, null)
                }
              }
            }
          }
        }
      }
    }

    dynamic "database" {
      for_each = var.spec.databases
      content {
        name         = database.value.name != "" ? database.value.name : null
        engine       = (database.value.engine == "" || endswith(database.value.engine, "_unspecified")) ? null : upper(database.value.engine)
        version      = database.value.version != "" ? database.value.version : null
        production   = database.value.production
        cluster_name = database.value.cluster_name != "" ? database.value.cluster_name : null
        db_name      = database.value.db_name != "" ? database.value.db_name : null
        db_user      = database.value.db_user != "" ? database.value.db_user : null
      }
    }

    dynamic "domain" {
      for_each = var.spec.domains
      content {
        name     = domain.value.name
        type     = domain.value.type != "" ? domain.value.type : null
        wildcard = domain.value.wildcard
        zone     = domain.value.zone != "" ? domain.value.zone : null
      }
    }

    dynamic "alert" {
      for_each = var.spec.alerts
      content {
        rule     = upper(alert.value.rule)
        disabled = alert.value.disabled
        dynamic "destinations" {
          for_each = alert.value.destinations != null ? [alert.value.destinations] : []
          content {
            emails = length(destinations.value.emails) > 0 ? destinations.value.emails : null
            dynamic "slack_webhooks" {
              for_each = destinations.value.slack_webhooks
              content {
                channel = slack_webhooks.value.channel
                url     = slack_webhooks.value.url
              }
            }
          }
        }
      }
    }

    dynamic "ingress" {
      for_each = var.spec.ingress != null ? [var.spec.ingress] : []
      content {
        dynamic "rule" {
          for_each = ingress.value.rules
          content {
            dynamic "match" {
              for_each = rule.value.match != null ? [rule.value.match] : []
              content {
                dynamic "path" {
                  for_each = match.value.path_prefix != "" ? [match.value.path_prefix] : []
                  content { prefix = path.value }
                }
                dynamic "authority" {
                  for_each = match.value.authority_exact != "" ? [match.value.authority_exact] : []
                  content { exact = authority.value }
                }
              }
            }
            dynamic "component" {
              for_each = rule.value.component != null ? [rule.value.component] : []
              content {
                name                 = component.value.name
                preserve_path_prefix = component.value.preserve_path_prefix
                rewrite              = component.value.rewrite != "" ? component.value.rewrite : null
              }
            }
            dynamic "redirect" {
              for_each = rule.value.redirect != null ? [rule.value.redirect] : []
              content {
                uri           = redirect.value.uri != "" ? redirect.value.uri : null
                authority     = redirect.value.authority != "" ? redirect.value.authority : null
                port          = redirect.value.port
                scheme        = redirect.value.scheme != "" ? redirect.value.scheme : null
                redirect_code = redirect.value.redirect_code
              }
            }
            dynamic "cors" {
              for_each = rule.value.cors != null ? [rule.value.cors] : []
              content {
                dynamic "allow_origins" {
                  for_each = cors.value.allow_origins != null ? [cors.value.allow_origins] : []
                  content {
                    exact = allow_origins.value.exact != "" ? allow_origins.value.exact : null
                    regex = allow_origins.value.regex != "" ? allow_origins.value.regex : null
                  }
                }
                allow_methods     = length(cors.value.allow_methods) > 0 ? cors.value.allow_methods : null
                allow_headers     = length(cors.value.allow_headers) > 0 ? cors.value.allow_headers : null
                expose_headers    = length(cors.value.expose_headers) > 0 ? cors.value.expose_headers : null
                max_age           = cors.value.max_age != "" ? cors.value.max_age : null
                allow_credentials = cors.value.allow_credentials
              }
            }
          }
        }
        dynamic "secure_header" {
          for_each = ingress.value.secure_header != null ? [ingress.value.secure_header] : []
          content {
            key   = secure_header.value.key
            value = secure_header.value.value
          }
        }
      }
    }

    dynamic "egress" {
      for_each = (var.spec.egress == "" || endswith(var.spec.egress, "_unspecified")) ? [] : [var.spec.egress]
      content {
        type = upper(egress.value)
      }
    }

    dynamic "maintenance" {
      for_each = var.spec.maintenance != null ? [var.spec.maintenance] : []
      content {
        enabled          = maintenance.value.enabled
        archive          = maintenance.value.archive
        offline_page_url = maintenance.value.offline_page_url != "" ? maintenance.value.offline_page_url : null
      }
    }

    dynamic "vpc" {
      for_each = var.spec.vpc != "" ? [var.spec.vpc] : []
      content {
        id = vpc.value
      }
    }
  }
}
