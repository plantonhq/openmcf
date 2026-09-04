# Every Auth0 and OpenFGA kind gets a glyph of its own, and both providers join the judged-provider pin

## What changed

- **All six Auth0 logos and all three OpenFGA logos are judged under the logo law.** Before, the six Auth0 kinds shared one file -- and that file was Auth0's own badge trademark, the one thing the law forbids copying -- and the three OpenFGA kinds shared a hand-drawn imitation of the OpenFGA mark's dots and dumbbells on its dark tile. Nine kinds shared a glyph and nine named no provenance. Now no two kinds share bytes, every file opens with its provenance, and no file copies or imitates anyone's mark.
- **Both providers are drawn in full, because none of their kinds IS the provider.** A client, a connection, a resource server, a role, an action, and an event stream are objects inside an Auth0 tenant; a store, an authorization model, and a relationship tuple are objects inside an OpenFGA server. The law gives an official mark only to a kind that IS a product or built-in object of its provider where that provider offers icons for the purpose; Auth0's brand assets sit behind Okta's partner-credentialed brand portal and Auth0 publishes no icon set for its objects; OpenFGA's marks live in `cncf/artwork` under the CNCF trademark policy (attribution, linking, and no-recolor conditions the law declines to carry onto every console page) and it publishes no icon set for its objects either. Forge rule `021-logo`'s provider list records both findings.
- **The glyphs, from the shared vocabulary, in each provider's own palette.** Auth0, in its brand orange `#eb5424` with `#f27b53` as the second tone (a paler `#f5936c` faded on the light plinth): an application that signs users in is an app window with a person inside, a connection users arrive through an open door with an arrow entering, a protected API a server block with a shield standing at its side, a role an identity badge card listing the permissions it bundles, a login action a bolt between code brackets, an event stream events travelling down a pipe out to a destination. OpenFGA, in a deep green `#1fa971` with `#5fd49b` as the second tone (the hue of its mark's neon `#8bff95`, which itself vanishes on a light plinth): a store a chest holding the relationships written into it, an authorization model a schema of object types joined by their relations, a relationship tuple a user related to an object. Every glyph judged on the contact sheet at 18, 26, 34, and 48 px on both washes; every one distinct from its siblings at 18 px. The vocabulary in rule 021 gains the nine concepts so the next identity or authorization provider draws from the same list.
- **Auth0 and OpenFGA join the judged-provider pin.** `judgedProviders` is `gcp, cloudflare, kubernetes, auth0, openfga`; the 18 baseline entries are gone (656 → 638) and nothing else in the baseline moved. Copying one glyph onto a sibling fails the gate and the pin at both kinds by name.

## Why

On a diagram the icon is the one identity signal nothing else supplies, and a customer's sign-in and authorization pictures said it with one trademark repeated six times and one imitation repeated three. The provider's own products would wear its marks where offered; these kinds are not the provider, so they say what they are, in the provider's color, and the catalog carries no one else's mark or obligations.

## How to check

```bash
md5 -q catalog/auth0/*/logo.svg catalog/openfga/*/logo.svg | sort | uniq -d   # prints nothing
grep -L '<desc>' catalog/auth0/*/logo.svg catalog/openfga/*/logo.svg          # prints nothing
go test ./pkg/cataloglogo/ ./pkg/anatomy/                                      # green, both providers absent from the baseline and pinned
go run ./tools/catalog-logo-sheet/ -provider auth0                             # then open the HTML it names; the same for openfga
```

## What comes next

- The logos publish to the CDN at their versionless keys with the next release tag; consoles pick them up after the cache turns over (four hours at most), with no platform change.
- AWS, Azure, and DigitalOcean are judged the same way as their catalog work reaches them.
