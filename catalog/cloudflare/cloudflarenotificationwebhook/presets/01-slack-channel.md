# Slack channel

The chat-mirror shape: a Slack incoming-webhook URL registered as an alert destination, with no secret (the URL itself is the credential Slack issues, which is why it belongs in a managed secret if your policy treats webhook URLs as sensitive). Cloudflare detects "slack" from the URL and formats messages for it -- a generic proxy in front of Slack would get generic payloads instead. Reference the resulting `webhook_id` from any `CloudflareNotificationPolicy`.
