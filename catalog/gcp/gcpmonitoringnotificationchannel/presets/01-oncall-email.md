# On-call Email

The simplest paging path: incidents from every referencing alert policy
land in one inbox. No external credentials, no webhook plumbing.

## What it configures

- `type: email` with the destination address in
  `channelLabels.email_address` — the only configuration key the email
  type takes.

## Adjust before deploying

- **email_address** — replace with the real on-call address (a shared
  alias or rotation-backed address beats an individual inbox).
- Email channels require one-time verification before they deliver —
  check the `verification_status` output after the first deploy and
  complete the verification from the address's inbox.

## When to choose something else

Email is a slow paging medium. For real on-call response, pair this
channel with a **PagerDuty Service** or **Slack Channel** preset on the
same alert policies — policies accept multiple channels.
