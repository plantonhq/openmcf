package main

import (
	kubernetesgatewayapicrdsv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesgatewayapicrds/v1alpha1"
	"github.com/plantonhq/planton/catalog/kubernetes/kubernetesgatewayapicrds/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &kubernetesgatewayapicrdsv1alpha1.KubernetesGatewayApiCrdsStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
