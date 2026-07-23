# Kubernetes ServiceAccount: Research Documentation

## Introduction

Every pod in a Kubernetes cluster runs as a ServiceAccount — the namespace's `default` one unless the pod spec says otherwise. That identity is what the kube-apiserver sees when the pod makes API calls, what RBAC grants attach to, what the kubelet consults for registry pull credentials, and — on managed clouds — what federates to a cloud identity so pods can call cloud APIs without stored keys.

The object itself is deceptively thin: a name, some annotations, an `imagePullSecrets` list, an automount flag. The engineering weight is in what hangs off it. A ServiceAccount is the anchor of three distinct concerns:

1. **API authentication** — pods authenticate to the kube-apiserver as `system:serviceaccount:<namespace>:<name>`, and RBAC bindings grant permissions to that subject
2. **Registry authentication** — pull secrets attached to the ServiceAccount are presented by the kubelet for every pod running as it, removing per-pod `imagePullSecrets` boilerplate
3. **Cloud identity federation** — cloud-specific annotations bind the ServiceAccount to a GCP service account, AWS IAM role, or Azure managed identity, making pod-to-cloud access keyless

Planton's **KubernetesServiceAccount** component models all three as typed, validated, reference-aware fields — most notably replacing hand-written cloud annotation strings with a typed per-cloud `workloadIdentity` choice.

## Evolution and Historical Context

### Token secrets era (pre-1.22)

Originally, every ServiceAccount automatically owned a Secret containing a long-lived, never-expiring API token, and the `secrets` list on the object tracked it. These static tokens were a standing credential-theft risk: exfiltrate the Secret once, hold cluster access forever.

### TokenRequest API and projected tokens (1.22–1.24)

Kubernetes moved to short-lived, audience-bound tokens minted on demand by the TokenRequest API and projected into pods. Kubernetes 1.24 stopped auto-creating token secrets entirely. The `secrets` field on ServiceAccount lost its main purpose; its one remaining behavior (mountable-secrets enforcement via an annotation) was deprecated in v1.32. This is why the Planton spec deliberately does not model the `secrets` list — it is legacy surface with no forward story.

### Cloud workload identity (2019+)

The largest shift was external: clusters gained OIDC issuers, and clouds learned to trust them. GKE Workload Identity, EKS IRSA (IAM Roles for Service Accounts), and Azure AD Workload Identity all follow the same shape — the cluster's OIDC issuer vouches for a ServiceAccount, the cloud validates that token against a pre-configured trust, and exchanges it for cloud credentials. No keys are stored anywhere in the cluster.

Each cloud signals the binding through ServiceAccount annotations:

| Cloud | ServiceAccount annotation | Mechanism |
|-------|---------------------------|-----------|
| GKE | `iam.gke.io/gcp-service-account: <email>` | GKE Workload Identity |
| EKS | `eks.amazonaws.com/role-arn: <role-arn>` | IRSA |
| AKS | `azure.workload.identity/client-id: <client-id>` (+ optional `azure.workload.identity/tenant-id`) | Azure AD Workload Identity |

AKS adds one extra requirement: the **pod** must carry the `azure.workload.identity/use: "true"` label, because Azure's mutating webhook only injects the token volume and env vars into labeled pods. GKE and EKS key entirely off the ServiceAccount.

### Security baselines and automount hardening

Security benchmarks (CIS, NSA/CISA hardening guidance) converged on two ServiceAccount rules: dedicated identity per workload (never `default`), and `automountServiceAccountToken: false` for pods that don't call the kube-apiserver — an unused mounted token is pure attack surface. The field is tri-state in the API (unset defers to cluster/pod defaults), which the spec preserves with an `optional bool`.

## Deployment Methods Landscape

### Level 0: Manual (kubectl)

```bash
kubectl create serviceaccount app-identity -n production

# Cloud binding: hand-written annotation
kubectl annotate serviceaccount app-identity -n production \
  eks.amazonaws.com/role-arn=arn:aws:iam::123456789012:role/app-role

# Pull secret attachment: JSON patch
kubectl patch serviceaccount app-identity -n production \
  -p '{"imagePullSecrets": [{"name": "registry-cred"}]}'
```

**Pros:**
- Immediate; fine for experiments

**Cons:**
- The interesting parts (annotations, pull secrets) require annotate/patch follow-ups
- Magic annotation strings typed by hand — a one-character typo fails silently (pods just don't get cloud access)
- No reproducibility, no drift detection

**Verdict:** Debugging only.

### Level 1: Declarative YAML Manifests

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: app-identity
  namespace: production
  annotations:
    eks.amazonaws.com/role-arn: arn:aws:iam::123456789012:role/app-role
imagePullSecrets:
  - name: registry-cred
automountServiceAccountToken: false
```

**Pros:**
- Declarative, version-controllable, full surface

**Cons:**
- Cloud bindings are still untyped annotation strings — nothing validates the key or the value shape
- No connection to the cloud-side trust (IAM binding / trust policy / federated credential) that must exist for the annotation to mean anything
- No plan/preview, no state

**Verdict:** The baseline; correctness rests entirely on the author.

### Level 2: Terraform

```hcl
resource "kubernetes_service_account_v1" "app" {
  metadata {
    name      = "app-identity"
    namespace = "production"
    annotations = {
      "eks.amazonaws.com/role-arn" = aws_iam_role.app.arn
    }
  }

  image_pull_secret {
    name = kubernetes_secret_v1.registry.metadata[0].name
  }

  automount_service_account_token = false
}
```

**Pros:**
- Full IaC lifecycle; can reference the IAM role resource directly, closing the typo class

**Cons:**
- The annotation KEY is still a hand-typed magic string
- Cross-provider wiring (IAM trust policy referencing the cluster OIDC issuer and the SA subject string) is assembled by hand with `format()`/interpolation

**Verdict:** Production-grade, but the cloud-federation contract is still stringly-typed.

### Level 3: Pulumi

```go
sa, err := corev1.NewServiceAccount(ctx, "app-identity", &corev1.ServiceAccountArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name:      pulumi.String("app-identity"),
        Namespace: pulumi.String("production"),
        Annotations: pulumi.StringMap{
            "iam.gke.io/gcp-service-account": gsa.Email,
        },
    },
    AutomountServiceAccountToken: pulumi.Bool(false),
})
```

**Pros:**
- Full programming language; outputs flow between resources naturally

**Cons:**
- Same magic-string annotation keys
- Same hand-assembled trust subjects on the cloud side

**Verdict:** Excellent IaC; the federation contract remains the author's problem.

### Other Methods

**Helm:** charts commonly template a ServiceAccount with an `annotations` values passthrough — flexible, unvalidated.

**eksctl (`iamserviceaccount`):** creates the IAM role, trust policy, AND annotated ServiceAccount together — the right instinct (both halves in one operation), but EKS-only and CloudFormation-backed.

## Comparative Analysis

| Aspect | kubectl | YAML | Terraform | Pulumi | Planton |
|--------|---------|------|-----------|--------|---------|
| Cloud binding | Hand annotation | Hand annotation | Value referenced, key hand-typed | Value referenced, key hand-typed | Typed per-cloud arm |
| Identity handle validation | None | None | None | None | Schema (required, typed, ref-aware) |
| RBAC subject string | Assembled by hand | Assembled by hand | `format()` by hand | `Sprintf` by hand | Exported as output |
| Pull secrets as references | No | No | Yes | Yes | Yes (KubernetesSecret refs) |
| Multi-cloud shape | Per-cloud recipes | Per-cloud recipes | Per-cloud recipes | Per-cloud recipes | One oneof, three arms |
| Dual IaC | N/A | N/A | TF only | Pulumi only | Both |

## The Planton Approach

### Typed workload identity via `oneof`

The core design move: cloud federation is a `oneof` with three arms (`gke`, `eks`, `aks`), shared across every Kubernetes component that needs it. Each arm carries exactly the field(s) that cloud requires:

- **`gke.serviceAccountEmail`** → emits `iam.gke.io/gcp-service-account`
- **`eks.roleArn`** → emits `eks.amazonaws.com/role-arn`
- **`aks.clientId`** (+ optional `aks.tenantId`) → emits `azure.workload.identity/client-id` (+ `azure.workload.identity/tenant-id` when set)

This yields mutual exclusion by construction (a ServiceAccount binds to at most one cloud), required-field enforcement per arm, and zero hand-typed annotation keys. Every arm is a `StringValueOrRef`, so the identity handle can be a literal or a reference to the cloud identity resource's output — in an infra chart, the cluster, the cloud identity, and the ServiceAccount deploy in one run with the handle flowing through the graph.

The cloud-side half of the trust (IAM binding, trust policy, federated credential) is owned by the referenced cloud identity resource; the component's docs state each cloud's requirement precisely so neither half is guesswork.

### The RBAC subject as an output

`system:serviceaccount:<namespace>:<name>` appears in AWS trust policy conditions, Azure federated credential subjects, GKE workload-identity members, and RBAC bindings. The component exports it as `rbac_subject` so no downstream resource re-assembles it — the class of bug where a trust policy says `system:serviceaccount:prod:app` while the SA lives in `production` simply cannot ship.

### What is deliberately absent

The upstream `secrets` list is not modeled. Long-lived token secrets are superseded by the TokenRequest API (1.24+), and mountable-secrets enforcement is deprecated since v1.32. Modeling it would invite users into a legacy pattern.

## Implementation Landscape

### Pulumi Module Architecture

The Pulumi module (`iac/pulumi/module/`) follows the standard Planton pattern:

- **`main.go`**: Orchestrates resource creation — Kubernetes provider and the ServiceAccount resource
- **`locals.go`**: Computes merged labels/annotations, including translating the `workloadIdentity` arm into its annotation set, and resolves namespace and pull-secret references
- **`outputs.go`**: Exports `service_account_name`, `namespace`, `rbac_subject`, and `workload_identity_handle`

### Terraform Module Architecture

The Terraform module (`iac/tf/`) mirrors the Pulumi logic:

- **`variables.tf`**: Mirrors `spec.proto` fields as Terraform variables
- **`locals.tf`**: Computes the annotation set from the workload-identity variables and merges labels
- **`main.tf`**: Creates the `kubernetes_service_account_v1` resource
- **`outputs.tf`**: Exports name, namespace, RBAC subject, and identity handle

### Resource Count

This is a lean component — it creates exactly **one Kubernetes resource**: the ServiceAccount itself. The complexity is in the annotation translation and cross-resource references, not in orchestration.

## Production Best Practices

### Identity hygiene

1. **One ServiceAccount per workload**: Grants to the namespace `default` ServiceAccount leak to every pod in the namespace
2. **Set the namespace explicitly**: It participates in the RBAC subject and every cloud federation subject; implicit `default` placement is almost never intended
3. **Treat name+namespace as a contract**: Cloud trust matches on them exactly; a rename breaks federation with no error on the Kubernetes side — pods just stop getting cloud credentials

### Token hardening

1. **Disable automount by default**: `automountServiceAccountToken: false` for any identity whose pods never call the kube-apiserver; individual pods can re-enable
2. **Never mint long-lived tokens**: Use the TokenRequest API (`kubectl create token`) for ad-hoc needs; static token secrets are a standing exfiltration risk

### Workload identity

1. **Express bindings through the typed field**: Hand-written annotations in `spec.annotations` bypass validation and hide the binding from the resource graph; the component treats colliding annotations as an error, not a merge
2. **Configure both halves**: The annotation without the cloud-side trust does nothing; the trust without the annotation does nothing. Charts that create both in one run eliminate the drift window
3. **AKS: label the pods**: The `azure.workload.identity/use: "true"` pod label is required on every pod using the identity — the webhook ignores unlabeled pods

### Registry credentials

1. **Attach pull secrets at the identity, not per pod**: One attachment covers every pod running as the ServiceAccount
2. **Same-namespace rule**: Pull secrets must live in the ServiceAccount's own namespace; plan secret placement with identity placement

## Conclusion

KubernetesServiceAccount turns a thin object with high-stakes annotations into a typed, validated, reference-aware identity resource. The three concerns it anchors — RBAC, registry pulls, cloud federation — each become first-class fields, the fragile magic strings become generated output, and the `system:serviceaccount:...` subject that every cloud trust matches on becomes an exported fact instead of a hand-assembled one.

## References

- [Kubernetes ServiceAccounts Documentation](https://kubernetes.io/docs/concepts/security/service-accounts/)
- [Configure Service Accounts for Pods](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/)
- [ServiceAccount Token Volume Projection](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#serviceaccount-token-volume-projection)
- [GKE Workload Identity](https://cloud.google.com/kubernetes-engine/docs/how-to/workload-identity)
- [EKS IAM Roles for Service Accounts](https://docs.aws.amazon.com/eks/latest/userguide/iam-roles-for-service-accounts.html)
- [Azure AD Workload Identity](https://azure.github.io/azure-workload-identity/docs/)
- [Pulumi Kubernetes ServiceAccount](https://www.pulumi.com/registry/packages/kubernetes/api-docs/core/v1/serviceaccount/)
- [Terraform kubernetes_service_account_v1](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs/resources/service_account_v1)
