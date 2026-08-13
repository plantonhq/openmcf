# Amazon Bedrock flow: a node graph orchestrating prompts, agents,
# knowledge bases, Lambdas, and control-flow logic into one invocable
# pipeline. The module owns the flow's mutable DRAFT definition.
#
# AWS validates the GRAPH server-side (unreachable nodes, type mismatches,
# missing connections) at create/update -- the spec validates shapes, the
# service validates topology.

resource "aws_bedrockagent_flow" "this" {
  # Create-time naming basis; doubles as the Name tag. metadata.name on
  # both engines.
  name = local.flow_name

  execution_role_arn = var.spec.execution_role_arn

  description                 = var.spec.description != "" ? var.spec.description : null
  customer_encryption_key_arn = var.spec.customer_encryption_key_arn != "" ? var.spec.customer_encryption_key_arn : null

  dynamic "definition" {
    for_each = local.has_definition ? [var.spec.definition] : []
    content {

      dynamic "node" {
        for_each = definition.value.nodes
        content {
          name = node.value.name
          type = node.value.type

          dynamic "input" {
            for_each = try(node.value.inputs, [])
            content {
              name       = input.value.name
              expression = input.value.expression
              type       = input.value.type
              category   = try(input.value.category, "") != "" ? input.value.category : null
            }
          }

          dynamic "output" {
            for_each = try(node.value.outputs, [])
            content {
              name = output.value.name
              type = output.value.type
            }
          }

          # Exactly one AWS union member per configurable class; the
          # structural classes carry an EMPTY member derived from the
          # type, and the Loop family none (provider gap at the pin).
          dynamic "configuration" {
            for_each = contains(local.empty_config_types, node.value.type) || try(node.value.agent, null) != null || try(node.value.prompt, null) != null || try(node.value.knowledge_base, null) != null || try(node.value.lambda_function, null) != null || try(node.value.lex, null) != null || try(node.value.condition, null) != null || try(node.value.inline_code, null) != null || try(node.value.retrieval, null) != null || try(node.value.storage, null) != null ? [node.value] : []
            content {

              dynamic "input" {
                for_each = configuration.value.type == "Input" ? [1] : []
                content {}
              }

              dynamic "output" {
                for_each = configuration.value.type == "Output" ? [1] : []
                content {}
              }

              dynamic "iterator" {
                for_each = configuration.value.type == "Iterator" ? [1] : []
                content {}
              }

              dynamic "collector" {
                for_each = configuration.value.type == "Collector" ? [1] : []
                content {}
              }

              dynamic "agent" {
                for_each = try(configuration.value.agent, null) != null ? [configuration.value.agent] : []
                content {
                  agent_alias_arn = agent.value.agent_alias_arn
                }
              }

              dynamic "lambda_function" {
                for_each = try(configuration.value.lambda_function, null) != null ? [configuration.value.lambda_function] : []
                content {
                  lambda_arn = lambda_function.value.lambda_arn
                }
              }

              dynamic "lex" {
                for_each = try(configuration.value.lex, null) != null ? [configuration.value.lex] : []
                content {
                  bot_alias_arn = lex.value.bot_alias_arn
                  locale_id     = lex.value.locale_id
                }
              }

              dynamic "inline_code" {
                for_each = try(configuration.value.inline_code, null) != null ? [configuration.value.inline_code] : []
                content {
                  code = inline_code.value.code
                  # Python_3 is the only language AWS defines -- the
                  # module owns the constant.
                  language = "Python_3"
                }
              }

              dynamic "condition" {
                for_each = try(configuration.value.condition, null) != null ? [configuration.value.condition] : []
                content {
                  dynamic "condition" {
                    for_each = condition.value.conditions
                    content {
                      name       = condition.value.name
                      expression = try(condition.value.expression, "") != "" ? condition.value.expression : null
                    }
                  }
                }
              }

              dynamic "knowledge_base" {
                for_each = try(configuration.value.knowledge_base, null) != null ? [configuration.value.knowledge_base] : []
                content {
                  knowledge_base_id = knowledge_base.value.knowledge_base_id
                  model_id          = try(knowledge_base.value.model_id, "") != "" ? knowledge_base.value.model_id : null
                  number_of_results = try(knowledge_base.value.number_of_results, 0) != 0 ? knowledge_base.value.number_of_results : null

                  dynamic "guardrail_configuration" {
                    for_each = try(knowledge_base.value.guardrail, null) != null ? [knowledge_base.value.guardrail] : []
                    content {
                      guardrail_identifier = guardrail_configuration.value.guardrail_id
                      guardrail_version    = guardrail_configuration.value.version
                    }
                  }

                  dynamic "inference_configuration" {
                    for_each = try(knowledge_base.value.inference_configuration, null) != null ? [knowledge_base.value.inference_configuration] : []
                    content {
                      text {
                        max_tokens     = try(inference_configuration.value.max_tokens, null)
                        stop_sequences = length(try(inference_configuration.value.stop_sequences, [])) > 0 ? inference_configuration.value.stop_sequences : null
                        temperature    = try(inference_configuration.value.temperature, null)
                        top_p          = try(inference_configuration.value.top_p, null)
                      }
                    }
                  }
                }
              }

              dynamic "retrieval" {
                for_each = try(configuration.value.retrieval, null) != null ? [configuration.value.retrieval] : []
                content {
                  service_configuration {
                    # S3 is the only retrieval service AWS defines -- the
                    # module owns the union member.
                    s3 {
                      bucket_name = retrieval.value.bucket_name
                    }
                  }
                }
              }

              dynamic "storage" {
                for_each = try(configuration.value.storage, null) != null ? [configuration.value.storage] : []
                content {
                  service_configuration {
                    # S3 is the only storage service AWS defines.
                    s3 {
                      bucket_name = storage.value.bucket_name
                    }
                  }
                }
              }

              dynamic "prompt" {
                for_each = try(configuration.value.prompt, null) != null ? [configuration.value.prompt] : []
                content {

                  dynamic "guardrail_configuration" {
                    for_each = try(prompt.value.guardrail, null) != null ? [prompt.value.guardrail] : []
                    content {
                      guardrail_identifier = guardrail_configuration.value.guardrail_id
                      guardrail_version    = guardrail_configuration.value.version
                    }
                  }

                  source_configuration {

                    dynamic "resource" {
                      for_each = try(prompt.value.prompt_arn, "") != "" ? [prompt.value.prompt_arn] : []
                      content {
                        prompt_arn = resource.value
                      }
                    }

                    dynamic "inline" {
                      for_each = try(prompt.value.inline, null) != null ? [prompt.value.inline] : []
                      content {
                        model_id = inline.value.model_id
                        # TEXT or CHAT is derived from which template arm
                        # is set (exactly one, per the spec's CEL guard).
                        template_type = try(inline.value.text, null) != null ? "TEXT" : "CHAT"

                        additional_model_request_fields = try(inline.value.additional_model_request_fields, null) != null ? jsonencode(inline.value.additional_model_request_fields) : null

                        dynamic "inference_configuration" {
                          for_each = try(inline.value.inference_configuration, null) != null ? [inline.value.inference_configuration] : []
                          content {
                            text {
                              max_tokens     = try(inference_configuration.value.max_tokens, null)
                              stop_sequences = length(try(inference_configuration.value.stop_sequences, [])) > 0 ? inference_configuration.value.stop_sequences : null
                              temperature    = try(inference_configuration.value.temperature, null)
                              top_p          = try(inference_configuration.value.top_p, null)
                            }
                          }
                        }

                        template_configuration {

                          dynamic "text" {
                            for_each = try(inline.value.text, null) != null ? [inline.value.text] : []
                            content {
                              text = text.value.text
                              dynamic "cache_point" {
                                for_each = try(text.value.cache_point, false) ? [1] : []
                                content {
                                  # "default" is the only cache point type
                                  # AWS defines.
                                  type = "default"
                                }
                              }
                              dynamic "input_variable" {
                                for_each = try(text.value.input_variables, [])
                                content {
                                  name = input_variable.value
                                }
                              }
                            }
                          }

                          dynamic "chat" {
                            for_each = try(inline.value.chat, null) != null ? [inline.value.chat] : []
                            content {
                              dynamic "message" {
                                for_each = chat.value.messages
                                content {
                                  role = message.value.role
                                  content {
                                    text = try(message.value.text, "") != "" ? message.value.text : null
                                    dynamic "cache_point" {
                                      for_each = try(message.value.cache_point, false) ? [1] : []
                                      content {
                                        type = "default"
                                      }
                                    }
                                  }
                                }
                              }

                              dynamic "system" {
                                for_each = try(chat.value.system, [])
                                content {
                                  text = try(system.value.text, "") != "" ? system.value.text : null
                                  dynamic "cache_point" {
                                    for_each = try(system.value.cache_point, false) ? [1] : []
                                    content {
                                      type = "default"
                                    }
                                  }
                                }
                              }

                              dynamic "input_variable" {
                                for_each = try(chat.value.input_variables, [])
                                content {
                                  name = input_variable.value
                                }
                              }

                              dynamic "tool_configuration" {
                                for_each = try(chat.value.tool_configuration, null) != null ? [chat.value.tool_configuration] : []
                                content {
                                  dynamic "tool" {
                                    for_each = tool_configuration.value.tools
                                    content {
                                      dynamic "cache_point" {
                                        for_each = try(tool.value.cache_point, false) ? [1] : []
                                        content {
                                          type = "default"
                                        }
                                      }
                                      dynamic "tool_spec" {
                                        for_each = try(tool.value.spec, null) != null ? [tool.value.spec] : []
                                        content {
                                          name        = tool_spec.value.name
                                          description = try(tool_spec.value.description, "") != "" ? tool_spec.value.description : null
                                          dynamic "input_schema" {
                                            for_each = try(tool_spec.value.input_schema, null) != null ? [tool_spec.value.input_schema] : []
                                            content {
                                              json = jsonencode(input_schema.value)
                                            }
                                          }
                                        }
                                      }
                                    }
                                  }

                                  # tool_choice is `any`-typed in variables.tf (its member is
                                  # literally named "any", an HCL type keyword) -- read with try().
                                  dynamic "tool_choice" {
                                    for_each = tool_configuration.value.tool_choice != null ? [tool_configuration.value.tool_choice] : []
                                    content {
                                      dynamic "any" {
                                        for_each = try(tool_choice.value.any, false) ? [1] : []
                                        content {}
                                      }
                                      dynamic "auto" {
                                        for_each = try(tool_choice.value.auto, false) ? [1] : []
                                        content {}
                                      }
                                      dynamic "tool" {
                                        for_each = try(tool_choice.value.tool_name, "") != "" ? [tool_choice.value.tool_name] : []
                                        content {
                                          name = tool.value
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
                    }
                  }
                }
              }
            }
          }
        }
      }

      dynamic "connection" {
        for_each = definition.value.connections
        content {
          name   = connection.value.name
          source = connection.value.source
          target = connection.value.target
          # Data or Conditional is derived from which arm is set (exactly
          # one, per the spec's CEL guard).
          type = connection.value.data != null ? "Data" : "Conditional"

          configuration {
            dynamic "data" {
              for_each = connection.value.data != null ? [connection.value.data] : []
              content {
                source_output = data.value.source_output
                target_input  = data.value.target_input
              }
            }
            dynamic "conditional" {
              for_each = connection.value.conditional != null ? [connection.value.conditional] : []
              content {
                condition = conditional.value.condition
              }
            }
          }
        }
      }
    }
  }

  tags = local.aws_tags
}
