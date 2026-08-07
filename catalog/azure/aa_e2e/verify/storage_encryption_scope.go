package verify

import (
	"context"
	"errors"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	pkgerrors "github.com/pkg/errors"
)

// storageEncryptionScopeVerifier verifies an AzureStorageEncryptionScope,
// keyed on the scope's full ARM ID
// (.../storageAccounts/{account}/encryptionScopes/{name}).
//
// STATE-AWARE by necessity: ARM has no true delete for encryption scopes
// -- destroy flips the scope's state to Disabled and the object remains
// GETtable, so a plain 404 probe would report a destroyed scope as still
// existing. This verifier reads the scope's properties.state and treats
// Disabled as absent, mirroring the providers' own read semantics.
type storageEncryptionScopeVerifier struct{}

// IDOutputKey is the encryption scope's full ARM ID.
func (*storageEncryptionScopeVerifier) IDOutputKey() string {
	return "encryption_scope_id"
}

// encryptionScopeEnabled GETs the scope and reports (exists, enabled).
// A typed 404 is genuine absence; an existing scope is "enabled" only
// when properties.state says so.
func encryptionScopeEnabled(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) (bool, bool, error) {
	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return false, false, err
	}
	resp, err := client.GetByID(ctx, id, storageAPIVersion, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			return false, false, nil
		}
		return false, false, err
	}
	properties, ok := resp.Properties.(map[string]interface{})
	if !ok {
		return true, false, pkgerrors.Errorf("encryption scope %q returned no readable properties", id)
	}
	state, _ := properties["state"].(string)
	return true, strings.EqualFold(state, "Enabled"), nil
}

func (*storageEncryptionScopeVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, enabled, err := encryptionScopeEnabled(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestorageencryptionscope verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurestorageencryptionscope %q not found after deploy", id)
	}
	if !enabled {
		return pkgerrors.Errorf("azurestorageencryptionscope %q exists but is not Enabled after deploy", id)
	}
	return nil
}

func (*storageEncryptionScopeVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, enabled, err := encryptionScopeEnabled(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurestorageencryptionscope verify-absent failed for %q", id)
	}
	if exists && enabled {
		return pkgerrors.Errorf("azurestorageencryptionscope %q is still Enabled after destroy", id)
	}
	return nil
}
