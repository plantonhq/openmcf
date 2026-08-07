// Package main provides the Pulumi program entrypoint for the AWS Planton
// Runner appliance: a standing, outbound-only runner on ECS Fargate that
// executes deploy and cloud operations from inside the VPC.
package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/aws/awsplantonrunner/iac/pulumi/module"
	awsplantonrunnerv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsplantonrunner/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awsplantonrunnerv1alpha1.AwsPlantonRunnerStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
