package module

import (
	"strings"

	"github.com/pkg/errors"
)

// parseMongoDatabaseId extracts the database name, account name, and
// resource-group name from a Cosmos DB Mongo database ARM ID
// (/subscriptions/{sub}/resourceGroups/{rg}/providers/
// Microsoft.DocumentDB/databaseAccounts/{account}/mongodbDatabases/{db}).
// The ID must END with /mongodbDatabases/{db} and carry the account and
// resource-group segments -- the same anchored semantics as the
// Terraform module's regexes -- so a malformed ID fails loudly instead
// of computing wrong names.
func parseMongoDatabaseId(mongoDatabaseId string) (databaseName string, accountName string, resourceGroupName string, err error) {
	databaseParts := strings.Split(mongoDatabaseId, "/mongodbDatabases/")
	if len(databaseParts) != 2 || databaseParts[1] == "" || strings.Contains(databaseParts[1], "/") {
		return "", "", "", errors.Errorf("mongo_database_id %q is not a Cosmos DB Mongo database ARM id", mongoDatabaseId)
	}
	databaseName = databaseParts[1]

	accountParts := strings.Split(databaseParts[0], "/databaseAccounts/")
	if len(accountParts) != 2 || accountParts[1] == "" || strings.Contains(accountParts[1], "/") {
		return "", "", "", errors.Errorf("mongo_database_id %q carries no database-account segment", mongoDatabaseId)
	}
	accountName = accountParts[1]

	rgParts := strings.Split(mongoDatabaseId, "/resourceGroups/")
	if len(rgParts) != 2 || rgParts[1] == "" || !strings.Contains(rgParts[1], "/") {
		return "", "", "", errors.Errorf("mongo_database_id %q carries no resource-group segment", mongoDatabaseId)
	}
	resourceGroupName = rgParts[1][:strings.Index(rgParts[1], "/")]
	if resourceGroupName == "" {
		return "", "", "", errors.Errorf("mongo_database_id %q carries an empty resource-group name", mongoDatabaseId)
	}

	return databaseName, accountName, resourceGroupName, nil
}
