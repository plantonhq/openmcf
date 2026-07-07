package module

import (
	"strings"

	"github.com/pkg/errors"
)

// parseSqlDatabaseId extracts the database name, account name, and
// resource-group name from a Cosmos DB SQL database ARM ID
// (/subscriptions/{sub}/resourceGroups/{rg}/providers/
// Microsoft.DocumentDB/databaseAccounts/{account}/sqlDatabases/{db}).
// The ID must END with /sqlDatabases/{db} and carry the account and
// resource-group segments -- the same anchored semantics as the
// Terraform module's regexes -- so a malformed ID fails loudly instead
// of computing wrong names.
func parseSqlDatabaseId(sqlDatabaseId string) (databaseName string, accountName string, resourceGroupName string, err error) {
	databaseParts := strings.Split(sqlDatabaseId, "/sqlDatabases/")
	if len(databaseParts) != 2 || databaseParts[1] == "" || strings.Contains(databaseParts[1], "/") {
		return "", "", "", errors.Errorf("sql_database_id %q is not a Cosmos DB SQL database ARM id", sqlDatabaseId)
	}
	databaseName = databaseParts[1]

	accountParts := strings.Split(databaseParts[0], "/databaseAccounts/")
	if len(accountParts) != 2 || accountParts[1] == "" || strings.Contains(accountParts[1], "/") {
		return "", "", "", errors.Errorf("sql_database_id %q carries no database-account segment", sqlDatabaseId)
	}
	accountName = accountParts[1]

	rgParts := strings.Split(sqlDatabaseId, "/resourceGroups/")
	if len(rgParts) != 2 || rgParts[1] == "" || !strings.Contains(rgParts[1], "/") {
		return "", "", "", errors.Errorf("sql_database_id %q carries no resource-group segment", sqlDatabaseId)
	}
	resourceGroupName = rgParts[1][:strings.Index(rgParts[1], "/")]
	if resourceGroupName == "" {
		return "", "", "", errors.Errorf("sql_database_id %q carries an empty resource-group name", sqlDatabaseId)
	}

	return databaseName, accountName, resourceGroupName, nil
}
