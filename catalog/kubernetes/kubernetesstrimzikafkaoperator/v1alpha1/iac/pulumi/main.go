package main

import (
	"github.com/pkg/errors"
	kubernetesstrimzikafkaoperatorv1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesstrimzikafkaoperator/v1alpha1"
	"github.com/plantonhq/planton/catalog/kubernetes/kubernetesstrimzikafkaoperator/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &kubernetesstrimzikafkaoperatorv1alpha1.KubernetesStrimziKafkaOperatorStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
