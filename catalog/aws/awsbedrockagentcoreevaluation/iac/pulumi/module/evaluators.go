package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// evaluators creates the bundle's scoring definitions and exports
// outputs. It returns the created resources keyed by entry name so the
// online-config renders can resolve in-bundle evaluator references.
//
// Lifecycle facts the renders below depend on:
//   - the evaluator name is its identity seed (AWS derives
//     "<name>-<10 chars>" from it) and has no rename -- the provider
//     replaces on change;
//   - creation waits CREATING -> ACTIVE; Lambda-backed evaluators retry
//     briefly upstream while IAM/Lambda permissions propagate;
//   - an evaluator in use by an active online evaluation is locked for
//     modification (AWS returns conflicts; upstream retries deletes).
func evaluators(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (map[string]*createdEvaluator, error) {
	spec := locals.Spec

	created := map[string]*createdEvaluator{}
	evaluatorIds := pulumi.StringMap{}
	evaluatorArns := pulumi.StringMap{}

	// Iteration is name-sorted for deterministic previews.
	for _, e := range sortedEvaluators(spec.Evaluators) {
		config := &bedrock.AgentcoreEvaluatorEvaluatorConfigArgs{}

		if e.LlmAsAJudge != nil {
			judge := e.LlmAsAJudge
			model := &bedrock.AgentcoreEvaluatorEvaluatorConfigLlmAsAJudgeModelConfigBedrockEvaluatorModelConfigArgs{
				ModelId: pulumi.String(judge.Model.ModelId),
			}
			// Model-specific passthrough fields travel as the JSON the
			// model's own API documents.
			if judge.Model.AdditionalModelRequestFields != nil {
				extra, err := json.Marshal(judge.Model.AdditionalModelRequestFields.AsMap())
				if err != nil {
					return nil, errors.Wrapf(err, "marshal additional model request fields for evaluator %q", e.Name)
				}
				model.AdditionalModelRequestFields = pulumi.String(string(extra))
			}
			if judge.Model.Inference != nil {
				inference := &bedrock.AgentcoreEvaluatorEvaluatorConfigLlmAsAJudgeModelConfigBedrockEvaluatorModelConfigInferenceConfigArgs{}
				if judge.Model.Inference.MaxTokens > 0 {
					inference.MaxTokens = pulumi.Int(int(judge.Model.Inference.MaxTokens))
				}
				if len(judge.Model.Inference.StopSequences) > 0 {
					inference.StopSequences = pulumi.ToStringArray(judge.Model.Inference.StopSequences)
				}
				// Presence-typed: only an explicit choice is sent, so the
				// module never fights the model's own defaults.
				if judge.Model.Inference.Temperature != nil {
					inference.Temperature = pulumi.Float64(*judge.Model.Inference.Temperature)
				}
				if judge.Model.Inference.TopP != nil {
					inference.TopP = pulumi.Float64(*judge.Model.Inference.TopP)
				}
				model.InferenceConfig = inference
			}

			// Exactly one scale shape (spec-validated).
			scale := &bedrock.AgentcoreEvaluatorEvaluatorConfigLlmAsAJudgeRatingScaleArgs{}
			if len(judge.RatingScale.Categorical) > 0 {
				var categoricals bedrock.AgentcoreEvaluatorEvaluatorConfigLlmAsAJudgeRatingScaleCategoricalArray
				for _, c := range judge.RatingScale.Categorical {
					categoricals = append(categoricals, &bedrock.AgentcoreEvaluatorEvaluatorConfigLlmAsAJudgeRatingScaleCategoricalArgs{
						Label:      pulumi.String(c.Label),
						Definition: pulumi.String(c.Definition),
					})
				}
				scale.Categoricals = categoricals
			}
			if len(judge.RatingScale.Numerical) > 0 {
				var numericals bedrock.AgentcoreEvaluatorEvaluatorConfigLlmAsAJudgeRatingScaleNumericalArray
				for _, n := range judge.RatingScale.Numerical {
					numericals = append(numericals, &bedrock.AgentcoreEvaluatorEvaluatorConfigLlmAsAJudgeRatingScaleNumericalArgs{
						Label:      pulumi.String(n.Label),
						Definition: pulumi.String(n.Definition),
						Value:      pulumi.Float64(n.Value),
					})
				}
				scale.Numericals = numericals
			}

			config.LlmAsAJudge = &bedrock.AgentcoreEvaluatorEvaluatorConfigLlmAsAJudgeArgs{
				Instructions: pulumi.String(judge.Instructions),
				ModelConfig: &bedrock.AgentcoreEvaluatorEvaluatorConfigLlmAsAJudgeModelConfigArgs{
					BedrockEvaluatorModelConfig: model,
				},
				RatingScale: scale,
			}
		}

		if e.CodeBased != nil {
			lambdaConfig := &bedrock.AgentcoreEvaluatorEvaluatorConfigCodeBasedLambdaConfigArgs{
				LambdaArn: pulumi.String(e.CodeBased.LambdaArn.GetValue()),
			}
			if e.CodeBased.TimeoutSeconds > 0 {
				lambdaConfig.LambdaTimeoutInSeconds = pulumi.Int(int(e.CodeBased.TimeoutSeconds))
			}
			config.CodeBased = &bedrock.AgentcoreEvaluatorEvaluatorConfigCodeBasedArgs{
				LambdaConfig: lambdaConfig,
			}
		}

		args := &bedrock.AgentcoreEvaluatorArgs{
			EvaluatorName:   pulumi.String(e.Name),
			Level:           pulumi.String(e.Level),
			EvaluatorConfig: config,
			Tags:            pulumi.ToStringMap(locals.AwsTags),
		}
		if e.Description != "" {
			args.Description = pulumi.String(e.Description)
		}
		if e.KmsKeyArn.GetValue() != "" {
			args.KmsKeyArn = pulumi.String(e.KmsKeyArn.GetValue())
		}

		resource, err := bedrock.NewAgentcoreEvaluator(ctx, "evaluator-"+e.Name, args, pulumi.Provider(provider))
		if err != nil {
			return nil, errors.Wrapf(err, "create evaluator %q", e.Name)
		}
		created[e.Name] = &createdEvaluator{resource: resource}
		evaluatorIds[e.Name] = resource.EvaluatorId
		evaluatorArns[e.Name] = resource.EvaluatorArn
	}

	ctx.Export(OpEvaluatorIds, evaluatorIds)
	ctx.Export(OpEvaluatorArns, evaluatorArns)
	return created, nil
}
