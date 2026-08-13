# GcpKmsKey Guide

The judgment this guide protects: everything encrypted under this key
is exactly as durable as the key's versions. Most resources can be
recreated; destroyed key versions cannot, and neither can the data
under them.

## Set PREVENT before the key matters

`deletionPolicy: DELETE` (the provider default) makes destroy schedule
EVERY key version for destruction and disable rotation — after the
`destroyScheduledDuration` window (default 30 days) elapses, data
encrypted under this key is permanently unrecoverable. The key SHELL
survives (KMS keys have no delete API), which makes the damage easy to
misread in the console: the key is still listed while its versions die.
For any key protecting data you cannot afford to lose, set
`deletionPolicy: PREVENT` the day the key is created — the destroy then
fails instead of starting a 30-day countdown someone must notice.
`ABANDON` is the clean exit: the key leaves management with every
version intact and data stays decryptable.

## The recovery window is your last undo

`destroyScheduledDuration` is the gap between "destroy version" and
"gone forever" — versions can be restored inside it. Shortening it
below the default trades blast-radius recovery time for compliance
optics; lengthen it for keys where a mistaken destroy must survive a
long vacation.

## Names are permanent; plan for keys, not key edits

A key name can never be reused within its ring, even after every
version is destroyed. Combined with immutable purpose and protection
level, the practical pattern is: new requirement → new key (and
re-encrypt), not "adjust the existing key". Rotation
(`rotationPeriod`, ENCRYPT_DECRYPT only) mints new primary versions on
cadence — old versions stay decryptable, so rotation is cheap; version
destruction is the expensive act.

## Protection levels are a procurement decision

SOFTWARE covers most workloads; HSM buys FIPS 140-2 Level 3 attestation
at higher cost; EXTERNAL/EXTERNAL_VPC hand the key material to an
external manager — after which GCP availability depends on the EKM
being reachable. Choose by compliance requirement, not by "HSM sounds
safer".
