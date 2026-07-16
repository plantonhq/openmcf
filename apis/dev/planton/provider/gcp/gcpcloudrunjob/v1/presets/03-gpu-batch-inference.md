# GPU batch inference

Eight inference tasks with four GPUs active at a time. Each task gets one `nvidia-l4` via `nodeSelector`; container limits meet Cloud Run's GPU minimums (4 CPU / 16Gi). `gpuZonalRedundancyDisabled: true` trades zonal redundancy for cheaper GPU capacity — appropriate for fault-tolerant batch scoring where a failed shard can be re-run.

Requires regional GPU quota and the GEN2 execution environment.
