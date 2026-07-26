# Team SSO + RBAC preset

Argo CD the way a team should meet it: logins through the
organization's identity provider, authorization through Argo CD's own
RBAC (read-only by default, admin for the platform group), and the
local admin user retired once SSO carries the load. The client secret
never touches the manifest — Argo CD resolves `$<secret>:<key>`
references at runtime from a labeled Secret composed next to it, so
the credential lives where credentials belong.

The order of operations matters and this preset encodes it: `domain`
tells SSO redirects and CLI hints the public name, `server.insecure`
hands TLS to the composed edge (an Ingress or Gateway kind referencing
the exported `server_service`), and `admin_enabled: false` only after
the OIDC path works — disabling admin before SSO functions locks
everyone out of the UI.

Change first: tighten `rbac.policy_csv` from one admin group to real
per-project roles as teams onboard, and consider `exec_enabled` only
with a matching RBAC `exec` rule — an interactive shell through the
API server is a deliberate grant, never a default.

See [02-team-sso-rbac.yaml](./02-team-sso-rbac.yaml) for the manifest.
