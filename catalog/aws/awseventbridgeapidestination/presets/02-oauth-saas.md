# OAuth SaaS Integration

Client-credentials OAuth against a SaaS token endpoint: EventBridge fetches the token itself (grant_type and scope ride the token request's body — most OAuth servers require at least grant_type) and invokes the activity endpoint under it. The client secret is secret-typed and lands in the AWS-owned secret.

Every credential-bearing value in this preset is a `$secret/<slug>` org-secret reference — the backend rejects plaintext in sensitive fields at create. That includes the token-request body values (`grant_type`, `scope`): AWS masks every connection http parameter on read, so the spec secret-types them all. Create the referenced org secrets (grant_type's secret holds the literal `client_credentials`, scope's holds your scope string), or pick your own in the wizard.
