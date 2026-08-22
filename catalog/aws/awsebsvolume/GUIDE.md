# AwsEbsVolume — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## The zone is the contract

A volume lives and dies in one availability zone, and EBS refuses cross-AZ attachments. Decide the zone from the compute's placement FIRST — moving later means snapshot + restore + re-attach, not an edit.

## Copies inherit what you cannot override

`copy_from` lands the copy in the SOURCE volume's zone with the source's encryption posture and snapshot lineage — the provider offers no override, so the spec forbids the create-arm fields on a copy. Need a different zone or key? Snapshot, then restore through the snapshot kind's copy arm.

## Multi-attach is a filesystem decision, not a flag

io1/io2 multi-attach lets several instances see the same block device — with NO write coordination. A regular ext4/xfs mount on two instances corrupts data silently. Only enable it for cluster-aware filesystems (GFS2, OCFS2) or applications doing their own fencing.

## In-place elasticity has a cooldown

Size, type, IOPS, and throughput all modify in place, but AWS allows one modification per volume per ~6 hours — batch your dials into one change instead of trickling them.

## final_snapshot is the teardown safety net nobody sees

It snapshots on destroy, but the snapshot is untracked (config-only at AWS, invisible to imports) — budget for it in cleanup sweeps, and name-search snapshots by the volume's tags when reclaiming space.
