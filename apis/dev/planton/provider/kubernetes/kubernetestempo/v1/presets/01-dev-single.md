# Dev single-node Tempo

One monolithic Tempo replica on a persistent volume with OTLP receivers.
The smallest honest trace store — traces survive pod restarts (unlike the
chart's emptyDir default), but storage is local to the one replica.

**When to use:** local development, demos, a single-node cluster's own
traces.

**When to move on:** for scale or durability beyond one node, switch to
`02-production-s3` (object storage decouples trace blocks from any one
pod).
