# Minimal — evaluation registry, zero dependencies

The smallest honest Harbor: the chart's in-cluster PostgreSQL and Redis
(single-node, evaluation-grade by upstream's own position), artifact
blobs on a 20Gi PersistentVolumeClaim, Trivy scanning on (the chart
default), and a ClusterIP front door reached through the exported
`port_forward_command`.

What you get without declaring it: every credential is generated per
install — the admin password lands in the exported `<name>-admin-auth`
Secret (key `HARBOR_ADMIN_PASSWORD`), and none of the chart's publicly
documented defaults (`Harbor12345`, `changeit`, `not-a-secure-key`)
ever ship.

KNOW THIS before pushing images: `externalUrl` is load-bearing —
Harbor embeds it in the auth-token URL returned to every OCI client,
so `docker login/push/pull` must dial exactly that address. This
preset pins it to the port-forward address; change it in lockstep with
whatever exposure you compose in front.

Graduation path: swap `database.internal` for the external arm
composing a KubernetesPostgres, `cache.internal` for a
KubernetesValkey, and `storage.filesystem` for an object-storage
backend — see the production preset.
