package main

import (
	gcpdataprocautoscalingpolicyv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpdataprocautoscalingpolicy/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpdataprocautoscalingpolicy/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpdataprocautoscalingpolicyv1alpha1.GcpDataprocAutoscalingPolicyStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
