package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// containerAppJobVerifier verifies an AzureContainerAppJob via the generic
// ARM resources GetByID (see armResourceExists), keyed on the job's full
// ARM ID.
type containerAppJobVerifier struct{}

// IDOutputKey is the job's full ARM ID.
func (*containerAppJobVerifier) IDOutputKey() string {
	return "job_id"
}

func (*containerAppJobVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappjob verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurecontainerappjob %q not found after deploy", id)
	}
	return nil
}

func (*containerAppJobVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, containerAppsAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurecontainerappjob verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurecontainerappjob %q still exists after destroy", id)
	}
	return nil
}
