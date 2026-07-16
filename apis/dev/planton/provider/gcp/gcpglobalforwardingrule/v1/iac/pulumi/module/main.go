package module

import (
	"github.com/pkg/errors"
	gcpglobalforwardingrulev1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpglobalforwardingrule/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func Resources(ctx *pulumi.Context, stackInput *gcpglobalforwardingrulev1.GcpGlobalForwardingRuleStackInput) error {
	locals := initializeLocals(ctx, stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if err := globalForwardingRule(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create global forwarding rule")
	}

	return nil
}
