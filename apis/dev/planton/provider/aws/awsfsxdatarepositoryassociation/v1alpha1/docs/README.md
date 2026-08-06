# AwsFsxDataRepositoryAssociation — Technical Reference

## Resource Mapping

One `aws_fsx_data_repository_association` resource. The provider models the
sync policies inside a single-purpose `s3 { auto_import_policy { events } /
auto_export_policy { events } }` wrapper; the spec exposes the two event lists
directly — the wrapper carries no information of its own.

## Identity and Lifecycle

An association's identity is the (file system, file-system path, S3 URI)
triple: `file_system_id`, `file_system_path`, and `data_repository_path` are
all ForceNew. Everything else updates in place:

- `auto_import_events` / `auto_export_events` — sync policies can be added,
  changed, or removed on a live association.
- `imported_file_chunk_size` — stripe size for newly imported files.

`batch_import_meta_data_on_create` and `delete_data_in_filesystem` are
create-/delete-time behaviors, not persistent AWS state.

## HSM Semantics

Data repository associations use Lustre's Hierarchical Storage Management:

1. **Import** brings object *metadata* into the namespace; file data stays in
   S3 until first read (lazy hydration adds latency to the first access).
2. **Auto-import** (`NEW`/`CHANGED`/`DELETED`) keeps the namespace tracking
   the bucket; without it the namespace reflects only creation-time state
   (`batch_import_meta_data_on_create`) plus manual import tasks.
3. **Auto-export** writes file-system changes back to the bucket as they
   happen; without it, results return to S3 via manual export tasks
   (`lfs hsm_archive` or FSx data repository tasks addressed by
   `association_id`).

## Constraints

- Up to **8 associations per file system**, **25 per account** (AWS quotas).
- File-system paths must not overlap between associations on the same file
  system — each directory subtree belongs to at most one repository.
- A file system carrying the legacy in-spec `import_path` link cannot also
  carry associations (the two S3-link generations are mutually exclusive).
- PERSISTENT_2 file systems support associations exclusively; SCRATCH and
  PERSISTENT_1 support both generations (prefer associations).

## Composition

- `file_system_id` references `AwsFsxLustreFileSystem.status.outputs.file_system_id`.
- `association_id` is the join key for FSx data repository tasks; the ARN is
  the IAM policy target.

## Appendix: Quick Reference

| Spec Field | Default | ForceNew |
|------------|---------|----------|
| region | (required) | Yes |
| file_system_id | (required) | Yes |
| file_system_path | (required) | Yes |
| data_repository_path | (required) | Yes |
| auto_import_events | (none) | No |
| auto_export_events | (none) | No |
| imported_file_chunk_size | 1024 | No |
| batch_import_meta_data_on_create | false | Create-time |
| delete_data_in_filesystem | false | Delete-time |
