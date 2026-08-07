# Backend TLS Policy on Kubernetes

Encrypts the hop between your gateway and the Services behind it, and verifies
that the backend really is who it claims to be.

A route decides *where* traffic goes. This policy decides *how* the
gateway-to-backend leg is secured — the leg most setups leave in plaintext
because the client-facing side already has TLS.

## What it does

Attaches directly to one or more Services in its own namespace. When a gateway
forwards a request to a targeted Service it now:

1. Opens a TLS connection instead of a plaintext one.
2. Sends the configured hostname as the SNI so the backend can pick a
   certificate.
3. Verifies that certificate against your trust anchor — your own CA bundle, or
   the implementation's system store.
4. Optionally checks the certificate's Subject Alternative Names, for backends
   whose identity is a SPIFFE URI rather than a DNS name.

## Before you rely on it

**Your gateway controller must implement BackendTLSPolicy.** Support is still
uneven across implementations. A policy attached to a controller that does not
implement it is accepted by the API server and then *silently ignored* — the
gateway keeps sending plaintext, with nothing in the events to say so. Check the
policy's `Accepted` and `ResolvedRefs` conditions after the first deploy.

**The CRDs must exist first.** BackendTLSPolicy ships with the Gateway API
standard channel. Install the CRDs (KubernetesGatewayApiCrds) before the policy.

**Same namespace only.** Upstream forbids cross-namespace targets here, so the
policy, the Services it secures, and any CA-bundle ConfigMap all live together.

## Choosing a trust anchor

Exactly one of two, never both and never neither:

| Arm | When | Shape |
|---|---|---|
| CA bundle | The backend's certificate is issued by your own CA | One ConfigMap in this namespace with the PEM chain under the key `ca.crt` |
| System store | The backend serves a publicly-issued certificate | `wellKnownCACertificates: System` |

A cert-manager CA chain composes directly into the first arm: the root
Certificate's ConfigMap — or a trust-manager Bundle target — is exactly the
referent this expects.

## Hostname does double duty

The hostname is sent as the SNI *and*, unless Subject Alternative Names are
listed, it is the identity the certificate must match. Once you list SANs the
hostname only selects the certificate; if it should still be accepted as an
identity, add it as a `Hostname` SAN too. Getting this wrong produces a
handshake failure that names neither the field nor the policy.

## Works with

| Kind | Relationship |
|---|---|
| KubernetesService | The policy's targets — referenced by name, optionally narrowed to one port |
| KubernetesConfigMap | Carries the CA bundle when you bring your own trust anchor |
| KubernetesCertificate | cert-manager issues the backend certificates this policy verifies |
| KubernetesGatewayApiCrds | Installs the CRD this policy is served by |
| KubernetesGateway | The gateway whose controller enforces the policy |

## Outputs

| Output | Description |
|---|---|
| `policy_name` | The BackendTLSPolicy object's name in the cluster |
| `namespace` | The namespace it was created in |

Whether the policy actually attached is controller state rather than a stack
output — it is reconciled asynchronously, so read the `Accepted` and
`ResolvedRefs` conditions with kubectl rather than expecting them here.
