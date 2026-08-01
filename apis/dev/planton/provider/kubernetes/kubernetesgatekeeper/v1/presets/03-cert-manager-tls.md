# cert-manager TLS preset

Gatekeeper serving its webhook with a certificate issued by
cert-manager instead of the engine's embedded rotator. Organizations
standardizing certificate issuance (one CA, one audit trail, one
renewal pipeline) point `externalCert.secretName` at the Secret a
KubernetesCertificate materializes — the reference above composes the
two kinds so the dependency edge is explicit in the resource graph.

Order matters: the Certificate must be issued BEFORE this install
(the chart mounts the Secret; a missing one holds the rollout), and
its DNS names must cover
`gatekeeper-webhook-service.gatekeeper-system.svc` (plus the
fully-qualified `...svc.cluster.local` form). The Secret name in the
manifest is a NAME reference only — no certificate material ever
renders into chart values.

What the module handles for you (chart-truth at the pin): enabling
external injection auto-disables cert rotation on the audit
controller but NOT on the controller-manager — the module sets
`disableCertRotation` there explicitly, or the embedded rotator would
keep overwriting your issued certificate.

Change first: pair this with the production-enforce posture; if you
run cert-manager you are past the audit-only stage.

See [03-cert-manager-tls.yaml](./03-cert-manager-tls.yaml) for the
manifest.
