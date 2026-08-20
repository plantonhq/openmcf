# Planton on Kubernetes

Run Planton itself on your own cluster. One deploy installs the Planton
operator (once per cluster) and declares a complete self-hosted platform
— control plane, web console, identity server, databases, secrets
manager, and an in-cluster deployment runner — in its own namespace.
Zero-config by design: the version is the one required choice, the
built-in gateway serves console and sign-in over a single port-forward
(the exact command is a deploy output), and the first console visitor
becomes the admin using a setup code read from a Secret (also an output).

| Resource | Kind | Purpose | Conditional on |
|---|---|---|---|
| `<env>-planton-ns` | KubernetesNamespace | The platform's home — owned once, joined by the platform | always |
| `<env>-planton-operator` | KubernetesPlantonOperator | The manager (one per cluster, self-enforced) in `planton-operator` | `install_operator` |
| `<env>-planton` | KubernetesPlantonPlatform | The platform declaration the operator reconciles | always |

## The shape

- **One operator, many platforms.** The operator watches every namespace;
  each platform lives in its own namespace with its own URL, identity
  server, and databases. To run a SECOND platform on the same cluster,
  deploy this chart again with a different `namespace` and
  `install_operator: false` — the resident operator picks it up. (Two
  honest cluster-wide facts: builds can feed only one build-enabled
  platform per cluster — Tekton's events sink is a singleton — and the
  cluster runs one operator/CRD schema version, while each platform pins
  its own `version`.)
- **Exposure is opt-in refinement.** The default door is the built-in
  gateway over `kubectl port-forward`. `ingress_enabled` +
  `hostname` (+ `cert_cluster_issuer` for HTTPS via cert-manager) serve a
  real URL — set the hostname BEFORE the first sign-in; the identity
  server bakes the URL into its realm at first boot. Deeper exposure
  shapes (ALB/ACM annotations, brought certificates) live on the platform
  resource's own spec after deploy.
- **Storage is one dial.** `storage_class` + `storage_size` lift every
  platform volume at once — built for backends with minimum-size floors.
  Unset, the cluster's default class serves, and the operator verifies it
  can actually provision before deploying.

## After deploy

The platform resource's outputs carry the two commands you need:
`port_forward_command` opens the door on your machine, and
`setup_code_command` reads the code the console's first-visit setup page
asks for. With ingress enabled, `kubectl get plantonplatforms -A` shows
the URL the moment it is admitted.

## Version and upgrades

`version` pins the platform's image line and is never defaulted by the
underlying kind — the chart's param default is the line this chart
release was validated against. Upgrading a running platform is editing
the value and redeploying; the operator rolls the platform to the new
line.
