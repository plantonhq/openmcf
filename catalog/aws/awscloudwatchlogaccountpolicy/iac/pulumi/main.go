package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/aws/awscloudwatchlogaccountpolicy/iac/pulumi/module"
	awscloudwatchlogaccountpolicyv1alpha1 "github.com/plantonhq/planton/catalog/aws/awscloudwatchlogaccountpolicy/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awscloudwatchlogaccountpolicyv1alpha1.AwsCloudwatchLogAccountPolicyStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}
		return module.Resources(ctx, stackInput)
	})
}
