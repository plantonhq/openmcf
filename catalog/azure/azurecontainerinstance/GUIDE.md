# Azure Container Instance -- Operational Guide

Judgment calls that matter when you run container groups in production.

## Design the group's shape before shipping -- almost nothing updates in place

Azure applies exactly two changes to a live group: identity and tags. Everything else -- images, resources, ports, volumes, probes, networking -- replaces the group. Three fields are worse than ForceNew: `cpuLimit`, `memoryLimit`, and `keyVaultUserAssignedIdentityId` are ACCEPTED on update and silently never applied (the provider's update path skips them), so a manifest change there looks successful and does nothing. Treat all three as create-only; when they must change, change something ForceNew alongside (or recreate deliberately).

## This is a group, not a pod with a scheduler

There is no rescheduling, no rolling update, no horizontal scaling. A liveness probe restarts a container IN PLACE on the same host; if the host dies, an "Always" group restarts on another, but a "Never" job does not rerun. For anything needing orchestration -- scaling, rollouts, self-healing across hosts -- reach for Container Apps or AKS. Reach for THIS kind when the unit of work IS the group: jobs, burst workers, simple services, sidecar pairs.

## The secrets never come back

Azure returns none of the group's secrets on reads: secure environment variables, volume storage keys, inline secret files, registry passwords, and the Log Analytics workspace key are all write-only. Two operational consequences: an IMPORT of an existing group cannot recover them (re-supply them in the manifest -- that is expected, not drift), and any identity update re-sends them from your configuration, so keep the manifest's secret references live (a rotated storage key must be updated in the manifest before the next apply, or the update call fails).

## Networking: pick one of three postures and mean it

**Public** (the default) gives an internet-facing IP; add `dnsNameLabel` for a stable name and set a reuse policy stricter than the default "Unsecure" if anything long-lived points at the label -- a released label with "Unsecure" is claimable by anyone (dangling-DNS takeover). **Private** joins a subnet delegated to `Microsoft.ContainerInstance/containerGroups`; Azure serializes group operations per subnet, so parallel deploys into one subnet queue up. **None** means no group IP at all -- and it is the only posture Spot priority accepts. A DNS label on a "None" group would be silently discarded by Azure; the spec rejects the combination outright.

## Volumes: four forms, one choice each

`azureFile` is the only persistent form -- and it needs the storage account KEY (managed-identity mounts are not in the provider's surface), so treat that manifest as secret-bearing. `emptyDir` is group-lifetime scratch; giving the SAME empty_dir name to several containers is the sharing mechanism (init seeds, worker reads). `gitRepo` clones at start -- pin a `revision` or every replacement group gets whatever the branch head is that day. `secret` mounts inline files; values are BASE64 of the file content, and a plain-text value fails at mount time, not at validation.

## Jobs: Never + exit codes, and the read-back quirks

For run-once work set `restartPolicy: Never`, let the container exit, and read the exit state before deleting the group. Two read-back behaviors are expected, not drift: `commands` echoes the image's own entrypoint when you did not override it, and the group-level exposed ports echo the containers' ports when you omitted `exposedPorts`.
