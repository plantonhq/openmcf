# PAC files only

The proxy-configuration shape: PAC files managed as rows without touching the settings or logging singletons at all. Each row is a real resource -- removing one deletes the file. The `slug` is baked into the file's public download URL and forces replacement on change: set it deliberately and never touch it. Replace the proxy endpoint in `contents` with your account's own Gateway proxy hostname.
