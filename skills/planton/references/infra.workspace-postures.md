---
title: The Three Workspace Postures — What This Folder IS
description: The full protocol behind the folder-identity check — composing in your own workspace (charts as top-level subfolders, checkouts, loose manifests and manifest SETS, the canvas rules), working a deployed project's working copy, and composing when the folder itself is the chart. Read when the identity check's one-line answers need their complete choreography — where files go, how checkouts land, when a set beats a chart, and how the canvas follows your writes.
---

# The Three Workspace Postures — What This Folder IS

The identity check (`ls .planton/ 2>/dev/null` — its files never appear in the workspace tree) gives three answers. This is each answer's complete choreography.

## `.planton/workspace.yaml` — this is YOUR WORKSPACE

A plain working folder, like a colleague's directory: you fill it with whatever the request calls for. Every chart you compose is its own TOP-LEVEL subfolder, named for the chart (`gke-cluster/Chart.yaml`, `gke-cluster/values.yaml`, `gke-cluster/templates/`), and one request may produce several charts side by side. Loose resource manifests and notes may live at the workspace root. Never place chart files at the workspace root — the root is the surface that HOLDS charts, not a chart. Build each chart from its own folder: `planton chart build <chart-dir>`. The files are the user's: when they ask, copy anything to a destination they name (`cp`/`mv` — a path the user gives you is an invitation under the workspace boundary).

**What already exists is checked out, never re-typed.** When the work starts from a published chart or a deployed project, pull the real files into their own top-level subfolder and work there: `planton chart checkout <slug> --output-dir <chart-dir>` lays out a chart from the catalog ready to customize, and `planton infra project checkout <id-or-name> --output-dir <project-dir>` lays out a deployed project's working copy (its hidden binding included, so the folder follows `references/infra.deployed-projects.md`). A checkout writes many files at once, so the composing declaration comes FIRST, exactly as for files you author yourself. And it lands where your SHELL is: `cd` to the folder you were given before running it, so the subfolder is born inside the workspace, never beside it.

**A chart is for a parameterized architecture; loose manifests cover the rest.** When the request names one thing ("an S3 bucket for our assets"), write ONE manifest at the workspace root, check it with `planton validate <file>`, and offer the next step — apply it now (`planton apply -f <file>`, a mutation with the usual one confirmation) or grow it when the request grows. Several resources wired together deploy WITHOUT chart ceremony too: a folder of manifests (or one multi-document file) applies as a SET — `planton apply -f <dir>` runs a preflight report first (schema, cross-references, state backend, credentials, modules — everything verifiable, one report, before any IaC handoff), then deploys in dependency order derived from the manifests' own `valueFrom` references, each resource's captured outputs resolving the next ones' references. Exit codes tell CI the truth: 2 means refused at preflight (nothing ran), 1 means a deploy failed (state advanced; re-running the same command re-applies completed resources as no-ops and continues). Reach for a chart when the architecture wants parameterized reuse or per-environment templating. On the platform-tools arm there is no offline validator for a loose manifest: ground it meticulously from its reference page, say plainly that validation happens at apply, and apply through the platform's own apply tool (the same mutation, the same one confirmation).

**The canvas follows your writes.** When several charts exist and the user asks to look at or change "the chart" without naming one, confirm which they mean before writing — writing into the wrong chart drags their view there. When you move your own work between charts, say so in one line first.

## `.planton/project.yaml` — the working copy of a DEPLOYED project

The folder is bound to a running deployment: your edits target THAT project, saving starts a real deployment pipeline, and the workflow bends — read `references/infra.deployed-projects.md` before doing anything.

## Neither exists — the folder itself is the chart

A chart pulled from git or with `planton chart checkout`, a scaffold, a folder the user picked: compose in place at its root, exactly as the chart anatomy shows.
