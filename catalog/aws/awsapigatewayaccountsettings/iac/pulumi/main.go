package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/aws/awsapigatewayaccountsettings/iac/pulumi/module"
	awsapigatewayaccountsettingsv1alpha1 "github.com/plantonhq/planton/catalog/aws/awsapigatewayaccountsettings/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awsapigatewayaccountsettingsv1alpha1.AwsApiGatewayAccountSettingsStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}
		return module.Resources(ctx, stackInput)
	})
}
