# CodePipeline -- the release-orchestration unit. The pipeline is a
# metadata-only control-plane resource: create/update/delete are single API
# calls that complete in seconds. The only operational wait is IAM eventual
# consistency on a freshly created pipeline role, which the provider absorbs
# with a bounded retry on create.

resource "aws_codepipeline" "this" {
  name     = local.pipeline_name
  role_arn = var.spec.role_arn
  # pipeline_type/execution_mode pass through unmodified: the spec's
  # V2/SUPERSEDED defaults are materialized by the platform when the
  # manifest is loaded, so the module never re-derives them (one source of
  # truth; same pass-through in the Pulumi module). A raw tfvars invocation
  # that bypasses manifest loading and omits pipeline_type gets the
  # PROVIDER default, which is V1.
  pipeline_type  = var.spec.pipeline_type
  execution_mode = var.spec.execution_mode
  tags           = local.tags

  # --- Artifact stores --------------------------------------------------------
  # Two render arms because AWS models single-region (no region field allowed)
  # and cross-region (region required per store) differently.

  dynamic "artifact_store" {
    for_each = local.is_single_region ? var.spec.artifact_stores : []
    content {
      location = artifact_store.value.location
      type     = "S3"

      dynamic "encryption_key" {
        for_each = artifact_store.value.encryption_key_id != "" ? [artifact_store.value.encryption_key_id] : []
        content {
          id   = encryption_key.value
          type = "KMS"
        }
      }
    }
  }

  dynamic "artifact_store" {
    for_each = local.is_single_region ? [] : var.spec.artifact_stores
    content {
      location = artifact_store.value.location
      type     = "S3"
      region   = artifact_store.value.region

      dynamic "encryption_key" {
        for_each = artifact_store.value.encryption_key_id != "" ? [artifact_store.value.encryption_key_id] : []
        content {
          id   = encryption_key.value
          type = "KMS"
        }
      }
    }
  }

  # --- Stages -----------------------------------------------------------------

  dynamic "stage" {
    for_each = var.spec.stages
    content {
      name = stage.value.name

      dynamic "action" {
        for_each = stage.value.actions
        content {
          name             = action.value.name
          category         = action.value.category
          owner            = action.value.owner
          provider         = action.value.provider
          version          = action.value.version
          configuration    = length(action.value.configuration) > 0 ? action.value.configuration : null
          input_artifacts  = length(action.value.input_artifacts) > 0 ? action.value.input_artifacts : null
          output_artifacts = length(action.value.output_artifacts) > 0 ? action.value.output_artifacts : null
          namespace        = action.value.namespace != "" ? action.value.namespace : null
          region           = action.value.region != "" ? action.value.region : null
          role_arn         = action.value.role_arn != "" ? action.value.role_arn : null
          run_order        = action.value.run_order > 0 ? action.value.run_order : null
          # Per-action timeout override (omit for the provider default).
          timeout_in_minutes = action.value.timeout_in_minutes > 0 ? action.value.timeout_in_minutes : null

          # Compute-action surface: inline shell commands with exported
          # variables and file-based output artifacts (Compute actions use
          # these INSTEAD of plain output_artifacts — the spec's CEL
          # enforces the split before the provider's plan-time check).
          commands         = length(action.value.commands) > 0 ? action.value.commands : null
          output_variables = length(action.value.output_variables) > 0 ? action.value.output_variables : null

          dynamic "output_artifacts_for_compute_action" {
            for_each = action.value.output_artifacts_for_compute_action
            content {
              name  = output_artifacts_for_compute_action.value.name
              files = length(output_artifacts_for_compute_action.value.files) > 0 ? output_artifacts_for_compute_action.value.files : null
            }
          }
        }
      }

      # Entry gate: rules that must pass before the stage starts (e.g., a
      # DeploymentWindow rule admitting executions only in business hours).
      dynamic "before_entry" {
        for_each = stage.value.before_entry != null ? [stage.value.before_entry] : []
        content {
          condition {
            result = before_entry.value.result != "" ? before_entry.value.result : null

            dynamic "rule" {
              for_each = before_entry.value.rules
              content {
                name = rule.value.name

                # category/owner are presence-carrying optionals whose spec
                # defaults ("Rule"/"AWS") are the only values AWS accepts today;
                # applied here, never left to the provider (same as Pulumi).
                rule_type_id {
                  category = rule.value.rule_type_id.category != null ? rule.value.rule_type_id.category : "Rule"
                  owner    = rule.value.rule_type_id.owner != null ? rule.value.rule_type_id.owner : "AWS"
                  provider = rule.value.rule_type_id.provider
                  version  = rule.value.rule_type_id.version != "" ? rule.value.rule_type_id.version : null
                }

                configuration      = length(rule.value.configuration) > 0 ? rule.value.configuration : null
                commands           = length(rule.value.commands) > 0 ? rule.value.commands : null
                input_artifacts    = length(rule.value.input_artifacts) > 0 ? rule.value.input_artifacts : null
                region             = rule.value.region != "" ? rule.value.region : null
                role_arn           = rule.value.role_arn != "" ? rule.value.role_arn : null
                timeout_in_minutes = rule.value.timeout_in_minutes > 0 ? rule.value.timeout_in_minutes : null
              }
            }
          }
        }
      }

      # Post-success verification: a failing rule fails the stage despite
      # successful actions (e.g., a post-deploy CloudWatchAlarm check).
      dynamic "on_success" {
        for_each = stage.value.on_success != null ? [stage.value.on_success] : []
        content {
          condition {
            result = on_success.value.result != "" ? on_success.value.result : null

            dynamic "rule" {
              for_each = on_success.value.rules
              content {
                name = rule.value.name

                # category/owner are presence-carrying optionals whose spec
                # defaults ("Rule"/"AWS") are the only values AWS accepts today;
                # applied here, never left to the provider (same as Pulumi).
                rule_type_id {
                  category = rule.value.rule_type_id.category != null ? rule.value.rule_type_id.category : "Rule"
                  owner    = rule.value.rule_type_id.owner != null ? rule.value.rule_type_id.owner : "AWS"
                  provider = rule.value.rule_type_id.provider
                  version  = rule.value.rule_type_id.version != "" ? rule.value.rule_type_id.version : null
                }

                configuration      = length(rule.value.configuration) > 0 ? rule.value.configuration : null
                commands           = length(rule.value.commands) > 0 ? rule.value.commands : null
                input_artifacts    = length(rule.value.input_artifacts) > 0 ? rule.value.input_artifacts : null
                region             = rule.value.region != "" ? rule.value.region : null
                role_arn           = rule.value.role_arn != "" ? rule.value.role_arn : null
                timeout_in_minutes = rule.value.timeout_in_minutes > 0 ? rule.value.timeout_in_minutes : null
              }
            }
          }
        }
      }

      # Failure handling: automatic rollback to the last successful state,
      # automatic retry, or rule-gated handling.
      dynamic "on_failure" {
        for_each = stage.value.on_failure != null ? [stage.value.on_failure] : []
        content {
          result = on_failure.value.result != "" ? on_failure.value.result : null

          dynamic "retry_configuration" {
            for_each = on_failure.value.retry_configuration != null ? [on_failure.value.retry_configuration] : []
            content {
              retry_mode = retry_configuration.value.retry_mode
            }
          }

          dynamic "condition" {
            for_each = on_failure.value.condition != null ? [on_failure.value.condition] : []
            content {
              result = condition.value.result != "" ? condition.value.result : null

              dynamic "rule" {
                for_each = condition.value.rules
                content {
                  name = rule.value.name

                  rule_type_id {
                    category = rule.value.rule_type_id.category != null ? rule.value.rule_type_id.category : "Rule"
                    owner    = rule.value.rule_type_id.owner != null ? rule.value.rule_type_id.owner : "AWS"
                    provider = rule.value.rule_type_id.provider
                    version  = rule.value.rule_type_id.version != "" ? rule.value.rule_type_id.version : null
                  }

                  configuration      = length(rule.value.configuration) > 0 ? rule.value.configuration : null
                  commands           = length(rule.value.commands) > 0 ? rule.value.commands : null
                  input_artifacts    = length(rule.value.input_artifacts) > 0 ? rule.value.input_artifacts : null
                  region             = rule.value.region != "" ? rule.value.region : null
                  role_arn           = rule.value.role_arn != "" ? rule.value.role_arn : null
                  timeout_in_minutes = rule.value.timeout_in_minutes > 0 ? rule.value.timeout_in_minutes : null
                }
              }
            }
          }
        }
      }
    }
  }

  # --- Triggers (V2 git-event execution) --------------------------------------

  dynamic "trigger" {
    for_each = var.spec.triggers
    content {
      provider_type = trigger.value.provider_type

      git_configuration {
        source_action_name = trigger.value.git_configuration.source_action_name

        dynamic "push" {
          for_each = trigger.value.git_configuration.push
          content {
            dynamic "branches" {
              for_each = push.value.branches != null ? [push.value.branches] : []
              content {
                includes = length(branches.value.includes) > 0 ? branches.value.includes : null
                excludes = length(branches.value.excludes) > 0 ? branches.value.excludes : null
              }
            }
            dynamic "file_paths" {
              for_each = push.value.file_paths != null ? [push.value.file_paths] : []
              content {
                includes = length(file_paths.value.includes) > 0 ? file_paths.value.includes : null
                excludes = length(file_paths.value.excludes) > 0 ? file_paths.value.excludes : null
              }
            }
            dynamic "tags" {
              for_each = push.value.tags != null ? [push.value.tags] : []
              content {
                includes = length(tags.value.includes) > 0 ? tags.value.includes : null
                excludes = length(tags.value.excludes) > 0 ? tags.value.excludes : null
              }
            }
          }
        }

        dynamic "pull_request" {
          for_each = trigger.value.git_configuration.pull_request
          content {
            events = length(pull_request.value.events) > 0 ? pull_request.value.events : null

            dynamic "branches" {
              for_each = pull_request.value.branches != null ? [pull_request.value.branches] : []
              content {
                includes = length(branches.value.includes) > 0 ? branches.value.includes : null
                excludes = length(branches.value.excludes) > 0 ? branches.value.excludes : null
              }
            }
            dynamic "file_paths" {
              for_each = pull_request.value.file_paths != null ? [pull_request.value.file_paths] : []
              content {
                includes = length(file_paths.value.includes) > 0 ? file_paths.value.includes : null
                excludes = length(file_paths.value.excludes) > 0 ? file_paths.value.excludes : null
              }
            }
          }
        }
      }
    }
  }

  # --- Pipeline-level variables (V2) -------------------------------------------

  dynamic "variable" {
    for_each = var.spec.variables
    content {
      name          = variable.value.name
      default_value = variable.value.default_value != "" ? variable.value.default_value : null
      description   = variable.value.description != "" ? variable.value.description : null
    }
  }
}
