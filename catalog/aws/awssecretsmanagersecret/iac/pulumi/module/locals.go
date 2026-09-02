package module

import (
	"strconv"

	awssecretsmanagersecretv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssecretsmanagersecret/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/aws/awstagkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds pre-computed values derived from the stack input.
type Locals struct {
	Target *awssecretsmanagersecretv1alpha1.AwsSecretsManagerSecret
	Spec   *awssecretsmanagersecretv1alpha1.AwsSecretsManagerSecretSpec

	// SecretName is spec.secret_name when set (hierarchical paths and
	// service-required prefixes like "ecr-pullthroughcache/..." that
	// metadata.name cannot carry), else metadata.name. The AWS secret name
	// is create-time immutable, and both engines share this resolution so
	// a manifest deploys identically on either.
	SecretName string

	AwsTags map[string]string
}

func initializeLocals(_ *pulumi.Context, in *awssecretsmanagersecretv1alpha1.AwsSecretsManagerSecretStackInput) *Locals {
	locals := &Locals{}
	locals.Target = in.Target
	locals.Spec = in.Target.Spec

	metadata := in.Target.Metadata
	locals.SecretName = metadata.Name
	if in.Target.Spec.GetSecretName() != "" {
		locals.SecretName = in.Target.Spec.GetSecretName()
	}

	// Resource-identity tags match the Terraform module key-for-key.
	locals.AwsTags = map[string]string{
		awstagkeys.Name:         metadata.Name,
		awstagkeys.Resource:     strconv.FormatBool(true),
		awstagkeys.Organization: metadata.Org,
		awstagkeys.Environment:  metadata.Env,
		awstagkeys.ResourceKind: cloudresourcekind.CloudResourceKind_AwsSecretsManagerSecret.String(),
		awstagkeys.ResourceId:   metadata.Id,
	}

	return locals
}
