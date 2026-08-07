package verify

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	pkgerrors "github.com/pkg/errors"
)

// managedRedisGeoReplicationVerifier verifies an
// AzureManagedRedisGeoReplication, keyed on the managing cluster's ARM
// ID (the group has no ARM object of its own -- membership lives on
// every member's default database).
//
// STATE-AWARE by necessity: the group's existence cannot be probed by a
// 404 -- creating a group does not create a resource, and destroying one
// unlinks the members without deleting anything. This verifier GETs the
// managing cluster's default database and reads
// properties.geoReplication.linkedDatabases: a group exists when the
// database reports MORE THAN ONE linked database (the database itself is
// always its own first member once a group name is set), and it is
// absent when the linked list has collapsed back to at most the
// database itself.
type managedRedisGeoReplicationVerifier struct{}

// IDOutputKey is the managing cluster's ARM ID.
func (*managedRedisGeoReplicationVerifier) IDOutputKey() string {
	return "geo_replication_id"
}

// geoReplicationLinkCount GETs the cluster's default database and
// reports (databaseExists, linkedDatabaseCount).
func geoReplicationLinkCount(ctx context.Context, cred azcore.TokenCredential, subscriptionID, clusterID string) (bool, int, error) {
	client, err := armresources.NewClient(subscriptionID, cred, nil)
	if err != nil {
		return false, 0, err
	}
	databaseID := clusterID + "/databases/default"
	resp, err := client.GetByID(ctx, databaseID, managedRedisAPIVersion, nil)
	if err != nil {
		var respErr *azcore.ResponseError
		if errors.As(err, &respErr) && respErr.StatusCode == 404 {
			return false, 0, nil
		}
		return false, 0, err
	}
	properties, ok := resp.Properties.(map[string]interface{})
	if !ok {
		return true, 0, pkgerrors.Errorf("managed redis database %q returned no readable properties", databaseID)
	}
	geoReplication, ok := properties["geoReplication"].(map[string]interface{})
	if !ok {
		return true, 0, nil
	}
	linkedDatabases, ok := geoReplication["linkedDatabases"].([]interface{})
	if !ok {
		return true, 0, nil
	}
	return true, len(linkedDatabases), nil
}

func (*managedRedisGeoReplicationVerifier) VerifyExists(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, linkCount, err := geoReplicationLinkCount(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremanagedredisgeoreplication verify-exists failed for %q", id)
	}
	if !exists {
		return pkgerrors.Errorf("azuremanagedredisgeoreplication: managing database for %q not found after deploy", id)
	}
	if linkCount < 2 {
		return pkgerrors.Errorf("azuremanagedredisgeoreplication: %q reports %d linked database(s) after deploy -- expected at least 2 (the group did not link)", id, linkCount)
	}
	return nil
}

func (*managedRedisGeoReplicationVerifier) VerifyAbsent(ctx context.Context, cred azcore.TokenCredential, subscriptionID, id string) error {
	exists, linkCount, err := geoReplicationLinkCount(ctx, cred, subscriptionID, id)
	if err != nil {
		return pkgerrors.Wrapf(err, "azuremanagedredisgeoreplication verify-absent failed for %q", id)
	}
	if exists && linkCount > 1 {
		return pkgerrors.Errorf("azuremanagedredisgeoreplication: %q still reports %d linked databases after destroy", id, linkCount)
	}
	return nil
}
