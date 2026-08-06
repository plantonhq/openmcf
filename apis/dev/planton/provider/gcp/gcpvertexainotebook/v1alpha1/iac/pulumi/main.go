package main

import (
	gcpvertexainotebookv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpvertexainotebook/v1alpha1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpvertexainotebook/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpvertexainotebookv1alpha1.GcpVertexAiNotebookStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
