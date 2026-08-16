package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cfg"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// conformancePack creates the pack at the spec's scope - the
// account-scope resource or the organization-scope one, never both -
// and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - deploying a pack REQUIRES an active Config recorder in the
//     region (AWS rejects it otherwise) - the recorder is the
//     consumer's contract (AwsConfigRecorder), never this module's;
//   - AWS never reports the template back, so template drift is
//     undetectable by design (both provider resources document it);
//   - the pack service-linked role creates the pack's rules; deleting
//     the pack deletes them (org packs unwind from member accounts,
//     which can take minutes - the provider waits);
//   - conformance packs carry NO tags at the provider.
func conformancePack(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	if spec.OrganizationScope {
		// Organization packs cap the name at 128 chars; the delivery
		// bucket must begin with "awsconfigconforms" (AWS naming
		// contracts enforced at deploy). Exactly one template form
		// arrives here (the spec CEL guarantees it).
		args := &cfg.OrganizationConformancePackArgs{
			Name: pulumi.String(locals.Target.Metadata.Name),
		}
		if spec.DeliveryS3Bucket.GetValue() != "" {
			args.DeliveryS3Bucket = pulumi.String(spec.DeliveryS3Bucket.GetValue())
		}
		if spec.DeliveryS3KeyPrefix != "" {
			args.DeliveryS3KeyPrefix = pulumi.String(spec.DeliveryS3KeyPrefix)
		}
		if spec.TemplateBody != "" {
			args.TemplateBody = pulumi.String(spec.TemplateBody)
		}
		if spec.TemplateS3Uri != "" {
			args.TemplateS3Uri = pulumi.String(spec.TemplateS3Uri)
		}
		if len(spec.ExcludedAccounts) > 0 {
			args.ExcludedAccounts = pulumi.ToStringArray(spec.ExcludedAccounts)
		}
		var params cfg.OrganizationConformancePackInputParameterArray
		for _, p := range spec.InputParameters {
			params = append(params, &cfg.OrganizationConformancePackInputParameterArgs{
				ParameterName:  pulumi.String(p.ParameterName),
				ParameterValue: pulumi.String(p.ParameterValue),
			})
		}
		if len(params) > 0 {
			args.InputParameters = params
		}

		createdPack, err := cfg.NewOrganizationConformancePack(ctx, "organization-conformance-pack", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create organization conformance pack")
		}

		ctx.Export(OpPackName, createdPack.Name)
		ctx.Export(OpPackArn, createdPack.Arn)
		return nil
	}

	// Account packs accept both template forms at once (AWS prefers
	// the S3 one); the spec CEL guarantees at least one arrives here.
	args := &cfg.ConformancePackArgs{
		Name: pulumi.String(locals.Target.Metadata.Name),
	}
	if spec.DeliveryS3Bucket.GetValue() != "" {
		args.DeliveryS3Bucket = pulumi.String(spec.DeliveryS3Bucket.GetValue())
	}
	if spec.DeliveryS3KeyPrefix != "" {
		args.DeliveryS3KeyPrefix = pulumi.String(spec.DeliveryS3KeyPrefix)
	}
	if spec.TemplateBody != "" {
		args.TemplateBody = pulumi.String(spec.TemplateBody)
	}
	if spec.TemplateS3Uri != "" {
		args.TemplateS3Uri = pulumi.String(spec.TemplateS3Uri)
	}
	var params cfg.ConformancePackInputParameterArray
	for _, p := range spec.InputParameters {
		params = append(params, &cfg.ConformancePackInputParameterArgs{
			ParameterName:  pulumi.String(p.ParameterName),
			ParameterValue: pulumi.String(p.ParameterValue),
		})
	}
	if len(params) > 0 {
		args.InputParameters = params
	}

	createdPack, err := cfg.NewConformancePack(ctx, "conformance-pack", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create conformance pack")
	}

	ctx.Export(OpPackName, createdPack.Name)
	ctx.Export(OpPackArn, createdPack.Arn)
	return nil
}
