# Governance-Tagged, Regionally Isolated Identity

This preset creates an identity shaped for regulated environments: regional
isolation restricts token issuance to the identity's own region (a
data-residency / blast-radius control), and governance tags carry the
ownership and cost-attribution conventions Azure Policy can enforce and
Microsoft Cost Management groups by.

User tags merge over the Planton-derived resource tags (organization,
environment, resource id), and a user tag with the same key wins -- so your
org's conventions always take precedence.

## When to Use

- Environments with data-residency requirements where an identity must not
  be usable outside its region
- Organizations enforcing tag compliance (cost center, owner, data
  classification) through Azure Policy
- Any identity that should be attributable in cost and ownership reports

## Key Configuration Choices

- **`isolationScope: REGIONAL`** -- opt-in; most deployments omit it (the
  default identity is usable from any region). It updates in place, so you
  can adopt it later without replacing the identity
- **Tags as governance, not documentation** -- align keys with what your
  Azure Policy assignments actually check

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<azure-region>` | Azure region | Your regional deployment strategy |
| `<your-resource-group-name>` | Name of the resource group | Azure portal or `AzureResourceGroup` status outputs |
| `<your-identity-name>` | Name for the managed identity (3-128 chars) | Your naming convention |
| `<cost-center>` / `<owning-team>` | Your governance tag values | Your tagging policy |
