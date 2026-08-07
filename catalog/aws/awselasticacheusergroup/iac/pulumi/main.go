// Package main provides the Pulumi program entrypoint for AWS ElastiCache user group deployment.
package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/aws/awselasticacheusergroup/iac/pulumi/module"
	awselasticacheusergroupv1alpha1 "github.com/plantonhq/planton/catalog/aws/awselasticacheusergroup/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awselasticacheusergroupv1alpha1.AwsElasticacheUserGroupStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
