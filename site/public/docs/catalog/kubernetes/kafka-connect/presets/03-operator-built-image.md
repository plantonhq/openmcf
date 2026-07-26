---
title: "Operator-built image preset"
description: "The build arm end-to-end: the operator builds a Connect image containing the Debezium Postgres connector (pinned by Maven coordinates, resolved at build time) with Kaniko, pushes it to your registry..."
type: "preset"
rank: "03"
presetSlug: "03-operator-built-image"
componentSlug: "kafka-connect"
componentTitle: "Kafka Connect"
provider: "kubernetes"
icon: "package"
order: 3
---

# Operator-built image preset

The build arm end-to-end: the operator builds a Connect image
containing the Debezium Postgres connector (pinned by Maven
coordinates, resolved at build time) with Kaniko, pushes it to your
registry using the named push-credentials Secret, and runs the
workers on the result. There is deliberately NO `image` field — the
spec enforces the image-XOR-build rule because the operator would
silently override a declared image with the one it builds.

Choose this arm when you need an exact, reproducible artifact set —
several connectors at pinned versions in one image — without
maintaining your own image pipeline. URL artifacts (`jar`, `tgz`,
`zip`) work too; give them a `sha512sum` so a tampered download fails
the build instead of running in the workers.

What the build needs from the environment: a registry the build pod
can push to, and the `registry-push-creds` docker-registry Secret in
this namespace (the field carries the Secret's NAME — no credential
ever rides the manifest). The built image reference is then a natural
`image` value for other Connect clusters that want the same plugins
without rebuilding.

See [03-operator-built-image.yaml](./03-operator-built-image.yaml)
for the manifest.
