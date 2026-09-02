---
display_name: On-call Paging
---

# On-call paging

The wake-someone-up shape: a tunnel turning unhealthy pages PagerDuty and copies the on-call inbox, with a 30-minute refire floor so a flapping tunnel does not become a pager storm. The PagerDuty service UUID must come from an integration already connected in the Cloudflare dashboard -- Cloudflare accepts an unknown UUID at deploy and fails at delivery. Swap `tunnel_health_event` for `load_balancing_health_alert` or `http_alert_origin_error` to page on those instead.
