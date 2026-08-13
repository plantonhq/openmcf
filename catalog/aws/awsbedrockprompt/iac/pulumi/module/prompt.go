package module

import (
	"encoding/json"
	"sort"

	"github.com/pkg/errors"
	awsbedrockpromptv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockprompt/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// prompt creates the Bedrock prompt's mutable DRAFT version with its
// variants and exports outputs.
func prompt(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &bedrock.AgentPromptArgs{
		// Create-time naming basis; doubles as the Name tag. metadata.name
		// on both engines.
		Name: pulumi.String(locals.PromptName),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.CustomerEncryptionKeyArn.GetValue() != "" {
		args.CustomerEncryptionKeyArn = pulumi.String(spec.CustomerEncryptionKeyArn.GetValue())
	}
	if spec.DefaultVariant != "" {
		args.DefaultVariant = pulumi.String(spec.DefaultVariant)
	}

	// Variants iterate name-sorted for deterministic previews.
	var variants bedrock.AgentPromptVariantArray
	for _, v := range sortedVariants(spec.Variants) {
		variant, err := variantArgs(v)
		if err != nil {
			return errors.Wrapf(err, "render variant %q", v.Name)
		}
		variants = append(variants, variant)
	}
	args.Variants = variants

	createdPrompt, err := bedrock.NewAgentPrompt(ctx, locals.PromptName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create prompt")
	}

	ctx.Export(OpPromptId, createdPrompt.ID())
	ctx.Export(OpPromptArn, createdPrompt.Arn)
	ctx.Export(OpDraftVersion, pulumi.String("DRAFT"))

	return nil
}

// variantArgs renders one variant -- the execution target, the template
// arm (with AWS's template_type derived), tools, and inference settings.
func variantArgs(v *awsbedrockpromptv1alpha1.AwsBedrockPromptVariant) (*bedrock.AgentPromptVariantArgs, error) {
	variant := &bedrock.AgentPromptVariantArgs{
		Name: pulumi.String(v.Name),
	}

	// TEXT or CHAT is derived from which template arm is set (exactly
	// one, per the spec's CEL guard) -- as is the model-vs-agent target.
	if v.Text != nil {
		variant.TemplateType = pulumi.String("TEXT")
	} else {
		variant.TemplateType = pulumi.String("CHAT")
	}

	if v.ModelId != "" {
		variant.ModelId = pulumi.String(v.ModelId)
	}
	if v.AgentAliasArn.GetValue() != "" {
		variant.GenAiResource = &bedrock.AgentPromptVariantGenAiResourceArgs{
			Agent: &bedrock.AgentPromptVariantGenAiResourceAgentArgs{
				AgentIdentifier: pulumi.String(v.AgentAliasArn.GetValue()),
			},
		}
	}

	if v.AdditionalModelRequestFields != nil {
		fieldsJson, err := json.Marshal(v.AdditionalModelRequestFields.AsMap())
		if err != nil {
			return nil, errors.Wrap(err, "marshal additional model request fields")
		}
		variant.AdditionalModelRequestFields = pulumi.String(string(fieldsJson))
	}

	if v.InferenceConfiguration != nil {
		text := &bedrock.AgentPromptVariantInferenceConfigurationTextArgs{}
		if v.InferenceConfiguration.MaxTokens != nil {
			text.MaxTokens = pulumi.Int(int(*v.InferenceConfiguration.MaxTokens))
		}
		if len(v.InferenceConfiguration.StopSequences) > 0 {
			text.StopSequences = pulumi.ToStringArray(v.InferenceConfiguration.StopSequences)
		}
		if v.InferenceConfiguration.Temperature != nil {
			text.Temperature = pulumi.Float64(*v.InferenceConfiguration.Temperature)
		}
		if v.InferenceConfiguration.TopP != nil {
			text.TopP = pulumi.Float64(*v.InferenceConfiguration.TopP)
		}
		variant.InferenceConfiguration = &bedrock.AgentPromptVariantInferenceConfigurationArgs{
			Text: text,
		}
	}

	if len(v.Metadata) > 0 {
		var metadataEntries bedrock.AgentPromptVariantMetadataArray
		keys := make([]string, 0, len(v.Metadata))
		for k := range v.Metadata {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			metadataEntries = append(metadataEntries, &bedrock.AgentPromptVariantMetadataArgs{
				Key:   pulumi.String(k),
				Value: pulumi.String(v.Metadata[k]),
			})
		}
		variant.Metadatas = metadataEntries
	}

	template := &bedrock.AgentPromptVariantTemplateConfigurationArgs{}
	if v.Text != nil {
		text := &bedrock.AgentPromptVariantTemplateConfigurationTextArgs{
			Text: pulumi.String(v.Text.Text),
		}
		if v.Text.CachePoint {
			// "default" is the only cache point type AWS defines.
			text.CachePoint = &bedrock.AgentPromptVariantTemplateConfigurationTextCachePointArgs{
				Type: pulumi.String("default"),
			}
		}
		var variables bedrock.AgentPromptVariantTemplateConfigurationTextInputVariableArray
		for _, name := range v.Text.InputVariables {
			variables = append(variables, &bedrock.AgentPromptVariantTemplateConfigurationTextInputVariableArgs{
				Name: pulumi.String(name),
			})
		}
		text.InputVariables = variables
		template.Text = text
	}
	if v.Chat != nil {
		chat := &bedrock.AgentPromptVariantTemplateConfigurationChatArgs{}

		var messages bedrock.AgentPromptVariantTemplateConfigurationChatMessageArray
		for _, m := range v.Chat.Messages {
			message := &bedrock.AgentPromptVariantTemplateConfigurationChatMessageArgs{
				Role: pulumi.String(m.Role),
			}
			content := &bedrock.AgentPromptVariantTemplateConfigurationChatMessageContentArgs{}
			if m.Text != "" {
				content.Text = pulumi.String(m.Text)
			}
			if m.CachePoint {
				content.CachePoint = &bedrock.AgentPromptVariantTemplateConfigurationChatMessageContentCachePointArgs{
					Type: pulumi.String("default"),
				}
			}
			message.Content = content
			messages = append(messages, message)
		}
		chat.Messages = messages

		var systems bedrock.AgentPromptVariantTemplateConfigurationChatSystemArray
		for _, s := range v.Chat.System {
			system := &bedrock.AgentPromptVariantTemplateConfigurationChatSystemArgs{}
			if s.Text != "" {
				system.Text = pulumi.String(s.Text)
			}
			if s.CachePoint {
				system.CachePoint = &bedrock.AgentPromptVariantTemplateConfigurationChatSystemCachePointArgs{
					Type: pulumi.String("default"),
				}
			}
			systems = append(systems, system)
		}
		chat.Systems = systems

		var variables bedrock.AgentPromptVariantTemplateConfigurationChatInputVariableArray
		for _, name := range v.Chat.InputVariables {
			variables = append(variables, &bedrock.AgentPromptVariantTemplateConfigurationChatInputVariableArgs{
				Name: pulumi.String(name),
			})
		}
		chat.InputVariables = variables

		if v.Chat.ToolConfiguration != nil {
			toolConfiguration := &bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationArgs{}
			var tools bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationToolArray
			for _, t := range v.Chat.ToolConfiguration.Tools {
				tool := &bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationToolArgs{}
				if t.CachePoint {
					tool.CachePoint = &bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationToolCachePointArgs{
						Type: pulumi.String("default"),
					}
				}
				if t.Spec != nil {
					toolSpec := &bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationToolToolSpecArgs{
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
						toolSpec.InputSchema = &bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationToolToolSpecInputSchemaArgs{
							Json: pulumi.String(string(schemaJson)),
						}
					}
					tool.ToolSpec = toolSpec
				}
				tools = append(tools, tool)
			}
			toolConfiguration.Tools = tools

			if v.Chat.ToolConfiguration.ToolChoice != nil {
				choice := &bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationToolChoiceArgs{}
				if v.Chat.ToolConfiguration.ToolChoice.Any {
					choice.Any = &bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationToolChoiceAnyArgs{}
				}
				if v.Chat.ToolConfiguration.ToolChoice.Auto {
					choice.Auto = &bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationToolChoiceAutoArgs{}
				}
				if v.Chat.ToolConfiguration.ToolChoice.ToolName != "" {
					choice.Tool = &bedrock.AgentPromptVariantTemplateConfigurationChatToolConfigurationToolChoiceToolArgs{
						Name: pulumi.String(v.Chat.ToolConfiguration.ToolChoice.ToolName),
					}
				}
				toolConfiguration.ToolChoice = choice
			}
			chat.ToolConfiguration = toolConfiguration
		}
		template.Chat = chat
	}
	variant.TemplateConfiguration = template

	return variant, nil
}

func sortedVariants(in []*awsbedrockpromptv1alpha1.AwsBedrockPromptVariant) []*awsbedrockpromptv1alpha1.AwsBedrockPromptVariant {
	out := append([]*awsbedrockpromptv1alpha1.AwsBedrockPromptVariant{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
