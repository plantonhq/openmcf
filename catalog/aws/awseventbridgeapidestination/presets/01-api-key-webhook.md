# API-Key Webhook

Both arms in one instance: an api-key connection (the key lands in an AWS-owned Secrets Manager secret, never in state) and a POST destination with a static source header and a 50/s rate cap. Point an EventBridge rule, pipe, or schedule at the destination's ARN output.

Every credential-bearing value in this preset is a `$secret/<slug>` org-secret reference — the backend rejects plaintext in sensitive fields at create. Create the referenced org secrets (the API key, and the source-header value: AWS masks every connection http parameter on read, so the spec secret-types them all), or pick your own in the wizard.
