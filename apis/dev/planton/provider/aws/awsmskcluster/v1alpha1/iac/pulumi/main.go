package main

import (
	"github.com/pkg/errors"
	awsmskclusterv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsmskcluster/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/aws/awsmskcluster/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awsmskclusterv1alpha1.AwsMskClusterStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}
		return module.Resources(ctx, stackInput)
	})
}
