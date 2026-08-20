# Bearer-authenticated third-party server

The vendor shape: Cloudflare holds the vendor's bearer token (a managed-secret reference -- the value is write-only at Cloudflare and never returned) and presents it on every upstream call, so users never touch the credential. Upstream traffic additionally runs through Secure Web Gateway policies. Rotate by changing the managed secret and re-applying -- dashboard rotation is invisible to IaC.
