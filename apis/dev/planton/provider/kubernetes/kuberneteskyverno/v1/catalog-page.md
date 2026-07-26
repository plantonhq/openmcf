# Kyverno

Policy as Kubernetes resources. Kyverno validates, mutates, generates,
and cleans up cluster resources with policies written as plain
Kubernetes YAML — no new policy language — making it the fastest path
from "we need guardrails" to enforced guardrails. The CNCF-graduated
engine behind Pod Security replacement, image verification, and
auto-generation of defaults across thousands of production clusters.

## Highlights

- **Policies are YAML, not a language** — platform engineers write and
  review policy the same way they write everything else on the
  cluster; agents configure it from the spec alone.
- **Four typed controllers** — admission, background, cleanup, and
  reports, each independently sized and scheduled; HPA on the
  admission path; ServiceMonitor fan-out across all four.
- **Safety-first lifecycle** — the pre-delete hook that removes the
  runtime-registered webhooks stays on by default, the CRD-cascade
  destroy warning lives on the field that controls it, and
  fail-open-everywhere is one typed flag when you want it.
- **Air-gap in one field** — a single registry override reroutes every
  Kyverno image, including the CRD-migration hook's non-obvious
  reg.kyverno.io home.
- **Native VAP generation** — offloads eligible validations to the API
  server's own ValidatingAdmissionPolicy, cutting webhook round-trips.

---

© Planton. Licensed under [Apache-2.0](https://github.com/plantonhq/planton/blob/main/LICENSE).
