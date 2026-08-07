package main

import (
	gcpfilestoreinstancev1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpfilestoreinstance/v1alpha1"
	"github.com/plantonhq/planton/catalog/gcp/gcpfilestoreinstance/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpfilestoreinstancev1alpha1.GcpFilestoreInstanceStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
