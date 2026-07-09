package main

import (
	gcpvertexaiindexendpointv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpvertexaiindexendpoint/v1"
	"github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpvertexaiindexendpoint/v1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &gcpvertexaiindexendpointv1.GcpVertexAiIndexEndpointStackInput{}
		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return err
		}
		return module.Resources(ctx, stackInput)
	})
}
