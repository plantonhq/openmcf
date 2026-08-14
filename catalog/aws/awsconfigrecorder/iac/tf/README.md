# AwsConfigRecorder — Terraform/OpenTofu module

Manages the region's Config recording posture:
`aws_config_configuration_recorder`, the folded
`aws_config_configuration_recorder_status` toggle,
`aws_config_delivery_channel`, and
`aws_config_retention_configuration`.

Module facts worth knowing before editing:

- **The names are constants.** Recorder and channel render as
  `default` (AWS's regional singleton convention); the retention
  singleton's name is API-computed. `metadata.name` never reaches AWS.
- **Ordering is encoded.** Channel depends on recorder (AWS refuses a
  channel without one); the status toggle depends on the channel
  (starting without one fails); deletion reverses the chain and the
  provider retries channel deletion while the recorder stop lands.
- **The status resource is the folded toggle.** Unset
  `recording_enabled` renders `is_enabled = true` — the posture the
  component exists for.
- **No tags.** None of the four resources carries tags upstream — the
  catalog's identity-tag map does not apply here.

Outputs mirror the Pulumi module key-for-key: `recorder_name`,
`delivery_channel_name`, `recording_enabled`.
