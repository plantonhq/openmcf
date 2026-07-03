package module

import (
	"github.com/pkg/errors"
	gcpcloudsqlv1 "github.com/plantonhq/planton/apis/dev/planton/provider/gcp/gcpcloudsql/v1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the Pulumi program entry-point for the GcpCloudSql component.
func Resources(ctx *pulumi.Context, stackInput *gcpcloudsqlv1.GcpCloudSqlStackInput) error {
	locals := initializeLocals(stackInput)

	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	if _, err := databaseInstance(ctx, locals, gcpProvider); err != nil {
		return errors.Wrap(err, "failed to create cloud sql instance")
	}

	return nil
}
