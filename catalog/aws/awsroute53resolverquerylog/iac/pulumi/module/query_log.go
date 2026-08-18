package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// queryLogConfig creates the logging configuration and its VPC
// associations, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - name and destination_arn are both ForceNew - the configuration
//     is immutable except tags (existing log data survives
//     replacement in the destination);
//   - associations are pure joins (config, vpc) - every argument
//     ForceNew, no update path;
//   - an association can FAIL asynchronously after a clean apply when
//     the resolver cannot write to the destination - the provider's
//     waiter surfaces the association's error code.
func queryLogConfig(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	createdConfig, err := route53.NewResolverQueryLogConfig(ctx, "config",
		&route53.ResolverQueryLogConfigArgs{
			Name:           pulumi.String(locals.Target.Metadata.Name),
			DestinationArn: pulumi.String(spec.DestinationArn.GetValue()),
			Tags:           pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create config")
	}

	associationIds := pulumi.StringMap{}
	for _, vpcId := range spec.VpcIds {
		resolvedVpcId := vpcId.GetValue()
		createdAssociation, err := route53.NewResolverQueryLogConfigAssociation(ctx,
			fmt.Sprintf("association-%s", resolvedVpcId),
			&route53.ResolverQueryLogConfigAssociationArgs{
				ResolverQueryLogConfigId: createdConfig.ID(),
				ResourceId:               pulumi.String(resolvedVpcId),
			}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "associate vpc %s", resolvedVpcId)
		}
		associationIds[resolvedVpcId] = createdAssociation.ID().ToStringOutput()
	}

	ctx.Export(OpQueryLogConfigId, createdConfig.ID())
	ctx.Export(OpQueryLogConfigArn, createdConfig.Arn)
	ctx.Export(OpShareStatus, createdConfig.ShareStatus)
	ctx.Export(OpAssociationIds, associationIds)
	return nil
}
