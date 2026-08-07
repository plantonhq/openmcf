package module

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/artifactregistry"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func repository(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpArtifactRegistryRepo.Spec

	// Enable the Artifact Registry API — the control plane that owns
	// repositories. DisableOnDestroy stays false: tearing down one
	// repository must never disable the API for everything else in the
	// project (other repositories keep serving pulls for running
	// workloads).
	artifactregistryApiArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("artifactregistry.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
		DisableOnDestroy:         pulumi.BoolPtr(false),
	}
	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		artifactregistryApiArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdArtifactregistryApi, err := projects.NewService(ctx,
		"gcpart-artifactregistry.googleapis.com", artifactregistryApiArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable artifactregistry.googleapis.com api")
	}

	// The Artifact Registry repository — one repository, one format, one
	// serving mode. Sharp edges, all taught by the API rather than
	// invented here:
	//
	//   - format, mode, location, project, and kms_key_name are immutable:
	//     changing any of them replaces the repository AND everything
	//     stored in it. There is no in-place migration between formats or
	//     modes.
	//
	//   - The remote_repository_config block is immutable as a whole, but
	//     the upstream credentials and disable_upstream_validation inside
	//     it are in the API's update mask — credential rotation updates in
	//     place.
	//
	//   - Deleting a repository deletes every artifact version in it.
	//     Unlike GCS buckets there is no force_destroy gate; protect
	//     precious repositories with KEEP cleanup policies and IAM, not
	//     with the module.
	args := &artifactregistry.RepositoryArgs{
		RepositoryId: pulumi.String(locals.RepositoryId),
		Location:     pulumi.String(spec.Location),
		Format:       pulumi.String(spec.Format),
		Labels:       pulumi.ToStringMap(locals.GcpLabels),

		// PARITY: the bridged provider tracks a newer line than the
		// released Terraform google provider and carries a client-side
		// deletion_policy knob (DELETE / ABANDON / PREVENT) that the
		// released 6.x TF resource does not have. Its default, DELETE,
		// is exactly the released TF destroy behavior (delete the
		// repository and its artifacts). Pin it explicitly so a future
		// bridged-default change can never make destroy behave
		// differently on one engine only.
		DeletionPolicy: pulumi.StringPtr("DELETE"),
	}

	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
	}

	// Mode defaults to STANDARD_REPOSITORY server-side; send only when set.
	if spec.Mode != "" {
		args.Mode = pulumi.StringPtr(spec.Mode)
	}

	if spec.Description != "" {
		args.Description = pulumi.StringPtr(spec.Description)
	}

	// CMEK: the Artifact Registry service agent must hold
	// roles/cloudkms.cryptoKeyEncrypterDecrypter on this key before create.
	if spec.KmsKeyName.GetValue() != "" {
		args.KmsKeyName = pulumi.StringPtr(spec.KmsKeyName.GetValue())
	}

	if spec.DockerConfig != nil {
		args.DockerConfig = &artifactregistry.RepositoryDockerConfigArgs{
			ImmutableTags: pulumi.BoolPtr(spec.DockerConfig.ImmutableTags),
		}
	}

	if spec.MavenConfig != nil {
		mavenArgs := &artifactregistry.RepositoryMavenConfigArgs{
			AllowSnapshotOverwrites: pulumi.BoolPtr(spec.MavenConfig.AllowSnapshotOverwrites),
		}
		if spec.MavenConfig.VersionPolicy != "" {
			mavenArgs.VersionPolicy = pulumi.StringPtr(spec.MavenConfig.VersionPolicy)
		}
		args.MavenConfig = mavenArgs
	}

	// Cleanup: DELETE policies remove matching versions; KEEP policies
	// protect them (KEEP wins on overlap). Dry-run logs matches without
	// deleting.
	if spec.CleanupPolicyDryRun {
		args.CleanupPolicyDryRun = pulumi.BoolPtr(true)
	}

	if len(spec.CleanupPolicies) > 0 {
		policies := artifactregistry.RepositoryCleanupPolicyArray{}
		for _, policy := range spec.CleanupPolicies {
			policyArgs := &artifactregistry.RepositoryCleanupPolicyArgs{
				Id:     pulumi.String(policy.Id),
				Action: pulumi.StringPtr(policy.Action),
			}
			if policy.Condition != nil {
				conditionArgs := &artifactregistry.RepositoryCleanupPolicyConditionArgs{}
				if policy.Condition.NewerThan != "" {
					conditionArgs.NewerThan = pulumi.StringPtr(policy.Condition.NewerThan)
				}
				if policy.Condition.OlderThan != "" {
					conditionArgs.OlderThan = pulumi.StringPtr(policy.Condition.OlderThan)
				}
				if len(policy.Condition.PackageNamePrefixes) > 0 {
					conditionArgs.PackageNamePrefixes = pulumi.ToStringArray(policy.Condition.PackageNamePrefixes)
				}
				if len(policy.Condition.TagPrefixes) > 0 {
					conditionArgs.TagPrefixes = pulumi.ToStringArray(policy.Condition.TagPrefixes)
				}
				if policy.Condition.TagState != "" {
					conditionArgs.TagState = pulumi.StringPtr(policy.Condition.TagState)
				}
				if len(policy.Condition.VersionNamePrefixes) > 0 {
					conditionArgs.VersionNamePrefixes = pulumi.ToStringArray(policy.Condition.VersionNamePrefixes)
				}
				policyArgs.Condition = conditionArgs
			}
			if policy.MostRecentVersions != nil {
				mrvArgs := &artifactregistry.RepositoryCleanupPolicyMostRecentVersionsArgs{}
				if policy.MostRecentVersions.KeepCount > 0 {
					mrvArgs.KeepCount = pulumi.IntPtr(int(policy.MostRecentVersions.KeepCount))
				}
				if len(policy.MostRecentVersions.PackageNamePrefixes) > 0 {
					mrvArgs.PackageNamePrefixes = pulumi.ToStringArray(policy.MostRecentVersions.PackageNamePrefixes)
				}
				policyArgs.MostRecentVersions = mrvArgs
			}
			policies = append(policies, policyArgs)
		}
		args.CleanupPolicies = policies
	}

	// REMOTE_REPOSITORY: a pull-through cache of exactly one upstream. The
	// spec enforces mode↔config coherence and exactly-one-upstream
	// pre-deploy.
	if spec.RemoteRepositoryConfig != nil {
		remote := spec.RemoteRepositoryConfig
		remoteArgs := &artifactregistry.RepositoryRemoteRepositoryConfigArgs{
			DisableUpstreamValidation: pulumi.BoolPtr(remote.DisableUpstreamValidation),
		}
		if remote.Description != "" {
			remoteArgs.Description = pulumi.StringPtr(remote.Description)
		}
		if remote.DockerPublicRepository != "" {
			remoteArgs.DockerRepository = &artifactregistry.RepositoryRemoteRepositoryConfigDockerRepositoryArgs{
				PublicRepository: pulumi.StringPtr(remote.DockerPublicRepository),
			}
		}
		if remote.MavenPublicRepository != "" {
			remoteArgs.MavenRepository = &artifactregistry.RepositoryRemoteRepositoryConfigMavenRepositoryArgs{
				PublicRepository: pulumi.StringPtr(remote.MavenPublicRepository),
			}
		}
		if remote.NpmPublicRepository != "" {
			remoteArgs.NpmRepository = &artifactregistry.RepositoryRemoteRepositoryConfigNpmRepositoryArgs{
				PublicRepository: pulumi.StringPtr(remote.NpmPublicRepository),
			}
		}
		if remote.PythonPublicRepository != "" {
			remoteArgs.PythonRepository = &artifactregistry.RepositoryRemoteRepositoryConfigPythonRepositoryArgs{
				PublicRepository: pulumi.StringPtr(remote.PythonPublicRepository),
			}
		}
		if remote.AptRepository != nil {
			remoteArgs.AptRepository = &artifactregistry.RepositoryRemoteRepositoryConfigAptRepositoryArgs{
				PublicRepository: &artifactregistry.RepositoryRemoteRepositoryConfigAptRepositoryPublicRepositoryArgs{
					RepositoryBase: pulumi.String(remote.AptRepository.RepositoryBase),
					RepositoryPath: pulumi.String(remote.AptRepository.RepositoryPath),
				},
			}
		}
		if remote.YumRepository != nil {
			remoteArgs.YumRepository = &artifactregistry.RepositoryRemoteRepositoryConfigYumRepositoryArgs{
				PublicRepository: &artifactregistry.RepositoryRemoteRepositoryConfigYumRepositoryPublicRepositoryArgs{
					RepositoryBase: pulumi.String(remote.YumRepository.RepositoryBase),
					RepositoryPath: pulumi.String(remote.YumRepository.RepositoryPath),
				},
			}
		}
		// Custom upstream: another AR repository or any registry URI.
		if remote.CommonRepository != nil && remote.CommonRepository.Uri.GetValue() != "" {
			remoteArgs.CommonRepository = &artifactregistry.RepositoryRemoteRepositoryConfigCommonRepositoryArgs{
				Uri: pulumi.String(remote.CommonRepository.Uri.GetValue()),
			}
		}
		// Credential rotation updates in place — the password itself lives
		// in Secret Manager; only the secret-version PATH passes through
		// here.
		if remote.UpstreamCredentials != nil {
			remoteArgs.UpstreamCredentials = &artifactregistry.RepositoryRemoteRepositoryConfigUpstreamCredentialsArgs{
				UsernamePasswordCredentials: &artifactregistry.RepositoryRemoteRepositoryConfigUpstreamCredentialsUsernamePasswordCredentialsArgs{
					Username:              pulumi.StringPtr(remote.UpstreamCredentials.Username),
					PasswordSecretVersion: pulumi.StringPtr(remote.UpstreamCredentials.PasswordSecretVersion),
				},
			}
		}
		args.RemoteRepositoryConfig = remoteArgs
	}

	// VIRTUAL_REPOSITORY: priority-ordered aggregation of other AR
	// repositories (highest priority wins on conflicts).
	if spec.VirtualRepositoryConfig != nil {
		policies := artifactregistry.RepositoryVirtualRepositoryConfigUpstreamPolicyArray{}
		for _, policy := range spec.VirtualRepositoryConfig.UpstreamPolicies {
			policies = append(policies, &artifactregistry.RepositoryVirtualRepositoryConfigUpstreamPolicyArgs{
				Id:         pulumi.StringPtr(policy.Id),
				Repository: pulumi.StringPtr(policy.Repository.GetValue()),
				Priority:   pulumi.IntPtr(int(policy.Priority)),
			})
		}
		args.VirtualRepositoryConfig = &artifactregistry.RepositoryVirtualRepositoryConfigArgs{
			UpstreamPolicies: policies,
		}
	}

	if spec.VulnerabilityScanningEnablement != "" {
		args.VulnerabilityScanningConfig = &artifactregistry.RepositoryVulnerabilityScanningConfigArgs{
			EnablementConfig: pulumi.StringPtr(spec.VulnerabilityScanningEnablement),
		}
	}

	createdRepository, err := artifactregistry.NewRepository(ctx, "repository", args,
		pulumi.Provider(gcpProvider),
		pulumi.DependsOn([]pulumi.Resource{createdArtifactregistryApi}))
	if err != nil {
		return errors.Wrap(err, "failed to create artifact registry repository")
	}

	// Additive IAM grants: one (role, member) pair per resource, merging
	// into the repository's policy without touching grants made elsewhere —
	// authoritative bindings/policies are deliberately not used. Resource
	// names key on (role, member) — the grant's identity — matching the
	// Terraform module's for_each keys.
	for _, iamMember := range spec.IamMembers {
		member := iamMember.Member.GetValue()
		iamArgs := &artifactregistry.RepositoryIamMemberArgs{
			Repository: createdRepository.RepositoryId,
			Location:   createdRepository.Location,
			Role:       pulumi.String(iamMember.Role),
			Member:     pulumi.String(member),
		}
		if spec.ProjectId.GetValue() != "" {
			iamArgs.Project = pulumi.StringPtr(spec.ProjectId.GetValue())
		}
		if iamMember.Condition != nil {
			conditionArgs := &artifactregistry.RepositoryIamMemberConditionArgs{
				Title:      pulumi.String(iamMember.Condition.Title),
				Expression: pulumi.String(iamMember.Condition.Expression),
			}
			if iamMember.Condition.Description != "" {
				conditionArgs.Description = pulumi.StringPtr(iamMember.Condition.Description)
			}
			iamArgs.Condition = conditionArgs
		}
		_, err := artifactregistry.NewRepositoryIamMember(ctx,
			fmt.Sprintf("iam-member-%s-%s", iamMember.Role, member), iamArgs,
			pulumi.Provider(gcpProvider),
			pulumi.Parent(createdRepository))
		if err != nil {
			return errors.Wrapf(err, "failed to grant %s to %s on the repository", iamMember.Role, member)
		}
	}

	// Short name of the repository (the repository ID).
	ctx.Export(OpName, createdRepository.Name)

	// The resource ID is the fully qualified repository path
	// (projects/{project}/locations/{location}/repositories/{repo}) — the
	// exact string every composing resource consumes: a Cloud Function's
	// docker_repository, a virtual repository's upstream policy, and a
	// remote repository's common upstream.
	ctx.Export(OpRepositoryPath, createdRepository.ID())

	// The registry endpoint clients push to and pull from. Constructed
	// from resolved attributes ({location}-{format}.pkg.dev/{project}/
	// {repo}) because the released 6.x Terraform provider does not export
	// a registry URI attribute — the Terraform module builds the identical
	// string.
	ctx.Export(OpRegistryUri, pulumi.All(
		createdRepository.Location,
		createdRepository.Format,
		createdRepository.Project,
		createdRepository.Name,
	).ApplyT(func(parts []interface{}) string {
		return fmt.Sprintf("%s-%s.pkg.dev/%s/%s",
			strings.ToLower(parts[0].(string)),
			strings.ToLower(parts[1].(string)),
			parts[2].(string),
			parts[3].(string),
		)
	}).(pulumi.StringOutput))

	ctx.Export(OpLocation, createdRepository.Location)

	return nil
}
