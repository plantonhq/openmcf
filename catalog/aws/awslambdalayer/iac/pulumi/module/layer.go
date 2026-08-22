package module

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// layer publishes the layer version from its S3 archive and attaches
// the share grants, then exports outputs.
//
// Lifecycle facts the render below depends on:
//   - every layer-version argument is ForceNew - a config change
//     publishes a NEW version (functions keep the exact version ARN
//     they pinned, so a replacement never breaks consumers mid-run);
//   - skip_destroy leaves the previous version available in AWS on
//     replacement/destroy (dormant versions bill nothing);
//   - permissions are per-VERSION policy statements - they replace
//     alongside the layer version they grant access to;
//   - source_code_hash is a local change detector only - AWS never
//     reports it back.
func layer(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	layerArgs := &lambda.LayerVersionArgs{
		LayerName: pulumi.String(locals.Target.Metadata.Name),
		S3Bucket:  pulumi.String(spec.Code.Bucket.GetValue()),
		S3Key:     pulumi.String(spec.Code.Key),
	}
	if spec.Code.Version != "" {
		layerArgs.S3ObjectVersion = pulumi.String(spec.Code.Version)
	}
	if spec.SourceCodeHash != "" {
		layerArgs.SourceCodeHash = pulumi.String(spec.SourceCodeHash)
	}
	if spec.Description != "" {
		layerArgs.Description = pulumi.String(spec.Description)
	}
	if len(spec.CompatibleRuntimes) > 0 {
		layerArgs.CompatibleRuntimes = pulumi.ToStringArray(spec.CompatibleRuntimes)
	}
	if len(spec.CompatibleArchitectures) > 0 {
		layerArgs.CompatibleArchitectures = pulumi.ToStringArray(spec.CompatibleArchitectures)
	}
	if spec.LicenseInfo != "" {
		layerArgs.LicenseInfo = pulumi.String(spec.LicenseInfo)
	}
	if spec.SkipDestroy {
		layerArgs.SkipDestroy = pulumi.Bool(true)
	}

	createdLayer, err := lambda.NewLayerVersion(ctx, "layer", layerArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create layer version")
	}

	// The permission resource takes the version as a NUMBER while the
	// layer resource reports it as a string.
	versionNumber := createdLayer.Version.ApplyT(func(v string) (int, error) {
		return strconv.Atoi(v)
	}).(pulumi.IntOutput)

	permissionRevisionIds := pulumi.StringMap{}
	for _, permission := range spec.Permissions {
		permissionArgs := &lambda.LayerVersionPermissionArgs{
			// lambda:GetLayerVersion is the only action AWS supports
			// on layers - pinned here, never spec surface.
			Action:        pulumi.String("lambda:GetLayerVersion"),
			LayerName:     createdLayer.LayerArn,
			VersionNumber: versionNumber,
			StatementId:   pulumi.String(permission.StatementId),
			Principal:     pulumi.String(permission.Principal),
		}
		if permission.OrganizationId != "" {
			permissionArgs.OrganizationId = pulumi.String(permission.OrganizationId)
		}
		if permission.SkipDestroy {
			permissionArgs.SkipDestroy = pulumi.Bool(true)
		}
		createdPermission, err := lambda.NewLayerVersionPermission(ctx,
			fmt.Sprintf("permission-%s", permission.StatementId),
			permissionArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrapf(err, "create permission %s", permission.StatementId)
		}
		permissionRevisionIds[permission.StatementId] = createdPermission.RevisionId
	}

	ctx.Export(OpLayerArn, createdLayer.LayerArn)
	ctx.Export(OpLayerVersionArn, createdLayer.Arn)
	ctx.Export(OpVersion, createdLayer.Version)
	ctx.Export(OpCodeSha256, createdLayer.CodeSha256)
	ctx.Export(OpPermissionRevisionIds, permissionRevisionIds)
	return nil
}
