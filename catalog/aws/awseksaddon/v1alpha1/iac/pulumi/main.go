// Package main provides the Pulumi program entrypoint for AWS EKS Add-on deployment.
package main

import (
	"github.com/pkg/errors"
	awseksaddonv1alpha1 "github.com/plantonhq/planton/catalog/aws/awseksaddon/v1alpha1"
	"github.com/plantonhq/planton/catalog/aws/awseksaddon/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awseksaddonv1alpha1.AwsEksAddonStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
