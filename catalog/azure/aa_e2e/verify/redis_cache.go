package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// redisAPIVersion pins the Microsoft.Cache resource-provider API version
// for the whole Redis family (cache, linked server, access policies and
// their assignments) -- the same version line the IaC providers target.
const redisAPIVersion = "2024-11-01"

// redisCacheVerifier verifies an AzureRedisCache via the generic ARM
// resources GetByID, keyed on the cache's full ARM ID.
type redisCacheVerifier struct{}

// IDOutputKey is the cache's full ARM ID.
func (*redisCacheVerifier) IDOutputKey() string {
	return "redis_cache_id"
}

func (*redisCacheVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, redisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurerediscache verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azurerediscache %q not found after deploy", id)
	}
	return nil
}

func (*redisCacheVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, redisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azurerediscache verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azurerediscache %q still exists after destroy", id)
	}
	return nil
}
