package main

import (
	"github.com/plantonhq/planton/catalog/gcp/gcpbigquerydataset/iac/pulumi/module"
	gcpbigquerydatasetv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpbigquerydataset/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpbigquerydatasetv1alpha1.GcpBigQueryDatasetStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
