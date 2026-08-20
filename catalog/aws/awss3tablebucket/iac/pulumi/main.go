package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/aws/awss3tablebucket/iac/pulumi/module"
	awss3tablebucketv1alpha1 "github.com/plantonhq/planton/catalog/aws/awss3tablebucket/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awss3tablebucketv1alpha1.AwsS3TableBucketStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}
		return module.Resources(ctx, stackInput)
	})
}
