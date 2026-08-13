# App Identity Access

This preset grants an application's managed identity access to a Mongo vCore cluster -- the passwordless wiring: the app authenticates to MongoDB with its Entra token, and nothing needs rotating.

## When to Use

- Wiring a workload (AKS pod, App Service, Function) to its database through its managed identity
- Replacing shared admin connection strings with per-app revocable access

## Key Configuration Choices

- **`servicePrincipal` as the principal type** -- managed identities bind through their service principal; `user` is for humans
- **The PRINCIPAL id, not the client id** -- referencing the identity's principal output avoids the classic copy-the-wrong-UUID failure
- **`root` on `admin`** -- cluster-wide access, Azure's only role today; the cluster must list `MicrosoftEntraID` in its `authenticationMethods` first

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<your-mongo-cluster>` | The Planton name of your `AzureMongoCluster` resource | Planton console (or replace `valueFrom` with `value:` and the cluster's ARM ID) |
| `<your-app-identity>` | The Planton name of the app's `AzureUserAssignedIdentity` | Planton console (or replace `valueFrom` with `value:` and the principal's object id) |
