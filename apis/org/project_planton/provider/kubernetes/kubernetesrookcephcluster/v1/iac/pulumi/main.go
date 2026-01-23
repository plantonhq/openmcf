package main

import (
	kubernetesrookcephclusterv1 "github.com/plantonhq/project-planton/apis/org/project_planton/provider/kubernetes/kubernetesrookcephcluster/v1"
	"github.com/plantonhq/project-planton/apis/org/project_planton/provider/kubernetes/kubernetesrookcephcluster/v1/iac/pulumi/module"
	"github.com/plantonhq/project-planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &kubernetesrookcephclusterv1.KubernetesRookCephClusterStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}

		return module.Resources(ctx, stackInput)
	})
}
