package module

import (
	"github.com/pkg/errors"
	gcpgkeclusterv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/gcp/gcpgkecluster/v1"
	"github.com/plantonhq/openmcf/pkg/iac/pulumi/pulumimodule/provider/gcp/pulumigoogleprovider"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Resources is the Pulumi program entry‑point invoked by the OpenMCF
func Resources(
	ctx *pulumi.Context,
	stackInput *gcpgkeclusterv1.GcpGkeClusterStackInput,
) error {
	// gather locals (Terraform‑style “locals”)
	locals := initializeLocals(stackInput)

	// configure a GCP provider from the given credential
	gcpProvider, err := pulumigoogleprovider.Get(ctx, stackInput.ProviderConfig)
	if err != nil {
		return errors.Wrap(err, "failed to setup google provider")
	}

	// Cluster.
	_, err = cluster(ctx, locals, gcpProvider)
	if err != nil {
		return errors.Wrap(err, "cluster creation failed")
	}

	return nil
}
