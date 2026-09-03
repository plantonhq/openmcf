# keycloak-login-theme

The Planton design system translated to the self-hosted identity server's pages -- the sign-in form, forced password update, profile completion, error, and logout screens. On a self-hosted install these ARE the product's front door: the console's own `/login` route is only a redirect interstitial, so the identity server's login page is what every teammate on an adopting team sees every working day.

## How Keycloak theming works (primer)

Keycloak renders every user-facing screen server-side from a **theme**: a named directory of static files under `/opt/keycloak/themes/<name>/`, split by page family (`login/`, `account/`, `email/`). A theme is **not code** -- it is:

1. **`theme.properties`** -- the manifest. Its `parent=` line is the core mechanism: a theme extends a built-in one and overrides only what it wants. Everything not overridden (OTP enrollment, recovery codes, WebAuthn, future LDAP-federation screens) keeps the parent's templates and styling.
2. **CSS + images + fonts** under `resources/` -- for a visual reskin, all that is needed.
3. Optionally, FreeMarker page templates -- only to change page *structure*. This theme deliberately overrides **none**: overridden templates rot when Keycloak upgrades, and CSS reached everything.

One realm setting -- `loginTheme` -- selects the theme. The operator sets it in the realm import it already generates.

## Why this theme extends `keycloak.v2`

`parent=keycloak.v2` (Keycloak's bundled PatternFly-v5 login theme) is the safety property: every login-family screen Keycloak can ever serve stays functional and styled, and this theme lays the Planton palette on top. A from-scratch theme (`parent=base`) would leave any screen nobody thought to override as unstyled raw HTML -- on exactly the surface where an unstyled page reads as "something is broken".

Two verified-against-26.3 subtleties encoded in `theme_properties.go`:

- `styles=` entries resolve child-first, falling back to the parent -- so the list re-declares the parent's `css/styles.css` and appends ours.
- The parent ships OS-preference dark/light auto-switching; `darkMode=false` pins the page to the Planton dark design (dark by design, matching the console's default -- not by OS accident).

## Why delivery is a ConfigMap, not a custom image

The operator materializes `Files()` into a ConfigMap mounted at `/opt/keycloak/themes/planton/`. The alternative -- baking the theme into a custom Keycloak image -- was rejected deliberately: the identity server is a security-critical component, and adopters should be able to verify that what runs is the **official, pinned, unmodified upstream image**. A theme tweak is an operator release, never an identity-server rebuild, and `Hash()` lets the operator roll the pod when the theme changes (Keycloak caches themes; a restart is what makes a new version take effect).

## The token-translation contract

Every color here is a **deliberate, documented copy** of the console's design tokens (the platform's theme package). Generating CSS from the TypeScript tokens would mean a cross-language build pipeline for a dozen hex values that change rarely -- disproportionate. The contract instead: **change a token in the console theme, change it here in the same commit.** `theme_test.go` pins the load-bearing values byte-exact so drift fails a test instead of shipping.

One value flows the other way on purpose: `--planton-focus-border` (`#696741`) is the design system's single non-monochrome accent, defined in the console theme as the focused-input border and reused here so "you are typing here" looks identical on both surfaces.

## Layout rule: one named Go file per shipped artifact

Each text artifact is a single constant in a file named for it (`theme_properties.go`, `styles_css.go`, `logo_svg.go`). **The package listing is the manifest of what ships into Keycloak** -- no grab-bag assets directory, no opaque embedded blob. Adding an artifact means adding one named file and one `Files()` entry; a test keeps the two in sync.

The one exception, stated rather than hidden: fonts are binary (woff2) and cannot be honest string constants, so `inter-latin.woff2` lives as a real file beside `fonts.go` (which carries exactly one embed directive per font). Bundling the font keeps the sign-in page fetching **nothing** from third parties -- no font-CDN call from a credential screen, identical rendering air-gapped. Inter is licensed under the [SIL Open Font License 1.1](https://github.com/rsms/inter/blob/master/LICENSE.txt).

## API

- `ThemeName` -- the theme directory / `loginTheme` value (`planton`).
- `Files() map[string][]byte` -- every file, keyed by path under the theme root (`login/...`). Kubernetes-free by design; the operator owns the ConfigMap shape.
- `Hash() string` -- deterministic content fingerprint for restart-on-change.

## Consumers

- the operator's identity component mounts the theme and selects it in the realm import.
