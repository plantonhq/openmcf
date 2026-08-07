package main

import (
	"github.com/pkg/errors"
	azurekeyvaultkeyv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurekeyvaultkey/v1alpha1"
	"github.com/plantonhq/planton/catalog/azure/azurekeyvaultkey/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azurekeyvaultkeyv1alpha1.AzureKeyVaultKeyStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
