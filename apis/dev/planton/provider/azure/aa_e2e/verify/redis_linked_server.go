package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// redisLinkedServerVerifier verifies an AzureRedisLinkedServer via the
// generic ARM resources GetByID, keyed on the link's full ARM ID
// (.../redis/{primary}/linkedServers/{secondary}).
type redisLinkedServerVerifier struct{}

// IDOutputKey is the linked server's full ARM ID.
func (*redisLinkedServerVerifier) IDOutputKey() string {
	return "linked_server_id"
}

func (*redisLinkedServerVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, redisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureredislinkedserver verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azureredislinkedserver %q not found after deploy", id)
	}
	return nil
}

func (*redisLinkedServerVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, redisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azureredislinkedserver verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azureredislinkedserver %q still exists after destroy", id)
	}
	return nil
}
