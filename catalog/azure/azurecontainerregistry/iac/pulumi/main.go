// Package main provides the Pulumi program entrypoint for Azure Container Registry deployment.
// Binary releases are gzip-compressed to reduce download size.
package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/azure/azurecontainerregistry/iac/pulumi/module"
	azurecontainerregistryv1alpha1 "github.com/plantonhq/planton/catalog/azure/azurecontainerregistry/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &azurecontainerregistryv1alpha1.AzureContainerRegistryStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
