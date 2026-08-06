# Air-gapped mirror preset

Kyverno for clusters that cannot reach public registries. One field
does the heavy lifting: `imageRegistry` sets the chart's global
registry, rerouting EVERY Kyverno image — the four controllers, the
pre-delete webhook-cleanup hook, and the CRD migration hook (whose
default home is reg.kyverno.io, a registry air-gap allowlists often
miss) — to your mirror. Repository paths and tags stay chart-managed,
so your mirror only needs to host the same paths under its own
hostname. `imagePullSecrets` names an existing pull secret in the
install namespace.

The second mirror concern is what Kyverno does to OTHER workloads'
images: `config.defaultRegistry` is where the engine's default image
mutation points bare image references (`nginx` becomes
`mirror.example.com/nginx`). Aligning it with your mirror makes the
policy engine itself enforce the air-gap for every pod it admits.

Change first: your mirror's actual pull-secret name, and add
`resourceFiltersExclude` entries only after the engine is stable —
policing more of the cluster from day one widens the blast radius of
a mirror outage.

See [03-airgapped-mirror.yaml](./03-airgapped-mirror.yaml) for the
manifest.
