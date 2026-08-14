package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/opensearch"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// collection creates the OpenSearch Serverless collection and exports
// outputs.
func collection(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	createdEncryptionPolicy *opensearch.ServerlessSecurityPolicy,
) error {
	spec := locals.Spec

	args := &opensearch.ServerlessCollectionArgs{
		// The AWS collection name is create-time immutable and doubles as
		// the Pulumi resource name -- metadata.name on both engines.
		Name: pulumi.String(locals.CollectionName),
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Workload type (ForceNew). The default (TIMESERIES) is materialized by
	// the manifest loader, so the value always arrives -- sent explicitly
	// on both engines rather than leaning on AWS's own default.
	if spec.Type != nil && *spec.Type != "" {
		args.Type = pulumi.String(*spec.Type)
	}

	// Standby replicas (ForceNew). Default ENABLED materialized -- explicit
	// send keeps the engines symmetric and the manifest's intent visible in
	// state.
	if spec.StandbyReplicas != nil && *spec.StandbyReplicas != "" {
		args.StandbyReplicas = pulumi.String(*spec.StandbyReplicas)
	}

	// Collection-group membership (ForceNew). The group's standby-replicas
	// setting must match the collection's -- AWS rejects the mismatch at
	// create.
	if spec.CollectionGroupName != "" {
		args.CollectionGroupName = pulumi.String(spec.CollectionGroupName)
	}

	// GPU-accelerated vector capacity -- VECTORSEARCH collections only
	// (CEL enforces the coupling at manifest time).
	if spec.ServerlessVectorAcceleration != "" {
		args.VectorOptions = opensearch.ServerlessCollectionVectorOptionArray{
			&opensearch.ServerlessCollectionVectorOptionArgs{
				ServerlessVectorAcceleration: pulumi.String(spec.ServerlessVectorAcceleration),
			},
		}
	}

	// The encryption KEY CHOICE lives on the module-rendered encryption
	// security policy (see policies.go), which must match the collection
	// by name BEFORE it is created -- the dependency below enforces the
	// ordering on create and reverses it on destroy. The collection's own
	// inline encryption_config argument is the same setting's
	// collection-group-era twin and is deliberately not sent (recorded in
	// the parity manifest).
	c, err := opensearch.NewServerlessCollection(ctx, locals.CollectionName, args,
		pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdEncryptionPolicy}))
	if err != nil {
		return errors.Wrap(err, "create serverless collection")
	}

	ctx.Export(OpCollectionId, c.ID())
	ctx.Export(OpCollectionArn, c.Arn)
	ctx.Export(OpCollectionName, c.Name)
	ctx.Export(OpCollectionEndpoint, c.CollectionEndpoint)
	ctx.Export(OpDashboardEndpoint, c.DashboardEndpoint)
	// The effective key -- the customer-managed key when spec.encryption
	// chose one, else the AWS-owned key AWS reports back.
	ctx.Export(OpKmsKeyArn, c.KmsKeyArn)

	return nil
}
