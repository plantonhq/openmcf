# Amazon Bedrock prompt (Prompt Management): a reusable prompt with one or
# more variants, each targeting a model or an agent alias with a text or
# chat template. The module owns the prompt's mutable DRAFT version.

resource "aws_bedrockagent_prompt" "this" {
  # Create-time naming basis; doubles as the Name tag. metadata.name on
  # both engines.
  name = local.prompt_name

  description                 = var.spec.description != "" ? var.spec.description : null
  customer_encryption_key_arn = var.spec.customer_encryption_key_arn != "" ? var.spec.customer_encryption_key_arn : null
  default_variant             = var.spec.default_variant != "" ? var.spec.default_variant : null

  dynamic "variant" {
    for_each = var.spec.variants
    content {
      name = variant.value.name
      # TEXT or CHAT is derived from which template arm is set (exactly
      # one, per the spec's CEL guard).
      template_type = try(variant.value.text, null) != null ? "TEXT" : "CHAT"

      model_id = try(variant.value.model_id, "") != "" ? variant.value.model_id : null

      additional_model_request_fields = try(variant.value.additional_model_request_fields, null) != null ? jsonencode(variant.value.additional_model_request_fields) : null

      # Execute through an agent alias instead of a model (exactly one of
      # the two, per the spec's CEL guard).
      dynamic "gen_ai_resource" {
        for_each = try(variant.value.agent_alias_arn, "") != "" ? [variant.value.agent_alias_arn] : []
        content {
          agent {
            agent_identifier = gen_ai_resource.value
          }
        }
      }

      # Bedrock stores temperature/top_p as float32 -- non-float32-exact
      # values (0.9) read back widened (0.8999999761581421). Applies are
      # unaffected (state keeps the config value); blind imports plan a
      # one-time reconcile on exactly those leaves (declared
      # write-normalized in the provider import catalog).
      dynamic "inference_configuration" {
        for_each = try(variant.value.inference_configuration, null) != null ? [variant.value.inference_configuration] : []
        content {
          text {
            max_tokens     = try(inference_configuration.value.max_tokens, null)
            stop_sequences = length(try(inference_configuration.value.stop_sequences, [])) > 0 ? inference_configuration.value.stop_sequences : null
            temperature    = try(inference_configuration.value.temperature, null)
            top_p          = try(inference_configuration.value.top_p, null)
          }
        }
      }

      dynamic "metadata" {
        for_each = try(variant.value.metadata, {})
        content {
          key   = metadata.key
          value = metadata.value
        }
      }

      template_configuration {

        dynamic "text" {
          for_each = try(variant.value.text, null) != null ? [variant.value.text] : []
          content {
            text = text.value.text
            dynamic "cache_point" {
              for_each = try(text.value.cache_point, false) ? [1] : []
              content {
                # "default" is the only cache point type AWS defines --
                # the spec models the checkpoint as a bool and the module
                # owns the constant.
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
          for_each = try(variant.value.chat, null) != null ? [variant.value.chat] : []
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

  tags = local.aws_tags
}
