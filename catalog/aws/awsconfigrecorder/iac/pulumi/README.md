# AwsConfigRecorder — Pulumi module

Manages the region's Config recording posture:
`aws:cfg/recorder:Recorder`, the folded
`aws:cfg/recorderStatus:RecorderStatus` toggle,
`aws:cfg/deliveryChannel:DeliveryChannel`, and
`aws:cfg/retentionConfiguration:RetentionConfiguration`.

Module facts worth knowing before editing:

- **The names are constants.** Recorder and channel render as
  `default` (AWS's regional singleton convention); the retention
  singleton's name is API-computed. `metadata.name` never reaches AWS.
- **Ordering is encoded.** Channel depends on recorder (AWS refuses a
  channel without one); the status toggle depends on the channel
  (starting without one fails); deletion reverses the chain and
  upstream retries channel deletion while the recorder stop lands.
- **The status resource is the folded toggle.** Unset
  `recording_enabled` sends `isEnabled: true` — the posture the
  component exists for.
- **No tags.** None of the four resources carries tags upstream — the
  catalog's identity-tag map does not apply here.

Outputs mirror the Terraform module key-for-key: `recorder_name`,
`delivery_channel_name`, `recording_enabled`.
