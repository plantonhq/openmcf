// Package main provides the Pulumi program entrypoint for AWS EKS Node Group deployment.
// Auto-release test: Multi-provider Pulumi change (AWS component).
package main

import (
	"github.com/pkg/errors"
	awseksnodegroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awseksnodegroup/v1alpha1"
	"github.com/plantonhq/planton/catalog/aws/awseksnodegroup/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awseksnodegroupv1alpha1.AwsEksNodeGroupStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
