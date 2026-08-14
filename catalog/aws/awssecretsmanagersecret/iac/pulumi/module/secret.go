package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/secretsmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// secret creates the Secrets Manager secret plus its satellites (resource
// policy, managed version, rotation configuration) and exports outputs.
func secret(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &secretsmanager.SecretArgs{
		// The AWS secret name is create-time immutable and doubles as the
		// Pulumi resource name -- metadata.name on both engines (never
		// provider auto-naming, which would suffix a random token and
		// diverge from Terraform).
		Name: pulumi.String(locals.SecretName),
		// Description is ALWAYS sent explicitly so the two engines never
		// inject differing defaults into state.
		Description: pulumi.String(spec.Description),
		// Always sent: false fails replication loudly on a name collision
		// in a replica region (the safe posture); true overwrites. Explicit
		// on both engines so the send surface stays symmetric.
		ForceOverwriteReplicaSecret: pulumi.Bool(spec.ForceOverwriteReplicaSecret),
		Tags:                        pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.String(spec.KmsKeyId.GetValue())
	}

	// Consumed only at delete time: 0 forces immediate permanent deletion,
	// 7-30 keeps the soft-delete recovery window. The default (30) is
	// materialized by the manifest loader, so the value always arrives.
	if spec.RecoveryWindowInDays != nil {
		args.RecoveryWindowInDays = pulumi.Int(int(*spec.RecoveryWindowInDays))
	}

	// Cross-region replicas. Each replica encrypts under its own region's
	// key (the referenced customer key, or that region's AWS-managed key).
	// Two delete-time truths both engines inherit from the provider
	// (live-verified 2026-08-13): AWS deletes replica secrets
	// ASYNCHRONOUSLY after RemoveRegionsFromReplication and the provider
	// does not wait for it, so destroying with recovery_window 0 can
	// strand a replica as a live standalone secret (a recovery window
	// lets the async deletion complete); and replication is never waited
	// on at create either, so a Failed replication (e.g. against a
	// stranded same-name ex-replica, which force_overwrite does NOT
	// clear) is silent at apply.
	if len(spec.ReplicaRegions) > 0 {
		var replicas secretsmanager.SecretReplicaArray
		for _, r := range spec.ReplicaRegions {
			replica := &secretsmanager.SecretReplicaArgs{
				Region: pulumi.String(r.Region),
			}
			if r.KmsKeyId.GetValue() != "" {
				replica.KmsKeyId = pulumi.String(r.KmsKeyId.GetValue())
			}
			replicas = append(replicas, replica)
		}
		args.Replicas = replicas
	}

	// Managed external secret partner identifier (ForceNew). Empty means an
	// ordinary self-managed secret; the argument is omitted entirely then --
	// AWS treats an absent Type and an empty Type differently in error
	// paths, and omission is the documented shape for self-managed secrets.
	if spec.Type != "" {
		args.Type = pulumi.String(spec.Type)
	}

	createdSecret, err := secretsmanager.NewSecret(ctx, locals.SecretName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create secret")
	}

	ctx.Export(OpSecretArn, createdSecret.Arn)
	ctx.Export(OpSecretName, createdSecret.Name)

	// Resource policy -- rendered through the standalone policy resource
	// (not the secret's inline policy argument) because only the standalone
	// resource carries block_public_policy, the PutResourcePolicy guard
	// that rejects policies granting anonymous access.
	if spec.Policy != nil {
		policyJSON, err := json.Marshal(spec.Policy.AsMap())
		if err != nil {
			return errors.Wrap(err, "marshal resource policy to JSON")
		}
		policyArgs := &secretsmanager.SecretPolicyArgs{
			SecretArn: createdSecret.Arn,
			Policy:    pulumi.String(string(policyJSON)),
		}
		// Default true (materialized by the manifest loader): reject
		// public policies unless the manifest deliberately opts out.
		if spec.BlockPublicPolicy != nil {
			policyArgs.BlockPublicPolicy = pulumi.Bool(*spec.BlockPublicPolicy)
		}
		if _, err := secretsmanager.NewSecretPolicy(ctx, "policy", policyArgs, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "create secret policy")
		}
	}

	// The managed version -- created only when a value arm is set (a shell
	// secret with no value is legal; an application or rotation function
	// writes the first version then). version_id is exported in every arm
	// (empty for a shell secret) so both engines emit the same output set.
	var createdVersion *secretsmanager.SecretVersion
	if spec.StringValue != "" || spec.BinaryValue != "" {
		versionArgs := &secretsmanager.SecretVersionArgs{
			SecretId: createdSecret.Arn,
		}
		if spec.StringValue != "" {
			versionArgs.SecretString = pulumi.String(spec.StringValue)
		}
		if spec.BinaryValue != "" {
			// The provider expects base64 in secret_binary and decodes it
			// before calling PutSecretValue (CEL already guaranteed the
			// encoding at manifest time).
			versionArgs.SecretBinary = pulumi.String(spec.BinaryValue)
		}
		// Custom staging labels ride ALONGSIDE AWSCURRENT: when the
		// manifest adds labels, the module sends AWSCURRENT explicitly too
		// -- providing version_stages REPLACES the automatic AWSCURRENT
		// assignment, and dropping it would leave the secret with no
		// current version.
		if len(spec.VersionStages) > 0 {
			stages := pulumi.StringArray{pulumi.String("AWSCURRENT")}
			for _, s := range spec.VersionStages {
				stages = append(stages, pulumi.String(s))
			}
			versionArgs.VersionStages = stages
		}
		createdVersion, err = secretsmanager.NewSecretVersion(ctx, "value", versionArgs, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "create secret version")
		}
		ctx.Export(OpVersionId, createdVersion.VersionId)
	} else {
		ctx.Export(OpVersionId, pulumi.String(""))
	}

	// Rotation. Ordered after the version: with rotate_immediately (the
	// default) AWS invokes the rotation mechanism as soon as rotation is
	// configured, and the rotation function reads the current value -- so
	// the value must exist first.
	if spec.Rotation != nil {
		rotation := spec.Rotation
		rotationArgs := &secretsmanager.SecretRotationArgs{
			SecretId: createdSecret.Arn,
		}
		if rotation.RotationLambdaArn.GetValue() != "" {
			rotationArgs.RotationLambdaArn = pulumi.String(rotation.RotationLambdaArn.GetValue())
		}
		if rotation.ExternalRotationRoleArn.GetValue() != "" {
			rotationArgs.ExternalSecretRotationRoleArn = pulumi.String(rotation.ExternalRotationRoleArn.GetValue())
		}
		if len(rotation.ExternalRotationMetadata) > 0 {
			var items secretsmanager.SecretRotationExternalSecretRotationMetadataArray
			for _, m := range rotation.ExternalRotationMetadata {
				items = append(items, &secretsmanager.SecretRotationExternalSecretRotationMetadataArgs{
					Key:   pulumi.String(m.Key),
					Value: pulumi.String(m.Value),
				})
			}
			rotationArgs.ExternalSecretRotationMetadatas = items
		}
		// Default true (materialized): rotate once as soon as rotation is
		// configured. Explicit false only tests the configuration.
		if rotation.RotateImmediately != nil {
			rotationArgs.RotateImmediately = pulumi.Bool(*rotation.RotateImmediately)
		}
		rules := &secretsmanager.SecretRotationRotationRulesArgs{}
		if rotation.AutomaticallyAfterDays != nil {
			rules.AutomaticallyAfterDays = pulumi.Int(int(*rotation.AutomaticallyAfterDays))
		}
		if rotation.ScheduleExpression != "" {
			rules.ScheduleExpression = pulumi.String(rotation.ScheduleExpression)
		}
		if rotation.Duration != "" {
			rules.Duration = pulumi.String(rotation.Duration)
		}
		rotationArgs.RotationRules = rules

		var rotationDeps []pulumi.Resource
		if createdVersion != nil {
			rotationDeps = append(rotationDeps, createdVersion)
		}
		if _, err := secretsmanager.NewSecretRotation(ctx, "rotation", rotationArgs,
			pulumi.Provider(provider), pulumi.DependsOn(rotationDeps)); err != nil {
			return errors.Wrap(err, "create secret rotation")
		}
	}

	return nil
}
