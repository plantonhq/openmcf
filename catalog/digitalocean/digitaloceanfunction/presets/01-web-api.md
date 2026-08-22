# Hello-World Function

This preset deploys a DigitalOcean Functions app from DigitalOcean's public Node.js hello-world sample. Both engines create an App Platform app with one functions component -- there is no standalone Functions resource.

Runtime, memory, timeout, entrypoint, and cron schedules are **not** on this spec. They live in `project.yml` inside `sourceDirectory` (here, `packages`). App Platform reads that file at deploy time. Putting those knobs on the spec would silently do nothing.

The sample uses a public git clone URL so the preset applies without a linked GitHub account. Switch to `github` (owner/repo plus optional `deployOnPush`) only when the DigitalOcean account has GitHub connected.

## When to Use

- HTTP functions whose source is a Functions-style repo (`project.yml` + a `packages/` tree)
- Accounts that do not have GitHub linked yet
- Checking that a Functions deploy works before pointing at your own repo

## Key Configuration Choices

- **Function name** (`functionName`) -- the functions component name inside the app, max 32 characters.
- **Git clone URL** (`git`) -- public HTTPS clone. No GitHub OAuth required.
- **Source directory** (`sourceDirectory: packages`) -- directory that contains `project.yml` and the packages tree.
- **Environment variables** (`envs`) -- `plaintext` or `secret`. Secrets are stored in App Platform's secret store.

To change runtime or add a schedule, edit `project.yml` in the repo, not this manifest.
