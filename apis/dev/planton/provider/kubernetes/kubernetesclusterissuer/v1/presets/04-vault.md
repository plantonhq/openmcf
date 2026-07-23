# Vault PKI ClusterIssuer

This preset signs certificates through a HashiCorp Vault (or OpenBao) PKI
secrets engine using the Kubernetes auth method — the keyless path: Vault
validates a ServiceAccount token instead of holding a static credential.

## When to Use

- Your organization centralizes PKI in Vault/OpenBao and cluster
  certificates must chain to it
- You want issuance audited and policy-controlled outside the cluster
- You do not want static Vault tokens or AppRole secrets on the cluster

## How the Keyless Auth Works

1. Vault's Kubernetes auth method is configured to trust the cluster
   (its OIDC issuer / token reviewer)
2. A Vault role binds the referenced ServiceAccount to a Vault policy
   allowing `<pki-mount>/sign/<role-name>`
3. cert-manager presents the ServiceAccount token; Vault signs

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<vault-host>` | Vault server address | Your Vault deployment |
| `<pki-mount>/sign/<role-name>` | PKI signing path | `vault secrets list` + your PKI role |
| `<vault-role-bound-to-cert-manager>` | Vault Kubernetes-auth role | `vault list auth/kubernetes/role` |
| `<service-account-name>` | ServiceAccount (in the cert-manager namespace) Vault trusts | Compose with a KubernetesServiceAccount resource |

## Alternative Auth Methods

Static token (`tokenAuth`) and AppRole (`appRoleAuth`) are also modeled —
use them only where the Kubernetes auth method is unavailable.
