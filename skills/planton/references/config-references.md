# Variables and Secrets — the `$var` / `$secret` Reference Grammar

Planton manifests never carry credentials, and rarely carry org-specific
config literals. Instead, any string field can hold a REFERENCE to a value
managed in the platform's config manager: a variable (plain config) or a
secret (encrypted, resolved just-in-time inside the deployment runner, never
stored in the manifest or readable back out of it). This file is the grammar,
when to use it, and how to ground references against what actually exists.

This is a DIFFERENT instrument from `valueFrom` (`dependencies.md`):

| Value comes from | Wire it with |
|---|---|
| Another resource's deployment output (an id, ARN, endpoint the platform creates) | `valueFrom` — see `dependencies.md` |
| An operator-managed config value (a region-independent setting, a team-owned constant) | `$var/...` |
| A credential or any sensitive value (password, API key, token, private key) | `$secret/...` — the ONLY thing a sensitive field accepts |

## The grammar

Both kinds share one shape. Scope is encoded in the reference itself:

```
$var/<slug>                    # organization-scoped variable
$var/@<env>/<slug>             # environment-scoped variable
$var/<group>/<entry>           # entry in a variable group (org scope)
$var/@<env>/<group>/<entry>    # entry in a variable group (env scope)

$secret/<slug>                 # organization-scoped secret
$secret/@<env>/<slug>          # environment-scoped secret
$secret/<name>/<key>           # key inside a key-value secret (org scope)
$secret/@<env>/<name>/<key>    # key inside a key-value secret (env scope)
```

The `@` sigil as the FIRST path segment is the scope switch, and it is
strict: an org-scoped reference never resolves an environment-scoped record
and vice versa — there is no fallback. A secret that exists only in the
`staging` environment is `$secret/@staging/db-password`; writing
`$secret/db-password` for it fails validation with "secret not found",
exactly like a misspelled slug. Scope is part of the address.

Example — a sensitive field on a database user, environment-scoped:

```yaml
spec:
  userName: app
  password: $secret/@staging/db-password
```

## Sensitive fields accept ONLY a secret reference

Every field a component's schema marks sensitive (the explain report and the
component reference page flag these) rejects plaintext before anything
deploys — the control plane validates that the field holds a well-formed
`$secret/...` reference to an EXISTING secret. So for a password/key/token
field there are exactly two failure modes to avoid:

- a plaintext value (rejected outright — never write one, not even a
  placeholder), and
- a reference whose slug or scope does not match a real secret (rejected as
  not found).

Never expose a chart param that asks the user to paste a credential — that
is the same failure as a plaintext literal, one step removed. The param, if
any, carries the REFERENCE (see "In chart templates" below).

## Ground before you write — never invent a slug

References are validated against the org's real records, so look up what
exists before writing one, exactly as you ground field names with
`planton explain`:

```
planton secret list -o json      # every secret; each record's "env" field
planton variable list -o json    #   distinguishes org- from env-scoped
```

Use `-o json`: the JSON records carry each entry's `env` (empty = org
scope), which the human-readable table does not show — and the scope decides
which reference form you write. On the platform-tools arm (no CLI), check
your roster for a config-manager search/list tool; when none exists, you
cannot verify existence — treat every secret you reference as
possibly-missing and follow the missing-secret protocol below. Never fake a
lookup or claim a secret exists unverified.

## Referencing a secret that does not exist yet

Normal, not an error: the reference is the design; the secret's value is
deployment-time material only its owner should supply. When the lookup finds
no matching secret (or no lookup instrument exists), still write the
reference — never plaintext, never a placeholder — and then say so plainly
in the explain-after:

1. name every reference you wrote whose secret does not exist yet,
2. say the secrets must be created before deploy, and
3. hand the user the exact ready-to-run commands, offering to run them
   yourself once they supply the values (a mutation — one confirmation), or
   point at the console's Secrets page as the click path:

```
planton secret set db-password value=<the-value> --env staging   # env-scoped
planton secret set api-key value=<the-value>                     # org-scoped
```

A composed chart with declared-but-uncreated secrets is honest and
deployable after one step; a chart with an invented value is a silent
failure.

## Creating variables and secrets (mutations — confirm first)

```
planton secret set <slug> value=<value> [--env <env>]     # key=value pairs
planton secret set <slug> --from-file value=./key.json    # file-backed value
planton variable set <slug> <value> [--env <env>]         # value is POSITIONAL
```

Note the asymmetry: `secret set` takes `key=value` pairs (a secret can be a
key-value map); `variable set` takes the value as a positional argument.
`--env` makes the record environment-scoped; without it the record is
org-scoped. Re-running `secret set` on an existing secret writes a new
version (rotation), never a duplicate.

## In chart templates

Slugs are DNS-compatible (`[a-z0-9-]`) and can never contain `@`, so the env
sigil composes cleanly with Jinja — the canonical pattern for a chart that
deploys per-environment:

```yaml
spec:
  password: $secret/@{{ values.env }}/db-password
```

Prefer carrying the WHOLE reference through a param over hardcoding a slug,
so the chart stays portable across orgs whose secrets are named differently
(param defaults are plain values — template expressions live in templates/,
never in values.yaml):

```yaml
# values.yaml
params:
  - name: db_password_secret_ref
    description: Secret reference for the database password (e.g. $secret/@staging/db-password)
    value: $secret/db-password

# templates/database.yaml
spec:
  password: "{{ values.db_password_secret_ref }}"
```

Either way the manifest field ends up holding a `$secret/...` string — the
param is a naming convenience, never a channel for the raw value.
