# Partner SFTP User with Key Authentication

This preset onboards one exchange partner with SSH public-key
authentication and a full-access grant on the partner's own container
-- the per-partner isolation pattern that lets one account serve many
partners.

## When to Use

- B2B file exchange where each partner lands files in its own container
- Automated pipelines (MFT tools, cron jobs) that authenticate with a
  key pair

## Key Configuration Choices

- **Key auth only** -- prefer this over passwords: the partner rotates
  its own key pair, and nothing secret transits between you
- **The scope is the isolation boundary** -- the user sees and touches
  ONLY the granted container; `homeDirectory` drops sessions straight
  into it
- **The login is `{account}.{user}`** (the `sftp_username` output), on
  `{account}.blob.core.windows.net` port 22
- The ACCOUNT needs `sftpEnabled: true` + `isHnsEnabled: true` for
  logins to work

## Placeholders to Replace

| Placeholder | Description | Where to Find |
| --- | --- | --- |
| `<storage-account-resource-name>` | The SFTP AzureStorageAccount's Planton resource name | Your exchange composition |
| `partneracme` | 3-64 lowercase letters and digits | Your partner naming convention |
| `ssh-ed25519 AAAA...` | The partner's OpenSSH public key ("ssh-ed25519 AAAA...") | Provided by the partner |
| `<containerName>` / `<container-resource-name>` | The partner's container | Your exchange composition |
