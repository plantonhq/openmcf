# CloudflareZeroTrustGatewaySettings guide

The judgment this guide protects you from: one component, three lifecycles -- and two provider defects the manifests must be designed around, not discovered by.

## Three lifecycles under one spec

The `settings` and `logging` surfaces are account SINGLETONS: create and update are the same PUT, and destroy is a NO-OP that abandons the live configuration. The `pac_files` rows are real resources: removing a row deletes the file. Plan teardown accordingly -- destroying this component removes PAC files but leaves every Gateway setting exactly as last applied.

## Unset means unmanaged (settings only)

A settings sub-object you don't declare is never sent -- dashboard values survive. This makes partial adoption safe, but removing a sub-object from the manifest does NOT revert it: apply the previous value explicitly. The `logging` surface is the deliberate opposite: when declared, the COMPLETE tree ships (unset switches become false, Cloudflare's default), because Cloudflare reports drift on partially-sent logging that would never converge.

## The 2211 trap: certificate before decryption

`tls_decrypt`, `fips.tls`, and deep `body_scanning` all require an ACTIVATED Gateway certificate on the account. Without one, the API rejects the write with error 2211 -- and Cloudflare's own test tooling documents accounts left erroring in that state. The certificate lifecycle is not yet a catalog kind: activate one out-of-band before flipping these switches, and flip them in a change window.

## The block-page drift fields

The provider has a recorded defect: `block_page.mode`, `include_context`, `suppress_footer`, and `target_uri` drop from state on refresh when absent from configuration -- its own tests accept permanent non-empty plans on them. If you manage the block page, DECLARE those four fields explicitly (even at their defaults) so refresh has nothing to drop.

## PAC slugs are URLs

A PAC file's `slug` is baked into its public download URL and forces replacement on change -- every device configured with the old URL breaks. Set the slug deliberately on day one (or accept the name-derived one) and never touch it.

## Pairs well with

- [CloudflareZeroTrustGatewayPolicy](../cloudflarezerotrustgatewaypolicy/README.md) -- the rules this panel sets the behavior for.
- [CloudflareZeroTrustDnsLocation](../cloudflarezerotrustdnslocation/README.md) -- the entry points for DNS filtering.
- [CloudflareZeroTrustOrganization](../cloudflarezerotrustorganization/README.md) -- the login half of Zero Trust.
