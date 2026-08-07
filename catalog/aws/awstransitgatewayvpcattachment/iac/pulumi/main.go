// Package main provides the Pulumi program entrypoint for AWS Transit
// Gateway VPC attachment deployment.
package main

import (
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/catalog/aws/awstransitgatewayvpcattachment/iac/pulumi/module"
	awstgwattachv1 "github.com/plantonhq/planton/catalog/aws/awstransitgatewayvpcattachment/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awstgwattachv1.AwsTransitGatewayVpcAttachmentStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
