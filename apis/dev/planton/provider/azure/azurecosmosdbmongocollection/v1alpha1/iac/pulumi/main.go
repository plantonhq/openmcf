package main

import (
	"github.com/pkg/errors"
	azurecosmosdbmongocollectionv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurecosmosdbmongocollection/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/azure/azurecosmosdbmongocollection/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azurecosmosdbmongocollectionv1alpha1.AzureCosmosdbMongoCollectionStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
