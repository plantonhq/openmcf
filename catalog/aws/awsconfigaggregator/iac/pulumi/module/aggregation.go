package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/cfg"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// aggregation creates the aggregation arms - the aggregator
// (collector side) and/or the reciprocal grants (source-account
// side) - and exports outputs.
//
// Lifecycle facts the renders below depend on:
//   - the provider replaces the aggregator only when a source block
//     APPEARS on an existing aggregator (absent -> present); content
//     changes and block removal update in place;
//   - the spec CEL guarantees exactly one source shape arrives here;
//   - grants are keyed "{account_id}:{authorized_aws_region}" (the
//     provider's own import ID), so reordering the spec list never
//     churns them;
//   - the provider's deprecated "region" alias on the grant is
//     deliberately not rendered - authorized_aws_region is the
//     surviving argument.
func aggregation(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	if spec.Aggregation != nil {
		args := &cfg.ConfigurationAggregatorArgs{
			// metadata.name is the aggregator name on both engines.
			Name: pulumi.String(locals.Target.Metadata.Name),
			Tags: pulumi.ToStringMap(locals.AwsTags),
		}

		if src := spec.Aggregation.AccountSource; src != nil {
			srcArgs := &cfg.ConfigurationAggregatorAccountAggregationSourceArgs{
				AccountIds: pulumi.ToStringArray(src.AccountIds),
			}
			if src.AllRegions {
				srcArgs.AllRegions = pulumi.Bool(true)
			}
			if len(src.Regions) > 0 {
				srcArgs.Regions = pulumi.ToStringArray(src.Regions)
			}
			args.AccountAggregationSource = srcArgs
		}

		if src := spec.Aggregation.OrganizationSource; src != nil {
			srcArgs := &cfg.ConfigurationAggregatorOrganizationAggregationSourceArgs{
				RoleArn: pulumi.String(src.RoleArn.GetValue()),
			}
			if src.AllRegions {
				srcArgs.AllRegions = pulumi.Bool(true)
			}
			if len(src.Regions) > 0 {
				srcArgs.Regions = pulumi.ToStringArray(src.Regions)
			}
			args.OrganizationAggregationSource = srcArgs
		}

		createdAggregator, err := cfg.NewConfigurationAggregator(ctx, "aggregator", args, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create configuration aggregator")
		}

		ctx.Export(OpAggregatorName, createdAggregator.Name)
		ctx.Export(OpAggregatorArn, createdAggregator.Arn)
	} else {
		// A grants-only deployment still exports the aggregator arms
		// empty so consumers read a stable output shape.
		ctx.Export(OpAggregatorName, pulumi.String(""))
		ctx.Export(OpAggregatorArn, pulumi.String(""))
	}

	// The source-account side: each grant authorizes ONE aggregator
	// (account+region) to collect this account's Config data.
	authorizationArns := pulumi.StringMap{}
	for _, grant := range spec.Authorizations {
		key := fmt.Sprintf("%s:%s", grant.AccountId, grant.AuthorizedAwsRegion)
		createdGrant, err := cfg.NewAggregateAuthorization(ctx, "grant-"+key, &cfg.AggregateAuthorizationArgs{
			AccountId:           pulumi.String(grant.AccountId),
			AuthorizedAwsRegion: pulumi.String(grant.AuthorizedAwsRegion),
			Tags:                pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create aggregate authorization %s", key)
		}
		authorizationArns[key] = createdGrant.Arn
	}
	ctx.Export(OpAuthorizationArns, authorizationArns)

	return nil
}
