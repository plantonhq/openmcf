package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/backup"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// framework creates the Backup Audit Manager framework and exports
// outputs.
//
// Lifecycle facts the render below depends on:
//   - AWS framework names forbid hyphens, so the name is
//     spec.framework_name (an explicit field), never metadata.name;
//   - evaluations run on AWS Config: without an ACTIVE recorder
//     recording the backup types, deployment lands FAILED - and the
//     provider treats FAILED as a completed apply (the failure shows
//     in deployment_status, not as an error).
func framework(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	var controls backup.FrameworkControlArray
	for _, c := range spec.Controls {
		controlArgs := &backup.FrameworkControlArgs{
			Name: pulumi.String(c.Name),
		}
		var params backup.FrameworkControlInputParameterArray
		for _, p := range c.InputParameters {
			params = append(params, &backup.FrameworkControlInputParameterArgs{
				Name:  pulumi.String(p.Name),
				Value: pulumi.String(p.Value),
			})
		}
		if len(params) > 0 {
			controlArgs.InputParameters = params
		}
		if c.Scope != nil {
			scopeArgs := &backup.FrameworkControlScopeArgs{}
			if len(c.Scope.ComplianceResourceIds) > 0 {
				scopeArgs.ComplianceResourceIds = pulumi.ToStringArray(c.Scope.ComplianceResourceIds)
			}
			if len(c.Scope.ComplianceResourceTypes) > 0 {
				scopeArgs.ComplianceResourceTypes = pulumi.ToStringArray(c.Scope.ComplianceResourceTypes)
			}
			if len(c.Scope.Tags) > 0 {
				scopeArgs.Tags = pulumi.ToStringMap(c.Scope.Tags)
			}
			controlArgs.Scope = scopeArgs
		}
		controls = append(controls, controlArgs)
	}

	createdFramework, err := backup.NewFramework(ctx, "framework", &backup.FrameworkArgs{
		Name:        pulumi.String(spec.FrameworkName),
		Description: pulumi.String(spec.Description),
		Controls:    controls,
		Tags:        pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create framework")
	}

	ctx.Export(OpFrameworkArn, createdFramework.Arn)
	// Frameworks are addressed by region + name; consumers (and the
	// harness verifier) reaching the framework off the ambient region
	// need the resolved region alongside the ARN.
	ctx.Export(OpRegion, createdFramework.Region)
	return nil
}
