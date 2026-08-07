# Composite Filter-and-Sort Index

The standard multi-field index: equality filter on one field, sort on
another — the shape Firestore error-message links usually suggest.

## When to use

Any query combining a where-clause on one field with orderBy on another
within the same collection.

## What to customize

- `fields` — list fields in query order (equality filters first, then
  inequality/sort fields).
- `collection` — the collection ID your queries target.

## Composes with

`GcpFirestoreDatabase` upstream (reference its `database_name` output).
