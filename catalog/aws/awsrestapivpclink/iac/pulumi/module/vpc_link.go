package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/apigateway"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// vpcLink creates the REST API VPC link and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the target NLB is immutable in AWS (the provider replaces the
//     link when it changes) -- expected: a link IS its network
//     attachment;
//   - AWS takes exactly one balancer per link (the provider's
//     target_arns list caps at one; this engine's SDK flattens it to
//     the singular target_arn) -- the spec is singular by design;
//   - creation waits for the attachment to reach AVAILABLE (up to ~20
//     minutes upstream) before integrations can reference the link.
func vpcLink(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &apigateway.VpcLinkArgs{
		// metadata.name is the naming basis on both engines.
		Name:      pulumi.String(locals.Target.Metadata.Name),
		TargetArn: pulumi.String(spec.TargetArn.GetValue()),
		Tags:      pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	created, err := apigateway.NewVpcLink(ctx, "vpc-link", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create vpc link")
	}

	ctx.Export(OpVpcLinkId, created.ID())
	ctx.Export(OpVpcLinkArn, created.Arn)
	return nil
}
