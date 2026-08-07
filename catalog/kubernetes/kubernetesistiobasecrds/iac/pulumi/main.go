package main

import (
	"github.com/plantonhq/planton/catalog/kubernetes/kubernetesistiobasecrds/iac/pulumi/module"
	kubernetesistiobasecrdsv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesistiobasecrds/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &kubernetesistiobasecrdsv1alpha1.KubernetesIstioBaseCrdsStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
