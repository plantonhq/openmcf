package module

import (
	awss3objectsetv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awss3objectset/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	AwsS3ObjectSet *awss3objectsetv1.AwsS3ObjectSet
	Labels         map[string]string
}

func initializeLocals(ctx *pulumi.Context, stackInput *awss3objectsetv1.AwsS3ObjectSetStackInput) *Locals {
	locals := &Locals{}

	locals.AwsS3ObjectSet = stackInput.Target

	target := locals.AwsS3ObjectSet

	// Resource-identity tags follow the catalog convention; user labels merge
	// in first so they can never override the identity keys. Every object in
	// the set carries them (S3 object tags), attributing each object back to
	// the owning set for auditing and orphan cleanup — the same key set the
	// Terraform module emits, keeping cross-engine tag parity.
	locals.Labels = map[string]string{}
	for k, v := range target.Metadata.Labels {
		locals.Labels[k] = v
	}
	locals.Labels["Name"] = target.Metadata.Name
	locals.Labels["planton.ai/resource"] = "true"
	locals.Labels["planton.ai/organization"] = target.Metadata.Org
	locals.Labels["planton.ai/environment"] = target.Metadata.Env
	locals.Labels["planton.ai/resource-kind"] = "AwsS3ObjectSet"
	locals.Labels["planton.ai/resource-id"] = target.Metadata.Id

	return locals
}
