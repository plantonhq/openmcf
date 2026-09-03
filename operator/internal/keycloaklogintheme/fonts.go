package keycloaklogintheme

import _ "embed"

// The one deliberate exception to this package's one-constant-per-file
// rule: fonts are binary (woff2) and cannot be honest Go string constants,
// so they live as real files beside this one -- still individually named
// and countable in the package listing.
//
// interLatinWOFF2 is the Inter variable font (latin subset, weights
// 100-900 in one file), the console's UI typeface. Bundling it keeps the
// sign-in page fetching NOTHING from third parties -- no font CDN call
// from a credential screen, and air-gapped installs render identically.
// Inter is licensed under the SIL Open Font License 1.1 (see README).
//
//go:embed inter-latin.woff2
var interLatinWOFF2 []byte
