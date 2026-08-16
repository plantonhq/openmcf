# AwsManagedPrefixList — Operational Guide

Live-earned judgment lands here as proof runs and adopter operations teach it; the notes below are the forge-time seed.

## max_entries is a quota decision, not a guess

A security-group rule referencing this list consumes `max_entries` slots of the rules-per-group quota (default 60) — a 20-capacity list in one rule spends a third of the group's quota even holding two entries. Size capacity to real growth, not round numbers.

## The family is forever

`address_family` replaces the list on change, which breaks every rule and route referencing the old `pl-` id — there is no in-place IPv4→IPv6 story. Run dual-stack as two lists.

## Description edits look expensive because they are

AWS has no "update entry description" call — the provider removes and re-adds the entry across two API round trips. Expected plan noise, not drift; entries keep their CIDR identity throughout.

## Concurrent writers serialize on the version

Every modification carries the list's current version, and AWS rejects stale writers (`PrefixListVersionMismatch`). Both engines' single-owner in-line model avoids the race entirely — resist adding out-of-band entries with the standalone entry resource, which the provider itself serializes behind a mutex for this reason.

## Deletion drains

A deleted list drains through `delete-in-progress` — and AWS refuses deletion while any rule or route still references the `pl-` id. Deleting this component after its consumers is the clean order (charts get this right by dependency).
