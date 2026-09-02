# CloudflareZeroTrustOrganization guide

The judgment this guide protects you from: this resource is a singleton upsert over the account's REAL login infrastructure -- there is no create, no delete, and one field (`auth_domain`) whose change logs every user out of everything.

## The organization exists before you do

Cloudflare creates the organization when Zero Trust is first enabled on the account (the team-name onboarding step). This resource CONFIGURES that organization; it never creates one. Applying against a fresh account without Zero Trust onboarding fails at the API -- enable Zero Trust once, by hand or onboarding automation, then manage everything here.

## auth_domain is the blast radius

The team domain is the URL every user signs in through and every WARP client enrolls against (the full `*.cloudflareaccess.com` hostname -- declare it exactly as the dashboard shows it). Changing it invalidates active sessions and breaks bookmarks, IdP redirect URIs, and WARP enrollment in one apply. Treat it like a production hostname: set it once, correctly, and never touch it casually. And it is REQUIRED on every write (Cloudflare rejects an organization update without it, API error 11004, live-measured): every manifest managing this singleton declares the current team domain even when it never changes.

## Destroy reverts nothing

Deleting this resource is a NO-OP at Cloudflare -- the organization keeps the last-applied configuration forever. There is no "undo by destroy": to revert a setting, apply the previous value. The folded key-rotation cadence behaves identically.

## Unset means unmanaged

Fields you leave unset are never sent -- dashboard-set values survive. This makes partial adoption safe (manage only the login design, say), but it also means REMOVING a field from the manifest does not clear it at Cloudflare; apply the empty/default value explicitly instead.

Two exceptions, both live-measured: `auth_domain` is REQUIRED on every write (Cloudflare rejects the update with error 11004 without it -- declare your current team domain even when it never changes), and `name` should be declared on any account whose organization already carries one (the API echoes the existing name on every read, so an omitted name shows a permanent `name -> null` plan diff).

## MFA and PIV pairing rules

`mfa_required_for_all_apps` needs MFA enabled with at least one authenticator and a session duration configured. And Cloudflare rejects `allowed_authenticators: [ssh_piv_key]` alone when the organization has any non-infrastructure applications -- PIV keys pair only with infrastructure apps. Both rules are API-side; the spec comments carry them so the manifest author sees them first.

## The dashboard lock

`is_ui_read_only: true` makes the API (and therefore IaC) the only write path -- the natural companion of managing the organization from a manifest. Set `ui_read_only_toggle_reason` so dashboard users see why.

## Pairs well with

- [CloudflareZeroTrustAccessIdentityProvider](../cloudflarezerotrustaccessidentityprovider/README.md) -- the sign-in methods behind the login page.
- [CloudflareZeroTrustAccessApplication](../cloudflarezerotrustaccessapplication/README.md) -- the doors this login guards.
- [CloudflareZeroTrustGatewaySettings](../cloudflarezerotrustgatewaysettings/README.md) -- the traffic-filtering half of Zero Trust.
