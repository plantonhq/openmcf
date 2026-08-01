# Operator preset

The official Keycloak Operator in its default posture: installed into
the `keycloak` namespace with the NAMESPACED watch — the operator
reconciles only Keycloak declarations living beside it, so several
teams can run isolated operator+Keycloak stacks in separate namespaces.
Installing the operator alone deploys NO Keycloak server; declare
`KubernetesKeycloak` resources in this namespace to get servers.

Every resource in the bundle carries upstream's fixed names
(`keycloak-operator` and friends), so exactly one operator install
fits per namespace. There is no version knob by design: the module
pins the release the `KubernetesKeycloak` declaration kind renders its
CR against — a selectable operator version would drift the CRD schema
away from what the declaration renders. Upgrades arrive as module
updates.

Cluster-wide alternative: set `clusterWide: true` to watch ALL
namespaces (run at most one per cluster, and know the upstream
constraint that custom ServiceAccounts on Keycloak pod templates are
refused in that mode).

Change first: nothing for most installs. Air-gapped clusters point
`operatorImage` / `defaultKeycloakImage` at their mirrors, keeping the
pinned tags.

See [01-operator.yaml](./01-operator.yaml) for the manifest.
