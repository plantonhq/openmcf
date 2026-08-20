package module

import (
	"sort"

	"github.com/pkg/errors"
	awsbedrockagentcorememoryv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsbedrockagentcorememory/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// memory creates the AgentCore memory and its folded strategy satellites
// and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - strategy writes go through the parent memory's update API -- AWS
//     serializes them per memory and the provider holds a per-memory
//     lock, so strategy operations can take tens of minutes (the
//     provider's own timeout is 45m);
//   - the deprecated strategy-level memory_execution_role_arn and
//     namespaces arguments are never sent (excluded-deprecated; the
//     memory-level role and namespace_templates are the living
//     surfaces);
//   - MEMORY_RECORDS is the only stream content type AWS defines -- the
//     module owns the constant.
func memory(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &bedrock.AgentcoreMemoryArgs{
		// AWS's memory-name charset (letter first, then letters/digits/_)
		// is stricter than metadata.name conventions, so the name is an
		// explicit spec field. Changing it replaces the memory.
		Name: pulumi.String(spec.MemoryName),
		// Required by AWS: days raw session events survive (the
		// short-term window).
		EventExpiryDuration: pulumi.Int(int(spec.EventExpiryDays)),
		Tags:                pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	// Changing the key replaces the memory (provider-enforced).
	if spec.EncryptionKeyArn.GetValue() != "" {
		args.EncryptionKeyArn = pulumi.String(spec.EncryptionKeyArn.GetValue())
	}
	if spec.ExecutionRoleArn.GetValue() != "" {
		args.MemoryExecutionRoleArn = pulumi.String(spec.ExecutionRoleArn.GetValue())
	}

	// Metadata keys indexed for filtered retrieval (1-10; changing the
	// set replaces the memory, provider-enforced).
	var indexedKeys bedrock.AgentcoreMemoryIndexedKeyArray
	for _, k := range spec.IndexedKeys {
		indexedKeys = append(indexedKeys, &bedrock.AgentcoreMemoryIndexedKeyArgs{
			Key:  pulumi.String(k.Key),
			Type: pulumi.String(k.Type),
		})
	}
	if len(indexedKeys) > 0 {
		args.IndexedKeys = indexedKeys
	}

	// Stream long-term records to Kinesis as they are written.
	if spec.KinesisDelivery != nil {
		content := &bedrock.AgentcoreMemoryStreamDeliveryResourcesResourceKinesisContentConfigurationArgs{
			// MEMORY_RECORDS is the only content type AWS defines.
			Type: pulumi.String("MEMORY_RECORDS"),
		}
		if spec.KinesisDelivery.ContentLevel != "" {
			content.Level = pulumi.String(spec.KinesisDelivery.ContentLevel)
		}
		args.StreamDeliveryResources = &bedrock.AgentcoreMemoryStreamDeliveryResourcesArgs{
			Resource: &bedrock.AgentcoreMemoryStreamDeliveryResourcesResourceArgs{
				Kinesis: &bedrock.AgentcoreMemoryStreamDeliveryResourcesResourceKinesisArgs{
					DataStreamArn:        pulumi.String(spec.KinesisDelivery.DataStreamArn.GetValue()),
					ContentConfiguration: content,
				},
			},
		}
	}

	createdMemory, err := bedrock.NewAgentcoreMemory(ctx, spec.MemoryName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create memory")
	}

	ctx.Export(OpMemoryId, createdMemory.ID())
	ctx.Export(OpMemoryArn, createdMemory.Arn)

	// Strategies keyed by their stable entry names. Iteration is
	// name-sorted for deterministic previews; AWS serializes the writes
	// per memory regardless.
	strategyIds := pulumi.StringMap{}
	for _, s := range sortedStrategies(spec.Strategies) {
		strategyArgs := &bedrock.AgentcoreMemoryStrategyArgs{
			MemoryId: createdMemory.ID(),
			Name:     pulumi.String(s.Name),
			// Changing the type replaces the strategy
			// (provider-enforced).
			Type: pulumi.String(s.Type),
		}
		if s.Description != "" {
			strategyArgs.Description = pulumi.String(s.Description)
		}
		// Required by the provider: its deprecated `namespaces` twin and
		// this field are an exactly-one pair, and the living surface is
		// this one.
		strategyArgs.NamespaceTemplates = pulumi.ToStringArray(s.NamespaceTemplates)
		// Prompt/model overrides -- present exactly when type is CUSTOM
		// (spec-validated).
		if s.Custom != nil {
			configuration := &bedrock.AgentcoreMemoryStrategyConfigurationArgs{
				Type: pulumi.String(s.Custom.Type),
			}
			if s.Custom.Extraction != nil {
				configuration.Extraction = &bedrock.AgentcoreMemoryStrategyConfigurationExtractionArgs{
					AppendToPrompt: pulumi.String(s.Custom.Extraction.AppendToPrompt),
					ModelId:        pulumi.String(s.Custom.Extraction.ModelId),
				}
			}
			if s.Custom.Consolidation != nil {
				configuration.Consolidation = &bedrock.AgentcoreMemoryStrategyConfigurationConsolidationArgs{
					AppendToPrompt: pulumi.String(s.Custom.Consolidation.AppendToPrompt),
					ModelId:        pulumi.String(s.Custom.Consolidation.ModelId),
				}
			}
			if s.Custom.Reflection != nil {
				configuration.Reflection = &bedrock.AgentcoreMemoryStrategyConfigurationReflectionArgs{
					AppendToPrompt:     pulumi.String(s.Custom.Reflection.AppendToPrompt),
					ModelId:            pulumi.String(s.Custom.Reflection.ModelId),
					NamespaceTemplates: pulumi.ToStringArray(s.Custom.Reflection.NamespaceTemplates),
				}
			}
			strategyArgs.Configuration = configuration
		}
		// EPISODIC reflection namespaces -- only legal on EPISODIC
		// strategies (spec-validated).
		if len(s.ReflectionNamespaceTemplates) > 0 {
			strategyArgs.ReflectionConfiguration = &bedrock.AgentcoreMemoryStrategyReflectionConfigurationArgs{
				NamespaceTemplates: pulumi.ToStringArray(s.ReflectionNamespaceTemplates),
			}
		}
		createdStrategy, err := bedrock.NewAgentcoreMemoryStrategy(ctx, "strategy-"+s.Name, strategyArgs,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdMemory}))
		if err != nil {
			return errors.Wrapf(err, "create strategy %q", s.Name)
		}
		strategyIds[s.Name] = createdStrategy.MemoryStrategyId
	}
	ctx.Export(OpStrategyIds, strategyIds)

	return nil
}

func sortedStrategies(in []*awsbedrockagentcorememoryv1alpha1.AwsBedrockAgentCoreMemoryStrategy) []*awsbedrockagentcorememoryv1alpha1.AwsBedrockAgentCoreMemoryStrategy {
	out := append([]*awsbedrockagentcorememoryv1alpha1.AwsBedrockAgentCoreMemoryStrategy{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
