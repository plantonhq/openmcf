# Web Service with Autoscaling and Zero-Downtime Rollouts

This preset is the production baseline for a stateless web service: CPU-based horizontal autoscaling between 2 and 10 replicas, a rollout strategy that never drops below the desired replica count, a pod disruption budget that survives node drains, and a pre-stop sleep that drains connections before termination.

## When to Use

- Production web services and APIs with variable traffic
- Any service where a dropped request during deploys or node maintenance is unacceptable

## Key Configuration Choices

- **Autoscaling 2→10 on 70% CPU** — `replicas: 2` is the floor the autoscaler never goes below; scale-out triggers when average CPU across replicas exceeds 70% of requests
- **`maxUnavailable: "0"` + `maxSurge: "1"`** — rolling updates create the new pod before removing an old one; combined with the readiness probe, deploys are zero-downtime
- **PDB `minAvailable: "1"`** — voluntary disruptions (drains, upgrades) always leave at least one serving pod
- **Pre-stop sleep of 10s** — the kubelet-native sleep hook keeps the pod serving while endpoint removal propagates through load balancers; no sleep binary needed in the image

## Placeholders to Replace

| Placeholder | Description | Where to Find |
|---|---|---|
| `<your-namespace>` | Target namespace for the deployment | Your namespace management or `KubernetesNamespace` resource |
| `<your-container-registry>/<your-image>` | Container image repository | Your container registry |
| `<your-image-tag>` | Image tag or version | Your CI/CD pipeline output |

## Related Presets

- **01-web-service** — Minimal single-replica starting point
- **04-hardened-production** — Adds restricted-profile security hardening and topology spread
