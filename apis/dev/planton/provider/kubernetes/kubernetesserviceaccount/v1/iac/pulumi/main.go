// Pulumi entrypoint for the KubernetesServiceAccount component.
// Loads the stack input and delegates all resource creation to the module package.
package main

import (
	"github.com/pkg/errors"
	kubernetesserviceaccountv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesserviceaccount/v1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesserviceaccount/v1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &kubernetesserviceaccountv1.KubernetesServiceAccountStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
