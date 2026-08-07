// Package main provides the Pulumi program entrypoint for AWS Transit
// Gateway route table deployment.
package main

import (
	"github.com/pkg/errors"
	awstgwrtv1 "github.com/plantonhq/planton/catalog/aws/awstransitgatewayroutetable/v1alpha1"
	"github.com/plantonhq/planton/catalog/aws/awstransitgatewayroutetable/v1alpha1/iac/pulumi/module"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/stackinput"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		stackInput := &awstgwrtv1.AwsTransitGatewayRouteTableStackInput{}

		if err := stackinput.LoadStackInput(ctx, stackInput); err != nil {
			return errors.Wrap(err, "failed to load stack-input")
		}

		return module.Resources(ctx, stackInput)
	})
}
