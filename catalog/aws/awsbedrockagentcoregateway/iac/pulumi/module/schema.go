package module

import (
	"sort"

	awsbedrockagentcoregatewayv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockagentcoregateway/v1alpha1"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The provider bridges the lambda tool schema's three-level recursion
// with a DISTINCT Go type per tree position (input vs output schema,
// items-under-property vs items-under-root, ...), so the builders below
// are position-specific mirrors of one shape: type/description at every
// node, properties XOR items (spec-validated), raw-JSON leaves at the
// bottom -- exactly where AWS's own configuration surface bottoms out.

func inputSchemaArgs(in *awsbedrockagentcoregatewayv1alpha1.AwsBedrockAgentCoreGatewaySchemaDefinition) (*bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaArgs, error) {
	out := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaArgs{
		Type: pulumi.String(in.Type),
	}
	if in.Description != "" {
		out.Description = pulumi.String(in.Description)
	}
	var properties bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaPropertyArray
	for _, p := range in.Properties {
		property := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaPropertyArgs{
			Name:     pulumi.String(p.Name),
			Type:     pulumi.String(p.Type),
			Required: pulumi.Bool(p.Required),
		}
		if p.Description != "" {
			property.Description = pulumi.String(p.Description)
		}
		if p.Items != nil {
			items := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaPropertyItemsArgs{
				Type: pulumi.String(p.Items.Type),
			}
			if p.Items.Description != "" {
				items.Description = pulumi.String(p.Items.Description)
			}
			if p.Items.Items != nil {
				leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaPropertyItemsItemsArgs{
					Type: pulumi.String(p.Items.Items.Type),
				}
				if p.Items.Items.Description != "" {
					leaf.Description = pulumi.String(p.Items.Items.Description)
				}
				if p.Items.Items.ItemsJson != nil {
					jsonStr, err := structToJson(p.Items.Items.ItemsJson)
					if err != nil {
						return nil, err
					}
					leaf.ItemsJson = pulumi.String(jsonStr)
				}
				if p.Items.Items.PropertiesJson != nil {
					jsonStr, err := structToJson(p.Items.Items.PropertiesJson)
					if err != nil {
						return nil, err
					}
					leaf.PropertiesJson = pulumi.String(jsonStr)
				}
				items.Items = leaf
			}
			var leafProperties bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaPropertyItemsPropertyArray
			for _, lp := range p.Items.Properties {
				leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaPropertyItemsPropertyArgs{
					Name:     pulumi.String(lp.Name),
					Type:     pulumi.String(lp.Type),
					Required: pulumi.Bool(lp.Required),
				}
				if lp.Description != "" {
					leaf.Description = pulumi.String(lp.Description)
				}
				if lp.ItemsJson != nil {
					jsonStr, err := structToJson(lp.ItemsJson)
					if err != nil {
						return nil, err
					}
					leaf.ItemsJson = pulumi.String(jsonStr)
				}
				if lp.PropertiesJson != nil {
					jsonStr, err := structToJson(lp.PropertiesJson)
					if err != nil {
						return nil, err
					}
					leaf.PropertiesJson = pulumi.String(jsonStr)
				}
				leafProperties = append(leafProperties, leaf)
			}
			if len(leafProperties) > 0 {
				items.Properties = leafProperties
			}
			property.Items = items
		}
		var leafProperties bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaPropertyPropertyArray
		for _, lp := range p.Properties {
			leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaPropertyPropertyArgs{
				Name:     pulumi.String(lp.Name),
				Type:     pulumi.String(lp.Type),
				Required: pulumi.Bool(lp.Required),
			}
			if lp.Description != "" {
				leaf.Description = pulumi.String(lp.Description)
			}
			if lp.ItemsJson != nil {
				jsonStr, err := structToJson(lp.ItemsJson)
				if err != nil {
					return nil, err
				}
				leaf.ItemsJson = pulumi.String(jsonStr)
			}
			if lp.PropertiesJson != nil {
				jsonStr, err := structToJson(lp.PropertiesJson)
				if err != nil {
					return nil, err
				}
				leaf.PropertiesJson = pulumi.String(jsonStr)
			}
			leafProperties = append(leafProperties, leaf)
		}
		if len(leafProperties) > 0 {
			property.Properties = leafProperties
		}
		properties = append(properties, property)
	}
	if len(properties) > 0 {
		out.Properties = properties
	}
	if in.Items != nil {
		items := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaItemsArgs{
			Type: pulumi.String(in.Items.Type),
		}
		if in.Items.Description != "" {
			items.Description = pulumi.String(in.Items.Description)
		}
		if in.Items.Items != nil {
			leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaItemsItemsArgs{
				Type: pulumi.String(in.Items.Items.Type),
			}
			if in.Items.Items.Description != "" {
				leaf.Description = pulumi.String(in.Items.Items.Description)
			}
			if in.Items.Items.ItemsJson != nil {
				jsonStr, err := structToJson(in.Items.Items.ItemsJson)
				if err != nil {
					return nil, err
				}
				leaf.ItemsJson = pulumi.String(jsonStr)
			}
			if in.Items.Items.PropertiesJson != nil {
				jsonStr, err := structToJson(in.Items.Items.PropertiesJson)
				if err != nil {
					return nil, err
				}
				leaf.PropertiesJson = pulumi.String(jsonStr)
			}
			items.Items = leaf
		}
		var leafProperties bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaItemsPropertyArray
		for _, lp := range in.Items.Properties {
			leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadInputSchemaItemsPropertyArgs{
				Name:     pulumi.String(lp.Name),
				Type:     pulumi.String(lp.Type),
				Required: pulumi.Bool(lp.Required),
			}
			if lp.Description != "" {
				leaf.Description = pulumi.String(lp.Description)
			}
			if lp.ItemsJson != nil {
				jsonStr, err := structToJson(lp.ItemsJson)
				if err != nil {
					return nil, err
				}
				leaf.ItemsJson = pulumi.String(jsonStr)
			}
			if lp.PropertiesJson != nil {
				jsonStr, err := structToJson(lp.PropertiesJson)
				if err != nil {
					return nil, err
				}
				leaf.PropertiesJson = pulumi.String(jsonStr)
			}
			leafProperties = append(leafProperties, leaf)
		}
		if len(leafProperties) > 0 {
			items.Properties = leafProperties
		}
		out.Items = items
	}
	return out, nil
}

func outputSchemaArgs(in *awsbedrockagentcoregatewayv1alpha1.AwsBedrockAgentCoreGatewaySchemaDefinition) (*bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaArgs, error) {
	out := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaArgs{
		Type: pulumi.String(in.Type),
	}
	if in.Description != "" {
		out.Description = pulumi.String(in.Description)
	}
	var properties bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaPropertyArray
	for _, p := range in.Properties {
		property := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaPropertyArgs{
			Name:     pulumi.String(p.Name),
			Type:     pulumi.String(p.Type),
			Required: pulumi.Bool(p.Required),
		}
		if p.Description != "" {
			property.Description = pulumi.String(p.Description)
		}
		if p.Items != nil {
			items := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaPropertyItemsArgs{
				Type: pulumi.String(p.Items.Type),
			}
			if p.Items.Description != "" {
				items.Description = pulumi.String(p.Items.Description)
			}
			if p.Items.Items != nil {
				leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaPropertyItemsItemsArgs{
					Type: pulumi.String(p.Items.Items.Type),
				}
				if p.Items.Items.Description != "" {
					leaf.Description = pulumi.String(p.Items.Items.Description)
				}
				if p.Items.Items.ItemsJson != nil {
					jsonStr, err := structToJson(p.Items.Items.ItemsJson)
					if err != nil {
						return nil, err
					}
					leaf.ItemsJson = pulumi.String(jsonStr)
				}
				if p.Items.Items.PropertiesJson != nil {
					jsonStr, err := structToJson(p.Items.Items.PropertiesJson)
					if err != nil {
						return nil, err
					}
					leaf.PropertiesJson = pulumi.String(jsonStr)
				}
				items.Items = leaf
			}
			var leafProperties bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaPropertyItemsPropertyArray
			for _, lp := range p.Items.Properties {
				leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaPropertyItemsPropertyArgs{
					Name:     pulumi.String(lp.Name),
					Type:     pulumi.String(lp.Type),
					Required: pulumi.Bool(lp.Required),
				}
				if lp.Description != "" {
					leaf.Description = pulumi.String(lp.Description)
				}
				if lp.ItemsJson != nil {
					jsonStr, err := structToJson(lp.ItemsJson)
					if err != nil {
						return nil, err
					}
					leaf.ItemsJson = pulumi.String(jsonStr)
				}
				if lp.PropertiesJson != nil {
					jsonStr, err := structToJson(lp.PropertiesJson)
					if err != nil {
						return nil, err
					}
					leaf.PropertiesJson = pulumi.String(jsonStr)
				}
				leafProperties = append(leafProperties, leaf)
			}
			if len(leafProperties) > 0 {
				items.Properties = leafProperties
			}
			property.Items = items
		}
		var leafProperties bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaPropertyPropertyArray
		for _, lp := range p.Properties {
			leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaPropertyPropertyArgs{
				Name:     pulumi.String(lp.Name),
				Type:     pulumi.String(lp.Type),
				Required: pulumi.Bool(lp.Required),
			}
			if lp.Description != "" {
				leaf.Description = pulumi.String(lp.Description)
			}
			if lp.ItemsJson != nil {
				jsonStr, err := structToJson(lp.ItemsJson)
				if err != nil {
					return nil, err
				}
				leaf.ItemsJson = pulumi.String(jsonStr)
			}
			if lp.PropertiesJson != nil {
				jsonStr, err := structToJson(lp.PropertiesJson)
				if err != nil {
					return nil, err
				}
				leaf.PropertiesJson = pulumi.String(jsonStr)
			}
			leafProperties = append(leafProperties, leaf)
		}
		if len(leafProperties) > 0 {
			property.Properties = leafProperties
		}
		properties = append(properties, property)
	}
	if len(properties) > 0 {
		out.Properties = properties
	}
	if in.Items != nil {
		items := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaItemsArgs{
			Type: pulumi.String(in.Items.Type),
		}
		if in.Items.Description != "" {
			items.Description = pulumi.String(in.Items.Description)
		}
		if in.Items.Items != nil {
			leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaItemsItemsArgs{
				Type: pulumi.String(in.Items.Items.Type),
			}
			if in.Items.Items.Description != "" {
				leaf.Description = pulumi.String(in.Items.Items.Description)
			}
			if in.Items.Items.ItemsJson != nil {
				jsonStr, err := structToJson(in.Items.Items.ItemsJson)
				if err != nil {
					return nil, err
				}
				leaf.ItemsJson = pulumi.String(jsonStr)
			}
			if in.Items.Items.PropertiesJson != nil {
				jsonStr, err := structToJson(in.Items.Items.PropertiesJson)
				if err != nil {
					return nil, err
				}
				leaf.PropertiesJson = pulumi.String(jsonStr)
			}
			items.Items = leaf
		}
		var leafProperties bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaItemsPropertyArray
		for _, lp := range in.Items.Properties {
			leaf := &bedrock.AgentcoreGatewayTargetTargetConfigurationMcpLambdaToolSchemaInlinePayloadOutputSchemaItemsPropertyArgs{
				Name:     pulumi.String(lp.Name),
				Type:     pulumi.String(lp.Type),
				Required: pulumi.Bool(lp.Required),
			}
			if lp.Description != "" {
				leaf.Description = pulumi.String(lp.Description)
			}
			if lp.ItemsJson != nil {
				jsonStr, err := structToJson(lp.ItemsJson)
				if err != nil {
					return nil, err
				}
				leaf.ItemsJson = pulumi.String(jsonStr)
			}
			if lp.PropertiesJson != nil {
				jsonStr, err := structToJson(lp.PropertiesJson)
				if err != nil {
					return nil, err
				}
				leaf.PropertiesJson = pulumi.String(jsonStr)
			}
			leafProperties = append(leafProperties, leaf)
		}
		if len(leafProperties) > 0 {
			items.Properties = leafProperties
		}
		out.Items = items
	}
	return out, nil
}

func svrsToStringArray(in []*foreignkeyv1.StringValueOrRef) pulumi.StringArray {
	out := pulumi.StringArray{}
	for _, ref := range in {
		out = append(out, pulumi.String(ref.GetValue()))
	}
	return out
}

func sortedTargets(in []*awsbedrockagentcoregatewayv1alpha1.AwsBedrockAgentCoreGatewayTarget) []*awsbedrockagentcoregatewayv1alpha1.AwsBedrockAgentCoreGatewayTarget {
	out := append([]*awsbedrockagentcoregatewayv1alpha1.AwsBedrockAgentCoreGatewayTarget{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
