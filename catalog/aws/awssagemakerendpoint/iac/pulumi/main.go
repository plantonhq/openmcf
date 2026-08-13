package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/aws/awssagemakerendpoint/iac/pulumi/module"
	awssagemakerendpointv1alpha1 "github.com/plantonhq/planton/catalog/aws/awssagemakerendpoint/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awssagemakerendpointv1alpha1.AwsSagemakerEndpointStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}
		return module.Resources(ctx, stackInput)
	})
}
