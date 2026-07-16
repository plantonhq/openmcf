package verify

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	pkgerrors "github.com/pkg/errors"
)

// managedRedisAPIVersion pins the Microsoft.Cache redisEnterprise
// resource-provider API version for the whole Managed Redis family
// (cluster, geo-replication, access policy assignments) -- the same
// version line the IaC providers target. Distinct from the classic
// Redis family's pin: redisEnterprise is a different ARM type with its
// own API version history.
const managedRedisAPIVersion = "2025-07-01"

// managedRedisVerifier verifies an AzureManagedRedis via the generic ARM
// resources GetByID, keyed on the cluster's full ARM ID.
type managedRedisVerifier struct{}

// IDOutputKey is the cluster's full ARM ID.
func (*managedRedisVerifier) IDOutputKey() string {
	return "managed_redis_id"
}

func (*managedRedisVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, managedRedisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremanagedredis verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremanagedredis %q not found after deploy", id)
	}
	return nil
}

func (*managedRedisVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, err := armResourceExists(ctx, cred, subscriptionID, id, managedRedisAPIVersion)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremanagedredis verify-absent failed for %q", id)
	}
	if exists {
		return pkgerrors.Errorf("azuremanagedredis %q still exists after destroy", id)
	}
	return nil
}
