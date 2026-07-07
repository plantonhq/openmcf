package module

import (
	"strings"

	"github.com/pkg/errors"
)

// parseCosmosdbAccountId extracts the account name and resource-group
// name from a Cosmos DB account ARM ID
// (/subscriptions/{sub}/resourceGroups/{rg}/providers/
// Microsoft.DocumentDB/databaseAccounts/{name}). The ID must END with
// /databaseAccounts/{name} and carry a /resourceGroups/{rg}/ segment --
// the same anchored semantics as the Terraform module's regexes -- so a
// malformed ID fails loudly instead of computing wrong names.
func parseCosmosdbAccountId(cosmosdbAccountId string) (accountName string, resourceGroupName string, err error) {
	accountParts := strings.Split(cosmosdbAccountId, "/databaseAccounts/")
	if len(accountParts) != 2 || accountParts[1] == "" || strings.Contains(accountParts[1], "/") {
		return "", "", errors.Errorf("cosmosdb_account_id %q is not a Cosmos DB account ARM id", cosmosdbAccountId)
	}
	accountName = accountParts[1]

	rgParts := strings.Split(cosmosdbAccountId, "/resourceGroups/")
	if len(rgParts) != 2 || rgParts[1] == "" || !strings.Contains(rgParts[1], "/") {
		return "", "", errors.Errorf("cosmosdb_account_id %q carries no resource-group segment", cosmosdbAccountId)
	}
	resourceGroupName = rgParts[1][:strings.Index(rgParts[1], "/")]
	if resourceGroupName == "" {
		return "", "", errors.Errorf("cosmosdb_account_id %q carries an empty resource-group name", cosmosdbAccountId)
	}

	return accountName, resourceGroupName, nil
}
