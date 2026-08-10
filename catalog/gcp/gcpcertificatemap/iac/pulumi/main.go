package main

import (
	"github.com/plantonhq/planton/catalog/gcp/gcpcertificatemap/iac/pulumi/module"
	gcpcertificatemapv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpcertificatemap/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpcertificatemapv1alpha1.GcpCertificateMapStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
