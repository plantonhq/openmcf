package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	awscodebuildprojectv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awscodebuildproject/v1"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/codebuild"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// lambdaEnvironmentTypes are the environment types AWS caps itself: build and
// queued timeouts are not supported there (the spec's CEL rejects explicit
// values; this guard additionally keeps the spec-level defaults from being
// sent, matching the Terraform module).
var lambdaEnvironmentTypes = map[string]bool{
	"LINUX_LAMBDA_CONTAINER": true,
	"ARM_LAMBDA_CONTAINER":   true,
}

// project provisions the CodeBuild project -- a metadata-only control-plane
// resource: creating it provisions nothing until a build starts, so
// create/update/delete are near-instant. The only operational wait is IAM
// eventual consistency on a freshly created service role, which the provider
// absorbs with a bounded retry.
func project(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
) (*codebuild.Project, error) {
	spec := locals.AwsCodeBuildProject.Spec

	// --- Environment ---
	env := spec.Environment
	envArgs := &codebuild.ProjectEnvironmentArgs{
		Type:        pulumi.String(env.Type),
		ComputeType: pulumi.String(env.ComputeType),
		Image:       pulumi.String(env.Image),
	}
	if env.Certificate != "" {
		envArgs.Certificate = pulumi.StringPtr(env.Certificate)
	}
	if env.PrivilegedMode {
		envArgs.PrivilegedMode = pulumi.BoolPtr(true)
	}
	if env.GetImagePullCredentialsType() != "" && env.GetImagePullCredentialsType() != "CODEBUILD" {
		envArgs.ImagePullCredentialsType = pulumi.StringPtr(env.GetImagePullCredentialsType())
	}
	if len(env.EnvironmentVariables) > 0 {
		var envVars codebuild.ProjectEnvironmentEnvironmentVariableArray
		for _, ev := range env.EnvironmentVariables {
			evArgs := &codebuild.ProjectEnvironmentEnvironmentVariableArgs{
				Name:  pulumi.String(ev.Name),
				Value: pulumi.String(ev.Value),
			}
			if ev.GetType() != "" && ev.GetType() != "PLAINTEXT" {
				evArgs.Type = pulumi.StringPtr(ev.GetType())
			}
			envVars = append(envVars, evArgs)
		}
		envArgs.EnvironmentVariables = envVars
	}
	if env.RegistryCredential != nil {
		envArgs.RegistryCredential = &codebuild.ProjectEnvironmentRegistryCredentialArgs{
			Credential:         pulumi.String(env.RegistryCredential.Credential),
			CredentialProvider: pulumi.String(env.RegistryCredential.CredentialProvider),
		}
	}
	// Persistent, dedicated Docker server: layer state survives across
	// builds, unlike the per-build daemon privileged_mode provides.
	if env.DockerServer != nil {
		dsArgs := &codebuild.ProjectEnvironmentDockerServerArgs{
			ComputeType: pulumi.String(env.DockerServer.ComputeType),
		}
		if len(env.DockerServer.SecurityGroupIds) > 0 {
			var sgIds pulumi.StringArray
			for _, sg := range env.DockerServer.SecurityGroupIds {
				if sg.GetValue() != "" {
					sgIds = append(sgIds, pulumi.String(sg.GetValue()))
				}
			}
			dsArgs.SecurityGroupIds = sgIds
		}
		envArgs.DockerServer = dsArgs
	}
	// Reserved-capacity fleet membership -- pre-provisioned, always-warm
	// build machines. The fleet is a shared account-level resource; the
	// project only references its ARN.
	if env.FleetArn != "" {
		envArgs.Fleet = &codebuild.ProjectEnvironmentFleetArgs{
			FleetArn: pulumi.StringPtr(env.FleetArn),
		}
	}

	// --- Project args ---
	args := &codebuild.ProjectArgs{
		Name:        pulumi.StringPtr(locals.ProjectName),
		ServiceRole: pulumi.String(spec.ServiceRole.GetValue()),
		Source:      buildSourceArgs(spec.Source),
		Environment: envArgs,
		Artifacts:   buildArtifactsArgs(spec.Artifacts),
		Tags:        pulumi.ToStringMap(locals.Labels),
	}

	// --- Secondary sources / versions / artifacts ---
	if len(spec.SecondarySources) > 0 {
		var secondaries codebuild.ProjectSecondarySourceArray
		for _, s := range spec.SecondarySources {
			secondaries = append(secondaries, buildSecondarySourceArgs(s))
		}
		args.SecondarySources = secondaries
	}
	if len(spec.SecondarySourceVersions) > 0 {
		var versions codebuild.ProjectSecondarySourceVersionArray
		for _, v := range spec.SecondarySourceVersions {
			versions = append(versions, &codebuild.ProjectSecondarySourceVersionArgs{
				SourceIdentifier: pulumi.String(v.SourceIdentifier),
				SourceVersion:    pulumi.String(v.SourceVersion),
			})
		}
		args.SecondarySourceVersions = versions
	}
	if len(spec.SecondaryArtifacts) > 0 {
		var secondaries codebuild.ProjectSecondaryArtifactArray
		for _, a := range spec.SecondaryArtifacts {
			secondaries = append(secondaries, buildSecondaryArtifactArgs(a))
		}
		args.SecondaryArtifacts = secondaries
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}
	if spec.EncryptionKey != nil && spec.EncryptionKey.GetValue() != "" {
		args.EncryptionKey = pulumi.StringPtr(spec.EncryptionKey.GetValue())
	}
	// Lambda environments ignore timeouts entirely (AWS caps them itself);
	// sending the spec defaults would create a permanent diff there.
	if !lambdaEnvironmentTypes[env.Type] {
		if spec.GetBuildTimeout() != 0 {
			args.BuildTimeout = pulumi.IntPtr(int(spec.GetBuildTimeout()))
		}
		if spec.GetQueuedTimeout() != 0 {
			args.QueuedTimeout = pulumi.IntPtr(int(spec.GetQueuedTimeout()))
		}
	}
	if spec.ConcurrentBuildLimit > 0 {
		args.ConcurrentBuildLimit = pulumi.IntPtr(int(spec.ConcurrentBuildLimit))
	}
	// Additional automatic retries after a failed build (not total attempts).
	if spec.AutoRetryLimit > 0 {
		args.AutoRetryLimit = pulumi.IntPtr(int(spec.AutoRetryLimit))
	}
	if spec.BadgeEnabled {
		args.BadgeEnabled = pulumi.BoolPtr(true)
	}
	if spec.SourceVersion != "" {
		args.SourceVersion = pulumi.StringPtr(spec.SourceVersion)
	}
	// Public visibility is managed by a separate AWS API call under the hood
	// (UpdateProjectVisibility); the provider sequences it after project
	// create/update. The resource access role is what CodeBuild uses to read
	// the logs/artifacts it re-exposes publicly.
	if spec.GetProjectVisibility() != "" && spec.GetProjectVisibility() != "PRIVATE" {
		args.ProjectVisibility = pulumi.StringPtr(spec.GetProjectVisibility())
	}
	if spec.ResourceAccessRole != nil && spec.ResourceAccessRole.GetValue() != "" {
		args.ResourceAccessRole = pulumi.StringPtr(spec.ResourceAccessRole.GetValue())
	}

	// --- Cache ---
	// An absent cache block and an explicit NO_CACHE deploy identically; the
	// cache type is a presence-carrying optional, so unset means NO_CACHE.
	if spec.Cache != nil && spec.Cache.GetType() != "" && spec.Cache.GetType() != "NO_CACHE" {
		cacheArgs := &codebuild.ProjectCacheArgs{
			Type: pulumi.StringPtr(spec.Cache.GetType()),
		}
		if spec.Cache.Location != nil && spec.Cache.Location.GetValue() != "" {
			cacheArgs.Location = pulumi.StringPtr(spec.Cache.Location.GetValue())
		}
		if spec.Cache.GetType() == "LOCAL" && len(spec.Cache.Modes) > 0 {
			cacheArgs.Modes = pulumi.ToStringArray(spec.Cache.Modes)
		}
		if spec.Cache.CacheNamespace != "" {
			cacheArgs.CacheNamespace = pulumi.StringPtr(spec.Cache.CacheNamespace)
		}
		args.Cache = cacheArgs
	}

	// --- Logs config ---
	if spec.LogsConfig != nil {
		logsArgs := &codebuild.ProjectLogsConfigArgs{}
		if spec.LogsConfig.CloudwatchLogs != nil {
			cwArgs := &codebuild.ProjectLogsConfigCloudwatchLogsArgs{}
			if spec.LogsConfig.CloudwatchLogs.GetStatus() != "" {
				cwArgs.Status = pulumi.StringPtr(spec.LogsConfig.CloudwatchLogs.GetStatus())
			}
			if spec.LogsConfig.CloudwatchLogs.GroupName != nil && spec.LogsConfig.CloudwatchLogs.GroupName.GetValue() != "" {
				cwArgs.GroupName = pulumi.StringPtr(spec.LogsConfig.CloudwatchLogs.GroupName.GetValue())
			}
			if spec.LogsConfig.CloudwatchLogs.StreamName != "" {
				cwArgs.StreamName = pulumi.StringPtr(spec.LogsConfig.CloudwatchLogs.StreamName)
			}
			logsArgs.CloudwatchLogs = cwArgs
		}
		if spec.LogsConfig.S3Logs != nil {
			s3Args := &codebuild.ProjectLogsConfigS3LogsArgs{}
			if spec.LogsConfig.S3Logs.GetStatus() != "" {
				s3Args.Status = pulumi.StringPtr(spec.LogsConfig.S3Logs.GetStatus())
			}
			if spec.LogsConfig.S3Logs.Location != nil && spec.LogsConfig.S3Logs.Location.GetValue() != "" {
				s3Args.Location = pulumi.StringPtr(spec.LogsConfig.S3Logs.Location.GetValue())
			}
			if spec.LogsConfig.S3Logs.EncryptionDisabled {
				s3Args.EncryptionDisabled = pulumi.BoolPtr(true)
			}
			if spec.LogsConfig.S3Logs.BucketOwnerAccess != "" {
				s3Args.BucketOwnerAccess = pulumi.StringPtr(spec.LogsConfig.S3Logs.BucketOwnerAccess)
			}
			logsArgs.S3Logs = s3Args
		}
		args.LogsConfig = logsArgs
	}

	// --- VPC placement ---
	if spec.VpcConfig != nil {
		var subnetIds pulumi.StringArray
		for _, s := range spec.VpcConfig.SubnetIds {
			if s.GetValue() != "" {
				subnetIds = append(subnetIds, pulumi.String(s.GetValue()))
			}
		}
		var sgIds pulumi.StringArray
		for _, sg := range spec.VpcConfig.SecurityGroupIds {
			if sg.GetValue() != "" {
				sgIds = append(sgIds, pulumi.String(sg.GetValue()))
			}
		}
		args.VpcConfig = &codebuild.ProjectVpcConfigArgs{
			VpcId:            pulumi.String(spec.VpcConfig.VpcId.GetValue()),
			Subnets:          subnetIds,
			SecurityGroupIds: sgIds,
		}
	}

	// --- EFS mounts (shared caches that outlive individual builds) ---
	if len(spec.FileSystemLocations) > 0 {
		var fsLocations codebuild.ProjectFileSystemLocationArray
		for _, fs := range spec.FileSystemLocations {
			fsArgs := &codebuild.ProjectFileSystemLocationArgs{
				Type:       pulumi.StringPtr(fs.GetType()),
				Identifier: pulumi.StringPtr(fs.Identifier),
				Location:   pulumi.StringPtr(fs.Location),
				MountPoint: pulumi.StringPtr(fs.MountPoint),
			}
			if fs.MountOptions != "" {
				fsArgs.MountOptions = pulumi.StringPtr(fs.MountOptions)
			}
			fsLocations = append(fsLocations, fsArgs)
		}
		args.FileSystemLocations = fsLocations
	}

	// --- Batch builds ---
	if spec.BuildBatchConfig != nil {
		batchArgs := &codebuild.ProjectBuildBatchConfigArgs{
			ServiceRole: pulumi.String(spec.BuildBatchConfig.ServiceRole.GetValue()),
		}
		if spec.BuildBatchConfig.CombineArtifacts {
			batchArgs.CombineArtifacts = pulumi.BoolPtr(true)
		}
		if spec.BuildBatchConfig.TimeoutInMins > 0 {
			batchArgs.TimeoutInMins = pulumi.IntPtr(int(spec.BuildBatchConfig.TimeoutInMins))
		}
		if spec.BuildBatchConfig.Restrictions != nil {
			restrictionArgs := &codebuild.ProjectBuildBatchConfigRestrictionsArgs{}
			if len(spec.BuildBatchConfig.Restrictions.ComputeTypesAllowed) > 0 {
				restrictionArgs.ComputeTypesAlloweds = pulumi.ToStringArray(spec.BuildBatchConfig.Restrictions.ComputeTypesAllowed)
			}
			if spec.BuildBatchConfig.Restrictions.MaximumBuildsAllowed > 0 {
				restrictionArgs.MaximumBuildsAllowed = pulumi.IntPtr(int(spec.BuildBatchConfig.Restrictions.MaximumBuildsAllowed))
			}
			batchArgs.Restrictions = restrictionArgs
		}
		args.BuildBatchConfig = batchArgs
	}

	created, err := codebuild.NewProject(ctx, "codebuild-project", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "create codebuild project")
	}

	return created, nil
}

// resourcePolicy attaches the folded resource-based IAM policy to the
// project -- the cross-account access mechanism. One document per project,
// keyed by the project ARN.
func resourcePolicy(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
	proj *codebuild.Project,
) error {
	policyBytes, err := json.Marshal(locals.AwsCodeBuildProject.Spec.ResourcePolicy.AsMap())
	if err != nil {
		return errors.Wrap(err, "marshal resource policy to JSON")
	}

	_, err = codebuild.NewResourcePolicy(ctx, "codebuild-resource-policy", &codebuild.ResourcePolicyArgs{
		ResourceArn: proj.Arn,
		Policy:      pulumi.String(string(policyBytes)),
	}, pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{proj}))
	if err != nil {
		return errors.Wrap(err, "create codebuild resource policy")
	}

	return nil
}

// buildSourceCommon maps the shared source fields. The primary source and
// secondary sources use the same spec shape; only the SDK arg types differ.
func buildSourceArgs(src *awscodebuildprojectv1.AwsCodeBuildSource) *codebuild.ProjectSourceArgs {
	args := &codebuild.ProjectSourceArgs{
		Type: pulumi.String(src.Type),
	}
	if src.Location != "" {
		args.Location = pulumi.StringPtr(src.Location)
	}
	if src.Buildspec != "" {
		args.Buildspec = pulumi.StringPtr(src.Buildspec)
	}
	if src.GitCloneDepth > 0 {
		args.GitCloneDepth = pulumi.IntPtr(int(src.GitCloneDepth))
	}
	if src.InsecureSsl {
		args.InsecureSsl = pulumi.BoolPtr(true)
	}
	if src.ReportBuildStatus {
		args.ReportBuildStatus = pulumi.BoolPtr(true)
	}
	if src.GitSubmodulesConfig != nil {
		args.GitSubmodulesConfig = &codebuild.ProjectSourceGitSubmodulesConfigArgs{
			FetchSubmodules: pulumi.Bool(src.GitSubmodulesConfig.FetchSubmodules),
		}
	}
	if src.BuildStatusConfig != nil {
		bsArgs := &codebuild.ProjectSourceBuildStatusConfigArgs{}
		if src.BuildStatusConfig.Context != "" {
			bsArgs.Context = pulumi.StringPtr(src.BuildStatusConfig.Context)
		}
		if src.BuildStatusConfig.TargetUrl != "" {
			bsArgs.TargetUrl = pulumi.StringPtr(src.BuildStatusConfig.TargetUrl)
		}
		args.BuildStatusConfig = bsArgs
	}
	if src.Auth != nil {
		args.Auth = &codebuild.ProjectSourceAuthArgs{
			Type:     pulumi.String(src.Auth.Type),
			Resource: pulumi.String(src.Auth.Resource),
		}
	}
	return args
}

func buildSecondarySourceArgs(src *awscodebuildprojectv1.AwsCodeBuildSource) *codebuild.ProjectSecondarySourceArgs {
	args := &codebuild.ProjectSecondarySourceArgs{
		SourceIdentifier: pulumi.String(src.SourceIdentifier),
		Type:             pulumi.String(src.Type),
	}
	if src.Location != "" {
		args.Location = pulumi.StringPtr(src.Location)
	}
	if src.Buildspec != "" {
		args.Buildspec = pulumi.StringPtr(src.Buildspec)
	}
	if src.GitCloneDepth > 0 {
		args.GitCloneDepth = pulumi.IntPtr(int(src.GitCloneDepth))
	}
	if src.InsecureSsl {
		args.InsecureSsl = pulumi.BoolPtr(true)
	}
	if src.ReportBuildStatus {
		args.ReportBuildStatus = pulumi.BoolPtr(true)
	}
	if src.GitSubmodulesConfig != nil {
		args.GitSubmodulesConfig = &codebuild.ProjectSecondarySourceGitSubmodulesConfigArgs{
			FetchSubmodules: pulumi.Bool(src.GitSubmodulesConfig.FetchSubmodules),
		}
	}
	if src.BuildStatusConfig != nil {
		bsArgs := &codebuild.ProjectSecondarySourceBuildStatusConfigArgs{}
		if src.BuildStatusConfig.Context != "" {
			bsArgs.Context = pulumi.StringPtr(src.BuildStatusConfig.Context)
		}
		if src.BuildStatusConfig.TargetUrl != "" {
			bsArgs.TargetUrl = pulumi.StringPtr(src.BuildStatusConfig.TargetUrl)
		}
		args.BuildStatusConfig = bsArgs
	}
	if src.Auth != nil {
		args.Auth = &codebuild.ProjectSecondarySourceAuthArgs{
			Type:     pulumi.String(src.Auth.Type),
			Resource: pulumi.String(src.Auth.Resource),
		}
	}
	return args
}

func buildArtifactsArgs(art *awscodebuildprojectv1.AwsCodeBuildArtifacts) *codebuild.ProjectArtifactsArgs {
	args := &codebuild.ProjectArtifactsArgs{
		Type: pulumi.String(art.Type),
	}
	if art.ArtifactIdentifier != "" {
		args.ArtifactIdentifier = pulumi.StringPtr(art.ArtifactIdentifier)
	}
	if art.Location != nil && art.Location.GetValue() != "" {
		args.Location = pulumi.StringPtr(art.Location.GetValue())
	}
	if art.Name != "" {
		args.Name = pulumi.StringPtr(art.Name)
	}
	if art.Path != "" {
		args.Path = pulumi.StringPtr(art.Path)
	}
	if art.Packaging != "" {
		args.Packaging = pulumi.StringPtr(art.Packaging)
	}
	if art.NamespaceType != "" {
		args.NamespaceType = pulumi.StringPtr(art.NamespaceType)
	}
	if art.EncryptionDisabled {
		args.EncryptionDisabled = pulumi.BoolPtr(true)
	}
	if art.OverrideArtifactName {
		args.OverrideArtifactName = pulumi.BoolPtr(true)
	}
	if art.BucketOwnerAccess != "" {
		args.BucketOwnerAccess = pulumi.StringPtr(art.BucketOwnerAccess)
	}
	return args
}

func buildSecondaryArtifactArgs(art *awscodebuildprojectv1.AwsCodeBuildArtifacts) *codebuild.ProjectSecondaryArtifactArgs {
	args := &codebuild.ProjectSecondaryArtifactArgs{
		ArtifactIdentifier: pulumi.String(art.ArtifactIdentifier),
		Type:               pulumi.String(art.Type),
	}
	if art.Location != nil && art.Location.GetValue() != "" {
		args.Location = pulumi.StringPtr(art.Location.GetValue())
	}
	if art.Name != "" {
		args.Name = pulumi.StringPtr(art.Name)
	}
	if art.Path != "" {
		args.Path = pulumi.StringPtr(art.Path)
	}
	if art.Packaging != "" {
		args.Packaging = pulumi.StringPtr(art.Packaging)
	}
	if art.NamespaceType != "" {
		args.NamespaceType = pulumi.StringPtr(art.NamespaceType)
	}
	if art.EncryptionDisabled {
		args.EncryptionDisabled = pulumi.BoolPtr(true)
	}
	if art.OverrideArtifactName {
		args.OverrideArtifactName = pulumi.BoolPtr(true)
	}
	if art.BucketOwnerAccess != "" {
		args.BucketOwnerAccess = pulumi.StringPtr(art.BucketOwnerAccess)
	}
	return args
}
