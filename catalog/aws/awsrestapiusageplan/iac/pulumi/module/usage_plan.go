package module

import (
	"sort"

	"github.com/pkg/errors"
	awsrestapiusageplanv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsrestapiusageplan/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// usagePlan creates the plan, its API keys, and the key memberships,
// and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - the plan's api_stages PATCH away cleanly on update, and upstream
//     detaches every stage before deleting the plan -- stage coverage is
//     freely editable;
//   - product_code cannot be set at create (AWS rejects it); the
//     provider applies it via a follow-up PATCH -- no ordering concern
//     for this module, but worth knowing when reading apply logs;
//   - a usage_plan_key is pure membership (key <-> plan); every field
//     change replaces it, which is free and instant;
//   - key VALUES are secrets: AWS generates them unless the spec pins
//     one, and this module never exports them.
func usagePlan(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &apigateway.UsagePlanArgs{
		// metadata.name is the naming basis on both engines.
		Name: pulumi.String(locals.Target.Metadata.Name),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.ProductCode != "" {
		args.ProductCode = pulumi.String(spec.ProductCode)
	}

	// Stage coverage with optional per-method throttles.
	var apiStages apigateway.UsagePlanApiStageArray
	for _, s := range spec.ApiStages {
		stage := &apigateway.UsagePlanApiStageArgs{
			ApiId: pulumi.String(s.RestApiId.GetValue()),
			Stage: pulumi.String(s.StageName.GetValue()),
		}
		var throttles apigateway.UsagePlanApiStageThrottleArray
		for _, t := range s.MethodThrottles {
			throttles = append(throttles, &apigateway.UsagePlanApiStageThrottleArgs{
				Path:       pulumi.String(t.Path),
				BurstLimit: pulumi.Int(int(t.BurstLimit)),
				RateLimit:  pulumi.Float64(t.RateLimit),
			})
		}
		if len(throttles) > 0 {
			stage.Throttles = throttles
		}
		apiStages = append(apiStages, stage)
	}
	if len(apiStages) > 0 {
		args.ApiStages = apiStages
	}

	if spec.Quota != nil {
		args.QuotaSettings = &apigateway.UsagePlanQuotaSettingsArgs{
			Limit:  pulumi.Int(int(spec.Quota.Limit)),
			Period: pulumi.String(spec.Quota.Period),
			Offset: pulumi.Int(int(spec.Quota.Offset)),
		}
	}
	if spec.Throttle != nil {
		args.ThrottleSettings = &apigateway.UsagePlanThrottleSettingsArgs{
			BurstLimit: pulumi.Int(int(spec.Throttle.BurstLimit)),
			RateLimit:  pulumi.Float64(spec.Throttle.RateLimit),
		}
	}

	plan, err := apigateway.NewUsagePlan(ctx, "usage-plan", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create usage plan")
	}

	// API keys and their plan memberships. Iteration is name-sorted for
	// deterministic previews.
	keyIds := pulumi.StringMap{}
	keyArns := pulumi.StringMap{}
	for _, k := range sortedApiKeys(spec.ApiKeys) {
		keyArgs := &apigateway.ApiKeyArgs{
			Name: pulumi.String(k.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}
		if k.Description != "" {
			keyArgs.Description = pulumi.String(k.Description)
		}
		// Rendered only on an explicit choice so the module never fights
		// AWS's enabled-by-default.
		if k.Enabled != nil {
			keyArgs.Enabled = pulumi.Bool(*k.Enabled)
		}
		if k.CustomerId != "" {
			keyArgs.CustomerId = pulumi.String(k.CustomerId)
		}
		// Omitted = AWS generates the value (recommended); a pinned value
		// arrives resolved from the managed-secret reference.
		if k.Value != "" {
			keyArgs.Value = pulumi.String(k.Value)
		}
		key, err := apigateway.NewApiKey(ctx, "api-key-"+k.Name, keyArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create api key %q", k.Name)
		}

		// The membership attaching the key to this plan.
		_, err = apigateway.NewUsagePlanKey(ctx, "usage-plan-key-"+k.Name, &apigateway.UsagePlanKeyArgs{
			KeyId:       key.ID(),
			KeyType:     pulumi.String("API_KEY"),
			UsagePlanId: plan.ID(),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "attach api key %q", k.Name)
		}

		keyIds[k.Name] = key.ID().ToStringOutput()
		keyArns[k.Name] = key.Arn
	}

	ctx.Export(OpUsagePlanId, plan.ID())
	ctx.Export(OpUsagePlanArn, plan.Arn)
	ctx.Export(OpApiKeyIds, keyIds)
	ctx.Export(OpApiKeyArns, keyArns)
	return nil
}

func sortedApiKeys(in []*awsrestapiusageplanv1alpha1.AwsRestApiUsagePlanApiKey) []*awsrestapiusageplanv1alpha1.AwsRestApiUsagePlanApiKey {
	out := append([]*awsrestapiusageplanv1alpha1.AwsRestApiUsagePlanApiKey{}, in...)
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}
