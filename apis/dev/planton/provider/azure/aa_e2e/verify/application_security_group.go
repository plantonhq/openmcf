package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// applicationSecurityGroupAPIVersion is the stable Microsoft.Network API
// version the generic existence probe is pinned to.
const applicationSecurityGroupAPIVersion = "2024-05-01"

// applicationSecurityGroupVerifier verifies an AzureApplicationSecurityGroup
// via the generic ARM resources GetByID (see armResourceExists), keyed on the
// group's full ARM ID.
type applicationSecurityGroupVerifier struct{}

// IDOutputKey is the ASG's full ARM ID.
func (*applicationSecurityGroupVerifier) IDOutputKey() string {
	return "application_security_group_id"
}

func (*applicationSecurityGroupVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, applicationSecurityGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureapplicationsecuritygroup verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureapplicationsecuritygroup %q not found after deploy", id)
	}
	return nil
}

func (*applicationSecurityGroupVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, applicationSecurityGroupAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureapplicationsecuritygroup verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureapplicationsecuritygroup %q still exists after destroy", id)
	}
	return nil
}
