package module

import (
	"github.com/pkg/errors"
	awssagemakerdomainv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssagemakerdomain/v1alpha1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// domain creates the SageMaker Domain. The spec's presence semantics map
// one-to-one onto the SDK: absent optional scalars and messages are simply not
// sent, so AWS applies its own defaults and the preview stays clean -- never
// an always-sent zero value, which would pin defaults and create phantom
// diffs. Same contract as the Terraform module.
func domain(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*sagemaker.Domain, error) {
	spec := locals.Spec

	var subnetIds pulumi.StringArray
	for _, s := range spec.SubnetIds {
		subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
	}

	args := &sagemaker.DomainArgs{
		// The domain's cloud name is metadata.name (create-time-immutable);
		// set explicitly so both engines deploy the identical name instead of
		// relying on Pulumi auto-naming.
		DomainName:          pulumi.String(locals.DomainName),
		AuthMode:            pulumi.String(spec.AuthMode),
		VpcId:               pulumi.String(spec.VpcId.GetValue()),
		SubnetIds:           subnetIds,
		DefaultUserSettings: buildDefaultUserSettings(spec.DefaultUserSettings),
		Tags:                pulumi.ToStringMap(locals.AwsTags),
	}

	if spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.StringPtr(spec.KmsKeyId.GetValue())
	}

	if spec.AppNetworkAccessType != nil {
		args.AppNetworkAccessType = pulumi.StringPtr(spec.GetAppNetworkAccessType())
	}

	// Only honored by AWS when RStudio is configured (the spec enforces the
	// pairing, so a silent no-op cannot reach this argument).
	if spec.AppSecurityGroupManagement != nil {
		args.AppSecurityGroupManagement = pulumi.StringPtr(spec.GetAppSecurityGroupManagement())
	}

	// Whether the domain's tags propagate to apps/spaces/user profiles
	// created inside it. Absent defers to AWS's default (DISABLED).
	if spec.TagPropagation != nil {
		args.TagPropagation = pulumi.StringPtr(spec.GetTagPropagation())
	}

	// What happens to the domain's auto-created EFS file system on destroy.
	// AWS's default is Retain, which leaves a billing orphan behind -- the
	// spec surfaces the decision so ephemeral domains can opt into Delete.
	if spec.HomeEfsRetentionPolicy != nil {
		args.RetentionPolicy = &sagemaker.DomainRetentionPolicyArgs{
			HomeEfsFileSystem: pulumi.StringPtr(spec.GetHomeEfsRetentionPolicy()),
		}
	}

	if domainSettings := buildDomainSettings(spec); domainSettings != nil {
		args.DomainSettings = domainSettings
	}

	if spec.DefaultSpaceSettings != nil {
		args.DefaultSpaceSettings = buildDefaultSpaceSettings(spec.DefaultSpaceSettings)
	}

	createdDomain, err := sagemaker.NewDomain(ctx, "sagemaker-domain", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create sagemaker domain")
	}

	return createdDomain, nil
}

// buildDomainSettings assembles the provider's domain_settings block from the
// spec's top-level domain-administration fields. The spec models them flat
// (everything in a Domain spec is a "domain setting" -- the wrapper adds no
// information for manifest authors); the block is reconstructed here and only
// sent when at least one dial is set.
func buildDomainSettings(spec *awssagemakerdomainv1alpha1.AwsSagemakerDomainSpec) *sagemaker.DomainDomainSettingsArgs {
	hasDomainSettings := false
	dsArgs := &sagemaker.DomainDomainSettingsArgs{}

	if len(spec.DomainSecurityGroupIds) > 0 {
		var sgIds pulumi.StringArray
		for _, sg := range spec.DomainSecurityGroupIds {
			sgIds = append(sgIds, pulumi.String(sg.GetValue()))
		}
		dsArgs.SecurityGroupIds = sgIds
		hasDomainSettings = true
	}

	if spec.DockerSettings != nil {
		docker := spec.DockerSettings
		dockerArgs := &sagemaker.DomainDomainSettingsDockerSettingsArgs{}
		if docker.EnableDockerAccess != "" {
			dockerArgs.EnableDockerAccess = pulumi.StringPtr(docker.EnableDockerAccess)
		}
		if len(docker.VpcOnlyTrustedAccounts) > 0 {
			var accounts pulumi.StringArray
			for _, acct := range docker.VpcOnlyTrustedAccounts {
				accounts = append(accounts, pulumi.String(acct))
			}
			dockerArgs.VpcOnlyTrustedAccounts = accounts
		}
		dsArgs.DockerSettings = dockerArgs
		hasDomainSettings = true
	}

	// USER_PROFILE_NAME stamps each session's sts:SourceIdentity with the
	// acting user profile, so CloudTrail can attribute actions taken through
	// the shared execution role to a human.
	if spec.ExecutionRoleIdentityConfig != nil {
		dsArgs.ExecutionRoleIdentityConfig = pulumi.StringPtr(spec.GetExecutionRoleIdentityConfig())
		hasDomainSettings = true
	}

	// RStudio (Posit) Workbench activation: configuring the domain execution
	// role is what turns the RStudio app plane on for the domain.
	if rstudio := spec.RStudioServerProDomainSettings; rstudio != nil {
		rsArgs := &sagemaker.DomainDomainSettingsRStudioServerProDomainSettingsArgs{
			DomainExecutionRoleArn: pulumi.String(rstudio.DomainExecutionRoleArn.GetValue()),
		}
		if rstudio.RStudioConnectUrl != "" {
			rsArgs.RStudioConnectUrl = pulumi.StringPtr(rstudio.RStudioConnectUrl)
		}
		if rstudio.RStudioPackageManagerUrl != "" {
			rsArgs.RStudioPackageManagerUrl = pulumi.StringPtr(rstudio.RStudioPackageManagerUrl)
		}
		if rs := rstudio.DefaultResourceSpec; rs != nil {
			specArgs := &sagemaker.DomainDomainSettingsRStudioServerProDomainSettingsDefaultResourceSpecArgs{}
			if rs.InstanceType != "" {
				specArgs.InstanceType = pulumi.StringPtr(rs.InstanceType)
			}
			if rs.LifecycleConfigArn != "" {
				specArgs.LifecycleConfigArn = pulumi.StringPtr(rs.LifecycleConfigArn)
			}
			if rs.SagemakerImageArn != "" {
				specArgs.SagemakerImageArn = pulumi.StringPtr(rs.SagemakerImageArn)
			}
			if rs.SagemakerImageVersionAlias != "" {
				specArgs.SagemakerImageVersionAlias = pulumi.StringPtr(rs.SagemakerImageVersionAlias)
			}
			if rs.SagemakerImageVersionArn != "" {
				specArgs.SagemakerImageVersionArn = pulumi.StringPtr(rs.SagemakerImageVersionArn)
			}
			rsArgs.DefaultResourceSpec = specArgs
		}
		dsArgs.RStudioServerProDomainSettings = rsArgs
		hasDomainSettings = true
	}

	// Trusted identity propagation: the spec's CEL guarantees ENABLED only
	// ever reaches AWS on an SSO domain (AWS rejects it under IAM auth).
	if spec.TrustedIdentityPropagationStatus != nil {
		dsArgs.TrustedIdentityPropagationSettings = &sagemaker.DomainDomainSettingsTrustedIdentityPropagationSettingsArgs{
			Status: pulumi.String(spec.GetTrustedIdentityPropagationStatus()),
		}
		hasDomainSettings = true
	}

	if !hasDomainSettings {
		return nil
	}

	return dsArgs
}
