package main

import (
	"github.com/pkg/errors"
	testcloudresourcegenericv1alpha2 "github.com/plantonhq/planton/catalog/_test/testcloudresourcegeneric/v1alpha2"
	"github.com/plantonhq/planton/catalog/_test/testcloudresourcegeneric/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &testcloudresourcegenericv1alpha2.TestCloudResourceGenericStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
