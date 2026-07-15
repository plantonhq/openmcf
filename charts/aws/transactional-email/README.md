# AWS Transactional Email

Production email sending on SES, set up the way deliverability
consultants charge for: a DKIM-signed domain identity, a custom MAIL
FROM domain so SPF aligns for DMARC, a configuration set enforcing TLS
and suppression, a bounce/complaint feedback topic your application
subscribes to, and reputation alarms that page BEFORE your account hits
the thresholds SES suspends senders over.

Password resets, receipts, magic links — the mail your product cannot
function without — lands in inboxes instead of spam folders because of
exactly the records and feedback loops this chart wires.

## Architecture

```
 application ──▶ SES SendEmail (from anything@email_domain)
                   │
        AwsSesEmailIdentity (domain, Easy DKIM 2048)
                   │ inherits                    envelope (mail_from_enabled)
        AwsSesConfigurationSet ◀─────────────── MAIL FROM mail.email_domain
                   │ TLS posture, suppression,
                   │ reputation metrics
     ┌─────────────┴─────────────┐
     ▼ (feedback_notifications)  ▼ (alarms_enabled)
 AwsSnsTopic feedback         AwsCloudwatchAlarm bounce-rate ≥ 5%
   bounce/complaint events    AwsCloudwatchAlarm complaint-rate ≥ 0.1%
   (your app subscribes)         │
                              AwsSnsTopic alerts ──▶ email

 AwsRoute53DnsRecord (dns_records_enabled, bring-your-own zone):
   _dmarc.email_domain TXT        — DMARC policy + aggregate reports
   mail.email_domain   MX         — SES regional feedback endpoint (mail_from)
   mail.email_domain   TXT        — SPF for the envelope domain     (mail_from)
 DKIM CNAMEs ×3 — post-deploy from the identity's dkim_tokens output
```

## Included Cloud Resources

| Resource | Kind | Purpose |
|----------|------|---------|
| Configuration set | `AwsSesConfigurationSet` | TLS/suppression posture, reputation metrics, event publishing |
| Domain identity | `AwsSesEmailIdentity` | The verified sender: Easy DKIM, config-set inheritance, MAIL FROM |
| Feedback topic | `AwsSnsTopic` | Machine channel: every bounce/complaint as JSON (conditional) |
| DMARC record | `AwsRoute53DnsRecord` | The domain's authentication policy + reporting (conditional) |
| MAIL FROM MX + SPF | `AwsRoute53DnsRecord` ×2 | Envelope-domain verification and alignment (conditional) |
| Alerts topic + subscription | `AwsSnsTopic` + `AwsSnsSubscription` | Human channel for the alarms (conditional) |
| Reputation alarms | `AwsCloudwatchAlarm` ×2 | Bounce ≥ 5%, complaint ≥ 0.1% — AWS's attention thresholds (conditional) |

## Parameters

| Name | Description | Default | Required |
|------|-------------|---------|----------|
| `aws_region` | SES region (MAIL FROM MX points at its feedback endpoint) | `us-east-1` | yes |
| `aws_account_id` | Account id scoping the feedback topic's SES publish grant | `123456789012` | yes |
| `sender_name` | Name prefix for every resource | `mailer` | yes |
| `email_domain` | The domain to verify and send from — immutable | `example.com` | yes |
| `mail_from_enabled` | Custom envelope domain for DMARC-aligned SPF | `true` | no |
| `mail_from_subdomain` | Label for the MAIL FROM domain | `mail` | when MAIL FROM on |
| `dns_records_enabled` | Publish DMARC/MX/SPF into an existing Route 53 zone | `false` | no |
| `zone_id` | The hosted zone that owns `email_domain` | placeholder | when DNS on |
| `dmarc_policy` | `none` → `quarantine` → `reject` ratchet | `none` | when DNS on |
| `dmarc_report_email` | DMARC aggregate report (rua) destination | `dmarc-reports@example.com` | when DNS on |
| `tls_required` | Refuse plaintext delivery | `true` | no |
| `feedback_notifications_enabled` | Bounce/complaint events onto the feedback topic | `true` | no |
| `alarms_enabled` | Reputation alarms wired to email | `true` | no |
| `alert_email` | Alert destination (confirm AWS's first email) | `ops@example.com` | when alarms on |

## Completing verification (post-deploy)

The identity deploys in `PENDING` verification and SES flips it to
`SUCCESS` when it sees the DKIM records in DNS. The three CNAME values
are generated WITH the identity, which is why no template can publish
them — read them from the identity's outputs and publish each:

```bash
# The identity's dkim_tokens output carries three tokens; for each one:
#   name : <token>._domainkey.<email_domain>
#   type : CNAME
#   value: <token>.dkim.amazonses.com
```

With Route 53, that is three `AwsRoute53DnsRecord` resources (or three
console entries) — one per token, TTL 3600. Verification typically
completes within minutes of the records resolving; `verification_status`
in the identity's outputs tracks it.

If `dns_records_enabled` is off, also publish the chart's render-stable
records manually: the `_dmarc` TXT, and (with MAIL FROM) the MX
(`10 feedback-smtp.<region>.amazonses.com`) and SPF TXT
(`v=spf1 include:amazonses.com ~all`) on the MAIL FROM subdomain.

## First send

New accounts start in the SES sandbox: sending is limited to verified
addresses until AWS approves production access (SES console → Account
dashboard → Request production access; approval usually lands within a
day). Then:

```bash
aws sesv2 send-email --region us-east-1 \
  --from-email-address noreply@example.com \
  --destination ToAddresses=you@example.com \
  --content 'Simple={Subject={Data="hello"},Body={Text={Data="from SES"}}}'
```

The identity's default configuration set applies automatically — TLS
posture, suppression, metrics, and feedback events all inherited.

## The feedback loop

With `feedback_notifications_enabled`, every bounce and complaint lands
on the feedback topic as JSON. Subscribe your application (an SQS queue
is the durable shape — create the queue, grant `sns.amazonaws.com`
`sqs:SendMessage` scoped to the topic in the queue's policy, and add an
`AwsSnsSubscription`) and act on the events: hard bounce → delete the
address; complaint → unsubscribe and never mail again. The account-level
suppression list catches repeats, but the feedback loop is what keeps
your OWN lists clean.

## Day-2 guidance

- **Ratchet DMARC.** Run at `p=none` for a few weeks and read the
  aggregate reports at `dmarc_report_email`. When every legitimate
  source aligns, move to `quarantine`, then `reject` — the point where
  spoofing your domain stops working.
- **Enforce MAIL FROM.** Once the MAIL FROM domain shows verified, set
  `behaviorOnMxFailure: REJECT_MESSAGE` on the identity's `mailFrom`
  block so a future DNS regression fails loudly instead of silently
  de-aligning your mail.
- **Tighten the feedback topic grant.** The topic policy scopes SES
  publishing by configuration-set ARN pattern; after deploy you can pin
  it to the exact ARN from the set's outputs.
- **Open/click tracking.** Add `trackingOptions` with a custom redirect
  domain (a CNAME to the regional `awstrack.me` endpoint) to keep
  tracking links on your domain, then add `OPEN`/`CLICK` to the event
  destination's types.
- **Multi-region sending.** Deploy this chart once per sending region —
  identities and configuration sets are regional; the DNS records
  differ only in the MX target's region.
- **Dedicated IPs.** High-volume senders (many millions/month) can add
  a dedicated IP pool and set `deliveryOptions.sendingPoolName` on the
  configuration set; below that volume, SES's shared pools with VDM
  deliver better than a cold dedicated IP.
