package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/hetznercloud/hetznercloudcertificate/iac/pulumi/module"
	hetznercloudcertificatev1alpha1 "github.com/plantonhq/planton/catalog/hetznercloud/hetznercloudcertificate/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &hetznercloudcertificatev1alpha1.HetznerCloudCertificateStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
