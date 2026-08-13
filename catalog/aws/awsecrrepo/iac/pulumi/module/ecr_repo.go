package module

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/ecr"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// ecrRepo creates the ECR repository plus its two folded 1:1 satellites (the
// lifecycle policy and the repository policy) from AwsEcrRepoSpec.
func ecrRepo(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.AwsEcrRepo.Spec

	// Tag mutability — the spec default is MUTABLE (the AWS default). The
	// exclusion filters only exist in the *_WITH_EXCLUSION modes
	// (CEL-enforced); WILDCARD is the only filter type AWS supports today, so
	// it is materialized here rather than modeled in the spec.
	imageTagMutability := spec.GetImageTagMutability()
	if imageTagMutability == "" {
		imageTagMutability = "MUTABLE"
	}
	var exclusionFilters ecr.RepositoryImageTagMutabilityExclusionFilterArray
	for _, filter := range spec.ImageTagMutabilityExclusionFilters {
		exclusionFilters = append(exclusionFilters, &ecr.RepositoryImageTagMutabilityExclusionFilterArgs{
			Filter:     pulumi.String(filter),
			FilterType: pulumi.String("WILDCARD"),
		})
	}

	// Encryption — the whole configuration is create-time (ForceNew). A
	// customer-managed key only applies to the KMS types (CEL-enforced).
	encryptionType := spec.GetEncryptionType()
	if encryptionType == "" {
		encryptionType = "AES256"
	}
	encryptionConfig := &ecr.RepositoryEncryptionConfigurationArgs{
		EncryptionType: pulumi.String(encryptionType),
	}
	if spec.KmsKeyId.GetValue() != "" {
		encryptionConfig.KmsKey = pulumi.String(spec.KmsKeyId.GetValue())
	}

	// scan_on_push defaults to true (the spec's recommended security posture).
	scanOnPush := true
	if spec.ScanOnPush != nil {
		scanOnPush = spec.GetScanOnPush()
	}

	repo, err := ecr.NewRepository(ctx, locals.AwsEcrRepo.Metadata.Name, &ecr.RepositoryArgs{
		// The repository name is the immutable registry path (ForceNew).
		// Changing it replaces the repository; images are NOT migrated.
		Name:                               pulumi.String(spec.RepositoryName),
		ImageTagMutability:                 pulumi.String(imageTagMutability),
		ImageTagMutabilityExclusionFilters: exclusionFilters,
		ImageScanningConfiguration: &ecr.RepositoryImageScanningConfigurationArgs{
			ScanOnPush: pulumi.Bool(scanOnPush),
		},
		ForceDelete: pulumi.Bool(spec.ForceDelete),
		EncryptionConfigurations: ecr.RepositoryEncryptionConfigurationArray{
			encryptionConfig,
		},
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "unable to create ECR repository")
	}

	// Lifecycle policy — a separate AWS API resource keyed 1:1 by the
	// repository, folded into the spec as structured rules.
	if len(spec.LifecycleRules) > 0 {
		if err := createLifecyclePolicy(ctx, locals, repo, provider); err != nil {
			return errors.Wrap(err, "unable to create lifecycle policy")
		}
	}

	// Repository policy — resource-based access control (cross-account
	// pulls, service principals like Lambda). Also a 1:1 folded satellite.
	if spec.RepositoryPolicy != nil {
		policyJson, err := spec.RepositoryPolicy.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "failed to marshal repository policy")
		}
		_, err = ecr.NewRepositoryPolicy(ctx, fmt.Sprintf("%s-policy", locals.AwsEcrRepo.Metadata.Name), &ecr.RepositoryPolicyArgs{
			Repository: repo.Name,
			Policy:     pulumi.String(policyJson),
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "failed to create repository policy")
		}
	}

	ctx.Export(OpEcrRepoName, repo.Name)
	ctx.Export(OpEcrRepoUrl, repo.RepositoryUrl)
	ctx.Export(OpEcrRepoArn, repo.Arn)
	ctx.Export(OpEcrRepoRegistryId, repo.RegistryId)

	return nil
}

// createLifecyclePolicy rebuilds the exact lifecycle policy JSON document AWS
// expects from the spec's structured rules. The "days" count unit exists for
// all three "since..." count types and is materialized here. Rules either
// expire images (the default) or transition them to the archive storage
// class; the transition/target coupling and the storageClass/countType
// pairing are CEL-enforced. Optional members are added only when set — both
// engines must hand AWS the SAME document. tagged rules carry exactly one
// selector list (CEL-enforced); untagged/any rules carry none.
func createLifecyclePolicy(ctx *pulumi.Context, locals *Locals, repo *ecr.Repository, provider *aws.Provider) error {
	rules := make([]map[string]interface{}, 0, len(locals.AwsEcrRepo.Spec.LifecycleRules))
	for _, rule := range locals.AwsEcrRepo.Spec.LifecycleRules {
		selection := map[string]interface{}{
			"tagStatus":   rule.TagStatus,
			"countType":   rule.CountType,
			"countNumber": rule.CountNumber,
		}
		switch rule.CountType {
		case "sinceImagePushed", "sinceImagePulled", "sinceImageTransitioned":
			selection["countUnit"] = "days"
		}
		if rule.StorageClass != "" {
			selection["storageClass"] = rule.StorageClass
		}
		if len(rule.TagPrefixes) > 0 {
			selection["tagPrefixList"] = rule.TagPrefixes
		}
		if len(rule.TagPatterns) > 0 {
			selection["tagPatternList"] = rule.TagPatterns
		}

		action := map[string]interface{}{"type": "expire"}
		if rule.ActionType != "" {
			action["type"] = rule.ActionType
		}
		if rule.TargetStorageClass != "" {
			action["targetStorageClass"] = rule.TargetStorageClass
		}

		entry := map[string]interface{}{
			"rulePriority": rule.RulePriority,
			"selection":    selection,
			"action":       action,
		}
		if rule.Description != "" {
			entry["description"] = rule.Description
		}
		rules = append(rules, entry)
	}

	policyJson, err := json.Marshal(map[string]interface{}{"rules": rules})
	if err != nil {
		return errors.Wrap(err, "failed to marshal lifecycle policy JSON")
	}

	_, err = ecr.NewLifecyclePolicy(ctx, fmt.Sprintf("%s-lifecycle", locals.AwsEcrRepo.Metadata.Name), &ecr.LifecyclePolicyArgs{
		Repository: repo.Name,
		Policy:     pulumi.String(policyJson),
	}, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create lifecycle policy")
	}

	return nil
}
