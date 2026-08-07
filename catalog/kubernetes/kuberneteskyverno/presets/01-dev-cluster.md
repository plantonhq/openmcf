# Dev cluster preset

The smallest declarable Kyverno: the chart defaults end to end. All
four controllers run single-replica, the engine generates and rotates
its own webhook certificates (no prerequisites), the policy CRDs
install with the release, and the webhooks exclude kube-system and the
kyverno namespace itself. This is a complete, working policy engine —
apply a ClusterPolicy afterwards and it enforces immediately.

Know the two lifecycle facts before you rely on it. First, the
webhooks Kyverno registers at runtime default each policy rule to
fail-CLOSED (`failurePolicy: Fail`): if the admission controller goes
down, resources matched by your policies stop admitting until it
returns. On a dev cluster that trade is usually right — you want to
SEE enforcement — but it is why the production preset runs three
replicas. Second, destroy deletes the policy CRDs, which
cascade-deletes every policy and policy report on the cluster; set
`crds.keepOnUninstall` if your policies must survive an engine
reinstall.

Change first: `features.forceFailurePolicyIgnore` if you would rather
the engine fail-open while you learn (enforcement gaps instead of
blocked admissions), and `admissionController.replicas: 3` the moment
this cluster hosts anything you care about.

See [01-dev-cluster.yaml](./01-dev-cluster.yaml) for the manifest.
