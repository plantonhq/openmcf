---
title: Reading a Service's Repository — Two Lanes, One Judgment
description: How to read a service's repository the way an engineer troubleshoots — the platform-mediated file read for one or a few files (list_repo_files, read_repo_files, planton service repo ls|cat), the sandbox clone at an exact commit for depth (a one-hour token used once, never stored), which commit to read (the run's sha, never today's HEAD, when explaining a failure), validating a pipeline straight from the repository, discovering a repository before it is registered, and every refusal's next step. Read when a failed build or deploy needs the actual files explained, when someone asks what a repository contains or whether its pipeline compiles, when you need to run or search a repository rather than glance at it, or when a repository read is refused.
---

# Reading a Service's Repository — Two Lanes, One Judgment

A service's truth is its repository, and the platform gives you two ways in. **The file lane**: the platform reads files for you through the organization's GitHub connection — the tree at a commit, then up to twenty files' text in one call — with no credential leaving the control plane. **The clone lane**: you ask the platform for the connection's one-hour installation token and shallow-clone the repository at an exact commit into your workspace, where you can run, search, and diff. Choose the way an engineer chooses: start from the run's own words and open the one or two files it names; clone only when the answer needs a working copy (running the install, grepping a tree, diffing two commits). A glance never needs a clone; a reproduction never fits in a glance.

**Which commit, always.** A failed run built ONE commit: `spec.git_commit.sha` on the run record (`get_service_pipeline`; CLI `planton service pipelines`/`runs`). Read the files AT THAT SHA — `ref` on every read below — and explain the failure from what the build actually saw. Today's default branch may already differ; reading it to explain yesterday's failure is guessing with extra steps. "Would this work now?" and "what does the repository contain?" read HEAD (empty `ref`).

## The file lane

Address the repository one of two ways, and never both: **by registered service** (`org` + `service` slug — the record supplies the connection, owner, repository, and default branch), or **by coordinates** for a repository no service binds yet (`org` + `git_connection` + `owner_name` + `repo_name`; see "Before a service exists").

- **See what is there**: `list_repo_files` (CLI: `planton service repo ls <service> [--ref] [--prefix]`) answers every file under an optional prefix with its size, plus `commit_sha` — the commit actually read — and `truncated`. Use a prefix when you know the neighbourhood (`.planton/` for the pipeline and its tasks, `_kustomize/` for the deploy tree); sizes tell you what fits a read and what wants a clone.
- **Read the files**: `read_repo_files` (CLI: `planton service repo cat <service> <path> [--ref]`) takes one to twenty repository-relative paths and answers each on its own: `found` with `content`; found but not carried (`unreadable_reason` says why — over one megabyte, or not text — and names the clone lane); or `found: false`. One absent or oversized file never fails the batch, so ask for the whole neighbourhood at once: the Dockerfile with the lockfile it installs, the pipeline with its tasks, the overlay with its base.
- **Validate the pipeline straight from the repository**: `validate_service_pipeline` with the repository as its source (`org` + `service`, or the coordinates, plus an optional `ref`) reads `.planton/pipeline.yaml` (or `pipeline_path`) and the tasks beside it and compiles them exactly as dispatch would; the report's `pin` is the commit it read. `references/service.pipeline-authoring.md` owns the compile loop; this is its "nothing pasted" entrance.

What the platform will never do in this lane: execute anything it reads, write anything back (writing has its own door — a pull request, `references/service.opening-a-pull-request.md`), or hand out a credential. What it needs from the caller: read access on the connection (the same grant detection uses), or credential-read on it — either passes.

## The clone lane

When the files are not the answer — the install must run, the tree must be searched, two commits must be diffed — clone at the exact commit.

**Platform-tools arm** (the web console; `git` in your sandbox, no CLI):

1. Mint the token: `get_github_installation_token` with `org` + `slug` (the slug is the service's `spec.gitRepo.gitConnection`). Refused? You lack credential-read on that connection — say so in one sentence and continue in the file lane; never ask around it.
2. Clone with the token supplied ONCE, on the fetch's own command line, into a directory inside your workspace:

```sh
git init --quiet storefront && cd storefront
git remote add origin https://github.com/<owner>/<repo>.git
git -c http.extraheader="AUTHORIZATION: basic $(printf 'x-access-token:%s' "$TOKEN" | base64)" fetch --quiet --depth 1 origin <sha-or-branch>
git checkout --quiet --detach FETCH_HEAD
```

   `-c` applies to that one invocation, so nothing lands in `.git/config` or the remote URL; the checkout carries no credential (a later `git fetch` will ask for one — mint again if you need it). Fetching by sha is how you get exactly the commit a run built.
3. Hygiene, always: the token appears in this conversation's tool output and expires within the hour; it reaches every repository the installation reaches. Never write it to a file, never repeat it in prose, never put it in a remote URL, never commit anything. The sandbox is disposable; treat the clone as one.

**CLI arm** (desktop): `planton service repo clone <service> [directory] [--ref <sha-or-branch>]` does the whole sequence with the same hygiene — token minted, used once on the fetch, never stored; shallow, detached at the commit. Then `planton service pipeline validate` on the checkout, `grep`, `docker build`, whatever the question needs.

The clone lands in your workspace and nowhere else — the hard boundary on reading outside the workspace still holds.

## Before a service exists

Registering a repository (`references/service.detection-first-registration.md`) and reading one the platform has never seen use the coordinates door. Find them in two reads: `search_api_resources_by_kind` on kind `github_connection` in the organization gives each connection's `slug`; `list_github_repositories` with `org` + `slug` lists every repository that connection's installation reaches, each with `owner_name`, `name`, and `default_branch` — exactly the coordinates `detect_service_from_repo`, `list_repo_files`, and `read_repo_files` take. A repository absent from that list is one the installation was not granted; the organization's GitHub administrator adds it on the installation's settings.

## Walking a refusal

Every refusal names its next step; relay it, never retry around it.

- **"you are not allowed to read repositories through this connection"** (file lane): the caller holds neither read nor credential-read on the connection. Name the connection and stop; access is granted on the connection, not the service.
- **Token refused** (clone lane): credential-read on the connection is the stricter grant. Continue in the file lane and say what you could not do.
- **`truncated: true`**: the repository is too large for the provider to list whole. Narrow with a prefix, or clone.
- **`unreadable_reason`**: the file is over one megabyte or not text. Clone to read it; say which.
- **`found: false`**: the path does not exist at that commit. List the directory before guessing another spelling.
- **"is not a branch, tag, or commit of …"**: the ref is wrong for this repository; the run record's `spec.git_commit.sha` is the authoritative one for a failure.
- **"deploys from a GitLab repository"**: the platform's file lane reads GitHub connections only today; the CLI on the person's own checkout, or a clone with their own credential, is the way in.
- **A provider sentence relayed verbatim** (rate limit, revoked installation, repository moved): GitHub's own words are the diagnosis; relay them and name the administrator when the installation is the cause.

## Writing back

Neither lane writes. When a read ends in a fix, the fix leaves through one door only: a pull request the platform opens on a branch of its own, after the person has seen the exact files and the pull request's text and said yes — `references/service.opening-a-pull-request.md`. On the desktop, edit the clone in your workspace and hand the edited files to `planton service repo pr`; never push from the clone with the minted token. A person who would rather commit themselves gets the exact file or diff, explained against the lines you read.
