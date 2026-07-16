// Package main provides the Pulumi program entrypoint for AWS EKS Access Entry deployment.
package main

import (
	"github.com/pkg/errors"
	awseksaccessentryv1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awseksaccessentry/v1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/aws/awseksaccessentry/v1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awseksaccessentryv1.AwsEksAccessEntryStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
