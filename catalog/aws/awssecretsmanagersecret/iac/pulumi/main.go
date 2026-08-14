package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/aws/awssecretsmanagersecret/iac/pulumi/module"
	awssecretsmanagersecretv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssecretsmanagersecret/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awssecretsmanagersecretv1alpha1.AwsSecretsManagerSecretStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}
		return module.Resources(ctx, stackInput)
	})
}
