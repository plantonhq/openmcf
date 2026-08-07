package main

import (
	"github.com/pkg/errors"
	alicloudsaeapplicationv1alpha1 "github.com/plantonhq/planton/catalog/alicloud/alicloudsaeapplication/v1alpha1"
	"github.com/plantonhq/planton/catalog/alicloud/alicloudsaeapplication/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &alicloudsaeapplicationv1alpha1.AliCloudSaeApplicationStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
