# Generic receiver with a shared secret

The own-your-pipeline shape: alerts POSTed to a service you run, with a shared secret your receiver checks so it can reject anything that did not come from Cloudflare. Cloudflare classifies this as a "generic" destination and sends its standard payload. Point `secret` at a managed secret rather than a literal -- it is write-only at Cloudflare, so the manifest is the only place the value is ever readable, and re-applying is how you rotate it.
