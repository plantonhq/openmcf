package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	awsbedrockflowv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockflow/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// flow creates the Bedrock flow's mutable DRAFT definition and exports
// outputs.
//
// AWS validates the GRAPH server-side (unreachable nodes, type
// mismatches, missing connections) at create/update -- the spec validates
// shapes, the service validates topology.
func flow(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &bedrock.AgentFlowArgs{
		// Create-time naming basis; doubles as the Name tag. metadata.name
		// on both engines.
		Name:             pulumi.String(locals.FlowName),
		ExecutionRoleArn: pulumi.String(spec.ExecutionRoleArn.GetValue()),
		Tags:             pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.CustomerEncryptionKeyArn.GetValue() != "" {
		args.CustomerEncryptionKeyArn = pulumi.String(spec.CustomerEncryptionKeyArn.GetValue())
	}

	if spec.Definition != nil {
		definition := &bedrock.AgentFlowDefinitionArgs{}

		var nodes bedrock.AgentFlowDefinitionNodeArray
		for _, n := range spec.Definition.Nodes {
			node, err := flowNodeArgs(n)
			if err != nil {
				return errors.Wrapf(err, "render node %q", n.Name)
			}
			nodes = append(nodes, node)
		}
		definition.Nodes = nodes

		var connections bedrock.AgentFlowDefinitionConnectionArray
		for _, c := range spec.Definition.Connections {
			connection := &bedrock.AgentFlowDefinitionConnectionArgs{
				Name:   pulumi.String(c.Name),
				Source: pulumi.String(c.Source),
				Target: pulumi.String(c.Target),
			}
			// Data or Conditional is derived from which arm is set
			// (exactly one, per the spec's CEL guard).
			configuration := &bedrock.AgentFlowDefinitionConnectionConfigurationArgs{}
			if c.Data != nil {
				connection.Type = pulumi.String("Data")
				configuration.Data = &bedrock.AgentFlowDefinitionConnectionConfigurationDataArgs{
					SourceOutput: pulumi.String(c.Data.SourceOutput),
					TargetInput:  pulumi.String(c.Data.TargetInput),
				}
			}
			if c.Conditional != nil {
				connection.Type = pulumi.String("Conditional")
				configuration.Conditional = &bedrock.AgentFlowDefinitionConnectionConfigurationConditionalArgs{
					Condition: pulumi.String(c.Conditional.Condition),
				}
			}
			connection.Configuration = configuration
			connections = append(connections, connection)
		}
		definition.Connections = connections

		args.Definition = definition
	}

	createdFlow, err := bedrock.NewAgentFlow(ctx, locals.FlowName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create flow")
	}

	ctx.Export(OpFlowId, createdFlow.ID())
	ctx.Export(OpFlowArn, createdFlow.Arn)
	ctx.Export(OpDraftVersion, pulumi.String("DRAFT"))

	return nil
}

// flowNodeArgs renders one node -- inputs, outputs, and exactly the AWS
// configuration union member matching the node class. Structural classes
// (Input, Output, Iterator, Collector) render an EMPTY member; the Loop
// family renders no configuration at all (not expressible at the pinned
// provider -- an upstream gap).
func flowNodeArgs(n *awsbedrockflowv1alpha1.AwsBedrockFlowNode) (*bedrock.AgentFlowDefinitionNodeArgs, error) {
	node := &bedrock.AgentFlowDefinitionNodeArgs{
		Name: pulumi.String(n.Name),
		Type: pulumi.String(n.Type),
	}

	var inputs bedrock.AgentFlowDefinitionNodeInputTypeArray
	for _, in := range n.Inputs {
		input := &bedrock.AgentFlowDefinitionNodeInputTypeArgs{
			Name:       pulumi.String(in.Name),
			Expression: pulumi.String(in.Expression),
			Type:       pulumi.String(in.Type),
		}
		if in.Category != "" {
			input.Category = pulumi.String(in.Category)
		}
		inputs = append(inputs, input)
	}
	node.Inputs = inputs

	var outputs bedrock.AgentFlowDefinitionNodeOutputTypeArray
	for _, out := range n.Outputs {
		outputs = append(outputs, &bedrock.AgentFlowDefinitionNodeOutputTypeArgs{
			Name: pulumi.String(out.Name),
			Type: pulumi.String(out.Type),
		})
	}
	node.Outputs = outputs

	configuration := &bedrock.AgentFlowDefinitionNodeConfigurationArgs{}
	configured := true

	switch n.Type {
	case "Input":
		configuration.Input = &bedrock.AgentFlowDefinitionNodeConfigurationInputTypeArgs{}
	case "Output":
		configuration.Output = &bedrock.AgentFlowDefinitionNodeConfigurationOutputTypeArgs{}
	case "Iterator":
		configuration.Iterator = &bedrock.AgentFlowDefinitionNodeConfigurationIteratorArgs{}
	case "Collector":
		configuration.Collector = &bedrock.AgentFlowDefinitionNodeConfigurationCollectorArgs{}
	case "Agent":
		configuration.Agent = &bedrock.AgentFlowDefinitionNodeConfigurationAgentArgs{
			AgentAliasArn: pulumi.String(n.Agent.AgentAliasArn.GetValue()),
		}
	case "LambdaFunction":
		configuration.LambdaFunction = &bedrock.AgentFlowDefinitionNodeConfigurationLambdaFunctionArgs{
			LambdaArn: pulumi.String(n.LambdaFunction.LambdaArn.GetValue()),
		}
	case "Lex":
		configuration.Lex = &bedrock.AgentFlowDefinitionNodeConfigurationLexArgs{
			BotAliasArn: pulumi.String(n.Lex.BotAliasArn),
			LocaleId:    pulumi.String(n.Lex.LocaleId),
		}
	case "InlineCode":
		configuration.InlineCode = &bedrock.AgentFlowDefinitionNodeConfigurationInlineCodeArgs{
			Code: pulumi.String(n.InlineCode.Code),
			// Python_3 is the only language AWS defines -- the module owns
			// the constant.
			Language: pulumi.String("Python_3"),
		}
	case "Condition":
		var conditions bedrock.AgentFlowDefinitionNodeConfigurationConditionConditionArray
		for _, c := range n.Condition.Conditions {
			condition := &bedrock.AgentFlowDefinitionNodeConfigurationConditionConditionArgs{
				Name: pulumi.String(c.Name),
			}
			if c.Expression != "" {
				condition.Expression = pulumi.String(c.Expression)
			}
			conditions = append(conditions, condition)
		}
		configuration.Condition = &bedrock.AgentFlowDefinitionNodeConfigurationConditionArgs{
			Conditions: conditions,
		}
	case "KnowledgeBase":
		kb := &bedrock.AgentFlowDefinitionNodeConfigurationKnowledgeBaseArgs{
			KnowledgeBaseId: pulumi.String(n.KnowledgeBase.KnowledgeBaseId.GetValue()),
		}
		if n.KnowledgeBase.ModelId != "" {
			kb.ModelId = pulumi.String(n.KnowledgeBase.ModelId)
		}
		if n.KnowledgeBase.NumberOfResults != 0 {
			kb.NumberOfResults = pulumi.Int(int(n.KnowledgeBase.NumberOfResults))
		}
		if n.KnowledgeBase.Guardrail != nil {
			kb.GuardrailConfiguration = &bedrock.AgentFlowDefinitionNodeConfigurationKnowledgeBaseGuardrailConfigurationArgs{
				GuardrailIdentifier: pulumi.String(n.KnowledgeBase.Guardrail.GuardrailId.GetValue()),
				GuardrailVersion:    pulumi.String(n.KnowledgeBase.Guardrail.Version),
			}
		}
		if n.KnowledgeBase.InferenceConfiguration != nil {
			kb.InferenceConfiguration = &bedrock.AgentFlowDefinitionNodeConfigurationKnowledgeBaseInferenceConfigurationArgs{
				Text: knowledgeBaseInferenceTextArgs(n.KnowledgeBase.InferenceConfiguration),
			}
		}
		configuration.KnowledgeBase = kb
	case "Retrieval":
		configuration.Retrieval = &bedrock.AgentFlowDefinitionNodeConfigurationRetrievalArgs{
			ServiceConfiguration: &bedrock.AgentFlowDefinitionNodeConfigurationRetrievalServiceConfigurationArgs{
				// S3 is the only retrieval service AWS defines.
				S3: &bedrock.AgentFlowDefinitionNodeConfigurationRetrievalServiceConfigurationS3Args{
					BucketName: pulumi.String(n.Retrieval.BucketName.GetValue()),
				},
			},
		}
	case "Storage":
		configuration.Storage = &bedrock.AgentFlowDefinitionNodeConfigurationStorageArgs{
			ServiceConfiguration: &bedrock.AgentFlowDefinitionNodeConfigurationStorageServiceConfigurationArgs{
				// S3 is the only storage service AWS defines.
				S3: &bedrock.AgentFlowDefinitionNodeConfigurationStorageServiceConfigurationS3Args{
					BucketName: pulumi.String(n.Storage.BucketName.GetValue()),
				},
			},
		}
	case "Prompt":
		prompt, err := flowPromptNodeArgs(n.Prompt)
		if err != nil {
			return nil, err
		}
		configuration.Prompt = prompt
	default:
		// The Loop family: accepted by AWS, no configuration member
		// expressible at the pinned provider.
		configured = false
	}

	if configured {
		node.Configuration = configuration
	}
	return node, nil
}

// flowPromptNodeArgs renders the prompt node -- a Prompt Management
// resource or an inline template (mirrors the AwsBedrockPrompt module
// arm-for-arm).
func flowPromptNodeArgs(p *awsbedrockflowv1alpha1.AwsBedrockFlowPromptNode) (*bedrock.AgentFlowDefinitionNodeConfigurationPromptArgs, error) {
	prompt := &bedrock.AgentFlowDefinitionNodeConfigurationPromptArgs{}

	if p.Guardrail != nil {
		prompt.GuardrailConfiguration = &bedrock.AgentFlowDefinitionNodeConfigurationPromptGuardrailConfigurationArgs{
			GuardrailIdentifier: pulumi.String(p.Guardrail.GuardrailId.GetValue()),
			GuardrailVersion:    pulumi.String(p.Guardrail.Version),
		}
	}

	source := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationArgs{}
	if p.PromptArn.GetValue() != "" {
		source.Resource = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationResourceArgs{
			PromptArn: pulumi.String(p.PromptArn.GetValue()),
		}
	}
	if p.Inline != nil {
		inline := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineArgs{
			ModelId: pulumi.String(p.Inline.ModelId),
		}
		// TEXT or CHAT is derived from which template arm is set (exactly
		// one, per the spec's CEL guard).
		if p.Inline.Text != nil {
			inline.TemplateType = pulumi.String("TEXT")
		} else {
			inline.TemplateType = pulumi.String("CHAT")
		}

		if p.Inline.AdditionalModelRequestFields != nil {
			fieldsJson, err := json.Marshal(p.Inline.AdditionalModelRequestFields.AsMap())
			if err != nil {
				return nil, errors.Wrap(err, "marshal additional model request fields")
			}
			inline.AdditionalModelRequestFields = pulumi.String(string(fieldsJson))
		}

		if p.Inline.InferenceConfiguration != nil {
			inline.InferenceConfiguration = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineInferenceConfigurationArgs{
				Text: inlineInferenceTextArgs(p.Inline.InferenceConfiguration),
			}
		}

		template := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationArgs{}
		if p.Inline.Text != nil {
			text := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationTextArgs{
				Text: pulumi.String(p.Inline.Text.Text),
			}
			if p.Inline.Text.CachePoint {
				// "default" is the only cache point type AWS defines.
				text.CachePoint = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationTextCachePointArgs{
					Type: pulumi.String("default"),
				}
			}
			var variables bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationTextInputVariableArray
			for _, v := range p.Inline.Text.InputVariables {
				variables = append(variables, &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationTextInputVariableArgs{
					Name: pulumi.String(v),
				})
			}
			text.InputVariables = variables
			template.Text = text
		}
		if p.Inline.Chat != nil {
			chat := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatArgs{}

			var messages bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatMessageArray
			for _, m := range p.Inline.Chat.Messages {
				message := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatMessageArgs{
					Role: pulumi.String(m.Role),
				}
				content := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatMessageContentArgs{}
				if m.Text != "" {
					content.Text = pulumi.String(m.Text)
				}
				if m.CachePoint {
					content.CachePoint = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatMessageContentCachePointArgs{
						Type: pulumi.String("default"),
					}
				}
				message.Content = content
				messages = append(messages, message)
			}
			chat.Messages = messages

			var systems bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatSystemArray
			for _, s := range p.Inline.Chat.System {
				system := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatSystemArgs{}
				if s.Text != "" {
					system.Text = pulumi.String(s.Text)
				}
				if s.CachePoint {
					system.CachePoint = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatSystemCachePointArgs{
						Type: pulumi.String("default"),
					}
				}
				systems = append(systems, system)
			}
			chat.Systems = systems

			var variables bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatInputVariableArray
			for _, v := range p.Inline.Chat.InputVariables {
				variables = append(variables, &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatInputVariableArgs{
					Name: pulumi.String(v),
				})
			}
			chat.InputVariables = variables

			if p.Inline.Chat.ToolConfiguration != nil {
				toolConfiguration := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationArgs{}
				var tools bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationToolArray
				for _, t := range p.Inline.Chat.ToolConfiguration.Tools {
					tool := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationToolArgs{}
					if t.CachePoint {
						tool.CachePoint = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationToolCachePointArgs{
							Type: pulumi.String("default"),
						}
					}
					if t.Spec != nil {
						toolSpec := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationToolToolSpecArgs{
							Name: pulumi.String(t.Spec.Name),
						}
						if t.Spec.Description != "" {
							toolSpec.Description = pulumi.String(t.Spec.Description)
						}
						if t.Spec.InputSchema != nil {
							schemaJson, err := json.Marshal(t.Spec.InputSchema.AsMap())
							if err != nil {
								return nil, errors.Wrap(err, "marshal tool input schema")
							}
							toolSpec.InputSchema = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationToolToolSpecInputSchemaArgs{
								Json: pulumi.String(string(schemaJson)),
							}
						}
						tool.ToolSpec = toolSpec
					}
					tools = append(tools, tool)
				}
				toolConfiguration.Tools = tools

				if p.Inline.Chat.ToolConfiguration.ToolChoice != nil {
					choice := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationToolChoiceArgs{}
					if p.Inline.Chat.ToolConfiguration.ToolChoice.Any {
						choice.Any = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationToolChoiceAnyArgs{}
					}
					if p.Inline.Chat.ToolConfiguration.ToolChoice.Auto {
						choice.Auto = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationToolChoiceAutoArgs{}
					}
					if p.Inline.Chat.ToolConfiguration.ToolChoice.ToolName != "" {
						choice.Tool = &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineTemplateConfigurationChatToolConfigurationToolChoiceToolArgs{
							Name: pulumi.String(p.Inline.Chat.ToolConfiguration.ToolChoice.ToolName),
						}
					}
					toolConfiguration.ToolChoice = choice
				}
				chat.ToolConfiguration = toolConfiguration
			}
			template.Chat = chat
		}
		inline.TemplateConfiguration = template
		source.Inline = inline
	}
	prompt.SourceConfiguration = source

	return prompt, nil
}

// Bedrock stores temperature/top_p as float32 -- non-float32-exact
// values (0.2) read back widened (0.20000000298023224). Applies are
// unaffected (state keeps the config value); blind imports plan a
// one-time reconcile on exactly those leaves (declared write-normalized
// in the provider import catalog). Same class as AwsBedrockPrompt.
func knowledgeBaseInferenceTextArgs(c *awsbedrockflowv1alpha1.AwsBedrockFlowInferenceConfiguration) *bedrock.AgentFlowDefinitionNodeConfigurationKnowledgeBaseInferenceConfigurationTextArgs {
	text := &bedrock.AgentFlowDefinitionNodeConfigurationKnowledgeBaseInferenceConfigurationTextArgs{}
	if c.MaxTokens != nil {
		text.MaxTokens = pulumi.Int(int(*c.MaxTokens))
	}
	if len(c.StopSequences) > 0 {
		text.StopSequences = pulumi.ToStringArray(c.StopSequences)
	}
	if c.Temperature != nil {
		text.Temperature = pulumi.Float64(*c.Temperature)
	}
	if c.TopP != nil {
		text.TopP = pulumi.Float64(*c.TopP)
	}
	return text
}

// Same float32 class as knowledgeBaseInferenceTextArgs (and AwsBedrockPrompt).
func inlineInferenceTextArgs(c *awsbedrockflowv1alpha1.AwsBedrockFlowInferenceConfiguration) *bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineInferenceConfigurationTextArgs {
	text := &bedrock.AgentFlowDefinitionNodeConfigurationPromptSourceConfigurationInlineInferenceConfigurationTextArgs{}
	if c.MaxTokens != nil {
		text.MaxTokens = pulumi.Int(int(*c.MaxTokens))
	}
	if len(c.StopSequences) > 0 {
		text.StopSequences = pulumi.ToStringArray(c.StopSequences)
	}
	if c.Temperature != nil {
		text.Temperature = pulumi.Float64(*c.Temperature)
	}
	if c.TopP != nil {
		text.TopP = pulumi.Float64(*c.TopP)
	}
	return text
}
