# GCP Firestore Backup Schedule — Deep Dive

## The problem this resource solves

Point-in-time recovery (PITR) covers the last 7 days of document versions — enough for recent mistakes, not for compliance archives or recovery after database deletion. Firestore's managed backup schedules extend protection to 14 weeks and survive database deletion. Declaring schedules as infrastructure makes a deployment's backup posture reviewable and reproducible instead of console-only.

## Daily-plus-weekly is two resources

Firestore allows one daily and one weekly schedule per database. The common production pattern — short-retention daily backups beside longer-retention weekly archives — is exactly two `GcpFirestoreBackupSchedule` resources on the same database, not one resource with two cadences.

## Mutability profile

| Surface | Mutability |
|---|---|
| `database`, recurrence (`daily` / `weeklyRecurrence`) | Immutable |
| `retention` | Mutable — applies to backups created after the change |

## Backups outlive the schedule

Deleting this resource stops future backups but never deletes existing ones. They age out per the retention they were created with. This is why teardown of a database with active backups may need those backups to expire first.

## Timing is chosen by Firestore

Daily schedules run every day but the exact time within the day is chosen by Firestore. Weekly schedules pin a day; timing within that day is still Firestore-managed.

## No labels surface

Firestore backup schedules do not support GCP labels — both engines skip label merge identically.

## IAM and API prerequisites

- `firestore.googleapis.com` enabled (both modules enable it with `disable_on_destroy=false`).
- Firestore Admin permissions on the project.

## Deliberately not modeled

- **`deletion_policy`** — client-side Terraform lever conflicting with Planton-managed destroy (catalog-wide decision).
