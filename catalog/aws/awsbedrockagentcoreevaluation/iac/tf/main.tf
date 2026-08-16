# Amazon Bedrock AgentCore Evaluations: evaluators (LLM-judge or
# Lambda scorers), harnesses (repeatable agent test benches), and
# online evaluation configs (continuous scoring of sampled production
# sessions).
#
# Lifecycle facts the renders below depend on:
#   - evaluator and online-config names are identity seeds (AWS derives
#     "<name>-<10 chars>") with no rename; the harness name likewise
#     replaces on change -- everything else updates in place;
#   - creation waits on the provider's state machines (evaluator/config
#     CREATING -> ACTIVE; harness CREATING -> READY), with brief
#     upstream retries while IAM roles and Lambda permissions
#     propagate;
#   - the harness models `environment` and `truncation` as ATTRIBUTES
#     (`= [{...}]`), not blocks -- the provider generates them by
#     reflection; the rest of the harness surface is ordinary blocks;
#   - the spec's single `enabled` knob on online configs fans out to
#     the provider's TWO lifecycle fields (enable_on_create +
#     execution_status) -- one declarative knob, both engines wire it
#     identically;
#   - non-Bedrock model vendors authenticate via a Secrets Manager ARN
#     the harness reads at run time -- the module sends the reference,
#     never a key.

resource "aws_bedrockagentcore_evaluator" "this" {
  for_each = local.evaluators

  evaluator_name = each.value.name
  level          = each.value.level
  description    = try(each.value.description, "") != "" ? each.value.description : null
  kms_key_arn    = try(each.value.kms_key_arn, "") != "" ? each.value.kms_key_arn : null

  evaluator_config {
    # A Bedrock model judges against instructions and a rating scale.
    # evaluators is typed `any` in variables.tf (a Struct on the judge
    # model forces it); try() is how missing XOR arms stay unrendered.
    dynamic "llm_as_a_judge" {
      for_each = try(each.value.llm_as_a_judge, null) != null ? [each.value.llm_as_a_judge] : []
      content {
        instructions = llm_as_a_judge.value.instructions

        model_config {
          bedrock_evaluator_model_config {
            model_id = llm_as_a_judge.value.model.model_id

            # Model-specific passthrough fields travel as the JSON the
            # model's own API documents.
            additional_model_request_fields = try(llm_as_a_judge.value.model.additional_model_request_fields, null) != null ? jsonencode(llm_as_a_judge.value.model.additional_model_request_fields) : null

            dynamic "inference_config" {
              for_each = try(llm_as_a_judge.value.model.inference, null) != null ? [llm_as_a_judge.value.model.inference] : []
              content {
                max_tokens     = try(inference_config.value.max_tokens, 0) > 0 ? inference_config.value.max_tokens : null
                stop_sequences = try(length(inference_config.value.stop_sequences), 0) > 0 ? inference_config.value.stop_sequences : null
                # Presence-typed: only an explicit choice is sent, so the
                # module never fights the model's own defaults.
                temperature = try(inference_config.value.temperature, null)
                top_p       = try(inference_config.value.top_p, null)
              }
            }
          }
        }

        rating_scale {
          # Exactly one scale shape (spec-validated).
          dynamic "categorical" {
            for_each = try(llm_as_a_judge.value.rating_scale.categorical, [])
            content {
              label      = categorical.value.label
              definition = categorical.value.definition
            }
          }
          dynamic "numerical" {
            for_each = try(llm_as_a_judge.value.rating_scale.numerical, [])
            content {
              label      = numerical.value.label
              definition = numerical.value.definition
              value      = numerical.value.value
            }
          }
        }
      }
    }

    # Your Lambda function computes the score.
    dynamic "code_based" {
      for_each = try(each.value.code_based, null) != null ? [each.value.code_based] : []
      content {
        lambda_config {
          lambda_arn                = code_based.value.lambda_arn
          lambda_timeout_in_seconds = try(code_based.value.timeout_seconds, 0) > 0 ? code_based.value.timeout_seconds : null
        }
      }
    }
  }

  tags = local.aws_tags
}

resource "aws_bedrockagentcore_harness" "this" {
  for_each = local.harnesses

  harness_name       = each.value.name
  execution_role_arn = each.value.execution_role_arn

  # The model under test -- exactly one vendor arm (spec-validated).
  model {
    dynamic "bedrock_model_config" {
      for_each = each.value.model.bedrock != null ? [each.value.model.bedrock] : []
      content {
        model_id    = bedrock_model_config.value.model_id
        max_tokens  = bedrock_model_config.value.max_tokens > 0 ? bedrock_model_config.value.max_tokens : null
        temperature = bedrock_model_config.value.temperature
        top_p       = bedrock_model_config.value.top_p
      }
    }
    dynamic "gemini_model_config" {
      for_each = each.value.model.gemini != null ? [each.value.model.gemini] : []
      content {
        api_key_arn = gemini_model_config.value.api_key_arn
        model_id    = gemini_model_config.value.model_id
        max_tokens  = gemini_model_config.value.max_tokens > 0 ? gemini_model_config.value.max_tokens : null
        temperature = gemini_model_config.value.temperature
        top_p       = gemini_model_config.value.top_p
        top_k       = gemini_model_config.value.top_k
      }
    }
    dynamic "openai_model_config" {
      for_each = each.value.model.openai != null ? [each.value.model.openai] : []
      content {
        api_key_arn = openai_model_config.value.api_key_arn
        model_id    = openai_model_config.value.model_id
        max_tokens  = openai_model_config.value.max_tokens > 0 ? openai_model_config.value.max_tokens : null
        temperature = openai_model_config.value.temperature
        top_p       = openai_model_config.value.top_p
      }
    }
  }

  # System prompts prepended to every run, in order.
  dynamic "system_prompt" {
    for_each = each.value.system_prompts
    content {
      text = system_prompt.value.text
    }
  }

  # Tools the agent under test may call. Each entry's config arm was
  # spec-validated against its declared type (stricter than the
  # provider, which silently takes the first configured arm).
  dynamic "tool" {
    for_each = each.value.tools
    content {
      name = tool.value.name != "" ? tool.value.name : null
      type = tool.value.type

      dynamic "config" {
        for_each = tool.value.remote_mcp != null || tool.value.agentcore_browser != null || tool.value.agentcore_gateway != null || tool.value.inline_function != null || tool.value.agentcore_code_interpreter != null ? [tool.value] : []
        content {
          dynamic "remote_mcp" {
            for_each = config.value.remote_mcp != null ? [config.value.remote_mcp] : []
            content {
              url     = remote_mcp.value.url
              headers = length(remote_mcp.value.headers) > 0 ? remote_mcp.value.headers : null
            }
          }
          dynamic "agentcore_browser" {
            for_each = config.value.agentcore_browser != null ? [config.value.agentcore_browser] : []
            content {
              browser_arn = agentcore_browser.value.browser_arn != "" ? agentcore_browser.value.browser_arn : null
            }
          }
          dynamic "agentcore_code_interpreter" {
            for_each = config.value.agentcore_code_interpreter != null ? [config.value.agentcore_code_interpreter] : []
            content {
              code_interpreter_arn = agentcore_code_interpreter.value.code_interpreter_arn != "" ? agentcore_code_interpreter.value.code_interpreter_arn : null
            }
          }
          dynamic "agentcore_gateway" {
            for_each = config.value.agentcore_gateway != null ? [config.value.agentcore_gateway] : []
            content {
              gateway_arn = agentcore_gateway.value.gateway_arn
              dynamic "outbound_auth" {
                for_each = agentcore_gateway.value.outbound_auth != null ? [agentcore_gateway.value.outbound_auth] : []
                content {
                  aws_iam = outbound_auth.value.aws_iam ? true : null
                  none    = outbound_auth.value.no_auth ? true : null
                  dynamic "oauth" {
                    for_each = outbound_auth.value.oauth != null ? [outbound_auth.value.oauth] : []
                    content {
                      provider_arn       = oauth.value.provider_arn
                      scopes             = oauth.value.scopes
                      custom_parameters  = length(oauth.value.custom_parameters) > 0 ? oauth.value.custom_parameters : null
                      default_return_url = oauth.value.default_return_url != "" ? oauth.value.default_return_url : null
                      grant_type         = oauth.value.grant_type != "" ? oauth.value.grant_type : null
                    }
                  }
                }
              }
            }
          }
          dynamic "inline_function" {
            for_each = config.value.inline_function != null ? [config.value.inline_function] : []
            content {
              description  = inline_function.value.description
              input_schema = inline_function.value.input_schema
            }
          }
        }
      }
    }
  }

  # Skill bundle paths loaded into the harness.
  dynamic "skill" {
    for_each = each.value.skill_paths
    content {
      path = skill.value
    }
  }

  # AgentCore memory the harness reads/writes during runs.
  dynamic "memory" {
    for_each = each.value.memory != null ? [each.value.memory] : []
    content {
      agentcore_memory_configuration {
        arn            = memory.value.memory_arn
        actor_id       = memory.value.actor_id != "" ? memory.value.actor_id : null
        messages_count = memory.value.messages_count > 0 ? memory.value.messages_count : null
        dynamic "retrieval_config" {
          for_each = memory.value.retrieval != null ? [memory.value.retrieval] : []
          content {
            map_block_key   = retrieval_config.value.namespace
            relevance_score = retrieval_config.value.relevance_score
            strategy_id     = retrieval_config.value.strategy_id != "" ? retrieval_config.value.strategy_id : null
            top_k           = retrieval_config.value.top_k > 0 ? retrieval_config.value.top_k : null
          }
        }
      }
    }
  }

  environment_variables = length(each.value.environment_variables) > 0 ? each.value.environment_variables : null

  # The explicit runtime environment. The provider models this as an
  # ATTRIBUTE (reflection-generated), so it renders as nested object
  # lists rather than blocks.
  environment = each.value.runtime_environment != null ? [{
    agentcore_runtime_environment = [{
      agent_runtime_arn  = each.value.runtime_environment.agent_runtime_arn != "" ? each.value.runtime_environment.agent_runtime_arn : null
      agent_runtime_id   = null
      agent_runtime_name = null
      filesystem_configuration = [for f in each.value.runtime_environment.filesystems : {
        efs_access_point = f.efs_access_point_arn != "" ? [{
          access_point_arn = f.efs_access_point_arn
          mount_path       = f.mount_path
        }] : []
        s3_files_access_point = f.s3_access_point_arn != "" ? [{
          access_point_arn = f.s3_access_point_arn
          mount_path       = f.mount_path
        }] : []
        session_storage = f.session_storage ? [{
          mount_path = f.mount_path
        }] : []
      }]
      lifecycle_configuration = each.value.runtime_environment.lifecycle != null ? [{
        idle_runtime_session_timeout = each.value.runtime_environment.lifecycle.idle_runtime_session_timeout_seconds > 0 ? each.value.runtime_environment.lifecycle.idle_runtime_session_timeout_seconds : null
        max_lifetime                 = each.value.runtime_environment.lifecycle.max_lifetime_seconds > 0 ? each.value.runtime_environment.lifecycle.max_lifetime_seconds : null
      }] : []
      network_configuration = each.value.runtime_environment.network != null ? [{
        network_mode = each.value.runtime_environment.network.mode
        network_mode_config = each.value.runtime_environment.network.vpc_config != null ? [{
          subnets                     = each.value.runtime_environment.network.vpc_config.subnets
          security_groups             = each.value.runtime_environment.network.vpc_config.security_groups
          require_service_s3_endpoint = each.value.runtime_environment.network.vpc_config.require_service_s3_endpoint
        }] : []
      }] : []
    }]
  }] : null

  # The harness's own container image (environment artifact).
  dynamic "environment_artifact" {
    for_each = each.value.container_image_uri != "" ? [each.value.container_image_uri] : []
    content {
      container_configuration {
        container_uri = environment_artifact.value
      }
    }
  }

  # Inbound JWT authorization (the shared AgentCore authorizer shape;
  # AWS requires the endpoint address type on the harness's
  # managed-VPC endpoint).
  dynamic "authorizer_configuration" {
    for_each = each.value.custom_jwt_authorizer != null ? [each.value.custom_jwt_authorizer] : []
    content {
      custom_jwt_authorizer {
        discovery_url    = authorizer_configuration.value.discovery_url
        allowed_audience = length(authorizer_configuration.value.allowed_audience) > 0 ? authorizer_configuration.value.allowed_audience : null
        allowed_clients  = length(authorizer_configuration.value.allowed_clients) > 0 ? authorizer_configuration.value.allowed_clients : null
        allowed_scopes   = length(authorizer_configuration.value.allowed_scopes) > 0 ? authorizer_configuration.value.allowed_scopes : null

        dynamic "allowed_workload_configuration" {
          for_each = authorizer_configuration.value.allowed_workloads != null ? [authorizer_configuration.value.allowed_workloads] : []
          content {
            workload_identities = length(allowed_workload_configuration.value.workload_identities) > 0 ? allowed_workload_configuration.value.workload_identities : null
            dynamic "hosting_environment" {
              for_each = allowed_workload_configuration.value.hosting_environment_arns
              content {
                arn = hosting_environment.value
              }
            }
          }
        }

        dynamic "custom_claim" {
          for_each = authorizer_configuration.value.custom_claims
          content {
            inbound_token_claim_name       = custom_claim.value.claim_name
            inbound_token_claim_value_type = custom_claim.value.value_type
            authorizing_claim_match_value {
              claim_match_operator = custom_claim.value.match_operator
              claim_match_value {
                match_value_string      = custom_claim.value.match_value != "" ? custom_claim.value.match_value : null
                match_value_string_list = length(custom_claim.value.match_values) > 0 ? custom_claim.value.match_values : null
              }
            }
          }
        }

        dynamic "private_endpoint" {
          for_each = authorizer_configuration.value.private_endpoint != null ? [authorizer_configuration.value.private_endpoint] : []
          content {
            dynamic "managed_vpc_resource" {
              for_each = private_endpoint.value.managed_vpc != null ? [private_endpoint.value.managed_vpc] : []
              content {
                vpc_identifier           = managed_vpc_resource.value.vpc_id
                subnet_ids               = managed_vpc_resource.value.subnet_ids
                security_group_ids       = length(managed_vpc_resource.value.security_group_ids) > 0 ? managed_vpc_resource.value.security_group_ids : null
                endpoint_ip_address_type = managed_vpc_resource.value.endpoint_ip_address_type
                routing_domain           = managed_vpc_resource.value.routing_domain != "" ? managed_vpc_resource.value.routing_domain : null
                tags                     = length(managed_vpc_resource.value.tags) > 0 ? managed_vpc_resource.value.tags : null
              }
            }
            dynamic "self_managed_lattice_resource" {
              for_each = private_endpoint.value.self_managed_lattice != null ? [private_endpoint.value.self_managed_lattice] : []
              content {
                resource_configuration_identifier = self_managed_lattice_resource.value.resource_configuration_id
              }
            }
          }
        }

        dynamic "private_endpoint_overrides" {
          for_each = authorizer_configuration.value.private_endpoint_overrides
          content {
            domain = private_endpoint_overrides.value.domain
            private_endpoint {
              dynamic "managed_vpc_resource" {
                for_each = private_endpoint_overrides.value.private_endpoint.managed_vpc != null ? [private_endpoint_overrides.value.private_endpoint.managed_vpc] : []
                content {
                  vpc_identifier           = managed_vpc_resource.value.vpc_id
                  subnet_ids               = managed_vpc_resource.value.subnet_ids
                  security_group_ids       = length(managed_vpc_resource.value.security_group_ids) > 0 ? managed_vpc_resource.value.security_group_ids : null
                  endpoint_ip_address_type = managed_vpc_resource.value.endpoint_ip_address_type
                  routing_domain           = managed_vpc_resource.value.routing_domain != "" ? managed_vpc_resource.value.routing_domain : null
                  tags                     = length(managed_vpc_resource.value.tags) > 0 ? managed_vpc_resource.value.tags : null
                }
              }
              dynamic "self_managed_lattice_resource" {
                for_each = private_endpoint_overrides.value.private_endpoint.self_managed_lattice != null ? [private_endpoint_overrides.value.private_endpoint.self_managed_lattice] : []
                content {
                  resource_configuration_identifier = self_managed_lattice_resource.value.resource_configuration_id
                }
              }
            }
          }
        }
      }
    }
  }

  allowed_tools   = length(each.value.allowed_tools) > 0 ? each.value.allowed_tools : null
  max_iterations  = each.value.max_iterations > 0 ? each.value.max_iterations : null
  max_tokens      = each.value.max_tokens > 0 ? each.value.max_tokens : null
  timeout_seconds = each.value.timeout_seconds > 0 ? each.value.timeout_seconds : null

  # Truncation is an ATTRIBUTE like environment (see the header
  # comment).
  truncation = each.value.truncation != null ? [{
    strategy = each.value.truncation.strategy
    config = each.value.truncation.sliding_window != null || each.value.truncation.summarization != null ? [{
      sliding_window = each.value.truncation.sliding_window != null ? [{
        messages_count = each.value.truncation.sliding_window.messages_count
      }] : []
      summarization = each.value.truncation.summarization != null ? [{
        summary_ratio               = each.value.truncation.summarization.summary_ratio
        preserve_recent_messages    = each.value.truncation.summarization.preserve_recent_messages > 0 ? each.value.truncation.summarization.preserve_recent_messages : null
        summarization_system_prompt = each.value.truncation.summarization.summarization_system_prompt != "" ? each.value.truncation.summarization.summarization_system_prompt : null
      }] : []
    }] : []
  }] : null

  tags = local.aws_tags
}

resource "aws_bedrockagentcore_online_evaluation_config" "this" {
  for_each = local.online_configs

  online_evaluation_config_name = each.value.name
  description                   = each.value.description != "" ? each.value.description : null
  evaluation_execution_role_arn = each.value.execution_role_arn

  # The spec's single `enabled` knob fans out to the provider's TWO
  # lifecycle fields (create-time intent + post-create status) -- one
  # declarative knob, both engines wire it identically.
  enable_on_create = each.value.enabled != null ? each.value.enabled : true
  execution_status = (each.value.enabled != null ? each.value.enabled : true) ? "ENABLED" : "DISABLED"

  data_source_config {
    cloudwatch_logs {
      log_group_names = each.value.data_source.log_group_names
      service_names   = each.value.data_source.service_names
    }
  }

  # In-bundle evaluator names resolve to the created evaluator's
  # AWS-generated ID (and gain the dependency edge); builtins and full
  # custom IDs pass through as literals.
  dynamic "evaluator" {
    for_each = each.value.evaluator_ids
    content {
      evaluator_id = contains(local.bundle_evaluator_names, evaluator.value) ? aws_bedrockagentcore_evaluator.this[evaluator.value].evaluator_id : evaluator.value
    }
  }

  rule {
    sampling_config {
      sampling_percentage = each.value.rule.sampling_percentage
    }
    dynamic "filter" {
      for_each = each.value.rule.filters
      content {
        key      = filter.value.key
        operator = filter.value.operator
        value {
          # Exactly one typed value (spec-validated).
          string_value  = filter.value.string_value != "" ? filter.value.string_value : null
          boolean_value = filter.value.boolean_value
          double_value  = filter.value.double_value
        }
      }
    }
    dynamic "session_config" {
      for_each = each.value.rule.session_timeout_minutes > 0 ? [each.value.rule.session_timeout_minutes] : []
      content {
        session_timeout_minutes = session_config.value
      }
    }
  }

  tags = local.aws_tags
}
