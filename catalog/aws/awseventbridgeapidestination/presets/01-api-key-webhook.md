# API-Key Webhook

Both arms in one instance: an api-key connection (the key lands in an AWS-owned Secrets Manager secret, never in state) and a POST destination with a static source header and a 50/s rate cap. Point an EventBridge rule, pipe, or schedule at the destination's ARN output.
