// Package main provides the Pulumi program entrypoint for AWS ECS Task Definition deployment.
package main

import (
	"github.com/pkg/errors"
	awsecstaskdefinitionv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsecstaskdefinition/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsecstaskdefinition/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awsecstaskdefinitionv1alpha1.AwsEcsTaskDefinitionStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
