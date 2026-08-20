# HTTP requests to R2

The traffic-analytics shape: one zone's HTTP request logs into an R2 bucket in the same account, batched every minute, with an explicit field list rather than the dataset's default. Same-account R2 is the destination that skips Cloudflare's ownership handshake, so this shape deploys in one pass. Trim `field_names` to what your queries actually read -- fewer fields is less storage and faster scans.
