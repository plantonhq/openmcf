package verify

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	pkgerrors "github.com/pkg/errors"
)

// roleAssignmentVerifier verifies an AzureRoleAssignment via
// RoleAssignmentsClient.GetByID, keyed on the assignment's fully-scoped ARM ID
// ({scope}/providers/Microsoft.Authorization/roleAssignments/{guid}) -- the one
// identifier that pins both the scope and the assignment GUID, so the probe
// cannot match a different grant at another scope. A typed 404 ResponseError is
// the absence signal; any other failure (auth, network) surfaces as a real error
// rather than masquerading as "absent".
type roleAssignmentVerifier struct{}

// IDOutputKey is the fully-scoped ARM ID of the assignment: GetByID resolves it
// directly, independent of which scope the role was granted at.
func (*roleAssignmentVerifier) IDOutputKey() string { return "role_assignment_id" }

func (*roleAssignmentVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := roleAssignmentExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureroleassignment verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureroleassignment %q not found after deploy", id)
	}
	return nil
}

func (*roleAssignmentVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := roleAssignmentExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureroleassignment verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureroleassignment %q still exists after destroy", id)
	}
	return nil
}

func roleAssignmentExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, roleAssignmentID string) (bool, error) {
	client, err := armauthorization.NewRoleAssignmentsClient(subscriptionID, cred, nil)
	if err != nil {
		return false, err
	}
	if _, err := client.GetByID(ctx, roleAssignmentID, nil); err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
