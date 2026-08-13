package module

import (
	"fmt"

	"github.com/pkg/errors"
	gcpsecretmanagersecretv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpsecretmanagersecret/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/secretmanager"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// secret provisions the Secret Manager secret plus its optional first
// version and secret-scoped IAM grants.
//
// One kind, two GCP API surfaces: an empty spec.region creates the GLOBAL
// trio (Secret / SecretVersion / SecretIamMember) and a set region creates
// the REGIONAL trio (RegionalSecret / RegionalSecretVersion /
// RegionalSecretIamMember) — mirroring the Terraform module's count guards.
// The surfaces differ exactly where the spec does: replication is
// global-only (a regional secret's payloads live in its one region), and
// the regional secret takes CMEK directly via customer_managed_encryption.
//
// An OMITTED replication message renders the API's `auto {}` mode — the
// provider REQUIRES a replication block on the global secret, and automatic
// placement is the right default when no residency regime applies (the
// spec comment documents this contract).
//
// initial_version.enabled is sent EXPLICITLY: it is Optional in the
// provider with default true, and a spec transition true -> false must
// reach the API rather than being omitted (the send-true-or-omit class).
func secret(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpSecretManagerSecret.Spec

	// Enable the Secret Manager API so a fresh project can host the
	// secret. disable_on_destroy stays false (the provider default):
	// tearing down one secret must never disable the API for everything
	// else in the project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("secretmanager.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"secret-secretmanager.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable secretmanager.googleapis.com api")
	}

	if spec.Region == "" {
		return globalSecret(ctx, locals, gcpProvider, createdProjectService)
	}
	return regionalSecret(ctx, locals, gcpProvider, createdProjectService)
}

func globalSecret(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, apiService pulumi.Resource) error {
	spec := locals.GcpSecretManagerSecret.Spec

	args := &secretmanager.SecretArgs{
		SecretId:    pulumi.String(locals.SecretId),
		Replication: expandReplication(spec.Replication),
		Labels:      pulumi.ToStringMap(locals.GcpLabels),
	}

	if len(spec.Annotations) > 0 {
		args.Annotations = pulumi.ToStringMap(spec.Annotations)
	}
	if len(spec.Tags) > 0 {
		args.Tags = pulumi.ToStringMap(spec.Tags)
	}
	if spec.ExpireTime != "" {
		args.ExpireTime = pulumi.StringPtr(spec.ExpireTime)
	}
	if spec.Ttl != "" {
		args.Ttl = pulumi.StringPtr(spec.Ttl)
	}
	// GCP validates aliases against EXISTING versions at create/update, so an
	// alias cannot land in the same apply that seeds its version — add aliases
	// on a subsequent apply (live API: "Aliases cannot be assigned to versions
	// that don't exist").
	if len(spec.VersionAliases) > 0 {
		args.VersionAliases = pulumi.ToStringMap(spec.VersionAliases)
	}
	if spec.VersionDestroyTtl != "" {
		args.VersionDestroyTtl = pulumi.StringPtr(spec.VersionDestroyTtl)
	}
	if spec.Rotation != nil {
		rotationArgs := &secretmanager.SecretRotationArgs{}
		if spec.Rotation.RotationPeriod != "" {
			rotationArgs.RotationPeriod = pulumi.StringPtr(spec.Rotation.RotationPeriod)
		}
		if spec.Rotation.NextRotationTime != "" {
			rotationArgs.NextRotationTime = pulumi.StringPtr(spec.Rotation.NextRotationTime)
		}
		args.Rotation = rotationArgs
	}
	if len(spec.Topics) > 0 {
		topics := secretmanager.SecretTopicArray{}
		for _, topic := range spec.Topics {
			topics = append(topics, &secretmanager.SecretTopicArgs{
				Name: pulumi.String(topic.GetValue()),
			})
		}
		args.Topics = topics
	}
	if spec.DeletionProtection {
		args.DeletionProtection = pulumi.BoolPtr(true)
	}
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdSecret, err := secretmanager.NewSecret(ctx, "secret", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{apiService}))
	if err != nil {
		return errors.Wrap(err, "failed to create secret")
	}

	ctx.Export(OpSecretName, createdSecret.Name)
	ctx.Export(OpSecretId, createdSecret.SecretId)

	if spec.InitialVersion != nil {
		versionArgs := &secretmanager.SecretVersionArgs{
			Secret:     createdSecret.ID(),
			SecretData: pulumi.String(spec.InitialVersion.Data.GetValue()),
			// Explicit send — see the function comment.
			Enabled: pulumi.Bool(spec.InitialVersion.Enabled == nil || spec.InitialVersion.GetEnabled()),
		}
		if spec.InitialVersion.IsBase64 {
			versionArgs.IsSecretDataBase64 = pulumi.BoolPtr(true)
		}
		if spec.InitialVersion.DeletionPolicy != "" {
			versionArgs.DeletionPolicy = pulumi.String(spec.InitialVersion.DeletionPolicy)
		}
		createdVersion, err := secretmanager.NewSecretVersion(ctx, "secret-version", versionArgs,
			pulumi.Provider(gcpProvider))
		if err != nil {
			return errors.Wrap(err, "failed to create secret version")
		}
		ctx.Export(OpLatestVersionName, createdVersion.Name)
	} else {
		// The output key must exist either way so the outputs transformer
		// maps a stable shape; empty means "no initial version configured".
		ctx.Export(OpLatestVersionName, pulumi.String(""))
	}

	for index, member := range spec.IamMembers {
		memberArgs := &secretmanager.SecretIamMemberArgs{
			SecretId: createdSecret.SecretId,
			Role:     pulumi.String(member.Role),
			Member:   pulumi.String(member.Member.GetValue()),
		}
		if member.Condition != nil {
			conditionArgs := &secretmanager.SecretIamMemberConditionArgs{
				Title:      pulumi.String(member.Condition.Title),
				Expression: pulumi.String(member.Condition.Expression),
			}
			if member.Condition.Description != "" {
				conditionArgs.Description = pulumi.StringPtr(member.Condition.Description)
			}
			memberArgs.Condition = conditionArgs
		}
		if spec.ProjectId.GetValue() != "" {
			memberArgs.Project = pulumi.String(spec.ProjectId.GetValue())
		}
		if _, err := secretmanager.NewSecretIamMember(ctx,
			fmt.Sprintf("iam-member-%d", index), memberArgs, pulumi.Provider(gcpProvider)); err != nil {
			return errors.Wrapf(err, "failed to create iam member %d", index)
		}
	}

	return nil
}

func regionalSecret(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider, apiService pulumi.Resource) error {
	spec := locals.GcpSecretManagerSecret.Spec

	args := &secretmanager.RegionalSecretArgs{
		SecretId: pulumi.String(locals.SecretId),
		Location: pulumi.String(spec.Region),
		Labels:   pulumi.ToStringMap(locals.GcpLabels),
	}

	if spec.CustomerManagedEncryption != nil {
		args.CustomerManagedEncryption = &secretmanager.RegionalSecretCustomerManagedEncryptionArgs{
			KmsKeyName: pulumi.String(spec.CustomerManagedEncryption.KmsKey.GetValue()),
		}
	}
	if len(spec.Annotations) > 0 {
		args.Annotations = pulumi.ToStringMap(spec.Annotations)
	}
	if len(spec.Tags) > 0 {
		args.Tags = pulumi.ToStringMap(spec.Tags)
	}
	if spec.ExpireTime != "" {
		args.ExpireTime = pulumi.StringPtr(spec.ExpireTime)
	}
	if spec.Ttl != "" {
		args.Ttl = pulumi.StringPtr(spec.Ttl)
	}
	// Same alias temporal constraint as the global variant above.
	if len(spec.VersionAliases) > 0 {
		args.VersionAliases = pulumi.ToStringMap(spec.VersionAliases)
	}
	if spec.VersionDestroyTtl != "" {
		args.VersionDestroyTtl = pulumi.StringPtr(spec.VersionDestroyTtl)
	}
	if spec.Rotation != nil {
		rotationArgs := &secretmanager.RegionalSecretRotationArgs{}
		if spec.Rotation.RotationPeriod != "" {
			rotationArgs.RotationPeriod = pulumi.StringPtr(spec.Rotation.RotationPeriod)
		}
		if spec.Rotation.NextRotationTime != "" {
			rotationArgs.NextRotationTime = pulumi.StringPtr(spec.Rotation.NextRotationTime)
		}
		args.Rotation = rotationArgs
	}
	if len(spec.Topics) > 0 {
		topics := secretmanager.RegionalSecretTopicArray{}
		for _, topic := range spec.Topics {
			topics = append(topics, &secretmanager.RegionalSecretTopicArgs{
				Name: pulumi.String(topic.GetValue()),
			})
		}
		args.Topics = topics
	}
	if spec.DeletionProtection {
		args.DeletionProtection = pulumi.BoolPtr(true)
	}
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	createdSecret, err := secretmanager.NewRegionalSecret(ctx, "secret", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{apiService}))
	if err != nil {
		return errors.Wrap(err, "failed to create regional secret")
	}

	ctx.Export(OpSecretName, createdSecret.Name)
	ctx.Export(OpSecretId, createdSecret.SecretId)

	if spec.InitialVersion != nil {
		versionArgs := &secretmanager.RegionalSecretVersionArgs{
			Secret:     createdSecret.ID(),
			SecretData: pulumi.String(spec.InitialVersion.Data.GetValue()),
			// Explicit send — see the function comment.
			Enabled: pulumi.Bool(spec.InitialVersion.Enabled == nil || spec.InitialVersion.GetEnabled()),
		}
		if spec.InitialVersion.IsBase64 {
			versionArgs.IsSecretDataBase64 = pulumi.BoolPtr(true)
		}
		if spec.InitialVersion.DeletionPolicy != "" {
			versionArgs.DeletionPolicy = pulumi.String(spec.InitialVersion.DeletionPolicy)
		}
		createdVersion, err := secretmanager.NewRegionalSecretVersion(ctx, "secret-version", versionArgs,
			pulumi.Provider(gcpProvider))
		if err != nil {
			return errors.Wrap(err, "failed to create regional secret version")
		}
		ctx.Export(OpLatestVersionName, createdVersion.Name)
	} else {
		ctx.Export(OpLatestVersionName, pulumi.String(""))
	}

	for index, member := range spec.IamMembers {
		memberArgs := &secretmanager.RegionalSecretIamMemberArgs{
			SecretId: createdSecret.SecretId,
			Location: pulumi.StringPtr(spec.Region),
			Role:     pulumi.String(member.Role),
			Member:   pulumi.String(member.Member.GetValue()),
		}
		if member.Condition != nil {
			conditionArgs := &secretmanager.RegionalSecretIamMemberConditionArgs{
				Title:      pulumi.String(member.Condition.Title),
				Expression: pulumi.String(member.Condition.Expression),
			}
			if member.Condition.Description != "" {
				conditionArgs.Description = pulumi.StringPtr(member.Condition.Description)
			}
			memberArgs.Condition = conditionArgs
		}
		if spec.ProjectId.GetValue() != "" {
			memberArgs.Project = pulumi.String(spec.ProjectId.GetValue())
		}
		if _, err := secretmanager.NewRegionalSecretIamMember(ctx,
			fmt.Sprintf("iam-member-%d", index), memberArgs, pulumi.Provider(gcpProvider)); err != nil {
			return errors.Wrapf(err, "failed to create regional iam member %d", index)
		}
	}

	return nil
}

// expandReplication renders the GLOBAL secret's replication block. The
// provider requires the block, so an omitted spec message becomes `auto {}`
// — automatic placement, the documented default.
func expandReplication(replication *gcpsecretmanagersecretv1alpha1.GcpSecretManagerSecretReplication) *secretmanager.SecretReplicationArgs {
	if replication == nil {
		return &secretmanager.SecretReplicationArgs{
			Auto: &secretmanager.SecretReplicationAutoArgs{},
		}
	}

	if replication.Auto != nil {
		autoArgs := &secretmanager.SecretReplicationAutoArgs{}
		if replication.Auto.CustomerManagedEncryption != nil {
			autoArgs.CustomerManagedEncryption = &secretmanager.SecretReplicationAutoCustomerManagedEncryptionArgs{
				KmsKeyName: pulumi.String(replication.Auto.CustomerManagedEncryption.KmsKey.GetValue()),
			}
		}
		return &secretmanager.SecretReplicationArgs{Auto: autoArgs}
	}

	replicas := secretmanager.SecretReplicationUserManagedReplicaArray{}
	for _, replica := range replication.UserManaged.Replicas {
		replicaArgs := &secretmanager.SecretReplicationUserManagedReplicaArgs{
			Location: pulumi.String(replica.Location),
		}
		if replica.CustomerManagedEncryption != nil {
			replicaArgs.CustomerManagedEncryption = &secretmanager.SecretReplicationUserManagedReplicaCustomerManagedEncryptionArgs{
				KmsKeyName: pulumi.String(replica.CustomerManagedEncryption.KmsKey.GetValue()),
			}
		}
		replicas = append(replicas, replicaArgs)
	}
	return &secretmanager.SecretReplicationArgs{
		UserManaged: &secretmanager.SecretReplicationUserManagedArgs{Replicas: replicas},
	}
}
