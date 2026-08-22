# DigitalOcean Database Kafka Schema -- Operational Guide

What experience with this component teaches that the field reference cannot.

## Every change drops ALL prior versions -- the kind's loudest warning

The provider has no update path: changing the definition (or type, or subject name) destroys the subject and re-registers the new document as version 1. Consumers that pin older schema versions lose them the moment the replacement lands. If your consumers rely on registry-mediated compatibility across versions, do NOT evolve schemas through this resource -- evolve them through your producers' registry client (which appends versions), and use this resource only to declare the founding schema.

## Whitespace is a replacement

The definition is compared verbatim -- no JSON normalization, no key reordering, nothing. A reformatted document is a "change" that triggers the destroy-and-drop above. Keep the manifest's schema string byte-stable: single-line, machine-formatted, never hand-prettified after the fact.

## Compatibility level is not here

The registry's subject compatibility level (BACKWARD, FULL, etc.) has no surface in the provider at the pinned version -- DigitalOcean's API carries it, the provider does not. What this resource registers is the document itself; compatibility policy stays whatever the registry defaults to (or whatever was set out-of-band).

## Imports do not work

The upstream importer is broken at the pinned provider version (it never restores the subject name, so the post-import read addresses nothing). Adopting an existing subject means recreating it through this resource -- which, per the first warning, drops its version history. Plan adoption accordingly.

## What is deliberately NOT here

Multi-version subject management and compatibility levels (no provider surface -- see above); topics (their own kind, DigitalOceanDatabaseKafkaTopic); and registry credentials (the cluster's users).
