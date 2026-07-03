package verify

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/authorization/armauthorization/v2"
	pkgerrors "github.com/pkg/errors"
)

// roleDefinitionVerifier verifies an AzureRoleDefinition via
// RoleDefinitionsClient.GetByID, keyed on the definition's fully-scoped ARM ID
// ({scope}/providers/Microsoft.Authorization/roleDefinitions/{guid}) -- the one
// identifier that pins both the scope and the definition GUID, so the probe
// cannot match a same-named role at another scope. A typed 404 ResponseError is
// the absence signal; any other failure (auth, network) surfaces as a real error
// rather than masquerading as "absent".
type roleDefinitionVerifier struct{}

// IDOutputKey is the fully-scoped ARM ID of the definition: GetByID resolves it
// directly, independent of which scope the definition was created at.
func (*roleDefinitionVerifier) IDOutputKey() string { return "role_definition_id" }

func (*roleDefinitionVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := roleDefinitionExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureroledefinition verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureroledefinition %q not found after deploy", id)
	}
	return nil
}

func (*roleDefinitionVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := roleDefinitionExists(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureroledefinition verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureroledefinition %q still exists after destroy", id)
	}
	return nil
}

func roleDefinitionExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, roleDefinitionID string) (bool, error) {
	client, err := armauthorization.NewRoleDefinitionsClient(cred, nil)
	if err != nil {
		return false, err
	}
	if _, err := client.GetByID(ctx, roleDefinitionID, nil); err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
