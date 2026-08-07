# KubernetesBackendTlsPolicy

> Generated from the protobuf schema by `make generate-reference` -- do not
> edit by hand. To change a fact on this page, change the proto field comment
> or validation rule it is derived from, then regenerate.

**apiVersion**: `kubernetes.planton.dev/v1alpha1`

KubernetesBackendTlsPolicySpec defines a Kubernetes Gateway API
BackendTLSPolicy: the standard way to tell a Gateway implementation to
originate TLS to the backends BEHIND the gateway and verify the
certificate they present. Routes decide WHERE traffic goes; this policy
decides HOW the gateway-to-backend hop is secured. It is a Direct policy
attachment: it targets Services (Core support), optionally narrowed to a
single named Service port through section_name.

The trust anchor comes from exactly one of two arms: ca_certificate_refs
(bring your own CA bundle — the Core shape is one ConfigMap carrying the
PEM bundle in a key named `ca.crt`, which is exactly what a cert-manager
CA chain materializes) or well_known_ca_certificates: "System" (trust the
implementation's system store — for backends serving publicly-issued
certificates). The hostname is mandatory and is used as the SNI for the
backend connection and, unless subject_alt_names is set, as the identity
the backend certificate must prove.

IMPORTANT: the policy only takes effect when the Gateway's controller
IMPLEMENTS BackendTLSPolicy — support is still uneven across
implementations (verify your gateway controller's documentation). A
policy targeting a Service behind a non-implementing gateway is accepted
by the API server and then silently ignored: the gateway keeps sending
plaintext (or its own default TLS posture) to the backend.

100% fidelity with the upstream Gateway API v1.6.1 BackendTLSPolicySpec
(kubernetes-sigs/gateway-api apis/v1/backendtlspolicy_types.go), standard
channel. BackendTLSPolicy is served as gateway.networking.k8s.io/v1 (the
v1alpha3 version is deprecated upstream and no longer served).

## Example

```yaml
apiVersion: kubernetes.planton.dev/v1alpha1
kind: KubernetesBackendTlsPolicy
metadata:
  name: test-backend-tls-policy
spec:
  namespace:
    value: test-namespace
  targetRefs:
    - group: ""
      kind: Service
      name:
        value: test-backend-service
      sectionName: https
    - group: ""
      kind: Service
      name:
        value: test-backend-service
      sectionName: grpc
  validation:
    caCertificateRefs:
      - group: ""
        kind: ConfigMap
        name:
          value: test-ca-bundle
    hostname: test-backend-service.test-namespace.svc.cluster.local
    subjectAltNames:
      - type: Hostname
        hostname: test-backend-service.test-namespace.svc.cluster.local
      - type: URI
        uri: spiffe://cluster.example.com/ns/test-namespace/sa/test-backend-sa
  options:
    gateway.example.com/min-tls-version: TLSv1_3
```

## Spec Fields

| Path | Type | Required | Default | References |
|---|---|---|---|---|
| `spec.namespace` | `string \| valueFrom` | yes |  | KubernetesNamespace (`spec.name`) |
| `spec.targetRefs` | `[]KubernetesBackendTlsPolicyTargetReference` | yes |  |  |
| `spec.targetRefs[].group` | `string` | yes |  |  |
| `spec.targetRefs[].kind` | `string` | yes |  |  |
| `spec.targetRefs[].name` | `string \| valueFrom` | yes |  | KubernetesService (`status.outputs.service_name`) |
| `spec.targetRefs[].sectionName` | `string` | yes |  |  |
| `spec.validation` | `KubernetesBackendTlsPolicyValidation` | yes |  |  |
| `spec.validation.caCertificateRefs` | `[]KubernetesBackendTlsPolicyCaCertificateReference` |  |  |  |
| `spec.validation.caCertificateRefs[].group` | `string` | yes |  |  |
| `spec.validation.caCertificateRefs[].kind` | `string` | yes |  |  |
| `spec.validation.caCertificateRefs[].name` | `string \| valueFrom` | yes |  | KubernetesConfigMap (`status.outputs.configmap_name`) |
| `spec.validation.wellKnownCACertificates` | `string` | yes |  |  |
| `spec.validation.hostname` | `string` | yes |  |  |
| `spec.validation.subjectAltNames` | `[]KubernetesBackendTlsPolicySubjectAltName` |  |  |  |
| `spec.validation.subjectAltNames[].type` | `string` | yes |  |  |
| `spec.validation.subjectAltNames[].hostname` | `string` | yes |  |  |
| `spec.validation.subjectAltNames[].uri` | `string` | yes |  |  |
| `spec.options` | `map<string, string>` |  |  |  |

## Field Details

### spec.namespace

`string | valueFrom` · required

Namespace in which the BackendTLSPolicy is created. A BackendTLSPolicy
can only target Services in ITS OWN namespace (upstream forbids
cross-namespace targetRefs for this policy), so create the policy in the
namespace of the backend Services it secures.

- references: KubernetesNamespace (`spec.name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesNamespace, name: <that resource's name>, fieldPath: spec.name}} -- a bare string does not parse

### spec.targetRefs

`[]KubernetesBackendTlsPolicyTargetReference` · required

The backend Services this policy applies to, in the policy's own
namespace. Core support targets a Service; section_name narrows a
reference to one named Service port (omit it to cover every port).
Upstream note: implementations SHOULD support a single targetRef —
multiple entries are accepted by the API but the safest portable shape
is one.

The two list rules mirror the CRD's own CEL exactly: refs to the same
target must all carry a section_name (or none of them), and no two refs
to the same target may carry the same section_name.

- rule: section_name must be specified on every reference (or on none) when target_refs includes 2 or more references to the same target
- rule: section_name must be unique when target_refs includes 2 or more references to the same target
- rule: {"repeated":{"minItems":"1","maxItems":"16"}}

### spec.targetRefs[].group

`string` · required · optional (explicit presence)

Group of the target resource. Services live in the core API group —
the empty string.

Upstream requires the KEY to be present, but its Group type explicitly
allows the empty value -- so this is a proto3 `optional` string with a
presence `required` rule: it must be SET (and is therefore emitted to
the rendered CR, whose CRD rejects a missing key) but may be empty.
The `optional` is what keeps the projection faithful: protojson omits
unset proto3 scalars, so a non-optional empty-string group would be
dropped from the manifest and rejected by the API server.

Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"required":true,"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.targetRefs[].kind

`string` · required

Kind of the target resource. Core support: "Service".

Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.targetRefs[].name

`string | valueFrom` · required

Name of the target resource. Defaults to a KubernetesService foreign
key — wire it with valueFrom in an infra chart and the policy deploys
after the Service it secures. Pass the literal name with `value:` when
the target is not Planton-managed.

- references: KubernetesService (`status.outputs.service_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesService, name: <that resource's name>, fieldPath: status.outputs.service_name}} -- a bare string does not parse

### spec.targetRefs[].sectionName

`string` · required · optional (explicit presence)

Name of a section within the target resource. When unspecified, the
policy targets the entire resource. For a Service target this is a
PORT NAME — the policy then applies only to connections to that port.
A section_name that does not exist on the target makes the policy fail
to attach (surfaced through the policy's ResolvedRefs condition).

Upstream SectionName constraints: 1-253 chars, lowercase RFC 1123
subdomain.

- rule: {"string":{"minLen":"1","maxLen":"253","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.validation

`KubernetesBackendTlsPolicyValidation` · required

How the gateway validates the TLS handshake with the backend: the trust
anchor (CA bundle XOR system store), the SNI/authentication hostname,
and optional Subject Alternative Names.

- rule: {"required":true}
- rule: ca_certificate_refs and well_known_ca_certificates are mutually exclusive — bring your own CA bundle OR trust the system store, not both
- rule: either ca_certificate_refs (a CA-bundle ConfigMap) or well_known_ca_certificates ("System") must be specified — without a trust anchor the gateway cannot verify the backend

### spec.validation.caCertificateRefs

`[]KubernetesBackendTlsPolicyCaCertificateReference`

References to same-namespace objects carrying the PEM-encoded CA bundle
that signs the backend's serving certificate. Core support: ONE
ConfigMap with the bundle in a key named `ca.crt` (more refs, other
kinds, or multi-certificate bundles are implementation-specific). A
cert-manager CA chain composes here: the root Certificate's ConfigMap
(or a trust-manager Bundle target) is the natural referent.

- rule: {"repeated":{"maxItems":"8"}}

### spec.validation.caCertificateRefs[].group

`string` · required · optional (explicit presence)

Group of the referent. ConfigMaps live in the core API group — the
empty string.

Upstream requires the KEY to be present, but its Group type explicitly
allows the empty value -- so this is a proto3 `optional` string with a
presence `required` rule: it must be SET (and is therefore emitted to
the rendered CR, whose CRD rejects a missing key) but may be empty.

Group pattern: empty or an RFC 1123 subdomain (max 253).

- rule: {"required":true,"string":{"maxLen":"253","pattern":"^$|^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.validation.caCertificateRefs[].kind

`string` · required

Kind of the referent. Core support: "ConfigMap" (the CA bundle in a key
named `ca.crt`); other kinds are implementation-specific.

Kind pattern: 1-63 chars, ^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$

- rule: {"required":true,"string":{"minLen":"1","maxLen":"63","pattern":"^[a-zA-Z]([-a-zA-Z0-9]*[a-zA-Z0-9])?$"}}

### spec.validation.caCertificateRefs[].name

`string | valueFrom` · required

Name of the referent, in the policy's own namespace (cross-namespace CA
references are upstream-invalid for this policy). Defaults to a
KubernetesConfigMap foreign key — wire it with valueFrom against the
ConfigMap carrying the CA bundle.

- references: KubernetesConfigMap (`status.outputs.configmap_name`)
- rule: {"required":true}
- rule: write as {value: <literal>} or {valueFrom: {kind: KubernetesConfigMap, name: <that resource's name>, fieldPath: status.outputs.configmap_name}} -- a bare string does not parse

### spec.validation.wellKnownCACertificates

`string` · required · optional (explicit presence)

Trust a well-known CA set instead of bringing a bundle. "System" (the
one upstream-defined value) trusts the implementation's system
certificate store — the arm for backends serving publicly-issued
certificates. Implementations may define their own domain-prefixed sets
(e.g. "mycompany.com/my-custom-ca-certificates"). Upstream support:
Implementation-specific.

json_name pins the CRD's exact key: protojson would derive
"wellKnownCaCertificates" from the field name, but the upstream key is
"wellKnownCACertificates" (capital CA) and the API server rejects the
miscased key as undeclared — the faithful-projection contract includes
acronym casing.

- rule: {"string":{"minLen":"1","maxLen":"253","pattern":"^(System|([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]))$"}}

### spec.validation.hostname

`string` · required

Hostname for the backend connection, doing double duty (upstream Core):
it is sent as the SNI (RFC 6066), and — unless subject_alt_names is set
— it is the identity the backend certificate must match. When
subject_alt_names IS set, this hostname only selects the certificate;
add it to subject_alt_names if it should also authenticate.

- rule: {"required":true,"string":{"minLen":"1","maxLen":"253","pattern":"^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.validation.subjectAltNames

`[]KubernetesBackendTlsPolicySubjectAltName`

Subject Alternative Names the backend certificate must contain at least
one of. Set this when the backend's certificate identity differs from
the SNI hostname — the SPIFFE-ID pattern in mTLS meshes is the common
case. Upstream support: Extended.

- rule: {"repeated":{"maxItems":"5"}}
- rule: hostname is required when type is 'Hostname'
- rule: hostname must not be set when type is not 'Hostname'
- rule: uri is required when type is 'URI'
- rule: uri must not be set when type is not 'URI'

### spec.validation.subjectAltNames[].type

`string` · required

Format of this Subject Alternative Name. Closed upstream enum:
"Hostname" | "URI". Always required.

- rule: type must be either 'Hostname' or 'URI'
- rule: {"required":true}

### spec.validation.subjectAltNames[].hostname

`string` · required · optional (explicit presence)

SAN in DNS name format (wildcard first label allowed). Required when
type is "Hostname", forbidden otherwise.

- rule: {"string":{"minLen":"1","maxLen":"253","pattern":"^(\\*\\.)?[a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$"}}

### spec.validation.subjectAltNames[].uri

`string` · required · optional (explicit presence)

SAN in full URI format — scheme and scheme-specific part are both
mandatory (e.g. "spiffe://mycluster.example.com/ns/myns/sa/svc1sa").
Required when type is "URI", forbidden otherwise.

- rule: {"string":{"minLen":"1","maxLen":"253","pattern":"^(([^:/?#]+):)(//([^/?#]*))([^?#]*)(\\?([^#]*))?(#(.*))?"}}

### spec.options

`map<string, string>`

Implementation-specific TLS options for the gateway-to-backend
connection (for example a vendor's minimum-TLS-version knob). Keys
should be domain-prefixed to avoid ambiguity. Upstream support:
Implementation-specific.

Upstream key/value: Gateway API AnnotationKey (1-253) / AnnotationValue
(0-4096).

- rule: {"map":{"maxPairs":"16","keys":{"string":{"minLen":"1","maxLen":"253","pattern":"^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?([A-Za-z0-9][-A-Za-z0-9_.]{0,61})?[A-Za-z0-9]$"}},"values":{"string":{"maxLen":"4096"}}}}

## Outputs

Reference an output from another manifest as `valueFrom: {kind: KubernetesBackendTlsPolicy, name: <resource-name>, fieldPath: status.outputs.<output>}`.

| Output | Type | Description |
|---|---|---|
| `status.outputs.policy_name` | `string` | Name of the created BackendTLSPolicy (equals metadata.name). In InfraCharts this orders the policy after the Services and CA ConfigMaps it references. |
| `status.outputs.namespace` | `string` | Namespace the BackendTLSPolicy was created in (the resolved spec.namespace). The policy can only target Services in this namespace. |

## References

Fields that can point at another resource's outputs:

| Field | Kind | Output |
|---|---|---|
| `spec.namespace` | KubernetesNamespace | `spec.name` |
| `spec.targetRefs[].name` | KubernetesService | `status.outputs.service_name` |
| `spec.validation.caCertificateRefs[].name` | KubernetesConfigMap | `status.outputs.configmap_name` |

## See Also

- [Overview](../README.md)
