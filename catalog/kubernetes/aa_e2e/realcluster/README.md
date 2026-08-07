# Real-Cluster E2E Batches

Some Kubernetes E2E assertions need cloud fabric no local kind cluster can
provide: cloud load balancers, IRSA identity hops, snapshot-capable CSI
storage, real node autoscaling. Those scenarios carry a real-cluster value in
the `planton.dev/e2e-cluster-profile` annotation (e.g. `aws-eks`) — they skip
with a reason on local runs and execute against a batch-provisioned real
cluster through the harness's external-cluster lane.

## Lane contract

```bash
# 1. Provision the batch cluster + IRSA roles + batch test resources:
./aws-eks/bootstrap.sh

# 2. Source the generated env file (lane selection + identifier exports):
source ~/.planton-e2e/planton-e2e-eks/env.sh

# 3. Run targeted per-component lanes (NEVER tier-wide sweeps):
go test -tags=e2e -timeout=30m -v -count=1 -run 'Test.*KubernetesVelero_' ./e2e/

# 4. Tear down and verify zero residue:
./aws-eks/teardown.sh && ./aws-eks/audit.sh
```

The env file exports two kinds of variables:

- **Lane selection** — `PLANTON_E2E_KUBECONFIG` (adopts the external cluster)
  and `PLANTON_E2E_CLUSTER_PROFILE` (declares what the cluster IS, so
  profile-annotated scenarios are matched rather than assumed; a `cilium-cni`
  scenario can never run onto an EKS batch cluster by accident).
- **Batch identifiers** — `PLANTON_E2E_*` values (IRSA role ARNs, bucket,
  queue, zone id) that scenario manifests reference through
  `${E2E_ENV:PLANTON_E2E_*}` tokens. Committed manifests stay honest: they
  never hardcode one test account's identifiers.

## Why targeted invocations

Component-level create/verify/destroy is identical to the kind lanes; only
the cluster outlives the run. Run one component's entrypoints at a time:
capacity managers (Karpenter, ClusterAutoscaler) must never overlap, and a
tier-wide sweep wastes the batch on components already proven on kind.

## aws-eks batch

| Asset | Purpose |
|---|---|
| `cluster.eksctl.yaml` | Cluster shape: system + tainted `ca-scale` (min-0) node groups, OIDC, EBS CSI + snapshot-controller addons, Karpenter node-role identity mapping |
| `karpenter-prerequisites.cloudformation.yaml` | Node role, controller policies, interruption SQS queue — vendored from the pinned upstream Karpenter getting-started template |
| `bootstrap.sh` | Stack + cluster + discovery tags + IRSA roles + zone/secret/bucket/queue + storage classes + env file |
| `storage.yaml` | gp3 StorageClass + VolumeSnapshotClass (CSI snapshot lanes) |
| `teardown.sh` | Deletes everything, handling the non-obvious residues (zone records, bucket objects, tagged snapshots) |
| `audit.sh` | Enumerates every resource class by tag/name; fails on any survivor |

IRSA trust policies bind to each chart's **rendered** ServiceAccount name
(documented in the modules' `vars.go`/`locals.go`): `kube-system/karpenter`,
`cluster-autoscaler/cluster-autoscaler-aws-cluster-autoscaler`,
`velero/velero-server`, `keda/keda-operator`, plus the scenario-named
external-dns SA and the ESO store fixture's ServiceAccount. Renaming a
scenario that owns one of these means updating `bootstrap.sh` to match.
